// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Package main contains users main function to start the users service.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	chclient "github.com/mainflux/callhome/pkg/client"
	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/internal"
	jaegerclient "github.com/mainflux/mainflux/internal/clients/jaeger"
	pgclient "github.com/mainflux/mainflux/internal/clients/postgres"
	redisclient "github.com/mainflux/mainflux/internal/clients/redis"
	"github.com/mainflux/mainflux/internal/email"
	"github.com/mainflux/mainflux/internal/env"
	mfgroups "github.com/mainflux/mainflux/internal/groups"
	gapi "github.com/mainflux/mainflux/internal/groups/api"
	gcache "github.com/mainflux/mainflux/internal/groups/redis"
	gtracing "github.com/mainflux/mainflux/internal/groups/tracing"
	"github.com/mainflux/mainflux/internal/postgres"
	"github.com/mainflux/mainflux/internal/server"
	httpserver "github.com/mainflux/mainflux/internal/server/http"
	mflog "github.com/mainflux/mainflux/logger"
	mfclients "github.com/mainflux/mainflux/pkg/clients"
	"github.com/mainflux/mainflux/pkg/groups"
	gpostgres "github.com/mainflux/mainflux/pkg/groups/postgres"
	"github.com/mainflux/mainflux/pkg/uuid"
	capi "github.com/mainflux/mainflux/users/api"
	"github.com/mainflux/mainflux/users/clients"
	"github.com/mainflux/mainflux/users/clients/emailer"
	uclients "github.com/mainflux/mainflux/users/clients/postgres"
	ucache "github.com/mainflux/mainflux/users/clients/redis"
	ctracing "github.com/mainflux/mainflux/users/clients/tracing"
	"github.com/mainflux/mainflux/users/hasher"
	"github.com/mainflux/mainflux/users/jwt"
	clientspg "github.com/mainflux/mainflux/users/postgres"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
)

const (
	svcName        = "users"
	envPrefixDB    = "MF_USERS_DB_"
	envPrefixES    = "MF_USERS_ES_"
	envPrefixHTTP  = "MF_USERS_HTTP_"
	envPrefixGrpc  = "MF_USERS_GRPC_"
	defDB          = "users"
	defSvcHTTPPort = "9002"
	defSvcGRPCPort = "9192"
)

type config struct {
	LogLevel        string `env:"MF_USERS_LOG_LEVEL"              envDefault:"info"`
	SecretKey       string `env:"MF_USERS_SECRET_KEY"             envDefault:"secret"`
	AdminEmail      string `env:"MF_USERS_ADMIN_EMAIL"            envDefault:""`
	AdminPassword   string `env:"MF_USERS_ADMIN_PASSWORD"         envDefault:""`
	PassRegexText   string `env:"MF_USERS_PASS_REGEX"             envDefault:"^.{8,}$"`
	AccessDuration  string `env:"MF_USERS_ACCESS_TOKEN_DURATION"  envDefault:"15m"`
	RefreshDuration string `env:"MF_USERS_REFRESH_TOKEN_DURATION" envDefault:"24h"`
	ResetURL        string `env:"MF_TOKEN_RESET_ENDPOINT"         envDefault:"/reset-request"`
	JaegerURL       string `env:"MF_JAEGER_URL"                   envDefault:"http://jaeger:14268/api/traces"`
	SendTelemetry   bool   `env:"MF_SEND_TELEMETRY"               envDefault:"true"`
	InstanceID      string `env:"MF_USERS_INSTANCE_ID"            envDefault:""`
	PassRegex       *regexp.Regexp
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load %s configuration : %s", svcName, err.Error())
	}
	passRegex, err := regexp.Compile(cfg.PassRegexText)
	if err != nil {
		log.Fatalf("invalid password validation rules %s\n", cfg.PassRegexText)
	}
	cfg.PassRegex = passRegex

	logger, err := mflog.New(os.Stdout, cfg.LogLevel)
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed to init logger: %s", err.Error()))
	}

	var exitCode int
	defer mflog.ExitWithError(&exitCode)

	if cfg.InstanceID == "" {
		if cfg.InstanceID, err = uuid.New().ID(); err != nil {
			logger.Error(fmt.Sprintf("failed to generate instanceID: %s", err))
			exitCode = 1
			return
		}
	}

	ec := email.Config{}
	if err := env.Parse(&ec); err != nil {
		logger.Error(fmt.Sprintf("failed to load email configuration : %s", err.Error()))
		exitCode = 1
		return
	}

	dbConfig := pgclient.Config{Name: defDB}
	if err := dbConfig.LoadEnv(envPrefixDB); err != nil {
		logger.Fatal(err.Error())
	}
	db, err := pgclient.SetupWithConfig(envPrefixDB, *clientspg.Migration(), dbConfig)
	if err != nil {
		logger.Error(err.Error())
		exitCode = 1
		return
	}
	defer db.Close()
	fmt.Println("JG", cfg.JaegerURL, cfg.InstanceID, svcName)

	tp, err := jaegerclient.NewProvider(svcName, cfg.JaegerURL, cfg.InstanceID)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to init Jaeger: %s", err))
		exitCode = 1
		return
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Error(fmt.Sprintf("error shutting down tracer provider: %v", err))
		}
	}()
	tracer := tp.Tracer(svcName)
	fmt.Println("TRACER", trace.FlagsSampled)
	// Setup new redis event store client
	esClient, err := redisclient.Setup(envPrefixES)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer esClient.Close()

	csvc, gsvc := newService(ctx, db, dbConfig, esClient, tracer, cfg, ec, logger)

	httpServerConfig := server.Config{Port: defSvcHTTPPort}
	if err := env.Parse(&httpServerConfig, env.Options{Prefix: envPrefixHTTP}); err != nil {
		logger.Error(fmt.Sprintf("failed to load %s HTTP server configuration : %s", svcName, err.Error()))
		exitCode = 1
		return
	}
	mux := chi.NewRouter()
	httpSvc := httpserver.New(ctx, cancel, svcName, httpServerConfig, capi.MakeHandler(csvc, gsvc, mux, logger, cfg.InstanceID), logger)

	if cfg.SendTelemetry {
		chc := chclient.New(svcName, mainflux.Version, logger, cancel)
		go chc.CallHome(ctx)
	}

	g.Go(func() error {
		return httpSvc.Start()
	})

	g.Go(func() error {
		return server.StopSignalHandler(ctx, cancel, logger, svcName, httpSvc)
	})

	if err := g.Wait(); err != nil {
		logger.Error(fmt.Sprintf("users service terminated: %s", err))
	}
}

func newService(ctx context.Context, db *sqlx.DB, dbConfig pgclient.Config, esClient *redis.Client, tracer trace.Tracer, c config, ec email.Config, logger mflog.Logger) (clients.Service, groups.Service) {
	database := postgres.NewDatabase(db, dbConfig, tracer)
	cRepo := uclients.NewRepository(database)
	gRepo := gpostgres.New(database)

	idp := uuid.New()
	hsr := hasher.New()

	aDuration, err := time.ParseDuration(c.AccessDuration)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to parse access token duration: %s", err.Error()))
	}
	rDuration, err := time.ParseDuration(c.RefreshDuration)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to parse refresh token duration: %s", err.Error()))
	}
	tokenizer := jwt.NewRepository([]byte(c.SecretKey), aDuration, rDuration)

	emailer, err := emailer.New(c.ResetURL, &ec)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to configure e-mailing util: %s", err.Error()))
	}
	csvc := clients.NewService(cRepo, tokenizer, emailer, hsr, idp, c.PassRegex)
	gsvc := mfgroups.NewService(gRepo, idp)

	csvc = ucache.NewEventStoreMiddleware(ctx, csvc, esClient)
	gsvc = gcache.NewEventStoreMiddleware(ctx, gsvc, esClient)

	csvc = ctracing.New(csvc, tracer)
	csvc = capi.LoggingMiddleware(csvc, logger)
	counter, latency := internal.MakeMetrics(svcName, "api")
	csvc = capi.MetricsMiddleware(csvc, counter, latency)

	gsvc = gtracing.New(gsvc, tracer)
	gsvc = gapi.LoggingMiddleware(gsvc, logger)
	counter, latency = internal.MakeMetrics("groups", "api")
	gsvc = gapi.MetricsMiddleware(gsvc, counter, latency)

	if err := createAdmin(ctx, c, cRepo, hsr, csvc); err != nil {
		logger.Error(fmt.Sprintf("failed to create admin client: %s", err))
	}
	return csvc, gsvc
}

func createAdmin(ctx context.Context, c config, crepo uclients.Repository, hsr clients.Hasher, svc clients.Service) error {
	id, err := uuid.New().ID()
	if err != nil {
		return err
	}
	hash, err := hsr.Hash(c.AdminPassword)
	if err != nil {
		return err
	}

	client := mfclients.Client{
		ID:   id,
		Name: "admin",
		Credentials: mfclients.Credentials{
			Identity: c.AdminEmail,
			Secret:   hash,
		},
		Metadata: mfclients.Metadata{
			"role": "admin",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Role:      mfclients.AdminRole,
		Status:    mfclients.EnabledStatus,
	}

	if _, err := crepo.RetrieveByIdentity(ctx, client.Credentials.Identity); err == nil {
		return nil
	}

	// Create an admin
	if _, err = crepo.Save(ctx, client); err != nil {
		return err
	}
	if _, err = svc.IssueToken(ctx, c.AdminEmail, c.AdminPassword); err != nil {
		return err
	}

	return nil
}

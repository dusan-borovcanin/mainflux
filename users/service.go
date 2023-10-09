// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package users

import (
	"context"
	"regexp"
	"time"

	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/internal/apiutil"
	mfclients "github.com/mainflux/mainflux/pkg/clients"
	"github.com/mainflux/mainflux/pkg/errors"
	"github.com/mainflux/mainflux/users/jwt"
	"github.com/mainflux/mainflux/users/postgres"
)

const ()

const (
	administratorRelationKey = "administrator"
	directMemberRelation     = "direct_member"
	createRelation           = "create"
	ownerRelation            = "owner"

	adminPermission      = "admin"
	memberPermission     = "member"
	createUserPermission = "create_user"
	deletePermission     = "delete"
	updatePermission     = "update"
	viewPermission       = "view"

	userKind  = "users"
	tokenKind = "token"

	userType  = "user"
	groupType = "group"
	// organizationType = "organization"

	mainfluxObject = "mainflux"
	anyBodySubject = "_any_body"
)

var (
	// ErrMissingResetToken indicates malformed or missing reset token
	// for reseting password.
	ErrMissingResetToken = errors.New("missing reset token")

	// ErrRecoveryToken indicates error in generating password recovery token.
	ErrRecoveryToken = errors.New("failed to generate password recovery token")

	// ErrGetToken indicates error in getting signed token.
	ErrGetToken = errors.New("failed to fetch signed token")

	// ErrPasswordFormat indicates weak password.
	ErrPasswordFormat = errors.New("password does not meet the requirements")
)

// Service unites Clients and JWT services.
type Service interface {
	ClientService
	jwt.Service
}

type service struct {
	clients    postgres.Repository
	idProvider mainflux.IDProvider
	auth       mainflux.AuthServiceClient
	hasher     Hasher
	// tokens     jwt.Repository
	email     Emailer
	passRegex *regexp.Regexp
}

// NewService returns a new Users service implementation.
func NewService(c postgres.Repository, a mainflux.AuthServiceClient, e Emailer, h Hasher, idp mainflux.IDProvider, pr *regexp.Regexp) Service {
	return service{
		clients:    c,
		auth:       a,
		hasher:     h,
		email:      e,
		idProvider: idp,
		passRegex:  pr,
	}
}

func (svc service) RegisterClient(ctx context.Context, token string, cli mfclients.Client) (mfclients.Client, error) {
	// We don't check the error currently since we can register client with empty token
	ownerID, _ := svc.Identify(ctx, token)

	clientID, err := svc.idProvider.ID()
	if err != nil {
		return mfclients.Client{}, err
	}
	if cli.Owner == "" && ownerID != "" {
		cli.Owner = ownerID
	}
	if cli.Credentials.Secret == "" {
		return mfclients.Client{}, apiutil.ErrMissingSecret
	}
	hash, err := svc.hasher.Hash(cli.Credentials.Secret)
	if err != nil {
		return mfclients.Client{}, errors.Wrap(errors.ErrMalformedEntity, err)
	}
	cli.Credentials.Secret = hash
	if cli.Status != mfclients.DisabledStatus && cli.Status != mfclients.EnabledStatus {
		return mfclients.Client{}, apiutil.ErrInvalidStatus
	}
	if cli.Role != mfclients.UserRole && cli.Role != mfclients.AdminRole {
		return mfclients.Client{}, apiutil.ErrInvalidRole
	}
	cli.ID = clientID
	cli.CreatedAt = time.Now()

	client, err := svc.clients.Save(ctx, cli)
	if err != nil {
		return mfclients.Client{}, err
	}

	return client, nil
}

func (svc service) IssueToken(ctx context.Context, identity, secret string) (jwt.Token, error) {
	dbUser, err := svc.clients.RetrieveByIdentity(ctx, identity)
	if err != nil {
		return jwt.Token{}, err
	}
	if err := svc.hasher.Compare(secret, dbUser.Credentials.Secret); err != nil {
		return jwt.Token{}, errors.Wrap(errors.ErrLogin, err)
	}
	tkn, err := svc.auth.Issue(ctx, &mainflux.IssueReq{Id: dbUser.ID, Email: dbUser.Credentials.Identity, Type: 0})
	if err != nil {
		return jwt.Token{}, errors.Wrap(errors.ErrNotFound, err)
	}
	return parseToken(tkn)
}

func (svc service) RefreshToken(ctx context.Context, refreshToken string) (jwt.Token, error) {
	tkn, err := svc.auth.Refresh(ctx, &mainflux.RefreshReq{Value: refreshToken})
	if err != nil {
		return jwt.Token{}, err
	}
	return parseToken(tkn)
}

func (svc service) ViewClient(ctx context.Context, token string, id string) (mfclients.Client, error) {
	tokenUserID, err := svc.Identify(ctx, token)
	if err != nil {
		return mfclients.Client{}, err
	}

	if tokenUserID != id {
		if err := svc.isOwner(ctx, id, tokenUserID); err != nil {
			return mfclients.Client{}, err
		}
	}

	client, err := svc.clients.RetrieveByID(ctx, id)
	if err != nil {
		return mfclients.Client{}, err
	}
	client.Credentials.Secret = ""

	return client, nil
}

func (svc service) ViewProfile(ctx context.Context, token string) (mfclients.Client, error) {
	id, err := svc.Identify(ctx, token)
	if err != nil {
		return mfclients.Client{}, err
	}
	client, err := svc.clients.RetrieveByID(ctx, id)
	if err != nil {
		return mfclients.Client{}, err
	}
	client.Credentials.Secret = ""

	return client, nil
}

func (svc service) ListClients(ctx context.Context, token string, pm mfclients.Page) (mfclients.ClientsPage, error) {
	id, err := svc.Identify(ctx, token)
	if err != nil {
		return mfclients.ClientsPage{}, err
	}

	// switch err := svc.authorize(ctx, id, clientsObjectKey, listRelationKey); err {
	// // If the user is admin, fetch all users from database.
	// case nil:
	// 	switch {
	// 	// visibility = all
	// 	case pm.SharedBy == myKey && pm.Owner == myKey:
	// 		pm.SharedBy = ""
	// 		pm.Owner = ""
	// 	// visibility = shared
	// 	case pm.SharedBy == myKey && pm.Owner != myKey:
	// 		pm.SharedBy = id
	// 		pm.Owner = ""
	// 	// visibility = mine
	// 	case pm.Owner == myKey && pm.SharedBy != myKey:
	// 		pm.Owner = id
	// 		pm.SharedBy = ""
	// 	}

	// // If the user is not admin, fetch users that they own or are shared with them.
	// default:
	// 	switch {
	// 	// visibility = all
	// 	case pm.SharedBy == myKey && pm.Owner == myKey:
	// 		pm.SharedBy = id
	// 		pm.Owner = id
	// 	// visibility = shared
	// 	case pm.SharedBy == myKey && pm.Owner != myKey:
	// 		pm.SharedBy = id
	// 		pm.Owner = ""
	// 	// visibility = mine
	// 	case pm.Owner == myKey && pm.SharedBy != myKey:
	// 		pm.Owner = id
	// 		pm.SharedBy = ""
	// 	default:
	// 		pm.Owner = id
	// 	}
	// 	pm.Action = listRelationKey
	// }
	pm.Owner = id

	clients, err := svc.clients.RetrieveAll(ctx, pm)
	if err != nil {
		return mfclients.ClientsPage{}, err
	}

	return clients, nil
}

func (svc service) UpdateClient(ctx context.Context, token string, cli mfclients.Client) (mfclients.Client, error) {
	tokenUserID, err := svc.Identify(ctx, token)
	if err != nil {
		return mfclients.Client{}, err
	}

	if tokenUserID != cli.ID {
		if err := svc.isOwner(ctx, cli.ID, tokenUserID); err != nil {
			return mfclients.Client{}, err
		}
	}

	client := mfclients.Client{
		ID:        cli.ID,
		Name:      cli.Name,
		Metadata:  cli.Metadata,
		UpdatedAt: time.Now(),
		UpdatedBy: tokenUserID,
	}

	return svc.clients.Update(ctx, client)
}

func (svc service) UpdateClientTags(ctx context.Context, token string, cli mfclients.Client) (mfclients.Client, error) {
	tokenUserID, err := svc.Identify(ctx, token)
	if err != nil {
		return mfclients.Client{}, err
	}

	if tokenUserID != cli.ID {
		if err := svc.isOwner(ctx, cli.ID, tokenUserID); err != nil {
			return mfclients.Client{}, err
		}
	}

	client := mfclients.Client{
		ID:        cli.ID,
		Tags:      cli.Tags,
		UpdatedAt: time.Now(),
		UpdatedBy: tokenUserID,
	}

	return svc.clients.UpdateTags(ctx, client)
}

func (svc service) UpdateClientIdentity(ctx context.Context, token, clientID, identity string) (mfclients.Client, error) {
	tokenUserID, err := svc.Identify(ctx, token)
	if err != nil {
		return mfclients.Client{}, err
	}

	if tokenUserID != clientID {
		if err := svc.isOwner(ctx, clientID, tokenUserID); err != nil {
			return mfclients.Client{}, err
		}
	}

	cli := mfclients.Client{
		ID: clientID,
		Credentials: mfclients.Credentials{
			Identity: identity,
		},
		UpdatedAt: time.Now(),
		UpdatedBy: tokenUserID,
	}
	return svc.clients.UpdateIdentity(ctx, cli)
}

func (svc service) GenerateResetToken(ctx context.Context, email, host string) error {
	// client, err := svc.clients.RetrieveByIdentity(ctx, email)
	// if err != nil || client.Credentials.Identity == "" {
	// 	return errors.ErrNotFound
	// }
	// claims := jwt.Claims{
	// 	ClientID: client.ID,
	// 	Email:    client.Credentials.Identity,
	// }
	// t, err := svc.tokens.Issue(ctx, claims)
	// if err != nil {
	// 	return errors.Wrap(ErrRecoveryToken, err)
	// }
	// return svc.SendPasswordReset(ctx, host, email, client.Name, t.AccessToken)
	return nil
}

func (svc service) ResetSecret(ctx context.Context, resetToken, secret string) error {
	id, err := svc.Identify(ctx, resetToken)
	if err != nil {
		return errors.Wrap(errors.ErrAuthentication, err)
	}
	c, err := svc.clients.RetrieveByID(ctx, id)
	if err != nil {
		return err
	}
	if c.Credentials.Identity == "" {
		return errors.ErrNotFound
	}
	if !svc.passRegex.MatchString(secret) {
		return ErrPasswordFormat
	}
	secret, err = svc.hasher.Hash(secret)
	if err != nil {
		return err
	}
	c = mfclients.Client{
		Credentials: mfclients.Credentials{
			Identity: c.Credentials.Identity,
			Secret:   secret,
		},
		UpdatedAt: time.Now(),
		UpdatedBy: id,
	}
	if _, err := svc.clients.UpdateSecret(ctx, c); err != nil {
		return err
	}
	return nil
}

func (svc service) UpdateClientSecret(ctx context.Context, token, oldSecret, newSecret string) (mfclients.Client, error) {
	id, err := svc.Identify(ctx, token)
	if err != nil {
		return mfclients.Client{}, err
	}
	if !svc.passRegex.MatchString(newSecret) {
		return mfclients.Client{}, ErrPasswordFormat
	}
	dbClient, err := svc.clients.RetrieveByID(ctx, id)
	if err != nil {
		return mfclients.Client{}, err
	}
	if _, err := svc.IssueToken(ctx, dbClient.Credentials.Identity, oldSecret); err != nil {
		return mfclients.Client{}, err
	}
	newSecret, err = svc.hasher.Hash(newSecret)
	if err != nil {
		return mfclients.Client{}, err
	}
	dbClient.Credentials.Secret = newSecret
	dbClient.UpdatedAt = time.Now()
	dbClient.UpdatedBy = id

	return svc.clients.UpdateSecret(ctx, dbClient)
}

func (svc service) SendPasswordReset(_ context.Context, host, email, user, token string) error {
	to := []string{email}
	return svc.email.SendPasswordReset(to, host, user, token)
}

func (svc service) UpdateClientOwner(ctx context.Context, token string, cli mfclients.Client) (mfclients.Client, error) {
	tokenUserID, err := svc.Identify(ctx, token)
	if err != nil {
		return mfclients.Client{}, err
	}

	if tokenUserID != cli.ID {
		if err := svc.isOwner(ctx, cli.ID, tokenUserID); err != nil {
			return mfclients.Client{}, err
		}
	}
	client := mfclients.Client{
		ID:        cli.ID,
		Owner:     cli.Owner,
		UpdatedAt: time.Now(),
		UpdatedBy: tokenUserID,
	}

	return svc.clients.UpdateOwner(ctx, client)
}

func (svc service) EnableClient(ctx context.Context, token, id string) (mfclients.Client, error) {
	client := mfclients.Client{
		ID:        id,
		UpdatedAt: time.Now(),
		Status:    mfclients.EnabledStatus,
	}
	client, err := svc.changeClientStatus(ctx, token, client)
	if err != nil {
		return mfclients.Client{}, errors.Wrap(mfclients.ErrEnableClient, err)
	}

	return client, nil
}

func (svc service) DisableClient(ctx context.Context, token, id string) (mfclients.Client, error) {
	client := mfclients.Client{
		ID:        id,
		UpdatedAt: time.Now(),
		Status:    mfclients.DisabledStatus,
	}
	client, err := svc.changeClientStatus(ctx, token, client)
	if err != nil {
		return mfclients.Client{}, errors.Wrap(mfclients.ErrDisableClient, err)
	}

	return client, nil
}

func (svc service) changeClientStatus(ctx context.Context, token string, client mfclients.Client) (mfclients.Client, error) {
	tokenUserID, err := svc.Identify(ctx, token)
	if err != nil {
		return mfclients.Client{}, err
	}
	if tokenUserID != client.ID {
		if err := svc.isOwner(ctx, client.ID, tokenUserID); err != nil {
			return mfclients.Client{}, err
		}
	}
	dbClient, err := svc.clients.RetrieveByID(ctx, client.ID)
	if err != nil {
		return mfclients.Client{}, err
	}
	if dbClient.Status == client.Status {
		return mfclients.Client{}, mfclients.ErrStatusAlreadyAssigned
	}
	client.UpdatedBy = tokenUserID
	return svc.clients.ChangeStatus(ctx, client)
}

func (svc service) ListMembers(ctx context.Context, token, groupID string, pm mfclients.Page) (mfclients.MembersPage, error) {
	if _, err := svc.authorize(ctx, userType, tokenKind, token, pm.Permission, groupType, groupID); err != nil {
		return mfclients.MembersPage{}, err
	}
	uids, err := svc.auth.ListAllSubjects(ctx, &mainflux.ListSubjectsReq{
		SubjectType: userType,
		Permission:  pm.Permission,
		Object:      groupID,
		ObjectType:  groupType,
	})
	if err != nil {
		return mfclients.MembersPage{}, err
	}

	pm.IDs = uids.Policies

	cp, err := svc.clients.RetrieveAll(ctx, pm)
	if err != nil {
		return mfclients.MembersPage{}, err
	}
	return mfclients.MembersPage{
		Page:    cp.Page,
		Members: cp.Clients,
	}, nil
}

func (svc *service) isOwner(ctx context.Context, clientID, ownerID string) error {
	return svc.clients.IsOwner(ctx, clientID, ownerID)
}

func (svc *service) authorize(ctx context.Context, subjType, subjKind, subj, perm, objType, obj string) (string, error) {
	req := &mainflux.AuthorizeReq{
		SubjectType: subjType,
		SubjectKind: subjKind,
		Subject:     subj,
		Permission:  perm,
		ObjectType:  objType,
		Object:      obj,
	}
	res, err := svc.auth.Authorize(ctx, req)
	if err != nil {
		return "", errors.Wrap(errors.ErrAuthorization, err)
	}

	if !res.GetAuthorized() {
		return "", errors.ErrAuthorization
	}
	return res.GetId(), nil
}

func (svc service) Identify(ctx context.Context, token string) (string, error) {
	user, err := svc.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return "", err
	}
	return user.GetId(), nil
}

// Auth helpers
func (svc service) issue(ctx context.Context, id, email string, keyType uint32) (jwt.Token, error) {
	tkn, err := svc.auth.Issue(ctx, &mainflux.IssueReq{Id: id, Email: email, Type: keyType})
	if err != nil {
		return jwt.Token{}, errors.Wrap(errors.ErrNotFound, err)
	}
	extra := tkn.Extra.AsMap()["refresh_token"]
	refresh, ok := extra.(string)
	if !ok {
		return jwt.Token{}, errors.ErrAuthentication
	}
	ret := jwt.Token{
		AccessToken:  tkn.GetValue(),
		RefreshToken: refresh,
		AccessType:   "bearer",
	}

	return ret, nil
}

func parseToken(tkn *mainflux.Token) (jwt.Token, error) {
	if tkn == nil {
		return jwt.Token{}, errors.New("invalid token")
	}
	extra := tkn.Extra.AsMap()["refresh_token"]
	refresh, ok := extra.(string)
	if !ok {
		return jwt.Token{}, errors.ErrAuthentication
	}
	ret := jwt.Token{
		AccessToken:  tkn.GetValue(),
		RefreshToken: refresh,
		AccessType:   "bearer",
	}

	return ret, nil
}

func (svc service) claimOwnership(ctx context.Context, subjectType, subject, relation, permission, objectType, object string) error {
	req := &mainflux.AddPolicyReq{
		SubjectType: subjectType,
		Subject:     subject,
		Relation:    relation,
		Permission:  permission,
		Object:      object,
		ObjectType:  objectType,
	}
	res, err := svc.auth.AddPolicy(ctx, req)
	if err != nil {
		return errors.Wrap(errors.ErrAuthorization, err)
	}
	if !res.GetAuthorized() {
		return errors.ErrAuthorization
	}
	return nil
}

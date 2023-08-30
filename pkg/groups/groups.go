// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package groups

import (
	"context"
	"time"

	"github.com/mainflux/mainflux/pkg/clients"
)

const (
	// MaxLevel represents the maximum group hierarchy level.
	MaxLevel = uint64(5)
	// MinLevel represents the minimum group hierarchy level.
	MinLevel = uint64(0)
)

// Group represents the group of Clients.
// Indicates a level in tree hierarchy. Root node is level 1.
// Path in a tree consisting of group IDs
// Paths are unique per owner.
type Group struct {
	ID          string           `json:"id"`
	Owner       string           `json:"owner_id"`
	Parent      string           `json:"parent_id,omitempty"`
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Metadata    clients.Metadata `json:"metadata,omitempty"`
	Level       int              `json:"level,omitempty"`
	Path        string           `json:"path,omitempty"`
	Children    []*Group         `json:"children,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at,omitempty"`
	UpdatedBy   string           `json:"updated_by,omitempty"`
	Status      clients.Status   `json:"status"`
}

// Memberships contains page related metadata as well as list of memberships that
// belong to this page.
type Memberships struct {
	PageMeta
	Groups []Group
}

// Page contains page related metadata as well as list
// of Groups that belong to the page.
type Page struct {
	PageMeta
	Path      string
	Level     uint64
	ID        string
	Direction int64 // ancestors (-1) or descendants (+1)
	Groups    []Group
}

// Repository specifies a group persistence API.
type Repository interface {
	// Save group.
	Save(ctx context.Context, g Group) (Group, error)

	// Update a group.
	Update(ctx context.Context, g Group) (Group, error)

	// RetrieveByID retrieves group by its id.
	RetrieveByID(ctx context.Context, id string) (Group, error)

	// RetrieveAll retrieves all groups.
	RetrieveAll(ctx context.Context, gm Page) (Page, error)

	// Memberships retrieves everything that is assigned to a group identified by clientID.
	Memberships(ctx context.Context, clientID string, gm Page) (Memberships, error)

	// ChangeStatus changes groups status to active or inactive
	ChangeStatus(ctx context.Context, group Group) (Group, error)
}

type Service interface {
	// CreateGroup creates new  group.
	CreateGroup(ctx context.Context, token string, g Group) (Group, error)

	// UpdateGroup updates the group identified by the provided ID.
	UpdateGroup(ctx context.Context, token string, g Group) (Group, error)

	// ViewGroup retrieves data about the group identified by ID.
	ViewGroup(ctx context.Context, token, id string) (Group, error)

	// ListGroups retrieves
	ListGroups(ctx context.Context, token string, gm Page) (Page, error)

	// ListMemberships retrieves everything that is assigned to a group identified by clientID.
	ListMemberships(ctx context.Context, token, clientID string, gm Page) (Memberships, error)

	// EnableGroup logically enables the group identified with the provided ID.
	EnableGroup(ctx context.Context, token, id string) (Group, error)

	// DisableGroup logically disables the group identified with the provided ID.
	DisableGroup(ctx context.Context, token, id string) (Group, error)
}

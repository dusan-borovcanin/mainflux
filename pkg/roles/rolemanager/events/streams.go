// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package events

import (
	"context"

	"github.com/absmach/supermq/pkg/authn"
	"github.com/absmach/supermq/pkg/events"
	"github.com/absmach/supermq/pkg/roles"
)

var _ roles.RoleManager = (*RoleManagerEventStore)(nil)

type RoleManagerEventStore struct {
	events.Publisher
	svc             roles.RoleManager
	operationPrefix string
	svcName         string
}

// NewEventStoreMiddleware returns wrapper around auth service that sends
// events to event store.
func NewRoleManagerEventStore(svcName, operationPrefix string, svc roles.RoleManager, publisher events.Publisher) RoleManagerEventStore {
	return RoleManagerEventStore{
		svcName:   svcName,
		svc:       svc,
		Publisher: publisher,
	}
}

func (rmes *RoleManagerEventStore) AddRole(ctx context.Context, session authn.Session, entityID, roleName string, optionalActions []string, optionalMembers []string) (roles.Role, error) {
	ro, err := rmes.svc.AddRole(ctx, session, entityID, roleName, optionalActions, optionalMembers)
	if err != nil {
		return ro, err
	}

	e := addRoleEvent{
		operationPrefix: rmes.operationPrefix,
		Role:            ro,
	}
	if err := rmes.Publish(ctx, e); err != nil {
		return ro, err
	}
	return ro, nil
}

func (rmes *RoleManagerEventStore) RemoveRole(ctx context.Context, session authn.Session, entityID, roleName string) error {
	if err := rmes.svc.RemoveRole(ctx, session, entityID, roleName); err != nil {
		return err
	}
	e := removeRoleEvent{
		operationPrefix: rmes.operationPrefix,
		roleName:        roleName,
		entityID:        entityID,
	}
	if err := rmes.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}

func (rmes *RoleManagerEventStore) UpdateRoleName(ctx context.Context, session authn.Session, entityID, oldRoleName, newRoleName string) (roles.Role, error) {
	ro, err := rmes.svc.UpdateRoleName(ctx, session, entityID, oldRoleName, newRoleName)
	if err != nil {
		return ro, err
	}

	e := updateRoleEvent{
		operationPrefix: rmes.operationPrefix,
		Role:            ro,
	}
	if err := rmes.Publish(ctx, e); err != nil {
		return ro, err
	}
	return ro, nil
}

func (rmes *RoleManagerEventStore) RetrieveRole(ctx context.Context, session authn.Session, entityID, roleName string) (roles.Role, error) {
	ro, err := rmes.svc.RetrieveRole(ctx, session, entityID, roleName)
	if err != nil {
		return ro, err
	}
	e := retrieveRoleEvent{
		operationPrefix: rmes.operationPrefix,
		Role:            ro,
	}
	if err := rmes.Publish(ctx, e); err != nil {
		return ro, err
	}
	return ro, nil
}

func (rmes *RoleManagerEventStore) RetrieveAllRoles(ctx context.Context, session authn.Session, entityID string, limit, offset uint64) (roles.RolePage, error) {
	rp, err := rmes.svc.RetrieveAllRoles(ctx, session, entityID, limit, offset)
	if err != nil {
		return rp, err
	}

	e := retrieveAllRolesEvent{
		operationPrefix: rmes.operationPrefix,
		entityID:        entityID,
		limit:           limit,
		offset:          offset,
	}
	if err := rmes.Publish(ctx, e); err != nil {
		return rp, err
	}
	return rp, nil
}

func (rmes *RoleManagerEventStore) ListAvailableActions(ctx context.Context, session authn.Session) ([]string, error) {
	actions, err := rmes.svc.ListAvailableActions(ctx, session)
	if err != nil {
		return actions, err
	}
	e := listAvailableActionsEvent{
		operationPrefix: rmes.operationPrefix,
	}
	if err := rmes.Publish(ctx, e); err != nil {
		return actions, err
	}
	return actions, nil
}

func (rmes *RoleManagerEventStore) RoleAddActions(ctx context.Context, session authn.Session, entityID, roleName string, actions []string) ([]string, error) {
	actions, err := rmes.svc.RoleAddActions(ctx, session, entityID, roleName, actions)
	if err != nil {
		return actions, err
	}
	e := roleAddActionsEvent{
		operationPrefix: rmes.operationPrefix,
		entityID:        entityID,
		roleName:        roleName,
		actions:         actions,
	}
	if err := rmes.Publish(ctx, e); err != nil {
		return actions, err
	}
	return actions, nil
}

func (rmes *RoleManagerEventStore) RoleListActions(ctx context.Context, session authn.Session, entityID, roleName string) ([]string, error) {
	actions, err := rmes.svc.RoleListActions(ctx, session, entityID, roleName)
	if err != nil {
		return actions, err
	}

	e := roleListActionsEvent{
		operationPrefix: rmes.operationPrefix,
		entityID:        entityID,
		roleName:        roleName,
	}
	if err := rmes.Publish(ctx, e); err != nil {
		return actions, err
	}
	return actions, nil
}

func (rmes *RoleManagerEventStore) RoleCheckActionsExists(ctx context.Context, session authn.Session, entityID, roleName string, actions []string) (bool, error) {
	isAllExists, err := rmes.svc.RoleCheckActionsExists(ctx, session, entityID, roleName, actions)
	if err != nil {
		return isAllExists, err
	}

	e := roleCheckActionsExistsEvent{
		operationPrefix: rmes.operationPrefix,
		entityID:        entityID,
		roleName:        roleName,
		actions:         actions,
		isAllExists:     isAllExists,
	}
	if err := rmes.Publish(ctx, e); err != nil {
		return isAllExists, err
	}
	return isAllExists, nil
}

func (rmes *RoleManagerEventStore) RoleRemoveActions(ctx context.Context, session authn.Session, entityID, roleName string, actions []string) (err error) {
	if err := rmes.svc.RoleRemoveActions(ctx, session, entityID, roleName, actions); err != nil {
		return err
	}

	e := roleRemoveActionsEvent{
		operationPrefix: rmes.operationPrefix,
		entityID:        entityID,
		roleName:        roleName,
		actions:         actions,
	}
	if err := rmes.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}

func (rmes *RoleManagerEventStore) RoleRemoveAllActions(ctx context.Context, session authn.Session, entityID, roleName string) error {
	if err := rmes.svc.RoleRemoveAllActions(ctx, session, entityID, roleName); err != nil {
		return err
	}

	e := roleRemoveAllActionsEvent{
		operationPrefix: rmes.operationPrefix,
		entityID:        entityID,
		roleName:        roleName,
	}
	if err := rmes.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}

func (rmes *RoleManagerEventStore) RoleAddMembers(ctx context.Context, session authn.Session, entityID, roleName string, members []string) ([]string, error) {
	mems, err := rmes.svc.RoleAddMembers(ctx, session, entityID, roleName, members)
	if err != nil {
		return mems, err
	}

	e := roleAddMembersEvent{
		operationPrefix: rmes.operationPrefix,
		entityID:        entityID,
		roleName:        roleName,
		members:         members,
	}
	if err := rmes.Publish(ctx, e); err != nil {
		return mems, err
	}
	return mems, nil
}

func (rmes *RoleManagerEventStore) RoleListMembers(ctx context.Context, session authn.Session, entityID, roleName string, limit, offset uint64) (roles.MembersPage, error) {
	mp, err := rmes.svc.RoleListMembers(ctx, session, entityID, roleName, limit, offset)
	if err != nil {
		return mp, err
	}

	e := roleListMembersEvent{
		operationPrefix: rmes.operationPrefix,
		entityID:        entityID,
		roleName:        roleName,
		limit:           limit,
		offset:          offset,
	}
	if err := rmes.Publish(ctx, e); err != nil {
		return mp, err
	}
	return mp, nil
}

func (rmes *RoleManagerEventStore) RoleCheckMembersExists(ctx context.Context, session authn.Session, entityID, roleName string, members []string) (bool, error) {
	isAllExists, err := rmes.svc.RoleCheckMembersExists(ctx, session, entityID, roleName, members)
	if err != nil {
		return isAllExists, err
	}

	e := roleCheckMembersExistsEvent{
		operationPrefix: rmes.operationPrefix,
		entityID:        entityID,
		roleName:        roleName,
		members:         members,
	}
	if err := rmes.Publish(ctx, e); err != nil {
		return isAllExists, err
	}
	return isAllExists, nil
}

func (rmes *RoleManagerEventStore) RoleRemoveMembers(ctx context.Context, session authn.Session, entityID, roleName string, members []string) (err error) {
	if err := rmes.svc.RoleRemoveMembers(ctx, session, entityID, roleName, members); err != nil {
		return err
	}

	e := roleRemoveMembersEvent{
		operationPrefix: rmes.operationPrefix,
		entityID:        entityID,
		roleName:        roleName,
		members:         members,
	}
	if err := rmes.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}

func (rmes *RoleManagerEventStore) RoleRemoveAllMembers(ctx context.Context, session authn.Session, entityID, roleName string) (err error) {
	if err := rmes.svc.RoleRemoveAllMembers(ctx, session, entityID, roleName); err != nil {
		return err
	}

	e := roleRemoveAllMembersEvent{
		operationPrefix: rmes.operationPrefix,
		entityID:        entityID,
		roleName:        roleName,
	}
	if err := rmes.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}

func (rmes *RoleManagerEventStore) RemoveMemberFromAllRoles(ctx context.Context, session authn.Session, memberID string) (err error) {
	if err := rmes.svc.RemoveMemberFromAllRoles(ctx, session, memberID); err != nil {
		return err
	}

	e := removeMemberFromAllRolesEvent{
		operationPrefix: rmes.operationPrefix,
		memberID:        memberID,
	}
	if err := rmes.Publish(ctx, e); err != nil {
		return err
	}
	return nil
}
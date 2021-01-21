package org

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	uniqueOrgname           = "org_name"
	OrgAddedEventType       = orgEventTypePrefix + "added"
	OrgChangedEventType     = orgEventTypePrefix + "changed"
	OrgDeactivatedEventType = orgEventTypePrefix + "deactivated"
	OrgReactivatedEventType = orgEventTypePrefix + "reactivated"
	OrgRemovedEventType     = orgEventTypePrefix + "removed"
)

type OrgnameUniqueConstraint struct {
	uniqueType string
	orgName    string
	action     eventstore.UniqueConstraintAction
}

func NewAddOrgnameUniqueConstraint(orgName string) *eventstore.EventUniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		uniqueOrgname,
		orgName,
		"Errors.Org.AlreadyExists")
}

func NewRemoveUsernameUniqueConstraint(orgName string) *eventstore.EventUniqueConstraint {
	return eventstore.NewRemoveEventUniqueConstraint(
		uniqueOrgname,
		orgName)
}

type OrgAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name string `json:"name,omitempty"`
}

func (e *OrgAddedEvent) Data() interface{} {
	return e
}

func (e *OrgAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddOrgnameUniqueConstraint(e.Name)}
}

func NewOrgAddedEvent(ctx context.Context, name string) *OrgAddedEvent {
	return &OrgAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OrgAddedEventType,
		),
		Name: name,
	}
}

func OrgAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	orgAdded := &OrgAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, orgAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "ORG-Bren2", "unable to unmarshal org added")
	}

	return orgAdded, nil
}

type OrgChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name string `json:"name,omitempty"`
}

func (e *OrgChangedEvent) Data() interface{} {
	return e
}

func (e *OrgChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewOrgChangedEvent(ctx context.Context, name string) *OrgChangedEvent {
	return &OrgChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OrgChangedEventType,
		),
		Name: name,
	}
}

func OrgChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	orgChanged := &OrgChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, orgChanged)
	if err != nil {
		return nil, errors.ThrowInternal(err, "ORG-Bren2", "unable to unmarshal org added")
	}

	return orgChanged, nil
}

type OrgDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *OrgDeactivatedEvent) Data() interface{} {
	return e
}

func (e *OrgDeactivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewOrgDeactivatedEvent(ctx context.Context) *OrgDeactivatedEvent {
	return &OrgDeactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OrgDeactivatedEventType,
		),
	}
}

func OrgDeactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	orgChanged := &OrgDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, orgChanged)
	if err != nil {
		return nil, errors.ThrowInternal(err, "ORG-DAfbs", "unable to unmarshal org deactivated")
	}

	return orgChanged, nil
}

type OrgReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *OrgReactivatedEvent) Data() interface{} {
	return e
}

func (e *OrgReactivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewOrgReactivatedEvent(ctx context.Context) *OrgReactivatedEvent {
	return &OrgReactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OrgReactivatedEventType,
		),
	}
}

func OrgReactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	orgChanged := &OrgReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, orgChanged)
	if err != nil {
		return nil, errors.ThrowInternal(err, "ORG-DAfbs", "unable to unmarshal org deactivated")
	}

	return orgChanged, nil
}

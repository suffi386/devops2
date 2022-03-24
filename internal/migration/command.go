package migration

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/service"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

//SetupStep is the command pushed on the eventstore
type SetupStep struct {
	eventstore.BaseEvent
	migration Migration
	Name      string `json:"name"`
	Error     error  `json:"error,omitempty"`
}

func setupStartedCmd(migration Migration) eventstore.Command {
	ctx := authz.SetCtxData(service.WithService(context.Background(), "system"), authz.CtxData{UserID: "system", OrgID: "SYSTEM", ResourceOwner: "SYSTEM"})
	return &SetupStep{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			eventstore.NewAggregate(ctx, aggregateID, aggregateType, "v1"),
			startedType),
		migration: migration,
		Name:      migration.String(),
	}
}

func setupDoneCmd(migration Migration, err error) eventstore.Command {
	ctx := authz.SetCtxData(service.WithService(context.Background(), "system"), authz.CtxData{UserID: "system", OrgID: "SYSTEM", ResourceOwner: "SYSTEM"})
	s := &SetupStep{
		migration: migration,
		Name:      migration.String(),
		Error:     err,
	}

	typ := doneType
	if err != nil {
		typ = failedType
	}

	s.BaseEvent = *eventstore.NewBaseEventForPush(
		ctx,
		eventstore.NewAggregate(ctx, aggregateID, aggregateType, "v1"),
		typ)

	return s
}

func (s *SetupStep) Data() interface{} {
	return s
}

func (s *SetupStep) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	switch s.Type() {
	case startedType:
		return []*eventstore.EventUniqueConstraint{
			eventstore.NewAddEventUniqueConstraint("migration_started", s.migration.String(), "Errors.Step.Started.AlreadyExists"),
		}
	case failedType:
		return []*eventstore.EventUniqueConstraint{
			eventstore.NewRemoveEventUniqueConstraint("migration_started", s.migration.String()),
		}
	default:
		return []*eventstore.EventUniqueConstraint{
			eventstore.NewAddEventUniqueConstraint("migration_done", s.migration.String(), "Errors.Step.Done.AlreadyExists"),
		}
	}
}

func RegisterMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(startedType, SetupMapper)
	es.RegisterFilterEventMapper(doneType, SetupMapper)
	es.RegisterFilterEventMapper(failedType, SetupMapper)
}

func SetupMapper(event *repository.Event) (eventstore.Event, error) {
	step := &SetupStep{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	if len(event.Data) == 0 {
		return step, nil
	}
	err := json.Unmarshal(event.Data, step)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-O6rVg", "unable to unmarshal step")
	}

	return step, nil
}

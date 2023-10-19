package handler

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	handler2 "github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/static"
)

type Config struct {
	Client     *database.DB
	Eventstore *eventstore.Eventstore

	BulkLimit             uint64
	FailureCountUntilSkip uint64
	HandleActiveInstances time.Duration
	TransactionDuration   time.Duration
	Handlers              map[string]*ConfigOverwrites
}

type ConfigOverwrites struct {
	MinimumCycleDuration time.Duration
}

func Register(ctx context.Context, config Config, view *view.View, static static.Storage) {
	if static == nil {
		return
	}

	newStyling(ctx,
		config.overwrite("Styling"),
		static,
		view,
	).Start(ctx)
}

func (config Config) overwrite(viewModel string) handler2.Config {
	c := handler2.Config{
		Client:                config.Client,
		Eventstore:            config.Eventstore,
		BulkLimit:             uint16(config.BulkLimit),
		RequeueEvery:          3 * time.Minute,
		HandleActiveInstances: config.HandleActiveInstances,
		MaxFailureCount:       uint8(config.FailureCountUntilSkip),
		TransactionDuration:   config.TransactionDuration,
	}
	overwrite, ok := config.Handlers[viewModel]
	if !ok {
		return c
	}
	if overwrite.MinimumCycleDuration > 0 {
		c.RequeueEvery = overwrite.MinimumCycleDuration
	}
	return c
}

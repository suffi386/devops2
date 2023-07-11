package query

import (
	"context"
	"strings"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type OIDCSessionAccessTokenReadModel struct {
	eventstore.WriteModel

	UserID                string
	SessionID             string
	ClientID              string
	Audience              []string
	Scope                 []string
	AuthMethodsReferences []string
	AuthTime              time.Time
	State                 domain.OIDCSessionState
	AccessTokenID         string
	AccessTokenCreation   time.Time
	AccessTokenExpiration time.Time
}

func newOIDCSessionAccessTokenWriteModel(id string) *OIDCSessionAccessTokenReadModel {
	return &OIDCSessionAccessTokenReadModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: id,
		},
	}
}

func (wm *OIDCSessionAccessTokenReadModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *oidcsession.AddedEvent:
			wm.reduceAdded(e)
		case *oidcsession.AccessTokenAddedEvent:
			wm.reduceAccessTokenAdded(e)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *OIDCSessionAccessTokenReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AllowTimeTravel().
		AddQuery().
		AggregateTypes(oidcsession.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			oidcsession.AddedType,
			oidcsession.AccessTokenAddedType,
		).
		Builder()
}

func (wm *OIDCSessionAccessTokenReadModel) reduceAdded(e *oidcsession.AddedEvent) {
	wm.UserID = e.UserID
	wm.SessionID = e.SessionID
	wm.ClientID = e.ClientID
	wm.Audience = e.Audience
	wm.Scope = e.Scope
	wm.AuthMethodsReferences = e.AuthMethodsReferences
	wm.AuthTime = e.AuthTime
	wm.State = domain.OIDCSessionStateActive
}

func (wm *OIDCSessionAccessTokenReadModel) reduceAccessTokenAdded(e *oidcsession.AccessTokenAddedEvent) {
	wm.AccessTokenID = e.ID
	wm.AccessTokenCreation = e.CreationDate()
	wm.AccessTokenExpiration = e.CreationDate().Add(e.Lifetime)
}

func (q *Queries) GetAccessToken(ctx context.Context, token string) (model *OIDCSessionAccessTokenReadModel, err error) {
	split := strings.Split(token, "-")
	if len(split) != 2 {
		return nil, caos_errs.ThrowPermissionDenied(nil, "QUERY-SAhtk", "Errors.OIDCSession.Token.Invalid")
	}
	return q.GetAccessTokenByOIDCSessionAndTokenID(ctx, split[0], split[1])
}

func (q *Queries) GetAccessTokenByOIDCSessionAndTokenID(ctx context.Context, oidcSessionID, tokenID string) (model *OIDCSessionAccessTokenReadModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	model = newOIDCSessionAccessTokenWriteModel(oidcSessionID)
	if err = q.eventstore.FilterToQueryReducer(ctx, model); err != nil {
		return nil, caos_errs.ThrowPermissionDenied(err, "QUERY-ASfe2", "Errors.OIDCSession.Token.Invalid")
	}
	if model.AccessTokenID != tokenID {
		return nil, caos_errs.ThrowPermissionDenied(nil, "QUERY-M2u9w", "Errors.OIDCSession.Token.Invalid")
	}
	return model, nil
}

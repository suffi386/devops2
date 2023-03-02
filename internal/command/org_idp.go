package command

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func (c *Commands) AddOrgGenericOAuthProvider(ctx context.Context, resourceOwner string, provider GenericOAuthProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewOAuthOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgOAuthProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgGenericOAuthProvider(ctx context.Context, resourceOwner, id string, provider GenericOAuthProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewOAuthOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgOAuthProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddOrgGenericOIDCProvider(ctx context.Context, resourceOwner string, provider GenericOIDCProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewOIDCOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgOIDCProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgGenericOIDCProvider(ctx context.Context, resourceOwner, id string, provider GenericOIDCProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewOIDCOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgOIDCProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddOrgJWTProvider(ctx context.Context, resourceOwner string, provider JWTProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewJWTOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgJWTProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgJWTProvider(ctx context.Context, resourceOwner, id string, provider JWTProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewJWTOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgJWTProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddOrgGitHubProvider(ctx context.Context, resourceOwner string, provider GitHubProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgGitHubProvider(orgAgg, resourceOwner, id, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgGitHubProvider(ctx context.Context, resourceOwner, id string, provider GitHubProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgGitHubProvider(orgAgg, resourceOwner, id, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddOrgGitHubEnterpriseProvider(ctx context.Context, resourceOwner string, provider GitHubEnterpriseProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgGitHubEnterpriseProvider(orgAgg, resourceOwner, id, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgGitHubEnterpriseProvider(ctx context.Context, resourceOwner, id string, provider GitHubEnterpriseProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgGitHubEnterpriseProvider(orgAgg, resourceOwner, id, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddOrgGoogleProvider(ctx context.Context, resourceOwner string, provider GoogleProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewGoogleOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgGoogleProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgGoogleProvider(ctx context.Context, resourceOwner, id string, provider GoogleProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewGoogleOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgGoogleProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddOrgLDAPProvider(ctx context.Context, resourceOwner string, provider LDAPProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewLDAPOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgLDAPProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgLDAPProvider(ctx context.Context, resourceOwner, id string, provider LDAPProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewLDAPOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgLDAPProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) DeleteOrgProvider(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareDeleteOrgProvider(orgAgg, resourceOwner, id))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func ExistsOrgIDP(ctx context.Context, filter preparation.FilterToQueryReducer, id, orgID string) (exists bool, err error) {
	writeModel := NewOrgIDPRemoveWriteModel(orgID, id)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return false, err
	}

	if len(events) == 0 {
		return false, nil
	}
	writeModel.AppendEvents(events...)
	if err := writeModel.Reduce(); err != nil {
		return false, err
	}
	return writeModel.State.Exists(), nil
}

func (c *Commands) prepareAddOrgOAuthProvider(a *org.Aggregate, writeModel *OrgOAuthIDPWriteModel, provider GenericOAuthProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-D32ef", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Dbgzf", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-DF4ga", "Errors.Invalid.Argument")
		}
		if provider.AuthorizationEndpoint = strings.TrimSpace(provider.AuthorizationEndpoint); provider.AuthorizationEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-B23bs", "Errors.Invalid.Argument")
		}
		if provider.TokenEndpoint = strings.TrimSpace(provider.TokenEndpoint); provider.TokenEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-D2gj8", "Errors.Invalid.Argument")
		}
		if provider.UserEndpoint = strings.TrimSpace(provider.UserEndpoint); provider.UserEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Fb8jk", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewOAuthIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.ClientID,
					secret,
					provider.AuthorizationEndpoint,
					provider.TokenEndpoint,
					provider.UserEndpoint,
					provider.Scopes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgOAuthProvider(a *org.Aggregate, writeModel *OrgOAuthIDPWriteModel, provider GenericOAuthProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-asfsa", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-D32ef", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Dbgzf", "Errors.Invalid.Argument")
		}
		if provider.AuthorizationEndpoint = strings.TrimSpace(provider.AuthorizationEndpoint); provider.AuthorizationEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-B23bs", "Errors.Invalid.Argument")
		}
		if provider.TokenEndpoint = strings.TrimSpace(provider.TokenEndpoint); provider.TokenEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-D2gj8", "Errors.Invalid.Argument")
		}
		if provider.UserEndpoint = strings.TrimSpace(provider.UserEndpoint); provider.UserEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Fb8jk", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "ORG-JNsd3", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				provider.Name,
				provider.ClientID,
				provider.ClientSecret,
				c.idpConfigEncryption,
				provider.AuthorizationEndpoint,
				provider.TokenEndpoint,
				provider.UserEndpoint,
				provider.Scopes,
				provider.IDPOptions,
			)
			if err != nil || event == nil {
				return nil, err
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddOrgOIDCProvider(a *org.Aggregate, writeModel *OrgOIDCIDPWriteModel, provider GenericOIDCProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Sgtj5", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Hz6zj", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-fb5jm", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Sfdf4", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewOIDCIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.Issuer,
					provider.ClientID,
					secret,
					provider.Scopes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgOIDCProvider(a *org.Aggregate, writeModel *OrgOIDCIDPWriteModel, provider GenericOIDCProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-SAfd3", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Dvf4f", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-BDfr3", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Db3bs", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "ORG-Dg331", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				provider.Name,
				provider.Issuer,
				provider.ClientID,
				provider.ClientSecret,
				c.idpConfigEncryption,
				provider.Scopes,
				provider.IDPOptions,
			)
			if err != nil || event == nil {
				return nil, err
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddOrgJWTProvider(a *org.Aggregate, writeModel *OrgJWTIDPWriteModel, provider JWTProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-JLKef", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-WNJK3", "Errors.Invalid.Argument")
		}
		if provider.JWTEndpoint = strings.TrimSpace(provider.JWTEndpoint); provider.JWTEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-NJKSD", "Errors.Invalid.Argument")
		}
		if provider.KeyEndpoint = strings.TrimSpace(provider.KeyEndpoint); provider.KeyEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-NJKE3", "Errors.Invalid.Argument")
		}
		if provider.HeaderName = strings.TrimSpace(provider.HeaderName); provider.HeaderName == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-2rlks", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewJWTIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.Issuer,
					provider.JWTEndpoint,
					provider.KeyEndpoint,
					provider.HeaderName,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgJWTProvider(a *org.Aggregate, writeModel *OrgJWTIDPWriteModel, provider JWTProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-HUe3q", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-JKLS2", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-JKs3f", "Errors.Invalid.Argument")
		}
		if provider.JWTEndpoint = strings.TrimSpace(provider.JWTEndpoint); provider.JWTEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-NJKS2", "Errors.Invalid.Argument")
		}
		if provider.KeyEndpoint = strings.TrimSpace(provider.KeyEndpoint); provider.KeyEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-SJk2d", "Errors.Invalid.Argument")
		}
		if provider.HeaderName = strings.TrimSpace(provider.HeaderName); provider.HeaderName == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-SJK2f", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "ORG-Bhju5", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				provider.Name,
				provider.Issuer,
				provider.JWTEndpoint,
				provider.KeyEndpoint,
				provider.HeaderName,
				provider.IDPOptions,
			)
			if err != nil || event == nil {
				return nil, err
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddOrgGitHubProvider(a *org.Aggregate, resourceOwner, id string, provider GitHubProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Jdsgf", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-dsgz3", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewGitHubOrgIDPWriteModel(resourceOwner, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewGitHubIDPAddedEvent(
					ctx,
					&a.Aggregate,
					id,
					provider.Name,
					provider.ClientID,
					secret,
					provider.Scopes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgGitHubProvider(a *org.Aggregate, resourceOwner, id string, provider GitHubProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if id = strings.TrimSpace(id); id == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-sdf4h", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-fdh5z", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewGitHubOrgIDPWriteModel(resourceOwner, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "ORG-Dr1gs", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				id,
				provider.Name,
				provider.ClientID,
				provider.ClientSecret,
				c.idpConfigEncryption,
				provider.Scopes,
				provider.IDPOptions,
			)
			if err != nil {
				return nil, err
			}
			if event == nil {
				return nil, nil
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddOrgGitHubEnterpriseProvider(a *org.Aggregate, resourceOwner, id string, provider GitHubEnterpriseProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Dg4td", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-dgj53", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Ghjjs", "Errors.Invalid.Argument")
		}
		if provider.AuthorizationEndpoint = strings.TrimSpace(provider.AuthorizationEndpoint); provider.AuthorizationEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-sani2", "Errors.Invalid.Argument")
		}
		if provider.TokenEndpoint = strings.TrimSpace(provider.TokenEndpoint); provider.TokenEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-agj42", "Errors.Invalid.Argument")
		}
		if provider.UserEndpoint = strings.TrimSpace(provider.UserEndpoint); provider.UserEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-sd5hn", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewGitHubEnterpriseOrgIDPWriteModel(resourceOwner, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewGitHubEnterpriseIDPAddedEvent(
					ctx,
					&a.Aggregate,
					id,
					provider.Name,
					provider.ClientID,
					secret,
					provider.AuthorizationEndpoint,
					provider.TokenEndpoint,
					provider.UserEndpoint,
					provider.Scopes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgGitHubEnterpriseProvider(a *org.Aggregate, resourceOwner, id string, provider GitHubEnterpriseProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if id = strings.TrimSpace(id); id == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-sdfh3", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-shj42", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-sdh73", "Errors.Invalid.Argument")
		}
		if provider.AuthorizationEndpoint = strings.TrimSpace(provider.AuthorizationEndpoint); provider.AuthorizationEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-acx2w", "Errors.Invalid.Argument")
		}
		if provider.TokenEndpoint = strings.TrimSpace(provider.TokenEndpoint); provider.TokenEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-dgj6q", "Errors.Invalid.Argument")
		}
		if provider.UserEndpoint = strings.TrimSpace(provider.UserEndpoint); provider.UserEndpoint == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-ybj62", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewGitHubEnterpriseOrgIDPWriteModel(resourceOwner, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "ORG-GBr42", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				id,
				provider.Name,
				provider.ClientID,
				provider.ClientSecret,
				c.idpConfigEncryption,
				provider.AuthorizationEndpoint,
				provider.TokenEndpoint,
				provider.UserEndpoint,
				provider.Scopes,
				provider.IDPOptions,
			)
			if err != nil {
				return nil, err
			}
			if event == nil {
				return nil, nil
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddOrgGoogleProvider(a *org.Aggregate, writeModel *OrgGoogleIDPWriteModel, provider GoogleProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-D3fvs", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-W2vqs", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewGoogleIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.ClientID,
					secret,
					provider.Scopes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgGoogleProvider(a *org.Aggregate, writeModel *OrgGoogleIDPWriteModel, provider GoogleProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-S32t1", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-ds432", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "ORG-Dqrg1", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				provider.Name,
				provider.ClientID,
				provider.ClientSecret,
				c.idpConfigEncryption,
				provider.Scopes,
				provider.IDPOptions,
			)
			if err != nil || event == nil {
				return nil, err
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddOrgLDAPProvider(a *org.Aggregate, writeModel *OrgLDAPIDPWriteModel, provider LDAPProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-SAfdd", "Errors.Invalid.Argument")
		}
		if provider.Host = strings.TrimSpace(provider.Host); provider.Host == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-SDVg2", "Errors.Invalid.Argument")
		}
		if provider.BaseDN = strings.TrimSpace(provider.BaseDN); provider.BaseDN == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-sv31s", "Errors.Invalid.Argument")
		}
		if provider.UserObjectClass = strings.TrimSpace(provider.UserObjectClass); provider.UserObjectClass == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-sdgf4", "Errors.Invalid.Argument")
		}
		if provider.UserUniqueAttribute = strings.TrimSpace(provider.UserUniqueAttribute); provider.UserUniqueAttribute == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-AEG2w", "Errors.Invalid.Argument")
		}
		if provider.Admin = strings.TrimSpace(provider.Admin); provider.Admin == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-SAD5n", "Errors.Invalid.Argument")
		}
		if provider.Password = strings.TrimSpace(provider.Password); provider.Password == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-sdf5h", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.Password), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewLDAPIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.Host,
					provider.Port,
					provider.TLS,
					provider.BaseDN,
					provider.UserObjectClass,
					provider.UserUniqueAttribute,
					provider.Admin,
					secret,
					provider.LDAPAttributes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgLDAPProvider(a *org.Aggregate, writeModel *OrgLDAPIDPWriteModel, provider LDAPProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Dgdbs", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Sffgd", "Errors.Invalid.Argument")
		}
		if provider.Host = strings.TrimSpace(provider.Host); provider.Host == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-Dz62d", "Errors.Invalid.Argument")
		}
		if provider.BaseDN = strings.TrimSpace(provider.BaseDN); provider.BaseDN == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-vb3ss", "Errors.Invalid.Argument")
		}
		if provider.UserObjectClass = strings.TrimSpace(provider.UserObjectClass); provider.UserObjectClass == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-hbere", "Errors.Invalid.Argument")
		}
		if provider.UserUniqueAttribute = strings.TrimSpace(provider.UserUniqueAttribute); provider.UserUniqueAttribute == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-ASFt6", "Errors.Invalid.Argument")
		}
		if provider.Admin = strings.TrimSpace(provider.Admin); provider.Admin == "" {
			return nil, caos_errs.ThrowInvalidArgument(nil, "ORG-DG45z", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "ORG-ASF3F", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				writeModel.Name,
				provider.Name,
				provider.Host,
				provider.Port,
				provider.TLS,
				provider.BaseDN,
				provider.UserObjectClass,
				provider.UserUniqueAttribute,
				provider.Admin,
				provider.Password,
				c.idpConfigEncryption,
				provider.LDAPAttributes,
				provider.IDPOptions,
			)
			if err != nil || event == nil {
				return nil, err
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareDeleteOrgProvider(a *org.Aggregate, resourceOwner, id string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewOrgIDPRemoveWriteModel(resourceOwner, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, caos_errs.ThrowNotFound(nil, "ORG-Se3tg", "Errors.Org.IDPConfig.NotExisting")
			}
			return []eventstore.Command{org.NewIDPRemovedEvent(ctx, &a.Aggregate, id)}, nil
		}, nil
	}
}

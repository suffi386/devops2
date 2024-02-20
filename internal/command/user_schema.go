package command

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/santhosh-tekuri/jsonschema"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/user/schema"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type CreateUserSchema struct {
	Type                   string
	Schema                 map[string]any
	PossibleAuthenticators []domain.AuthenticatorType
}

func (s *CreateUserSchema) Valid() error {
	if s.Type == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMA-DGFj3", "Errors.UserSchema.Type.Missing") // TODO: i18n
	}
	if err := validateUserSchema(s.Schema); err != nil {
		return err
	}
	for _, authenticator := range s.PossibleAuthenticators {
		if authenticator == domain.AuthenticatorTypeUnspecified {
			return zerrors.ThrowInvalidArgument(nil, "COMMA-Gh652", "Errors.UserSchema.Authenticator.Invalid") // TODO: i18n
		}
	}
	return nil
}

type UpdateUserSchema struct {
	ID                     string
	Type                   *string
	Schema                 map[string]any
	PossibleAuthenticators []domain.AuthenticatorType
}

func (s *UpdateUserSchema) Valid() error {
	if s.ID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMA-H5421", "Errors.IDMissing")
	}
	if s.Type != nil && *s.Type == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMA-G43gn", "Errors.UserSchema.Type.Missing") // TODO: i18n
	}
	if err := validateUserSchema(s.Schema); err != nil {
		return err
	}
	for _, authenticator := range s.PossibleAuthenticators {
		if authenticator == domain.AuthenticatorTypeUnspecified {
			return zerrors.ThrowInvalidArgument(nil, "COMMA-WF4hg", "Errors.UserSchema.Authenticator.Invalid") // TODO: i18n
		}
	}
	return nil
}

func (c *Commands) CreateUserSchema(ctx context.Context, userSchema *CreateUserSchema) (string, *domain.ObjectDetails, error) {
	if err := userSchema.Valid(); err != nil {
		return "", nil, err
	}
	id, err := c.idGenerator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewUserSchemaWriteModel(id, authz.GetInstance(ctx).InstanceID())
	err = c.pushAppendAndReduce(ctx, writeModel,
		schema.NewCreatedEvent(ctx,
			UserSchemaAggregateFromWriteModel(&writeModel.WriteModel),
			userSchema.Type, userSchema.Schema, userSchema.PossibleAuthenticators,
		),
	)
	if err != nil {
		return "", nil, err
	}
	return id, writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) UpdateUserSchema(ctx context.Context, userSchema *UpdateUserSchema) (*domain.ObjectDetails, error) {
	if err := userSchema.Valid(); err != nil {
		return nil, err
	}
	writeModel := NewUserSchemaWriteModel(userSchema.ID, "")
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	if writeModel.State != domain.UserSchemaStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMA-HB3e1", "Errors.UserSchema.NotActive") // TODO: i18n
	}
	updatedEvent := writeModel.NewUpdatedEvent(
		ctx,
		UserSchemaAggregateFromWriteModel(&writeModel.WriteModel),
		userSchema.Type,
		userSchema.Schema,
		userSchema.PossibleAuthenticators,
	)
	if updatedEvent == nil {
		return writeModelToObjectDetails(&writeModel.WriteModel), nil
	}
	if err := c.pushAppendAndReduce(ctx, writeModel, updatedEvent); err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) DeactivateUserSchema(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMA-Vvf3w", "Errors.IDMissing")
	}
	writeModel := NewUserSchemaWriteModel(id, "")
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	if writeModel.State != domain.UserSchemaStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMA-E4t4z", "Errors.UserSchema.NotActive") // TODO: i18n
	}
	err := c.pushAppendAndReduce(ctx, writeModel,
		schema.NewDeactivatedEvent(ctx, UserSchemaAggregateFromWriteModel(&writeModel.WriteModel)),
	)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) ReactivateUserSchema(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMA-wq3Gw", "Errors.IDMissing")
	}
	writeModel := NewUserSchemaWriteModel(id, "")
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	if writeModel.State != domain.UserSchemaStateInactive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMA-DGzh5", "Errors.UserSchema.NotInactive") // TODO: i18n
	}
	err := c.pushAppendAndReduce(ctx, writeModel,
		schema.NewReactivatedEvent(ctx, UserSchemaAggregateFromWriteModel(&writeModel.WriteModel)),
	)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func (c *Commands) DeleteUserSchema(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	if id == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMA-E22gg", "Errors.IDMissing")
	}
	writeModel := NewUserSchemaWriteModel(id, "")
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	if !writeModel.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMA-Grg41", "Errors.UserSchema.NotExists") // TODO: i18n
	}
	// TODO: check for users based on that schema
	err := c.pushAppendAndReduce(ctx, writeModel,
		schema.NewDeletedEvent(ctx, UserSchemaAggregateFromWriteModel(&writeModel.WriteModel), writeModel.SchemaType),
	)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

func validateUserSchema(schema map[string]any) error {
	jsonSchema, err := json.Marshal(schema)
	if err != nil {
		return zerrors.ThrowInvalidArgument(err, "COMMA-SFerg", "Errors.UserSchema.Schema.Invalid")
	}
	c := jsonschema.NewCompiler()
	err = c.AddResource("", bytes.NewReader(jsonSchema))
	if err != nil {
		return zerrors.ThrowInvalidArgument(err, "COMMA-Frh42", "Errors.UserSchema.Schema.Invalid")
	}
	_, err = c.Compile("")
	if err != nil {
		return zerrors.ThrowInvalidArgument(err, "COMMA-W21tg", "Errors.UserSchema.Schema.Invalid")
	}
	return nil
}

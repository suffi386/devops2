package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/errors"
)

func domainPolicyWriteModel(ctx context.Context, filter preparation.FilterToQueryReducer) (*PolicyDomainWriteModel, error) {
	wm, err := orgDomainPolicy(ctx, filter)
	if err != nil {
		return nil, err
	}
	if wm != nil && wm.State.Exists() {
		return &wm.PolicyDomainWriteModel, err
	}
	instanceWriteModel, err := instanceDomainPolicy(ctx, filter)
	if err != nil {
		return nil, err
	}
	if instanceWriteModel != nil && instanceWriteModel.State.Exists() {
		return &instanceWriteModel.PolicyDomainWriteModel, err
	}
	return nil, errors.ThrowInternal(nil, "USER-Ggk9n", "Errors.Internal")
}

func orgDomainPolicy(ctx context.Context, filter preparation.FilterToQueryReducer) (*OrgDomainPolicyWriteModel, error) {
	policy := NewOrgDomainPolicyWriteModel(authz.GetCtxData(ctx).OrgID)
	events, err := filter(ctx, policy.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return policy, nil
	}
	policy.AppendEvents(events...)
	err = policy.Reduce()
	return policy, err
}

func instanceDomainPolicy(ctx context.Context, filter preparation.FilterToQueryReducer) (*InstanceDomainPolicyWriteModel, error) {
	policy := NewInstanceDomainPolicyWriteModel(ctx)
	events, err := filter(ctx, policy.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return policy, nil
	}
	policy.AppendEvents(events...)
	err = policy.Reduce()
	return policy, err
}

func domainPolicyUsernames(ctx context.Context, filter preparation.FilterToQueryReducer, orgID string) (*DomainPolicyUsernamesWriteModel, error) {
	policy := NewDomainPolicyUsernamesWriteModel(orgID)
	events, err := filter(ctx, policy.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return policy, nil
	}
	policy.AppendEvents(events...)
	err = policy.Reduce()
	return policy, err
}

func domainPolicyOrgs(ctx context.Context, filter preparation.FilterToQueryReducer) (*DomainPolicyOrgsWriteModel, error) {
	policy := NewDomainPolicyOrgsWriteModel()
	events, err := filter(ctx, policy.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return policy, nil
	}
	policy.AppendEvents(events...)
	err = policy.Reduce()
	return policy, err
}

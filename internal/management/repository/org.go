package repository

import (
	"context"

	org_model "github.com/caos/zitadel/internal/org/model"
)

type OrgRepository interface {
	OrgByID(ctx context.Context, id string) (*org_model.Org, error)
	OrgByDomainGlobal(ctx context.Context, domain string) (*org_model.OrgView, error)
	UpdateOrg(ctx context.Context, org *org_model.Org) (*org_model.Org, error)
	DeactivateOrg(ctx context.Context, id string) (*org_model.Org, error)
	ReactivateOrg(ctx context.Context, id string) (*org_model.Org, error)

	SearchMyOrgDomains(ctx context.Context, request *org_model.OrgDomainSearchRequest) (*org_model.OrgDomainSearchResponse, error)
	AddMyOrgDomain(ctx context.Context, domain *org_model.OrgDomain) (*org_model.OrgDomain, error)
	RemoveMyOrgDomain(ctx context.Context, domain string) error

	SearchOrgMembers(ctx context.Context, request *org_model.OrgMemberSearchRequest) (*org_model.OrgMemberSearchResponse, error)
	AddOrgMember(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error)
	ChangeOrgMember(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error)
	RemoveOrgMember(ctx context.Context, orgID, userID string) error

	GetOrgMemberRoles() []string
}

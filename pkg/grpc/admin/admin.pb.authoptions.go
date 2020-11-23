// Code generated by protoc-gen-authmethod. DO NOT EDIT.

package admin

import (
	"github.com/caos/zitadel/internal/api/authz"
)

/**
 * AdminService
 */

const AdminService_MethodPrefix = "caos.zitadel.admin.api.v1.AdminService"

var AdminService_AuthMethods = authz.MethodMapping{

	"/caos.zitadel.admin.api.v1.AdminService/IsOrgUnique": authz.Option{
		Permission: "iam.read",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/GetOrgByID": authz.Option{
		Permission: "iam.read",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/SearchOrgs": authz.Option{
		Permission: "iam.read",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/SetUpOrg": authz.Option{
		Permission: "iam.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/GetDefaultOrgIamPolicy": authz.Option{
		Permission: "iam.policy.read",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/UpdateDefaultOrgIamPolicy": authz.Option{
		Permission: "iam.policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/GetOrgIamPolicy": authz.Option{
		Permission: "iam.policy.read",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/CreateOrgIamPolicy": authz.Option{
		Permission: "iam.policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/UpdateOrgIamPolicy": authz.Option{
		Permission: "iam.policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/RemoveOrgIamPolicy": authz.Option{
		Permission: "iam.policy.delete",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/GetIamMemberRoles": authz.Option{
		Permission: "iam.member.read",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/AddIamMember": authz.Option{
		Permission: "iam.member.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/ChangeIamMember": authz.Option{
		Permission: "iam.member.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/RemoveIamMember": authz.Option{
		Permission: "iam.member.delete",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/SearchIamMembers": authz.Option{
		Permission: "iam.member.read",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/GetViews": authz.Option{
		Permission: "iam.read",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/ClearView": authz.Option{
		Permission: "iam.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/GetFailedEvents": authz.Option{
		Permission: "iam.read",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/RemoveFailedEvent": authz.Option{
		Permission: "iam.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/IdpByID": authz.Option{
		Permission: "iam.idp.read",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/CreateOidcIdp": authz.Option{
		Permission: "iam.idp.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/UpdateIdpConfig": authz.Option{
		Permission: "iam.idp.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/DeactivateIdpConfig": authz.Option{
		Permission: "iam.idp.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/ReactivateIdpConfig": authz.Option{
		Permission: "iam.idp.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/RemoveIdpConfig": authz.Option{
		Permission: "iam.idp.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/UpdateOidcIdpConfig": authz.Option{
		Permission: "iam.idp.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/SearchIdps": authz.Option{
		Permission: "iam.idp.read",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/GetDefaultLabelPolicy": authz.Option{
		Permission: "iam.policy.read",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/UpdateDefaultLabelPolicy": authz.Option{
		Permission: "iam.policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/GetDefaultMailTemplate": authz.Option{
		Permission: "iam.policy.read",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/UpdateDefaultMailTemplate": authz.Option{
		Permission: "iam.policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/GetDefaultMailTexts": authz.Option{
		Permission: "iam.policy.read",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/UpdateDefaultMailText": authz.Option{
		Permission: "iam.policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/GetDefaultLoginPolicy": authz.Option{
		Permission: "iam.policy.read",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/UpdateDefaultLoginPolicy": authz.Option{
		Permission: "iam.policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/GetDefaultLoginPolicyIdpProviders": authz.Option{
		Permission: "iam.policy.read",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/AddIdpProviderToDefaultLoginPolicy": authz.Option{
		Permission: "iam.policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/RemoveIdpProviderFromDefaultLoginPolicy": authz.Option{
		Permission: "iam.policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/GetDefaultPasswordComplexityPolicy": authz.Option{
		Permission: "iam.policy.read",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/UpdateDefaultPasswordComplexityPolicy": authz.Option{
		Permission: "iam.policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/GetDefaultPasswordAgePolicy": authz.Option{
		Permission: "iam.policy.read",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/UpdateDefaultPasswordAgePolicy": authz.Option{
		Permission: "iam.policy.write",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/GetDefaultPasswordLockoutPolicy": authz.Option{
		Permission: "iam.policy.read",
		CheckParam: "",
	},

	"/caos.zitadel.admin.api.v1.AdminService/UpdateDefaultPasswordLockoutPolicy": authz.Option{
		Permission: "iam.policy.write",
		CheckParam: "",
	},
}

package management

import (
	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes"
)

func loginPolicyRequestToModel(policy *management.LoginPolicyRequest) *iam_model.LoginPolicy {
	return &iam_model.LoginPolicy{
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowExternalIdp:      policy.AllowExternalIdp,
		AllowRegister:         policy.AllowRegister,
		ForceMFA:              policy.ForceMfa,
	}
}

func loginPolicyFromModel(policy *iam_model.LoginPolicy) *management.LoginPolicy {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-2Fsm8").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-3Flo0").OnError(err).Debug("date parse failed")

	return &management.LoginPolicy{
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowExternalIdp:      policy.AllowExternalIdp,
		AllowRegister:         policy.AllowRegister,
		CreationDate:          creationDate,
		ChangeDate:            changeDate,
		ForceMfa:              policy.ForceMFA,
	}
}

func loginPolicyViewFromModel(policy *iam_model.LoginPolicyView) *management.LoginPolicyView {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-5Tsm8").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-8dJgs").OnError(err).Debug("date parse failed")

	return &management.LoginPolicyView{
		Default:               policy.Default,
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowExternalIdp:      policy.AllowExternalIDP,
		AllowRegister:         policy.AllowRegister,
		CreationDate:          creationDate,
		ChangeDate:            changeDate,
		ForceMfa:              policy.ForceMFA,
	}
}

func idpProviderSearchRequestToModel(request *management.IdpProviderSearchRequest) *iam_model.IDPProviderSearchRequest {
	return &iam_model.IDPProviderSearchRequest{
		Limit:  request.Limit,
		Offset: request.Offset,
	}
}

func idpProviderSearchResponseFromModel(response *iam_model.IDPProviderSearchResponse) *management.IdpProviderSearchResponse {
	return &management.IdpProviderSearchResponse{
		Limit:       response.Limit,
		Offset:      response.Offset,
		TotalResult: response.TotalResult,
		Result:      idpProviderViewsFromModel(response.Result),
	}
}

func idpProviderToModel(provider *management.IdpProviderID) *iam_model.IDPProvider {
	return &iam_model.IDPProvider{
		IdpConfigID: provider.IdpConfigId,
		Type:        iam_model.IDPProviderTypeSystem,
	}
}

func idpProviderAddToModel(provider *management.IdpProviderAdd) *iam_model.IDPProvider {
	return &iam_model.IDPProvider{
		IdpConfigID: provider.IdpConfigId,
		Type:        idpProviderTypeToModel(provider.IdpProviderType),
	}
}

func idpProviderIDFromModel(provider *iam_model.IDPProvider) *management.IdpProviderID {
	return &management.IdpProviderID{
		IdpConfigId: provider.IdpConfigID,
	}
}

func idpProviderFromModel(provider *iam_model.IDPProvider) *management.IdpProvider {
	return &management.IdpProvider{
		IdpConfigId:      provider.IdpConfigID,
		IdpProvider_Type: idpProviderTypeFromModel(provider.Type),
	}
}

func idpProviderViewsFromModel(providers []*iam_model.IDPProviderView) []*management.IdpProviderView {
	converted := make([]*management.IdpProviderView, len(providers))
	for i, provider := range providers {
		converted[i] = idpProviderViewFromModel(provider)
	}

	return converted
}

func idpProviderViewFromModel(provider *iam_model.IDPProviderView) *management.IdpProviderView {
	return &management.IdpProviderView{
		IdpConfigId: provider.IDPConfigID,
		Name:        provider.Name,
		Type:        idpConfigTypeToModel(provider.IDPConfigType),
	}
}

func idpConfigTypeToModel(providerType iam_model.IdpConfigType) management.IdpType {
	switch providerType {
	case iam_model.IDPConfigTypeOIDC:
		return management.IdpType_IDPTYPE_OIDC
	case iam_model.IDPConfigTypeSAML:
		return management.IdpType_IDPTYPE_SAML
	default:
		return management.IdpType_IDPTYPE_UNSPECIFIED
	}
}

func idpProviderTypeToModel(providerType management.IdpProviderType) iam_model.IDPProviderType {
	switch providerType {
	case management.IdpProviderType_IDPPROVIDERTYPE_SYSTEM:
		return iam_model.IDPProviderTypeSystem
	case management.IdpProviderType_IDPPROVIDERTYPE_ORG:
		return iam_model.IDPProviderTypeOrg
	default:
		return iam_model.IDPProviderTypeSystem
	}
}

func idpProviderTypeFromModel(providerType iam_model.IDPProviderType) management.IdpProviderType {
	switch providerType {
	case iam_model.IDPProviderTypeSystem:
		return management.IdpProviderType_IDPPROVIDERTYPE_SYSTEM
	case iam_model.IDPProviderTypeOrg:
		return management.IdpProviderType_IDPPROVIDERTYPE_ORG
	default:
		return management.IdpProviderType_IDPPROVIDERTYPE_UNSPECIFIED
	}
}

func softwareMFAResultFromModel(result *iam_model.SoftwareMFASearchResponse) *management.SoftwareMFAResult {
	converted := make([]management.SoftwareMFAType, len(result.Result))
	for i, mfaType := range result.Result {
		converted[i] = softwareMFATypeFromModel(mfaType)
	}
	return &management.SoftwareMFAResult{
		Mfas: converted,
	}
}

func softwareMFAFromModel(mfaType iam_model.SoftwareMFAType) *management.SoftwareMFA {
	return &management.SoftwareMFA{
		Mfa: softwareMFATypeFromModel(mfaType),
	}
}

func softwareMFATypeFromModel(mfaType iam_model.SoftwareMFAType) management.SoftwareMFAType {
	switch mfaType {
	case iam_model.SoftwareMFATypeOTP:
		return management.SoftwareMFAType_SOFTWAREMFATYPE_OTP
	default:
		return management.SoftwareMFAType_SOFTWAREMFATYPE_UNSPECIFIED
	}
}

func softwareMFATypeToModel(mfaType *management.SoftwareMFA) iam_model.SoftwareMFAType {
	switch mfaType.Mfa {
	case management.SoftwareMFAType_SOFTWAREMFATYPE_OTP:
		return iam_model.SoftwareMFATypeOTP
	default:
		return iam_model.SoftwareMFATypeUnspecified
	}
}

func hardwareMFAResultFromModel(result *iam_model.HardwareMFASearchResponse) *management.HardwareMFAResult {
	converted := make([]management.HardwareMFAType, len(result.Result))
	for i, mfaType := range result.Result {
		converted[i] = hardwareMFATypeFromModel(mfaType)
	}
	return &management.HardwareMFAResult{
		Mfas: converted,
	}
}

func hardwareMFAFromModel(mfaType iam_model.HardwareMFAType) *management.HardwareMFA {
	return &management.HardwareMFA{
		Mfa: hardwareMFATypeFromModel(mfaType),
	}
}

func hardwareMFATypeFromModel(mfaType iam_model.HardwareMFAType) management.HardwareMFAType {
	switch mfaType {
	case iam_model.HardwareMFATypeU2F:
		return management.HardwareMFAType_HARDWAREMFATYPE_U2F
	default:
		return management.HardwareMFAType_HARDWAREMFATYPE_UNSPECIFIED
	}
}

func hardwareMFATypeToModel(mfaType *management.HardwareMFA) iam_model.HardwareMFAType {
	switch mfaType.Mfa {
	case management.HardwareMFAType_HARDWAREMFATYPE_U2F:
		return iam_model.HardwareMFATypeU2F
	default:
		return iam_model.HardwareMFATypeUnspecified
	}
}

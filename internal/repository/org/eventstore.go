package org

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, OrgAddedEventType, OrgAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, OrgChangedEventType, OrgChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, OrgDeactivatedEventType, OrgDeactivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, OrgReactivatedEventType, OrgReactivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, OrgRemovedEventType, OrgRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, OrgDomainAddedEventType, DomainAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, OrgDomainVerificationAddedEventType, DomainVerificationAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, OrgDomainVerificationFailedEventType, DomainVerificationFailedEventMapper).
		RegisterFilterEventMapper(AggregateType, OrgDomainVerifiedEventType, DomainVerifiedEventMapper).
		RegisterFilterEventMapper(AggregateType, OrgDomainPrimarySetEventType, DomainPrimarySetEventMapper).
		RegisterFilterEventMapper(AggregateType, OrgDomainRemovedEventType, DomainRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, MemberAddedEventType, MemberAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, MemberChangedEventType, MemberChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, MemberRemovedEventType, MemberRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, MemberCascadeRemovedEventType, MemberCascadeRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyAddedEventType, LabelPolicyAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyChangedEventType, LabelPolicyChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyActivatedEventType, LabelPolicyActivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyRemovedEventType, LabelPolicyRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyLogoAddedEventType, LabelPolicyLogoAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyLogoRemovedEventType, LabelPolicyLogoRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyIconAddedEventType, LabelPolicyIconAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyIconRemovedEventType, LabelPolicyIconRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyLogoDarkAddedEventType, LabelPolicyLogoDarkAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyLogoDarkRemovedEventType, LabelPolicyLogoDarkRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyIconDarkAddedEventType, LabelPolicyIconDarkAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyIconDarkRemovedEventType, LabelPolicyIconDarkRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyFontAddedEventType, LabelPolicyFontAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyFontRemovedEventType, LabelPolicyFontRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LabelPolicyAssetsRemovedEventType, LabelPolicyAssetsRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LoginPolicyAddedEventType, LoginPolicyAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LoginPolicyChangedEventType, LoginPolicyChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, LoginPolicyRemovedEventType, LoginPolicyRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LoginPolicySecondFactorAddedEventType, SecondFactorAddedEventEventMapper).
		RegisterFilterEventMapper(AggregateType, LoginPolicySecondFactorRemovedEventType, SecondFactorRemovedEventEventMapper).
		RegisterFilterEventMapper(AggregateType, LoginPolicyMultiFactorAddedEventType, MultiFactorAddedEventEventMapper).
		RegisterFilterEventMapper(AggregateType, LoginPolicyMultiFactorRemovedEventType, MultiFactorRemovedEventEventMapper).
		RegisterFilterEventMapper(AggregateType, LoginPolicyIDPProviderAddedEventType, IdentityProviderAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LoginPolicyIDPProviderRemovedEventType, IdentityProviderRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LoginPolicyIDPProviderCascadeRemovedEventType, IdentityProviderCascadeRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, DomainPolicyAddedEventType, DomainPolicyAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, DomainPolicyChangedEventType, DomainPolicyChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, DomainPolicyRemovedEventType, DomainPolicyRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, PasswordAgePolicyAddedEventType, PasswordAgePolicyAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, PasswordAgePolicyChangedEventType, PasswordAgePolicyChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, PasswordAgePolicyRemovedEventType, PasswordAgePolicyRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, PasswordComplexityPolicyAddedEventType, PasswordComplexityPolicyAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, PasswordComplexityPolicyChangedEventType, PasswordComplexityPolicyChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, PasswordComplexityPolicyRemovedEventType, PasswordComplexityPolicyRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, LockoutPolicyAddedEventType, LockoutPolicyAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LockoutPolicyChangedEventType, LockoutPolicyChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, LockoutPolicyRemovedEventType, LockoutPolicyRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, PrivacyPolicyAddedEventType, PrivacyPolicyAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, PrivacyPolicyChangedEventType, PrivacyPolicyChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, PrivacyPolicyRemovedEventType, PrivacyPolicyRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, MailTemplateAddedEventType, MailTemplateAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, MailTemplateChangedEventType, MailTemplateChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, MailTemplateRemovedEventType, MailTemplateRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, MailTextAddedEventType, MailTextAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, MailTextChangedEventType, MailTextChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, MailTextRemovedEventType, MailTextRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, CustomTextSetEventType, CustomTextSetEventMapper).
		RegisterFilterEventMapper(AggregateType, CustomTextRemovedEventType, CustomTextRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, CustomTextTemplateRemovedEventType, CustomTextTemplateRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPConfigAddedEventType, IDPConfigAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPConfigChangedEventType, IDPConfigChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPConfigRemovedEventType, IDPConfigRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPConfigDeactivatedEventType, IDPConfigDeactivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPConfigReactivatedEventType, IDPConfigReactivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPOIDCConfigAddedEventType, IDPOIDCConfigAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPOIDCConfigChangedEventType, IDPOIDCConfigChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPJWTConfigAddedEventType, IDPJWTConfigAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPJWTConfigChangedEventType, IDPJWTConfigChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, LDAPIDPAddedEventType, LDAPIDPAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, LDAPIDPChangedEventType, LDAPIDPChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, IDPRemovedEventType, IDPRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, TriggerActionsSetEventType, TriggerActionsSetEventMapper).
		RegisterFilterEventMapper(AggregateType, TriggerActionsCascadeRemovedEventType, TriggerActionsCascadeRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, FlowClearedEventType, FlowClearedEventMapper).
		RegisterFilterEventMapper(AggregateType, MetadataSetType, MetadataSetEventMapper).
		RegisterFilterEventMapper(AggregateType, MetadataRemovedType, MetadataRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, MetadataRemovedAllType, MetadataRemovedAllEventMapper).
		RegisterFilterEventMapper(AggregateType, NotificationPolicyAddedEventType, NotificationPolicyAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, NotificationPolicyChangedEventType, NotificationPolicyChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, NotificationPolicyRemovedEventType, NotificationPolicyRemovedEventMapper)
}

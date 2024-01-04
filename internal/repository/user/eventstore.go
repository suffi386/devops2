package user

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func RegisterEventMappers(es *eventstore.Eventstore) {
	es.RegisterFilterEventMapper(AggregateType, UserV1AddedType, HumanAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1RegisteredType, HumanRegisteredEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1InitialCodeAddedType, HumanInitialCodeAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1InitialCodeSentType, HumanInitialCodeSentEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1InitializedCheckSucceededType, HumanInitializedCheckSucceededEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1InitializedCheckFailedType, HumanInitializedCheckFailedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1SignedOutType, HumanSignedOutEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1PasswordChangedType, HumanPasswordChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1PasswordCodeAddedType, HumanPasswordCodeAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1PasswordCodeSentType, HumanPasswordCodeSentEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1PasswordCheckSucceededType, HumanPasswordCheckSucceededEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1PasswordCheckFailedType, HumanPasswordCheckFailedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1EmailChangedType, HumanEmailChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1EmailVerifiedType, HumanEmailVerifiedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1EmailVerificationFailedType, HumanEmailVerificationFailedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1EmailCodeAddedType, HumanEmailCodeAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1EmailCodeSentType, HumanEmailCodeSentEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1PhoneChangedType, HumanPhoneChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1PhoneRemovedType, HumanPhoneRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1PhoneVerifiedType, HumanPhoneVerifiedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1PhoneVerificationFailedType, HumanPhoneVerificationFailedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1PhoneCodeAddedType, HumanPhoneCodeAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1PhoneCodeSentType, HumanPhoneCodeSentEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1ProfileChangedType, HumanProfileChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1AddressChangedType, HumanAddressChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1MFAInitSkippedType, HumanMFAInitSkippedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1MFAOTPAddedType, HumanOTPAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1MFAOTPVerifiedType, HumanOTPVerifiedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1MFAOTPRemovedType, HumanOTPRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1MFAOTPCheckSucceededType, HumanOTPCheckSucceededEventMapper).
		RegisterFilterEventMapper(AggregateType, UserV1MFAOTPCheckFailedType, HumanOTPCheckFailedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserLockedType, UserLockedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserUnlockedType, UserUnlockedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserDeactivatedType, UserDeactivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserReactivatedType, UserReactivatedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserRemovedType, UserRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserTokenAddedType, UserTokenAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserTokenRemovedType, UserTokenRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserDomainClaimedType, DomainClaimedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserDomainClaimedSentType, DomainClaimedSentEventMapper).
		RegisterFilterEventMapper(AggregateType, UserUserNameChangedType, UsernameChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, MetadataSetType, MetadataSetEventMapper).
		RegisterFilterEventMapper(AggregateType, MetadataRemovedType, MetadataRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, MetadataRemovedAllType, MetadataRemovedAllEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanAddedType, HumanAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanRegisteredType, HumanRegisteredEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanInitialCodeAddedType, HumanInitialCodeAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanInitialCodeSentType, HumanInitialCodeSentEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanInitializedCheckSucceededType, HumanInitializedCheckSucceededEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanInitializedCheckFailedType, HumanInitializedCheckFailedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanSignedOutType, HumanSignedOutEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPasswordChangedType, HumanPasswordChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPasswordCodeAddedType, HumanPasswordCodeAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPasswordCodeSentType, HumanPasswordCodeSentEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPasswordChangeSentType, HumanPasswordChangeSentEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPasswordCheckSucceededType, HumanPasswordCheckSucceededEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPasswordCheckFailedType, HumanPasswordCheckFailedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPasswordHashUpdatedType, eventstore.GenericEventMapper[HumanPasswordHashUpdatedEvent]).
		RegisterFilterEventMapper(AggregateType, UserIDPLinkAddedType, UserIDPLinkAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserIDPLinkRemovedType, UserIDPLinkRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserIDPLinkCascadeRemovedType, UserIDPLinkCascadeRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, UserIDPLoginCheckSucceededType, UserIDPCheckSucceededEventMapper).
		RegisterFilterEventMapper(AggregateType, UserIDPExternalIDMigratedType, eventstore.GenericEventMapper[UserIDPExternalIDMigratedEvent]).
		RegisterFilterEventMapper(AggregateType, UserIDPExternalUsernameChangedType, eventstore.GenericEventMapper[UserIDPExternalUsernameEvent]).
		RegisterFilterEventMapper(AggregateType, HumanEmailChangedType, HumanEmailChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanEmailVerifiedType, HumanEmailVerifiedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanEmailVerificationFailedType, HumanEmailVerificationFailedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanEmailCodeAddedType, HumanEmailCodeAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanEmailCodeSentType, HumanEmailCodeSentEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPhoneChangedType, HumanPhoneChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPhoneRemovedType, HumanPhoneRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPhoneVerifiedType, HumanPhoneVerifiedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPhoneVerificationFailedType, HumanPhoneVerificationFailedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPhoneCodeAddedType, HumanPhoneCodeAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPhoneCodeSentType, HumanPhoneCodeSentEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanProfileChangedType, HumanProfileChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanAvatarAddedType, HumanAvatarAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanAvatarRemovedType, HumanAvatarRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanAddressChangedType, HumanAddressChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanMFAInitSkippedType, HumanMFAInitSkippedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanMFAOTPAddedType, HumanOTPAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanMFAOTPVerifiedType, HumanOTPVerifiedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanMFAOTPRemovedType, HumanOTPRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanMFAOTPCheckSucceededType, HumanOTPCheckSucceededEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanMFAOTPCheckFailedType, HumanOTPCheckFailedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanOTPSMSAddedType, eventstore.GenericEventMapper[HumanOTPSMSAddedEvent]).
		RegisterFilterEventMapper(AggregateType, HumanOTPSMSRemovedType, eventstore.GenericEventMapper[HumanOTPSMSRemovedEvent]).
		RegisterFilterEventMapper(AggregateType, HumanOTPSMSCodeAddedType, eventstore.GenericEventMapper[HumanOTPSMSCodeAddedEvent]).
		RegisterFilterEventMapper(AggregateType, HumanOTPSMSCodeSentType, eventstore.GenericEventMapper[HumanOTPSMSCodeSentEvent]).
		RegisterFilterEventMapper(AggregateType, HumanOTPSMSCheckSucceededType, eventstore.GenericEventMapper[HumanOTPSMSCheckSucceededEvent]).
		RegisterFilterEventMapper(AggregateType, HumanOTPSMSCheckFailedType, eventstore.GenericEventMapper[HumanOTPSMSCheckFailedEvent]).
		RegisterFilterEventMapper(AggregateType, HumanOTPEmailAddedType, eventstore.GenericEventMapper[HumanOTPEmailAddedEvent]).
		RegisterFilterEventMapper(AggregateType, HumanOTPEmailRemovedType, eventstore.GenericEventMapper[HumanOTPEmailRemovedEvent]).
		RegisterFilterEventMapper(AggregateType, HumanOTPEmailCodeAddedType, eventstore.GenericEventMapper[HumanOTPEmailCodeAddedEvent]).
		RegisterFilterEventMapper(AggregateType, HumanOTPEmailCodeSentType, eventstore.GenericEventMapper[HumanOTPEmailCodeSentEvent]).
		RegisterFilterEventMapper(AggregateType, HumanOTPEmailCheckSucceededType, eventstore.GenericEventMapper[HumanOTPEmailCheckSucceededEvent]).
		RegisterFilterEventMapper(AggregateType, HumanOTPEmailCheckFailedType, eventstore.GenericEventMapper[HumanOTPEmailCheckFailedEvent]).
		RegisterFilterEventMapper(AggregateType, HumanU2FTokenAddedType, HumanU2FAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanU2FTokenVerifiedType, HumanU2FVerifiedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanU2FTokenSignCountChangedType, HumanU2FSignCountChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanU2FTokenRemovedType, HumanU2FRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanU2FTokenBeginLoginType, HumanU2FBeginLoginEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanU2FTokenCheckSucceededType, HumanU2FCheckSucceededEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanU2FTokenCheckFailedType, HumanU2FCheckFailedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPasswordlessTokenAddedType, HumanPasswordlessAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPasswordlessTokenVerifiedType, HumanPasswordlessVerifiedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPasswordlessTokenSignCountChangedType, HumanPasswordlessSignCountChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPasswordlessTokenRemovedType, HumanPasswordlessRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPasswordlessTokenBeginLoginType, HumanPasswordlessBeginLoginEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPasswordlessTokenCheckSucceededType, HumanPasswordlessCheckSucceededEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPasswordlessTokenCheckFailedType, HumanPasswordlessCheckFailedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPasswordlessInitCodeAddedType, HumanPasswordlessInitCodeAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPasswordlessInitCodeRequestedType, HumanPasswordlessInitCodeRequestedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPasswordlessInitCodeSentType, HumanPasswordlessInitCodeSentEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPasswordlessInitCodeCheckFailedType, HumanPasswordlessInitCodeCodeCheckFailedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanPasswordlessInitCodeCheckSucceededType, HumanPasswordlessInitCodeCodeCheckSucceededEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanRefreshTokenAddedType, HumanRefreshTokenAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanRefreshTokenRenewedType, HumanRefreshTokenRenewedEventEventMapper).
		RegisterFilterEventMapper(AggregateType, HumanRefreshTokenRemovedType, HumanRefreshTokenRemovedEventEventMapper).
		RegisterFilterEventMapper(AggregateType, MachineAddedEventType, MachineAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, MachineChangedEventType, MachineChangedEventMapper).
		RegisterFilterEventMapper(AggregateType, MachineKeyAddedEventType, MachineKeyAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, MachineKeyRemovedEventType, MachineKeyRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, PersonalAccessTokenAddedType, PersonalAccessTokenAddedEventMapper).
		RegisterFilterEventMapper(AggregateType, PersonalAccessTokenRemovedType, PersonalAccessTokenRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, MachineSecretSetType, MachineSecretSetEventMapper).
		RegisterFilterEventMapper(AggregateType, MachineSecretRemovedType, MachineSecretRemovedEventMapper).
		RegisterFilterEventMapper(AggregateType, MachineSecretCheckSucceededType, MachineSecretCheckSucceededEventMapper).
		RegisterFilterEventMapper(AggregateType, MachineSecretCheckFailedType, MachineSecretCheckFailedEventMapper)
}
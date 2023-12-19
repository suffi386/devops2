package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddUserHuman(t *testing.T) {
	type fields struct {
		eventstore         func(t *testing.T) *eventstore.Eventstore
		idGenerator        id.Generator
		userPasswordHasher *crypto.PasswordHasher
		newCode            cryptoCodeFunc
	}
	type args struct {
		ctx             context.Context
		orgID           string
		human           *AddHuman
		secretGenerator crypto.Generator
		allowInitMail   bool
		codeAlg         crypto.EncryptionAlgorithm
	}
	type res struct {
		want          *domain.ObjectDetails
		wantID        string
		wantEmailCode string
		err           func(error) bool
	}

	userAgg := user.NewAggregate("user1", "org1")

	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "orgid missing, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "",
				human: &AddHuman{
					Username:  "username",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@test.ch",
					},
				},
				allowInitMail: true,
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMA-095xh8fll1", "Errors.Internal"))
				},
			},
		},
		{
			name: "user invalid, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					FirstName: "firstname",
				},
				allowInitMail: true,
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "EMAIL-spblu", "Errors.User.Email.Empty"))
				},
			},
		},
		{
			name: "with id, already exists, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					ID:        "user1",
					Username:  "username",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@test.ch",
					},
					PreferredLanguage: language.English,
				},
				allowInitMail: true,
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-7yiox1isql", "Errors.User.AlreadyExisting"))
				},
			},
		},
		{
			name: "domain policy not found, precondition error",
			fields: fields{
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(),
					expectFilter(),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@test.ch",
					},
					PreferredLanguage: language.English,
				},
				allowInitMail: true,
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInternal(nil, "USER-Ggk9n", "Errors.Internal"))
				},
			},
		},
		{
			name: "password policy not found, precondition error",
			fields: fields{
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(),
					expectFilter(),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					FirstName: "firstname",
					LastName:  "lastname",
					Password:  "pass",
					Email: Email{
						Address:  "email@test.ch",
						Verified: true,
					},
					PreferredLanguage: language.English,
				},
				allowInitMail: true,
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInternal(nil, "USER-uQ96e", "Errors.Internal"))
				},
			},
		},
		{
			name: "add human (with initial code), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						user.NewHumanAddedEvent(context.Background(),
							&userAgg.Aggregate,
							"username",
							"firstname",
							"lastname",
							"",
							"firstname lastname",
							language.English,
							domain.GenderUnspecified,
							"email@test.ch",
							true,
						),
						user.NewHumanInitialCodeAddedEvent(context.Background(),
							&userAgg.Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("userinit"),
							},
							time.Hour*1,
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				newCode:     mockCode("userinit", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@test.ch",
					},
					PreferredLanguage: language.English,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human (with password and initial code), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						newAddHumanEvent("$plain$x$password", false, true, "", language.English),
						user.NewHumanInitialCodeAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("userinit"),
							},
							1*time.Hour,
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
				newCode:            mockCode("userinit", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					Password:  "password",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@test.ch",
					},
					PreferredLanguage: language.English,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human (with password and email code custom template), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						newAddHumanEvent("$plain$x$password", false, true, "", language.English),
						user.NewHumanEmailCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("emailCode"),
							},
							1*time.Hour,
							"https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}",
							false,
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
				newCode:            mockCode("emailCode", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					Password:  "password",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address:     "email@test.ch",
						URLTemplate: "https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}",
					},
					PreferredLanguage: language.English,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   false,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human (with password and return email code), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						newAddHumanEvent("$plain$x$password", false, true, "", language.English),
						user.NewHumanEmailCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("emailCode"),
							},
							1*time.Hour,
							"",
							true,
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
				newCode:            mockCode("emailCode", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					Password:  "password",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address:    "email@test.ch",
						ReturnCode: true,
					},
					PreferredLanguage: language.English,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   false,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID:        "user1",
				wantEmailCode: "emailCode",
			},
		},
		{
			name: "add human email verified, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					Password:  "password",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address:  "email@test.ch",
						Verified: true,
					},
					PreferredLanguage:      language.English,
					PasswordChangeRequired: true,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human email verified, trim spaces, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  " username ",
					Password:  "password",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address:  "email@test.ch",
						Verified: true,
					},
					PreferredLanguage:      language.English,
					PasswordChangeRequired: true,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human, email verified, userLoginMustBeDomain false, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								false,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						newAddHumanEvent("$plain$x$password", true, false, "", language.English),
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					Password:  "password",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address:  "email@test.ch",
						Verified: true,
					},
					PreferredLanguage:      language.English,
					PasswordChangeRequired: true,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human claimed domain, userLoginMustBeDomain false, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								false,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(context.Background(),
								&org.NewAggregate("org2").Aggregate,
								"test.ch",
							),
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username@test.ch",
					Password:  "password",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address:  "email@test.ch",
						Verified: true,
					},
					PreferredLanguage:      language.English,
					PasswordChangeRequired: true,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-SFd21", "Errors.User.DomainNotAllowedAsUsername"))
				},
			},
		},
		{
			name: "add human domain, userLoginMustBeDomain false, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								false,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"test.ch",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						func() eventstore.Command {
							event := user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username@test.ch",
								"firstname",
								"lastname",
								"",
								"firstname lastname",
								language.English,
								domain.GenderUnspecified,
								"email@test.ch",
								false,
							)
							event.AddPasswordData("$plain$x$password", true)
							return event
						}(),
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username@test.ch",
					Password:  "password",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address:  "email@test.ch",
						Verified: true,
					},
					PreferredLanguage:      language.English,
					PasswordChangeRequired: true,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},

			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human (with phone), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						newAddHumanEvent("$plain$x$password", false, true, "+41711234567", language.English),
						user.NewHumanEmailVerifiedEvent(
							context.Background(),
							&userAgg.Aggregate,
						),
						user.NewHumanPhoneCodeAddedEvent(context.Background(),
							&userAgg.Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("phonecode"),
							},
							time.Hour*1,
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
				newCode:            mockCode("phonecode", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					FirstName: "firstname",
					LastName:  "lastname",
					Password:  "password",
					Email: Email{
						Address:  "email@test.ch",
						Verified: true,
					},
					Phone: Phone{
						Number: "+41711234567",
					},
					PreferredLanguage: language.English,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human (with verified phone), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						newAddHumanEvent("", false, true, "+41711234567", language.English),
						user.NewHumanInitialCodeAddedEvent(
							context.Background(),
							&userAgg.Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("userinit"),
							},
							1*time.Hour,
						),
						user.NewHumanPhoneVerifiedEvent(
							context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				newCode:     mockCode("userinit", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@test.ch",
					},
					Phone: Phone{
						Number:   "+41711234567",
						Verified: true,
					},
					PreferredLanguage: language.English,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		}, {
			name: "add human (with return code), ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						newAddHumanEvent("$plain$x$password", false, true, "+41711234567", language.English),
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate),
						user.NewHumanPhoneCodeAddedEventV2(
							context.Background(),
							&userAgg.Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("phoneCode"),
							},
							1*time.Hour,
							true,
						),
					),
				),
				idGenerator:        id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				userPasswordHasher: mockPasswordHasher("x"),
				newCode:            mockCode("phoneCode", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					Password:  "password",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address:  "email@test.ch",
						Verified: true,
					},
					Phone: Phone{
						Number:     "+41711234567",
						ReturnCode: true,
					},
					PreferredLanguage: language.English,
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
		{
			name: "add human with metadata, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						newAddHumanEvent("", false, true, "", language.English),
						user.NewHumanInitialCodeAddedEvent(
							context.Background(),
							&userAgg.Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("userinit"),
							},
							1*time.Hour,
						),
						user.NewMetadataSetEvent(
							context.Background(),
							&userAgg.Aggregate,
							"testKey",
							[]byte("testValue"),
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "user1"),
				newCode:     mockCode("userinit", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &AddHuman{
					Username:  "username",
					FirstName: "firstname",
					LastName:  "lastname",
					Email: Email{
						Address: "email@test.ch",
					},
					PreferredLanguage: language.English,
					Metadata: []*AddMetadataEntry{
						{
							Key:   "testKey",
							Value: []byte("testValue"),
						},
					},
				},
				secretGenerator: GetMockSecretGenerator(t),
				allowInitMail:   true,
				codeAlg:         crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				wantID: "user1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:         tt.fields.eventstore(t),
				userPasswordHasher: tt.fields.userPasswordHasher,
				idGenerator:        tt.fields.idGenerator,
				newCode:            tt.fields.newCode,
			}
			err := r.AddUserHuman(tt.args.ctx, tt.args.orgID, tt.args.human, tt.args.allowInitMail, tt.args.codeAlg)
			if tt.res.err == nil {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			} else if !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, tt.args.human.Details)
				assert.Equal(t, tt.res.wantID, tt.args.human.ID)
				assert.Equal(t, tt.res.wantEmailCode, gu.Value(tt.args.human.EmailCode))
			}
		})
	}
}

func TestCommandSide_ChangeUserHuman(t *testing.T) {
	type fields struct {
		eventstore         func(t *testing.T) *eventstore.Eventstore
		userPasswordHasher *crypto.PasswordHasher
		newCode            cryptoCodeFunc
		checkPermission    domain.PermissionCheck
	}
	type args struct {
		ctx     context.Context
		orgID   string
		human   *ChangeHuman
		codeAlg crypto.EncryptionAlgorithm
	}
	type res struct {
		want          *domain.ObjectDetails
		wantEmailCode *string
		wantPhoneCode *string
		err           func(error) bool
	}

	userAgg := user.NewAggregate("user1", "org1")

	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "domain policy not found, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectFilter(),
					expectFilter(),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Username: gu.Ptr("changed"),
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-79pv6e1q62", "Errors.Org.DomainPolicy.NotExisting"))
				},
			},
		},
		{
			name: "change human username, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewDomainPolicyAddedEvent(context.Background(),
								&userAgg.Aggregate,
								true,
								true,
								true,
							),
						),
					),
					expectPush(
						user.NewUsernameChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"username",
							"changed",
							true,
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Username: gu.Ptr("changed"),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human username, no change",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Username: gu.Ptr("username"),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human profile, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectPush(
						func() eventstore.Command {
							cmd, _ := user.NewHumanProfileChangedEvent(context.Background(),
								&userAgg.Aggregate,
								[]user.ProfileChanges{
									user.ChangeFirstName("changedfn"),
									user.ChangeLastName("changedln"),
									user.ChangeNickName("changednn"),
									user.ChangeDisplayName("changeddn"),
									user.ChangePreferredLanguage(language.Afrikaans),
									user.ChangeGender(domain.GenderDiverse),
								},
							)
							return cmd
						}(),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Profile: &Profile{
						FirstName:         gu.Ptr("changedfn"),
						LastName:          gu.Ptr("changedln"),
						NickName:          gu.Ptr("changednn"),
						DisplayName:       gu.Ptr("changeddn"),
						PreferredLanguage: gu.Ptr(language.Afrikaans),
						Gender:            gu.Ptr(domain.GenderDiverse),
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human profile, no change",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Profile: &Profile{
						FirstName:         gu.Ptr("firstname"),
						LastName:          gu.Ptr("lastname"),
						NickName:          gu.Ptr(""),
						DisplayName:       gu.Ptr("firstname lastname"),
						PreferredLanguage: gu.Ptr(language.English),
						Gender:            gu.Ptr(domain.GenderUnspecified),
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human email, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectPush(
						user.NewHumanEmailChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"changed@example.com",
						),
						user.NewHumanEmailCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("emailCode"),
							},
							time.Hour,
							"",
							false,
						),
					),
				),
				newCode: mockCode("emailCode", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Email: &Email{
						Address: "changed@example.com",
					},
				},
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human email, no change",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Email: &Email{
						Address: "email@test.ch",
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human email verified, not allowed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Email: &Email{
						Address:  "changed@example.com",
						Verified: true,
					},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			name: "change human email verified, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectPush(
						user.NewHumanEmailChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"changed@example.com",
						),
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Email: &Email{
						Address:  "changed@example.com",
						Verified: true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human email isVerified, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectPush(
						user.NewHumanEmailVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Email: &Email{
						Verified: true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human email returnCode, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectPush(
						user.NewHumanEmailChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"changed@test.com",
						),
						user.NewHumanEmailCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("emailCode"),
							},
							time.Hour,
							"",
							true,
						),
					),
				),
				newCode: mockCode("emailCode", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Email: &Email{
						Address:    "changed@test.com",
						ReturnCode: true,
					},
				},
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
				wantEmailCode: gu.Ptr("emailCode"),
			},
		},
		{
			name: "change human phone, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectPush(
						user.NewHumanPhoneChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"+41791234567",
						),
						user.NewHumanPhoneCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("phoneCode"),
							},
							time.Hour,
							false,
						),
					),
				),
				newCode: mockCode("phoneCode", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Phone: &Phone{
						Number: "+41791234567",
					},
				},
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		}, {
			name: "change human phone verified, not allowed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Phone: &Phone{
						Number:   "+41791234567",
						Verified: true,
					},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			name: "change human phone verified, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectPush(
						user.NewHumanPhoneChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"+41791234567",
						),
						user.NewHumanPhoneVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Phone: &Phone{
						Number:   "+41791234567",
						Verified: true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human phone isVerified, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectPush(
						user.NewHumanPhoneVerifiedEvent(context.Background(),
							&userAgg.Aggregate,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Phone: &Phone{
						Verified: true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human phone returnCode, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
					),
					expectPush(
						user.NewHumanPhoneChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"+41791234567",
						),
						user.NewHumanPhoneCodeAddedEventV2(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("phoneCode"),
							},
							time.Hour,
							true,
						),
					),
				),
				newCode: mockCode("phoneCode", time.Hour),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Phone: &Phone{
						Number:     "+41791234567",
						ReturnCode: true,
					},
				},
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
				wantPhoneCode: gu.Ptr("phoneCode"),
			},
		},
		{
			name: "password change, no password, invalid argument error",
			fields: fields{
				eventstore:         expectEventstore(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-3klek4sbns", "Errors.User.Password.Empty"))
				},
			},
		},
		{
			name: "change human password, not initialized",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitialCodeAddedEvent(context.Background(),
								&userAgg.Aggregate,
								nil, time.Hour*1,
							),
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						Password:       gu.Ptr("password2"),
						OldPassword:    gu.Ptr("password"),
						ChangeRequired: true,
					},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-M9dse", "Errors.User.NotInitialised"))
				},
			},
		},
		{
			name: "change human password, not in complexity",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								1,
								false,
								true,
								false,
								false,
							),
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						Password:       gu.Ptr("password2"),
						OldPassword:    gu.Ptr("password"),
						ChangeRequired: true,
					},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "DOMAIN-VoaRj", "Errors.User.PasswordComplexityPolicy.HasUpper"))
				},
			},
		},
		{
			name: "change human password, empty",
			fields: fields{
				eventstore:         expectEventstore(),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						OldPassword:    gu.Ptr("password"),
						ChangeRequired: true,
					},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-3klek4sbns", "Errors.User.Password.Empty"))
				},
			},
		},
		{
			name: "change human password, not allowed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
				checkPermission:    newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						Password:       gu.Ptr("password2"),
						ChangeRequired: true,
					},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"))
				},
			},
		},
		{
			name: "change human password, permission, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						user.NewHumanPasswordChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"$plain$x$password2",
							true,
							"",
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
				checkPermission:    newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						Password:       gu.Ptr("password2"),
						ChangeRequired: true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human password, old password, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						user.NewHumanPasswordChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"$plain$x$password2",
							true,
							"",
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						Password:       gu.Ptr("password2"),
						OldPassword:    gu.Ptr("password"),
						ChangeRequired: true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human password, password code, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordCodeAddedEventV2(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour*1,
								domain.NotificationTypeEmail,
								"",
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						user.NewHumanPasswordChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"$plain$x$password2",
							true,
							"",
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						Password:       gu.Ptr("password2"),
						PasswordCode:   gu.Ptr("code"),
						ChangeRequired: true,
					},
				},
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human password encoded, password code, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordCodeAddedEventV2(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour*1,
								domain.NotificationTypeEmail,
								"",
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						user.NewHumanPasswordChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"$plain$x$password2",
							true,
							"",
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						EncodedPasswordHash: gu.Ptr("$plain$x$password2"),
						PasswordCode:        gu.Ptr("code"),
						ChangeRequired:      true,
					},
				},
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "change human password and password encoded, password code, encoded used",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							newAddHumanEvent("$plain$x$password", true, true, "", language.English),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&userAgg.Aggregate,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							user.NewHumanPasswordCodeAddedEventV2(context.Background(),
								&userAgg.Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("code"),
								},
								time.Hour*1,
								domain.NotificationTypeEmail,
								"",
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								1,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectPush(
						user.NewHumanPasswordChangedEvent(context.Background(),
							&userAgg.Aggregate,
							"$plain$x$password2",
							true,
							"",
						),
					),
				),
				userPasswordHasher: mockPasswordHasher("x"),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				human: &ChangeHuman{
					Password: &Password{
						Password:            gu.Ptr("passwordnotused"),
						EncodedPasswordHash: gu.Ptr("$plain$x$password2"),
						PasswordCode:        gu.Ptr("code"),
						ChangeRequired:      true,
					},
				},
				codeAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res: res{
				want: &domain.ObjectDetails{
					Sequence:      0,
					EventDate:     time.Time{},
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:         tt.fields.eventstore(t),
				userPasswordHasher: tt.fields.userPasswordHasher,
				newCode:            tt.fields.newCode,
				checkPermission:    tt.fields.checkPermission,
			}
			err := r.ChangeUserHuman(tt.args.ctx, tt.args.orgID, tt.args.human, tt.args.codeAlg)
			if tt.res.err == nil {
				if !assert.NoError(t, err) {
					t.FailNow()
				}
			} else if !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, tt.args.human.Details)
				assert.Equal(t, tt.res.wantEmailCode, tt.args.human.EmailCode)
				assert.Equal(t, tt.res.wantPhoneCode, tt.args.human.PhoneCode)
			}
		})
	}
}

func TestCommandSide_LockUserHuman(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx    context.Context
			orgID  string
			userID string
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "userid missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-agz3eczifm", "Errors.User.UserIDMissing"))
				},
			},
		},
		{
			name: "user not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-450yxuqrh1", "Errors.User.NotFound"))
				},
			},
		},
		{
			name: "user already locked, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewUserLockedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-lgws8wtsqf", "Errors.User.ShouldBeActiveOrInitial"))
				},
			},
		},
		{
			name: "user already locked, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewMachineAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"name",
								"description",
								true,
								domain.OIDCTokenTypeBearer,
							),
						),
						eventFromEventPusher(
							user.NewUserLockedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-lgws8wtsqf", "Errors.User.ShouldBeActiveOrInitial"))
				},
			},
		},
		{
			name: "lock user, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectPush(
						user.NewUserLockedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "lock user, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewMachineAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"name",
								"description",
								true,
								domain.OIDCTokenTypeBearer,
							),
						),
					),
					expectPush(
						user.NewUserLockedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.LockUserHuman(tt.args.ctx, tt.args.orgID, tt.args.userID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UnlockUserHuman(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx    context.Context
			orgID  string
			userID string
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "userid missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-a9ld4xckax", "Errors.User.UserIDMissing"))
				},
			},
		},
		{
			name: "user not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-x377t913pw", "Errors.User.NotFound"))
				},
			},
		},
		{
			name: "user already active, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-olb9vb0oca", "Errors.User.NotLocked"))
				},
			},
		},
		{
			name: "user already active, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						user.NewMachineAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"username",
							"name",
							"description",
							true,
							domain.OIDCTokenTypeBearer,
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-olb9vb0oca", "Errors.User.NotLocked"))
				},
			},
		},
		{
			name: "unlock user, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewUserLockedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate),
						),
					),
					expectPush(
						user.NewUserUnlockedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "unlock user machine, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewMachineAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"name",
								"description",
								true,
								domain.OIDCTokenTypeBearer,
							),
						),
						eventFromEventPusher(
							user.NewUserLockedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate),
						),
					),
					expectPush(
						user.NewUserUnlockedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.UnlockUserHuman(tt.args.ctx, tt.args.orgID, tt.args.userID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_DeactivateUserHuman(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx    context.Context
			orgID  string
			userID string
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{{
		name: "resourceowner missing, invalid argument error",
		fields: fields{
			eventstore: eventstoreExpect(
				t,
			),
		},
		args: args{
			ctx:    context.Background(),
			orgID:  "",
			userID: "",
		},
		res: res{
			err: func(err error) bool {
				return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMA-r1etf6eex9", "Errors.Internal"))
			},
		},
	},
		{
			name: "userid missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-78iiirat8y", "Errors.User.UserIDMissing"))
				},
			},
		},
		{
			name: "user not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-5gp2p62iin", "Errors.User.NotFound"))
				},
			},
		},
		{
			name: "user initial, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanInitialCodeAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								nil, time.Hour*1,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-gvx4kct9r2", "Errors.User.CantDeactivateInitial"))
				},
			},
		},
		{
			name: "user already inactive, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewUserDeactivatedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-5gunjw0cd7", "Errors.User.AlreadyInactive"))
				},
			},
		},
		{
			name: "deactivate user, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewHumanInitializedCheckSucceededEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
					expectPush(
						user.NewUserDeactivatedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},

		{
			name: "user machine already inactive, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewMachineAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"name",
								"description",
								true,
								domain.OIDCTokenTypeBearer,
							),
						),
						eventFromEventPusher(
							user.NewUserDeactivatedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-5gunjw0cd7", "Errors.User.AlreadyInactive"))
				},
			},
		},
		{
			name: "deactivate user machine, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewMachineAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"name",
								"description",
								true,
								domain.OIDCTokenTypeBearer,
							),
						),
					),
					expectPush(
						user.NewUserDeactivatedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.DeactivateUserHuman(tt.args.ctx, tt.args.orgID, tt.args.userID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_ReactivateUserHuman(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx    context.Context
			orgID  string
			userID string
		}
	)
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "resourceowner missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "",
				userID: "",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMA-0vlu1bb5uo", "Errors.Internal"))
				},
			},
		},
		{
			name: "userid missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "COMMAND-0nx1ie38fw", "Errors.User.UserIDMissing"))
				},
			},
		},
		{
			name: "user not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-9hy5kzbuk6", "Errors.User.NotFound"))
				},
			},
		},
		{
			name: "user already active, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-s5qqcz97hf", "Errors.User.NotInactive"))
				},
			},
		},
		{
			name: "user machine already active, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewMachineAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"name",
								"description",
								true,
								domain.OIDCTokenTypeBearer,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowPreconditionFailed(nil, "COMMAND-s5qqcz97hf", "Errors.User.NotInactive"))
				},
			},
		},
		{
			name: "reactivate user, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
						eventFromEventPusher(
							user.NewUserDeactivatedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate),
						),
					),
					expectPush(
						user.NewUserReactivatedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "reactivate user machine, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							user.NewMachineAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"name",
								"description",
								true,
								domain.OIDCTokenTypeBearer,
							),
						),
						eventFromEventPusher(
							user.NewUserDeactivatedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate),
						),
					),
					expectPush(
						user.NewUserReactivatedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				orgID:  "org1",
				userID: "user1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ReactivateUserHuman(tt.args.ctx, tt.args.orgID, tt.args.userID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}
//go:build integration

package user_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/repository/idp"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

var (
	CTX    context.Context
	ErrCTX context.Context
	Tester *integration.Tester
	Client user.UserServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, errCtx, cancel := integration.Contexts(time.Hour)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()

		CTX, ErrCTX = Tester.WithSystemAuthorization(ctx, integration.OrgOwner), errCtx
		Client = user.NewUserServiceClient(Tester.GRPCClientConn)
		return m.Run()
	}())
}

func TestServer_AddHumanUser(t *testing.T) {
	type args struct {
		ctx context.Context
		req *user.AddHumanUserRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.AddHumanUserResponse
		wantErr bool
	}{
		{
			name: "default verification",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Profile: &user.SetHumanProfile{
						FirstName:         "Donald",
						LastName:          "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			want: &user.AddHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "return verification code",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Profile: &user.SetHumanProfile{
						FirstName:         "Donald",
						LastName:          "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{
						Verification: &user.SetHumanEmail_ReturnCode{
							ReturnCode: &user.ReturnEmailVerificationCode{},
						},
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			want: &user.AddHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
				EmailCode: gu.Ptr("something"),
			},
		},
		{
			name: "custom template",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Profile: &user.SetHumanProfile{
						FirstName:         "Donald",
						LastName:          "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{
						Verification: &user.SetHumanEmail_SendCode{
							SendCode: &user.SendEmailVerificationCode{
								UrlTemplate: gu.Ptr("https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}"),
							},
						},
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			want: &user.AddHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "custom template error",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Profile: &user.SetHumanProfile{
						FirstName:         "Donald",
						LastName:          "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{
						Verification: &user.SetHumanEmail_SendCode{
							SendCode: &user.SendEmailVerificationCode{
								UrlTemplate: gu.Ptr("{{"),
							},
						},
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing REQUIRED profile",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Email: &user.SetHumanEmail{
						Verification: &user.SetHumanEmail_ReturnCode{
							ReturnCode: &user.ReturnEmailVerificationCode{},
						},
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing REQUIRED email",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Profile: &user.SetHumanProfile{
						FirstName:         "Donald",
						LastName:          "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := fmt.Sprint(time.Now().UnixNano() + int64(i))
			tt.args.req.UserId = &userID
			if email := tt.args.req.GetEmail(); email != nil {
				email.Email = fmt.Sprintf("%s@me.now", userID)
			}

			if tt.want != nil {
				tt.want.UserId = userID
			}

			got, err := Client.AddHumanUser(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want.GetUserId(), got.GetUserId())
			if tt.want.GetEmailCode() != "" {
				assert.NotEmpty(t, got.GetEmailCode())
			}
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func createProvider(t *testing.T) string {
	ctx := authz.WithInstance(context.Background(), Tester.Instance)
	id, _, err := Tester.Commands.AddOrgGenericOAuthProvider(ctx, Tester.Organisation.ID, command.GenericOAuthProvider{
		"idp",
		"clientID",
		"clientSecret",
		"https://example.com/oauth/v2/authorize",
		"https://example.com/oauth/v2/token",
		"https://api.example.com/user",
		[]string{"openid", "profile", "email"},
		"id",
		idp.Options{
			IsLinkingAllowed:  true,
			IsCreationAllowed: true,
			IsAutoCreation:    true,
			IsAutoUpdate:      true,
		},
	})
	require.NoError(t, err)
	return id
}

func TestServer_StartIdentityProviderFlow(t *testing.T) {
	idpID := createProvider(t)
	type args struct {
		ctx context.Context
		req *user.StartIdentityProviderFlowRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.StartIdentityProviderFlowResponse
		wantErr bool
	}{
		{
			name: "missing urls",
			args: args{
				CTX,
				&user.StartIdentityProviderFlowRequest{
					IdpId: idpID,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "next step auth url",
			args: args{
				CTX,
				&user.StartIdentityProviderFlowRequest{
					IdpId:      idpID,
					SuccessUrl: "https://example.com/success",
					FailureUrl: "https://example.com/failure",
				},
			},
			want: &user.StartIdentityProviderFlowResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
				NextStep: &user.StartIdentityProviderFlowResponse_AuthUrl{
					AuthUrl: "https://example.com/oauth/v2/authorize?client_id=clientID&prompt=select_account&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fidps%2Fcallback&response_type=code&scope=openid+profile+email&state=",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.StartIdentityProviderFlow(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if nextStep := tt.want.GetNextStep(); nextStep != nil {
				assert.True(t, strings.HasPrefix(got.GetAuthUrl(), tt.want.GetAuthUrl()))
			}
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

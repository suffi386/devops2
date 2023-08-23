package command

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestCommands_CreateOTPSMSChallengeReturnCode(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
		createCode cryptoCodeFunc
	}
	type res struct {
		err        error
		returnCode string
		commands   []eventstore.Command
	}
	tests := []struct {
		name   string
		fields fields
		res    res
	}{
		{
			name: "otp not ready, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			res: res{
				err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-BJ2g3", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "generate code",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPSMSAddedEvent(context.Background(), &user.NewAggregate("userID", "org").Aggregate),
						),
					),
				),
				createCode: mockCode("1234567", 5*time.Minute),
			},
			res: res{
				returnCode: "1234567",
				commands: []eventstore.Command{
					session.NewOTPSMSChallengedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
						&crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "enc",
							KeyID:      "id",
							Crypted:    []byte("1234567"),
						},
						5*time.Minute,
						true,
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{}
			var dst string
			cmd := c.CreateOTPSMSChallengeReturnCode(&dst)

			sessionModel := &SessionWriteModel{
				UserID:        "userID",
				UserCheckedAt: testNow,
				State:         domain.SessionStateActive,
				aggregate:     &session.NewAggregate("sessionID", "instanceID").Aggregate,
			}
			cmds := &SessionCommands{
				sessionCommands:   []SessionCommand{cmd},
				sessionWriteModel: sessionModel,
				eventstore:        tt.fields.eventstore(t),
				createCode:        tt.fields.createCode,
				now:               time.Now,
			}

			err := cmd(context.Background(), cmds)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.returnCode, dst)
			assert.Equal(t, tt.res.commands, cmds.eventCommands)
		})
	}
}

func TestCommands_CreateOTPSMSChallenge(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
		createCode cryptoCodeFunc
	}
	type res struct {
		err      error
		commands []eventstore.Command
	}
	tests := []struct {
		name   string
		fields fields
		res    res
	}{
		{
			name: "otp not ready, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			res: res{
				err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-BJ2g3", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "generate code",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPSMSAddedEvent(context.Background(), &user.NewAggregate("userID", "org").Aggregate),
						),
					),
				),
				createCode: mockCode("1234567", 5*time.Minute),
			},
			res: res{
				commands: []eventstore.Command{
					session.NewOTPSMSChallengedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
						&crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "enc",
							KeyID:      "id",
							Crypted:    []byte("1234567"),
						},
						5*time.Minute,
						false,
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{}

			cmd := c.CreateOTPSMSChallenge()

			sessionModel := &SessionWriteModel{
				UserID:        "userID",
				UserCheckedAt: testNow,
				State:         domain.SessionStateActive,
				aggregate:     &session.NewAggregate("sessionID", "instanceID").Aggregate,
			}
			cmds := &SessionCommands{
				sessionCommands:   []SessionCommand{cmd},
				sessionWriteModel: sessionModel,
				eventstore:        tt.fields.eventstore(t),
				createCode:        tt.fields.createCode,
				now:               time.Now,
			}

			err := cmd(context.Background(), cmds)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.commands, cmds.eventCommands)
		})
	}
}

func TestCommands_OTPSMSSent(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		sessionID     string
		resourceOwner string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "not challenged, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				sessionID:     "sessionID",
				resourceOwner: "instanceID",
			},
			wantErr: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-G3t31", "Errors.User.Code.NotFound"),
		},
		{
			name: "challenged and sent",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							session.NewOTPSMSChallengedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("1234567"),
								},
								5*time.Minute,
								false,
							),
						),
					),
					expectPush(
						eventPusherToEvents(
							session.NewOTPSMSSentEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				sessionID:     "sessionID",
				resourceOwner: "instanceID",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			err := c.OTPSMSSent(tt.args.ctx, tt.args.sessionID, tt.args.resourceOwner)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCommands_CreateOTPEmailChallengeURLTemplate(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
		createCode cryptoCodeFunc
	}
	type args struct {
		urlTmpl string
	}
	type res struct {
		templateError error
		err           error
		commands      []eventstore.Command
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid template, precondition error",
			args: args{
				urlTmpl: "https://example.com/mfa/email?userID={{.UserID}}&code={{.InvalidField}}",
			},
			fields: fields{
				eventstore: expectEventstore(),
			},
			res: res{
				templateError: caos_errs.ThrowInvalidArgument(nil, "DOMAIN-ieYa7", "Errors.User.InvalidURLTemplate"),
			},
		},
		{
			name: "otp not ready, precondition error",
			args: args{
				urlTmpl: "https://example.com/mfa/email?userID={{.UserID}}&code={{.Code}}",
			},
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			res: res{
				err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-JKLJ3", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "generate code",
			args: args{
				urlTmpl: "https://example.com/mfa/email?userID={{.UserID}}&code={{.Code}}",
			},
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPEmailAddedEvent(context.Background(), &user.NewAggregate("userID", "org").Aggregate),
						),
					),
				),
				createCode: mockCode("1234567", 5*time.Minute),
			},
			res: res{
				commands: []eventstore.Command{
					session.NewOTPEmailChallengedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
						&crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "enc",
							KeyID:      "id",
							Crypted:    []byte("1234567"),
						},
						5*time.Minute,
						false,
						"https://example.com/mfa/email?userID={{.UserID}}&code={{.Code}}",
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{}

			cmd, err := c.CreateOTPEmailChallengeURLTemplate(tt.args.urlTmpl)
			assert.ErrorIs(t, err, tt.res.templateError)
			if tt.res.templateError != nil {
				return
			}

			sessionModel := &SessionWriteModel{
				UserID:        "userID",
				UserCheckedAt: testNow,
				State:         domain.SessionStateActive,
				aggregate:     &session.NewAggregate("sessionID", "instanceID").Aggregate,
			}
			cmds := &SessionCommands{
				sessionCommands:   []SessionCommand{cmd},
				sessionWriteModel: sessionModel,
				eventstore:        tt.fields.eventstore(t),
				createCode:        tt.fields.createCode,
				now:               time.Now,
			}

			err = cmd(context.Background(), cmds)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.commands, cmds.eventCommands)
		})
	}
}

func TestCommands_CreateOTPEmailChallengeReturnCode(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
		createCode cryptoCodeFunc
	}
	type res struct {
		err        error
		returnCode string
		commands   []eventstore.Command
	}
	tests := []struct {
		name   string
		fields fields
		res    res
	}{
		{
			name: "otp not ready, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			res: res{
				err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-JKLJ3", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "generate code",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPEmailAddedEvent(context.Background(), &user.NewAggregate("userID", "org").Aggregate),
						),
					),
				),
				createCode: mockCode("1234567", 5*time.Minute),
			},
			res: res{
				returnCode: "1234567",
				commands: []eventstore.Command{
					session.NewOTPEmailChallengedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
						&crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "enc",
							KeyID:      "id",
							Crypted:    []byte("1234567"),
						},
						5*time.Minute,
						true,
						"",
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{}
			var dst string
			cmd := c.CreateOTPEmailChallengeReturnCode(&dst)

			sessionModel := &SessionWriteModel{
				UserID:        "userID",
				UserCheckedAt: testNow,
				State:         domain.SessionStateActive,
				aggregate:     &session.NewAggregate("sessionID", "instanceID").Aggregate,
			}
			cmds := &SessionCommands{
				sessionCommands:   []SessionCommand{cmd},
				sessionWriteModel: sessionModel,
				eventstore:        tt.fields.eventstore(t),
				createCode:        tt.fields.createCode,
				now:               time.Now,
			}

			err := cmd(context.Background(), cmds)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.returnCode, dst)
			assert.Equal(t, tt.res.commands, cmds.eventCommands)
		})
	}
}

func TestCommands_CreateOTPEmailChallenge(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
		createCode cryptoCodeFunc
	}
	type res struct {
		err      error
		commands []eventstore.Command
	}
	tests := []struct {
		name   string
		fields fields
		res    res
	}{
		{
			name: "otp not ready, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			res: res{
				err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-JKLJ3", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "generate code",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPEmailAddedEvent(context.Background(), &user.NewAggregate("userID", "org").Aggregate),
						),
					),
				),
				createCode: mockCode("1234567", 5*time.Minute),
			},
			res: res{
				commands: []eventstore.Command{
					session.NewOTPEmailChallengedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
						&crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "enc",
							KeyID:      "id",
							Crypted:    []byte("1234567"),
						},
						5*time.Minute,
						false,
						"",
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{}

			cmd := c.CreateOTPEmailChallenge()

			sessionModel := &SessionWriteModel{
				UserID:        "userID",
				UserCheckedAt: testNow,
				State:         domain.SessionStateActive,
				aggregate:     &session.NewAggregate("sessionID", "instanceID").Aggregate,
			}
			cmds := &SessionCommands{
				sessionCommands:   []SessionCommand{cmd},
				sessionWriteModel: sessionModel,
				eventstore:        tt.fields.eventstore(t),
				createCode:        tt.fields.createCode,
				now:               time.Now,
			}

			err := cmd(context.Background(), cmds)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.commands, cmds.eventCommands)
		})
	}
}

func TestCommands_OTPEmailSent(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		sessionID     string
		resourceOwner string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "not challenged, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				sessionID:     "sessionID",
				resourceOwner: "instanceID",
			},
			wantErr: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-SLr02", "Errors.User.Code.NotFound"),
		},
		{
			name: "challenged and sent",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							session.NewOTPEmailChallengedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("1234567"),
								},
								5*time.Minute,
								false,
								"",
							),
						),
					),
					expectPush(
						eventPusherToEvents(
							session.NewOTPEmailSentEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				sessionID:     "sessionID",
				resourceOwner: "instanceID",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			err := c.OTPEmailSent(tt.args.ctx, tt.args.sessionID, tt.args.resourceOwner)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCheckOTPSMS(t *testing.T) {
	type fields struct {
		eventstore       func(*testing.T) *eventstore.Eventstore
		userID           string
		otpCodeChallenge *OTPCode
		otpAlg           crypto.EncryptionAlgorithm
	}
	type args struct {
		code string
	}
	type res struct {
		err      error
		commands []eventstore.Command
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "missing userID",
			fields: fields{
				eventstore: expectEventstore(),
				userID:     "",
			},
			args: args{
				code: "code",
			},
			res: res{
				err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-VDrh3", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "missing challenge",
			fields: fields{
				eventstore:       expectEventstore(),
				userID:           "userID",
				otpCodeChallenge: nil,
			},
			args: args{
				code: "code",
			},
			res: res{
				err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-SF3tv", "Errors.User.Code.NotFound"),
			},
		},
		{
			name: "invalid code",
			fields: fields{
				eventstore: expectEventstore(),
				userID:     "userID",
				otpCodeChallenge: &OTPCode{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:       5 * time.Minute,
					CreationDate: testNow.Add(-10 * time.Minute),
				},
				otpAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				code: "code",
			},
			res: res{
				err: caos_errs.ThrowPreconditionFailed(nil, "CODE-QvUQ4P", "Errors.User.Code.Expired"),
			},
		},
		{
			name: "check ok",
			fields: fields{
				eventstore: expectEventstore(),
				userID:     "userID",
				otpCodeChallenge: &OTPCode{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:       5 * time.Minute,
					CreationDate: testNow,
				},
				otpAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				code: "code",
			},
			res: res{
				commands: []eventstore.Command{
					session.NewOTPSMSCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
						testNow,
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := CheckOTPSMS(tt.args.code)

			sessionModel := &SessionWriteModel{
				UserID:              tt.fields.userID,
				UserCheckedAt:       testNow,
				State:               domain.SessionStateActive,
				OTPSMSCodeChallenge: tt.fields.otpCodeChallenge,
				aggregate:           &session.NewAggregate("sessionID", "instanceID").Aggregate,
			}
			cmds := &SessionCommands{
				sessionCommands:   []SessionCommand{cmd},
				sessionWriteModel: sessionModel,
				eventstore:        tt.fields.eventstore(t),
				otpAlg:            tt.fields.otpAlg,
				now: func() time.Time {
					return testNow
				},
			}

			err := cmd(context.Background(), cmds)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.commands, cmds.eventCommands)
		})
	}
}

func TestCheckOTPEmail(t *testing.T) {
	type fields struct {
		eventstore       func(*testing.T) *eventstore.Eventstore
		userID           string
		otpCodeChallenge *OTPCode
		otpAlg           crypto.EncryptionAlgorithm
	}
	type args struct {
		code string
	}
	type res struct {
		err      error
		commands []eventstore.Command
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "missing userID",
			fields: fields{
				eventstore: expectEventstore(),
				userID:     "",
			},
			args: args{
				code: "code",
			},
			res: res{
				err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-ejo2w", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "missing challenge",
			fields: fields{
				eventstore:       expectEventstore(),
				userID:           "userID",
				otpCodeChallenge: nil,
			},
			args: args{
				code: "code",
			},
			res: res{
				err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-zF3g3", "Errors.User.Code.NotFound"),
			},
		},
		{
			name: "invalid code",
			fields: fields{
				eventstore: expectEventstore(),
				userID:     "userID",
				otpCodeChallenge: &OTPCode{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:       5 * time.Minute,
					CreationDate: testNow.Add(-10 * time.Minute),
				},
				otpAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				code: "code",
			},
			res: res{
				err: caos_errs.ThrowPreconditionFailed(nil, "CODE-QvUQ4P", "Errors.User.Code.Expired"),
			},
		},
		{
			name: "check ok",
			fields: fields{
				eventstore: expectEventstore(),
				userID:     "userID",
				otpCodeChallenge: &OTPCode{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:       5 * time.Minute,
					CreationDate: testNow,
				},
				otpAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				code: "code",
			},
			res: res{
				commands: []eventstore.Command{
					session.NewOTPEmailCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
						testNow,
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := CheckOTPEmail(tt.args.code)

			sessionModel := &SessionWriteModel{
				UserID:                tt.fields.userID,
				UserCheckedAt:         testNow,
				State:                 domain.SessionStateActive,
				OTPEmailCodeChallenge: tt.fields.otpCodeChallenge,
				aggregate:             &session.NewAggregate("sessionID", "instanceID").Aggregate,
			}
			cmds := &SessionCommands{
				sessionCommands:   []SessionCommand{cmd},
				sessionWriteModel: sessionModel,
				eventstore:        tt.fields.eventstore(t),
				otpAlg:            tt.fields.otpAlg,
				now: func() time.Time {
					return testNow
				},
			}

			err := cmd(context.Background(), cmds)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.commands, cmds.eventCommands)
		})
	}
}

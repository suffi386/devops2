package systemdefaults

import (
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/notification/providers/chat"
	"github.com/caos/zitadel/internal/notification/providers/email"
	"github.com/caos/zitadel/internal/notification/providers/twilio"
	"github.com/caos/zitadel/internal/notification/templates"
	pol "github.com/caos/zitadel/internal/policy"
)

type SystemDefaults struct {
	SecretGenerators    SecretGenerators
	UserVerificationKey *crypto.KeyConfig
	Multifactors        MultifactorConfig
	DefaultPolicies     DefaultPolicies
	Notifications       Notifications
}

type SecretGenerators struct {
	PasswordSaltCost         int
	ClientSecretGenerator    crypto.GeneratorConfig
	InitializeUserCode       crypto.GeneratorConfig
	EmailVerificationCode    crypto.GeneratorConfig
	PhoneVerificationCode    crypto.GeneratorConfig
	PasswordVerificationCode crypto.GeneratorConfig
}

type MultifactorConfig struct {
	OTP OTPConfig
}

type OTPConfig struct {
	Issuer          string
	VerificationKey *crypto.KeyConfig
}

type DefaultPolicies struct {
	Age        pol.PasswordAgePolicyDefault
	Complexity pol.PasswordComplexityPolicyDefault
	Lockout    pol.PasswordLockoutPolicyDefault
}

type Notifications struct {
	Debug        bool
	Endpoints    Endpoints
	Providers    Providers
	TemplateData TemplateData
}

type Endpoints struct {
	InitCode      string
	PasswordReset string
	VerifyEmail   string
}

type Providers struct {
	Chat   chat.ChatConfig
	Email  email.EmailConfig
	Twilio twilio.TwilioConfig
}

type TemplateData struct {
	InitCode      templates.TemplateData
	PasswordReset templates.TemplateData
}

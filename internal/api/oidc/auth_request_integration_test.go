//go:build integration

package oidc_test

import (
	"context"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/oidc/amr"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/integration"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2alpha"
	session "github.com/zitadel/zitadel/pkg/grpc/session/v2alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

var (
	CTX      context.Context
	CTXLOGIN context.Context
	Tester   *integration.Tester
	User     *user.AddHumanUserResponse
)

const (
	redirectURI         = "oidcIntegrationTest://callback"
	redirectURIImplicit = "http://localhost:9999/callback"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, errCtx, cancel := integration.Contexts(5 * time.Minute)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()

		CTX, _ = Tester.WithSystemAuthorization(ctx, integration.OrgOwner), errCtx
		User = Tester.CreateHumanUser(CTX)
		Tester.RegisterUserPasskey(CTX, User.GetUserId())
		CTXLOGIN, _ = Tester.WithSystemAuthorization(ctx, integration.Login), errCtx
		return m.Run()
	}())
}

func createClient(t testing.TB) string {
	app, err := Tester.CreateOIDCNativeClient(CTX, redirectURI)
	require.NoError(t, err)
	return app.GetClientId()
}

func createImplicitClient(t testing.TB) string {
	app, err := Tester.CreateOIDCImplicitFlowClient(CTX, redirectURIImplicit)
	require.NoError(t, err)
	return app.GetClientId()
}

func createAuthRequest(t testing.TB, clientID, redirectURI string, scope ...string) string {
	redURL, err := Tester.CreateOIDCAuthRequest(clientID, integration.LoginUser, redirectURI, scope...)
	require.NoError(t, err)
	return redURL
}

func createAuthRequestImplicit(t testing.TB, clientID, redirectURI string, scope ...string) string {
	redURL, err := Tester.CreateOIDCAuthRequestImplicit(clientID, integration.LoginUser, redirectURI, scope...)
	require.NoError(t, err)
	return redURL
}

func TestOPStorage_CreateAuthRequest(t *testing.T) {
	clientID := createClient(t)

	id := createAuthRequest(t, clientID, redirectURI)
	require.Contains(t, id, command.IDPrefixV2)
}

func TestOPStorage_CreateAccessToken_code(t *testing.T) {
	clientID := createClient(t)

	id := createAuthRequest(t, clientID, redirectURI)
	_ = id
	createResp, err := Tester.Client.SessionV2.CreateSession(CTXLOGIN, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{UserId: User.GetUserId()},
			},
		},
		Challenges: []session.ChallengeKind{
			session.ChallengeKind_CHALLENGE_KIND_PASSKEY,
		},
	})
	require.NoError(t, err)

	assertion, err := Tester.WebAuthN.CreateAssertionResponse(createResp.GetChallenges().GetPasskey().GetPublicKeyCredentialRequestOptions())
	require.NoError(t, err)

	updateResp, err := Tester.Client.SessionV2.SetSession(CTXLOGIN, &session.SetSessionRequest{
		SessionId:    createResp.GetSessionId(),
		SessionToken: createResp.GetSessionToken(),
		Checks: &session.Checks{
			Passkey: &session.CheckPasskey{
				CredentialAssertionData: assertion,
			},
		},
	})
	require.NoError(t, err)

	linkResp, err := Tester.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: id,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    createResp.GetSessionId(),
				SessionToken: updateResp.GetSessionToken(),
			},
		},
	})
	require.NoError(t, err)

	callback, err := integration.CheckRedirect(http_util.BuildOrigin(Tester.Host(), Tester.Config.ExternalSecure)+linkResp.GetCallbackUrl(), Tester.WithSystemAuthorizationHTTP(integration.Login))
	require.NoError(t, err)
	code := callback.Query().Get("code")
	require.NotEmpty(t, code)

	tokens, err := exchangeTokens(t, clientID, code)
	require.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.IDToken)
	assert.Empty(t, tokens.RefreshToken)
	assert.Equal(t, []string{amr.UserPresence, amr.MFA}, tokens.IDTokenClaims.AuthenticationMethodsReferences)
	assert.WithinRange(t, tokens.IDTokenClaims.AuthTime.AsTime().UTC(), createResp.Details.ChangeDate.AsTime().Add(-1*time.Second), updateResp.Details.ChangeDate.AsTime())

	// exchange with a used code must fail
	_, err = exchangeTokens(t, clientID, code)
	require.Error(t, err)
}

func TestOPStorage_CreateAccessToken_implicit(t *testing.T) {
	clientID := createImplicitClient(t)

	id := createAuthRequestImplicit(t, clientID, redirectURIImplicit)

	createResp, err := Tester.Client.SessionV2.CreateSession(CTXLOGIN, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{UserId: User.GetUserId()},
			},
		},
		Challenges: []session.ChallengeKind{
			session.ChallengeKind_CHALLENGE_KIND_PASSKEY,
		},
	})
	require.NoError(t, err)

	assertion, err := Tester.WebAuthN.CreateAssertionResponse(createResp.GetChallenges().GetPasskey().GetPublicKeyCredentialRequestOptions())
	require.NoError(t, err)

	updateResp, err := Tester.Client.SessionV2.SetSession(CTXLOGIN, &session.SetSessionRequest{
		SessionId:    createResp.GetSessionId(),
		SessionToken: createResp.GetSessionToken(),
		Checks: &session.Checks{
			Passkey: &session.CheckPasskey{
				CredentialAssertionData: assertion,
			},
		},
	})
	require.NoError(t, err)

	linkResp, err := Tester.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: id,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    createResp.GetSessionId(),
				SessionToken: updateResp.GetSessionToken(),
			},
		},
	})
	require.NoError(t, err)

	callbackURI := http_util.BuildOrigin(Tester.Host(), Tester.Config.ExternalSecure) + linkResp.GetCallbackUrl()
	callback, err := integration.CheckRedirect(callbackURI, Tester.WithSystemAuthorizationHTTP(integration.Login))
	require.NoError(t, err)

	values, err := url.ParseQuery(callback.Fragment)
	require.NoError(t, err)

	accessToken := values.Get("access_token")
	idToken := values.Get("id_token")
	refreshToken := values.Get("refresh_token")
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, idToken)
	assert.Empty(t, refreshToken)

	provider, err := Tester.CreateRelyingParty(clientID, redirectURIImplicit)
	require.NoError(t, err)

	claims, err := rp.VerifyTokens[*oidc.IDTokenClaims](context.Background(), accessToken, idToken, provider.IDTokenVerifier())
	require.NoError(t, err)
	assert.Equal(t, []string{amr.UserPresence, amr.MFA}, claims.AuthenticationMethodsReferences)
	assert.WithinRange(t, claims.AuthTime.AsTime().UTC(), createResp.Details.ChangeDate.AsTime().Add(-1*time.Second), updateResp.Details.ChangeDate.AsTime())

	// callback on a succeeded request must fail
	callback, err = integration.CheckRedirect(callbackURI, Tester.WithSystemAuthorizationHTTP(integration.Login))
	require.NoError(t, err)
	values, err = url.ParseQuery(callback.Fragment)
	require.NoError(t, err)
	require.NotEmpty(t, values.Get("error"))
}

func TestOPStorage_CreateAccessAndRefreshTokens_code(t *testing.T) {
	clientID := createClient(t)

	id := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeOfflineAccess)
	_ = id
	createResp, err := Tester.Client.SessionV2.CreateSession(CTXLOGIN, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{UserId: User.GetUserId()},
			},
		},
		Challenges: []session.ChallengeKind{
			session.ChallengeKind_CHALLENGE_KIND_PASSKEY,
		},
	})
	require.NoError(t, err)

	assertion, err := Tester.WebAuthN.CreateAssertionResponse(createResp.GetChallenges().GetPasskey().GetPublicKeyCredentialRequestOptions())
	require.NoError(t, err)

	updateResp, err := Tester.Client.SessionV2.SetSession(CTXLOGIN, &session.SetSessionRequest{
		SessionId:    createResp.GetSessionId(),
		SessionToken: createResp.GetSessionToken(),
		Checks: &session.Checks{
			Passkey: &session.CheckPasskey{
				CredentialAssertionData: assertion,
			},
		},
	})
	require.NoError(t, err)

	linkResp, err := Tester.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: id,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    createResp.GetSessionId(),
				SessionToken: updateResp.GetSessionToken(),
			},
		},
	})
	require.NoError(t, err)

	callback, err := integration.CheckRedirect(http_util.BuildOrigin(Tester.Host(), Tester.Config.ExternalSecure)+linkResp.GetCallbackUrl(), Tester.WithSystemAuthorizationHTTP(integration.Login))
	require.NoError(t, err)
	code := callback.Query().Get("code")
	require.NotEmpty(t, code)

	tokens, err := exchangeTokens(t, clientID, code)
	require.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.IDToken)
	assert.NotEmpty(t, tokens.RefreshToken)
	assert.Equal(t, []string{amr.UserPresence, amr.MFA}, tokens.IDTokenClaims.AuthenticationMethodsReferences)
	assert.WithinRange(t, tokens.IDTokenClaims.AuthTime.AsTime().UTC(), createResp.Details.ChangeDate.AsTime().Add(-1*time.Second), updateResp.Details.ChangeDate.AsTime())

	// exchange with a used code must fail
	_, err = exchangeTokens(t, clientID, code)
	require.Error(t, err)
}

func TestOPStorage_CreateAccessAndRefreshTokens_refresh(t *testing.T) {
	clientID := createClient(t)

	id := createAuthRequest(t, clientID, redirectURI, oidc.ScopeOpenID, oidc.ScopeOfflineAccess)
	_ = id
	createResp, err := Tester.Client.SessionV2.CreateSession(CTXLOGIN, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{UserId: User.GetUserId()},
			},
		},
		Challenges: []session.ChallengeKind{
			session.ChallengeKind_CHALLENGE_KIND_PASSKEY,
		},
	})
	require.NoError(t, err)

	assertion, err := Tester.WebAuthN.CreateAssertionResponse(createResp.GetChallenges().GetPasskey().GetPublicKeyCredentialRequestOptions())
	require.NoError(t, err)

	updateResp, err := Tester.Client.SessionV2.SetSession(CTXLOGIN, &session.SetSessionRequest{
		SessionId:    createResp.GetSessionId(),
		SessionToken: createResp.GetSessionToken(),
		Checks: &session.Checks{
			Passkey: &session.CheckPasskey{
				CredentialAssertionData: assertion,
			},
		},
	})
	require.NoError(t, err)

	linkResp, err := Tester.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
		AuthRequestId: id,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    createResp.GetSessionId(),
				SessionToken: updateResp.GetSessionToken(),
			},
		},
	})
	require.NoError(t, err)

	callback, err := integration.CheckRedirect(http_util.BuildOrigin(Tester.Host(), Tester.Config.ExternalSecure)+linkResp.GetCallbackUrl(), Tester.WithSystemAuthorizationHTTP(integration.Login))
	require.NoError(t, err)
	code := callback.Query().Get("code")
	require.NotEmpty(t, code)

	tokens, err := exchangeTokens(t, clientID, code)
	require.NoError(t, err)
	assert.NotEmpty(t, tokens.RefreshToken)

	provider, err := Tester.CreateRelyingParty(clientID, redirectURI)
	require.NoError(t, err)

	newTokens, err := refreshTokens(t, clientID, tokens.RefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, newTokens.AccessToken)
	assert.NotEmpty(t, newTokens.Extra("id_token"))
	assert.NotEmpty(t, newTokens.RefreshToken)

	// refresh with an old refresh_token must fail
	_, err = rp.RefreshAccessToken(provider, tokens.RefreshToken, "", "")
	require.Error(t, err)
}

func exchangeTokens(t testing.TB, clientID, code string) (*oidc.Tokens[*oidc.IDTokenClaims], error) {
	provider, err := Tester.CreateRelyingParty(clientID, redirectURI)
	require.NoError(t, err)

	codeVerifier := "codeVerifier"
	return rp.CodeExchange[*oidc.IDTokenClaims](context.Background(), code, provider, rp.WithCodeVerifier(codeVerifier))
}

func refreshTokens(t testing.TB, clientID, refreshToken string) (*oauth2.Token, error) {
	provider, err := Tester.CreateRelyingParty(clientID, redirectURI)
	require.NoError(t, err)

	return rp.RefreshAccessToken(provider, refreshToken, "", "")
}

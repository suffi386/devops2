package oidc

import (
	"context"
	"encoding/base64"
	"slices"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/zitadel/oidc/v3/pkg/crypto"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

/*
For each grant-type, tokens creation follows the same rough logical steps:

1. Information gathering: who is requesting the token, what do we put in the claims?
2. Decision making: is the request authorized? (valid exchange code, auth request completed, valid token etc...)
3. Build an OIDC session in storage: inform the eventstore we are creating tokens.
4. Use the OIDC session to encrypt and / or sign the requested tokens

In some cases step 1 till 3 are completely implemented in the command package,
for example the v2 code exchange and refresh token.
*/

func (s *Server) accessTokenResponseFromSession(ctx context.Context, client op.Client, session *command.OIDCSession, state, projectID string, projectRoleAssertion bool) (_ *oidc.AccessTokenResponse, err error) {
	getUserInfo := s.getUserInfoOnce(session.UserID, projectID, projectRoleAssertion, session.Scope)
	getSigner := s.getSignerOnce()

	resp := &oidc.AccessTokenResponse{
		TokenType:    oidc.BearerToken,
		RefreshToken: session.RefreshToken,
		ExpiresIn:    timeToOIDCExpiresIn(session.Expiration),
		State:        state,
	}

	// If the session does not have a token ID, it is an implicit ID-Token only response.
	if session.TokenID != "" {
		if client.AccessTokenType() == op.AccessTokenTypeJWT {
			resp.AccessToken, err = s.createJWT(ctx, client, session, getUserInfo, getSigner)
		} else {
			resp.AccessToken, err = op.CreateBearerToken(session.TokenID, session.UserID, s.opCrypto)
		}
		if err != nil {
			return nil, err
		}
	}

	if slices.Contains(session.Scope, oidc.ScopeOpenID) {
		resp.IDToken, _, err = s.createIDToken(ctx, client, getUserInfo, getSigner, resp.AccessToken, session.Audience, session.AuthMethods, session.AuthTime, session.Nonce, session.Actor)
	}
	return resp, err
}

// signerFunc is a getter function that allows add-hoc retrieval of the instance's signer.
type signerFunc func(ctx context.Context) (jose.Signer, jose.SignatureAlgorithm, error)

// getSignerOnce returns a function which retrieves the instance's signer from the database once.
// Repeated calls of the returned function return the same results.
func (s *Server) getSignerOnce() signerFunc {
	var (
		once    sync.Once
		signer  jose.Signer
		signAlg jose.SignatureAlgorithm
		err     error
	)
	return func(ctx context.Context) (jose.Signer, jose.SignatureAlgorithm, error) {
		once.Do(func() {
			ctx, span := tracing.NewSpan(ctx)
			defer func() { span.EndWithError(err) }()

			var signingKey op.SigningKey
			signingKey, err = s.Provider().Storage().SigningKey(ctx)
			if err != nil {
				return
			}
			signAlg = signingKey.SignatureAlgorithm()

			signer, err = op.SignerFromKey(signingKey)
			if err != nil {
				return
			}
		})
		return signer, signAlg, err
	}
}

// userInfoFunc is a getter function that allows add-hoc retrieval of a user.
type userInfoFunc func(ctx context.Context) (*oidc.UserInfo, error)

// getUserInfoOnce returns a function which retrieves userinfo from the database once.
// Repeated calls of the returned function return the same results.
func (s *Server) getUserInfoOnce(userID, projectID string, projectRoleAssertion bool, scope []string) userInfoFunc {
	var (
		once     sync.Once
		userInfo *oidc.UserInfo
		err      error
	)
	return func(ctx context.Context) (*oidc.UserInfo, error) {
		once.Do(func() {
			ctx, span := tracing.NewSpan(ctx)
			defer func() { span.EndWithError(err) }()
			userInfo, err = s.userInfo(ctx, userID, scope, projectID, projectRoleAssertion, false)
		})
		return userInfo, err
	}
}

func (*Server) createIDToken(ctx context.Context, client op.Client, getUserInfo userInfoFunc, getSigningKey signerFunc, accessToken string, audience []string, authMethods []domain.UserAuthMethodType, authTime time.Time, nonce string, actor *domain.TokenActor) (idToken string, exp uint64, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	userInfo, err := getUserInfo(ctx)
	if err != nil {
		return "", 0, err
	}

	signer, signAlg, err := getSigningKey(ctx)
	if err != nil {
		return "", 0, err
	}

	expTime := time.Now().Add(client.IDTokenLifetime()).Add(client.ClockSkew())
	claims := oidc.NewIDTokenClaims(
		op.IssuerFromContext(ctx),
		"",
		audience,
		expTime,
		authTime,
		nonce,
		"",
		AuthMethodTypesToAMR(authMethods),
		client.GetID(),
		client.ClockSkew(),
	)
	claims.Actor = actorDomainToClaims(actor)
	claims.SetUserInfo(userInfo)
	if accessToken != "" {
		claims.AccessTokenHash, err = oidc.ClaimHash(accessToken, signAlg)
		if err != nil {
			return "", 0, err
		}
	}
	idToken, err = crypto.Sign(claims, signer)
	return idToken, timeToOIDCExpiresIn(expTime), err
}

func timeToOIDCExpiresIn(exp time.Time) uint64 {
	return uint64(time.Until(exp) / time.Second)
}

func (*Server) createJWT(ctx context.Context, client op.Client, session *command.OIDCSession, getUserInfo userInfoFunc, getSigner signerFunc) (_ string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	userInfo, err := getUserInfo(ctx)
	if err != nil {
		return "", err
	}
	signer, _, err := getSigner(ctx)
	if err != nil {
		return "", err
	}

	expTime := session.Expiration.Add(client.ClockSkew())
	claims := oidc.NewAccessTokenClaims(
		op.IssuerFromContext(ctx),
		userInfo.Subject,
		session.Audience,
		expTime,
		session.TokenID,
		client.GetID(),
		client.ClockSkew(),
	)
	claims.Actor = actorDomainToClaims(session.Actor)
	claims.Claims = userInfo.Claims

	return crypto.Sign(claims, signer)
}

// decryptCode decrypts a code or refresh_token
func (s *Server) decryptCode(ctx context.Context, code string) (_ string, err error) {
	_, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	decoded, err := base64.RawURLEncoding.DecodeString(code)
	if err != nil {
		return "", err
	}
	return s.encAlg.DecryptString(decoded, s.encAlg.EncryptionKeyID())
}

package oidc

import (
	"context"
	"strings"

	"auth/internal/oauth"
)

const (
	microsoftIssuerPrefix = "https://login.microsoftonline.com/"
	microsoftTokenURL     = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	microsoftAuthURL      = "https://login.microsoftonline.com/common/oauth2/v2.0/authorize"
)

// MicrosoftOIDCProvider implements the Provider interface for Microsoft OIDC
type MicrosoftOIDCProvider struct {
	oauthProvider *oauth.MicrosoftProvider
}

// NewMicrosoftOIDCProvider creates a new Microsoft OIDC provider
func NewMicrosoftOIDCProvider(oauthProvider *oauth.MicrosoftProvider) *MicrosoftOIDCProvider {
	return &MicrosoftOIDCProvider{
		oauthProvider: oauthProvider,
	}
}

// GetAuthURL returns the URL for authentication with Microsoft OIDC
func (p *MicrosoftOIDCProvider) GetAuthURL(state string, nonce string) string {
	additionalParams := map[string]string{
		"response_mode": "form_post",
		"nonce":         nonce,
	}

	return oauth.BuildAuthURL(microsoftAuthURL, *p.oauthProvider.Config, state, additionalParams)
}

// Exchange exchanges the authorization code for an ID toke.n
func (p *MicrosoftOIDCProvider) Exchange(ctx context.Context, code string) (*oauth.Token, error) {
	return oauth.ExchangeToken(ctx, microsoftTokenURL, *p.oauthProvider.Config, code)
}

// VerifyIDToken verifies the ID token's signature and claims for Microsoft
func (p *MicrosoftOIDCProvider) VerifyIDToken(ctx context.Context, idToken string) (*IDTokenClaims, error) {
	// Parse the ID token
	claims, err := ParseIDToken(idToken)
	if err != nil {
		return nil, err
	}

	// Microsoft issuer is dynamic based on tenant ID, so we just check the prefix
	if !strings.HasPrefix(claims.Issuer, microsoftIssuerPrefix) {
		return nil, ErrInvalidIssuer
	}

	// Validate expiration and audience
	if err := ValidateBasicClaims(claims, "", p.oauthProvider.Config.ClientID, ""); err != nil {
		return nil, err
	}

	return claims, nil
}

// GetUserInfo retrieves user information using the access token
func (p *MicrosoftOIDCProvider) GetUserInfo(ctx context.Context, token *oauth.Token) (*oauth.UserInfo, error) {
	return p.oauthProvider.GetUserInfo(ctx, token)
}

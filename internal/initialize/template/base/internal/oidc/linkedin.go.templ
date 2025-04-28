package oidc

import (
	"context"

	"{{ .Name }}/internal/oauth"
)

const (
	linkedinIssuer   = "https://www.linkedin.com"
	linkedinTokenURL = "https://www.linkedin.com/oauth/v2/accessToken"
)

// LinkedInOIDCProvider implements the Provider interface for LinkedIn OIDC
type LinkedInOIDCProvider struct {
	oauthProvider *oauth.LinkedInProvider
}

// NewLinkedInOIDCProvider creates a new LinkedIn OIDC provider
func NewLinkedInOIDCProvider(oauthProvider *oauth.LinkedInProvider) *LinkedInOIDCProvider {
	return &LinkedInOIDCProvider{
		oauthProvider: oauthProvider,
	}
}

// GetAuthURL returns the URL for authentication with LinkedIn OIDC
func (p *LinkedInOIDCProvider) GetAuthURL(state string, nonce string) string {
	additionalParams := map[string]string{
		"nonce": nonce, // Add nonce parameter for OIDC
	}
	return oauth.BuildAuthURL("https://www.linkedin.com/oauth/v2/authorization", *p.oauthProvider.Config, state, additionalParams)
}

// Exchange exchanges the authorization code for an ID token
func (p *LinkedInOIDCProvider) Exchange(ctx context.Context, code string) (*oauth.Token, error) {
	return oauth.ExchangeToken(ctx, linkedinTokenURL, *p.oauthProvider.Config, code)
}

// VerifyIDToken verifies the ID token's signature and claims for LinkedIn
func (p *LinkedInOIDCProvider) VerifyIDToken(ctx context.Context, idToken string) (*IDTokenClaims, error) {
	return verifyIDToken(ctx, idToken, linkedinIssuer, p.oauthProvider.Config.ClientID)
}

// GetUserInfo retrieves user information using the access token
func (p *LinkedInOIDCProvider) GetUserInfo(ctx context.Context, token *oauth.Token) (*oauth.UserInfo, error) {
	return p.oauthProvider.GetUserInfo(ctx, token)
}

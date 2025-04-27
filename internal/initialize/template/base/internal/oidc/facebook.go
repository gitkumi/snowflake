package oidc

import (
	"context"

	"auth/internal/oauth"
)

const (
	facebookIssuer   = "https://www.facebook.com"
	facebookTokenURL = "https://graph.facebook.com/v14.0/oauth/access_token"
)

// FacebookOIDCProvider implements the Provider interface for Facebook OIDC
type FacebookOIDCProvider struct {
	oauthProvider *oauth.FacebookProvider
}

// NewFacebookOIDCProvider creates a new Facebook OIDC provider
func NewFacebookOIDCProvider(oauthProvider *oauth.FacebookProvider) *FacebookOIDCProvider {
	return &FacebookOIDCProvider{
		oauthProvider: oauthProvider,
	}
}

// GetAuthURL returns the URL for authentication with Facebook OIDC
func (p *FacebookOIDCProvider) GetAuthURL(state string, nonce string) string {
	additionalParams := map[string]string{
		"nonce": nonce, // Add nonce parameter for OIDC
	}
	return oauth.BuildAuthURL("https://www.facebook.com/v14.0/dialog/oauth", *p.oauthProvider.Config, state, additionalParams)
}

// Exchange exchanges the authorization code for an ID token
func (p *FacebookOIDCProvider) Exchange(ctx context.Context, code string) (*oauth.Token, error) {
	return oauth.ExchangeToken(ctx, facebookTokenURL, *p.oauthProvider.Config, code)
}

// VerifyIDToken verifies the ID token's signature and claims for Facebook
func (p *FacebookOIDCProvider) VerifyIDToken(ctx context.Context, idToken string) (*IDTokenClaims, error) {
	return verifyIDToken(ctx, idToken, facebookIssuer, p.oauthProvider.Config.ClientID)
}

// GetUserInfo retrieves user information using the access token
func (p *FacebookOIDCProvider) GetUserInfo(ctx context.Context, token *oauth.Token) (*oauth.UserInfo, error) {
	return p.oauthProvider.GetUserInfo(ctx, token)
}

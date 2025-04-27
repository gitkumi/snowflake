package oidc

import (
	"context"

	"auth/internal/oauth"
)

const (
	googleIssuer   = "https://accounts.google.com"
	googleAuthURL  = "https://accounts.google.com/o/oauth2/auth"
	googleTokenURL = "https://oauth2.googleapis.com/token"
)

// GoogleOIDCProvider implements the Provider interface for Google OIDC
type GoogleOIDCProvider struct {
	oauthProvider *oauth.GoogleProvider
}

// NewGoogleOIDCProvider creates a new Google OIDC provider
func NewGoogleOIDCProvider(oauthProvider *oauth.GoogleProvider) *GoogleOIDCProvider {
	return &GoogleOIDCProvider{
		oauthProvider: oauthProvider,
	}
}

// GetAuthURL returns the URL for authentication with Google OIDC
func (p *GoogleOIDCProvider) GetAuthURL(state string, nonce string) string {
	additionalParams := map[string]string{
		"access_type": "offline", // Request a refresh token
		"prompt":      "consent",
		"nonce":       nonce, // Add nonce parameter for OIDC
	}

	return oauth.BuildAuthURL(googleAuthURL, *p.oauthProvider.Config, state, additionalParams)
}

// Exchange exchanges the authorization code for an ID token
func (p *GoogleOIDCProvider) Exchange(ctx context.Context, code string) (*oauth.Token, error) {
	return oauth.ExchangeToken(ctx, googleTokenURL, *p.oauthProvider.Config, code)
}

// VerifyIDToken verifies the ID token's signature and claims
func (p *GoogleOIDCProvider) VerifyIDToken(ctx context.Context, idToken string) (*IDTokenClaims, error) {
	return verifyIDToken(ctx, idToken, googleIssuer, p.oauthProvider.Config.ClientID)
}

// GetUserInfo retrieves user information using the access token
func (p *GoogleOIDCProvider) GetUserInfo(ctx context.Context, token *oauth.Token) (*oauth.UserInfo, error) {
	return p.oauthProvider.GetUserInfo(ctx, token)
}

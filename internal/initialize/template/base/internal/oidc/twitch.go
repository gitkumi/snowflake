package oidc

import (
	"context"

	"auth/internal/oauth"
)

const (
	twitchIssuer   = "https://id.twitch.tv/oauth2"
	twitchTokenURL = "https://id.twitch.tv/oauth2/token"
)

// TwitchOIDCProvider implements the Provider interface for Twitch OIDC
type TwitchOIDCProvider struct {
	oauthProvider *oauth.TwitchProvider
}

// NewTwitchOIDCProvider creates a new Twitch OIDC provider
func NewTwitchOIDCProvider(oauthProvider *oauth.TwitchProvider) *TwitchOIDCProvider {
	return &TwitchOIDCProvider{
		oauthProvider: oauthProvider,
	}
}

// GetAuthURL returns the URL for authentication with Twitch OIDC
func (p *TwitchOIDCProvider) GetAuthURL(state string, nonce string) string {
	additionalParams := map[string]string{
		"force_verify": "true", // Always ask for verification
		"nonce":        nonce,  // Add nonce parameter for OIDC
	}

	return oauth.BuildAuthURL("https://id.twitch.tv/oauth2/authorize", *p.oauthProvider.Config, state, additionalParams)
}

// Exchange exchanges the authorization code for an ID token
func (p *TwitchOIDCProvider) Exchange(ctx context.Context, code string) (*oauth.Token, error) {
	return oauth.ExchangeToken(ctx, twitchTokenURL, *p.oauthProvider.Config, code)
}

// VerifyIDToken verifies the ID token's signature and claims for Twitch
func (p *TwitchOIDCProvider) VerifyIDToken(ctx context.Context, idToken string) (*IDTokenClaims, error) {
	return verifyIDToken(ctx, idToken, twitchIssuer, p.oauthProvider.Config.ClientID)
}

// GetUserInfo retrieves user information using the access token
func (p *TwitchOIDCProvider) GetUserInfo(ctx context.Context, token *oauth.Token) (*oauth.UserInfo, error) {
	return p.oauthProvider.GetUserInfo(ctx, token)
}

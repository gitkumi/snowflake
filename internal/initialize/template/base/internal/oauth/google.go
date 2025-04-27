package oauth

import (
	"context"
	"fmt"

	"auth/util"
)

const (
	googleAuthURL     = "https://accounts.google.com/o/oauth2/auth"
	googleTokenURL    = "https://oauth2.googleapis.com/token"
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"
)

// GoogleProvider implements the Provider interface for Google OAuth2
type GoogleProvider struct {
	Config *ProviderConfig
}

// NewGoogleProvider creates a new Google OAuth2 provider
func NewGoogleProvider(cfg *ProviderConfig) *GoogleProvider {
	return &GoogleProvider{
		Config: cfg,
	}
}

// Name returns the provider name
func (p *GoogleProvider) Name() string {
	return "google"
}

// GetAuthURL returns the Google OAuth2 authorization URL
func (p *GoogleProvider) GetAuthURL(state string) string {
	additionalParams := map[string]string{
		"access_type": "offline", // Request a refresh token
		"prompt":      "consent",
	}
	return BuildAuthURL(googleAuthURL, *p.Config, state, additionalParams)
}

// Exchange exchanges the authorization code for an access token
func (p *GoogleProvider) Exchange(ctx context.Context, code string) (*Token, error) {
	return ExchangeToken(ctx, googleTokenURL, *p.Config, code)
}

// GetUserInfo retrieves the user's information using the access token
func (p *GoogleProvider) GetUserInfo(ctx context.Context, token *Token) (*UserInfo, error) {
	var userInfoResp struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Locale        string `json:"locale"`
	}

	if err := util.MakeAuthenticatedRequest(ctx, "GET", googleUserInfoURL, token.AccessToken, &userInfoResp); err != nil {
		return nil, fmt.Errorf("failed to get Google user info: %w", err)
	}

	return &UserInfo{
		ID:            userInfoResp.Sub,
		Email:         userInfoResp.Email,
		VerifiedEmail: userInfoResp.EmailVerified,
		Name:          userInfoResp.Name,
		GivenName:     userInfoResp.GivenName,
		FamilyName:    userInfoResp.FamilyName,
		Picture:       userInfoResp.Picture,
		Locale:        userInfoResp.Locale,
		ProviderName:  p.Name(),
	}, nil
}

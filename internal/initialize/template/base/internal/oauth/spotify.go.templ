package oauth

import (
	"context"
	"fmt"

	"auth/util"
)

const (
	spotifyAuthURL     = "https://accounts.spotify.com/authorize"
	spotifyTokenURL    = "https://accounts.spotify.com/api/token"
	spotifyUserInfoURL = "https://api.spotify.com/v1/me"
)

// SpotifyProvider implements the Provider interface for Spotify OAuth2
type SpotifyProvider struct {
	Config *ProviderConfig
}

// NewSpotifyProvider creates a new Spotify OAuth2 provider
func NewSpotifyProvider(cfg *ProviderConfig) *SpotifyProvider {
	return &SpotifyProvider{
		Config: cfg,
	}
}

// Name returns the provider name
func (p *SpotifyProvider) Name() string {
	return "spotify"
}

// GetAuthURL returns the Spotify OAuth2 authorization URL
func (p *SpotifyProvider) GetAuthURL(state string) string {
	additionalParams := map[string]string{
		"show_dialog": "true", // Force user to approve the app each time
	}
	return BuildAuthURL(spotifyAuthURL, *p.Config, state, additionalParams)
}

// Exchange exchanges the authorization code for an access token
func (p *SpotifyProvider) Exchange(ctx context.Context, code string) (*Token, error) {
	// Spotify uses Basic Auth with client_id and client_secret
	return ExchangeTokenWithBasicAuth(ctx, spotifyTokenURL, *p.Config, code)
}

// GetUserInfo retrieves the user's information using the access token
func (p *SpotifyProvider) GetUserInfo(ctx context.Context, token *Token) (*UserInfo, error) {
	var userResp struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Display string `json:"display_name"`
		Images  []struct {
			URL string `json:"url"`
		} `json:"images"`
		Country         string `json:"country"`
		ExplicitContent struct {
			FilterEnabled bool `json:"filter_enabled"`
			FilterLocked  bool `json:"filter_locked"`
		} `json:"explicit_content"`
		ExternalURLs struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Followers struct {
			Total int `json:"total"`
		} `json:"followers"`
		Product string `json:"product"`
	}

	if err := util.MakeAuthenticatedRequest(ctx, "GET", spotifyUserInfoURL, token.AccessToken, &userResp); err != nil {
		return nil, fmt.Errorf("failed to get Spotify user info: %w", err)
	}

	// Get profile image if available
	picture := ""
	if len(userResp.Images) > 0 {
		picture = userResp.Images[0].URL
	}

	// Spotify doesn't provide first name and last name separately
	name := userResp.Display

	return &UserInfo{
		ID:            userResp.ID,
		Email:         userResp.Email,
		VerifiedEmail: userResp.Email != "", // If email is provided, it's verified
		Name:          name,
		Picture:       picture,
		Locale:        userResp.Country,
		ProviderName:  p.Name(),
		// Spotify doesn't provide these fields
		GivenName:  "",
		FamilyName: "",
	}, nil
}

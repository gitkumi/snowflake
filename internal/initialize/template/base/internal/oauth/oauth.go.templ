package oauth

import (
	"auth/util"
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// ProviderConfig contains the configuration for a specific OAuth provider
type ProviderConfig struct {
	Name         string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// Provider defines the interface for all OAuth providers
type Provider interface {
	// Name returns the name of the provider
	Name() string

	// GetAuthURL returns the URL to redirect the user to for authentication
	GetAuthURL(state string) string

	// Exchange exchanges the authorization code for an access token
	Exchange(ctx context.Context, code string) (*Token, error)

	// GetUserInfo retrieves user information using the access token
	GetUserInfo(ctx context.Context, token *Token) (*UserInfo, error)
}

// Token represents an OAuth access token
type Token struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       int64
	// For OpenID Connect providers like Google
	IDToken string
}

// UserInfo represents the standardized user information across providers
type UserInfo struct {
	ID            string
	Email         string
	VerifiedEmail bool
	Name          string
	GivenName     string
	FamilyName    string
	Picture       string
	Locale        string
	ProviderName  string
}

// OAuth maintains a registry of available OAuth providers
type OAuth struct {
	providers map[string]Provider
}

// NewOAuth creates a new provider registry with configured providers
func NewOAuth(configs []*ProviderConfig) *OAuth {
	registry := &OAuth{
		providers: make(map[string]Provider),
	}

	// Register each provider
	for _, config := range configs {
		switch config.Name {
		case "google":
			registry.providers[config.Name] = NewGoogleProvider(config)
		case "github":
			registry.providers[config.Name] = NewGitHubProvider(config)
		case "facebook":
			registry.providers[config.Name] = NewFacebookProvider(config)
		case "instagram":
			registry.providers[config.Name] = NewInstagramProvider(config)
		case "discord":
			registry.providers[config.Name] = NewDiscordProvider(config)
		case "linkedin":
			registry.providers[config.Name] = NewLinkedInProvider(config)
		case "reddit":
			registry.providers[config.Name] = NewRedditProvider(config)
		case "twitch":
			registry.providers[config.Name] = NewTwitchProvider(config)
		case "stripe":
			registry.providers[config.Name] = NewStripeProvider(config)
		case "x":
			registry.providers[config.Name] = NewXProvider(config)
		case "microsoft":
			registry.providers[config.Name] = NewMicrosoftProvider(config)
		case "slack":
			registry.providers[config.Name] = NewSlackProvider(config)
		case "spotify":
			registry.providers[config.Name] = NewSpotifyProvider(config)
		}
	}

	return registry
}

// Get retrieves a provider by name
func (r *OAuth) Get(name string) (Provider, bool) {
	provider, exists := r.providers[name]
	return provider, exists
}

// List returns all registered providers
func (r *OAuth) List() map[string]Provider {
	return r.providers
}

// BuildAuthURL creates an authorization URL with the appropriate parameters
func BuildAuthURL(authURL string, config ProviderConfig, state string, additionalParams map[string]string) string {
	v := url.Values{
		"client_id":     {config.ClientID},
		"redirect_uri":  {config.RedirectURL},
		"response_type": {"code"},
		"scope":         {strings.Join(config.Scopes, " ")},
		"state":         {state},
	}

	// Add any additional parameters
	for key, value := range additionalParams {
		v.Add(key, value)
	}

	return fmt.Sprintf("%s?%s", authURL, v.Encode())
}

// ExchangeToken exchanges the authorization code for an access token
// using standard OAuth2 parameters
func ExchangeToken(ctx context.Context, tokenURL string, config ProviderConfig, code string) (*Token, error) {
	form := url.Values{
		"code":          {code},
		"client_id":     {config.ClientID},
		"client_secret": {config.ClientSecret},
		"redirect_uri":  {config.RedirectURL},
		"grant_type":    {"authorization_code"},
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
		IDToken      string `json:"id_token"`
	}

	if err := util.MakeFormRequest(ctx, tokenURL, form, &tokenResp); err != nil {
		return nil, fmt.Errorf("oauth token exchange failed: %w", err)
	}

	return &Token{
		AccessToken:  tokenResp.AccessToken,
		TokenType:    tokenResp.TokenType,
		RefreshToken: tokenResp.RefreshToken,
		Expiry:       time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Unix(),
		IDToken:      tokenResp.IDToken,
	}, nil
}

// ExchangeTokenWithBasicAuth exchanges the authorization code for tokens
// using Basic Auth for client authentication (used by some providers like PayPal)
func ExchangeTokenWithBasicAuth(ctx context.Context, tokenURL string, config ProviderConfig, code string) (*Token, error) {
	form := url.Values{
		"code":         {code},
		"grant_type":   {"authorization_code"},
		"redirect_uri": {config.RedirectURL},
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token,omitempty"`
		IDToken      string `json:"id_token,omitempty"`
		Scope        string `json:"scope,omitempty"`
	}

	if err := util.MakeBasicAuthRequest(ctx, "POST", tokenURL, config.ClientID, config.ClientSecret, form, &tokenResp); err != nil {
		return nil, fmt.Errorf("oauth token exchange failed: %w", err)
	}

	return &Token{
		AccessToken:  tokenResp.AccessToken,
		TokenType:    tokenResp.TokenType,
		RefreshToken: tokenResp.RefreshToken,
		Expiry:       time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Unix(),
		IDToken:      tokenResp.IDToken,
	}, nil
}

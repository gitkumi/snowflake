package oidc

import (
	"context"

	"{{ .Name }}/internal/oauth"
)

// Provider defines the interface for all OIDC providers
type Provider interface {
	// GetAuthURL returns the URL for authentication with OIDC
	GetAuthURL(state string, nonce string) string

	// Exchange exchanges the authorization code for token and ID token
	Exchange(ctx context.Context, code string) (*oauth.Token, error)

	// GetUserInfo retrieves user information using the access token
	GetUserInfo(ctx context.Context, token *oauth.Token) (*oauth.UserInfo, error)

	// VerifyIDToken verifies the ID token's signature and claims
	VerifyIDToken(ctx context.Context, idToken string) (*IDTokenClaims, error)
}

// IDTokenClaims contains the claims from an ID token
type IDTokenClaims struct {
	// Standard claims
	Issuer         string `json:"iss"`
	Subject        string `json:"sub"`
	Audience       string `json:"aud"`
	ExpirationTime int64  `json:"exp"`
	IssuedAt       int64  `json:"iat"`
	AuthTime       int64  `json:"auth_time,omitempty"`
	Nonce          string `json:"nonce,omitempty"`

	// Profile claims
	Name              string `json:"name,omitempty"`
	GivenName         string `json:"given_name,omitempty"`
	FamilyName        string `json:"family_name,omitempty"`
	MiddleName        string `json:"middle_name,omitempty"`
	Nickname          string `json:"nickname,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	Profile           string `json:"profile,omitempty"`
	Picture           string `json:"picture,omitempty"`
	Website           string `json:"website,omitempty"`
	Gender            string `json:"gender,omitempty"`
	Birthdate         string `json:"birthdate,omitempty"`
	Zoneinfo          string `json:"zoneinfo,omitempty"`
	Locale            string `json:"locale,omitempty"`
	UpdatedAt         int64  `json:"updated_at,omitempty"`

	// Email claims
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`

	// Additional non-standard claims
	AdditionalClaims map[string]interface{} `json:"-"`

	// Raw token
	RawToken string `json:"-"`
}

// OIDC maintains a registry of available OIDC providers
type OIDC struct {
	providers map[string]Provider
}

func NewOIDC(configs []*oauth.ProviderConfig) *OIDC {
	registry := &OIDC{
		providers: make(map[string]Provider),
	}

	for _, config := range configs {
		providerName := config.Name

		switch providerName {
		case "google":
			googleProvider := oauth.NewGoogleProvider(config)
			registry.providers[providerName] = NewGoogleOIDCProvider(googleProvider)
		case "microsoft":
			microsoftProvider := oauth.NewMicrosoftProvider(config)
			registry.providers[providerName] = NewMicrosoftOIDCProvider(microsoftProvider)
		case "facebook":
			facebookProvider := oauth.NewFacebookProvider(config)
			registry.providers[providerName] = NewFacebookOIDCProvider(facebookProvider)
		case "linkedin":
			linkedinProvider := oauth.NewLinkedInProvider(config)
			registry.providers[providerName] = NewLinkedInOIDCProvider(linkedinProvider)
		case "twitch":
			twitchProvider := oauth.NewTwitchProvider(config)
			registry.providers[providerName] = NewTwitchOIDCProvider(twitchProvider)
		}
	}

	return registry
}

// Get retrieves a provider by name
func (r *OIDC) Get(name string) (Provider, bool) {
	provider, exists := r.providers[name]
	return provider, exists
}

// List returns all registered providers
func (r *OIDC) List() map[string]Provider {
	return r.providers
}

// Verify ID Tokens
func verifyIDToken(ctx context.Context, idToken string, issuer string, clientID string) (*IDTokenClaims, error) {
	// Parse the ID token
	claims, err := ParseIDToken(idToken)
	if err != nil {
		return nil, err
	}

	// Validate basic claims
	if err := ValidateBasicClaims(claims, issuer, clientID, ""); err != nil {
		return nil, err
	}

	return claims, nil
}

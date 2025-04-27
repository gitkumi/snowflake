package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	stripeAuthURL     = "https://connect.stripe.com/oauth/authorize"
	stripeTokenURL    = "https://connect.stripe.com/oauth/token"
	stripeUserInfoURL = "https://api.stripe.com/v1/accounts"
)

// StripeProvider implements the Provider interface for Stripe OAuth2
type StripeProvider struct {
	Config *ProviderConfig
}

// NewStripeProvider creates a new Stripe OAuth2 provider
func NewStripeProvider(cfg *ProviderConfig) *StripeProvider {
	return &StripeProvider{
		Config: cfg,
	}
}

// Name returns the provider name
func (p *StripeProvider) Name() string {
	return "stripe"
}

// GetAuthURL returns the Stripe OAuth2 authorization URL
func (p *StripeProvider) GetAuthURL(state string) string {
	v := url.Values{
		"client_id":     {p.Config.ClientID},
		"redirect_uri":  {p.Config.RedirectURL},
		"response_type": {"code"},
		"scope":         {strings.Join(p.Config.Scopes, " ")},
		"state":         {state},
	}

	return fmt.Sprintf("%s?%s", stripeAuthURL, v.Encode())
}

// Exchange exchanges the authorization code for an access token
func (p *StripeProvider) Exchange(ctx context.Context, code string) (*Token, error) {
	form := url.Values{
		"code":          {code},
		"client_secret": {p.Config.ClientSecret},
		"grant_type":    {"authorization_code"},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", stripeTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("oauth token exchange failed: %s", body)
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
		StripeUserID string `json:"stripe_user_id"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &Token{
		AccessToken: tokenResp.AccessToken,
		TokenType:   tokenResp.TokenType,
		// Stripe Connect tokens don't expire
	}, nil
}

// GetUserInfo retrieves the connected account information using the access token
func (p *StripeProvider) GetUserInfo(ctx context.Context, token *Token) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", stripeUserInfoURL+"/"+p.getCurrentAccountID(token), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform user info request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %s", body)
	}

	var accountResp struct {
		ID              string `json:"id"`
		Email           string `json:"email"`
		DisplayName     string `json:"display_name"`
		BusinessProfile struct {
			Name string `json:"name"`
			URL  string `json:"url"`
			Logo string `json:"logo"`
		} `json:"business_profile"`
		Settings struct {
			Dashboard struct {
				DisplayName string `json:"display_name"`
			} `json:"dashboard"`
		} `json:"settings"`
	}

	if err := json.Unmarshal(body, &accountResp); err != nil {
		return nil, fmt.Errorf("failed to parse user info response: %w", err)
	}

	// Use the most descriptive name available
	name := accountResp.DisplayName
	if name == "" {
		name = accountResp.BusinessProfile.Name
	}
	if name == "" {
		name = accountResp.Settings.Dashboard.DisplayName
	}

	return &UserInfo{
		ID:           accountResp.ID,
		Email:        accountResp.Email,
		Name:         name,
		Picture:      accountResp.BusinessProfile.Logo,
		ProviderName: p.Name(),
		// Stripe doesn't provide these fields in the standard response
		VerifiedEmail: true, // Stripe accounts have verified emails
		GivenName:     "",
		FamilyName:    "",
		Locale:        "",
	}, nil
}

// getCurrentAccountID returns the account ID from the Stripe API
func (p *StripeProvider) getCurrentAccountID(token *Token) string {
	// We need to get the current account ID first
	req, err := http.NewRequest("GET", stripeUserInfoURL+"/me", nil)
	if err != nil {
		return ""
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var accountResp struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(body, &accountResp); err != nil {
		return ""
	}

	return accountResp.ID
}

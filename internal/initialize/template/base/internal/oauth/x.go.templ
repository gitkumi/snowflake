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
	xAuthURL     = "https://twitter.com/i/oauth2/authorize"
	xTokenURL    = "https://api.twitter.com/2/oauth2/token"
	xUserInfoURL = "https://api.twitter.com/2/users/me"
)

// XProvider implements the Provider interface for X (Twitter) OAuth2
type XProvider struct {
	Config *ProviderConfig
}

// NewXProvider creates a new X (Twitter) OAuth2 provider
func NewXProvider(cfg *ProviderConfig) *XProvider {
	return &XProvider{
		Config: cfg,
	}
}

// Name returns the provider name
func (p *XProvider) Name() string {
	return "x"
}

// GetAuthURL returns the X (Twitter) OAuth2 authorization URL
func (p *XProvider) GetAuthURL(state string) string {
	v := url.Values{
		"client_id":             {p.Config.ClientID},
		"redirect_uri":          {p.Config.RedirectURL},
		"response_type":         {"code"},
		"scope":                 {strings.Join(p.Config.Scopes, " ")},
		"state":                 {state},
		"code_challenge_method": {"S256"},
		"code_challenge":        {"challenge"}, // This should be generated using PKCE in a real implementation
	}

	return fmt.Sprintf("%s?%s", xAuthURL, v.Encode())
}

// Exchange exchanges the authorization code for an access token
func (p *XProvider) Exchange(ctx context.Context, code string) (*Token, error) {
	form := url.Values{
		"code":          {code},
		"client_id":     {p.Config.ClientID},
		"client_secret": {p.Config.ClientSecret},
		"redirect_uri":  {p.Config.RedirectURL},
		"grant_type":    {"authorization_code"},
		"code_verifier": {"verifier"}, // This should match the challenge used in GetAuthURL in a real implementation
	}

	req, err := http.NewRequestWithContext(ctx, "POST", xTokenURL, strings.NewReader(form.Encode()))
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
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &Token{
		AccessToken:  tokenResp.AccessToken,
		TokenType:    tokenResp.TokenType,
		RefreshToken: tokenResp.RefreshToken,
		Expiry:       time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Unix(),
	}, nil
}

// GetUserInfo retrieves the user's information using the access token
func (p *XProvider) GetUserInfo(ctx context.Context, token *Token) (*UserInfo, error) {
	// Twitter API v2 requires specifying fields to include
	params := url.Values{
		"user.fields": {"id,name,username,profile_image_url,verified,location,description,entities,public_metrics"},
	}

	reqURL := fmt.Sprintf("%s?%s", xUserInfoURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
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

	var userResp struct {
		Data struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			Username        string `json:"username"`
			ProfileImageURL string `json:"profile_image_url"`
			Verified        bool   `json:"verified"`
			Location        string `json:"location"`
			Description     string `json:"description"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, fmt.Errorf("failed to parse user info response: %w", err)
	}

	// Extract first and last name if possible (X doesn't provide these separately)
	nameParts := strings.Split(userResp.Data.Name, " ")
	firstName := ""
	lastName := ""

	if len(nameParts) > 0 {
		firstName = nameParts[0]
	}
	if len(nameParts) > 1 {
		lastName = strings.Join(nameParts[1:], " ")
	}

	return &UserInfo{
		ID:            userResp.Data.ID,
		Name:          userResp.Data.Name,
		GivenName:     firstName,
		FamilyName:    lastName,
		Picture:       userResp.Data.ProfileImageURL,
		VerifiedEmail: userResp.Data.Verified,
		ProviderName:  p.Name(),
		// X (Twitter) doesn't provide these fields in its basic API
		Email:  "",
		Locale: userResp.Data.Location,
	}, nil
}

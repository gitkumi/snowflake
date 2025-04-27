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
	slackAuthURL     = "https://slack.com/oauth/v2/authorize"
	slackTokenURL    = "https://slack.com/api/oauth.v2.access"
	slackUserInfoURL = "https://slack.com/api/users.identity"
)

// SlackProvider implements the Provider interface for Slack OAuth2
type SlackProvider struct {
	Config *ProviderConfig
}

// NewSlackProvider creates a new Slack OAuth2 provider
func NewSlackProvider(cfg *ProviderConfig) *SlackProvider {
	return &SlackProvider{
		Config: cfg,
	}
}

// Name returns the provider name
func (p *SlackProvider) Name() string {
	return "slack"
}

// GetAuthURL returns the Slack OAuth2 authorization URL
func (p *SlackProvider) GetAuthURL(state string) string {
	v := url.Values{
		"client_id":    {p.Config.ClientID},
		"redirect_uri": {p.Config.RedirectURL},
		"scope":        {strings.Join(p.Config.Scopes, ",")},
		"state":        {state},
		"user_scope":   {"identity.basic,identity.email,identity.avatar"},
	}

	return fmt.Sprintf("%s?%s", slackAuthURL, v.Encode())
}

// Exchange exchanges the authorization code for an access token
func (p *SlackProvider) Exchange(ctx context.Context, code string) (*Token, error) {
	form := url.Values{
		"code":          {code},
		"client_id":     {p.Config.ClientID},
		"client_secret": {p.Config.ClientSecret},
		"redirect_uri":  {p.Config.RedirectURL},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", slackTokenURL, strings.NewReader(form.Encode()))
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
		Ok          bool   `json:"ok"`
		Error       string `json:"error,omitempty"`
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
		UserID      string `json:"user_id"`
		TeamID      string `json:"team_id"`
		Team        struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"team"`
		Authed_User struct {
			ID          string `json:"id"`
			AccessToken string `json:"access_token"`
			TokenType   string `json:"token_type"`
			Scope       string `json:"scope"`
		} `json:"authed_user"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	if !tokenResp.Ok {
		return nil, fmt.Errorf("slack oauth error: %s", tokenResp.Error)
	}

	// For Slack, we use the user token for user info
	return &Token{
		AccessToken: tokenResp.Authed_User.AccessToken,
		TokenType:   tokenResp.Authed_User.TokenType,
		// Slack tokens don't expire by default
	}, nil
}

// GetUserInfo retrieves the user's information using the access token
func (p *SlackProvider) GetUserInfo(ctx context.Context, token *Token) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", slackUserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))

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
		Ok    bool   `json:"ok"`
		Error string `json:"error,omitempty"`
		User  struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			Image_48  string `json:"image_48"`
			Image_72  string `json:"image_72"`
			Image_192 string `json:"image_192"`
			Image_512 string `json:"image_512"`
		} `json:"user"`
		Team struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"team"`
	}

	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, fmt.Errorf("failed to parse user info response: %w", err)
	}

	if !userResp.Ok {
		return nil, fmt.Errorf("slack user info error: %s", userResp.Error)
	}

	// Use the highest resolution image available
	picture := userResp.User.Image_512
	if picture == "" {
		picture = userResp.User.Image_192
	}
	if picture == "" {
		picture = userResp.User.Image_72
	}
	if picture == "" {
		picture = userResp.User.Image_48
	}

	return &UserInfo{
		ID:            userResp.User.ID,
		Email:         userResp.User.Email,
		VerifiedEmail: true, // Slack requires verified emails
		Name:          userResp.User.Name,
		Picture:       picture,
		ProviderName:  p.Name(),
		// Slack doesn't provide these fields
		GivenName:  "",
		FamilyName: "",
		Locale:     "",
	}, nil
}

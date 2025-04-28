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
	microsoftAuthURL     = "https://login.microsoftonline.com/common/oauth2/v2.0/authorize"
	microsoftTokenURL    = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	microsoftUserInfoURL = "https://graph.microsoft.com/v1.0/me"
	microsoftPhotoURL    = "https://graph.microsoft.com/v1.0/me/photo/$value"
)

// MicrosoftProvider implements the Provider interface for Microsoft OAuth2
type MicrosoftProvider struct {
	Config *ProviderConfig
}

// NewMicrosoftProvider creates a new Microsoft OAuth2 provider
func NewMicrosoftProvider(cfg *ProviderConfig) *MicrosoftProvider {
	return &MicrosoftProvider{
		Config: cfg,
	}
}

// Name returns the provider name
func (p *MicrosoftProvider) Name() string {
	return "microsoft"
}

// GetAuthURL returns the Microsoft OAuth2 authorization URL
func (p *MicrosoftProvider) GetAuthURL(state string) string {
	v := url.Values{
		"client_id":     {p.Config.ClientID},
		"redirect_uri":  {p.Config.RedirectURL},
		"response_type": {"code"},
		"scope":         {strings.Join(p.Config.Scopes, " ")},
		"state":         {state},
		"response_mode": {"query"},
	}
	return microsoftAuthURL + "?" + v.Encode()
}

// Exchange exchanges the authorization code for an access token
func (p *MicrosoftProvider) Exchange(ctx context.Context, code string) (*Token, error) {
	form := url.Values{
		"code":          {code},
		"client_id":     {p.Config.ClientID},
		"client_secret": {p.Config.ClientSecret},
		"redirect_uri":  {p.Config.RedirectURL},
		"grant_type":    {"authorization_code"},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", microsoftTokenURL, strings.NewReader(form.Encode()))
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
		IDToken      string `json:"id_token"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &Token{
		AccessToken:  tokenResp.AccessToken,
		TokenType:    tokenResp.TokenType,
		RefreshToken: tokenResp.RefreshToken,
		Expiry:       time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Unix(),
		IDToken:      tokenResp.IDToken,
	}, nil
}

// GetUserInfo retrieves the user's information using the access token
func (p *MicrosoftProvider) GetUserInfo(ctx context.Context, token *Token) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", microsoftUserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))
	req.Header.Set("Accept", "application/json")

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
		ID                string `json:"id"`
		DisplayName       string `json:"displayName"`
		GivenName         string `json:"givenName"`
		Surname           string `json:"surname"`
		Email             string `json:"mail"`
		UserPrincipalName string `json:"userPrincipalName"`
		PreferredLanguage string `json:"preferredLanguage"`
	}

	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, fmt.Errorf("failed to parse user info response: %w", err)
	}

	// Use the profile picture if available
	picture := ""
	photoReq, err := http.NewRequestWithContext(ctx, "GET", microsoftPhotoURL, nil)
	if err == nil {
		photoReq.Header.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))

		photoResp, err := client.Do(photoReq)
		if err == nil && photoResp.StatusCode == http.StatusOK {
			// In a production app, you might want to convert the photo to a data URL or store it
			// For simplicity, we'll just set a placeholder URL
			picture = fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/photo/$value", userResp.ID)
			photoResp.Body.Close()
		}
	}

	// If email is not in the mail field, use userPrincipalName
	email := userResp.Email
	if email == "" {
		email = userResp.UserPrincipalName
	}

	return &UserInfo{
		ID:            userResp.ID,
		Email:         email,
		VerifiedEmail: true, // Microsoft accounts have verified emails
		Name:          userResp.DisplayName,
		GivenName:     userResp.GivenName,
		FamilyName:    userResp.Surname,
		Picture:       picture,
		Locale:        userResp.PreferredLanguage,
		ProviderName:  p.Name(),
	}, nil
}

package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"auth/util"
)

const (
	redditAuthURL     = "https://www.reddit.com/api/v1/authorize"
	redditTokenURL    = "https://www.reddit.com/api/v1/access_token"
	redditUserInfoURL = "https://oauth.reddit.com/api/v1/me"
)

// RedditProvider implements the Provider interface for Reddit OAuth2
type RedditProvider struct {
	Config *ProviderConfig
}

// NewRedditProvider creates a new Reddit OAuth2 provider
func NewRedditProvider(cfg *ProviderConfig) *RedditProvider {
	return &RedditProvider{
		Config: cfg,
	}
}

// Name returns the provider name
func (p *RedditProvider) Name() string {
	return "reddit"
}

// GetAuthURL returns the Reddit OAuth2 authorization URL
func (p *RedditProvider) GetAuthURL(state string) string {
	additionalParams := map[string]string{
		"duration": "permanent", // Ask for a refresh token
	}
	return BuildAuthURL(redditAuthURL, *p.Config, state, additionalParams)
}

// Exchange exchanges the authorization code for an access token
func (p *RedditProvider) Exchange(ctx context.Context, code string) (*Token, error) {
	// Reddit requires Basic Auth with client_id and client_secret
	return ExchangeTokenWithBasicAuth(ctx, redditTokenURL, *p.Config, code)
}

// GetUserInfo retrieves the user's information using the access token
func (p *RedditProvider) GetUserInfo(ctx context.Context, token *Token) (*UserInfo, error) {
	var userResp struct {
		ID       string  `json:"id"`
		Name     string  `json:"name"`
		IconImg  string  `json:"icon_img"`
		Created  float64 `json:"created"`
		Verified bool    `json:"verified"`
	}

	// Reddit requires both a User-Agent header and the Authorization header
	// We'll need to make the request manually
	req, err := http.NewRequestWithContext(ctx, "GET", redditUserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Reddit user info request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))
	req.Header.Set("User-Agent", "AuthApplication/1.0") // Reddit requires a User-Agent

	client := util.DefaultHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform Reddit user info request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Reddit user info response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get Reddit user info: %s", body)
	}

	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, fmt.Errorf("failed to parse Reddit user info response: %w", err)
	}

	// Clean the icon URL if it has any query parameters
	iconURL := userResp.IconImg
	if strings.Contains(iconURL, "?") {
		iconURL = strings.Split(iconURL, "?")[0]
	}

	return &UserInfo{
		ID:            userResp.ID,
		Name:          userResp.Name,
		Picture:       iconURL,
		ProviderName:  p.Name(),
		VerifiedEmail: userResp.Verified,
		// Reddit doesn't provide these fields with basic permissions
		Email:      "",
		GivenName:  "",
		FamilyName: "",
		Locale:     "",
	}, nil
}

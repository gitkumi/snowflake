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
	linkedinAuthURL     = "https://www.linkedin.com/oauth/v2/authorization"
	linkedinTokenURL    = "https://www.linkedin.com/oauth/v2/accessToken"
	linkedinUserInfoURL = "https://api.linkedin.com/v2/me"
	linkedinEmailURL    = "https://api.linkedin.com/v2/emailAddress?q=members&projection=(elements*(handle~))"
)

// LinkedInProvider implements the Provider interface for LinkedIn OAuth2
type LinkedInProvider struct {
	Config *ProviderConfig
}

// NewLinkedInProvider creates a new LinkedIn OAuth2 provider
func NewLinkedInProvider(cfg *ProviderConfig) *LinkedInProvider {
	return &LinkedInProvider{
		Config: cfg,
	}
}

// Name returns the provider name
func (p *LinkedInProvider) Name() string {
	return "linkedin"
}

// GetAuthURL returns the LinkedIn OAuth2 authorization URL
func (p *LinkedInProvider) GetAuthURL(state string) string {
	v := url.Values{
		"client_id":     {p.Config.ClientID},
		"redirect_uri":  {p.Config.RedirectURL},
		"response_type": {"code"},
		"scope":         {strings.Join(p.Config.Scopes, " ")},
		"state":         {state},
	}

	return fmt.Sprintf("%s?%s", linkedinAuthURL, v.Encode())
}

// Exchange exchanges the authorization code for an access token
func (p *LinkedInProvider) Exchange(ctx context.Context, code string) (*Token, error) {
	form := url.Values{
		"code":          {code},
		"client_id":     {p.Config.ClientID},
		"client_secret": {p.Config.ClientSecret},
		"redirect_uri":  {p.Config.RedirectURL},
		"grant_type":    {"authorization_code"},
	}

	req, err := http.NewRequestWithContext(ctx, "POST", linkedinTokenURL, strings.NewReader(form.Encode()))
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
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &Token{
		AccessToken: tokenResp.AccessToken,
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Unix(),
	}, nil
}

// GetUserInfo retrieves the user's information using the access token
func (p *LinkedInProvider) GetUserInfo(ctx context.Context, token *Token) (*UserInfo, error) {
	// First get basic profile info
	userInfo, err := p.getUserProfile(ctx, token)
	if err != nil {
		return nil, err
	}

	// Then get email
	if err := p.enrichWithEmail(ctx, token, userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

// getUserProfile fetches the LinkedIn user profile
func (p *LinkedInProvider) getUserProfile(ctx context.Context, token *Token) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", linkedinUserInfoURL, nil)
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
		ID        string `json:"id"`
		FirstName struct {
			Localized       map[string]string `json:"localized"`
			PreferredLocale struct {
				Country  string `json:"country"`
				Language string `json:"language"`
			} `json:"preferredLocale"`
		} `json:"firstName"`
		LastName struct {
			Localized       map[string]string `json:"localized"`
			PreferredLocale struct {
				Country  string `json:"country"`
				Language string `json:"language"`
			} `json:"preferredLocale"`
		} `json:"lastName"`
		ProfilePicture struct {
			DisplayImage struct {
				Elements []struct {
					Identifiers []struct {
						Identifier string `json:"identifier"`
					} `json:"identifiers"`
				} `json:"elements"`
			} `json:"displayImage~"`
		} `json:"profilePicture"`
	}

	if err := json.Unmarshal(body, &userResp); err != nil {
		return nil, fmt.Errorf("failed to parse user info response: %w", err)
	}

	// Get localized first and last name using preferred locale
	locale := fmt.Sprintf("%s_%s", userResp.FirstName.PreferredLocale.Language, userResp.FirstName.PreferredLocale.Country)
	firstName := ""
	lastName := ""

	for k, v := range userResp.FirstName.Localized {
		if strings.HasPrefix(k, userResp.FirstName.PreferredLocale.Language) {
			firstName = v
			break
		}
	}

	for k, v := range userResp.LastName.Localized {
		if strings.HasPrefix(k, userResp.LastName.PreferredLocale.Language) {
			lastName = v
			break
		}
	}

	// Get profile picture if available
	picture := ""
	if len(userResp.ProfilePicture.DisplayImage.Elements) > 0 {
		if len(userResp.ProfilePicture.DisplayImage.Elements[0].Identifiers) > 0 {
			picture = userResp.ProfilePicture.DisplayImage.Elements[0].Identifiers[0].Identifier
		}
	}

	return &UserInfo{
		ID:           userResp.ID,
		Name:         fmt.Sprintf("%s %s", firstName, lastName),
		GivenName:    firstName,
		FamilyName:   lastName,
		Picture:      picture,
		Locale:       locale,
		ProviderName: p.Name(),
		// Email will be added by enrichWithEmail
		Email:         "",
		VerifiedEmail: true, // LinkedIn emails are verified
	}, nil
}

// enrichWithEmail fetches the user's email address
func (p *LinkedInProvider) enrichWithEmail(ctx context.Context, token *Token, userInfo *UserInfo) error {
	req, err := http.NewRequestWithContext(ctx, "GET", linkedinEmailURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create email request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform email request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read email response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get user email: %s", body)
	}

	var emailResp struct {
		Elements []struct {
			Handle struct {
				EmailAddress string `json:"emailAddress"`
			} `json:"handle~"`
		} `json:"elements"`
	}

	if err := json.Unmarshal(body, &emailResp); err != nil {
		return fmt.Errorf("failed to parse email response: %w", err)
	}

	// Get primary email
	if len(emailResp.Elements) > 0 {
		userInfo.Email = emailResp.Elements[0].Handle.EmailAddress
	}

	return nil
}

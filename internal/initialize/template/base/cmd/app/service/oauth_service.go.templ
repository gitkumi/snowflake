package service

import (
	"{{ .Name }}/internal/oauth"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// OAuthService handles OAuth authentication
type OAuthService struct {
	OAuth *oauth.OAuth
	Redis *redis.Client
}

// NewOAuthService creates a new OAuth service
func NewOAuthService(registry *oauth.OAuth, rdb *redis.Client) *OAuthService {
	return &OAuthService{
		OAuth: registry,
		Redis: rdb,
	}
}

// GetAuthURL returns the authorization URL for a specific provider
func (s *OAuthService) GetAuthURL(providerName string) (string, string, error) {
	provider, exists := s.OAuth.Get(providerName)
	if !exists {
		return "", "", errors.New("provider not found")
	}

	// Generate a random state
	state, err := generateRandomState()
	if err != nil {
		return "", "", err
	}

	// Store the state in Redis with expiration time (e.g., 10 minutes)
	ctx := context.Background()
	err = s.Redis.Set(ctx, fmt.Sprintf("oauth:%s:state:%s", providerName, state), "true", 10*time.Minute).Err()
	if err != nil {
		return "", "", fmt.Errorf("failed to store state: %w", err)
	}

	return provider.GetAuthURL(state), state, nil
}

// Exchange exchanges an authorization code for a token
func (s *OAuthService) Exchange(ctx context.Context, providerName, code, state string) (*oauth.Token, error) {
	// Validate the state using Redis
	stateKey := fmt.Sprintf("oauth:%s:state:%s", providerName, state)
	val, err := s.Redis.Get(ctx, stateKey).Result()
	if err == redis.Nil || val != "true" {
		return nil, errors.New("invalid state parameter")
	} else if err != nil {
		return nil, fmt.Errorf("error validating state: %w", err)
	}

	// Delete the state after validation
	err = s.Redis.Del(ctx, stateKey).Err()
	if err != nil {
		// Log this error but continue with the exchange
		fmt.Printf("Error deleting state from Redis: %v\n", err)
	}

	provider, exists := s.OAuth.Get(providerName)
	if !exists {
		return nil, errors.New("provider not found")
	}

	return provider.Exchange(ctx, code)
}

// GetUserInfo retrieves user information using a token
func (s *OAuthService) GetUserInfo(ctx context.Context, providerName string, token *oauth.Token) (*oauth.UserInfo, error) {
	provider, exists := s.OAuth.Get(providerName)
	if !exists {
		return nil, errors.New("provider not found")
	}

	return provider.GetUserInfo(ctx, token)
}

// ListProviders returns a list of available providers
func (s *OAuthService) ListProviders() []string {
	providers := s.OAuth.List()
	result := make([]string, 0, len(providers))

	for name := range providers {
		result = append(result, name)
	}

	return result
}

// generateRandomState generates a random state string for CSRF protection
func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random state: %w", err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

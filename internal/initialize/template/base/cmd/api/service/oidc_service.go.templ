package service

import (
	"{{ .Name }}/internal/oauth"
	"{{ .Name }}/internal/oidc"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// OIDCService handles OpenID Connect authentication
type OIDCService struct {
	OIDC  *oidc.OIDC
	Redis *redis.Client
}

// NewOIDCService creates a new OIDC service
func NewOIDCService(registry *oidc.OIDC, rdb *redis.Client) *OIDCService {
	return &OIDCService{
		OIDC:  registry,
		Redis: rdb,
	}
}

// GetAuthURL returns the authorization URL for a specific provider
func (s *OIDCService) GetAuthURL(providerName string) (string, string, string, error) {
	provider, exists := s.OIDC.Get(providerName)
	if !exists {
		return "", "", "", errors.New("provider not found")
	}

	// Generate a random state and nonce
	state, err := generateRandomState()
	if err != nil {
		return "", "", "", err
	}

	nonce, err := generateRandomState() // Reuse the same function for nonce
	if err != nil {
		return "", "", "", err
	}

	// Store the state and nonce in Redis with expiration time (e.g., 10 minutes)
	ctx := context.Background()

	// Store state
	stateKey := fmt.Sprintf("oidc:%s:state:%s", providerName, state)
	err = s.Redis.Set(ctx, stateKey, "true", 10*time.Minute).Err()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to store state: %w", err)
	}

	// Store nonce
	nonceKey := fmt.Sprintf("oidc:%s:nonce:%s", providerName, nonce)
	err = s.Redis.Set(ctx, nonceKey, "true", 10*time.Minute).Err()
	if err != nil {
		// Clean up state if nonce storage fails
		s.Redis.Del(ctx, stateKey)
		return "", "", "", fmt.Errorf("failed to store nonce: %w", err)
	}

	return provider.GetAuthURL(state, nonce), state, nonce, nil
}

// Exchange exchanges an authorization code for a token
func (s *OIDCService) Exchange(ctx context.Context, providerName, code, state string) (*oauth.Token, error) {
	// Validate the state using Redis
	stateKey := fmt.Sprintf("oidc:%s:state:%s", providerName, state)
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

	provider, exists := s.OIDC.Get(providerName)
	if !exists {
		return nil, errors.New("provider not found")
	}

	return provider.Exchange(ctx, code)
}

// VerifyIDToken verifies the ID token
func (s *OIDCService) VerifyIDToken(ctx context.Context, providerName, idToken, nonce string) (*oidc.IDTokenClaims, error) {
	provider, exists := s.OIDC.Get(providerName)
	if !exists {
		return nil, errors.New("provider not found")
	}

	// Validate the nonce if provided
	if nonce != "" {
		nonceKey := fmt.Sprintf("oidc:%s:nonce:%s", providerName, nonce)
		val, err := s.Redis.Get(ctx, nonceKey).Result()
		if err == redis.Nil || val != "true" {
			return nil, errors.New("invalid nonce parameter")
		} else if err != nil {
			return nil, fmt.Errorf("error validating nonce: %w", err)
		}

		// Delete the nonce after validation
		err = s.Redis.Del(ctx, nonceKey).Err()
		if err != nil {
			// Log this error but continue with the verification
			fmt.Printf("Error deleting nonce from Redis: %v\n", err)
		}
	}

	return provider.VerifyIDToken(ctx, idToken)
}

// GetUserInfo retrieves user information using a token
func (s *OIDCService) GetUserInfo(ctx context.Context, providerName string, token *oauth.Token) (*oauth.UserInfo, error) {
	provider, exists := s.OIDC.Get(providerName)
	if !exists {
		return nil, errors.New("provider not found")
	}

	return provider.GetUserInfo(ctx, token)
}

// ListProviders returns a list of available OIDC providers
func (s *OIDCService) ListProviders() []string {
	providers := s.OIDC.List()
	result := make([]string, 0, len(providers))

	for name := range providers {
		result = append(result, name)
	}

	return result
}

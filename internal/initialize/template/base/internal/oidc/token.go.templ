package oidc

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidToken    = errors.New("invalid token format")
	ErrTokenExpired    = errors.New("token has expired")
	ErrInvalidIssuer   = errors.New("invalid issuer")
	ErrInvalidAudience = errors.New("invalid audience")
	ErrInvalidNonce    = errors.New("invalid nonce")
)

// ParseIDToken parses an ID token without verification
func ParseIDToken(idToken string) (*IDTokenClaims, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	// Decode the payload (second part)
	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode token payload: %w", err)
	}

	var claims IDTokenClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse token claims: %w", err)
	}

	// Store the raw token
	claims.RawToken = idToken

	return &claims, nil
}

// ValidateBasicClaims validates the basic token claims
func ValidateBasicClaims(claims *IDTokenClaims, expectedIssuer, expectedAudience, expectedNonce string) error {
	// Check expiration
	if time.Now().Unix() > claims.ExpirationTime {
		return ErrTokenExpired
	}

	// Check issuer if provided
	if expectedIssuer != "" && claims.Issuer != expectedIssuer {
		return ErrInvalidIssuer
	}

	// Check audience if provided
	if expectedAudience != "" {
		// Audience can be a string or an array
		if claims.Audience != expectedAudience {
			// TODO: Handle audience as array if needed
			return ErrInvalidAudience
		}
	}

	// Check nonce if provided
	if expectedNonce != "" && claims.Nonce != expectedNonce {
		return ErrInvalidNonce
	}

	return nil
}

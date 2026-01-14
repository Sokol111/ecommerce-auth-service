package command

import (
	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
)

// PasswordHasher handles password hashing and comparison
type PasswordHasher interface {
	// Hash generates a hash from the given password
	Hash(password string) (string, error)
	// Compare compares a hash with a password, returns true if they match
	Compare(hash, password string) bool
}

// TokenClaims contains the claims extracted from a token
type TokenClaims struct {
	UserID string
	Email  string
	Role   adminuser.Role
}

// TokenService handles JWT/PASETO token operations
type TokenService interface {
	// GenerateTokenPair generates an access token and refresh token pair
	// Returns: accessToken, refreshToken, expiresInSeconds, error
	GenerateTokenPair(user *adminuser.AdminUser) (string, string, int, error)

	// ValidateAccessToken validates an access token and returns its claims
	ValidateAccessToken(token string) (*TokenClaims, error)

	// ValidateRefreshToken validates a refresh token and returns its claims
	ValidateRefreshToken(token string) (*TokenClaims, error)
}

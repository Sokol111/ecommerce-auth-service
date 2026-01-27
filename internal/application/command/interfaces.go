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

// TokenGenerator handles token generation operations
type TokenGenerator interface {
	// GenerateTokenPair generates an access token and refresh token pair
	// Returns: accessToken, refreshToken, refreshTokenID, expiresInSeconds, refreshExpiresInSeconds, error
	GenerateTokenPair(user *adminuser.AdminUser) (string, string, string, int, int, error)
}

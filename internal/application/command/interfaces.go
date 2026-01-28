package command

import (
	"time"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
)

// PasswordHasher handles password hashing and comparison
type PasswordHasher interface {
	// Hash generates a hash from the given password
	Hash(password string) (string, error)
	// Compare compares a hash with a password, returns true if they match
	Compare(hash, password string) bool
}

// TokenPairResult contains the result of token pair generation
type TokenPairResult struct {
	AccessToken      string
	RefreshToken     string
	RefreshTokenID   string
	ExpiresIn        int
	ExpiresAt        time.Time
	RefreshExpiresIn int
	RefreshExpiresAt time.Time
}

// TokenGenerator handles token generation operations
type TokenGenerator interface {
	// GenerateTokenPair generates an access token and refresh token pair
	GenerateTokenPair(user *adminuser.AdminUser) (*TokenPairResult, error)
}

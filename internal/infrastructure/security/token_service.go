package security

import (
	"encoding/hex"
	"errors"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/Sokol111/ecommerce-auth-service/internal/application/command"
	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
)

// TokenConfig holds configuration for token generation
type TokenConfig struct {
	// PrivateKey is the hex-encoded Ed25519 private key (64 bytes = 128 hex chars)
	PrivateKey           string        `yaml:"privateKey"`
	AccessTokenDuration  time.Duration `yaml:"accessTokenDuration"`
	RefreshTokenDuration time.Duration `yaml:"refreshTokenDuration"`
}

// DefaultTokenConfig returns default token configuration
func DefaultTokenConfig() TokenConfig {
	return TokenConfig{
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour, // 7 days
	}
}

// pasetoService implements TokenGenerator using PASETO v4 Public (asymmetric)
type pasetoService struct {
	privateKey paseto.V4AsymmetricSecretKey
	config     TokenConfig
}

// newPasetoService creates a new PASETO token service with asymmetric keys
func newPasetoService(config TokenConfig) (command.TokenGenerator, error) {
	var privateKey paseto.V4AsymmetricSecretKey

	if config.PrivateKey != "" {
		// Use provided private key
		keyBytes, err := hex.DecodeString(config.PrivateKey)
		if err != nil {
			return nil, errors.New("invalid private key hex encoding")
		}
		if len(keyBytes) != 64 {
			return nil, errors.New("private key must be 64 bytes (128 hex characters)")
		}

		privateKey, err = paseto.NewV4AsymmetricSecretKeyFromBytes(keyBytes)
		if err != nil {
			return nil, err
		}
	} else {
		// Generate a new key pair (for development only!)
		privateKey = paseto.NewV4AsymmetricSecretKey()
	}

	if config.AccessTokenDuration == 0 {
		config.AccessTokenDuration = DefaultTokenConfig().AccessTokenDuration
	}
	if config.RefreshTokenDuration == 0 {
		config.RefreshTokenDuration = DefaultTokenConfig().RefreshTokenDuration
	}

	return &pasetoService{
		privateKey: privateKey,
		config:     config,
	}, nil
}

// GetPublicKeyHex returns the public key as hex string.
// This can be shared with other services for token validation.
func (s *pasetoService) GetPublicKeyHex() string {
	return hex.EncodeToString(s.privateKey.Public().ExportBytes())
}

// GenerateTokenPair generates access and refresh tokens for a user
func (s *pasetoService) GenerateTokenPair(user *adminuser.AdminUser) (string, string, int, error) {
	accessToken, err := s.generateToken(user, s.config.AccessTokenDuration, "access")
	if err != nil {
		return "", "", 0, err
	}

	refreshToken, err := s.generateToken(user, s.config.RefreshTokenDuration, "refresh")
	if err != nil {
		return "", "", 0, err
	}

	expiresIn := int(s.config.AccessTokenDuration.Seconds())
	return accessToken, refreshToken, expiresIn, nil
}

func (s *pasetoService) generateToken(user *adminuser.AdminUser, duration time.Duration, tokenType string) (string, error) {
	now := time.Now()

	// Get permissions for the user's role
	permissions := user.GetPermissions()
	permStrings := make([]string, len(permissions))
	for i, p := range permissions {
		permStrings[i] = string(p)
	}

	tk := paseto.NewToken()
	tk.SetIssuedAt(now)
	tk.SetNotBefore(now)
	tk.SetExpiration(now.Add(duration))
	tk.SetSubject(user.ID)
	tk.SetString("role", string(user.Role))
	tk.SetString("type", tokenType)
	tk.Set("permissions", permStrings)

	// Sign with private key (V4 Public = asymmetric)
	return tk.V4Sign(s.privateKey, nil), nil
}

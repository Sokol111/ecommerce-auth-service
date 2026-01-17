package security

import (
	"encoding/hex"
	"errors"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/Sokol111/ecommerce-auth-service/internal/application/command"
	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
)

// tokenGenerator implements TokenGenerator using PASETO v4 Public (asymmetric)
type tokenGenerator struct {
	privateKey         paseto.V4AsymmetricSecretKey
	config             Config
	permissionProvider adminuser.RolePermissionProvider
}

// newTokenGenerator creates a new PASETO token generator with asymmetric keys
func newTokenGenerator(config Config, permissionProvider adminuser.RolePermissionProvider) (command.TokenGenerator, error) {
	keyBytes, err := hex.DecodeString(config.PrivateKey)
	if err != nil {
		return nil, errors.New("invalid private key hex encoding")
	}
	if len(keyBytes) != 64 {
		return nil, errors.New("private key must be 64 bytes (128 hex characters)")
	}

	privateKey, err := paseto.NewV4AsymmetricSecretKeyFromBytes(keyBytes)
	if err != nil {
		return nil, err
	}

	return &tokenGenerator{
		privateKey:         privateKey,
		config:             config,
		permissionProvider: permissionProvider,
	}, nil
}

// GetPublicKeyHex returns the public key as hex string.
// This can be shared with other services for token validation.
func (s *tokenGenerator) GetPublicKeyHex() string {
	return hex.EncodeToString(s.privateKey.Public().ExportBytes())
}

// GenerateTokenPair generates access and refresh tokens for a user
func (s *tokenGenerator) GenerateTokenPair(user *adminuser.AdminUser) (string, string, int, int, error) {
	accessToken, err := s.generateToken(user, s.config.AccessTokenDuration, "access")
	if err != nil {
		return "", "", 0, 0, err
	}

	refreshToken, err := s.generateToken(user, s.config.RefreshTokenDuration, "refresh")
	if err != nil {
		return "", "", 0, 0, err
	}

	expiresIn := int(s.config.AccessTokenDuration.Seconds())
	refreshExpiresIn := int(s.config.RefreshTokenDuration.Seconds())
	return accessToken, refreshToken, expiresIn, refreshExpiresIn, nil
}

func (s *tokenGenerator) generateToken(user *adminuser.AdminUser, duration time.Duration, tokenType string) (string, error) {
	now := time.Now()

	// Get permissions for the user's role
	permissions := s.permissionProvider.GetPermissionsForRole(user.Role)
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

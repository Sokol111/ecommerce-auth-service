package security

import (
	"encoding/hex"
	"errors"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/Sokol111/ecommerce-auth-service/internal/application/command"
	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
	"github.com/google/uuid"
	"github.com/samber/lo"
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
func (s *tokenGenerator) GenerateTokenPair(user *adminuser.AdminUser) (*command.TokenPairResult, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(s.config.AccessTokenDuration)
	refreshExpiresAt := now.Add(s.config.RefreshTokenDuration)

	accessToken, err := s.generateToken(user, s.config.AccessTokenDuration, "access", "")
	if err != nil {
		return nil, err
	}

	// Generate unique ID for refresh token to enable rotation/revocation
	refreshTokenID := uuid.New().String()
	refreshToken, err := s.generateToken(user, s.config.RefreshTokenDuration, "refresh", refreshTokenID)
	if err != nil {
		return nil, err
	}

	return &command.TokenPairResult{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		RefreshTokenID:   refreshTokenID,
		ExpiresIn:        int(s.config.AccessTokenDuration.Seconds()),
		ExpiresAt:        expiresAt,
		RefreshExpiresIn: int(s.config.RefreshTokenDuration.Seconds()),
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}

func (s *tokenGenerator) generateToken(user *adminuser.AdminUser, duration time.Duration, tokenType string, tokenID string) (string, error) {
	now := time.Now().UTC()

	permissions := s.permissionProvider.GetPermissionsForRole(user.Role)
	permStrings := lo.Map(permissions, func(p adminuser.Permission, _ int) string {
		return string(p)
	})

	tk := paseto.NewToken()
	tk.SetIssuedAt(now)
	tk.SetExpiration(now.Add(duration))
	tk.SetSubject(user.ID)
	tk.SetString("role", string(user.Role))
	tk.SetString("type", tokenType)
	tk.Set("permissions", permStrings)

	if tokenID != "" {
		tk.SetJti(tokenID)
	}

	return tk.V4Sign(s.privateKey, nil), nil
}

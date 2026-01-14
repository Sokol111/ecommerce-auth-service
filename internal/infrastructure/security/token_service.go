package security
package security

import (
	"errors"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/Sokol111/ecommerce-auth-service/internal/application/command"
	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// TokenConfig holds configuration for token generation
type TokenConfig struct {
	SecretKey             string        `yaml:"secretKey"`
	AccessTokenDuration   time.Duration `yaml:"accessTokenDuration"`
	RefreshTokenDuration  time.Duration `yaml:"refreshTokenDuration"`
}

// DefaultTokenConfig returns default token configuration
func DefaultTokenConfig() TokenConfig {
	return TokenConfig{
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour, // 7 days
	}
}

// PasetoService implements TokenService using PASETO v4
type PasetoService struct {
	symmetricKey paseto.V4SymmetricKey
	config       TokenConfig
}

// NewPasetoService creates a new PASETO token service
func NewPasetoService(config TokenConfig) (*PasetoService, error) {
	var key paseto.V4SymmetricKey
	
	if config.SecretKey != "" {
		// Use provided key (must be 32 bytes hex encoded)
		keyBytes, err := hexDecode(config.SecretKey)
		if err != nil {
			return nil, err
		}
		key = paseto.NewV4SymmetricKey()
		copy(key.ExportBytes()[:], keyBytes)
	} else {
		// Generate a new random key
		key = paseto.NewV4SymmetricKey()
	}

	if config.AccessTokenDuration == 0 {
		config.AccessTokenDuration = DefaultTokenConfig().AccessTokenDuration
	}
	if config.RefreshTokenDuration == 0 {
		config.RefreshTokenDuration = DefaultTokenConfig().RefreshTokenDuration
	}

	return &PasetoService{
		symmetricKey: key,
		config:       config,
	}, nil
}

// GenerateTokenPair generates access and refresh tokens for a user
func (s *PasetoService) GenerateTokenPair(user *adminuser.AdminUser) (string, string, int, error) {
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

// ValidateAccessToken validates an access token and returns claims
func (s *PasetoService) ValidateAccessToken(token string) (*command.TokenClaims, error) {
	return s.validateToken(token, "access")
}

// ValidateRefreshToken validates a refresh token and returns claims
func (s *PasetoService) ValidateRefreshToken(token string) (*command.TokenClaims, error) {
	return s.validateToken(token, "refresh")
}

func (s *PasetoService) generateToken(user *adminuser.AdminUser, duration time.Duration, tokenType string) (string, error) {
	now := time.Now()

	token := paseto.NewToken()
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(now.Add(duration))
	token.SetSubject(user.ID)
	token.SetString("email", user.Email)
	token.SetString("role", string(user.Role))
	token.SetString("type", tokenType)

	return token.V4Encrypt(s.symmetricKey, nil), nil
}

func (s *PasetoService) validateToken(tokenString string, expectedType string) (*command.TokenClaims, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())

	token, err := parser.ParseV4Local(s.symmetricKey, tokenString, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	tokenType, err := token.GetString("type")
	if err != nil || tokenType != expectedType {
		return nil, ErrInvalidToken
	}

	subject, err := token.GetSubject()
	if err != nil {
		return nil, ErrInvalidToken
	}

	email, err := token.GetString("email")
	if err != nil {
		return nil, ErrInvalidToken
	}

	roleStr, err := token.GetString("role")
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &command.TokenClaims{
		UserID: subject,
		Email:  email,
		Role:   adminuser.Role(roleStr),
	}, nil
}

func hexDecode(s string) ([]byte, error) {
	if len(s) != 64 {
		return nil, errors.New("secret key must be 64 hex characters (32 bytes)")
	}
	
	b := make([]byte, 32)
	for i := 0; i < 32; i++ {
		var val byte
		for j := 0; j < 2; j++ {
			c := s[i*2+j]
			switch {
			case c >= '0' && c <= '9':
				val = val*16 + (c - '0')
			case c >= 'a' && c <= 'f':
				val = val*16 + (c - 'a' + 10)
			case c >= 'A' && c <= 'F':
				val = val*16 + (c - 'A' + 10)
			default:
				return nil, errors.New("invalid hex character")
			}
		}
		b[i] = val
	}
	return b, nil
}

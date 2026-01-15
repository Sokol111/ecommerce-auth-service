package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds configuration for token generation
type Config struct {
	// PrivateKey is the hex-encoded Ed25519 private key (64 bytes = 128 hex chars)
	PrivateKey           string        `mapstructure:"private-key"`
	AccessTokenDuration  time.Duration `mapstructure:"access-token-duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh-token-duration"`
}

func newConfig(v *viper.Viper) (Config, error) {
	var cfg Config
	if err := v.Sub("token").UnmarshalExact(&cfg); err != nil {
		return cfg, fmt.Errorf("failed to load token config: %w", err)
	}

	if cfg.PrivateKey == "" {
		return cfg, errors.New("token private key is required")
	}
	if cfg.AccessTokenDuration <= 0 {
		return cfg, errors.New("token access token duration must be positive")
	}
	if cfg.RefreshTokenDuration <= 0 {
		return cfg, errors.New("token refresh token duration must be positive")
	}

	return cfg, nil
}

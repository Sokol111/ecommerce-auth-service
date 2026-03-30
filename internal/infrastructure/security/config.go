package security

import (
	"errors"
	"fmt"
	"time"

	"github.com/knadh/koanf/v2"
)

// Config holds configuration for token generation
type Config struct {
	// PrivateKey is the hex-encoded Ed25519 private key (64 bytes = 128 hex chars)
	PrivateKey           string        `koanf:"private-key"`
	AccessTokenDuration  time.Duration `koanf:"access-token-duration"`
	RefreshTokenDuration time.Duration `koanf:"refresh-token-duration"`
}

func newConfig(k *koanf.Koanf) (Config, error) {
	var cfg Config
	if err := k.Unmarshal("token", &cfg); err != nil {
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

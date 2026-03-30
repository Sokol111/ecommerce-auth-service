package http //nolint:revive // package name intentional

import (
	"fmt"
	"time"

	"github.com/knadh/koanf/v2"
)

// RateLimitConfig holds configuration for rate limiting
type RateLimitConfig struct {
	// Login rate limiting
	LoginTokens   uint64        `koanf:"login-tokens"`   // Number of allowed requests
	LoginInterval time.Duration `koanf:"login-interval"` // Time window for the tokens
}

func newRateLimitConfig(k *koanf.Koanf) (RateLimitConfig, error) {
	cfg := RateLimitConfig{
		// Defaults: 5 requests per minute
		LoginTokens:   5,
		LoginInterval: time.Minute,
	}

	if k.Exists("rate-limit") {
		if err := k.Unmarshal("rate-limit", &cfg); err != nil {
			return cfg, fmt.Errorf("failed to load rate-limit config: %w", err)
		}
	}

	if cfg.LoginTokens == 0 {
		cfg.LoginTokens = 5
	}
	if cfg.LoginInterval <= 0 {
		cfg.LoginInterval = time.Minute
	}

	return cfg, nil
}

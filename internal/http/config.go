package http //nolint:revive // package name intentional

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// RateLimitConfig holds configuration for rate limiting
type RateLimitConfig struct {
	// Login rate limiting
	LoginTokens   uint64        `mapstructure:"login-tokens"`   // Number of allowed requests
	LoginInterval time.Duration `mapstructure:"login-interval"` // Time window for the tokens
}

func newRateLimitConfig(v *viper.Viper) (RateLimitConfig, error) {
	cfg := RateLimitConfig{
		// Defaults: 5 requests per minute
		LoginTokens:   5,
		LoginInterval: time.Minute,
	}

	if sub := v.Sub("rate-limit"); sub != nil {
		if err := sub.UnmarshalExact(&cfg); err != nil {
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

package permissions

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds configuration for roles and their permissions
type Config struct {
	Roles []RoleConfig `mapstructure:"roles"`
}

// RoleConfig holds configuration for a single role
type RoleConfig struct {
	Name        string   `mapstructure:"name"`
	Description string   `mapstructure:"description"`
	Permissions []string `mapstructure:"permissions"`
}

func newConfig(v *viper.Viper) (Config, error) {
	var cfg Config

	sub := v.Sub("permissions")
	if sub == nil {
		return Config{}, fmt.Errorf("permissions configuration is required")
	}

	if err := sub.UnmarshalExact(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to load permissions config: %w", err)
	}

	return cfg, nil
}

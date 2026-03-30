package permissions

import (
	"fmt"

	"github.com/knadh/koanf/v2"
)

// Config holds configuration for roles and their permissions
type Config struct {
	Roles []RoleConfig `koanf:"roles"`
}

// RoleConfig holds configuration for a single role
type RoleConfig struct {
	Name        string   `koanf:"name"`
	Description string   `koanf:"description"`
	Permissions []string `koanf:"permissions"`
}

func newConfig(k *koanf.Koanf) (Config, error) {
	var cfg Config

	if !k.Exists("permissions") {
		return Config{}, fmt.Errorf("permissions configuration is required")
	}

	if err := k.Unmarshal("permissions", &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to load permissions config: %w", err)
	}

	return cfg, nil
}

package seeder

import (
	"github.com/knadh/koanf/v2"
)

// Config holds configuration for initial admin user seeding
type Config struct {
	Enabled   bool   `koanf:"enabled"`
	Email     string `koanf:"email"`
	Password  string `koanf:"password"`
	FirstName string `koanf:"first-name"`
	LastName  string `koanf:"last-name"`
}

func newConfig(k *koanf.Koanf) Config {
	var cfg Config

	if !k.Exists("admin.initial-user") {
		return Config{Enabled: false}
	}

	if err := k.Unmarshal("admin.initial-user", &cfg); err != nil {
		return Config{Enabled: false}
	}

	return cfg
}

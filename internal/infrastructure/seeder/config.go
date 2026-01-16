package seeder

import (
	"github.com/spf13/viper"
)

// Config holds configuration for initial admin user seeding
type Config struct {
	Enabled   bool   `mapstructure:"enabled"`
	Email     string `mapstructure:"email"`
	Password  string `mapstructure:"password"`
	FirstName string `mapstructure:"first-name"`
	LastName  string `mapstructure:"last-name"`
}

func newConfig(v *viper.Viper) Config {
	var cfg Config

	sub := v.Sub("admin.initial-user")
	if sub == nil {
		return Config{Enabled: false}
	}

	if err := sub.UnmarshalExact(&cfg); err != nil {
		return Config{Enabled: false}
	}

	return cfg
}

package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Port      int    `mapstructure:"port"`
	SecretKey string `mapstructure:"secretKey"`
	DB        DB     `mapstructure:"db"`
}

type DB struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"name"`
}

func LoadConfig(path string) Config {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}
	var c Config
	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal config")
	}
	return c
}

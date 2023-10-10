package configs

import (
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config represents the configuration settings for your application.
type Config struct {
	// Logger
	ZerologLevel zerolog.Level `env:"ZEROLOG_LEVEL" envDefault:"0"`

	// Transloadit
	TransloaditAuthKey    string `env:"TRANSLOADIT_AUTH_KEY" envDefault:"KEY"`
	TransloaditAuthSecret string `env:"TRANSLOADIT_AUTH_SECRET" envDefault:"SEC"`
}

// LoadConfigs loads the application configuration from environment variables.
func LoadConfigs() (Config, error) {
	cfg := Config{}

	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Error().Err(err).Msg("Error loading .env file")
		return Config{}, err
	}

	// Parse environment variables
	if err := env.Parse(&cfg); err != nil {
		log.Error().Err(err).Msg("Error parsing environment variables")
		return Config{}, err
	}

	// logger.Debug().Msg("Loaded config successfully")
	return cfg, nil
}

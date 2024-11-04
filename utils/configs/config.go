package configs

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Config holds the application wide configurations.
// The values are read by viper from the config file or environment variables.
type Config struct {
	PSQL PSQLConfig `mapstructure:"postgres"`
}

// PSQLConfig holds PostgreSQL connection configuration.
type PSQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

var c *Config

// C returns the loaded configuration globally.
func C() *Config {
	if c == nil {
		log.Fatal().Msg("Configuration not initialized. Call RunConfig() first.")
	}
	return c
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	// Set default values for the configuration.
	setDefaults(v)

	// Read from environment variables
	v.AutomaticEnv()

	// Try to read from config file, but continue if not found
	v.AddConfigPath(path)
	v.SetConfigName("config.example")
	v.SetConfigType("json")

	// Try to read config file, but log the error instead of failing
	if err := v.ReadInConfig(); err != nil {
		log.Warn().Msgf("No config file found; using environment variables: %v", err)
	}

	// Unmarshal the configuration into the config struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %w", err)
	}

	// Validate essential configuration values
	if err := validatePSQLConfig(config.PSQL); err != nil {
		return nil, err
	}

	return &config, nil
}

// RunConfig initializes and loads the configuration.
func RunConfig(path string) {
	config, err := LoadConfig(path)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}
	c = config
}

// setDefaults sets default configuration values in viper.
func setDefaults(v *viper.Viper) {
	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", "5432")
	v.SetDefault("postgres.user", "root")
	v.SetDefault("postgres.password", "secret")
	v.SetDefault("postgres.database", "psql_db")
	v.SetDefault("postgres.ssl_mode", "disable")
}

// validatePSQLConfig ensures that essential PSQL config values are present.
func validatePSQLConfig(psqlConfig PSQLConfig) error {
	if psqlConfig.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if psqlConfig.Port == "" {
		return fmt.Errorf("database port is required")
	}
	if psqlConfig.User == "" {
		return fmt.Errorf("database user is required")
	}
	if psqlConfig.Password == "" {
		return fmt.Errorf("database password is required")
	}
	if psqlConfig.Database == "" {
		return fmt.Errorf("database name is required")
	}
	return nil
}

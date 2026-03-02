// Package config provides configuration management for the GENAPP application.
// It uses Viper to load configuration from files, environment variables, and defaults.
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	// Database configuration
	Database DatabaseConfig `mapstructure:"database"`

	// Server configuration (for optional REST API)
	Server ServerConfig `mapstructure:"server"`

	// Application settings
	LogLevel string `mapstructure:"log_level"`
}

// DatabaseConfig holds database connection settings.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// DSN returns the PostgreSQL connection string.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
	)
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// Address returns the server address in host:port format.
func (s ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// Load reads configuration from files and environment variables.
func Load() (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Configuration file settings
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("$HOME/.genapp")

	// Environment variable settings
	v.SetEnvPrefix("GENAPP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read configuration file (optional)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found is not an error - we use defaults/env vars
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}

// setDefaults configures default values for all settings.
func setDefaults(v *viper.Viper) {
	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "genapp")
	v.SetDefault("database.password", "genapp")
	v.SetDefault("database.dbname", "genapp")
	v.SetDefault("database.sslmode", "disable")

	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)

	// Application defaults
	v.SetDefault("log_level", "info")
}

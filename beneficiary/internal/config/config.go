package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config holds all configuration for the service
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port         string        `envconfig:"PORT" default:"8080"`
	ReadTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"5s"`
	WriteTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"10s"`
}

type DatabaseConfig struct {
	Host         string `envconfig:"DB_HOST" required:"true"`
	Port         int    `envconfig:"DB_PORT" default:"5432"`
	User         string `envconfig:"DB_USER" required:"true"`
	Password     string `envconfig:"DB_PASSWORD" required:"true"`
	Database     string `envconfig:"DB_NAME" required:"true"`
	SSLMode      string `envconfig:"DB_SSLMODE" default:"disable"`
	MaxOpenConns int    `envconfig:"DB_MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns int    `envconfig:"DB_MAX_IDLE_CONNS" default:"25"`
}

type LoggingConfig struct {
	Level       string `envconfig:"LOG_LEVEL" default:"info"`
	PrettyPrint bool   `envconfig:"LOG_PRETTY_PRINT" default:"false"`
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}
	return &config, nil
}

// GetDSN returns the PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode,
	)
}

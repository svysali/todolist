package db

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

const (
	sslEnabled  = "required"
	sslDisabled = "disable"
)

// Config holds the database configuration from the environment
type Config struct {
	User               string `envconfig:"db_user"`
	Host               string `envconfig:"db_host"`
	Port               int    `envconfig:"db_port"`
	Database           string `envconfig:"db_name"`
	Password           string `envconfig:"db_pass"`
	SSLEnabled         bool   `envconfig:"db_ssl_enabled"`
	MaxConnections     int    `envconfig:"db_max_connections"`
	MaxIdleConnections int    `envconfig:"db_max_idle_connections"`
}

// SSLMode Returns the appropriate string constant
func (c *Config) SSLMode() string {
	var mode string
	if c.SSLEnabled {
		mode = sslEnabled
	} else {
		mode = sslDisabled
	}
	if log.GetLevel() >= log.InfoLevel {
		log.WithField("mode", mode).Debug("SSLMode set")
	}
	return mode
}

// ConnectionString gets the connection string from the environment variables
func (c *Config) ConnectionString() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		c.User,
		c.Password,
		c.Database,
		c.Host,
		c.Port,
		c.SSLMode(),
	)
}

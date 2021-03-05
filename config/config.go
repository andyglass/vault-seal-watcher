package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

var (
	// Cfg app global
	Cfg = Get()
)

// Config structure
type Config struct {
	ServerHost       string        `envconfig:"SERVER_HOST" default:"0.0.0.0"`
	ServerPort       string        `envconfig:"SERVER_PORT" default:"8000"`
	LogLevel         string        `envconfig:"SERVER_LOG_LEVEL" default:"info"`
	ReadTimeout      time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"10s"`
	WriteTimeout     time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"10s"`
	VaultAddr        string        `envconfig:"VAULT_ADDR" default:"http://localhost:8200"`
	VaultUnsealKeys  []string      `envconfig:"VAULT_UNSEAL_KEYS" default:""`
	VaultUnsealDelay time.Duration `envconfig:"VAULT_UNSEAL_DELAY" default:"200ms"`
	VaultTimeout     time.Duration `envconfig:"VAULT_TIMEOUT" default:"10s"`
	VaultWatchPeriod time.Duration `envconfig:"VAULT_WATCH_PERIOD" default:"60s"`
}

// Get config
func Get() *Config {
	cfg := &Config{}

	if err := envconfig.Process("", cfg); err != nil {
		panic(err)
	}

	return cfg
}

// Package config loads service configuration via viper. Base config.yaml
// plus whichever env-specific file CONFIG_FILE points at — see
// document/ARCHITECTURE.md "Configuration".
package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server Server `mapstructure:"server"`
	Logger Logger `mapstructure:"logger"`
	// Add CacheDatabase / Database / Kafka / Storage / Firebase / Downstreams
	// structs here as you uncomment the matching blocks in config.example.yaml.
}

type Server struct {
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	ServiceID    string        `mapstructure:"service_id"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type Logger struct {
	Level    string `mapstructure:"level"`
	Encoding string `mapstructure:"encoding"`
}

// Load reads the file at path (defaults to config.yaml), binds env vars for
// anything using "${VAR}" interpolation in the YAML, and unmarshals into Config.
func Load(path string) (*Config, error) {
	v := viper.New()
	if path == "" {
		path = "config/config.yaml"
	}
	v.SetConfigFile(path)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Package config contains the configuration parsing options.
package config

import (
	"io"
	"log/slog"
	"os"

	"github.com/BurntSushi/toml"
)

// Must panics if the error is not nil.
func Must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}

// OpensearchConfig is the configuration struct for the Opensearch server to monitor.
type OpensearchConfig struct {
	HealthURL                string `toml:"health_url"`
	SecondsInGreenForHealthy int    `toml:"seconds_in_green_for_healthy"`
	TickerInterval           int    `toml:"ticker_interval"`
}

// ServerConfig is the configuration struct for the server to expose the health metric.
type ServerConfig struct {
	Address    string `toml:"address"`
	Port       int    `toml:"port"`
	LogLevel   string `toml:"log_level"`
	LogHandler string `toml:"log_handler"`
}

// Config is the configuration struct.
type Config struct {
	Opensearch OpensearchConfig `toml:"opensearch"`
	Server     ServerConfig     `toml:"server"`
}

func defaultConfig() *Config {
	return &Config{
		Opensearch: OpensearchConfig{
			HealthURL:                "http://localhost:9200/_cluster/health",
			SecondsInGreenForHealthy: 60,
			TickerInterval:           1,
		},
		Server: ServerConfig{
			Address:    "localhost",
			Port:       8080,
			LogLevel:   slog.LevelInfo.String(),
			LogHandler: "json",
		},
	}
}

// InitConfig loads and initializes the configuration struct.
// This function will panic if the configuration is not valid.
func InitConfig(p string) *Config {
	c := defaultConfig()

	fd := Must(os.Open(p))
	defer fd.Close()

	data := Must(io.ReadAll(fd))

	err := toml.Unmarshal(data, &c)
	if err != nil {
		panic(err)
	}

	return c
}

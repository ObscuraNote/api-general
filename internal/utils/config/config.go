package config

import (
	"bufio"
	"os"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/philippe-berto/database/postgresdb"
	"github.com/philippe-berto/tracer"
)

type BaseConfig struct {
	DB               postgresdb.Config
	Debug            bool `env:"DEBUG" envDefault:"true"`
	Port             int  `env:"PORT" envDefault:"8080"`
	Metrics          MetricsConfig
	Tracer           tracer.Config
	Service          string `env:"APP_SERVICE" envDefault:"cryple_general"`
	Name             string `env:"APP_NAME" envDefault:"cryple"`
	EnableCORS       bool   `env:"ENABLE_CORS" envDefault:"true"`
	CorsAllowOrigins string `env:"CORS_ALLOW_ORIGINS" envDefault:"*"`
}

type MetricsConfig struct {
	Port   int64 `env:"METRICS_PORT"   envDefault:"80"`
	Enable bool  `env:"METRICS_ENABLE" envDefault:"0"`
}

func loadEnvFile() {
	file, err := os.Open(".env")
	if err != nil {
		return // .env file doesn't exist, ignore
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'')) {
				value = value[1 : len(value)-1]
			}
			os.Setenv(key, value)
		}
	}
}

func Load() (BaseConfig, error) {
	// Try to load .env file (ignore error if file doesn't exist)
	loadEnvFile()

	cfg := BaseConfig{}
	if err := env.Parse(&cfg); err != nil {
		return BaseConfig{}, err
	}

	return cfg, nil
}

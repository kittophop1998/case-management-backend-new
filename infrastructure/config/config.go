package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Services ServicesConfig `yaml:"services"`
}

type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type ServerConfig struct {
	Port     int    `yaml:"port"`
	GinMode  string `yaml:"gin_mode"`
	LogLevel string `yaml:"log_level"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type ServicesConfig struct {
	ConnectorAPI ConnectorAPIConfig `yaml:"connector_api"`
}

type ConnectorAPIConfig struct {
	BaseURL string `yaml:"base_url"`
}

// allowed environments for config loading
var allowedEnvs = map[string]struct{}{
	"dev":     {},
	"sit":     {},
	"uat":     {},
	"prod":    {},
	"railway": {},
}

var (
	cfg  *Config
	once sync.Once
)

func ValidateEnv(env string) error {
	if _, ok := allowedEnvs[env]; !ok {
		return fmt.Errorf("invalid environment: %s (allowed: %v)", env, keys(allowedEnvs))
	}
	return nil
}

func keys(m map[string]struct{}) []string {
	result := make([]string, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

// Load reads the config file for the given environment
func Load(env string) (*Config, error) {
	if err := ValidateEnv(env); err != nil {
		return nil, err
	}

	var loadErr error
	once.Do(func() {
		log.Printf("Loading config for environment: %s", env)

		filename := fmt.Sprintf("configs/%s.yaml", env)
		data, err := os.ReadFile(filename)
		if err != nil {
			loadErr = fmt.Errorf("failed to read config file: %w", err)
			return
		}

		var config Config
		if err := yaml.Unmarshal(data, &config); err != nil {
			loadErr = fmt.Errorf("failed to parse config: %w", err)
			return
		}

		cfg = &config
	})

	return cfg, loadErr
}

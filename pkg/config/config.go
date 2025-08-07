package config

import (
	"os"

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
	AgreementAPI AgreementAPIConfig `yaml:"agreement_api"`
}

type AgreementAPIConfig struct {
	BaseURL string `yaml:"base_url"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

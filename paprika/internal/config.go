package internal

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds runtime configuration for the CLI.
type Config struct {
	Token   string `yaml:"token"`
	BaseURL string `yaml:"base_url"`
}

// Load reads configuration with the following precedence (highest to lowest):
//  1. Environment variable PAPRIKA_TOKEN
//  2. Config file at ~/.config/<cliName>/config.yaml
func Load(cliName string) (*Config, error) {
	cfg := &Config{}

	// Attempt to load from the config file first (lowest precedence).
	configDir, err := os.UserConfigDir()
	if err == nil {
		configPath := filepath.Join(configDir, cliName, "config.yaml")
		data, readErr := os.ReadFile(configPath)
		if readErr == nil {
			_ = yaml.Unmarshal(data, cfg)
		}
	}

	// Environment variable overrides the config file.
	envKey := strings.ToUpper(strings.NewReplacer("-", "_", ".", "_").Replace(cliName)) + "_TOKEN"
	if token := os.Getenv(envKey); token != "" {
		cfg.Token = token
	}

	return cfg, nil
}

// Auth environment variables detected from the OpenAPI spec:
//   PAPRIKA_TOKEN

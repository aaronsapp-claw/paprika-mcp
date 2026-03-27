package internal

import (
	"fmt"
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

// commandspec:custom:start
// Save writes cfg to ~/.config/<cliName>/config.yaml, creating the directory
// if necessary. The file is written with 0600 permissions so the token is not
// world-readable.
func Save(cliName string, cfg *Config) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("locating config directory: %w", err)
	}

	dir := filepath.Join(configDir, cliName)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("creating config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("encoding config: %w", err)
	}

	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return "", fmt.Errorf("writing config file: %w", err)
	}

	return configPath, nil
}

// commandspec:custom:end

// Auth environment variables detected from the OpenAPI spec:
//
//	PAPRIKA_TOKEN

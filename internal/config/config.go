/*
Copyright © 2025 vijay papanaboina

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the CLI configuration
type Config struct {
	APIURL       string `mapstructure:"api_url" yaml:"api_url"`
	AccessToken  string `mapstructure:"access_token" yaml:"access_token"`
	RefreshToken string `mapstructure:"refresh_token" yaml:"refresh_token"`
	APIKey       string `mapstructure:"api_key" yaml:"api_key"`
}

var (
	viperInstance *viper.Viper
	configPath    string
)

const (
	configDirName  = ".cloud-storage-cli"
	configFileName  = "config.yaml"
	defaultAPIURL   = "http://localhost:8000"
	envVarPrefix    = "CLOUD_STORAGE"
)

// InitConfig initializes Viper with defaults and environment variable support
func InitConfig() error {
	viperInstance = viper.New()

	// Set defaults
	viperInstance.SetDefault("api_url", defaultAPIURL)
	viperInstance.SetDefault("access_token", "")
	viperInstance.SetDefault("refresh_token", "")
	viperInstance.SetDefault("api_key", "")

	// Set config file name and type
	viperInstance.SetConfigName(configFileName)
	viperInstance.SetConfigType("yaml")

	// Set config file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, configDirName)
	configPath = filepath.Join(configDir, configFileName)

	// Add config path
	viperInstance.AddConfigPath(configDir)

	// Set environment variable prefix
	viperInstance.SetEnvPrefix(envVarPrefix)
	viperInstance.AutomaticEnv()

	// Bind environment variables
	viperInstance.BindEnv("api_url", "CLOUD_STORAGE_API_URL")
	viperInstance.BindEnv("access_token", "CLOUD_STORAGE_ACCESS_TOKEN")
	viperInstance.BindEnv("refresh_token", "CLOUD_STORAGE_REFRESH_TOKEN")
	viperInstance.BindEnv("api_key", "CLOUD_STORAGE_API_KEY")

	// Read config file (ignore error if file doesn't exist)
	if err := viperInstance.ReadInConfig(); err != nil {
		// Config file not found is okay, we'll use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	return nil
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig() (*Config, error) {
	if viperInstance == nil {
		if err := InitConfig(); err != nil {
			return nil, err
		}
	}

	var cfg Config
	if err := viperInstance.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Ensure API URL has a default if empty
	if cfg.APIURL == "" {
		cfg.APIURL = defaultAPIURL
	}

	return &cfg, nil
}

// SaveConfig saves configuration to file
func SaveConfig(cfg *Config) error {
	if viperInstance == nil {
		if err := InitConfig(); err != nil {
			return err
		}
	}

	// Set values in Viper
	viperInstance.Set("api_url", cfg.APIURL)
	viperInstance.Set("access_token", cfg.AccessToken)
	viperInstance.Set("refresh_token", cfg.RefreshToken)
	viperInstance.Set("api_key", cfg.APIKey)

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write config file
	if err := viperInstance.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfig returns the current configuration
func GetConfig() (*Config, error) {
	return LoadConfig()
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(homeDir, configDirName, configFileName)
	}
	return configPath
}

// SetValue sets a configuration value by key
func SetValue(key, value string) error {
	if viperInstance == nil {
		if err := InitConfig(); err != nil {
			return err
		}
	}

	// Load current config
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	// Update the specified key
	switch key {
	case "api-url", "api_url":
		cfg.APIURL = value
	case "access-token", "access_token":
		cfg.AccessToken = value
	case "refresh-token", "refresh_token":
		cfg.RefreshToken = value
	case "api-key", "api_key":
		cfg.APIKey = value
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	// Save updated config
	return SaveConfig(cfg)
}

// GetValue gets a configuration value by key
func GetValue(key string) (string, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return "", err
	}

	switch key {
	case "api-url", "api_url":
		return cfg.APIURL, nil
	case "access-token", "access_token":
		return cfg.AccessToken, nil
	case "refresh-token", "refresh_token":
		return cfg.RefreshToken, nil
	case "api-key", "api_key":
		return cfg.APIKey, nil
	default:
		return "", fmt.Errorf("unknown config key: %s", key)
	}
}

// IsSensitiveKey returns true if the key contains sensitive information
func IsSensitiveKey(key string) bool {
	sensitiveKeys := []string{"access-token", "access_token", "refresh-token", "refresh_token", "api-key", "api_key"}
	for _, sk := range sensitiveKeys {
		if key == sk {
			return true
		}
	}
	return false
}

// MaskValue masks sensitive values for display
func MaskValue(value string) string {
	if value == "" {
		return "(not set)"
	}
	if len(value) <= 8 {
		return "***"
	}
	return value[:4] + "..." + value[len(value)-4:]
}


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
package auth

import (
	"fmt"

	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/config"
)

// SaveTokens saves access token and refresh token to configuration
func SaveTokens(accessToken, refreshToken string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	cfg.AccessToken = accessToken
	cfg.RefreshToken = refreshToken

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save tokens: %w", err)
	}

	return nil
}

// ClearTokens clears access token and refresh token from configuration
func ClearTokens() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	cfg.AccessToken = ""
	cfg.RefreshToken = ""

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to clear tokens: %w", err)
	}

	return nil
}

// GetStoredTokens retrieves stored access token and refresh token from configuration
func GetStoredTokens() (accessToken, refreshToken string, err error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", "", fmt.Errorf("failed to load config: %w", err)
	}

	return cfg.AccessToken, cfg.RefreshToken, nil
}


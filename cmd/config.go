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
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/config"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/util"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long: `Manage CLI configuration settings.

Configuration is stored in ~/.cloud-storage-cli/config.yaml

You can view, get, or set configuration values:
  - api-url: API base URL (default: http://localhost:8000)
  - access-token: JWT access token
  - refresh-token: JWT refresh token
  - api-key: API key for authentication

Examples:
  # Show all configuration values
  cloud-storage-api-cli config show

  # Get a specific configuration value
  cloud-storage-api-cli config get api-url

  # Set a configuration value
  cloud-storage-api-cli config set api-url http://api.example.com`,
}

// configShowCmd represents the config show command
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show all configuration values",
	Long:  `Display all configuration values. Sensitive values (tokens, API keys) are masked.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Check if JSON output is requested
		if jsonOutput {
			// For JSON output, create a struct with masked values
			type ConfigOutput struct {
				ConfigFile   string `json:"configFile"`
				APIURL       string `json:"apiUrl"`
				AccessToken  string `json:"accessToken"`
				RefreshToken string `json:"refreshToken"`
				APIKey       string `json:"apiKey"`
			}
			output := ConfigOutput{
				ConfigFile:   config.GetConfigPath(),
				APIURL:       cfg.APIURL,
				AccessToken:  config.MaskValue(cfg.AccessToken),
				RefreshToken: config.MaskValue(cfg.RefreshToken),
				APIKey:       config.MaskValue(cfg.APIKey),
			}
			return util.OutputJSON(output)
		}

		fmt.Println("Configuration:")
		fmt.Println("==============")
		fmt.Printf("Config file: %s\n\n", config.GetConfigPath())
		fmt.Printf("API URL:        %s\n", cfg.APIURL)
		fmt.Printf("Access Token:   %s\n", config.MaskValue(cfg.AccessToken))
		fmt.Printf("Refresh Token:  %s\n", config.MaskValue(cfg.RefreshToken))
		fmt.Printf("API Key:        %s\n", config.MaskValue(cfg.APIKey))

		return nil
	},
}

// configGetCmd represents the config get command
var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a specific configuration value",
	Long: `Get a specific configuration value by key.

Supported keys:
  - api-url
  - access-token
  - refresh-token
  - api-key

Sensitive values are masked when displayed.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value, err := config.GetValue(key)
		if err != nil {
			return fmt.Errorf("failed to get config value: %w", err)
		}

		// Mask sensitive values
		if config.IsSensitiveKey(key) {
			value = config.MaskValue(value)
		}

		fmt.Println(value)
		return nil
	},
}

// configSetCmd represents the config set command
var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value and save it to the config file.

Supported keys:
  - api-url: API base URL
  - access-token: JWT access token
  - refresh-token: JWT refresh token
  - api-key: API key for authentication

Examples:
  cloud-storage-api-cli config set api-url http://api.example.com
  cloud-storage-api-cli config set api-key your-api-key-here`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		// Validate key
		validKeys := []string{"api-url", "api_url", "access-token", "access_token", "refresh-token", "refresh_token", "api-key", "api_key"}
		isValid := false
		for _, vk := range validKeys {
			if key == vk {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("invalid key: %s. Supported keys: api-url, access-token, refresh-token, api-key", key)
		}

		if err := config.SetValue(key, value); err != nil {
			return fmt.Errorf("failed to set config value: %w", err)
		}

		// Display masked value for sensitive keys
		displayValue := value
		if config.IsSensitiveKey(key) {
			displayValue = config.MaskValue(value)
		}

		fmt.Printf("Configuration updated: %s = %s\n", key, displayValue)
		fmt.Printf("Config file: %s\n", config.GetConfigPath())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
}


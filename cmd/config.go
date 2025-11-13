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

You can view or get configuration values. API keys can only be set via the
'auth login' command, which validates the key before saving it.

Note: API URL is configured via the CLOUD_STORAGE_API_URL environment variable
or the --api-url flag. It cannot be set via config command.

Examples:
  # Show all configuration values
  cloud-storage-api-cli config show

  # Get a specific configuration value
  cloud-storage-api-cli config get api-key`,
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
				ConfigFile string `json:"configFile"`
				APIURL     string `json:"apiUrl"`
				APIKey     string `json:"apiKey"`
			}
			output := ConfigOutput{
				ConfigFile: config.GetConfigPath(),
				APIURL:     cfg.APIURL,
				APIKey:     config.MaskValue(cfg.APIKey),
			}
			return util.OutputJSON(output)
		}

		fmt.Println("Configuration:")
		fmt.Println("==============")
		fmt.Printf("Config file: %s\n\n", config.GetConfigPath())
		fmt.Printf("API URL:        %s\n", cfg.APIURL)
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

// configSetCmd is removed - API keys can only be set via 'auth login' command
// which validates the key before saving it. This prevents saving invalid keys.

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configGetCmd)
	// configSetCmd removed - API keys can only be set via 'auth login' command
}

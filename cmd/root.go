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
	"os"

	"github.com/spf13/cobra"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/config"
)

var (
	apiURL     string
	cfgFile    string
	verbose    bool
	jsonOutput bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloud-storage-api-cli",
	Short: "Cloud Storage API CLI - Command-line interface for cloud storage operations",
	Long: `Cloud Storage API CLI is a command-line tool for interacting with the Cloud Storage API.

It provides commands for:
  - Authentication (login, register, logout)
  - File operations (upload, download, list, search, update, delete, info)
  - Folder management (create, list, delete)
  - API key management
  - Batch job status

Examples:
  # Login to the API
  cloud-storage-api-cli auth login username

  # Upload a file
  cloud-storage-api-cli file upload ./document.pdf --folder-path /documents

  # List files
  cloud-storage-api-cli file list --page 0 --size 20

  # Download a file
  cloud-storage-api-cli file download <file-id> --output ./downloaded.pdf

For more information, use 'cloud-storage-api-cli <command> --help'`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Initialize configuration
	if err := config.InitConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize config: %v\n", err)
	}

	// API URL is hardcoded at compile time - always use the build-time value
	apiURL = config.GetAPIURL()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Persistent flags available to all subcommands
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cloud-storage-cli/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output in JSON format")
}

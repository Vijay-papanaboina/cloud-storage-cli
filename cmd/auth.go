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
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/client"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/config"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/util"
	"golang.org/x/term"
)

// Request/Response types matching API DTOs

// LoginRequest is no longer used - CLI uses API key authentication

// RegisterRequest, RefreshTokenRequest, RefreshTokenResponse are no longer used

// UserResponse represents user information
type UserResponse struct {
	ID          string     `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Active      bool       `json:"active"`
	CreatedAt   time.Time  `json:"createdAt"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
}

// AuthResponse is no longer used - CLI uses API key authentication

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
	Long: `Manage authentication using API keys.

Available commands:
  login  - Verify and store API key for authentication
  status - Show current authenticated user information`,
}

// readPassword securely reads a password from stdin without echoing
func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // New line after password input
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	return string(passwordBytes), nil
}

// authLoginCmd represents the auth login command
var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Verify and store API key for authentication",
	Long: `Verify an API key and save it to configuration for future use.

The API key will be prompted securely (not visible as you type).
After verification, the API key will be saved to the configuration file.

You can generate API keys from the web interface at the Settings page.

Examples:
  cloud-storage-api-cli auth login`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Prompt for API key securely
		apiKey, err := readPassword("API Key: ")
		if err != nil {
			return err
		}
		// Trim whitespace and remove all control characters from the API key
		apiKey = strings.TrimSpace(apiKey)
		// Remove any control characters (newlines, carriage returns, etc.)
		apiKey = strings.Map(func(r rune) rune {
			if r >= 32 && r != 127 { // Keep printable ASCII except DEL
				return r
			}
			return -1 // Remove control characters
		}, apiKey)
		if apiKey == "" {
			return fmt.Errorf("API key cannot be empty")
		}

		// Get API URL from config (LoadConfig ensures it's never empty)
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Create API client with the provided API key and base URL
		apiClient := client.NewClientWithConfig(cfg.APIURL, apiKey)

		// Verify API key by calling the verify endpoint
		var userResp UserResponse
		if err := apiClient.Post("/api/api-keys/verify", nil, &userResp); err != nil {
			return fmt.Errorf("API key verification failed: %w", err)
		}

		// Save API key to config
		if err := config.SetValue("api-key", apiKey); err != nil {
			return fmt.Errorf("failed to save API key: %w", err)
		}

		// Display success message
		fmt.Println("API key verified and saved successfully!")
		fmt.Printf("User: %s (%s)\n", userResp.Username, userResp.Email)
		fmt.Printf("User ID: %s\n", userResp.ID)
		fmt.Println("API key saved to configuration.")

		return nil
	},
}

// authStatusCmd represents the auth status command
var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current authenticated user information",
	Long: `Display information about the currently authenticated user.

This command uses the stored API key for authentication.

Examples:
  cloud-storage-api-cli auth status`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create API client (will use stored API key)
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Call verify endpoint to get user info
		var userResp UserResponse
		if err := apiClient.Post("/api/api-keys/verify", nil, &userResp); err != nil {
			return fmt.Errorf("failed to get user information: %w", err)
		}

		// Check if JSON output is requested
		if jsonOutput {
			return util.OutputJSON(userResp)
		}

		// Display user information
		fmt.Println("Current User Information:")
		fmt.Println("========================")
		fmt.Printf("ID:          %s\n", userResp.ID)
		fmt.Printf("Username:    %s\n", userResp.Username)
		fmt.Printf("Email:       %s\n", userResp.Email)
		fmt.Printf("Active:      %v\n", userResp.Active)
		fmt.Printf("Created At:  %s\n", userResp.CreatedAt.Format(time.RFC3339))
		if userResp.LastLoginAt != nil {
			fmt.Printf("Last Login:  %s\n", userResp.LastLoginAt.Format(time.RFC3339))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authStatusCmd)
}

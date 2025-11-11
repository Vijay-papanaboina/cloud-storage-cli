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
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/client"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/util"
)

// ApiKeyRequest represents a request to generate an API key
type ApiKeyRequest struct {
	Name      string     `json:"name"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

// ApiKeyResponse represents API key information
type ApiKeyResponse struct {
	ID         string     `json:"id"`
	Key        *string    `json:"key,omitempty"` // Only present on creation
	Name       string     `json:"name"`
	Active     bool       `json:"active"`
	CreatedAt  time.Time  `json:"createdAt"`
	ExpiresAt  *time.Time `json:"expiresAt,omitempty"`
	LastUsedAt *time.Time `json:"lastUsedAt,omitempty"`
}

// apikeyCmd represents the apikey command
var apikeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "API key management commands",
	Long: `Manage API keys for authentication.

Available commands:
  generate - Generate a new API key
  list     - List all API keys
  get      - Get API key details
  revoke   - Revoke an API key`,
}

// apikeyGenerateCmd represents the apikey generate command
var apikeyGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new API key",
	Long: `Generate a new API key for API authentication.

The API key will only be displayed once. Store it securely.
You can optionally set an expiration date.

Examples:
  cloud-storage-api-cli apikey generate --name "My API Key"
  cloud-storage-api-cli apikey generate --name "Temporary Key" --expires-at "2025-12-31T23:59:59Z"`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		expiresAtStr, _ := cmd.Flags().GetString("expires-at")

		// Validate name is provided
		if name == "" {
			return fmt.Errorf("--name is required")
		}

		// Parse expiration date if provided
		var expiresAt *time.Time
		if expiresAtStr != "" {
			parsed, err := time.Parse(time.RFC3339, expiresAtStr)
			if err != nil {
				return fmt.Errorf("invalid expiration date format (use RFC3339, e.g., 2025-12-31T23:59:59Z): %w", err)
			}
			expiresAt = &parsed
		}

		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Build generate request
		generateReq := ApiKeyRequest{
			Name: name,
		}
		if expiresAt != nil {
			generateReq.ExpiresAt = expiresAt
		}

		// Generate API key
		var apiKeyResp ApiKeyResponse
		if err := apiClient.Post("/api/auth/api-keys", generateReq, &apiKeyResp); err != nil {
			return fmt.Errorf("failed to generate API key: %w", err)
		}

		// Display success message with security warning
		fmt.Println("API key generated successfully!")
		fmt.Print("\n⚠️  SECURITY WARNING: This API key will only be displayed once.\n")
		fmt.Print("   Store it securely. You will not be able to retrieve it again.\n\n")
		fmt.Printf("API Key ID: %s\n", apiKeyResp.ID)
		fmt.Printf("Name: %s\n", apiKeyResp.Name)
		if apiKeyResp.Key != nil {
			fmt.Printf("API Key: %s\n", *apiKeyResp.Key)
		}
		fmt.Printf("Active: %v\n", apiKeyResp.Active)
		fmt.Printf("Created At: %s\n", apiKeyResp.CreatedAt.Format(time.RFC3339))
		if apiKeyResp.ExpiresAt != nil {
			fmt.Printf("Expires At: %s\n", apiKeyResp.ExpiresAt.Format(time.RFC3339))
		} else {
			fmt.Println("Expires At: Never")
		}
		fmt.Print("\nTo use this API key, set it in your config:\n")
		if apiKeyResp.Key == nil {
			return fmt.Errorf("API key was not returned by the server")
		}
		fmt.Printf("  cloud-storage-api-cli config set apiKey %s\n", *apiKeyResp.Key)
		return nil
	},
}

// apikeyListCmd represents the apikey list command
var apikeyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API keys",
	Long: `List all API keys for your account.

The API key values are not displayed for security reasons.

Examples:
  cloud-storage-api-cli apikey list`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Fetch API key list
		var apiKeys []ApiKeyResponse
		if err := apiClient.Get("/api/auth/api-keys", &apiKeys); err != nil {
			return fmt.Errorf("failed to list API keys: %w", err)
		}

		// Display results
		displayApiKeyList(apiKeys)

		return nil
	},
}

// apikeyGetCmd represents the apikey get command
var apikeyGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get API key details",
	Long: `Get detailed information about a specific API key.

The API key value is not displayed for security reasons.

Examples:
  cloud-storage-api-cli apikey get 660e8400-e29b-41d4-a716-446655440000`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKeyID := args[0]

		// Validate UUID format
		if err := util.ValidateUUID(apiKeyID); err != nil {
			return fmt.Errorf("invalid API key ID: %w", err)
		}
		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Get API key
		path := fmt.Sprintf("/api/auth/api-keys/%s", apiKeyID)
		var apiKeyResp ApiKeyResponse
		if err := apiClient.Get(path, &apiKeyResp); err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		// Display API key details
		displayApiKeyDetails(&apiKeyResp)

		return nil
	},
}

// apikeyRevokeCmd represents the apikey revoke command
var apikeyRevokeCmd = &cobra.Command{
	Use:   "revoke <id>",
	Short: "Revoke an API key",
	Long: `Revoke (deactivate) an API key.

This operation cannot be undone. The API key will no longer be usable for authentication.
You will be prompted for confirmation unless the --force flag is used.

Examples:
  cloud-storage-api-cli apikey revoke 660e8400-e29b-41d4-a716-446655440000
  cloud-storage-api-cli apikey revoke 660e8400-e29b-41d4-a716-446655440000 --force`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKeyID := args[0]
		force, _ := cmd.Flags().GetBool("force")

		// Basic UUID format validation
		if _, err := uuid.Parse(apiKeyID); err != nil {
			return fmt.Errorf("invalid API key ID format (expected UUID): %s", apiKeyID)
		}

		// Prompt for confirmation if not forced
		if !force {
			fmt.Printf("Are you sure you want to revoke API key %s? This cannot be undone. (y/N): ", apiKeyID)
			var response string
			if _, err := fmt.Scanln(&response); err != nil {
				fmt.Println("\nRevocation cancelled.")
				return nil
			}
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				fmt.Println("Revocation cancelled.")
				return nil
			}
		}

		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Revoke API key
		path := fmt.Sprintf("/api/auth/api-keys/%s", apiKeyID)
		if err := apiClient.Delete(path); err != nil {
			return fmt.Errorf("failed to revoke API key: %w", err)
		}

		// Display success message
		fmt.Printf("API key %s revoked successfully.\n", apiKeyID)

		return nil
	},
}

// displayApiKeyList displays the API key list in a formatted table
func displayApiKeyList(apiKeys []ApiKeyResponse) {
	if len(apiKeys) == 0 {
		fmt.Println("No API keys found.")
		return
	}

	// Print header
	fmt.Printf("\nAPI Keys (Total: %d)\n\n", len(apiKeys))

	// Print table header
	fmt.Printf("%-36s %-30s %-10s %-20s %-20s %-20s\n",
		"ID", "Name", "Active", "Created At", "Expires At", "Last Used At")
	fmt.Println(strings.Repeat("-", 140))

	// Sort by created date (newest first)
	sort.Slice(apiKeys, func(i, j int) bool {
		return apiKeys[i].CreatedAt.After(apiKeys[j].CreatedAt)
	})

	// Print table rows
	for _, key := range apiKeys {
		// Truncate ID if too long
		id := key.ID
		if len(id) > 36 {
			id = id[:36]
		}

		// Truncate name if too long
		name := key.Name
		if len(name) > 30 {
			name = name[:27] + "..."
		}

		// Format active status
		active := "Active"
		if !key.Active {
			active = "Inactive"
		}

		// Format dates
		createdAt := key.CreatedAt.Format("2006-01-02 15:04:05")
		expiresAt := "Never"
		if key.ExpiresAt != nil {
			expiresAt = key.ExpiresAt.Format("2006-01-02 15:04:05")
		}
		lastUsedAt := "Never"
		if key.LastUsedAt != nil {
			lastUsedAt = key.LastUsedAt.Format("2006-01-02 15:04:05")
		}

		fmt.Printf("%-36s %-30s %-10s %-20s %-20s %-20s\n",
			id, name, active, createdAt, expiresAt, lastUsedAt)
	}

	fmt.Println(strings.Repeat("-", 140))
	fmt.Println()
}

// displayApiKeyDetails displays API key details
func displayApiKeyDetails(apiKey *ApiKeyResponse) {
	fmt.Println("\nAPI Key Details")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("ID:         %s\n", apiKey.ID)
	fmt.Printf("Name:       %s\n", apiKey.Name)
	fmt.Printf("Active:     %v\n", apiKey.Active)
	fmt.Printf("Created At: %s\n", apiKey.CreatedAt.Format(time.RFC3339))
	if apiKey.ExpiresAt != nil {
		fmt.Printf("Expires At: %s\n", apiKey.ExpiresAt.Format(time.RFC3339))
	} else {
		fmt.Println("Expires At: Never")
	}
	if apiKey.LastUsedAt != nil {
		fmt.Printf("Last Used:  %s\n", apiKey.LastUsedAt.Format(time.RFC3339))
	} else {
		fmt.Println("Last Used:  Never")
	}
	if apiKey.Key != nil {
		fmt.Printf("API Key:    %s\n", *apiKey.Key)
		fmt.Println("\n⚠️  Note: API key value is only shown on creation.")
	}
	fmt.Println()
}

func init() {
	// Add apikey command to root
	rootCmd.AddCommand(apikeyCmd)

	// Add generate subcommand to apikey command
	apikeyCmd.AddCommand(apikeyGenerateCmd)

	// Add list subcommand to apikey command
	apikeyCmd.AddCommand(apikeyListCmd)

	// Add get subcommand to apikey command
	apikeyCmd.AddCommand(apikeyGetCmd)

	// Add revoke subcommand to apikey command
	apikeyCmd.AddCommand(apikeyRevokeCmd)

	// Add flags to generate command
	apikeyGenerateCmd.Flags().String("name", "", "API key name (required)")
	apikeyGenerateCmd.MarkFlagRequired("name")
	apikeyGenerateCmd.Flags().String("expires-at", "", "Expiration date in RFC3339 format (e.g., 2025-12-31T23:59:59Z)")

	// Add flags to revoke command
	apikeyRevokeCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
}

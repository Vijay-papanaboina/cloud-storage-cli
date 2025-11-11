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
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/auth"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/client"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/util"
	"golang.org/x/term"
)

// Request/Response types matching API DTOs

// LoginRequest represents a login request
type LoginRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	ClientType string `json:"clientType,omitempty"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// UserResponse represents user information
type UserResponse struct {
	ID          string     `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Active      bool       `json:"active"`
	CreatedAt   time.Time  `json:"createdAt"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	AccessToken      string       `json:"accessToken"`
	RefreshToken     string       `json:"refreshToken"`
	TokenType        string       `json:"tokenType"`
	ExpiresIn        int64        `json:"expiresIn"`
	RefreshExpiresIn int64        `json:"refreshExpiresIn"`
	ClientType       string       `json:"clientType"`
	User             UserResponse `json:"user"`
}

// RefreshTokenResponse represents a refresh token response
type RefreshTokenResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int64  `json:"expiresIn"`
}

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
	Long: `Manage authentication and user account.

Available commands:
  login    - Login with username and password
  register  - Register a new user account
  logout    - Logout and clear stored tokens
  refresh   - Refresh access token using refresh token
  me        - Show current authenticated user information`,
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
	Use:   "login <username>",
	Short: "Login with username and password",
	Long: `Login to the API and save authentication tokens.

The password will be prompted securely (not visible as you type).
The tokens will be saved to the configuration file for future use.
Use clientType "CLI" for longer token expiry (30 days for access token).

Examples:
  cloud-storage-api-cli auth login myuser
  cloud-storage-api-cli auth login admin`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]

		// Validate username
		if err := util.ValidateUsername(username); err != nil {
			return err
		}

		// Prompt for password securely
		password, err := readPassword("Password: ")
		if err != nil {
			return err
		}
		if password == "" {
			return fmt.Errorf("password cannot be empty")
		}

		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Prepare login request
		loginReq := LoginRequest{
			Username:   username,
			Password:   password,
			ClientType: "CLI", // Use CLI for longer token expiry
		}

		// Call login endpoint
		var authResp AuthResponse
		if err := apiClient.Post("/api/auth/login", loginReq, &authResp); err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		// Save tokens to config
		if err := auth.SaveTokens(authResp.AccessToken, authResp.RefreshToken); err != nil {
			return fmt.Errorf("failed to save tokens: %w", err)
		}

		// Display success message
		fmt.Println("Login successful!")
		fmt.Printf("User: %s (%s)\n", authResp.User.Username, authResp.User.Email)
		fmt.Printf("Access token expires in: %d seconds\n", authResp.ExpiresIn)
		fmt.Printf("Refresh token expires in: %d seconds\n", authResp.RefreshExpiresIn)
		fmt.Println("Tokens saved to configuration.")

		return nil
	},
}

// authRegisterCmd represents the auth register command
var authRegisterCmd = &cobra.Command{
	Use:   "register <username> <email>",
	Short: "Register a new user account",
	Long: `Register a new user account with the API.

The password will be prompted securely (not visible as you type).
After registration, you can login using the auth login command.

Examples:
  cloud-storage-api-cli auth register myuser user@example.com`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]
		email := args[1]

		// Validate username
		if err := util.ValidateUsername(username); err != nil {
			return err
		}
		// Validate email
		if err := util.ValidateEmail(email); err != nil {
			return err
		}

		// Prompt for password securely
		password, err := readPassword("Password: ")
		if err != nil {
			return err
		}
		if password == "" {
			return fmt.Errorf("password cannot be empty")
		}

		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Prepare register request
		registerReq := RegisterRequest{
			Username: username,
			Email:    email,
			Password: password,
		}

		// Call register endpoint
		var userResp UserResponse
		if err := apiClient.Post("/api/auth/register", registerReq, &userResp); err != nil {
			return fmt.Errorf("registration failed: %w", err)
		}

		// Display success message
		fmt.Println("Registration successful!")
		fmt.Printf("User ID: %s\n", userResp.ID)
		fmt.Printf("Username: %s\n", userResp.Username)
		fmt.Printf("Email: %s\n", userResp.Email)
		fmt.Println("\nYou can now login using: cloud-storage-api-cli auth login <username> <password>")

		return nil
	},
}

// authLogoutCmd represents the auth logout command
var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout and clear stored tokens",
	Long: `Logout from the API and clear stored authentication tokens.

This will:
  - Call the API logout endpoint to invalidate the refresh token (if available)
  - Clear stored tokens from configuration

Examples:
  cloud-storage-api-cli auth logout`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get stored refresh token
		_, refreshToken, err := auth.GetStoredTokens()
		if err != nil {
			return fmt.Errorf("failed to get stored tokens: %w", err)
		}

		// If refresh token exists, call logout endpoint
		if refreshToken != "" {
			apiClient, err := client.NewClient()
			if err != nil {
				// If client creation fails, still try to clear tokens
				fmt.Fprintf(os.Stderr, "Warning: Failed to create API client, clearing local tokens only\n")
			} else {
				// Call logout endpoint
				logoutReq := RefreshTokenRequest{
					RefreshToken: refreshToken,
				}
				if err := apiClient.Post("/api/auth/logout", logoutReq, nil); err != nil {
					// Log error but continue to clear local tokens
					fmt.Fprintf(os.Stderr, "Warning: Failed to call logout endpoint: %v\n", err)
				}
			}
		}

		// Clear tokens from config
		if err := auth.ClearTokens(); err != nil {
			return fmt.Errorf("failed to clear tokens: %w", err)
		}

		fmt.Println("Logged out successfully. Tokens cleared from configuration.")
		return nil
	},
}

// authRefreshCmd represents the auth refresh command
var authRefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh access token using refresh token",
	Long: `Refresh the access token using the stored refresh token.

The new access token will be saved to configuration automatically.

Examples:
  cloud-storage-api-cli auth refresh`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get stored refresh token
		_, refreshToken, err := auth.GetStoredTokens()
		if err != nil {
			return fmt.Errorf("failed to get stored tokens: %w", err)
		}

		if refreshToken == "" {
			return fmt.Errorf("no refresh token found. Please login first using 'auth login'")
		}

		// Create API client
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Prepare refresh request
		refreshReq := RefreshTokenRequest{
			RefreshToken: refreshToken,
		}

		// Call refresh endpoint
		var refreshResp RefreshTokenResponse
		if err := apiClient.Post("/api/auth/refresh", refreshReq, &refreshResp); err != nil {
			return fmt.Errorf("token refresh failed: %w", err)
		}

		// Save new access token (keep existing refresh token)
		if err := auth.SaveTokens(refreshResp.AccessToken, refreshToken); err != nil {
			return fmt.Errorf("failed to save new access token: %w", err)
		}

		fmt.Println("Token refreshed successfully!")
		fmt.Printf("New access token expires in: %d seconds\n", refreshResp.ExpiresIn)
		fmt.Println("New access token saved to configuration.")

		return nil
	},
}

// authMeCmd represents the auth me command
var authMeCmd = &cobra.Command{
	Use:   "me",
	Short: "Show current authenticated user information",
	Long: `Display information about the currently authenticated user.

This command uses the stored access token for authentication.

Examples:
  cloud-storage-api-cli auth me`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create API client (will use stored access token)
		apiClient, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create API client: %w", err)
		}

		// Call me endpoint
		var userResp UserResponse
		if err := apiClient.Get("/api/auth/me", &userResp); err != nil {
			return fmt.Errorf("failed to get user information: %w", err)
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
	authCmd.AddCommand(authRegisterCmd)
	authCmd.AddCommand(authLogoutCmd)
	authCmd.AddCommand(authRefreshCmd)
	authCmd.AddCommand(authMeCmd)
}

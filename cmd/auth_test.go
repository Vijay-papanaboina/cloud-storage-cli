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
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/client"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/testutil"
)

func TestAuthLogin_Integration(t *testing.T) {
	// Setup mock server
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			testutil.ErrorResponse(w, http.StatusMethodNotAllowed, "Expected POST")
			return
		}

		// Parse login request
		var loginReq LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
			testutil.ErrorResponse(w, http.StatusBadRequest, "Failed to decode request")
			return
		}

		// Verify login credentials
		if loginReq.Username != "testuser" {
			testutil.ErrorResponse(w, http.StatusBadRequest, "Expected username 'testuser'")
			return
		}

		if loginReq.Password != "testpass" {
			testutil.ErrorResponse(w, http.StatusBadRequest, "Expected password 'testpass'")
			return
		}

		// Return auth response
		response := AuthResponse{
			AccessToken:      "access-token-123",
			RefreshToken:     "refresh-token-456",
			TokenType:        "Bearer",
			ExpiresIn:        3600,
			RefreshExpiresIn: 86400,
			ClientType:       "CLI",
			User: UserResponse{
				ID:        "user-123",
				Username:  "testuser",
				Email:     "test@example.com",
				Active:    true,
				CreatedAt: time.Now(),
			},
		}

		testutil.JSONResponse(w, http.StatusOK, response)
	})
	defer server.Close()

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "", "")

	// Test login
	loginReq := LoginRequest{
		Username:   "testuser",
		Password:   "testpass",
		ClientType: "CLI",
	}

	var authResp AuthResponse
	err := apiClient.Post("/api/auth/login", loginReq, &authResp)

	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if authResp.AccessToken != "access-token-123" {
		t.Errorf("Expected access token 'access-token-123', got %q", authResp.AccessToken)
	}

	if authResp.User.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %q", authResp.User.Username)
	}
}

func TestAuthRegister_Integration(t *testing.T) {
	// Setup mock server
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}

		// Parse register request
		var registerReq RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			return
		}

		// Verify registration data
		if registerReq.Username != "newuser" {
			t.Errorf("Expected username 'newuser', got %q", registerReq.Username)
		}

		if registerReq.Email != "newuser@example.com" {
			t.Errorf("Expected email 'newuser@example.com', got %q", registerReq.Email)
		}

		// Return auth response
		response := AuthResponse{
			AccessToken:  "access-token-new",
			RefreshToken: "refresh-token-new",
			TokenType:    "Bearer",
			User: UserResponse{
				ID:        "user-new",
				Username:  "newuser",
				Email:     "newuser@example.com",
				Active:    true,
				CreatedAt: time.Now(),
			},
		}

		testutil.JSONResponse(w, http.StatusCreated, response)
	})
	defer server.Close()

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "", "")

	// Test register
	registerReq := RegisterRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "newpass123",
	}

	var authResp AuthResponse
	err := apiClient.Post("/api/auth/register", registerReq, &authResp)

	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if authResp.User.Username != "newuser" {
		t.Errorf("Expected username 'newuser', got %q", authResp.User.Username)
	}
}

func TestAuthRefresh_Integration(t *testing.T) {
	// Setup mock server
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}

		// Parse refresh request
		var refreshReq RefreshTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&refreshReq); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			return
		}

		if refreshReq.RefreshToken != "refresh-token-456" {
			t.Errorf("Expected refresh token 'refresh-token-456', got %q", refreshReq.RefreshToken)
		}

		// Return refresh response
		response := RefreshTokenResponse{
			AccessToken: "new-access-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		}

		testutil.JSONResponse(w, http.StatusOK, response)
	})
	defer server.Close()

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "", "")

	// Test refresh
	refreshReq := RefreshTokenRequest{
		RefreshToken: "refresh-token-456",
	}

	var refreshResp RefreshTokenResponse
	err := apiClient.Post("/api/auth/refresh", refreshReq, &refreshResp)

	if err != nil {
		t.Fatalf("Refresh failed: %v", err)
	}

	if refreshResp.AccessToken != "new-access-token" {
		t.Errorf("Expected access token 'new-access-token', got %q", refreshResp.AccessToken)
	}
}

func TestAuthMe_Integration(t *testing.T) {
	// Setup mock server
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}

		// Verify authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got %q", auth)
		}

		// Return user response
		response := UserResponse{
			ID:        "user-123",
			Username:  "testuser",
			Email:     "test@example.com",
			Active:    true,
			CreatedAt: time.Now(),
		}

		testutil.JSONResponse(w, http.StatusOK, response)
	})
	defer server.Close()

	// Create client with token
	apiClient := client.NewClientWithConfig(server.URL, "test-token", "")

	// Test me
	var userResp UserResponse
	err := apiClient.Get("/api/auth/me", &userResp)

	if err != nil {
		t.Fatalf("Me failed: %v", err)
	}

	if userResp.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %q", userResp.Username)
	}
}

func TestAuthLogin_ErrorHandling(t *testing.T) {
	// Setup mock server with error response
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		testutil.ErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
	})
	defer server.Close()

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "", "")

	// Test login with invalid credentials
	loginReq := LoginRequest{
		Username: "wronguser",
		Password: "wrongpass",
	}

	var authResp AuthResponse
	err := apiClient.Post("/api/auth/login", loginReq, &authResp)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestAuthRefresh_ErrorHandling(t *testing.T) {
	// Setup mock server with error response
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		testutil.ErrorResponse(w, http.StatusUnauthorized, "Invalid refresh token")
	})
	defer server.Close()

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "", "")

	// Test refresh with invalid token
	refreshReq := RefreshTokenRequest{
		RefreshToken: "invalid-token",
	}

	var refreshResp RefreshTokenResponse
	err := apiClient.Post("/api/auth/refresh", refreshReq, &refreshResp)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

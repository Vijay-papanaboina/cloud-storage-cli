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

		// Verify API key header
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != "test-api-key-123" {
			testutil.ErrorResponse(w, http.StatusUnauthorized, "Invalid API key")
			return
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

	// Create client with API key
	apiClient := client.NewClientWithConfig(server.URL, "test-api-key-123")

	// Test login (verify API key)
	var userResp UserResponse
	err := apiClient.Post("/api/api-keys/verify", nil, &userResp)

	if err != nil {
		t.Fatalf("API key verification failed: %v", err)
	}

	if userResp.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %q", userResp.Username)
	}

	if userResp.ID != "user-123" {
		t.Errorf("Expected user ID 'user-123', got %q", userResp.ID)
	}
}

func TestAuthStatus_Integration(t *testing.T) {
	// Setup mock server
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}

		// Verify API key header
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != "test-api-key-456" {
			testutil.ErrorResponse(w, http.StatusUnauthorized, "Invalid API key")
			return
		}

		// Return user response
		response := UserResponse{
			ID:        "user-456",
			Username:  "testuser",
			Email:     "test@example.com",
			Active:    true,
			CreatedAt: time.Now(),
		}

		testutil.JSONResponse(w, http.StatusOK, response)
	})
	defer server.Close()

	// Create client with API key
	apiClient := client.NewClientWithConfig(server.URL, "test-api-key-456")

	// Test status (verify API key and get user info)
	var userResp UserResponse
	err := apiClient.Post("/api/api-keys/verify", nil, &userResp)

	if err != nil {
		t.Fatalf("Status check failed: %v", err)
	}

	if userResp.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %q", userResp.Username)
	}

	if userResp.ID != "user-456" {
		t.Errorf("Expected user ID 'user-456', got %q", userResp.ID)
	}
}

// runInvalidAPIKeyTest sets up a test server that returns Unauthorized,
// creates a client with an invalid API key, performs a POST to /api/api-keys/verify,
// and asserts that an error is returned.
func runInvalidAPIKeyTest(t *testing.T) {
	// Setup mock server with error response
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		testutil.ErrorResponse(w, http.StatusUnauthorized, "Invalid API key")
	})
	defer server.Close()

	// Create client with invalid API key
	apiClient := client.NewClientWithConfig(server.URL, "invalid-key")

	// Test API key verification with invalid key
	var userResp UserResponse
	err := apiClient.Post("/api/api-keys/verify", nil, &userResp)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestAuthLogin_ErrorHandling(t *testing.T) {
	runInvalidAPIKeyTest(t)
}

func TestAuthStatus_ErrorHandling(t *testing.T) {
	runInvalidAPIKeyTest(t)
}

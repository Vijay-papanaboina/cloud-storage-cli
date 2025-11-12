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
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/file"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/testutil"
)

func TestFolderCreate_Integration(t *testing.T) {
	// Setup mock server
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}

		// Parse create request
		var createReq file.FolderCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			return
		}

		// Verify folder path
		if createReq.Path != "/documents" {
			t.Errorf("Expected path '/documents', got %q", createReq.Path)
		}

		// Return folder response
		response := file.FolderResponse{
			Path:      "/documents",
			FileCount: 0,
			CreatedAt: time.Now(),
		}

		testutil.JSONResponse(w, http.StatusCreated, response)
	})
	defer server.Close()

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "test-token", "")

	// Test create
	createReq := file.FolderCreateRequest{
		Path: "/documents",
	}

	var folderResp file.FolderResponse
	err := apiClient.Post("/api/folders", createReq, &folderResp)

	if err != nil {
		t.Fatalf("Create folder failed: %v", err)
	}

	if folderResp.Path != "/documents" {
		t.Errorf("Expected path '/documents', got %q", folderResp.Path)
	}
}

func TestFolderList_Integration(t *testing.T) {
	// Setup mock server
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}

		// Verify query parameter (if needed)
		_ = r.URL.Query().Get("parentPath")

		// Return folder list
		response := file.FolderListResponse{
			Folders: []file.FolderResponse{
				{
					Path:      "/documents",
					FileCount: 5,
					CreatedAt: time.Now(),
				},
				{
					Path:      "/photos",
					FileCount: 10,
					CreatedAt: time.Now(),
				},
			},
		}

		testutil.JSONResponse(w, http.StatusOK, response)
	})
	defer server.Close()

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "test-token", "")

	// Test list
	var listResp file.FolderListResponse
	err := apiClient.Get("/api/folders", &listResp)

	if err != nil {
		t.Fatalf("List folders failed: %v", err)
	}

	if len(listResp.Folders) != 2 {
		t.Errorf("Expected 2 folders, got %d", len(listResp.Folders))
	}
}

func TestFolderDelete_Integration(t *testing.T) {
	// Setup mock server
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}

		// Verify path parameter
		path := r.URL.Query().Get("path")
		if path != "/documents" {
			t.Errorf("Expected path '/documents', got %q", path)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "test-token", "")

	// Test delete
	err := apiClient.Delete("/api/folders?path=/documents")

	if err != nil {
		t.Fatalf("Delete folder failed: %v", err)
	}
}

func TestFolderStats_Integration(t *testing.T) {
	// Setup mock server
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}

		// Verify path parameter
		path := r.URL.Query().Get("path")
		if path != "/documents" {
			t.Errorf("Expected path '/documents', got %q", path)
		}

		// Return folder statistics
		response := file.FolderStatisticsResponse{
			Path:            "/documents",
			TotalFiles:      5,
			TotalSize:       102400,
			AverageFileSize: 20480,
			StorageUsed:     "100 KB",
			ByContentType: map[string]int64{
				"text/plain":      3,
				"application/pdf": 2,
			},
			CreatedAt: time.Now(),
		}

		testutil.JSONResponse(w, http.StatusOK, response)
	})
	defer server.Close()

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "test-token", "")

	// Test stats
	var statsResp file.FolderStatisticsResponse
	err := apiClient.Get("/api/folders/statistics?path=/documents", &statsResp)

	if err != nil {
		t.Fatalf("Folder stats failed: %v", err)
	}

	if statsResp.TotalFiles != 5 {
		t.Errorf("Expected 5 total files, got %d", statsResp.TotalFiles)
	}

	if statsResp.Path != "/documents" {
		t.Errorf("Expected path '/documents', got %q", statsResp.Path)
	}
}

func TestFolderCreate_ErrorHandling(t *testing.T) {
	// Setup mock server with error response
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		testutil.ErrorResponse(w, http.StatusBadRequest, "Invalid folder path")
	})
	defer server.Close()

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "test-token", "")

	// Test create with invalid path
	createReq := file.FolderCreateRequest{
		Path: "/invalid/../path",
	}

	var folderResp file.FolderResponse
	err := apiClient.Post("/api/folders", createReq, &folderResp)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestFolderDelete_ErrorHandling(t *testing.T) {
	// Setup mock server with error response
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		testutil.ErrorResponse(w, http.StatusNotFound, "Folder not found")
	})
	defer server.Close()

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "test-token", "")

	// Test delete with non-existent folder
	err := apiClient.Delete("/api/folders?path=/nonexistent")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestFolderStats_ErrorHandling(t *testing.T) {
	// Setup mock server with error response
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		testutil.ErrorResponse(w, http.StatusNotFound, "Folder not found")
	})
	defer server.Close()

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "test-token", "")

	// Test stats with non-existent folder
	var statsResp file.FolderStatisticsResponse
	err := apiClient.Get("/api/folders/statistics?path=/nonexistent", &statsResp)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

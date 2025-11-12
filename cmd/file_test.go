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
	"os"
	"testing"
	"time"

	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/client"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/file"
	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/testutil"
)

func TestFileUpload_Integration(t *testing.T) {
	// Create test file
	testFile := testutil.CreateTestFile(t, "test file content")
	defer os.Remove(testFile)

	// Setup mock server
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}

		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			t.Errorf("Failed to parse multipart form: %v", err)
			return
		}

		// Verify file was uploaded
		uploadedFile, header, err := r.FormFile("file")
		if err != nil {
			t.Errorf("Failed to get file: %v", err)
			return
		}
		defer uploadedFile.Close()

		if header.Filename == "" {
			t.Error("Expected filename in upload")
		}

		// Return success response
		response := file.FileResponse{
			ID:                  "test-id-123",
			Filename:            header.Filename,
			ContentType:         "text/plain",
			FileSize:            18,
			CloudinaryUrl:       "http://cloudinary.com/test.jpg",
			CloudinarySecureUrl: "https://cloudinary.com/test.jpg",
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		testutil.JSONResponse(w, http.StatusCreated, response)
	})
	defer server.Close()

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "test-token", "")

	// Test upload
	var fileResp file.FileResponse
	err := apiClient.UploadFile("/api/files/upload", testFile, "/documents", "", &fileResp)

	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}

	if fileResp.ID != "test-id-123" {
		t.Errorf("Expected file ID 'test-id-123', got %q", fileResp.ID)
	}
}

func TestFileList_Integration(t *testing.T) {
	// Setup mock server
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}

		// Verify query parameters
		page := r.URL.Query().Get("page")
		size := r.URL.Query().Get("size")

		if page == "" || size == "" {
			t.Error("Expected page and size query parameters")
		}

		// Return paginated response
		response := file.PageResponse{
			Content: []file.FileResponse{
				{
					ID:        "file-1",
					Filename:  "test1.txt",
					FileSize:  100,
					CreatedAt: time.Now(),
				},
				{
					ID:        "file-2",
					Filename:  "test2.txt",
					FileSize:  200,
					CreatedAt: time.Now(),
				},
			},
			TotalElements:    2,
			TotalPages:       1,
			First:            true,
			Last:             true,
			NumberOfElements: 2,
			Pageable: file.PageableResponse{
				PageNumber: 0,
				PageSize:   20,
			},
		}

		testutil.JSONResponse(w, http.StatusOK, response)
	})
	defer server.Close()

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "test-token", "")

	// Test list
	var pageResp file.PageResponse
	err := apiClient.Get("/api/files?page=0&size=20", &pageResp)

	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(pageResp.Content) != 2 {
		t.Errorf("Expected 2 files, got %d", len(pageResp.Content))
	}
}

func TestFileSearch_Integration(t *testing.T) {
	// Setup mock server
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}

		// Verify search query
		query := r.URL.Query().Get("q")
		if query != "test" {
			t.Errorf("Expected query 'test', got %q", query)
		}

		// Return search results
		response := file.PageResponse{
			Content: []file.FileResponse{
				{
					ID:       "file-1",
					Filename: "test-file.txt",
					FileSize: 100,
				},
			},
			TotalElements: 1,
		}

		testutil.JSONResponse(w, http.StatusOK, response)
	})
	defer server.Close()

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "test-token", "")

	// Test search
	var pageResp file.PageResponse
	err := apiClient.Get("/api/files/search?q=test&page=0&size=20", &pageResp)

	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(pageResp.Content) != 1 {
		t.Errorf("Expected 1 result, got %d", len(pageResp.Content))
	}
}

func TestFileDownload_Integration(t *testing.T) {
	fileContent := "downloaded file content"

	// Setup mock server
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Disposition", `attachment; filename="downloaded.txt"`)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fileContent))
	})
	defer server.Close()

	// Create output directory
	outputDir := testutil.CreateTestDir(t)

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "test-token", "")

	// Test download
	filePath, err := apiClient.DownloadFile("/api/files/123/download", outputDir)

	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}

	// Verify file was created
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	if string(content) != fileContent {
		t.Errorf("Expected content %q, got %q", fileContent, string(content))
	}
}

func TestFileUpdate_Integration(t *testing.T) {
	// Setup mock server
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT, got %s", r.Method)
		}

		// Parse request body
		var updateReq file.FileUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			return
		}

		// Verify update request
		if updateReq.Filename == nil || *updateReq.Filename != "new-name.txt" {
			t.Error("Expected filename in update request")
		}

		// Return updated file
		response := file.FileResponse{
			ID:       "file-123",
			Filename: "new-name.txt",
			FileSize: 100,
		}

		testutil.JSONResponse(w, http.StatusOK, response)
	})
	defer server.Close()

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "test-token", "")

	// Test update
	newName := "new-name.txt"
	updateReq := file.FileUpdateRequest{
		Filename: &newName,
	}

	var fileResp file.FileResponse
	err := apiClient.Put("/api/files/123", updateReq, &fileResp)

	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if fileResp.Filename != "new-name.txt" {
		t.Errorf("Expected filename 'new-name.txt', got %q", fileResp.Filename)
	}
}

func TestFileDelete_Integration(t *testing.T) {
	// Setup mock server
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "test-token", "")

	// Test delete
	err := apiClient.Delete("/api/files/123")

	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestFileInfo_Integration(t *testing.T) {
	// Setup mock server
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}

		// Return statistics
		response := file.FileStatisticsResponse{
			TotalFiles:      10,
			TotalSize:       102400,
			AverageFileSize: 10240,
			StorageUsed:     "100 KB",
			ByContentType: map[string]int64{
				"text/plain": 5,
				"image/jpeg": 5,
			},
			ByFolder: map[string]int64{
				"/documents": 5,
				"/photos":    5,
			},
		}

		testutil.JSONResponse(w, http.StatusOK, response)
	})
	defer server.Close()

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "test-token", "")

	// Test info
	var stats file.FileStatisticsResponse
	err := apiClient.Get("/api/files/info", &stats)

	if err != nil {
		t.Fatalf("Info failed: %v", err)
	}

	if stats.TotalFiles != 10 {
		t.Errorf("Expected 10 total files, got %d", stats.TotalFiles)
	}
}

func TestFileUpload_ErrorHandling(t *testing.T) {
	// Setup mock server with error response
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		testutil.ErrorResponse(w, http.StatusBadRequest, "File upload failed")
	})
	defer server.Close()

	// Create test file
	testFile := testutil.CreateTestFile(t, "test content")
	defer os.Remove(testFile)

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "test-token", "")

	// Test upload with error
	var fileResp file.FileResponse
	err := apiClient.UploadFile("/api/files/upload", testFile, "", "", &fileResp)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestFileDownload_ErrorHandling(t *testing.T) {
	// Setup mock server with 404 error
	server := testutil.SetupTestServer(func(w http.ResponseWriter, r *http.Request) {
		testutil.ErrorResponse(w, http.StatusNotFound, "File not found")
	})
	defer server.Close()

	// Create output directory
	outputDir := testutil.CreateTestDir(t)

	// Create client
	apiClient := client.NewClientWithConfig(server.URL, "test-token", "")

	// Test download with error
	_, err := apiClient.DownloadFile("/api/files/999/download", outputDir)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

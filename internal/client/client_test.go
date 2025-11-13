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
package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// setupTestServer creates a mock HTTP server for testing
func setupTestServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func TestClient_Get(t *testing.T) {
	tests := []struct {
		name         string
		responseBody interface{}
		statusCode   int
		wantErr      bool
		checkAuth    bool
		expectedAuth string
		authType     string // "token" or "apikey"
	}{
		{
			name: "successful GET request",
			responseBody: map[string]interface{}{
				"id":   "123",
				"name": "test",
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "GET with API key",
			responseBody: map[string]interface{}{
				"id": "123",
			},
			statusCode:   http.StatusOK,
			wantErr:      false,
			checkAuth:    true,
			expectedAuth: "test-api-key",
			authType:     "apikey",
		},
		{
			name: "GET with 404 error",
			responseBody: map[string]interface{}{
				"message": "Not found",
			},
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name: "GET with 500 error",
			responseBody: map[string]interface{}{
				"message": "Internal server error",
			},
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
				// Check authentication if needed
				if tt.checkAuth {
					if tt.authType == "token" {
						auth := r.Header.Get("Authorization")
						if auth != tt.expectedAuth {
							t.Errorf("Expected Authorization header %q, got %q", tt.expectedAuth, auth)
						}
					} else if tt.authType == "apikey" {
						apiKey := r.Header.Get("X-API-Key")
						if apiKey != tt.expectedAuth {
							t.Errorf("Expected X-API-Key header %q, got %q", tt.expectedAuth, apiKey)
						}
					}
				}

				// Check method
				if r.Method != http.MethodGet {
					t.Errorf("Expected method GET, got %s", r.Method)
				}

				// Set status code
				w.WriteHeader(tt.statusCode)

				// Write response
				if tt.responseBody != nil {
					json.NewEncoder(w).Encode(tt.responseBody)
				}
			})
			defer server.Close()

			// Create client
			var client *Client
			if tt.authType == "apikey" {
				client = NewClientWithConfig(server.URL, "test-api-key")
			} else {
				client = NewClientWithConfig(server.URL, "")
			}

			// Test GET
			var result map[string]interface{}
			err := client.Get("/api/test", &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result == nil {
				t.Error("Expected result to be populated, got nil")
			}
		})
	}
}

func TestClient_Post(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  interface{}
		responseBody interface{}
		statusCode   int
		wantErr      bool
	}{
		{
			name: "successful POST request",
			requestBody: map[string]interface{}{
				"name": "test",
			},
			responseBody: map[string]interface{}{
				"id":   "123",
				"name": "test",
			},
			statusCode: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "POST with 400 error",
			requestBody: map[string]interface{}{
				"invalid": "data",
			},
			responseBody: map[string]interface{}{
				"message": "Bad request",
			},
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
				// Check method
				if r.Method != http.MethodPost {
					t.Errorf("Expected method POST, got %s", r.Method)
				}

				// Check content type
				contentType := r.Header.Get("Content-Type")
				if !strings.Contains(contentType, "application/json") {
					t.Errorf("Expected Content-Type to contain application/json, got %s", contentType)
				}

				// Set status code
				w.WriteHeader(tt.statusCode)

				// Write response
				if tt.responseBody != nil {
					json.NewEncoder(w).Encode(tt.responseBody)
				}
			})
			defer server.Close()

			client := NewClientWithConfig(server.URL, "")
			var result map[string]interface{}
			err := client.Post("/api/test", tt.requestBody, &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_Put(t *testing.T) {
	server := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected method PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": "123"})
	})
	defer server.Close()

	client := NewClientWithConfig(server.URL, "")
	var result map[string]interface{}
	err := client.Put("/api/test", map[string]interface{}{"name": "test"}, &result)

	if err != nil {
		t.Errorf("Client.Put() error = %v", err)
	}
}

func TestClient_Delete(t *testing.T) {
	server := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected method DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	client := NewClientWithConfig(server.URL, "")
	err := client.Delete("/api/test")

	if err != nil {
		t.Errorf("Client.Delete() error = %v", err)
	}
}

func TestClient_UploadFile(t *testing.T) {
	tests := []struct {
		name         string
		fileContent  string
		folderPath   string
		filename     string
		statusCode   int
		wantErr      bool
		checkHeaders bool
	}{
		{
			name:        "successful file upload",
			fileContent: "test file content",
			folderPath:  "/documents",
			filename:    "test.txt",
			statusCode:  http.StatusCreated,
			wantErr:     false,
		},
		{
			name:        "upload without folder path",
			fileContent: "test content",
			statusCode:  http.StatusCreated,
			wantErr:     false,
		},
		{
			name:        "upload with 400 error",
			fileContent: "invalid",
			statusCode:  http.StatusBadRequest,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected method POST, got %s", r.Method)
				}

				// Check content type is multipart
				contentType := r.Header.Get("Content-Type")
				if !strings.HasPrefix(contentType, "multipart/form-data") {
					t.Errorf("Expected Content-Type to be multipart/form-data, got %s", contentType)
				}

				// Parse multipart form
				err := r.ParseMultipartForm(10 << 20) // 10MB
				if err != nil {
					t.Errorf("Failed to parse multipart form: %v", err)
					return
				}

				// Check file field
				file, header, err := r.FormFile("file")
				if err != nil {
					t.Errorf("Failed to get file from form: %v", err)
					return
				}
				defer file.Close()

				if header.Filename == "" {
					t.Error("Expected filename in form, got empty")
				}

				// Check folderPath if provided
				if tt.folderPath != "" {
					folderPath := r.FormValue("folderPath")
					if folderPath != tt.folderPath {
						t.Errorf("Expected folderPath %q, got %q", tt.folderPath, folderPath)
					}
				}

				// Check filename if provided
				if tt.filename != "" {
					filename := r.FormValue("filename")
					if filename != tt.filename {
						t.Errorf("Expected filename %q, got %q", tt.filename, filename)
					}
				}

				w.WriteHeader(tt.statusCode)
				if !tt.wantErr {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"id":       "123",
						"filename": header.Filename,
					})
				} else {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"message": "Upload failed",
					})
				}
			})
			defer server.Close()

			// Create temp file
			tmpFile := t.TempDir() + "/test.txt"
			err := os.WriteFile(tmpFile, []byte(tt.fileContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}

			client := NewClientWithConfig(server.URL, "")
			var result map[string]interface{}
			err = client.UploadFile("/api/files/upload", tmpFile, tt.folderPath, tt.filename, &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UploadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_DownloadFile(t *testing.T) {
	tests := []struct {
		name         string
		fileContent  string
		contentDispo string
		statusCode   int
		wantErr      bool
		outputPath   string
		checkFile    bool
	}{
		{
			name:         "successful file download",
			fileContent:  "downloaded content",
			contentDispo: `attachment; filename="test.txt"`,
			statusCode:   http.StatusOK,
			wantErr:      false,
			checkFile:    true,
		},
		{
			name:        "download to directory",
			fileContent: "test content",
			statusCode:  http.StatusOK,
			wantErr:     false,
			checkFile:   true,
		},
		{
			name:       "download with 404 error",
			statusCode: http.StatusNotFound,
			wantErr:    true,
			checkFile:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected method GET, got %s", r.Method)
				}

				if tt.contentDispo != "" {
					w.Header().Set("Content-Disposition", tt.contentDispo)
				}

				w.WriteHeader(tt.statusCode)
				if !tt.wantErr {
					w.Write([]byte(tt.fileContent))
				} else {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"message": "File not found",
					})
				}
			})
			defer server.Close()

			outputDir := t.TempDir()
			var outputPath string
			if tt.outputPath != "" {
				outputPath = outputDir + "/" + tt.outputPath
			} else {
				outputPath = outputDir
			}

			client := NewClientWithConfig(server.URL, "")
			filePath, err := client.DownloadFile("/api/files/123/download", outputPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DownloadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkFile {
				// Verify file was created
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Errorf("Failed to read downloaded file: %v", err)
					return
				}

				if string(content) != tt.fileContent {
					t.Errorf("Expected file content %q, got %q", tt.fileContent, string(content))
				}
			}
		})
	}
}

func TestClient_ErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
		wantErrMsg string
	}{
		{
			name:       "JSON error response",
			statusCode: http.StatusBadRequest,
			response:   `{"message": "Validation failed", "details": "Invalid input"}`,
			wantErrMsg: "Validation failed",
		},
		{
			name:       "plain text error response",
			statusCode: http.StatusInternalServerError,
			response:   "Internal Server Error",
			wantErrMsg: "Internal Server Error",
		},
		{
			name:       "empty error response",
			statusCode: http.StatusNotFound,
			response:   "",
			wantErrMsg: "Not Found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			})
			defer server.Close()

			client := NewClientWithConfig(server.URL, "")
			var result map[string]interface{}
			err := client.Get("/api/test", &result)

			if err == nil {
				t.Error("Expected error, got nil")
				return
			}

			apiErr, ok := err.(*APIError)
			if !ok {
				t.Errorf("Expected APIError, got %T", err)
				return
			}

			if apiErr.StatusCode != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, apiErr.StatusCode)
			}

			if !strings.Contains(apiErr.Error(), tt.wantErrMsg) {
				t.Errorf("Expected error message to contain %q, got %q", tt.wantErrMsg, apiErr.Error())
			}
		})
	}
}

func TestClient_APIKeyAuth(t *testing.T) {
	// Test that API key is sent in header
	server := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")

		if apiKey != "test-api-key" {
			t.Errorf("Expected X-API-Key header 'test-api-key', got %q", apiKey)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": "123"})
	})
	defer server.Close()

	client := NewClientWithConfig(server.URL, "test-api-key")
	var result map[string]interface{}
	err := client.Get("/api/test", &result)

	if err != nil {
		t.Errorf("Client.Get() error = %v", err)
	}
}

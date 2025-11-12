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
package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// SetupTestServer creates a mock HTTP server for testing
func SetupTestServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

// CreateTestFile creates a temporary test file with the given content
func CreateTestFile(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	return filePath
}

// CreateTestDir creates a temporary directory for testing
func CreateTestDir(t *testing.T) string {
	return t.TempDir()
}

// JSONResponse writes a JSON response to the HTTP response writer
func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	// Encode to buffer first to catch errors before writing headers
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		panic("failed to encode JSON response: " + err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	buf.WriteTo(w)
}

// ErrorResponse writes an error JSON response
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	JSONResponse(w, statusCode, map[string]interface{}{
		"message": message,
	})
}

// ErrorResponseWithDetails writes an error JSON response with details
func ErrorResponseWithDetails(w http.ResponseWriter, statusCode int, message, details string) {
	JSONResponse(w, statusCode, map[string]interface{}{
		"message": message,
		"details": details,
	})
}

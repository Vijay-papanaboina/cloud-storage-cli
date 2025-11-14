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

// SetupErrorServer creates a test server that always returns an error response
func SetupErrorServer(statusCode int, message string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ErrorResponse(w, statusCode, message)
	}))
}

// SetupJSONServer creates a test server that returns a JSON response
func SetupJSONServer(data interface{}, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		JSONResponse(w, statusCode, data)
	}))
}

// SetupAuthServer creates a test server that validates authentication headers
func SetupAuthServer(handler http.HandlerFunc, expectedToken, expectedAPIKey string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate token if expected
		if expectedToken != "" {
			auth := r.Header.Get("Authorization")
			if auth != expectedToken {
				ErrorResponse(w, http.StatusUnauthorized, "Invalid authorization token")
				return
			}
		}

		// Validate API key if expected
		if expectedAPIKey != "" {
			apiKey := r.Header.Get("X-API-Key")
			if apiKey != expectedAPIKey {
				ErrorResponse(w, http.StatusUnauthorized, "Invalid API key")
				return
			}
		}

		// Call the actual handler if auth passes
		handler(w, r)
	}))
}

// CreateTestFileWithContent creates a test file with specific filename and content
func CreateTestFileWithContent(t *testing.T, filename, content string) string {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, filename)
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	return filePath
}

// CreateTestFiles creates multiple test files in a temporary directory
func CreateTestFiles(t *testing.T, files map[string]string) string {
	tmpDir := t.TempDir()
	for filename, content := range files {
		filePath := filepath.Join(tmpDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}
	return tmpDir
}

// AssertFileExists checks if a file exists and fails the test if it doesn't
func AssertFileExists(t *testing.T, path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Expected file to exist: %s", path)
	}
}

// AssertFileContent checks if file content matches expected content
func AssertFileContent(t *testing.T, path, expected string) {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}
	if string(content) != expected {
		t.Errorf("File content mismatch for %s:\nExpected: %q\nGot: %q", path, expected, string(content))
	}
}

// AssertJSONEqual checks if two JSON values are equal
func AssertJSONEqual(t *testing.T, expected, actual interface{}) {
	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("Failed to marshal expected JSON: %v", err)
	}

	actualJSON, err := json.Marshal(actual)
	if err != nil {
		t.Fatalf("Failed to marshal actual JSON: %v", err)
	}

	var expectedNormalized, actualNormalized interface{}
	if err := json.Unmarshal(expectedJSON, &expectedNormalized); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}
	if err := json.Unmarshal(actualJSON, &actualNormalized); err != nil {
		t.Fatalf("Failed to unmarshal actual JSON: %v", err)
	}

	expectedStr, err := json.MarshalIndent(expectedNormalized, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal expected JSON for comparison: %v", err)
	}
	actualStr, err := json.MarshalIndent(actualNormalized, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal actual JSON for comparison: %v", err)
	}

	if string(expectedStr) != string(actualStr) {
		t.Errorf("JSON mismatch:\nExpected:\n%s\nGot:\n%s", string(expectedStr), string(actualStr))
	}
}

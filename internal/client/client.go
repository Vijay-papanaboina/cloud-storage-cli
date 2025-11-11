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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/vijay-papanaboina/cloud-storage-api-cli/internal/config"
)

const (
	defaultTimeout = 30 * time.Second
)

// Client represents an HTTP client for API communication
type Client struct {
	BaseURL     string
	HTTPClient  *http.Client
	AccessToken string
	APIKey      string
}

// NewClient creates a new API client instance
// It loads configuration and initializes the HTTP client
func NewClient() (*Client, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	client := &Client{
		BaseURL:     cfg.APIURL,
		AccessToken: cfg.AccessToken,
		APIKey:      cfg.APIKey,
		HTTPClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	return client, nil
}

// NewClientWithConfig creates a new API client with explicit configuration
// Useful for testing or when config needs to be overridden
func NewClientWithConfig(baseURL, accessToken, apiKey string) *Client {
	return &Client{
		BaseURL:     baseURL,
		AccessToken: accessToken,
		APIKey:      apiKey,
		HTTPClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// buildURL constructs the full URL from base URL and path
func (c *Client) buildURL(path string) (string, error) {
	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	baseURL := strings.TrimSuffix(c.BaseURL, "/")
	fullURL := baseURL + path

	// Validate URL
	_, err := url.Parse(fullURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	return fullURL, nil
}

// setAuthHeaders adds authentication headers to the request
// API key takes precedence if both are set (matches API behavior)
func (c *Client) setAuthHeaders(req *http.Request) {
	if c.APIKey != "" {
		req.Header.Set("X-API-Key", c.APIKey)
	} else if c.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	}
}

// parseErrorResponse parses an error response from the API
func (c *Client) parseErrorResponse(resp *http.Response, method, url string) *APIError {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		apiErr := NewAPIError(resp.StatusCode, "Failed to read error response")
		apiErr.Method = method
		apiErr.URL = url
		return apiErr
	}

	// Try to parse as JSON error response
	var apiErr APIError
	if err := json.Unmarshal(body, &apiErr); err == nil {
		apiErr.StatusCode = resp.StatusCode
		apiErr.Method = method
		apiErr.URL = url
		return &apiErr
	}

	// If not JSON, use response body as message
	message := string(body)
	if message == "" {
		message = resp.Status
	}

	apiErr = *NewAPIError(resp.StatusCode, message)
	apiErr.Method = method
	apiErr.URL = url
	return &apiErr
}

// doRequest performs an HTTP request with the given method, path, and body
func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	fullURL, err := c.buildURL(path)
	if err != nil {
		return nil, err
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add authentication headers
	c.setAuthHeaders(req)

	// Perform request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed [%s %s]: %w", method, fullURL, err)
	}

	// Check for error status codes
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, c.parseErrorResponse(resp, method, fullURL)
	}

	return resp, nil
}

// Get performs a GET request and unmarshals the response into result
func (c *Client) Get(path string, result interface{}) error {
	resp, err := c.doRequest(http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result == nil {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// Post performs a POST request with a JSON body and unmarshals the response into result
func (c *Client) Post(path string, body interface{}, result interface{}) error {
	resp, err := c.doRequest(http.MethodPost, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result == nil {
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(respBody, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// Put performs a PUT request with a JSON body and unmarshals the response into result
func (c *Client) Put(path string, body interface{}, result interface{}) error {
	resp, err := c.doRequest(http.MethodPut, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result == nil {
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(respBody, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// Delete performs a DELETE request
func (c *Client) Delete(path string) error {
	resp, err := c.doRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// UpdateAuth updates the authentication credentials in the client
func (c *Client) UpdateAuth(accessToken, apiKey string) {
	c.AccessToken = accessToken
	c.APIKey = apiKey
}

// UploadFile performs a multipart/form-data file upload request
// path: API endpoint path (e.g., "/api/files/upload")
// filePath: Local file path to upload
// folderPath: Optional folder path (can be empty string)
// result: Pointer to struct to unmarshal JSON response into
func (c *Client) UploadFile(path string, filePath string, folderPath string, result interface{}) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create multipart form data
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add file field
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("failed to create form file field: %w", err)
	}

	// Copy file content to form field
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// Add optional folderPath field
	if folderPath != "" {
		err = writer.WriteField("folderPath", folderPath)
		if err != nil {
			return fmt.Errorf("failed to write folderPath field: %w", err)
		}
	}

	// Close the multipart writer to finalize the form
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Build URL
	fullURL, err := c.buildURL(path)
	if err != nil {
		return err
	}

	// Create request
	req, err := http.NewRequest(http.MethodPost, fullURL, &body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set Content-Type header with boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.ContentLength = int64(body.Len())

	// Add authentication headers
	c.setAuthHeaders(req)

	// Perform request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed [POST %s]: %w", fullURL, err)
	}
	defer resp.Body.Close()

	// Check for error status codes
	if resp.StatusCode >= 400 {
		return c.parseErrorResponse(resp, http.MethodPost, fullURL)
	}

	// Parse response if result is provided
	if result != nil {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// extractFilenameFromContentDisposition extracts filename from Content-Disposition header
// Handles formats like: attachment; filename="filename.ext" or attachment; filename=filename.ext
func extractFilenameFromContentDisposition(header string) string {
	if header == "" {
		return ""
	}

	// Try to match filename="..." or filename=...
	// Pattern: filename="..." or filename=...
	re := regexp.MustCompile(`filename[*]?=(?:"([^"]+)"|([^;]+))`)
	matches := re.FindStringSubmatch(header)
	if len(matches) > 0 {
		// matches[1] is for quoted filename, matches[2] is for unquoted
		if matches[1] != "" {
			return strings.TrimSpace(matches[1])
		}
		if matches[2] != "" {
			return strings.TrimSpace(matches[2])
		}
	}

	return ""
}

// sanitizeFilename sanitizes a filename to prevent path traversal
func sanitizeFilename(filename string) string {
	// Remove path separators and other dangerous characters
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.ReplaceAll(filename, "..", "_")

	// Remove any remaining path components
	filename = filepath.Base(filename)

	// If empty after sanitization, use a default
	if filename == "" || filename == "." || filename == ".." {
		return "download"
	}

	return filename
}

// DownloadFile downloads a file from the API and saves it to the specified output path
// path: API endpoint path (e.g., "/api/files/{id}/download")
// outputPath: Local file path to save the downloaded file (can be directory or full path)
// Returns the final file path where the file was saved
func (c *Client) DownloadFile(path string, outputPath string) (string, error) {
	// Build URL
	fullURL, err := c.buildURL(path)
	if err != nil {
		return "", err
	}

	// Create request
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "*/*")

	// Add authentication headers
	c.setAuthHeaders(req)

	// Perform request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed [GET %s]: %w", fullURL, err)
	}
	defer resp.Body.Close()

	// Check for error status codes
	if resp.StatusCode >= 400 {
		return "", c.parseErrorResponse(resp, http.MethodGet, fullURL)
	}

	// Extract filename from Content-Disposition header
	contentDisposition := resp.Header.Get("Content-Disposition")
	filename := extractFilenameFromContentDisposition(contentDisposition)

	// Sanitize filename
	if filename != "" {
		filename = sanitizeFilename(filename)
	} else {
		// If no filename in header, extract from path or use default
		filename = "download"
		// Try to extract file ID from path as fallback
		parts := strings.Split(path, "/")
		if len(parts) > 0 {
			lastPart := parts[len(parts)-1]
			if lastPart != "download" && lastPart != "" {
				filename = lastPart
			}
		}
	}

	// Determine final output path
	var finalPath string
	outputPathInfo, err := os.Stat(outputPath)
	if err == nil && outputPathInfo.IsDir() {
		// Output path is a directory, combine with filename
		finalPath = filepath.Join(outputPath, filename)
	} else if outputPath != "" {
		// Output path is specified and not a directory, use it as-is
		finalPath = outputPath
		// Create parent directory if it doesn't exist
		parentDir := filepath.Dir(finalPath)
		if parentDir != "." && parentDir != "" {
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				return "", fmt.Errorf("failed to create output directory: %w", err)
			}
		}
	} else {
		// No output path specified, use current directory with filename
		finalPath = filename
	}

	// Create output file
	outFile, err := os.Create(finalPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Stream response body to file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		// Clean up file on error
		os.Remove(finalPath)
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return finalPath, nil
}

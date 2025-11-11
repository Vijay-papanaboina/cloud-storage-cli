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
	"net/http"
	"net/url"
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
func (c *Client) parseErrorResponse(resp *http.Response) *APIError {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewAPIError(resp.StatusCode, "Failed to read error response")
	}

	// Try to parse as JSON error response
	var apiErr APIError
	if err := json.Unmarshal(body, &apiErr); err == nil {
		apiErr.StatusCode = resp.StatusCode
		return &apiErr
	}

	// If not JSON, use response body as message
	message := string(body)
	if message == "" {
		message = resp.Status
	}

	return NewAPIError(resp.StatusCode, message)
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
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Check for error status codes
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, c.parseErrorResponse(resp)
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


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

import "fmt"

// APIError represents an error response from the API
type APIError struct {
	StatusCode int    `json:"statusCode,omitempty"`
	Message    string `json:"message,omitempty"`
	Details    string `json:"details,omitempty"`
	Method     string `json:"method,omitempty"` // HTTP method
	URL        string `json:"url,omitempty"`    // Request URL
}

// Error implements the error interface
func (e *APIError) Error() string {
	baseMsg := fmt.Sprintf("API error (%d)", e.StatusCode)
	if e.Method != "" && e.URL != "" {
		baseMsg = fmt.Sprintf("API error (%d) [%s %s]", e.StatusCode, e.Method, e.URL)
	}
	if e.Details != "" {
		return fmt.Sprintf("%s: %s - %s", baseMsg, e.Message, e.Details)
	}
	if e.Message != "" {
		if e.Message != "" {
			return fmt.Sprintf("%s: %s - %s", baseMsg, e.Message, e.Details)
		}
		return fmt.Sprintf("%s: %s", baseMsg, e.Details)
	}
	return baseMsg
}

// NewAPIError creates a new APIError instance
func NewAPIError(statusCode int, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// NewAPIErrorWithDetails creates a new APIError instance with details
func NewAPIErrorWithDetails(statusCode int, message, details string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
	}
}

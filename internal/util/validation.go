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
package util

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var (
	// uuidRegex matches standard UUID format (8-4-4-4-12 hex digits)
	uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
)

// ValidateUUID validates that a string is a valid UUID format
func ValidateUUID(id string) error {
	if id == "" {
		return fmt.Errorf("UUID cannot be empty")
	}
	if !uuidRegex.MatchString(id) {
		return fmt.Errorf("invalid UUID format: %s (expected format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)", id)
	}
	// Additional validation using google/uuid package
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}
	return nil
}

// ValidatePath validates a folder/file path
// Paths must start with '/' and use Unix-style separators
func ValidatePath(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("path must start with '/'")
	}
	// Check for path traversal attempts
	if strings.Contains(path, "..") {
		return fmt.Errorf("path cannot contain '..'")
	}
	// Check for backslashes (Windows-style paths not allowed)
	if strings.Contains(path, "\\") {
		return fmt.Errorf("path must use forward slashes, not backslashes")
	}
	// Check for null bytes
	if strings.Contains(path, "\x00") {
		return fmt.Errorf("path cannot contain null bytes")
	}
	// Check for control characters (including tab, newline, carriage return)
	// Only allow printable characters (code >= 32)
	for _, r := range path {
		if r < 32 {
			return fmt.Errorf("path cannot contain control characters")
		}
	}
	return nil
}

// ValidateFilename validates a filename
func ValidateFilename(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	// Platform-independent check: explicitly reject backslashes before normalization
	if strings.Contains(filename, "\\") {
		return fmt.Errorf("filename cannot contain path separators")
	}
	// Get base name to prevent path traversal
	baseName := filepath.Base(filename)
	if baseName != filename {
		return fmt.Errorf("filename cannot contain path separators")
	}
	if filename == "." || filename == ".." {
		return fmt.Errorf("filename cannot be '.' or '..'")
	}
	// Check for null bytes
	if strings.Contains(filename, "\x00") {
		return fmt.Errorf("filename cannot contain null bytes")
	}
	// Check for control characters (including tab, newline, carriage return)
	// Only allow printable characters (code >= 32)
	for _, r := range filename {
		if r < 32 {
			return fmt.Errorf("filename cannot contain control characters")
		}
	}
	// Windows reserved names
	reservedNames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}
	upperName := strings.ToUpper(baseName)
	for _, reserved := range reservedNames {
		if upperName == reserved || strings.HasPrefix(upperName, reserved+".") {
			return fmt.Errorf("filename cannot be a reserved name: %s", reserved)
		}
	}
	return nil
}

// ValidateEmail validates an email address format (basic validation)
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format: %s", email)
	}
	return nil
}

// ValidateUsername validates a username
func ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}
	if len(username) > 50 {
		return fmt.Errorf("username must be at most 50 characters")
	}
	// Allow alphanumeric, underscore, hyphen, dot
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	if !usernameRegex.MatchString(username) {
		return fmt.Errorf("username can only contain letters, numbers, underscores, hyphens, and dots")
	}
	return nil
}

// ValidatePageSize validates pagination page size
func ValidatePageSize(size int) error {
	if size <= 0 {
		return fmt.Errorf("page size must be greater than 0")
	}
	if size > 100 {
		return fmt.Errorf("page size must be at most 100")
	}
	return nil
}

// ValidatePageNumber validates pagination page number
func ValidatePageNumber(page int) error {
	if page < 0 {
		return fmt.Errorf("page number must be >= 0")
	}
	return nil
}

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
	"strings"
	"testing"
)

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"valid lowercase", "550e8400-e29b-41d4-a716-446655440000", false, ""},
		{"valid uppercase", "550E8400-E29B-41D4-A716-446655440000", false, ""},
		{"valid mixed case", "550e8400-E29b-41d4-A716-446655440000", false, ""},
		{"empty", "", true, "UUID cannot be empty"},
		{"missing hyphens", "550e8400e29b41d4a716446655440000", true, "invalid UUID format"},
		{"wrong separator", "550e8400_e29b_41d4_a716_446655440000", true, "invalid UUID format"},
		{"too short", "550e8400-e29b-41d4-a716-44665544000", true, "invalid UUID format"},
		{"invalid char", "550e8400-e29b-41d4-a716-44665544000g", true, "invalid UUID"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUUID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && (err == nil || !strings.Contains(err.Error(), tt.errMsg)) {
				t.Errorf("ValidateUUID() error = %v, want error containing %q", err, tt.errMsg)
			}
		})
	}
}

func TestValidatePath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"valid root", "/", false, ""},
		{"valid simple", "/documents", false, ""},
		{"valid nested", "/documents/photos/2024", false, ""},
		{"empty", "", true, "path cannot be empty"},
		{"no leading slash", "documents", true, "path must start with '/'"},
		{"path traversal", "/documents/../etc", true, "path cannot contain '..'"},
		{"backslash", "/documents\\photos", true, "path must use forward slashes"},
		{"null byte", "/documents\x00/photos", true, "path cannot contain null bytes"},
		{"tab char", "/documents\t/photos", true, "path cannot contain control characters"},
		{"newline", "/documents\n/photos", true, "path cannot contain control characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePath(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && (err == nil || !strings.Contains(err.Error(), tt.errMsg)) {
				t.Errorf("ValidatePath() error = %v, want error containing %q", err, tt.errMsg)
			}
		})
	}
}

func TestValidateFilename(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"valid simple", "document.pdf", false, ""},
		{"valid with underscore", "my_file.txt", false, ""},
		{"valid with hyphen", "my-file.txt", false, ""},
		{"empty", "", true, "filename cannot be empty"},
		{"with slash", "documents/file.txt", true, "filename cannot contain path separators"},
		{"current dir", ".", true, "filename cannot be '.' or '..'"},
		{"parent dir", "..", true, "filename cannot be '.' or '..'"},
		{"null byte", "file\x00.txt", true, "filename cannot contain null bytes"},
		{"tab char", "file\t.txt", true, "filename cannot contain control characters"},
		{"reserved CON", "CON", true, "filename cannot be a reserved name: CON"},
		{"reserved CON lowercase", "con", true, "filename cannot be a reserved name: CON"},
		{"reserved CON.ext", "CON.txt", true, "filename cannot be a reserved name: CON"},
		{"reserved PRN", "PRN", true, "filename cannot be a reserved name: PRN"},
		{"reserved COM1", "COM1", true, "filename cannot be a reserved name: COM1"},
		{"reserved LPT1", "LPT1", true, "filename cannot be a reserved name: LPT1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFilename(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFilename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && (err == nil || !strings.Contains(err.Error(), tt.errMsg)) {
				t.Errorf("ValidateFilename() error = %v, want error containing %q", err, tt.errMsg)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"valid simple", "user@example.com", false, ""},
		{"valid subdomain", "user@mail.example.com", false, ""},
		{"valid with plus", "user+tag@example.com", false, ""},
		{"valid with dot", "user.name@example.com", false, ""},
		{"empty", "", true, "email cannot be empty"},
		{"missing @", "userexample.com", true, "invalid email format"},
		{"missing domain", "user@", true, "invalid email format"},
		{"missing local", "@example.com", true, "invalid email format"},
		{"missing TLD", "user@example", true, "invalid email format"},
		{"TLD too short", "user@example.c", true, "invalid email format"},
		{"with space", "user name@example.com", true, "invalid email format"},
		{"multiple @", "user@name@example.com", true, "invalid email format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && (err == nil || !strings.Contains(err.Error(), tt.errMsg)) {
				t.Errorf("ValidateEmail() error = %v, want error containing %q", err, tt.errMsg)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"valid min length", "abc", false, ""},
		{"valid with numbers", "user123", false, ""},
		{"valid with underscore", "user_name", false, ""},
		{"valid with hyphen", "user-name", false, ""},
		{"valid with dot", "user.name", false, ""},
		{"valid max length", strings.Repeat("a", 50), false, ""},
		{"empty", "", true, "username cannot be empty"},
		{"too short", "ab", true, "username must be at least 3 characters"},
		{"too long", strings.Repeat("a", 51), true, "username must be at most 50 characters"},
		{"with space", "user name", true, "username can only contain letters"},
		{"with @", "user@name", true, "username can only contain letters"},
		{"with #", "user#name", true, "username can only contain letters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && (err == nil || !strings.Contains(err.Error(), tt.errMsg)) {
				t.Errorf("ValidateUsername() error = %v, want error containing %q", err, tt.errMsg)
			}
		})
	}
}

func TestValidatePageSize(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantErr bool
		errMsg  string
	}{
		{"valid min", 1, false, ""},
		{"valid middle", 50, false, ""},
		{"valid max", 100, false, ""},
		{"zero", 0, true, "page size must be greater than 0"},
		{"negative", -1, true, "page size must be greater than 0"},
		{"too large", 101, true, "page size must be at most 100"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePageSize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePageSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && (err == nil || !strings.Contains(err.Error(), tt.errMsg)) {
				t.Errorf("ValidatePageSize() error = %v, want error containing %q", err, tt.errMsg)
			}
		})
	}
}

func TestValidatePageNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantErr bool
		errMsg  string
	}{
		{"valid zero", 0, false, ""},
		{"valid positive", 1, false, ""},
		{"valid large", 100, false, ""},
		{"negative", -1, true, "page number must be >= 0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePageNumber(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePageNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && (err == nil || !strings.Contains(err.Error(), tt.errMsg)) {
				t.Errorf("ValidatePageNumber() error = %v, want error containing %q", err, tt.errMsg)
			}
		})
	}
}



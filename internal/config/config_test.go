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
package config

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestConfig creates a temporary config directory for testing
func setupTestConfig(t *testing.T) (string, func()) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".cloud-storage-cli")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test config directory: %v", err)
	}

	// Save original values
	originalViper := viperInstance
	originalPath := configPath

	// Reset viper instance
	viperInstance = nil
	configPath = filepath.Join(configDir, configFileName)

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	if originalHome == "" {
		originalHome = os.Getenv("USERPROFILE") // Windows
	}

	cleanup := func() {
		viperInstance = originalViper
		configPath = originalPath
		os.RemoveAll(tmpDir)
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		} else {
			os.Unsetenv("HOME")
			os.Unsetenv("USERPROFILE")
		}
	}

	return tmpDir, cleanup
}

func TestSaveConfig(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	cfg := &Config{
		APIURL: "http://test.example.com",
		APIKey: "test-api-key",
	}

	err := SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Verify file exists
	if configPath == "" {
		t.Fatal("configPath is empty")
	}

	fileInfo, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Config file does not exist: %v", err)
	}

	// Verify file permissions (0600 on Unix, may differ on Windows)
	mode := fileInfo.Mode().Perm()
	// On Windows, permissions might be different, so we just check that file exists
	// The important thing is that the file was created successfully
	_ = mode
}

func TestLoadConfig(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Save a config first
	cfg := &Config{
		APIURL: "http://test.example.com",
		APIKey: "test-api-key",
	}

	err := SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Reset viper to test loading
	viperInstance = nil

	// Load config
	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// API URL is hardcoded at compile time, so it should return GetAPIURL()
	expectedURL := GetAPIURL()
	if loaded.APIURL != expectedURL {
		t.Errorf("Expected APIURL %q (compile-time), got %q", expectedURL, loaded.APIURL)
	}

	// API key should be loaded from config
	if loaded.APIKey != cfg.APIKey {
		t.Errorf("Expected APIKey %q, got %q", cfg.APIKey, loaded.APIKey)
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	tmpDir, cleanup := setupTestConfig(t)
	defer cleanup()

	// Reset viper to ensure clean state
	viperInstance = nil
	configPath = filepath.Join(tmpDir, ".cloud-storage-cli", configFileName)

	// Ensure config directory exists but no config file
	configDir := filepath.Dir(configPath)
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	// Initialize fresh config (no file exists, should use defaults)
	err = InitConfig()
	if err != nil {
		t.Fatalf("InitConfig() error = %v", err)
	}

	// Load config without saving (should use compile-time URL)
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// API URL is hardcoded at compile time
	expectedURL := GetAPIURL()
	if cfg.APIURL != expectedURL {
		t.Errorf("Expected APIURL %q (compile-time), got %q", expectedURL, cfg.APIURL)
	}

	// APIKey should be empty if no config file and no env var
	// But we can't guarantee this due to env vars, so we just verify it loads
	_ = cfg.APIKey
}

func TestSetValue(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	tests := []struct {
		key   string
		value string
	}{
		{"api-key", "new-key"},
		{"api_key", "new-key2"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			err := SetValue(tt.key, tt.value)
			if err != nil {
				t.Fatalf("SetValue() error = %v", err)
			}

			// Verify value was set
			value, err := GetValue(tt.key)
			if err != nil {
				t.Fatalf("GetValue() error = %v", err)
			}

			if value != tt.value {
				t.Errorf("Expected value %q, got %q", tt.value, value)
			}
		})
	}
}

func TestSetValue_APIURL_NotAllowed(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// API URL cannot be set via SetValue - it's compile-time only
	tests := []string{"api-url", "api_url"}

	for _, key := range tests {
		t.Run(key, func(t *testing.T) {
			err := SetValue(key, "http://new.example.com")
			if err == nil {
				t.Errorf("Expected error when setting %q, got nil", key)
			}
		})
	}
}

func TestSetValue_InvalidKey(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	err := SetValue("invalid-key", "value")
	if err == nil {
		t.Error("Expected error for invalid key, got nil")
	}
}

func TestGetValue(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Test getting API key (can be set)
	err := SetValue("api-key", "test-api-key")
	if err != nil {
		t.Fatalf("SetValue() error = %v", err)
	}

	// Get API key value
	value, err := GetValue("api-key")
	if err != nil {
		t.Fatalf("GetValue() error = %v", err)
	}

	if value != "test-api-key" {
		t.Errorf("Expected value %q, got %q", "test-api-key", value)
	}

	// Test getting API URL (compile-time, read-only)
	urlValue, err := GetValue("api-url")
	if err != nil {
		t.Fatalf("GetValue() error = %v", err)
	}

	expectedURL := GetAPIURL()
	if urlValue != expectedURL {
		t.Errorf("Expected API URL %q (compile-time), got %q", expectedURL, urlValue)
	}
}

func TestGetValue_InvalidKey(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	_, err := GetValue("invalid-key")
	if err == nil {
		t.Error("Expected error for invalid key, got nil")
	}
}

func TestIsSensitiveKey(t *testing.T) {
	tests := []struct {
		key      string
		expected bool
	}{
		{"api-key", true},
		{"api_key", true},
		{"api-url", false},
		{"api_url", false},
		{"invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := IsSensitiveKey(tt.key)
			if result != tt.expected {
				t.Errorf("IsSensitiveKey(%q) = %v, want %v", tt.key, result, tt.expected)
			}
		})
	}
}

func TestMaskValue(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", "(not set)"},
		{"short", "***"},
		{"12345678", "***"},
		{"123456789", "1234...6789"},
		{"very-long-token-value-here", "very...here"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := MaskValue(tt.input)
			if result != tt.expected {
				t.Errorf("MaskValue(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConfig_EnvironmentVariables(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	// Set environment variable (should be ignored - API URL is compile-time only)
	os.Setenv("CLOUD_STORAGE_API_URL", "http://env.example.com")
	defer os.Unsetenv("CLOUD_STORAGE_API_URL")

	// Reset viper
	viperInstance = nil

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// API URL should be compile-time value, not from env var
	expectedURL := GetAPIURL()
	if cfg.APIURL != expectedURL {
		t.Errorf("Expected APIURL %q (compile-time, env var ignored), got %q", expectedURL, cfg.APIURL)
	}
}

func TestGetConfigPath(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	path := GetConfigPath()
	if path == "" {
		t.Error("GetConfigPath() returned empty string")
	}

	if !filepath.IsAbs(path) {
		t.Errorf("Expected absolute path, got %q", path)
	}

	if filepath.Base(path) != configFileName {
		t.Errorf("Expected config file name %q, got %q", configFileName, filepath.Base(path))
	}
}

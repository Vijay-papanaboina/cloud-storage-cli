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

// InvalidPaths returns a slice of invalid path examples for testing
func InvalidPaths() []string {
	return []string{
		"",                    // Empty path
		"documents",           // No leading slash
		"../etc/passwd",       // Path traversal
		"/documents/../etc",   // Path traversal
		"/documents/..",       // Path traversal
		"/documents\\photos",  // Backslash
		"/documents\x00/photos", // Null byte
		"/documents\t/photos",   // Tab character
		"/documents\n/photos",   // Newline
		"/documents//photos",    // Consecutive slashes (if not allowed)
	}
}

// InvalidFilenames returns a slice of invalid filename examples for testing
func InvalidFilenames() []string {
	return []string{
		"",                    // Empty filename
		"documents/file.txt",  // Path separator
		"folder\\test.txt",    // Backslash
		".",                   // Current directory
		"..",                  // Parent directory
		"file\x00.txt",        // Null byte
		"file\t.txt",          // Tab character
		"file\n.txt",          // Newline
		"CON",                 // Reserved name
		"con",                 // Reserved name (lowercase)
		"CON.txt",             // Reserved name with extension
		"PRN",                 // Reserved name
		"COM1",                // Reserved name
		"LPT1",                // Reserved name
	}
}

// InvalidUUIDs returns a slice of invalid UUID examples for testing
func InvalidUUIDs() []string {
	return []string{
		"",                                          // Empty UUID
		"550e8400e29b41d4a716446655440000",          // Missing hyphens
		"550e8400_e29b_41d4_a716_446655440000",      // Wrong separator
		"550e8400-e29b-41d4-a716-44665544000",       // Too short
		"550e8400-e29b-41d4-a716-44665544000g",      // Invalid character
		"550e8400-e29b-41d4-a716",                   // Incomplete
		"not-a-uuid",                                // Random string
		"12345",                                     // Too short
	}
}

// InvalidEmails returns a slice of invalid email examples for testing
func InvalidEmails() []string {
	return []string{
		"",                      // Empty email
		"userexample.com",       // Missing @
		"user@",                 // Missing domain
		"@example.com",          // Missing local part
		"user@example",          // Missing TLD
		"user@example.c",        // TLD too short
		"user name@example.com", // Space in email
		"user@name@example.com", // Multiple @
		"user@@example.com",     // Double @
	}
}

// ReservedFilenames returns a slice of Windows reserved filenames
func ReservedFilenames() []string {
	return []string{
		"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	}
}

// ValidPaths returns a slice of valid path examples for testing
func ValidPaths() []string {
	return []string{
		"/",
		"/documents",
		"/documents/photos",
		"/documents/photos/2024",
		"/a",
		"/very/long/path/with/many/segments",
	}
}

// ValidFilenames returns a slice of valid filename examples for testing
func ValidFilenames() []string {
	return []string{
		"document.pdf",
		"my_file.txt",
		"my-file.txt",
		"file123.txt",
		"test.file.name.txt",
		"a",
		"very-long-filename-with-many-characters.txt",
	}
}

// ValidUUIDs returns a slice of valid UUID examples for testing
func ValidUUIDs() []string {
	return []string{
		"550e8400-e29b-41d4-a716-446655440000", // Lowercase
		"550E8400-E29B-41D4-A716-446655440000", // Uppercase
		"550e8400-E29b-41d4-A716-446655440000", // Mixed case
		"00000000-0000-0000-0000-000000000000", // All zeros
		"ffffffff-ffff-ffff-ffff-ffffffffffff", // All Fs
	}
}

// ValidEmails returns a slice of valid email examples for testing
func ValidEmails() []string {
	return []string{
		"user@example.com",
		"user@mail.example.com",
		"user+tag@example.com",
		"user.name@example.com",
		"user_name@example.com",
		"user123@example.com",
		"test@subdomain.example.com",
	}
}


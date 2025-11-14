//go:build security || !unit

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

	"pgregory.net/rapid"
)

// TestValidatePath_PropertyBased tests that paths containing ".." always fail validation
func TestValidatePath_PropertyBased(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		path := rapid.String().Draw(t, "path")
		err := ValidatePath(path)

		// Property: If path contains "..", it should always error
		if strings.Contains(path, "..") && err == nil {
			t.Fatalf("Path with '..' should always error: %q", path)
		}

		// Property: If path contains backslash, it should always error
		if strings.Contains(path, "\\") && err == nil {
			t.Fatalf("Path with backslash should always error: %q", path)
		}

		// Property: If path contains null byte, it should always error
		if strings.Contains(path, "\x00") && err == nil {
			t.Fatalf("Path with null byte should always error: %q", path)
		}

		// Property: If path doesn't start with "/" and is not empty, it should error
		if path != "" && !strings.HasPrefix(path, "/") && err == nil {
			t.Fatalf("Path without leading slash should error: %q", path)
		}
	})
}

// TestValidateFilename_PropertyBased tests that filenames with path separators always fail
func TestValidateFilename_PropertyBased(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		filename := rapid.String().Draw(t, "filename")
		err := ValidateFilename(filename)

		// Property: If filename contains "/" or "\", it should always error
		// (Validation function now explicitly checks for backslashes before normalization)
		if (strings.Contains(filename, "/") || strings.Contains(filename, "\\")) && err == nil {
			t.Fatalf("Filename with path separator should always error: %q", filename)
		}

		// Property: If filename is "." or "..", it should always error
		if (filename == "." || filename == "..") && err == nil {
			t.Fatalf("Filename '.' or '..' should always error: %q", filename)
		}

		// Property: If filename contains null byte, it should always error
		if strings.Contains(filename, "\x00") && err == nil {
			t.Fatalf("Filename with null byte should always error: %q", filename)
		}

		// Property: Reserved names (case-insensitive) should always error
		reservedNames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
			"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}
		upperFilename := strings.ToUpper(filename)
		for _, reserved := range reservedNames {
			if upperFilename == reserved || strings.HasPrefix(upperFilename, reserved+".") {
				if err == nil {
					t.Fatalf("Reserved filename should always error: %q", filename)
				}
				break
			}
		}
	})
}

// TestValidateFilename_PropertyBased_Valid tests that valid filenames always pass validation
func TestValidateFilename_PropertyBased_Valid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate valid filename characters: letters, digits, '-', '_', '.'
		// Must not start with '.' and must not be reserved names
		validRuneGen := rapid.Custom(func(t *rapid.T) rune {
			choice := rapid.IntRange(0, 5).Draw(t, "choice")
			switch choice {
			case 0:
				return rune(rapid.Int32Range('a', 'z').Draw(t, "lower"))
			case 1:
				return rune(rapid.Int32Range('A', 'Z').Draw(t, "upper"))
			case 2:
				return rune(rapid.Int32Range('0', '9').Draw(t, "digit"))
			case 3:
				return '-'
			case 4:
				return '_'
			case 5:
				return '.'
			default:
				return 'a'
			}
		})

		// Generate first character (must not be '.')
		firstCharGen := rapid.Custom(func(t *rapid.T) rune {
			choice := rapid.IntRange(0, 4).Draw(t, "firstChoice")
			switch choice {
			case 0:
				return rune(rapid.Int32Range('a', 'z').Draw(t, "firstLower"))
			case 1:
				return rune(rapid.Int32Range('A', 'Z').Draw(t, "firstUpper"))
			case 2:
				return rune(rapid.Int32Range('0', '9').Draw(t, "firstDigit"))
			case 3:
				return '-'
			case 4:
				return '_'
			default:
				return 'a'
			}
		})

		firstChar := firstCharGen.Draw(t, "firstChar")

		// Generate remaining characters
		remaining := rapid.SliceOfN(validRuneGen, 0, 50).Draw(t, "remaining")

		// Build filename
		var builder strings.Builder
		builder.WriteRune(firstChar)
		for _, r := range remaining {
			builder.WriteRune(r)
		}
		filename := builder.String()

		// Filter out reserved names (case-insensitive)
		reservedNames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
			"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}
		upperName := strings.ToUpper(filename)
		for _, reserved := range reservedNames {
			if upperName == reserved || strings.HasPrefix(upperName, reserved+".") {
				t.Skipf("Skipping reserved name: %q", filename)
				return
			}
		}

		err := ValidateFilename(filename)
		if err != nil {
			t.Fatalf("Valid filename should pass validation: %q, error: %v", filename, err)
		}
	})
}

// TestValidateUUID_PropertyBased tests UUID validation properties
func TestValidateUUID_PropertyBased(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		uuidStr := rapid.String().Draw(t, "uuid")
		err := ValidateUUID(uuidStr)

		// Property: Empty string should always error
		if uuidStr == "" && err == nil {
			t.Fatalf("Empty UUID should always error")
		}

		// Property: UUID without hyphens should error
		if uuidStr != "" && !strings.Contains(uuidStr, "-") && err == nil {
			t.Fatalf("UUID without hyphens should error: %q", uuidStr)
		}

		// Property: UUID with wrong length should error
		// Valid UUID format: 8-4-4-4-12 = 36 characters total
		if uuidStr != "" && len(uuidStr) != 36 && err == nil {
			t.Fatalf("UUID with wrong length should error: %q (length: %d)", uuidStr, len(uuidStr))
		}
	})
}

// TestValidateEmail_PropertyBased tests email validation properties
func TestValidateEmail_PropertyBased(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		email := rapid.String().Draw(t, "email")
		err := ValidateEmail(email)

		// Property: Empty string should always error
		if email == "" && err == nil {
			t.Fatalf("Empty email should always error")
		}

		// Property: Email without "@" should error
		if email != "" && !strings.Contains(email, "@") && err == nil {
			t.Fatalf("Email without '@' should error: %q", email)
		}

		// Property: Email with multiple "@" should error
		atCount := strings.Count(email, "@")
		if atCount > 1 && err == nil {
			t.Fatalf("Email with multiple '@' should error: %q", email)
		}

		// Property: Email with spaces should error
		if strings.Contains(email, " ") && err == nil {
			t.Fatalf("Email with spaces should error: %q", email)
		}
	})
}

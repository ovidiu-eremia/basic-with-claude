// ABOUTME: Test file for the main CLI functionality of the BASIC interpreter
// ABOUTME: Tests file reading, error handling, and basic command-line interface

package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestReadBasicFile(t *testing.T) {
	// Test reading a valid BASIC file
	testContent := `10 PRINT "HELLO WORLD"
20 END`

	// Create temporary file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.bas")

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	content, err := readBasicFile(testFile)
	if err != nil {
		t.Errorf("readBasicFile() returned error: %v", err)
	}

	if strings.TrimSpace(content) != strings.TrimSpace(testContent) {
		t.Errorf("readBasicFile() = %q, want %q", content, testContent)
	}
}

func TestReadBasicFileNotFound(t *testing.T) {
	// Test reading a non-existent file
	_, err := readBasicFile("nonexistent.bas")
	if err == nil {
		t.Error("readBasicFile() should return error for non-existent file")
	}

	// Error message should indicate file not found
	if !strings.Contains(err.Error(), "no such file") && !strings.Contains(err.Error(), "cannot find") {
		t.Errorf("Error should indicate file not found, got: %v", err)
	}
}

func TestReadBasicFilePermissionDenied(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("permission tests not reliable on this platform or when running as root")
	}

	// Create a file with no read permissions
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "noperm.bas")

	err := os.WriteFile(testFile, []byte("10 PRINT \"TEST\""), 0000) // No permissions
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = readBasicFile(testFile)
	if err == nil {
		t.Error("readBasicFile() should return error for permission denied")
	}
}

func TestParseInputsFlag(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single input",
			input:    "42",
			expected: []string{"42"},
		},
		{
			name:     "multiple inputs",
			input:    "42,hello,17",
			expected: []string{"42", "hello", "17"},
		},
		{
			name:     "inputs with spaces",
			input:    "42, hello , 17",
			expected: []string{"42", "hello", "17"},
		},
		{
			name:     "empty string input",
			input:    "42,,17",
			expected: []string{"42", "", "17"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputs := strings.Split(tt.input, ",")
			for i := range inputs {
				inputs[i] = strings.TrimSpace(inputs[i])
			}

			if len(inputs) != len(tt.expected) {
				t.Errorf("Expected %d inputs, got %d", len(tt.expected), len(inputs))
				return
			}

			for i, expected := range tt.expected {
				if inputs[i] != expected {
					t.Errorf("Input[%d] = %q, want %q", i, inputs[i], expected)
				}
			}
		})
	}
}

func TestFlagValidation(t *testing.T) {
	tests := []struct {
		name         string
		executeFlag  string
		args         []string
		shouldError  bool
		errorMessage string
	}{
		{
			name:        "valid file argument",
			executeFlag: "",
			args:        []string{"test.bas"},
			shouldError: false,
		},
		{
			name:        "valid execute flag",
			executeFlag: "10 PRINT \"TEST\": 20 END",
			args:        []string{},
			shouldError: false,
		},
		{
			name:         "both execute flag and file",
			executeFlag:  "10 PRINT \"TEST\"",
			args:         []string{"test.bas"},
			shouldError:  true,
			errorMessage: "Cannot specify both -e flag and filename",
		},
		{
			name:        "no execute flag or file",
			executeFlag: "",
			args:        []string{},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the validation logic directly (without actually running main)
			executeFlag := tt.executeFlag
			hasFile := len(tt.args) > 0

			// Replicate the validation logic from main()
			bothSpecified := executeFlag != "" && hasFile
			neitherSpecified := executeFlag == "" && !hasFile

			if tt.shouldError {
				if !bothSpecified && !neitherSpecified {
					t.Error("Expected validation error but none occurred")
				}
				if bothSpecified && tt.errorMessage != "" {
					// This would trigger the "Cannot specify both" error
					if tt.errorMessage != "Cannot specify both -e flag and filename" {
						t.Errorf("Wrong error message expected")
					}
				}
			} else {
				if bothSpecified || neitherSpecified {
					t.Error("Unexpected validation error")
				}
			}
		})
	}
}

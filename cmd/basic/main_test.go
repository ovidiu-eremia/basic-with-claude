// ABOUTME: Test file for the main CLI functionality of the BASIC interpreter
// ABOUTME: Tests file reading, error handling, and basic command-line interface

package main

import (
	"os"
	"path/filepath"
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


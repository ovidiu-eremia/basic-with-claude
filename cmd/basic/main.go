// ABOUTME: Main CLI entry point for the BASIC interpreter
// ABOUTME: Handles command-line arguments and file loading with proper error handling

package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename.bas>\n", os.Args[0])
		os.Exit(1)
	}
	
	filename := os.Args[1]
	
	content, err := readBasicFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}
	
	// For now, just print success message
	// Later steps will parse and execute the content
	fmt.Println(formatSuccessMessage(filename))
	
	// Store content for future use (will be used in later steps)
	_ = content
}

// readBasicFile reads the contents of a BASIC program file
func readBasicFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// formatSuccessMessage creates a standard success message for file loading
func formatSuccessMessage(filename string) string {
	return fmt.Sprintf("Program loaded: %s", filename)
}
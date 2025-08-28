// ABOUTME: Standard runtime implementation for console I/O operations
// ABOUTME: Production runtime that uses os.Stdout and os.Stdin for real console interaction

package runtime

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// StandardRuntime implements Runtime interface for console I/O
type StandardRuntime struct {
	reader *bufio.Reader
}

// NewStandardRuntime creates a new StandardRuntime instance
func NewStandardRuntime() *StandardRuntime {
	return &StandardRuntime{
		reader: bufio.NewReader(os.Stdin),
	}
}

// Print outputs a string to stdout without a newline
func (sr *StandardRuntime) Print(value string) error {
	_, err := fmt.Print(value)
	return err
}

// PrintLine outputs a string to stdout with a newline
func (sr *StandardRuntime) PrintLine(value string) error {
	_, err := fmt.Println(value)
	return err
}

// Input prompts for user input and returns the entered string
func (sr *StandardRuntime) Input(prompt string) (string, error) {
	if prompt != "" {
		fmt.Print(prompt)
	}
	
	line, err := sr.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(line), nil
}

// Clear clears the screen (not implemented for console)
func (sr *StandardRuntime) Clear() error {
	return nil
}
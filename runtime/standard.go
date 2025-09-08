// ABOUTME: Standard runtime implementation for console I/O operations
// ABOUTME: Production runtime that uses os.Stdout and os.Stdin for real console interaction

package runtime

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// StandardRuntime implements Runtime interface for console I/O
type StandardRuntime struct {
	reader *bufio.Reader
	rng    *rand.Rand
}

// NewStandardRuntime creates a new StandardRuntime instance
func NewStandardRuntime() *StandardRuntime {
	return &StandardRuntime{
		reader: bufio.NewReader(os.Stdin),
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Print outputs a string to stdout without a newline
func (std *StandardRuntime) Print(value string) error {
	_, err := fmt.Print(value)
	return err
}

// PrintLine outputs a string to stdout with a newline
func (std *StandardRuntime) PrintLine(value string) error {
	_, err := fmt.Println(value)
	return err
}

// Input prompts for user input and returns the entered string
func (std *StandardRuntime) Input(prompt string) (string, error) {
	if prompt != "" {
		fmt.Print(prompt)
	}

	line, err := std.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(line), nil
}

// Clear clears the screen (not implemented for console)
func (std *StandardRuntime) Clear() error {
	return nil
}

// Random returns a random float64 in [0,1)
func (std *StandardRuntime) Random() float64 {
	return std.rng.Float64()
}

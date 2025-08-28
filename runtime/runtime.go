// ABOUTME: Runtime interface and implementations for I/O operations and system interaction
// ABOUTME: Provides abstraction layer for console I/O to enable testing and different runtime environments

package runtime

// Runtime provides an interface for all I/O operations
// This allows the interpreter to work with different environments (console, test, etc.)
type Runtime interface {
	// Print outputs a string without a newline
	Print(value string) error

	// PrintLine outputs a string with a newline
	PrintLine(value string) error

	// Input prompts for user input and returns the entered string
	Input(prompt string) (string, error)

	// Clear clears the output (if supported by the runtime)
	Clear() error
}

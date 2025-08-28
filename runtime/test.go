// ABOUTME: Test runtime implementation for testing BASIC program I/O without console interaction
// ABOUTME: Mock runtime that captures output and provides scripted input for automated testing

package runtime

import (
	"fmt"
)

// TestRuntime implements Runtime interface for testing
// It captures all output and provides scripted input
type TestRuntime struct {
	outputBuffer []string
	inputQueue   []string
	inputIndex   int
}

// NewTestRuntime creates a new TestRuntime instance
func NewTestRuntime() *TestRuntime {
	return &TestRuntime{
		outputBuffer: make([]string, 0),
		inputQueue:   make([]string, 0),
		inputIndex:   0,
	}
}

// Print captures output without a newline
func (test *TestRuntime) Print(value string) error {
	test.outputBuffer = append(test.outputBuffer, value)
	return nil
}

// PrintLine captures output with a newline
func (test *TestRuntime) PrintLine(value string) error {
	test.outputBuffer = append(test.outputBuffer, value+"\n")
	return nil
}

// Input returns scripted input from the queue
func (test *TestRuntime) Input(prompt string) (string, error) {
	if prompt != "" {
		test.outputBuffer = append(test.outputBuffer, prompt)
	}
	
	if test.inputIndex >= len(test.inputQueue) {
		return "", fmt.Errorf("no more input available in test queue")
	}
	
	result := test.inputQueue[test.inputIndex]
	test.inputIndex++
	return result, nil
}

// Clear clears the output buffer
func (test *TestRuntime) Clear() error {
	test.outputBuffer = make([]string, 0)
	return nil
}

// GetOutput returns all captured output
func (test *TestRuntime) GetOutput() []string {
	return test.outputBuffer
}

// SetInput sets the input queue for testing
func (test *TestRuntime) SetInput(inputs []string) {
	test.inputQueue = inputs
	test.inputIndex = 0
}

// ResetOutput clears the output buffer for testing
func (test *TestRuntime) ResetOutput() {
	test.outputBuffer = make([]string, 0)
}
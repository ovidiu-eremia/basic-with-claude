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
func (tr *TestRuntime) Print(value string) error {
	tr.outputBuffer = append(tr.outputBuffer, value)
	return nil
}

// PrintLine captures output with a newline
func (tr *TestRuntime) PrintLine(value string) error {
	tr.outputBuffer = append(tr.outputBuffer, value+"\n")
	return nil
}

// Input returns scripted input from the queue
func (tr *TestRuntime) Input(prompt string) (string, error) {
	if prompt != "" {
		tr.outputBuffer = append(tr.outputBuffer, prompt)
	}
	
	if tr.inputIndex >= len(tr.inputQueue) {
		return "", fmt.Errorf("no more input available in test queue")
	}
	
	result := tr.inputQueue[tr.inputIndex]
	tr.inputIndex++
	return result, nil
}

// Clear clears the output buffer
func (tr *TestRuntime) Clear() error {
	tr.outputBuffer = make([]string, 0)
	return nil
}

// GetOutput returns all captured output
func (tr *TestRuntime) GetOutput() []string {
	return tr.outputBuffer
}

// SetInput sets the input queue for testing
func (tr *TestRuntime) SetInput(inputs []string) {
	tr.inputQueue = inputs
	tr.inputIndex = 0
}

// ResetOutput clears the output buffer for testing
func (tr *TestRuntime) ResetOutput() {
	tr.outputBuffer = make([]string, 0)
}
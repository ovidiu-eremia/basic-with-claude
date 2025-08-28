// ABOUTME: Tests for runtime interface and implementations
// ABOUTME: Comprehensive test suite for I/O operations and runtime behavior

package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestRuntime_Print(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "simple string",
			value:    "HELLO",
			expected: "HELLO",
		},
		{
			name:     "empty string",
			value:    "",
			expected: "",
		},
		{
			name:     "string with spaces",
			value:    "HELLO WORLD",
			expected: "HELLO WORLD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runtime := NewTestRuntime()

			err := runtime.Print(tt.value)
			require.NoError(t, err)

			output := runtime.GetOutput()
			require.Len(t, output, 1)
			assert.Equal(t, tt.expected, output[0])
		})
	}
}

func TestTestRuntime_PrintLine(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "simple string with newline",
			value:    "HELLO",
			expected: "HELLO\n",
		},
		{
			name:     "empty string with newline",
			value:    "",
			expected: "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runtime := NewTestRuntime()

			err := runtime.PrintLine(tt.value)
			require.NoError(t, err)

			output := runtime.GetOutput()
			require.Len(t, output, 1)
			assert.Equal(t, tt.expected, output[0])
		})
	}
}

func TestTestRuntime_MultiplePrints(t *testing.T) {
	runtime := NewTestRuntime()

	err := runtime.Print("HELLO")
	require.NoError(t, err)

	err = runtime.Print(" ")
	require.NoError(t, err)

	err = runtime.PrintLine("WORLD")
	require.NoError(t, err)

	output := runtime.GetOutput()
	require.Len(t, output, 3)
	assert.Equal(t, "HELLO", output[0])
	assert.Equal(t, " ", output[1])
	assert.Equal(t, "WORLD\n", output[2])
}

func TestTestRuntime_Input(t *testing.T) {
	tests := []struct {
		name          string
		inputQueue    []string
		prompt        string
		expectedInput string
	}{
		{
			name:          "simple input",
			inputQueue:    []string{"test input"},
			prompt:        "? ",
			expectedInput: "test input",
		},
		{
			name:          "empty prompt",
			inputQueue:    []string{"hello"},
			prompt:        "",
			expectedInput: "hello",
		},
		{
			name:          "numeric input",
			inputQueue:    []string{"42"},
			prompt:        "Enter number: ",
			expectedInput: "42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runtime := NewTestRuntime()
			runtime.SetInput(tt.inputQueue)

			result, err := runtime.Input(tt.prompt)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedInput, result)
		})
	}
}

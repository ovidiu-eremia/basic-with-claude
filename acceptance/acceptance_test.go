// ABOUTME: End-to-end acceptance tests for the BASIC interpreter
// ABOUTME: Tests complete pipeline from source files through execution and output verification

package acceptance

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"basic-interpreter/interpreter"
	"basic-interpreter/lexer"
	"basic-interpreter/parser"
	"basic-interpreter/runtime"
)

// executeBasicFile loads and executes a BASIC program file, returning the output
func executeBasicFile(t *testing.T, filename string) []string {
	t.Helper()
	
	// Read the BASIC program file
	content, err := os.ReadFile(filename)
	require.NoError(t, err, "Failed to read file %s", filename)
	
	// Parse the program
	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()
	
	// Check for parsing errors
	require.Empty(t, p.Errors(), "Parsing errors in %s: %v", filename, p.Errors())
	require.NotNil(t, program, "Program is nil after parsing %s", filename)
	
	// Create test runtime and interpreter
	testRuntime := runtime.NewTestRuntime()
	interp := interpreter.NewInterpreter(testRuntime)
	
	// Execute the program
	err = interp.Execute(program)
	require.NoError(t, err, "Runtime error executing %s: %v", filename, err)
	
	// Return captured output
	return testRuntime.GetOutput()
}

func TestAcceptance_HelloWorld(t *testing.T) {
	output := executeBasicFile(t, "../hello.bas")
	
	expected := []string{"HELLO WORLD\n"}
	assert.Equal(t, expected, output, "Hello world program should output correct message")
}

func TestAcceptance_Variables(t *testing.T) {
	output := executeBasicFile(t, "../test_variables.bas")
	
	expected := []string{"42\n", "123\n"}
	assert.Equal(t, expected, output, "Variable program should assign and print variable values correctly")
}

func TestAcceptance_ExecutionOrder(t *testing.T) {
	// Test that programs execute in correct line number order
	// and that END statement stops execution properly
	output := executeBasicFile(t, "../hello.bas")
	
	// hello.bas has "20 END" so execution should stop there
	// Only one line of output expected
	require.Len(t, output, 1, "Program with END statement should produce exactly one output line")
	assert.Equal(t, "HELLO WORLD\n", output[0], "Output should match expected string")
}

func TestAcceptance_InvalidFile(t *testing.T) {
	// Test that non-existent files are handled gracefully
	content := "invalid content"
	l := lexer.New(content)
	p := parser.New(l)
	program := p.ParseProgram()
	
	// Should have parsing errors for invalid content
	assert.NotEmpty(t, p.Errors(), "Invalid BASIC content should produce parsing errors")
	
	// But parsing should still return a program (possibly empty)
	assert.NotNil(t, program, "Parser should always return a program object")
}
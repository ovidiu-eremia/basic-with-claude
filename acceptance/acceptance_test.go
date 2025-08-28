// ABOUTME: End-to-end acceptance tests for the BASIC interpreter
// ABOUTME: Tests complete pipeline from inline program strings through execution and output verification

package acceptance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"basic-interpreter/interpreter"
	"basic-interpreter/lexer"
	"basic-interpreter/parser"
	"basic-interpreter/runtime"
)

// executeBasicProgram parses and executes a BASIC program string, returning the output
func executeBasicProgram(t *testing.T, program string) []string {
	t.Helper()

	// Parse the program
	l := lexer.New(program)
	p := parser.New(l)
	ast := p.ParseProgram()

	// Check for parsing errors
	require.Empty(t, p.Errors(), "Parsing errors in program: %v", p.Errors())
	require.NotNil(t, ast, "Program is nil after parsing")

	// Create test runtime and interpreter
	testRuntime := runtime.NewTestRuntime()
	interp := interpreter.NewInterpreter(testRuntime)

	// Execute the program
	err := interp.Execute(ast)
	require.NoError(t, err, "Runtime error executing program: %v", err)

	// Return captured output
	return testRuntime.GetOutput()
}

func TestAcceptance_HelloWorld(t *testing.T) {
	program := `10 PRINT "HELLO WORLD"
20 END`

	output := executeBasicProgram(t, program)

	expected := []string{"HELLO WORLD\n"}
	assert.Equal(t, expected, output, "Hello world program should output correct message")
}

func TestAcceptance_Variables(t *testing.T) {
	program := `10 LET A = 42
20 PRINT A
30 X = 123
40 PRINT X`

	output := executeBasicProgram(t, program)

	expected := []string{"42\n", "123\n"}
	assert.Equal(t, expected, output, "Variable program should assign and print variable values correctly")
}

func TestAcceptance_ExecutionOrder(t *testing.T) {
	// Test that programs execute in correct line number order
	// and that END statement stops execution properly
	program := `10 PRINT "HELLO WORLD"
20 END
30 PRINT "THIS SHOULD NOT EXECUTE"`

	output := executeBasicProgram(t, program)

	// Program has "20 END" so execution should stop there
	// Only one line of output expected
	require.Len(t, output, 1, "Program with END statement should produce exactly one output line")
	assert.Equal(t, "HELLO WORLD\n", output[0], "Output should match expected string")
}

func TestAcceptance_StringVariables(t *testing.T) {
	program := `10 LET A$ = "HELLO"
20 PRINT A$
30 NAME$ = "WORLD"
40 PRINT NAME$`

	output := executeBasicProgram(t, program)

	expected := []string{"HELLO\n", "WORLD\n"}
	assert.Equal(t, expected, output, "String variables program should assign and print string values correctly")
}

func TestAcceptance_ArithmeticExpressions(t *testing.T) {
	program := `10 PRINT 2 + 3 * 4
20 PRINT (2 + 3) * 4
30 A = 5
40 B = 3
50 PRINT A * B + 1
60 PRINT A + B * 2
70 PRINT 10 / 2
80 PRINT 2 ^ 3`

	output := executeBasicProgram(t, program)

	expected := []string{
		"14\n", // 2 + 3 * 4 = 2 + 12 = 14 (precedence test)
		"20\n", // (2 + 3) * 4 = 5 * 4 = 20 (parentheses test)
		"16\n", // A * B + 1 = 5 * 3 + 1 = 15 + 1 = 16 (variables in expressions)
		"11\n", // A + B * 2 = 5 + 3 * 2 = 5 + 6 = 11 (precedence with variables)
		"5\n",  // 10 / 2 = 5 (division)
		"8\n",  // 2 ^ 3 = 8 (exponentiation)
	}
	assert.Equal(t, expected, output, "Arithmetic expressions should evaluate with correct precedence")
}

func TestAcceptance_InvalidProgram(t *testing.T) {
	// Test that invalid BASIC content is handled gracefully
	program := "invalid content"

	l := lexer.New(program)
	p := parser.New(l)
	ast := p.ParseProgram()

	// Should have parsing errors for invalid content
	assert.NotEmpty(t, p.Errors(), "Invalid BASIC content should produce parsing errors")

	// But parsing should still return a program (possibly empty)
	assert.NotNil(t, ast, "Parser should always return a program object")
}

func TestAcceptance_MixedFeatures(t *testing.T) {
	// Test a program that combines multiple features that are already implemented
	program := `10 LET MESSAGE$ = "Numbers: "
20 A = 10
30 B = 5
40 PRINT MESSAGE$
50 PRINT A + B * 2
60 END
70 PRINT "This should not execute"`

	output := executeBasicProgram(t, program)

	expected := []string{
		"Numbers: \n",
		"20\n", // A + B * 2 = 10 + 5 * 2 = 10 + 10 = 20
	}
	assert.Equal(t, expected, output, "Mixed features should work together correctly")
}

func TestAcceptance_ComplexArithmetic(t *testing.T) {
	// Test nested arithmetic with multiple levels of precedence
	program := `10 PRINT ((2 + 3) * 4 - 1) ^ 2
20 A = 2
30 B = 3  
40 C = 4
50 PRINT (A + B) * C - A ^ B`

	output := executeBasicProgram(t, program)

	expected := []string{
		"361\n", // ((2 + 3) * 4 - 1) ^ 2 = (5 * 4 - 1) ^ 2 = (20 - 1) ^ 2 = 19 ^ 2 = 361
		"12\n",  // (A + B) * C - A ^ B = (2 + 3) * 4 - 2 ^ 3 = 5 * 4 - 8 = 20 - 8 = 12
	}
	assert.Equal(t, expected, output, "Complex arithmetic expressions should evaluate correctly")
}

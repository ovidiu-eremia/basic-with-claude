// ABOUTME: Tests for FOR loops with colon-separated statements to verify statement-level positioning
// ABOUTME: Ensures FOR loops correctly continue execution after completing on multi-statement lines

package interpreter

import (
	"testing"

	"basic-interpreter/lexer"
	"basic-interpreter/parser"
	"basic-interpreter/runtime"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createLexer creates a lexer for the given program string
func createLexer(program string) *lexer.Lexer {
	return lexer.New(program)
}

func TestForLoopWithColonSeparatedStatements(t *testing.T) {
	tests := []struct {
		name           string
		program        string
		expectedOutput []string
	}{
		{
			name:    "Simple FOR loop with PRINT on same line",
			program: "10 FOR I = 1 TO 3: PRINT I: NEXT I",
			expectedOutput: []string{
				"1",
				"2",
				"3",
			},
		},
		{
			name:    "FOR loop with multiple statements after NEXT",
			program: "10 FOR I = 1 TO 2: PRINT I: NEXT I: PRINT \"DONE\"",
			expectedOutput: []string{
				"1",
				"2",
				"DONE",
			},
		},
		{
			name:    "FOR loop with statements before and after",
			program: "10 PRINT \"START\": FOR I = 1 TO 2: PRINT I: NEXT I: PRINT \"END\"",
			expectedOutput: []string{
				"START",
				"1",
				"2",
				"END",
			},
		},
		{
			name: "FOR loop spanning multiple lines with colon continuation",
			program: `10 FOR I = 1 TO 2: PRINT I
20 NEXT I: PRINT "FINISHED"`,
			expectedOutput: []string{
				"1",
				"2",
				"FINISHED",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRuntime := runtime.NewTestRuntime()
			interpreter := NewInterpreter(testRuntime)

			// Parse the program
			lexer := createLexer(tt.program)
			parser := parser.New(lexer)
			program := parser.ParseProgram()

			// Check for parsing errors
			if parser.ParseError() != nil {
				t.Fatalf("Parser error: %v", parser.ParseError())
			}

			// Execute the program
			err := interpreter.Execute(program)
			require.NoError(t, err, "Program execution should not fail")

			// Verify output - GetOutput returns strings with newlines, so we need to trim them
			rawOutput := testRuntime.GetOutput()
			actualOutput := make([]string, 0, len(rawOutput))
			for _, line := range rawOutput {
				if line == "\n" {
					actualOutput = append(actualOutput, "")
				} else if len(line) > 0 && line[len(line)-1] == '\n' {
					actualOutput = append(actualOutput, line[:len(line)-1])
				} else {
					actualOutput = append(actualOutput, line)
				}
			}
			assert.Equal(t, tt.expectedOutput, actualOutput, "Output should match expected")
		})
	}
}

func TestForLoopStatementPositioning(t *testing.T) {
	testRuntime := runtime.NewTestRuntime()
	interpreter := NewInterpreter(testRuntime)

	// Test that demonstrates the bug would have occurred with the old implementation
	// This program has a FOR loop followed by multiple statements on the same line
	program := `10 FOR I = 1 TO 2: PRINT "LOOP": PRINT I: NEXT I: PRINT "AFTER": PRINT "FINAL"`

	lexer := createLexer(program)
	parser := parser.New(lexer)
	parsedProgram := parser.ParseProgram()

	require.Nil(t, parser.ParseError(), "Should parse without errors")

	err := interpreter.Execute(parsedProgram)
	require.NoError(t, err, "Should execute without errors")

	expectedOutput := []string{
		"LOOP",
		"1",
		"LOOP",
		"2",
		"AFTER",
		"FINAL",
	}

	rawOutput := testRuntime.GetOutput()
	actualOutput := make([]string, 0, len(rawOutput))
	for _, line := range rawOutput {
		if line == "\n" {
			actualOutput = append(actualOutput, "")
		} else if len(line) > 0 && line[len(line)-1] == '\n' {
			actualOutput = append(actualOutput, line[:len(line)-1])
		} else {
			actualOutput = append(actualOutput, line)
		}
	}
	assert.Equal(t, expectedOutput, actualOutput,
		"Should execute all statements after NEXT I, not jump to next line")
}

func TestNestedForLoopsWithColons(t *testing.T) {
	testRuntime := runtime.NewTestRuntime()
	interpreter := NewInterpreter(testRuntime)

	// Test nested FOR loops with colon separation
	program := `10 FOR I = 1 TO 2: FOR J = 1 TO 2: PRINT I: PRINT J: NEXT J: NEXT I: PRINT "DONE"`

	lexer := createLexer(program)
	parser := parser.New(lexer)
	parsedProgram := parser.ParseProgram()

	require.Nil(t, parser.ParseError(), "Should parse without errors")

	err := interpreter.Execute(parsedProgram)
	require.NoError(t, err, "Should execute without errors")

	expectedOutput := []string{
		"1",
		"1",
		"1",
		"2",
		"2",
		"1",
		"2",
		"2",
		"DONE",
	}

	rawOutput := testRuntime.GetOutput()
	actualOutput := make([]string, 0, len(rawOutput))
	for _, line := range rawOutput {
		if line == "\n" {
			actualOutput = append(actualOutput, "")
		} else if len(line) > 0 && line[len(line)-1] == '\n' {
			actualOutput = append(actualOutput, line[:len(line)-1])
		} else {
			actualOutput = append(actualOutput, line)
		}
	}
	assert.Equal(t, expectedOutput, actualOutput, "Nested loops should work correctly with colons")
}

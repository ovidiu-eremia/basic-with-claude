package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"basic-interpreter/lexer"
)

func TestParser_FunctionCall(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *Program
	}{
		{
			name:  "LEN function call",
			input: `10 PRINT LEN("HELLO")`,
			expected: program(
				line(10, 1, printStmt(funcCall("LEN", []Expression{str("HELLO", 1)}, 1), 1)),
			),
		},
		{
			name:  "LEFT$ function call",
			input: `10 PRINT LEFT$("HELLO", 3)`,
			expected: program(
				line(10, 1, printStmt(funcCall("LEFT$", []Expression{str("HELLO", 1), num("3", 1)}, 1), 1)),
			),
		},
		{
			name:  "RIGHT$ function call",
			input: `10 PRINT RIGHT$("WORLD", 2)`,
			expected: program(
				line(10, 1, printStmt(funcCall("RIGHT$", []Expression{str("WORLD", 1), num("2", 1)}, 1), 1)),
			),
		},
		{
			name:  "Function call with variable argument",
			input: `10 PRINT LEN(A$)`,
			expected: program(
				line(10, 1, printStmt(funcCall("LEN", []Expression{varRef("A$", 1)}, 1), 1)),
			),
		},
		{
			name:  "Function call with no arguments",
			input: `10 PRINT RND()`,
			expected: program(
				line(10, 1, printStmt(funcCall("RND", []Expression{}, 1), 1)),
			),
		},
		{
			name:  "Function call in assignment",
			input: `10 LET L = LEN("TEST")`,
			expected: program(
				line(10, 1, letStmt("L", funcCall("LEN", []Expression{str("TEST", 1)}, 1), 1)),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			// Check for parsing errors
			if p.ParseError() != nil {
				t.Fatalf("parser error: %v", p.ParseError())
			}

			require.NotNil(t, program)
			assert.Equal(t, tt.expected, program)
		})
	}
}

func TestParser_FunctionCallErrors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "missing closing parenthesis",
			input:       `10 PRINT LEN("HELLO"`,
			expectError: true,
		},
		{
			name:        "missing argument in two-argument function",
			input:       `10 PRINT LEFT$("HELLO",)`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			if tt.expectError {
				assert.NotNil(t, p.ParseError(), "Expected parsing error but got none")
			} else {
				assert.Nil(t, p.ParseError(), "Unexpected parsing error: %v", p.ParseError())
				require.NotNil(t, program)
			}
		})
	}
}

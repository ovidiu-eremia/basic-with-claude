package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"basic-interpreter/lexer"
)

func TestParser_NestedFunctionCalls(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *Program
	}{
		{
			name:  "LEN_of_LEFT$",
			input: `10 PRINT LEN(LEFT$("HELLO", 2))`,
			expected: program(
				line(10, 1,
					printStmt(
						funcCall("LEN", []Expression{
							funcCall("LEFT$", []Expression{str("HELLO", 1), num("2", 1)}, 1),
						}, 1),
						1,
					),
				),
			),
		},
		{
			name:  "LEFT$_of_RIGHT$",
			input: `10 PRINT LEFT$(RIGHT$("HELLO", 4), 2)`,
			expected: program(
				line(10, 1,
					printStmt(
						funcCall("LEFT$", []Expression{
							funcCall("RIGHT$", []Expression{str("HELLO", 1), num("4", 1)}, 1),
							num("2", 1),
						}, 1),
						1,
					),
				),
			),
		},
		{
			name:  "RIGHT$_of_LEFT$",
			input: `10 PRINT RIGHT$(LEFT$("HELLO", 4), 2)`,
			expected: program(
				line(10, 1,
					printStmt(
						funcCall("RIGHT$", []Expression{
							funcCall("LEFT$", []Expression{str("HELLO", 1), num("4", 1)}, 1),
							num("2", 1),
						}, 1),
						1,
					),
				),
			),
		},
		{
			name:  "Triple_nesting_LEN_LEFT_RIGHT",
			input: `10 PRINT LEN(LEFT$(RIGHT$(LEFT$("ABCDE", 4), 3), 2))`,
			expected: program(
				line(10, 1,
					printStmt(
						funcCall("LEN", []Expression{
							funcCall("LEFT$", []Expression{
								funcCall("RIGHT$", []Expression{
									funcCall("LEFT$", []Expression{str("ABCDE", 1), num("4", 1)}, 1),
									num("3", 1),
								}, 1),
								num("2", 1),
							}, 1),
						}, 1),
						1,
					),
				),
			),
		},
		{
			name:  "Function_in_second_arg",
			input: `10 PRINT RIGHT$("ABCDE", LEN("XY"))`,
			expected: program(
				line(10, 1,
					printStmt(
						funcCall("RIGHT$", []Expression{
							str("ABCDE", 1),
							funcCall("LEN", []Expression{str("XY", 1)}, 1),
						}, 1),
						1,
					),
				),
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

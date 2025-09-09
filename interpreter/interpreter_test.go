// ABOUTME: Tests for the BASIC interpreter execution engine
// ABOUTME: Comprehensive test suite for AST execution and program state management

package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"basic-interpreter/parser"
	"basic-interpreter/runtime"
)

func TestInterpreter_ExecutePrintStatement(t *testing.T) {
	tests := []struct {
		name           string
		program        *parser.Program
		expectedOutput []string
	}{
		{
			name: "single print statement",
			program: &parser.Program{
				Lines: []*parser.Line{
					{
						Number: 10,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.StringLiteral{
									Value: "HELLO WORLD",
								},
							},
						},
					},
				},
			},
			expectedOutput: []string{"HELLO WORLD\n"},
		},
		{
			name: "numeric literal print",
			program: &parser.Program{
				Lines: []*parser.Line{
					{
						Number: 10,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.NumberLiteral{
									Value: "42",
								},
							},
						},
					},
				},
			},
			expectedOutput: []string{"42\n"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRuntime := runtime.NewTestRuntime()
			interpreter := NewInterpreter(testRuntime)

			err := interpreter.Execute(tt.program)
			require.NoError(t, err)

			output := testRuntime.GetOutput()
			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}

func TestInterpreter_NumericVariables(t *testing.T) {
	tests := []struct {
		name           string
		program        *parser.Program
		expectedOutput []string
	}{
		{
			name: "variable assignment with LET",
			program: &parser.Program{
				Lines: []*parser.Line{
					{
						Number: 10,
						Statements: []parser.Statement{
							&parser.LetStatement{
								Variable: "A",
								Expression: &parser.NumberLiteral{
									Value: "42",
								},
							},
						},
					},
					{
						Number: 20,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.VariableReference{
									Name: "A",
								},
							},
						},
					},
				},
			},
			expectedOutput: []string{"42\n"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRuntime := runtime.NewTestRuntime()
			interpreter := NewInterpreter(testRuntime)

			err := interpreter.Execute(tt.program)
			require.NoError(t, err)

			output := testRuntime.GetOutput()
			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}

func TestInterpreter_StringVariables(t *testing.T) {
	tests := []struct {
		name           string
		program        *parser.Program
		expectedOutput []string
	}{
		{
			name: "string variable assignment with LET",
			program: &parser.Program{
				Lines: []*parser.Line{
					{
						Number: 10,
						Statements: []parser.Statement{
							&parser.LetStatement{
								Variable: "A$",
								Expression: &parser.StringLiteral{
									Value: "HELLO",
								},
							},
						},
					},
					{
						Number: 20,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.VariableReference{
									Name: "A$",
								},
							},
						},
					},
				},
			},
			expectedOutput: []string{"HELLO\n"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRuntime := runtime.NewTestRuntime()
			interpreter := NewInterpreter(testRuntime)

			err := interpreter.Execute(tt.program)
			require.NoError(t, err)

			output := testRuntime.GetOutput()
			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}

func TestInterpreter_ArithmeticExpressions(t *testing.T) {
	tests := []struct {
		name           string
		program        *parser.Program
		expectedOutput []string
	}{
		{
			name: "simple addition",
			program: &parser.Program{
				Lines: []*parser.Line{
					{
						Number: 10,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.BinaryOperation{
									Left:     &parser.NumberLiteral{Value: "2"},
									Operator: "+",
									Right:    &parser.NumberLiteral{Value: "3"},
								},
							},
						},
					},
				},
			},
			expectedOutput: []string{"5\n"},
		},
		{
			name: "operator precedence",
			program: &parser.Program{
				Lines: []*parser.Line{
					{
						Number: 10,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.BinaryOperation{
									Left:     &parser.NumberLiteral{Value: "2"},
									Operator: "+",
									Right: &parser.BinaryOperation{
										Left:     &parser.NumberLiteral{Value: "3"},
										Operator: "*",
										Right:    &parser.NumberLiteral{Value: "4"},
									},
								},
							},
						},
					},
				},
			},
			expectedOutput: []string{"14\n"},
		},
		{
			name: "variables in expressions",
			program: &parser.Program{
				Lines: []*parser.Line{
					{
						Number: 10,
						Statements: []parser.Statement{
							&parser.LetStatement{
								Variable:   "A",
								Expression: &parser.NumberLiteral{Value: "5"},
							},
						},
					},
					{
						Number: 20,
						Statements: []parser.Statement{
							&parser.LetStatement{
								Variable:   "B",
								Expression: &parser.NumberLiteral{Value: "3"},
							},
						},
					},
					{
						Number: 30,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.BinaryOperation{
									Left:     &parser.VariableReference{Name: "A"},
									Operator: "*",
									Right: &parser.BinaryOperation{
										Left:     &parser.VariableReference{Name: "B"},
										Operator: "+",
										Right:    &parser.NumberLiteral{Value: "1"},
									},
								},
							},
						},
					},
				},
			},
			expectedOutput: []string{"20\n"},
		},
		{
			name: "division and power",
			program: &parser.Program{
				Lines: []*parser.Line{
					{
						Number: 10,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.BinaryOperation{
									Left:     &parser.NumberLiteral{Value: "10"},
									Operator: "/",
									Right:    &parser.NumberLiteral{Value: "2"},
								},
							},
						},
					},
					{
						Number: 20,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.BinaryOperation{
									Left:     &parser.NumberLiteral{Value: "2"},
									Operator: "^",
									Right:    &parser.NumberLiteral{Value: "3"},
								},
							},
						},
					},
				},
			},
			expectedOutput: []string{"5\n", "8\n"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRuntime := runtime.NewTestRuntime()
			interpreter := NewInterpreter(testRuntime)

			err := interpreter.Execute(tt.program)
			require.NoError(t, err)

			output := testRuntime.GetOutput()
			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}

func TestInterpreter_ArithmeticErrors(t *testing.T) {
	tests := []struct {
		name        string
		program     *parser.Program
		expectError bool
	}{
		{
			name: "division by zero",
			program: &parser.Program{
				Lines: []*parser.Line{
					{
						Number: 10,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.BinaryOperation{
									Left:     &parser.NumberLiteral{Value: "10"},
									Operator: "/",
									Right:    &parser.NumberLiteral{Value: "0"},
								},
							},
						},
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRuntime := runtime.NewTestRuntime()
			interpreter := NewInterpreter(testRuntime)

			err := interpreter.Execute(tt.program)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

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
									Line:  1,
								},
								Line: 1,
							},
						},
						SourceLine: 1,
					},
				},
			},
			expectedOutput: []string{"HELLO WORLD\n"},
		},
		{
			name: "multiple print statements",
			program: &parser.Program{
				Lines: []*parser.Line{
					{
						Number: 10,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.StringLiteral{
									Value: "HELLO",
									Line:  1,
								},
								Line: 1,
							},
						},
						SourceLine: 1,
					},
					{
						Number: 20,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.StringLiteral{
									Value: "WORLD",
									Line:  2,
								},
								Line: 2,
							},
						},
						SourceLine: 2,
					},
				},
			},
			expectedOutput: []string{"HELLO\n", "WORLD\n"},
		},
		{
			name: "empty string print",
			program: &parser.Program{
				Lines: []*parser.Line{
					{
						Number: 10,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.StringLiteral{
									Value: "",
									Line:  1,
								},
								Line: 1,
							},
						},
						SourceLine: 1,
					},
				},
			},
			expectedOutput: []string{"\n"},
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

func TestInterpreter_ExecuteWithEndStatement(t *testing.T) {
	testRuntime := runtime.NewTestRuntime()
	interpreter := NewInterpreter(testRuntime)
	
	program := &parser.Program{
		Lines: []*parser.Line{
			{
				Number: 10,
				Statements: []parser.Statement{
					&parser.PrintStatement{
						Expression: &parser.StringLiteral{
							Value: "BEFORE END",
							Line:  1,
						},
						Line: 1,
					},
				},
				SourceLine: 1,
			},
			{
				Number: 20,
				Statements: []parser.Statement{
					&parser.EndStatement{Line: 2},
				},
				SourceLine: 2,
			},
			{
				Number: 30,
				Statements: []parser.Statement{
					&parser.PrintStatement{
						Expression: &parser.StringLiteral{
							Value: "AFTER END",
							Line:  3,
						},
						Line: 3,
					},
				},
				SourceLine: 3,
			},
		},
	}
	
	err := interpreter.Execute(program)
	require.NoError(t, err)
	
	output := testRuntime.GetOutput()
	require.Len(t, output, 1)
	assert.Equal(t, "BEFORE END\n", output[0])
}

func TestInterpreter_EmptyProgram(t *testing.T) {
	testRuntime := runtime.NewTestRuntime()
	interpreter := NewInterpreter(testRuntime)
	
	program := &parser.Program{Lines: []*parser.Line{}}
	
	err := interpreter.Execute(program)
	require.NoError(t, err)
	
	output := testRuntime.GetOutput()
	assert.Len(t, output, 0)
}

func TestInterpreter_VariableAssignment(t *testing.T) {
	tests := []struct {
		name           string
		program        *parser.Program
		expectedOutput []string
	}{
		{
			name: "LET assignment and PRINT",
			program: &parser.Program{
				Lines: []*parser.Line{
					{
						Number: 10,
						Statements: []parser.Statement{
							&parser.LetStatement{
								Variable: "A",
								Expression: &parser.NumberLiteral{
									Value: "42",
									Line:  1,
								},
								Line: 1,
							},
						},
						SourceLine: 1,
					},
					{
						Number: 20,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.VariableReference{
									Name: "A",
									Line: 2,
								},
								Line: 2,
							},
						},
						SourceLine: 2,
					},
				},
			},
			expectedOutput: []string{"42\n"},
		},
		{
			name: "assignment without LET",
			program: &parser.Program{
				Lines: []*parser.Line{
					{
						Number: 10,
						Statements: []parser.Statement{
							&parser.LetStatement{
								Variable: "X",
								Expression: &parser.NumberLiteral{
									Value: "123",
									Line:  1,
								},
								Line: 1,
							},
						},
						SourceLine: 1,
					},
					{
						Number: 20,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.VariableReference{
									Name: "X",
									Line: 2,
								},
								Line: 2,
							},
						},
						SourceLine: 2,
					},
				},
			},
			expectedOutput: []string{"123\n"},
		},
		{
			name: "multiple variables",
			program: &parser.Program{
				Lines: []*parser.Line{
					{
						Number: 10,
						Statements: []parser.Statement{
							&parser.LetStatement{
								Variable: "A",
								Expression: &parser.NumberLiteral{
									Value: "10",
									Line:  1,
								},
								Line: 1,
							},
						},
						SourceLine: 1,
					},
					{
						Number: 20,
						Statements: []parser.Statement{
							&parser.LetStatement{
								Variable: "B",
								Expression: &parser.NumberLiteral{
									Value: "20",
									Line:  2,
								},
								Line: 2,
							},
						},
						SourceLine: 2,
					},
					{
						Number: 30,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.VariableReference{
									Name: "A",
									Line: 3,
								},
								Line: 3,
							},
						},
						SourceLine: 3,
					},
					{
						Number: 40,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.VariableReference{
									Name: "B",
									Line: 4,
								},
								Line: 4,
							},
						},
						SourceLine: 4,
					},
				},
			},
			expectedOutput: []string{"10\n", "20\n"},
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
									Line:  1,
								},
								Line: 1,
							},
						},
						SourceLine: 1,
					},
					{
						Number: 20,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.VariableReference{
									Name: "A$",
									Line: 2,
								},
								Line: 2,
							},
						},
						SourceLine: 2,
					},
				},
			},
			expectedOutput: []string{"HELLO\n"},
		},
		{
			name: "string variable assignment without LET",
			program: &parser.Program{
				Lines: []*parser.Line{
					{
						Number: 10,
						Statements: []parser.Statement{
							&parser.LetStatement{
								Variable: "NAME$",
								Expression: &parser.StringLiteral{
									Value: "JOHN DOE",
									Line:  1,
								},
								Line: 1,
							},
						},
						SourceLine: 1,
					},
					{
						Number: 20,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.VariableReference{
									Name: "NAME$",
									Line: 2,
								},
								Line: 2,
							},
						},
						SourceLine: 2,
					},
				},
			},
			expectedOutput: []string{"JOHN DOE\n"},
		},
		{
			name: "mixed numeric and string variables",
			program: &parser.Program{
				Lines: []*parser.Line{
					{
						Number: 10,
						Statements: []parser.Statement{
							&parser.LetStatement{
								Variable: "A",
								Expression: &parser.NumberLiteral{
									Value: "42",
									Line:  1,
								},
								Line: 1,
							},
						},
						SourceLine: 1,
					},
					{
						Number: 20,
						Statements: []parser.Statement{
							&parser.LetStatement{
								Variable: "B$",
								Expression: &parser.StringLiteral{
									Value: "ANSWER",
									Line:  2,
								},
								Line: 2,
							},
						},
						SourceLine: 2,
					},
					{
						Number: 30,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.VariableReference{
									Name: "B$",
									Line: 3,
								},
								Line: 3,
							},
						},
						SourceLine: 3,
					},
					{
						Number: 40,
						Statements: []parser.Statement{
							&parser.PrintStatement{
								Expression: &parser.VariableReference{
									Name: "A",
									Line: 4,
								},
								Line: 4,
							},
						},
						SourceLine: 4,
					},
				},
			},
			expectedOutput: []string{"ANSWER\n", "42\n"},
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
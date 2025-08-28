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
			interpreter := New(testRuntime)
			
			err := interpreter.Execute(tt.program)
			require.NoError(t, err)
			
			output := testRuntime.GetOutput()
			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}

func TestInterpreter_ExecuteWithEndStatement(t *testing.T) {
	testRuntime := runtime.NewTestRuntime()
	interpreter := New(testRuntime)
	
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
	interpreter := New(testRuntime)
	
	program := &parser.Program{Lines: []*parser.Line{}}
	
	err := interpreter.Execute(program)
	require.NoError(t, err)
	
	output := testRuntime.GetOutput()
	assert.Len(t, output, 0)
}
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

type AcceptanceTest struct {
	name        string
	program     string
	inputs      []string
	expected    []string
	wantErr     bool
	errContains string
}

// executeBasicProgram parses and executes a BASIC program string, returning the output
func executeBasicProgram(t *testing.T, program string) ([]string, error) {
	t.Helper()

	// Parse the program
	l := lexer.New(program)
	p := parser.New(l)
	ast := p.ParseProgram()

	// Check for parsing errors
	if len(p.Errors()) > 0 {
		return nil, assert.AnError
	}
	if ast == nil {
		return nil, assert.AnError
	}

	// Create test runtime and interpreter
	testRuntime := runtime.NewTestRuntime()
	interp := interpreter.NewInterpreter(testRuntime)

	// Execute the program
	err := interp.Execute(ast)
	if err != nil {
		return nil, err
	}

	// Return captured output
	return testRuntime.GetOutput(), nil
}

func TestAcceptance(t *testing.T) {
	tests := []AcceptanceTest{
		{
			name: "HelloWorld",
			program: `10 PRINT "HELLO WORLD"
20 END`,
			expected: []string{"HELLO WORLD\n"},
		},
		{
			name: "Variables",
			program: `10 LET A = 42
20 PRINT A
30 X = 123
40 PRINT X`,
			expected: []string{"42\n", "123\n"},
		},
		{
			name: "ExecutionOrder",
			program: `10 PRINT "HELLO WORLD"
20 END
30 PRINT "THIS SHOULD NOT EXECUTE"`,
			expected: []string{"HELLO WORLD\n"},
		},
		{
			name: "StringVariables",
			program: `10 LET A$ = "HELLO"
20 PRINT A$
30 NAME$ = "WORLD"
40 PRINT NAME$`,
			expected: []string{"HELLO\n", "WORLD\n"},
		},
		{
			name: "ArithmeticExpressions",
			program: `10 PRINT 2 + 3 * 4
20 PRINT (2 + 3) * 4
30 A = 5
40 B = 3
50 PRINT A * B + 1
60 PRINT A + B * 2
70 PRINT 10 / 2
80 PRINT 2 ^ 3`,
			expected: []string{
				"14\n", // 2 + 3 * 4 = 2 + 12 = 14 (precedence test)
				"20\n", // (2 + 3) * 4 = 5 * 4 = 20 (parentheses test)
				"16\n", // A * B + 1 = 5 * 3 + 1 = 15 + 1 = 16 (variables in expressions)
				"11\n", // A + B * 2 = 5 + 3 * 2 = 5 + 6 = 11 (precedence with variables)
				"5\n",  // 10 / 2 = 5 (division)
				"8\n",  // 2 ^ 3 = 8 (exponentiation)
			},
		},
		{
			name:        "InvalidProgram",
			program:     "invalid content",
			wantErr:     true,
			errContains: "",
		},
		{
			name: "MixedFeatures",
			program: `10 LET MESSAGE$ = "Numbers: "
20 A = 10
30 B = 5
40 PRINT MESSAGE$
50 PRINT A + B * 2
60 END
70 PRINT "This should not execute"`,
			expected: []string{
				"Numbers: \n",
				"20\n", // A + B * 2 = 10 + 5 * 2 = 10 + 10 = 20
			},
		},
		{
			name: "ComplexArithmetic",
			program: `10 PRINT ((2 + 3) * 4 - 1) ^ 2
20 A = 2
30 B = 3  
40 C = 4
50 PRINT (A + B) * C - A ^ B`,
			expected: []string{
				"361\n", // ((2 + 3) * 4 - 1) ^ 2 = (5 * 4 - 1) ^ 2 = (20 - 1) ^ 2 = 19 ^ 2 = 361
				"12\n",  // (A + B) * C - A ^ B = (2 + 3) * 4 - 2 ^ 3 = 5 * 4 - 8 = 20 - 8 = 12
			},
		},
		{
			name: "StopStatement",
			program: `10 PRINT "START"
20 STOP
30 PRINT "NEVER REACHED"`,
			expected: []string{
				"START\n",
			},
		},
		{
			name: "RunStatement",
			program: `10 PRINT "BEFORE RUN"
20 RUN
30 PRINT "AFTER RUN"`,
			expected: []string{
				"BEFORE RUN\n",
				"AFTER RUN\n",
			},
		},
		{
			name: "GotoStatement",
			program: `10 PRINT "BEFORE JUMP"
20 GOTO 50
30 PRINT "SKIPPED"
40 PRINT "ALSO SKIPPED"
50 PRINT "AFTER JUMP"`,
			expected: []string{
				"BEFORE JUMP\n",
				"AFTER JUMP\n",
			},
		},
		{
			name: "GotoBackward",
			program: `10 PRINT "FIRST"
20 GOTO 40
30 PRINT "NEVER"
40 PRINT "SECOND"
50 GOTO 70
60 PRINT "ALSO NEVER"
70 PRINT "THIRD"`,
			expected: []string{
				"FIRST\n",
				"SECOND\n",
				"THIRD\n",
			},
		},
		{
			name: "InvalidGoto",
			program: `10 PRINT "START"
20 GOTO 999
30 PRINT "END"`,
			wantErr:     true,
			errContains: "UNDEFINED STATEMENT ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := executeBasicProgram(t, tt.program)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, output)
			}
		})
	}
}

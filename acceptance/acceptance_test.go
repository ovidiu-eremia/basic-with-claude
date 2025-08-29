// ABOUTME: End-to-end acceptance tests for the BASIC interpreter
// ABOUTME: Tests complete pipeline from inline program strings through execution and output verification

package acceptance

import (
	"fmt"
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
	maxSteps    int // Custom max steps limit, 0 means use default
}

// executeBasicProgram parses and executes a BASIC program string, returning the output
func executeBasicProgram(t *testing.T, program string) ([]string, error) {
	return executeBasicProgramWithMaxSteps(t, program, 0) // Use default max steps
}

// executeBasicProgramWithMaxSteps parses and executes a BASIC program string with custom max steps
func executeBasicProgramWithMaxSteps(t *testing.T, program string, maxSteps int) ([]string, error) {
	t.Helper()

	// Parse the program
	l := lexer.New(program)
	p := parser.New(l)
	ast := p.ParseProgram()

	// Check for parsing errors
	if len(p.Errors()) > 0 {
		return nil, fmt.Errorf("parse errors: %v", p.Errors())
	}
	if ast == nil {
		return nil, fmt.Errorf("parsing returned nil AST")
	}

	// Create test runtime and interpreter
	testRuntime := runtime.NewTestRuntime()
	interp := interpreter.NewInterpreter(testRuntime)

	// Set custom max steps if specified
	if maxSteps > 0 {
		interp.SetMaxSteps(maxSteps)
	}

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
		// ---
		// BUG-DERIVED TESTS
		// ---
		{
			name:        "Bug_DivisionByZeroErrorMessage",
			program:     `10 PRINT 1/0`,
			wantErr:     true,
			errContains: "?DIVISION BY ZERO ERROR IN 10",
		},
		{
			name:    "Bug_StringConcatenation",
			program: `10 PRINT "HELLO" + " " + "WORLD"`,
			expected: []string{
				"HELLO WORLD\n",
			},
		},
		{
			name:    "Bug_CaseInsensitiveKeywords",
			program: `10 print "hello"`,
			expected: []string{
				"hello\n",
			},
		},
		{
			name:    "Bug_PrintStatementWithoutArgs",
			program: `10 PRINT`,
			expected: []string{
				"\n",
			},
		},
		{
			name: "Bug_VariableNameLength",
			program: `10 VA = 5
20 VARA = 10
30 PRINT VA`,
			expected: []string{
				"10\n",
			},
		},
		{
			name:    "Bug_UninitializedStringVariable",
			program: `10 PRINT A$`,
			expected: []string{
				"\n",
			},
		},
		{
			name:        "Bug_TypeMismatchOnAssignment",
			program:     `10 A = "hello"`,
			wantErr:     true,
			errContains: "?TYPE MISMATCH ERROR IN 10",
		},
		{
			name: "Bug_TypeMismatchOnAddition",
			program: `10 A$ = "10"
20 B = 5
30 C = A$ + B`,
			wantErr:     true,
			errContains: "?TYPE MISMATCH ERROR IN 30",
		},
		{
			name: "Bug_ParserArithmetic",
			program: `10 PRINT -10 + 5
20 PRINT 2.5 * 2
30 PRINT 2 * 3 + 4 * 5`,
			expected: []string{
				"-5\n",
				"5\n",
				"26\n",
			},
		},
		{
			name:        "InfiniteLoopProtection_SimpleLoop",
			program:     `10 GOTO 10`,
			wantErr:     true,
			errContains: "?INFINITE LOOP ERROR",
		},
		{
			name: "InfiniteLoopProtection_ComplexLoop",
			program: `10 A = A + 1
20 GOTO 10`,
			wantErr:     true,
			errContains: "?INFINITE LOOP ERROR",
		},
		{
			name: "InfiniteLoopProtection_LongButFinite",
			program: `10 A = A + 1
20 PRINT A
30 STOP`,
			expected: []string{
				"1\n",
			},
		},
		{
			name: "InfiniteLoopProtection_BackwardJump",
			program: `10 B = B + 1
20 PRINT B
30 GOTO 10`,
			wantErr:     true,
			errContains: "?INFINITE LOOP ERROR",
		},
		{
			name: "InfiniteLoopProtection_NestedGotos",
			program: `10 GOTO 30
20 GOTO 10
30 GOTO 20`,
			wantErr:     true,
			errContains: "?INFINITE LOOP ERROR",
		},
		{
			name: "InfiniteLoopProtection_CustomMaxSteps",
			program: `10 A = A + 1
20 GOTO 10`,
			wantErr:     true,
			errContains: "?INFINITE LOOP ERROR",
			maxSteps:    5, // Custom low limit
		},
		// ---
		// STEP 9A: BASIC IF...THEN TESTS
		// ---
		{
			name: "BasicIFThen_TrueCondition",
			program: `10 IF 1 THEN PRINT "TRUE"
20 PRINT "DONE"`,
			expected: []string{
				"TRUE\n",
				"DONE\n",
			},
		},
		{
			name: "BasicIFThen_FalseCondition",
			program: `10 IF 0 THEN PRINT "FALSE"
20 PRINT "DONE"`,
			expected: []string{
				"DONE\n",
			},
		},
		{
			name: "BasicIFThen_NonZeroTrue",
			program: `10 A = 5
20 IF A THEN PRINT "NON-ZERO"
30 PRINT "FINISHED"`,
			expected: []string{
				"NON-ZERO\n",
				"FINISHED\n",
			},
		},
		{
			name: "BasicIFThen_ZeroFalse",
			program: `10 A = 0
20 IF A THEN PRINT "NON-ZERO"
30 PRINT "FINISHED"`,
			expected: []string{
				"FINISHED\n",
			},
		},
		{
			name: "BasicIFThen_NegativeTrue",
			program: `10 IF -1 THEN PRINT "NEGATIVE IS TRUE"
20 PRINT "END"`,
			expected: []string{
				"NEGATIVE IS TRUE\n",
				"END\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output []string
			var err error
			if tt.maxSteps > 0 {
				output, err = executeBasicProgramWithMaxSteps(t, tt.program, tt.maxSteps)
			} else {
				output, err = executeBasicProgram(t, tt.program)
			}

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

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
func executeBasicProgram(t *testing.T, program string, inputs []string) ([]string, error) {
	return executeBasicProgramWithMaxSteps(t, program, inputs, 0) // Use default max steps
}

// executeBasicProgramWithMaxSteps parses and executes a BASIC program string with custom max steps
func executeBasicProgramWithMaxSteps(t *testing.T, program string, inputs []string, maxSteps int) ([]string, error) {
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
	if len(inputs) > 0 {
		testRuntime.SetInput(inputs)
	}
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
		// STEP 10: INPUT Statement
		{
			name: "Input_Numeric",
			program: `10 INPUT A
20 PRINT A`,
			inputs:   []string{"5"},
			expected: []string{"5\n"},
		},
		{
			name: "Input_String",
			program: `10 INPUT A$
20 PRINT A$`,
			inputs:   []string{"HELLO"},
			expected: []string{"HELLO\n"},
		},
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
		// ---
		// STEP 9B: NUMERIC COMPARISONS TESTS
		// ---
		{
			name: "NumericComparisons_Equality",
			program: `10 A = 5
20 B = 10
30 C = 5
40 IF A = C THEN PRINT "A EQUALS C"
50 IF A = B THEN PRINT "A EQUALS B"
60 PRINT "DONE"`,
			expected: []string{
				"A EQUALS C\n",
				"DONE\n",
			},
		},
		{
			name: "NumericComparisons_Inequality",
			program: `10 A = 5
20 B = 10
30 IF A <> B THEN PRINT "A NOT EQUAL B"
40 IF A <> 5 THEN PRINT "NOT REACHED"
50 PRINT "DONE"`,
			expected: []string{
				"A NOT EQUAL B\n",
				"DONE\n",
			},
		},
		{
			name: "NumericComparisons_LessThan",
			program: `10 A = 5
20 B = 10
30 IF A < B THEN PRINT "A LESS THAN B"
40 IF B < A THEN PRINT "NOT REACHED"
50 PRINT "DONE"`,
			expected: []string{
				"A LESS THAN B\n",
				"DONE\n",
			},
		},
		{
			name: "NumericComparisons_GreaterThan",
			program: `10 A = 5
20 B = 10
30 IF B > A THEN PRINT "B GREATER THAN A"
40 IF A > B THEN PRINT "NOT REACHED"
50 PRINT "DONE"`,
			expected: []string{
				"B GREATER THAN A\n",
				"DONE\n",
			},
		},
		{
			name: "NumericComparisons_LessOrEqual",
			program: `10 A = 5
20 B = 5
30 C = 10
40 IF A <= B THEN PRINT "A LESS OR EQUAL B"
50 IF A <= C THEN PRINT "A LESS OR EQUAL C"
60 IF C <= A THEN PRINT "NOT REACHED"
70 PRINT "DONE"`,
			expected: []string{
				"A LESS OR EQUAL B\n",
				"A LESS OR EQUAL C\n",
				"DONE\n",
			},
		},
		{
			name: "NumericComparisons_GreaterOrEqual",
			program: `10 A = 10
20 B = 5
30 C = 10
40 IF A >= B THEN PRINT "A GREATER OR EQUAL B"
50 IF A >= C THEN PRINT "A GREATER OR EQUAL C"
60 IF B >= A THEN PRINT "NOT REACHED"
70 PRINT "DONE"`,
			expected: []string{
				"A GREATER OR EQUAL B\n",
				"A GREATER OR EQUAL C\n",
				"DONE\n",
			},
		},
		{
			name: "NumericComparisons_WithExpressions",
			program: `10 A = 5
20 B = 10
30 IF A + 5 = B THEN PRINT "EXPRESSION EQUALS"
40 IF A * 2 > B THEN PRINT "NOT REACHED"
50 IF A * 3 > B THEN PRINT "EXPRESSION GREATER"
60 PRINT "DONE"`,
			expected: []string{
				"EXPRESSION EQUALS\n",
				"EXPRESSION GREATER\n",
				"DONE\n",
			},
		},
		{
			name: "NumericComparisons_AllOperators",
			program: `10 A = 5
20 B = 10
30 C = 5
40 IF A = C THEN PRINT "A EQUALS C"
50 IF A <> B THEN PRINT "A NOT EQUAL B"
60 IF A < B THEN PRINT "A LESS THAN B"
70 IF B > A THEN PRINT "B GREATER THAN A"
80 IF A <= C THEN PRINT "A LESS OR EQUAL C"
90 IF B >= A THEN PRINT "B GREATER OR EQUAL A"
100 PRINT "ALL TESTS PASSED"`,
			expected: []string{
				"A EQUALS C\n",
				"A NOT EQUAL B\n",
				"A LESS THAN B\n",
				"B GREATER THAN A\n",
				"A LESS OR EQUAL C\n",
				"B GREATER OR EQUAL A\n",
				"ALL TESTS PASSED\n",
			},
		},
		// ---
		// STEP 9C: STRING COMPARISONS AND MIXED EXPRESSIONS TESTS
		// ---
		{
			name: "StringComparisons_Equality",
			program: `10 A$ = "HELLO"
20 B$ = "WORLD"
30 C$ = "HELLO"
40 IF A$ = C$ THEN PRINT "A EQUALS C"
50 IF A$ = B$ THEN PRINT "NOT REACHED"
60 PRINT "DONE"`,
			expected: []string{
				"A EQUALS C\n",
				"DONE\n",
			},
		},
		{
			name: "StringComparisons_Inequality",
			program: `10 A$ = "HELLO"
20 B$ = "WORLD"
30 IF A$ <> B$ THEN PRINT "A NOT EQUAL B"
40 IF A$ <> "HELLO" THEN PRINT "NOT REACHED"
50 PRINT "DONE"`,
			expected: []string{
				"A NOT EQUAL B\n",
				"DONE\n",
			},
		},
		{
			name: "StringComparisons_LexicographicLess",
			program: `10 A$ = "APPLE"
20 B$ = "BANANA"
30 IF A$ < B$ THEN PRINT "APPLE BEFORE BANANA"
40 IF B$ < A$ THEN PRINT "NOT REACHED"
50 PRINT "DONE"`,
			expected: []string{
				"APPLE BEFORE BANANA\n",
				"DONE\n",
			},
		},
		{
			name: "StringComparisons_LexicographicGreater",
			program: `10 A$ = "ZEBRA"
20 B$ = "APPLE"
30 IF A$ > B$ THEN PRINT "ZEBRA AFTER APPLE"
40 IF B$ > A$ THEN PRINT "NOT REACHED"
50 PRINT "DONE"`,
			expected: []string{
				"ZEBRA AFTER APPLE\n",
				"DONE\n",
			},
		},
		{
			name: "StringComparisons_LessOrEqual",
			program: `10 A$ = "APPLE"
20 B$ = "APPLE"
30 C$ = "BANANA"
40 IF A$ <= B$ THEN PRINT "A LESS OR EQUAL B"
50 IF A$ <= C$ THEN PRINT "A LESS OR EQUAL C"
60 IF C$ <= A$ THEN PRINT "NOT REACHED"
70 PRINT "DONE"`,
			expected: []string{
				"A LESS OR EQUAL B\n",
				"A LESS OR EQUAL C\n",
				"DONE\n",
			},
		},
		{
			name: "StringComparisons_GreaterOrEqual",
			program: `10 A$ = "ZEBRA"
20 B$ = "APPLE"
30 C$ = "ZEBRA"
40 IF A$ >= B$ THEN PRINT "A GREATER OR EQUAL B"
50 IF A$ >= C$ THEN PRINT "A GREATER OR EQUAL C"
60 IF B$ >= A$ THEN PRINT "NOT REACHED"
70 PRINT "DONE"`,
			expected: []string{
				"A GREATER OR EQUAL B\n",
				"A GREATER OR EQUAL C\n",
				"DONE\n",
			},
		},
		{
			name: "StringComparisons_WithLiterals",
			program: `10 NAME$ = "ALICE"
20 IF NAME$ = "ALICE" THEN PRINT "HELLO ALICE"
30 IF NAME$ <> "BOB" THEN PRINT "NOT BOB"
40 IF "ALICE" = NAME$ THEN PRINT "LITERAL EQUALS VARIABLE"
50 PRINT "DONE"`,
			expected: []string{
				"HELLO ALICE\n",
				"NOT BOB\n",
				"LITERAL EQUALS VARIABLE\n",
				"DONE\n",
			},
		},
		{
			name: "StringComparisons_AllOperators",
			program: `10 A$ = "APPLE"
20 B$ = "BANANA"
30 C$ = "APPLE"
40 IF A$ = C$ THEN PRINT "A EQUALS C"
50 IF A$ <> B$ THEN PRINT "A NOT EQUAL B"  
60 IF A$ < B$ THEN PRINT "A LESS THAN B"
70 IF B$ > A$ THEN PRINT "B GREATER THAN A"
80 IF A$ <= C$ THEN PRINT "A LESS OR EQUAL C"
90 IF B$ >= A$ THEN PRINT "B GREATER OR EQUAL A"
100 PRINT "ALL STRING TESTS PASSED"`,
			expected: []string{
				"A EQUALS C\n",
				"A NOT EQUAL B\n",
				"A LESS THAN B\n",
				"B GREATER THAN A\n",
				"A LESS OR EQUAL C\n",
				"B GREATER OR EQUAL A\n",
				"ALL STRING TESTS PASSED\n",
			},
		},
		{
			name: "MixedComparisons_TypeMismatch",
			program: `10 A = 5
20 B$ = "5"
30 IF A = B$ THEN PRINT "NOT REACHED"
40 PRINT "TYPE MISMATCH HANDLED"`,
			wantErr:     true,
			errContains: "?TYPE MISMATCH ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output []string
			var err error
			if tt.maxSteps > 0 {
				output, err = executeBasicProgramWithMaxSteps(t, tt.program, tt.inputs, tt.maxSteps)
			} else {
				output, err = executeBasicProgram(t, tt.program, tt.inputs)
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

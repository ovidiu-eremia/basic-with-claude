package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"basic-interpreter/lexer"
)

func TestParser_StatementParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *Program
	}{
		// PRINT statements
		{
			name:     "single line with PRINT",
			input:    `10 PRINT "HELLO"`,
			expected: program(line(10, 1, printStmt(str("HELLO", 1), 1))),
		},
		{
			name:  "multiple lines",
			input: "10 PRINT \"LINE1\"\n20 PRINT \"LINE2\"",
			expected: program(
				line(10, 1, printStmt(str("LINE1", 1), 1)),
				line(20, 2, printStmt(str("LINE2", 2), 2)),
			),
		},
		{
			name:     "empty string",
			input:    `10 PRINT ""`,
			expected: program(line(10, 1, printStmt(str("", 1), 1))),
		},
		{
			name:     "PRINT variable",
			input:    `10 PRINT A`,
			expected: program(line(10, 1, printStmt(varRef("A", 1), 1))),
		},
		{
			name:     "PRINT string variable",
			input:    `10 PRINT A$`,
			expected: program(line(10, 1, printStmt(varRef("A$", 1), 1))),
		},

		// LET and assignment statements
		{
			name:     "LET assignment",
			input:    `10 LET A = 42`,
			expected: program(line(10, 1, letStmt("A", num("42", 1), 1))),
		},
		{
			name:     "assignment without LET",
			input:    `10 X = 123`,
			expected: program(line(10, 1, letStmt("X", num("123", 1), 1))),
		},
		{
			name:     "string variable assignment with LET",
			input:    `10 LET A$ = "HELLO"`,
			expected: program(line(10, 1, letStmt("A$", str("HELLO", 1), 1))),
		},
		{
			name:     "string variable assignment without LET",
			input:    `10 NAME$ = "JOHN DOE"`,
			expected: program(line(10, 1, letStmt("NAME$", str("JOHN DOE", 1), 1))),
		},

		// END statement
		{
			name:     "END statement",
			input:    "10 END",
			expected: program(line(10, 1, endStmt(1))),
		},
		{
			name:  "program with END",
			input: "10 PRINT \"START\"\n20 END\n30 PRINT \"NEVER REACHED\"",
			expected: program(
				line(10, 1, printStmt(str("START", 1), 1)),
				line(20, 2, endStmt(2)),
				line(30, 3, printStmt(str("NEVER REACHED", 3), 3)),
			),
		},

		// RUN and STOP statements
		{
			name:     "RUN statement",
			input:    "10 RUN",
			expected: program(line(10, 1, runStmt(1))),
		},
		{
			name:     "STOP statement",
			input:    "10 STOP",
			expected: program(line(10, 1, stopStmt(1))),
		},
		{
			name:  "program with STOP",
			input: "10 PRINT \"START\"\n20 STOP\n30 PRINT \"NEVER\"",
			expected: program(
				line(10, 1, printStmt(str("START", 1), 1)),
				line(20, 2, stopStmt(2)),
				line(30, 3, printStmt(str("NEVER", 3), 3)),
			),
		},

		// GOTO statements
		{
			name:     "GOTO statement",
			input:    "10 GOTO 50",
			expected: program(line(10, 1, gotoStmt(50, 1))),
		},
		{
			name:  "program with GOTO",
			input: "10 PRINT \"BEFORE\"\n20 GOTO 50\n30 PRINT \"SKIPPED\"\n50 PRINT \"AFTER\"",
			expected: program(
				line(10, 1, printStmt(str("BEFORE", 1), 1)),
				line(20, 2, gotoStmt(50, 2)),
				line(30, 3, printStmt(str("SKIPPED", 3), 3)),
				line(50, 4, printStmt(str("AFTER", 4), 4)),
			),
		},

		// GOSUB and RETURN statements
		{
			name:     "GOSUB statement",
			input:    "10 GOSUB 100",
			expected: program(line(10, 1, gosubStmt(100, 1))),
		},
		{
			name:     "RETURN statement",
			input:    "20 RETURN",
			expected: program(line(20, 1, returnStmt(1))),
		},
		{
			name:  "program with GOSUB and RETURN",
			input: "10 PRINT \"MAIN\"\n20 GOSUB 100\n30 PRINT \"BACK\"\n40 END\n100 PRINT \"SUB\"\n110 RETURN",
			expected: program(
				line(10, 1, printStmt(str("MAIN", 1), 1)),
				line(20, 2, gosubStmt(100, 2)),
				line(30, 3, printStmt(str("BACK", 3), 3)),
				line(40, 4, endStmt(4)),
				line(100, 5, printStmt(str("SUB", 5), 5)),
				line(110, 6, returnStmt(6)),
			),
		},

		// IF statements
		{
			name:  "simple IF THEN",
			input: "10 IF 1 THEN PRINT \"TRUE\"",
			expected: program(
				line(10, 1, ifStmt(num("1", 1), printStmt(str("TRUE", 1), 1), 1)),
			),
		},

		// INPUT statements
		{
			name:     "INPUT without prompt",
			input:    "10 INPUT A",
			expected: program(line(10, 1, inputStmt("", "A", 1))),
		},
		{
			name:     "INPUT with prompt",
			input:    "10 INPUT \"ENTER A NUMBER\"; N",
			expected: program(line(10, 1, inputStmt("ENTER A NUMBER", "N", 1))),
		},
		{
			name:     "INPUT with string variable and prompt",
			input:    "10 INPUT \"WHAT IS YOUR NAME\"; NAME$",
			expected: program(line(10, 1, inputStmt("WHAT IS YOUR NAME", "NAME$", 1))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			program := p.ParseProgram()

			require.NotNil(t, program, "ParseProgram() returned nil")
			require.Nil(t, p.ParseError(), "Parser error: %v", p.ParseError())

			assert.Equal(t, tt.expected, program)
		})
	}
}

func TestParser_ParseErrors(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "unterminated string",
			input:       `10 PRINT "HELLO`,
			expectError: true,
		},
		{
			name:        "missing line number",
			input:       `PRINT "HELLO"`,
			expectError: true,
		},
		{
			name:        "invalid syntax",
			input:       `10 INVALID "HELLO"`,
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
				assert.Nil(t, p.ParseError(), "Expected no parsing error but got: %v", p.ParseError())
				assert.NotNil(t, program)
			}
		})
	}
}

func TestParser_ArithmeticExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Expression
	}{
		{
			name:     "simple addition",
			input:    "2 + 3",
			expected: binaryOp(num("2", 1), "+", num("3", 1), 1),
		},
		{
			name:     "precedence: multiplication over addition",
			input:    "2 + 3 * 4",
			expected: binaryOp(num("2", 1), "+", binaryOp(num("3", 1), "*", num("4", 1), 1), 1),
		},
		{
			name:     "parentheses override precedence",
			input:    "(2 + 3) * 4",
			expected: binaryOp(binaryOp(num("2", 1), "+", num("3", 1), 1), "*", num("4", 1), 1),
		},
		{
			name:     "variables in expressions",
			input:    "A + B",
			expected: binaryOp(varRef("A", 1), "+", varRef("B", 1), 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			expr := p.parseExpression()

			require.Nil(t, p.ParseError(), "Parser error: %v", p.ParseError())
			assert.Equal(t, tt.expected, expr)
		})
	}
}

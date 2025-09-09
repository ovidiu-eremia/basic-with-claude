package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// assertToken is a helper function to reduce test boilerplate
func assertToken(t *testing.T, expected, actual Token, index int) {
	t.Helper()
	assert.Equal(t, expected.Type, actual.Type, "Token %d type mismatch", index)
	assert.Equal(t, expected.Literal, actual.Literal, "Token %d literal mismatch", index)
}

func TestLexer_NextToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "line number",
			input: "10",
			expected: []Token{
				{Type: NUMBER, Literal: "10"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "print keyword",
			input: "PRINT",
			expected: []Token{
				{Type: PRINT, Literal: "PRINT"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "string literal",
			input: `"HELLO WORLD"`,
			expected: []Token{
				{Type: STRING, Literal: "HELLO WORLD"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "arithmetic operators",
			input: "+ - * / ^ ( )",
			expected: []Token{
				{Type: PLUS, Literal: "+"},
				{Type: MINUS, Literal: "-"},
				{Type: MULTIPLY, Literal: "*"},
				{Type: DIVIDE, Literal: "/"},
				{Type: POWER, Literal: "^"},
				{Type: LPAREN, Literal: "("},
				{Type: RPAREN, Literal: ")"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "basic program",
			input: `10 PRINT "HELLO"`,
			expected: []Token{
				{Type: NUMBER, Literal: "10"},
				{Type: PRINT, Literal: "PRINT"},
				{Type: STRING, Literal: "HELLO"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "whitespace handling",
			input: `10  PRINT   "HELLO"`,
			expected: []Token{
				{Type: NUMBER, Literal: "10"},
				{Type: PRINT, Literal: "PRINT"},
				{Type: STRING, Literal: "HELLO"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "multiple lines",
			input: "10 PRINT \"LINE1\"\n20 PRINT \"LINE2\"",
			expected: []Token{
				{Type: NUMBER, Literal: "10"},
				{Type: PRINT, Literal: "PRINT"},
				{Type: STRING, Literal: "LINE1"},
				{Type: NEWLINE, Literal: "\n"},
				{Type: NUMBER, Literal: "20"},
				{Type: PRINT, Literal: "PRINT"},
				{Type: STRING, Literal: "LINE2"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "unterminated string",
			input: `10 PRINT "HELLO`,
			expected: []Token{
				{Type: NUMBER, Literal: "10"},
				{Type: PRINT, Literal: "PRINT"},
				{Type: ILLEGAL, Literal: "unterminated string"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "empty string",
			input: `""`,
			expected: []Token{
				{Type: STRING, Literal: ""},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "let assignment",
			input: `10 LET A = 42`,
			expected: []Token{
				{Type: NUMBER, Literal: "10"},
				{Type: LET, Literal: "LET"},
				{Type: IDENT, Literal: "A"},
				{Type: ASSIGN, Literal: "="},
				{Type: NUMBER, Literal: "42"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "assignment without LET",
			input: `10 X = 123`,
			expected: []Token{
				{Type: NUMBER, Literal: "10"},
				{Type: IDENT, Literal: "X"},
				{Type: ASSIGN, Literal: "="},
				{Type: NUMBER, Literal: "123"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "variable names with digits",
			input: `A1 = 5`,
			expected: []Token{
				{Type: IDENT, Literal: "A1"},
				{Type: ASSIGN, Literal: "="},
				{Type: NUMBER, Literal: "5"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "string variable names",
			input: `A$ = "HELLO"`,
			expected: []Token{
				{Type: IDENT, Literal: "A$"},
				{Type: ASSIGN, Literal: "="},
				{Type: STRING, Literal: "HELLO"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "string variable with LET",
			input: `10 LET NAME$ = "JOHN"`,
			expected: []Token{
				{Type: NUMBER, Literal: "10"},
				{Type: LET, Literal: "LET"},
				{Type: IDENT, Literal: "NAME$"},
				{Type: ASSIGN, Literal: "="},
				{Type: STRING, Literal: "JOHN"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "RUN keyword",
			input: "RUN",
			expected: []Token{
				{Type: RUN, Literal: "RUN"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "STOP keyword",
			input: "STOP",
			expected: []Token{
				{Type: STOP, Literal: "STOP"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "GOTO keyword",
			input: "GOTO",
			expected: []Token{
				{Type: GOTO, Literal: "GOTO"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "GOTO with line number",
			input: "GOTO 100",
			expected: []Token{
				{Type: GOTO, Literal: "GOTO"},
				{Type: NUMBER, Literal: "100"},
				{Type: EOF, Literal: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := New(tt.input)

			for i, expectedToken := range tt.expected {
				token := lexer.NextToken()
				assertToken(t, expectedToken, token, i)
			}
		})
	}
}

func TestLexer_GosubReturn(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "GOSUB statement",
			input: "GOSUB 100",
			expected: []Token{
				{Type: GOSUB, Literal: "GOSUB"},
				{Type: NUMBER, Literal: "100"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "RETURN statement",
			input: "RETURN",
			expected: []Token{
				{Type: RETURN, Literal: "RETURN"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "GOSUB RETURN program",
			input: "10 GOSUB 100\n20 PRINT \"BACK\"\n100 RETURN",
			expected: []Token{
				{Type: NUMBER, Literal: "10"},
				{Type: GOSUB, Literal: "GOSUB"},
				{Type: NUMBER, Literal: "100"},
				{Type: NEWLINE, Literal: "\n"},
				{Type: NUMBER, Literal: "20"},
				{Type: PRINT, Literal: "PRINT"},
				{Type: STRING, Literal: "BACK"},
				{Type: NEWLINE, Literal: "\n"},
				{Type: NUMBER, Literal: "100"},
				{Type: RETURN, Literal: "RETURN"},
				{Type: EOF, Literal: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := New(tt.input)

			for i, expectedToken := range tt.expected {
				token := lexer.NextToken()
				assertToken(t, expectedToken, token, i)
			}
		})
	}
}

func TestLexer_ComparisonOperators(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "all comparison operators",
			input: `< > <= >= <> =`,
			expected: []Token{
				{Type: LT, Literal: "<"},
				{Type: GT, Literal: ">"},
				{Type: LE, Literal: "<="},
				{Type: GE, Literal: ">="},
				{Type: NE, Literal: "<>"},
				{Type: ASSIGN, Literal: "="},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "IF THEN statement",
			input: `IF A > 5 THEN PRINT "BIG"`,
			expected: []Token{
				{Type: IF, Literal: "IF"},
				{Type: IDENT, Literal: "A"},
				{Type: GT, Literal: ">"},
				{Type: NUMBER, Literal: "5"},
				{Type: THEN, Literal: "THEN"},
				{Type: PRINT, Literal: "PRINT"},
				{Type: STRING, Literal: "BIG"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "IF THEN with not equal",
			input: `IF A <> 0 THEN GOTO 100`,
			expected: []Token{
				{Type: IF, Literal: "IF"},
				{Type: IDENT, Literal: "A"},
				{Type: NE, Literal: "<>"},
				{Type: NUMBER, Literal: "0"},
				{Type: THEN, Literal: "THEN"},
				{Type: GOTO, Literal: "GOTO"},
				{Type: NUMBER, Literal: "100"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "IF THEN with less than or equal",
			input: `IF X <= 10 THEN Y = 1`,
			expected: []Token{
				{Type: IF, Literal: "IF"},
				{Type: IDENT, Literal: "X"},
				{Type: LE, Literal: "<="},
				{Type: NUMBER, Literal: "10"},
				{Type: THEN, Literal: "THEN"},
				{Type: IDENT, Literal: "Y"},
				{Type: ASSIGN, Literal: "="},
				{Type: NUMBER, Literal: "1"},
				{Type: EOF, Literal: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := New(tt.input)

			for i, expectedToken := range tt.expected {
				token := lexer.NextToken()
				assertToken(t, expectedToken, token, i)
			}
		})
	}
}

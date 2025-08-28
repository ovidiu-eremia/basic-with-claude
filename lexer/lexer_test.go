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
	assert.Equal(t, expected.Line, actual.Line, "Token %d line mismatch", index)
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
				{Type: NUMBER, Literal: "10", Line: 1},
				{Type: EOF, Literal: "", Line: 1},
			},
		},
		{
			name:  "print keyword",
			input: "PRINT",
			expected: []Token{
				{Type: PRINT, Literal: "PRINT", Line: 1},
				{Type: EOF, Literal: "", Line: 1},
			},
		},
		{
			name:  "string literal",
			input: `"HELLO WORLD"`,
			expected: []Token{
				{Type: STRING, Literal: "HELLO WORLD", Line: 1},
				{Type: EOF, Literal: "", Line: 1},
			},
		},
		{
			name:  "basic program",
			input: `10 PRINT "HELLO"`,
			expected: []Token{
				{Type: NUMBER, Literal: "10", Line: 1},
				{Type: PRINT, Literal: "PRINT", Line: 1},
				{Type: STRING, Literal: "HELLO", Line: 1},
				{Type: EOF, Literal: "", Line: 1},
			},
		},
		{
			name:  "whitespace handling",
			input: `10  PRINT   "HELLO"`,
			expected: []Token{
				{Type: NUMBER, Literal: "10", Line: 1},
				{Type: PRINT, Literal: "PRINT", Line: 1},
				{Type: STRING, Literal: "HELLO", Line: 1},
				{Type: EOF, Literal: "", Line: 1},
			},
		},
		{
			name:  "multiple lines",
			input: "10 PRINT \"LINE1\"\n20 PRINT \"LINE2\"",
			expected: []Token{
				{Type: NUMBER, Literal: "10", Line: 1},
				{Type: PRINT, Literal: "PRINT", Line: 1},
				{Type: STRING, Literal: "LINE1", Line: 1},
				{Type: NEWLINE, Literal: "\n", Line: 1},
				{Type: NUMBER, Literal: "20", Line: 2},
				{Type: PRINT, Literal: "PRINT", Line: 2},
				{Type: STRING, Literal: "LINE2", Line: 2},
				{Type: EOF, Literal: "", Line: 2},
			},
		},
		{
			name:  "unterminated string",
			input: `10 PRINT "HELLO`,
			expected: []Token{
				{Type: NUMBER, Literal: "10", Line: 1},
				{Type: PRINT, Literal: "PRINT", Line: 1},
				{Type: ILLEGAL, Literal: "unterminated string", Line: 1},
				{Type: EOF, Literal: "", Line: 1},
			},
		},
		{
			name:  "empty string",
			input: `""`,
			expected: []Token{
				{Type: STRING, Literal: "", Line: 1},
				{Type: EOF, Literal: "", Line: 1},
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
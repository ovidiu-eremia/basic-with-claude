package lexer

import "testing"

func TestLexer_RemAndColon(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "REM keyword token",
			input: "REM",
			expected: []Token{
				{Type: REM, Literal: "REM", Line: 1},
				{Type: EOF, Literal: "", Line: 1},
			},
		},
		{
			name:  "colon separator",
			input: ":",
			expected: []Token{
				{Type: COLON, Literal: ":", Line: 1},
				{Type: EOF, Literal: "", Line: 1},
			},
		},
		{
			name:  "print then colon then print",
			input: "PRINT \"A\" : PRINT \"B\"",
			expected: []Token{
				{Type: PRINT, Literal: "PRINT", Line: 1},
				{Type: STRING, Literal: "A", Line: 1},
				{Type: COLON, Literal: ":", Line: 1},
				{Type: PRINT, Literal: "PRINT", Line: 1},
				{Type: STRING, Literal: "B", Line: 1},
				{Type: EOF, Literal: "", Line: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			for i, exp := range tt.expected {
				tok := l.NextToken()
				assertToken(t, exp, tok, i)
			}
		})
	}
}

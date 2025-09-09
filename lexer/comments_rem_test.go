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
				{Type: REM, Literal: "REM"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "colon separator",
			input: ":",
			expected: []Token{
				{Type: COLON, Literal: ":"},
				{Type: EOF, Literal: ""},
			},
		},
		{
			name:  "print then colon then print",
			input: "PRINT \"A\" : PRINT \"B\"",
			expected: []Token{
				{Type: PRINT, Literal: "PRINT"},
				{Type: STRING, Literal: "A"},
				{Type: COLON, Literal: ":"},
				{Type: PRINT, Literal: "PRINT"},
				{Type: STRING, Literal: "B"},
				{Type: EOF, Literal: ""},
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

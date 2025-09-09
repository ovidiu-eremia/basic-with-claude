package lexer

import "testing"

func TestLexer_DataReadAndComma(t *testing.T) {
	input := "DATA 10, \"HELLO\"\nREAD A, B$"
	l := New(input)

	tokens := []Token{
		{Type: DATA, Literal: "DATA"},
		{Type: NUMBER, Literal: "10"},
		{Type: COMMA, Literal: ","},
		{Type: STRING, Literal: "HELLO"},
		{Type: NEWLINE, Literal: "\n"},
		{Type: READ, Literal: "READ"},
		{Type: IDENT, Literal: "A"},
		{Type: COMMA, Literal: ","},
		{Type: IDENT, Literal: "B$"},
		{Type: EOF, Literal: ""},
	}

	for i := range tokens {
		tok := l.NextToken()
		if tok != tokens[i] {
			t.Fatalf("unexpected token %d: got %#v want %#v", i, tok, tokens[i])
		}
	}
}

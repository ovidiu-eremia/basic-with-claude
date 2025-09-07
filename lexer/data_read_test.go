package lexer

import "testing"

func TestLexer_DataReadAndComma(t *testing.T) {
	input := "DATA 10, \"HELLO\"\nREAD A, B$"
	l := New(input)

	tokens := []Token{
		{Type: DATA, Literal: "DATA", Line: 1},
		{Type: NUMBER, Literal: "10", Line: 1},
		{Type: COMMA, Literal: ",", Line: 1},
		{Type: STRING, Literal: "HELLO", Line: 1},
		{Type: NEWLINE, Literal: "\n", Line: 1},
		{Type: READ, Literal: "READ", Line: 2},
		{Type: IDENT, Literal: "A", Line: 2},
		{Type: COMMA, Literal: ",", Line: 2},
		{Type: IDENT, Literal: "B$", Line: 2},
		{Type: EOF, Literal: "", Line: 2},
	}

	for i := range tokens {
		tok := l.NextToken()
		if tok != tokens[i] {
			t.Fatalf("unexpected token %d: got %#v want %#v", i, tok, tokens[i])
		}
	}
}

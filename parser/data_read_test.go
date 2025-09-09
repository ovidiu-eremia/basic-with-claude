package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"basic-interpreter/lexer"
)

func TestParser_DataAndReadStatements(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "DATA with constants and READ with variables",
			input: "10 DATA 10, \"HELLO\"\n20 READ A, B$",
		},
		{
			name:  "READ before DATA (order independent)",
			input: "10 READ X, Y$\n20 DATA 1, \"S\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			prog := p.ParseProgram()
			require.NotNil(t, prog)
			require.Nil(t, p.ParseError())
		})
	}
}

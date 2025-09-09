package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"basic-interpreter/lexer"
)

func TestParser_ColonSeparatesStatements(t *testing.T) {
	input := "10 PRINT \"A\": PRINT \"B\""
	l := lexer.New(input)
	p := New(l)
	got := p.ParseProgram()
	require.Nil(t, p.ParseError())

	expected := program(
		line(10, 1,
			printStmt(str("A", 1), 1),
			printStmt(str("B", 1), 1),
		),
	)
	require.Equal(t, expected, got)
}

func TestParser_RemSkipsRestOfLine(t *testing.T) {
	input := "10 PRINT \"A\": REM ignore this: PRINT \"X\"\n20 PRINT \"B\""
	l := lexer.New(input)
	p := New(l)
	got := p.ParseProgram()
	require.Nil(t, p.ParseError())

	expected := program(
		line(10, 1,
			printStmt(str("A", 1), 1),
			remStmt(1),
		),
		line(20, 2,
			printStmt(str("B", 2), 2),
		),
	)
	require.Equal(t, expected, got)
}

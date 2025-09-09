package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"basic-interpreter/lexer"
)

func TestParseOnGoto(t *testing.T) {
	input := "10 ON X GOTO 100,200,300\n20 END\n"
	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()
	require.NotNil(t, prog)
	require.Nil(t, p.ParseError())
	require.Len(t, prog.Lines, 2)
	require.Len(t, prog.Lines[0].Statements, 1)
	st, ok := prog.Lines[0].Statements[0].(*OnGotoStatement)
	require.True(t, ok)
	require.Len(t, st.TargetLines, 3)
	assert.Equal(t, 200, st.TargetLines[1])
}

func TestParseOnGosub(t *testing.T) {
	input := "10 ON 2 GOSUB 100,200\n20 END\n"
	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()
	require.NotNil(t, prog)
	require.Nil(t, p.ParseError())
	require.Len(t, prog.Lines, 2)
	require.Len(t, prog.Lines[0].Statements, 1)
	st, ok := prog.Lines[0].Statements[0].(*OnGosubStatement)
	require.True(t, ok)
	require.Len(t, st.TargetLines, 2)
	assert.Equal(t, 100, st.TargetLines[0])
}

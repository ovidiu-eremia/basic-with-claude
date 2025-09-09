package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"basic-interpreter/lexer"
)

func TestParseReadArrayTargets(t *testing.T) {
	input := "10 READ S(1,2), X, S(J,K)\n20 END\n"
	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()
	require.NotNil(t, prog)
	require.Nil(t, p.ParseError())
	rs, ok := prog.Lines[0].Statements[0].(*ReadStatement)
	require.True(t, ok)
	require.Len(t, rs.Targets, 3)
	require.Len(t, rs.Targets[0].Indices, 2)
	require.Len(t, rs.Targets[1].Indices, 0)
	require.Len(t, rs.Targets[2].Indices, 2)
}

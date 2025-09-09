package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"basic-interpreter/lexer"
)

func TestParseDimWithTwoDimensions(t *testing.T) {
	input := "10 DIM S(20,3)\n20 END\n"
	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()
	require.NotNil(t, prog)
	require.Nil(t, p.ParseError())
	require.Len(t, prog.Lines, 2)

	// Line 10 should be a DIM statement with one declaration and two sizes
	require.Len(t, prog.Lines[0].Statements, 1)
	ds, ok := prog.Lines[0].Statements[0].(*DimStatement)
	require.True(t, ok)
	require.Len(t, ds.Declarations, 1)
	decl := ds.Declarations[0]
	assert.Equal(t, "S", decl.Name)
	require.Len(t, decl.Sizes, 2)
}

func TestParseArraySetTwoIndices(t *testing.T) {
	input := "10 S(1,2)=7\n20 END\n"
	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()
	require.NotNil(t, prog)
	require.Nil(t, p.ParseError())

	require.Len(t, prog.Lines, 2)
	require.Len(t, prog.Lines[0].Statements, 1)

	as, ok := prog.Lines[0].Statements[0].(*ArraySetStatement)
	require.True(t, ok)
	assert.Equal(t, "S", as.Name)
	require.Len(t, as.Indexes, 2)
}

func TestParseArrayRefTwoIndicesInPrint(t *testing.T) {
	input := "10 DIM S(2,3)\n20 PRINT S(1,2)\n30 END\n"
	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()
	require.NotNil(t, prog)
	require.Nil(t, p.ParseError())

	require.Len(t, prog.Lines, 3)
	require.Len(t, prog.Lines[1].Statements, 1)
	ps, ok := prog.Lines[1].Statements[0].(*PrintStatement)
	require.True(t, ok)
	// PRINT uses Expression or Items; handle either
	var expr Expression
	if len(ps.Items) > 0 {
		expr = ps.Items[0]
	} else {
		expr = ps.Expression
	}
	ar, ok := expr.(*ArrayReference)
	require.True(t, ok)
	assert.Equal(t, "S", ar.Name)
	require.Len(t, ar.Indices, 2)
}

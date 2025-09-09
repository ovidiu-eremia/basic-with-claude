package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"basic-interpreter/lexer"
	"basic-interpreter/parser"
	"basic-interpreter/runtime"
)

func TestInterpreter_ReadDataBasic(t *testing.T) {
	program := "" +
		"10 READ A, B, C$\n" +
		"20 PRINT A\n" +
		"25 PRINT B\n" +
		"28 PRINT C$\n" +
		"30 DATA 10, 20, \"HELLO\"\n" +
		"40 END\n"

	l := lexer.New(program)
	p := parser.New(l)
	ast := p.ParseProgram()
	require.Nil(t, p.ParseError())

	rt := runtime.NewTestRuntime()
	interp := NewInterpreter(rt)
	err := interp.Execute(ast)
	require.NoError(t, err)

	assert.Equal(t, []string{"10\n", "20\n", "HELLO\n"}, rt.GetOutput())
}

func TestInterpreter_ReadAcrossMultipleDataStatements(t *testing.T) {
	src := "" +
		"10 READ X, Y, Z, M$\n" +
		"20 PRINT X\n" +
		"30 PRINT Y\n" +
		"40 PRINT Z\n" +
		"50 PRINT M$\n" +
		"60 DATA 1, 2\n" +
		"70 DATA 3, \"DONE\"\n" +
		"80 END\n"

	l := lexer.New(src)
	p := parser.New(l)
	ast := p.ParseProgram()
	require.Nil(t, p.ParseError())

	rt := runtime.NewTestRuntime()
	interp := NewInterpreter(rt)
	err := interp.Execute(ast)
	require.NoError(t, err)
	assert.Equal(t, []string{"1\n", "2\n", "3\n", "DONE\n"}, rt.GetOutput())
}

func TestInterpreter_ReadOutOfData(t *testing.T) {
	src := "" +
		"10 READ A, B, C\n" +
		"20 DATA 5, 6\n"

	l := lexer.New(src)
	p := parser.New(l)
	ast := p.ParseProgram()
	require.Nil(t, p.ParseError())

	rt := runtime.NewTestRuntime()
	interp := NewInterpreter(rt)
	err := interp.Execute(ast)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "?OUT OF DATA ERROR")
}

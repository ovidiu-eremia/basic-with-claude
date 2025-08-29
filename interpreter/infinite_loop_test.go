package interpreter

import (
	"testing"

	"basic-interpreter/lexer"
	"basic-interpreter/parser"
	"basic-interpreter/runtime"
)

func TestInfiniteLoopProtection(t *testing.T) {
	t.Run("infinite loop with default protection", func(t *testing.T) {
		program := `10 GOTO 10`
		l := lexer.New(program)
		p := parser.New(l)
		ast := p.ParseProgram()

		testRuntime := runtime.NewTestRuntime()
		interp := NewInterpreter(testRuntime)

		err := interp.Execute(ast)
		if err == nil {
			t.Error("Expected infinite loop error but got nil")
		}
		if err.Error() != "?INFINITE LOOP ERROR" {
			t.Errorf("Expected '?INFINITE LOOP ERROR' but got '%s'", err.Error())
		}
	})

	t.Run("custom max steps limit", func(t *testing.T) {
		program := `10 A = A + 1
20 PRINT A
30 GOTO 10`
		l := lexer.New(program)
		p := parser.New(l)
		ast := p.ParseProgram()

		testRuntime := runtime.NewTestRuntime()
		interp := NewInterpreter(testRuntime)
		interp.SetMaxSteps(3) // Very low limit

		err := interp.Execute(ast)
		if err == nil {
			t.Error("Expected infinite loop error but got nil")
		}
		if err.Error() != "?INFINITE LOOP ERROR" {
			t.Errorf("Expected '?INFINITE LOOP ERROR' but got '%s'", err.Error())
		}
	})

	t.Run("finite program within limit", func(t *testing.T) {
		program := `10 A = A + 1
20 PRINT A
30 END`
		l := lexer.New(program)
		p := parser.New(l)
		ast := p.ParseProgram()

		testRuntime := runtime.NewTestRuntime()
		interp := NewInterpreter(testRuntime)

		err := interp.Execute(ast)
		if err != nil {
			t.Errorf("Expected no error but got '%s'", err.Error())
		}

		output := testRuntime.GetOutput()
		if len(output) != 1 || output[0] != "1\n" {
			t.Errorf("Expected output ['1\\n'] but got %v", output)
		}
	})
}

// ABOUTME: Tree-walking interpreter for BASIC AST execution and runtime state management
// ABOUTME: Executes parsed BASIC programs by walking the AST and managing program state

package interpreter

import (
	"basic-interpreter/parser"
	"basic-interpreter/runtime"
)

// Interpreter executes BASIC programs by walking the AST
type Interpreter struct {
	runtime runtime.Runtime
}

// New creates a new interpreter instance
func New(rt runtime.Runtime) *Interpreter {
	return &Interpreter{
		runtime: rt,
	}
}

// Execute runs a BASIC program
func (i *Interpreter) Execute(program *parser.Program) error {
	// Execute each line in sequence
	for _, line := range program.Lines {
		for _, stmt := range line.Statements {
			err := i.executeStatement(stmt)
			if err != nil {
				return err
			}
			
			// Check if this is an END statement - stop execution
			if _, isEnd := stmt.(*parser.EndStatement); isEnd {
				return nil
			}
		}
	}
	return nil
}

// executeStatement executes a single statement
func (i *Interpreter) executeStatement(stmt parser.Statement) error {
	switch s := stmt.(type) {
	case *parser.PrintStatement:
		return i.executePrintStatement(s)
	case *parser.EndStatement:
		// END statement - just return, handled in Execute
		return nil
	default:
		// For now, ignore unknown statement types
		return nil
	}
}

// executePrintStatement executes a PRINT statement
func (i *Interpreter) executePrintStatement(stmt *parser.PrintStatement) error {
	value, err := i.evaluateExpression(stmt.Expression)
	if err != nil {
		return err
	}
	
	return i.runtime.PrintLine(value)
}

// evaluateExpression evaluates an expression and returns its string value
func (i *Interpreter) evaluateExpression(expr parser.Expression) (string, error) {
	switch e := expr.(type) {
	case *parser.StringLiteral:
		return e.Value, nil
	default:
		return "", nil
	}
}
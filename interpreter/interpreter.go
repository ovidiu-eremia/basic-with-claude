// ABOUTME: Tree-walking interpreter for BASIC AST execution and runtime state management
// ABOUTME: Executes parsed BASIC programs by walking the AST and managing program state

package interpreter

import (
	"fmt"

	"basic-interpreter/parser"
	"basic-interpreter/runtime"
)

// Interpreter executes BASIC programs by walking the AST
type Interpreter struct {
	runtime   runtime.Runtime
	variables map[string]Value // Variable storage using proper Value types
}

// NewInterpreter creates a new interpreter instance
func NewInterpreter(rt runtime.Runtime) *Interpreter {
	return &Interpreter{
		runtime:   rt,
		variables: make(map[string]Value),
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
	case *parser.LetStatement:
		return i.executeLetStatement(s)
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
	
	return i.runtime.PrintLine(value.ToString())
}

// evaluateExpression evaluates an expression and returns its Value
func (i *Interpreter) evaluateExpression(expr parser.Expression) (Value, error) {
	switch e := expr.(type) {
	case *parser.StringLiteral:
		return NewStringValue(e.Value), nil
	case *parser.NumberLiteral:
		val, err := ParseValue(e.Value)
		if err != nil {
			return Value{}, err
		}
		return val, nil
	case *parser.VariableReference:
		if value, exists := i.variables[e.Name]; exists {
			return value, nil
		}
		return NewNumberValue(0), nil // Default value for uninitialized variables
	case *parser.BinaryOperation:
		return i.evaluateBinaryOperation(e)
	default:
		return Value{}, fmt.Errorf("unknown expression type")
	}
}

// evaluateBinaryOperation evaluates a binary arithmetic operation
func (i *Interpreter) evaluateBinaryOperation(expr *parser.BinaryOperation) (Value, error) {
	left, err := i.evaluateExpression(expr.Left)
	if err != nil {
		return Value{}, err
	}
	
	right, err := i.evaluateExpression(expr.Right)
	if err != nil {
		return Value{}, err
	}
	
	switch expr.Operator {
	case "+":
		return left.Add(right)
	case "-":
		return left.Subtract(right)
	case "*":
		return left.Multiply(right)
	case "/":
		return left.Divide(right)
	case "^":
		return left.Power(right)
	default:
		return Value{}, fmt.Errorf("unknown operator: %s", expr.Operator)
	}
}

// executeLetStatement executes a LET statement (variable assignment)
func (i *Interpreter) executeLetStatement(stmt *parser.LetStatement) error {
	value, err := i.evaluateExpression(stmt.Expression)
	if err != nil {
		return err
	}
	
	i.variables[stmt.Variable] = value
	return nil
}
// ABOUTME: Tree-walking interpreter for BASIC AST execution and runtime state management
// ABOUTME: Executes parsed BASIC programs by walking the AST and managing program state

package interpreter

import (
	"fmt"
	"math"
	"strconv"
	
	"basic-interpreter/parser"
	"basic-interpreter/runtime"
)

// Interpreter executes BASIC programs by walking the AST
type Interpreter struct {
	runtime   runtime.Runtime
	variables map[string]string // Variable storage
}

// NewInterpreter creates a new interpreter instance
func NewInterpreter(rt runtime.Runtime) *Interpreter {
	return &Interpreter{
		runtime:   rt,
		variables: make(map[string]string),
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
	
	return i.runtime.PrintLine(value)
}

// evaluateExpression evaluates an expression and returns its string value
func (i *Interpreter) evaluateExpression(expr parser.Expression) (string, error) {
	switch e := expr.(type) {
	case *parser.StringLiteral:
		return e.Value, nil
	case *parser.NumberLiteral:
		return e.Value, nil
	case *parser.VariableReference:
		if value, exists := i.variables[e.Name]; exists {
			return value, nil
		}
		return "0", nil // Default value for uninitialized variables
	case *parser.BinaryOperation:
		return i.evaluateBinaryOperation(e)
	default:
		return "", fmt.Errorf("unknown expression type")
	}
}

// evaluateBinaryOperation evaluates a binary arithmetic operation
func (i *Interpreter) evaluateBinaryOperation(expr *parser.BinaryOperation) (string, error) {
	leftStr, err := i.evaluateExpression(expr.Left)
	if err != nil {
		return "", err
	}
	
	rightStr, err := i.evaluateExpression(expr.Right)
	if err != nil {
		return "", err
	}
	
	// Convert to numbers for arithmetic
	left, err := strconv.ParseFloat(leftStr, 64)
	if err != nil {
		return "", fmt.Errorf("invalid number: %s", leftStr)
	}
	
	right, err := strconv.ParseFloat(rightStr, 64)
	if err != nil {
		return "", fmt.Errorf("invalid number: %s", rightStr)
	}
	
	var result float64
	switch expr.Operator {
	case "+":
		result = left + right
	case "-":
		result = left - right
	case "*":
		result = left * right
	case "/":
		if right == 0 {
			return "", fmt.Errorf("division by zero")
		}
		result = left / right
	case "^":
		result = math.Pow(left, right)
	default:
		return "", fmt.Errorf("unknown operator: %s", expr.Operator)
	}
	
	// Convert result back to string
	if result == float64(int64(result)) {
		// If it's a whole number, return as integer
		return strconv.FormatInt(int64(result), 10), nil
	}
	return strconv.FormatFloat(result, 'g', -1, 64), nil
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
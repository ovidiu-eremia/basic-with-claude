// ABOUTME: Tree-walking interpreter for BASIC AST execution and runtime state management
// ABOUTME: Executes parsed BASIC programs by walking the AST and managing program state

package interpreter

import (
	"fmt"
	"strings"

	"basic-interpreter/lexer"
	"basic-interpreter/parser"
	"basic-interpreter/runtime"
)

// RuntimeError represents an error that occurred during program execution
type RuntimeError struct {
	Message  string
	Position lexer.Position
}

// Error implements the error interface
func (re *RuntimeError) Error() string {
	return fmt.Sprintf("runtime error at line %d, column %d: %s", re.Position.Line, re.Position.Column, re.Message)
}

// binaryOperations maps operator strings to their corresponding Value methods
var binaryOperations = map[string]func(Value, Value) (Value, error){
	"+": Value.Add,
	"-": Value.Subtract,
	"*": Value.Multiply,
	"/": Value.Divide,
	"^": Value.Power,
}

// Interpreter executes BASIC programs by walking the AST
type Interpreter struct {
	runtime   runtime.Runtime
	variables map[string]Value     // Variable storage using proper Value types
	lineIndex map[int]*parser.Line // Maps line numbers to Line nodes for GOTO
	maxSteps  int                  // Maximum number of execution steps before infinite loop protection kicks in
	stepCount int                  // Current step count during execution
}

// NewInterpreter creates a new interpreter instance
func NewInterpreter(rt runtime.Runtime) *Interpreter {
	return &Interpreter{
		runtime:   rt,
		variables: make(map[string]Value),
		lineIndex: make(map[int]*parser.Line),
		maxSteps:  1000, // Default maximum steps
		stepCount: 0,
	}
}

// SetMaxSteps sets the maximum number of execution steps before infinite loop protection
func (i *Interpreter) SetMaxSteps(maxSteps int) {
	i.maxSteps = maxSteps
}

// Execute runs a BASIC program
func (i *Interpreter) Execute(program *parser.Program) error {
	// Reset step counter for new execution
	i.stepCount = 0

	// Build line number index for GOTO statements
	i.buildLineIndex(program)

	// Execute program with program counter for GOTO support
	return i.executeWithProgramCounter(program)
}

// buildLineIndex creates a map from line numbers to Line nodes
func (i *Interpreter) buildLineIndex(program *parser.Program) {
	for _, line := range program.Lines {
		i.lineIndex[line.Number] = line
	}
}

// executeWithProgramCounter executes program with support for GOTO jumps
func (i *Interpreter) executeWithProgramCounter(program *parser.Program) error {
	if len(program.Lines) == 0 {
		return nil
	}

	// Start execution at the first line
	currentLineIndex := 0

	for currentLineIndex < len(program.Lines) {
		line := program.Lines[currentLineIndex]

		for _, stmt := range line.Statements {
			// Increment step counter and check for infinite loop protection
			i.stepCount++
			if i.maxSteps > 0 && i.stepCount > i.maxSteps {
				return fmt.Errorf("?INFINITE LOOP ERROR")
			}

			err := i.executeStatement(stmt)
			if err != nil {
				return i.wrapErrorWithLine(err, line.Number)
			}

			// Check for flow control statements
			switch s := stmt.(type) {
			case *parser.EndStatement, *parser.StopStatement:
				return nil
			case *parser.GotoStatement:
				// Find the target line and jump to it
				targetLineIndex, found := i.findLineIndex(program, s.TargetLine)
				if !found {
					return fmt.Errorf("?UNDEFINED STATEMENT ERROR IN %d", line.Number)
				}
				currentLineIndex = targetLineIndex
				goto nextLine // Skip to the target line
			case *parser.IfStatement:
				// IF statement might execute a GOTO, need to check if flow changed
				if gotoStmt, ok := s.ThenStmt.(*parser.GotoStatement); ok {
					// Check if condition is true
					condition, err := i.evaluateExpression(s.Condition)
					if err != nil {
						return i.wrapErrorWithLine(err, line.Number)
					}
					if i.isConditionTrue(condition) {
						// Execute the GOTO
						targetLineIndex, found := i.findLineIndex(program, gotoStmt.TargetLine)
						if !found {
							return fmt.Errorf("?UNDEFINED STATEMENT ERROR IN %d", line.Number)
						}
						currentLineIndex = targetLineIndex
						goto nextLine // Skip to the target line
					}
				}
			}
		}

		// Move to next line
		currentLineIndex++
	nextLine:
	}

	return nil
}

// findLineIndex finds the index of a line with the given line number
func (i *Interpreter) findLineIndex(program *parser.Program, lineNumber int) (int, bool) {
	for index, line := range program.Lines {
		if line.Number == lineNumber {
			return index, true
		}
	}
	return 0, false
}

// executeStatement executes a single statement
func (i *Interpreter) executeStatement(stmt parser.Statement) error {
	switch s := stmt.(type) {
	case *parser.PrintStatement:
		return i.executePrintStatement(s)
	case *parser.LetStatement:
		return i.executeLetStatement(s)
	case *parser.EndStatement:
		// END statement - just return, handled in executeWithProgramCounter
		return nil
	case *parser.RunStatement:
		return i.executeRunStatement(s)
	case *parser.StopStatement:
		return i.executeStopStatement(s)
	case *parser.GotoStatement:
		// GOTO statement - handled in executeWithProgramCounter
		return nil
	case *parser.IfStatement:
		return i.executeIfStatement(s)
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
		normalizedName := i.normalizeVariableName(e.Name)
		if value, exists := i.variables[normalizedName]; exists {
			return value, nil
		}
		// Default value for uninitialized variables depends on type
		if strings.HasSuffix(e.Name, "$") {
			return NewStringValue(""), nil // String variables default to empty string
		}
		return NewNumberValue(0), nil // Numeric variables default to 0
	case *parser.BinaryOperation:
		return i.evaluateBinaryOperation(e)
	case *parser.UnaryOperation:
		return i.evaluateUnaryOperation(e)
	case *parser.ComparisonExpression:
		return i.evaluateComparisonExpression(e)
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

	if operation, exists := binaryOperations[expr.Operator]; exists {
		return operation(left, right)
	}

	return Value{}, fmt.Errorf("unknown operator: %s", expr.Operator)
}

// evaluateUnaryOperation evaluates a unary arithmetic operation
func (i *Interpreter) evaluateUnaryOperation(expr *parser.UnaryOperation) (Value, error) {
	operand, err := i.evaluateExpression(expr.Right)
	if err != nil {
		return Value{}, err
	}

	switch expr.Operator {
	case "-":
		// Negate the operand
		if operand.Type == NumberType {
			return NewNumberValue(-operand.Number), nil
		}
		return Value{}, fmt.Errorf("cannot negate non-numeric value")
	case "+":
		// Unary plus - just return the operand
		if operand.Type == NumberType {
			return operand, nil
		}
		return Value{}, fmt.Errorf("cannot apply unary plus to non-numeric value")
	default:
		return Value{}, fmt.Errorf("unknown unary operator: %s", expr.Operator)
	}
}

// evaluateComparisonExpression evaluates a comparison operation
func (i *Interpreter) evaluateComparisonExpression(expr *parser.ComparisonExpression) (Value, error) {
	left, err := i.evaluateExpression(expr.Left)
	if err != nil {
		return Value{}, err
	}

	right, err := i.evaluateExpression(expr.Right)
	if err != nil {
		return Value{}, err
	}

	// Perform the comparison based on operator
	result, err := i.compareValues(left, right, expr.Operator)
	if err != nil {
		return Value{}, err
	}

	// Return 1 for true, 0 for false (C64 BASIC convention)
	if result {
		return NewNumberValue(1), nil
	} else {
		return NewNumberValue(0), nil
	}
}

// compareValues compares two values using the specified operator
func (i *Interpreter) compareValues(left, right Value, operator string) (bool, error) {
	// Handle comparison based on types
	if left.Type == NumberType && right.Type == NumberType {
		// Numeric comparison
		return i.compareNumbers(left.Number, right.Number, operator), nil
	} else if left.Type == StringType && right.Type == StringType {
		// String comparison
		return i.compareStrings(left.String, right.String, operator), nil
	} else {
		// Type mismatch
		return false, fmt.Errorf("TYPE MISMATCH ERROR")
	}
}

// compareNumbers performs numeric comparison
func (i *Interpreter) compareNumbers(left, right float64, operator string) bool {
	switch operator {
	case "=":
		return left == right
	case "<>":
		return left != right
	case "<":
		return left < right
	case ">":
		return left > right
	case "<=":
		return left <= right
	case ">=":
		return left >= right
	default:
		return false // Invalid operator
	}
}

// compareStrings performs string comparison
func (i *Interpreter) compareStrings(left, right string, operator string) bool {
	switch operator {
	case "=":
		return left == right
	case "<>":
		return left != right
	case "<":
		return left < right
	case ">":
		return left > right
	case "<=":
		return left <= right
	case ">=":
		return left >= right
	default:
		return false // Invalid operator
	}
}

// executeLetStatement executes a LET statement (variable assignment)
func (i *Interpreter) executeLetStatement(stmt *parser.LetStatement) error {
	value, err := i.evaluateExpression(stmt.Expression)
	if err != nil {
		return err
	}

	// Type check: string variables can only hold strings, numeric variables can only hold numbers
	isStringVariable := strings.HasSuffix(stmt.Variable, "$")
	if isStringVariable && value.Type != StringType {
		return fmt.Errorf("TYPE MISMATCH ERROR")
	}
	if !isStringVariable && value.Type != NumberType {
		return fmt.Errorf("TYPE MISMATCH ERROR")
	}

	normalizedName := i.normalizeVariableName(stmt.Variable)
	i.variables[normalizedName] = value
	return nil
}

// executeRunStatement executes a RUN statement
func (i *Interpreter) executeRunStatement(stmt *parser.RunStatement) error {
	// RUN statement doesn't do anything during normal program execution
	// In a C64 BASIC, RUN would start program execution from the beginning,
	// but in our current architecture, we're already executing the program
	// so RUN is effectively a no-op when encountered in program flow
	return nil
}

// executeStopStatement executes a STOP statement
func (i *Interpreter) executeStopStatement(stmt *parser.StopStatement) error {
	// STOP statement - execution handled in Execute method
	return nil
}

// executeIfStatement executes an IF...THEN statement
func (i *Interpreter) executeIfStatement(stmt *parser.IfStatement) error {
	// Evaluate the condition
	condition, err := i.evaluateExpression(stmt.Condition)
	if err != nil {
		return err
	}

	// Check if condition is true
	if i.isConditionTrue(condition) {
		// Execute the THEN statement
		return i.executeStatement(stmt.ThenStmt)
	}

	// Condition is false, do nothing
	return nil
}

// isConditionTrue determines if a condition value evaluates to true
func (i *Interpreter) isConditionTrue(value Value) bool {
	switch value.Type {
	case NumberType:
		return value.Number != 0 // Non-zero numbers are true
	case StringType:
		return value.String != "" // Non-empty strings are true
	default:
		return false
	}
}

// normalizeVariableName truncates variable name to first 2 characters (C64 BASIC behavior)
func (i *Interpreter) normalizeVariableName(name string) string {
	if len(name) > 2 {
		return name[:2]
	}
	return name
}

// wrapErrorWithLine wraps an error with C64 BASIC format including line number
func (i *Interpreter) wrapErrorWithLine(err error, lineNumber int) error {
	// Check if it's already a C64 format error (starts with ?)
	errMsg := err.Error()
	if len(errMsg) > 0 && errMsg[0] == '?' {
		return err // Already formatted
	}

	// Convert common errors to C64 BASIC format
	switch {
	case strings.Contains(errMsg, "division by zero"):
		return fmt.Errorf("?DIVISION BY ZERO ERROR IN %d", lineNumber)
	case strings.Contains(errMsg, "TYPE MISMATCH ERROR"):
		return fmt.Errorf("?TYPE MISMATCH ERROR IN %d", lineNumber)
	default:
		return fmt.Errorf("?ERROR IN %d: %s", lineNumber, errMsg)
	}
}

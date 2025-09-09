// ABOUTME: Abstract Syntax Tree node definitions for BASIC language
// ABOUTME: Defines the structure of parsed BASIC programs as tree nodes

package parser

import (
	"basic-interpreter/types"
	"fmt"
	"strings"
)

// Node represents any node in the AST
type Node interface {
	GetLineNumber() int
}

// BaseNode provides common functionality for all AST nodes
type BaseNode struct {
	Line int // Source line number
}

func (bn BaseNode) GetLineNumber() int { return bn.Line }

// InterpreterOperations defines what AST nodes can ask the interpreter to do
// This interface enables double dispatch: AST nodes call back to interpreter
// operations without directly depending on the interpreter implementation
type InterpreterOperations interface {
	// Variable operations
	GetVariable(name string) (types.Value, error)
	SetVariable(name string, value types.Value) error

	// I/O operations
	Print(text string) error
	PrintLine(text string) error
	ReadInput(prompt string) (string, error)

	// Control flow requests
	RequestGoto(targetLine int) error
	RequestEnd() error
	RequestStop() error
	RequestGosub(targetLine int) error
	RequestReturn() error

	// Loop control for FOR/NEXT
	BeginFor(variable string, end types.Value, step types.Value) error
	IterateFor(variable string) error

	// Utility operations
	NormalizeVariableName(name string) string

	// Data management (READ/DATA)
	GetNextData() (types.Value, error)

	// Array management (DIM)
	DeclareArray(name string, size int, isString bool) error

	// Function evaluation
	EvaluateFunction(functionName string, args []Expression) (types.Value, error)

	// Array element operations
	GetArrayElement(name string, index int) (types.Value, error)
	SetArrayElement(name string, index int, value types.Value) error
}

// (No control error types are used for END/STOP; interpreter handles them statefully.)

// (No control error types for FOR/NEXT; interpreter handles via BeginFor/IterateFor.)

// Statement represents any statement node
type Statement interface {
	Node
	Execute(ops InterpreterOperations) error
}

// Expression represents any expression node
type Expression interface {
	Node
	Evaluate(ops InterpreterOperations) (types.Value, error)
}

// Program represents the root of the AST - a complete BASIC program
type Program struct {
	BaseNode
	Lines []*Line
}

// Program overrides GetLineNumber to return the first line's number
func (p *Program) GetLineNumber() int {
	if len(p.Lines) > 0 {
		return p.Lines[0].SourceLine
	}
	return 0
}

// Line represents a single line in a BASIC program
type Line struct {
	BaseNode
	Number     int         // BASIC line number (10, 20, etc.)
	Statements []Statement // Statements on this line
	SourceLine int         // Source line number for error reporting
}

// Line overrides GetLineNumber to return SourceLine
func (l *Line) GetLineNumber() int { return l.SourceLine }

// PrintStatement represents a PRINT statement
type PrintStatement struct {
	BaseNode
	// Legacy single expression (used when Items is empty)
	Expression Expression
	// Items is a list of expressions to print in sequence (semicolon/comma separated)
	Items []Expression
	// If true, suppress the trailing newline (trailing ';' in PRINT)
	NoNewline bool
}

func (ps *PrintStatement) Execute(ops InterpreterOperations) error {
	// If multiple items are present, concatenate them into a single output string
	if len(ps.Items) > 0 {
		var out string
		var prevType types.ValueType = -1
		for idx, it := range ps.Items {
			v, err := it.Evaluate(ops)
			if err != nil {
				return err
			}
			curr := v.ToString()
			// Insert a single space between items when either side is numeric,
			// but avoid double spaces if spacing is already present.
			if idx > 0 {
				if v.Type == types.NumberType || prevType == types.NumberType {
					needSpace := true
					if len(out) > 0 && out[len(out)-1] == ' ' {
						needSpace = false
					}
					if len(curr) > 0 && (curr[0] == ' ' || curr[0] == ',' || curr[0] == '.' || curr[0] == ';' || curr[0] == ':' || curr[0] == ')') {
						needSpace = false
					}
					if needSpace {
						out += " "
					}
				}
			}
			out += curr
			prevType = v.Type
		}
		if ps.NoNewline {
			return ops.Print(out)
		}
		return ops.PrintLine(out)
	}
	// Legacy behavior: single expression
	value, err := ps.Expression.Evaluate(ops)
	if err != nil {
		return err
	}
	return ops.PrintLine(value.ToString())
}

// StringLiteral represents a string literal expression
type StringLiteral struct {
	BaseNode
	Value string // The string value (without quotes)
}

func (sl *StringLiteral) Evaluate(ops InterpreterOperations) (types.Value, error) {
	return types.NewStringValue(sl.Value), nil
}

// EndStatement represents an END statement
type EndStatement struct {
	BaseNode
}

func (es *EndStatement) Execute(ops InterpreterOperations) error {
	return ops.RequestEnd()
}

// LetStatement represents a LET assignment statement
type LetStatement struct {
	BaseNode
	Variable   string     // Variable name
	Expression Expression // Value to assign
}

func (ls *LetStatement) Execute(ops InterpreterOperations) error {
	value, err := ls.Expression.Evaluate(ops)
	if err != nil {
		return err
	}
	return ops.SetVariable(ls.Variable, value)
}

// VariableReference represents a variable reference in an expression
type VariableReference struct {
	BaseNode
	Name string // Variable name
}

func (vr *VariableReference) Evaluate(ops InterpreterOperations) (types.Value, error) {
	return ops.GetVariable(vr.Name)
}

// NumberLiteral represents a numeric literal expression
type NumberLiteral struct {
	BaseNode
	Value string // The numeric value as string
}

func (nl *NumberLiteral) Evaluate(ops InterpreterOperations) (types.Value, error) {
	return types.ParseValue(nl.Value)
}

// BinaryOperation represents a binary arithmetic operation
type BinaryOperation struct {
	BaseNode
	Left     Expression // Left operand
	Operator string     // Operator (+, -, *, /, ^)
	Right    Expression // Right operand
}

func (bo *BinaryOperation) Evaluate(ops InterpreterOperations) (types.Value, error) {
	left, err := bo.Left.Evaluate(ops)
	if err != nil {
		return types.Value{}, err
	}

	right, err := bo.Right.Evaluate(ops)
	if err != nil {
		return types.Value{}, err
	}

	// Use the binary operations map from interpreter package
	switch bo.Operator {
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
		return types.Value{}, fmt.Errorf("unknown operator: %s", bo.Operator)
	}
}

// RunStatement represents a RUN statement
type RunStatement struct {
	BaseNode
}

func (rs *RunStatement) Execute(ops InterpreterOperations) error {
	// RUN statement doesn't do anything during normal program execution
	// In C64 BASIC, RUN would start program execution from the beginning,
	// but in our current architecture, we're already executing the program
	// so RUN is effectively a no-op when encountered in program flow
	return nil
}

// StopStatement represents a STOP statement
type StopStatement struct {
	BaseNode
}

func (ss *StopStatement) Execute(ops InterpreterOperations) error {
	return ops.RequestStop()
}

// InputStatement represents an INPUT statement
type InputStatement struct {
	BaseNode
	Prompt   string // Optional prompt string (empty for no prompt)
	Variable string // Variable name to read into
}

func (ins *InputStatement) Execute(ops InterpreterOperations) error {
	input, err := ops.ReadInput(ins.Prompt)
	if err != nil {
		return err
	}

	// Parse input based on variable type
	var value types.Value
	if strings.HasSuffix(ins.Variable, "$") {
		value = types.NewStringValue(input)
	} else {
		parsed, err := types.ParseValue(input)
		if err != nil || parsed.Type != types.NumberType {
			return types.ErrTypeMismatch
		}
		value = parsed
	}

	return ops.SetVariable(ins.Variable, value)
}

// GotoStatement represents a GOTO statement
type GotoStatement struct {
	BaseNode
	TargetLine int // Target line number to jump to
}

func (gs *GotoStatement) Execute(ops InterpreterOperations) error {
	return ops.RequestGoto(gs.TargetLine)
}

// IfStatement represents an IF...THEN statement
type IfStatement struct {
	BaseNode
	Condition Expression // The condition to evaluate
	ThenStmt  Statement  // The statement to execute if condition is true
}

func (is *IfStatement) Execute(ops InterpreterOperations) error {
	condition, err := is.Condition.Evaluate(ops)
	if err != nil {
		return err
	}

	if condition.IsTrue() {
		return is.ThenStmt.Execute(ops)
	}
	return nil
}

// UnaryOperation represents a unary arithmetic operation
type UnaryOperation struct {
	BaseNode
	Operator string     // Operator (-)
	Right    Expression // Right operand
}

func (uo *UnaryOperation) Evaluate(ops InterpreterOperations) (types.Value, error) {
	operand, err := uo.Right.Evaluate(ops)
	if err != nil {
		return types.Value{}, err
	}

	switch uo.Operator {
	case "-":
		// Negate the operand
		if operand.Type == types.NumberType {
			return types.NewNumberValue(-operand.Number), nil
		}
		return types.Value{}, fmt.Errorf("cannot negate non-numeric value")
	case "+":
		// Unary plus - just return the operand
		if operand.Type == types.NumberType {
			return operand, nil
		}
		return types.Value{}, fmt.Errorf("cannot apply unary plus to non-numeric value")
	default:
		return types.Value{}, fmt.Errorf("unknown unary operator: %s", uo.Operator)
	}
}

// ComparisonExpression represents a comparison operation (=, <>, <, >, <=, >=)
type ComparisonExpression struct {
	BaseNode
	Left     Expression // Left operand
	Operator string     // Comparison operator
	Right    Expression // Right operand
}

func (ce *ComparisonExpression) Evaluate(ops InterpreterOperations) (types.Value, error) {
	left, err := ce.Left.Evaluate(ops)
	if err != nil {
		return types.Value{}, err
	}

	right, err := ce.Right.Evaluate(ops)
	if err != nil {
		return types.Value{}, err
	}

	// Perform the comparison based on operator
	result, err := left.Compare(right, ce.Operator)
	if err != nil {
		return types.Value{}, err
	}

	// Return 1 for true, 0 for false (C64 BASIC convention)
	if result {
		return types.NewNumberValue(1), nil
	} else {
		return types.NewNumberValue(0), nil
	}
}

// ForStatement represents a FOR loop statement
type ForStatement struct {
	BaseNode
	Variable   string     // Loop variable name
	StartValue Expression // Starting value
	EndValue   Expression // Ending value
	StepValue  Expression // Optional step value (defaults to 1)
}

func (fs *ForStatement) Execute(ops InterpreterOperations) error {
	startVal, err := fs.StartValue.Evaluate(ops)
	if err != nil {
		return err
	}

	endVal, err := fs.EndValue.Evaluate(ops)
	if err != nil {
		return err
	}

	// Evaluate step value if provided, otherwise default to 1
	var stepVal types.Value
	if fs.StepValue != nil {
		s, err := fs.StepValue.Evaluate(ops)
		if err != nil {
			return err
		}
		if s.Type != types.NumberType {
			return fmt.Errorf("TYPE MISMATCH ERROR")
		}
		stepVal = s
	} else {
		stepVal = types.NewNumberValue(1)
	}

	// Initialize loop variable
	err = ops.SetVariable(fs.Variable, startVal)
	if err != nil {
		return err
	}

	// Begin the FOR loop with provided step
	return ops.BeginFor(fs.Variable, endVal, stepVal)
}

// NextStatement represents a NEXT statement
type NextStatement struct {
	BaseNode
	Variable string // Optional loop variable name (can be empty)
}

func (ns *NextStatement) Execute(ops InterpreterOperations) error {
	// Iterate the FOR loop via interpreter operations
	return ops.IterateFor(ns.Variable)
}

// GosubStatement represents a GOSUB statement
type GosubStatement struct {
	BaseNode
	TargetLine int // Target line number to call
}

func (gs *GosubStatement) Execute(ops InterpreterOperations) error {
	return ops.RequestGosub(gs.TargetLine)
}

// ReturnStatement represents a RETURN statement
type ReturnStatement struct {
	BaseNode
}

func (rs *ReturnStatement) Execute(ops InterpreterOperations) error {
	return ops.RequestReturn()
}

// DataStatement represents a DATA statement containing a list of constants
type DataStatement struct {
	BaseNode
	Values []Expression // Constants (numbers or strings)
}

// DATA is processed before execution by the interpreter; at runtime it's a no-op
func (ds *DataStatement) Execute(ops InterpreterOperations) error { return nil }

// ReadStatement represents a READ statement to read values from DATA
type ReadStatement struct {
	BaseNode
	Variables []string // Variable names to fill
}

func (rs *ReadStatement) Execute(ops InterpreterOperations) error {
	for _, vname := range rs.Variables {
		val, err := ops.GetNextData()
		if err != nil {
			return err
		}
		// Type check based on variable suffix
		if strings.HasSuffix(vname, "$") {
			if val.Type != types.StringType {
				return types.ErrTypeMismatch
			}
		} else {
			if val.Type != types.NumberType {
				return types.ErrTypeMismatch
			}
		}
		if err := ops.SetVariable(vname, val); err != nil {
			return err
		}
	}
	return nil
}

// RemStatement represents a REM (comment) statement; it is a no-op at runtime
type RemStatement struct {
	BaseNode
}

func (rs *RemStatement) Execute(ops InterpreterOperations) error { return nil }

// FunctionCall represents a function call expression
type FunctionCall struct {
	BaseNode
	FunctionName string       // Function name (LEN, LEFT$, RIGHT$, etc.)
	Arguments    []Expression // Function arguments
}

func (fc *FunctionCall) Evaluate(ops InterpreterOperations) (types.Value, error) {
	return ops.EvaluateFunction(fc.FunctionName, fc.Arguments)
}

// ArrayReference represents access to an array element, e.g., A(5)
type ArrayReference struct {
	BaseNode
	Name  string
	Index Expression
}

func (ar *ArrayReference) Evaluate(ops InterpreterOperations) (types.Value, error) {
	idxVal, err := ar.Index.Evaluate(ops)
	if err != nil {
		return types.Value{}, err
	}
	if idxVal.Type != types.NumberType {
		return types.Value{}, types.ErrTypeMismatch
	}
	n := idxVal.Number
	if n < 0 || float64(int(n)) != n {
		return types.Value{}, fmt.Errorf("?ILLEGAL QUANTITY ERROR")
	}
	return ops.GetArrayElement(ar.Name, int(n))
}

// ArraySetStatement assigns a value to an array element, e.g., A(5) = 42
type ArraySetStatement struct {
	BaseNode
	Name       string
	Index      Expression
	Expression Expression
}

func (as *ArraySetStatement) Execute(ops InterpreterOperations) error {
	idxVal, err := as.Index.Evaluate(ops)
	if err != nil {
		return err
	}
	if idxVal.Type != types.NumberType {
		return types.ErrTypeMismatch
	}
	n := idxVal.Number
	if n < 0 || float64(int(n)) != n {
		return fmt.Errorf("?ILLEGAL QUANTITY ERROR")
	}
	val, err := as.Expression.Evaluate(ops)
	if err != nil {
		return err
	}
	return ops.SetArrayElement(as.Name, int(n), val)
}

// DimDeclaration represents a single array declaration inside a DIM statement
type DimDeclaration struct {
	Name string
	Size Expression
}

// DimStatement represents a DIM statement declaring one or more arrays
type DimStatement struct {
	BaseNode
	Declarations []DimDeclaration
}

func (ds *DimStatement) Execute(ops InterpreterOperations) error {
	for _, d := range ds.Declarations {
		// Evaluate size
		val, err := d.Size.Evaluate(ops)
		if err != nil {
			return err
		}
		if val.Type != types.NumberType {
			return types.ErrTypeMismatch
		}
		n := val.Number
		// Size must be integer and >= 0
		if n < 0 || float64(int(n)) != n {
			return fmt.Errorf("?ILLEGAL QUANTITY ERROR")
		}
		isString := strings.HasSuffix(d.Name, "$")
		if err := ops.DeclareArray(d.Name, int(n), isString); err != nil {
			return err
		}
	}
	return nil
}

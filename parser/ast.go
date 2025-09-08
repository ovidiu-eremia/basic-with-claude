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

	// Function evaluation
	EvaluateFunction(functionName string, args []Expression) (types.Value, error)
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
	Lines []*Line
}

func (p *Program) GetLineNumber() int {
	if len(p.Lines) > 0 {
		return p.Lines[0].SourceLine
	}
	return 0
}

// Line represents a single line in a BASIC program
type Line struct {
	Number     int         // BASIC line number (10, 20, etc.)
	Statements []Statement // Statements on this line
	SourceLine int         // Source line number for error reporting
}

func (l *Line) GetLineNumber() int { return l.SourceLine }

// PrintStatement represents a PRINT statement
type PrintStatement struct {
	// Legacy single expression (used when Items is empty)
	Expression Expression
	// Items is a list of expressions to print in sequence (semicolon/comma separated)
	Items []Expression
	// If true, suppress the trailing newline (trailing ';' in PRINT)
	NoNewline bool
	Line      int // Source line number
}

func (ps *PrintStatement) GetLineNumber() int { return ps.Line }

func (ps *PrintStatement) Execute(ops InterpreterOperations) error {
	// If multiple items are present, concatenate them into a single output string
	if len(ps.Items) > 0 {
		var out string
		for _, it := range ps.Items {
			v, err := it.Evaluate(ops)
			if err != nil {
				return err
			}
			out += v.ToString()
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
	Value string // The string value (without quotes)
	Line  int    // Source line number
}

func (sl *StringLiteral) GetLineNumber() int { return sl.Line }

func (sl *StringLiteral) Evaluate(ops InterpreterOperations) (types.Value, error) {
	return types.NewStringValue(sl.Value), nil
}

// EndStatement represents an END statement
type EndStatement struct {
	Line int // Source line number
}

func (es *EndStatement) GetLineNumber() int { return es.Line }

func (es *EndStatement) Execute(ops InterpreterOperations) error {
	return ops.RequestEnd()
}

// LetStatement represents a LET assignment statement
type LetStatement struct {
	Variable   string     // Variable name
	Expression Expression // Value to assign
	Line       int        // Source line number
}

func (ls *LetStatement) GetLineNumber() int { return ls.Line }

func (ls *LetStatement) Execute(ops InterpreterOperations) error {
	value, err := ls.Expression.Evaluate(ops)
	if err != nil {
		return err
	}
	return ops.SetVariable(ls.Variable, value)
}

// VariableReference represents a variable reference in an expression
type VariableReference struct {
	Name string // Variable name
	Line int    // Source line number
}

func (vr *VariableReference) GetLineNumber() int { return vr.Line }

func (vr *VariableReference) Evaluate(ops InterpreterOperations) (types.Value, error) {
	return ops.GetVariable(vr.Name)
}

// NumberLiteral represents a numeric literal expression
type NumberLiteral struct {
	Value string // The numeric value as string
	Line  int    // Source line number
}

func (nl *NumberLiteral) GetLineNumber() int { return nl.Line }

func (nl *NumberLiteral) Evaluate(ops InterpreterOperations) (types.Value, error) {
	return types.ParseValue(nl.Value)
}

// BinaryOperation represents a binary arithmetic operation
type BinaryOperation struct {
	Left     Expression // Left operand
	Operator string     // Operator (+, -, *, /, ^)
	Right    Expression // Right operand
	Line     int        // Source line number
}

func (bo *BinaryOperation) GetLineNumber() int { return bo.Line }

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
	Line int // Source line number
}

func (rs *RunStatement) GetLineNumber() int { return rs.Line }

func (rs *RunStatement) Execute(ops InterpreterOperations) error {
	// RUN statement doesn't do anything during normal program execution
	// In C64 BASIC, RUN would start program execution from the beginning,
	// but in our current architecture, we're already executing the program
	// so RUN is effectively a no-op when encountered in program flow
	return nil
}

// StopStatement represents a STOP statement
type StopStatement struct {
	Line int // Source line number
}

func (ss *StopStatement) GetLineNumber() int { return ss.Line }

func (ss *StopStatement) Execute(ops InterpreterOperations) error {
	return ops.RequestStop()
}

// InputStatement represents an INPUT statement
type InputStatement struct {
	Prompt   string // Optional prompt string (empty for no prompt)
	Variable string // Variable name to read into
	Line     int    // Source line number
}

func (ins *InputStatement) GetLineNumber() int { return ins.Line }

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
	TargetLine int // Target line number to jump to
	Line       int // Source line number
}

func (gs *GotoStatement) GetLineNumber() int { return gs.Line }

func (gs *GotoStatement) Execute(ops InterpreterOperations) error {
	return ops.RequestGoto(gs.TargetLine)
}

// IfStatement represents an IF...THEN statement
type IfStatement struct {
	Condition Expression // The condition to evaluate
	ThenStmt  Statement  // The statement to execute if condition is true
	Line      int        // Source line number
}

func (is *IfStatement) GetLineNumber() int { return is.Line }

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
	Operator string     // Operator (-)
	Right    Expression // Right operand
	Line     int        // Source line number
}

func (uo *UnaryOperation) GetLineNumber() int { return uo.Line }

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
	Left     Expression // Left operand
	Operator string     // Comparison operator
	Right    Expression // Right operand
	Line     int        // Source line number
}

func (ce *ComparisonExpression) GetLineNumber() int { return ce.Line }

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
	Variable   string     // Loop variable name
	StartValue Expression // Starting value
	EndValue   Expression // Ending value
	StepValue  Expression // Optional step value (defaults to 1)
	Line       int        // Source line number
}

func (fs *ForStatement) GetLineNumber() int { return fs.Line }

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
	Variable string // Optional loop variable name (can be empty)
	Line     int    // Source line number
}

func (ns *NextStatement) GetLineNumber() int { return ns.Line }

func (ns *NextStatement) Execute(ops InterpreterOperations) error {
	// Iterate the FOR loop via interpreter operations
	return ops.IterateFor(ns.Variable)
}

// GosubStatement represents a GOSUB statement
type GosubStatement struct {
	TargetLine int // Target line number to call
	Line       int // Source line number
}

func (gs *GosubStatement) GetLineNumber() int { return gs.Line }

func (gs *GosubStatement) Execute(ops InterpreterOperations) error {
	return ops.RequestGosub(gs.TargetLine)
}

// ReturnStatement represents a RETURN statement
type ReturnStatement struct {
	Line int // Source line number
}

func (rs *ReturnStatement) GetLineNumber() int { return rs.Line }

func (rs *ReturnStatement) Execute(ops InterpreterOperations) error {
	return ops.RequestReturn()
}

// DataStatement represents a DATA statement containing a list of constants
type DataStatement struct {
	Values []Expression // Constants (numbers or strings)
	Line   int          // Source line number
}

func (ds *DataStatement) GetLineNumber() int { return ds.Line }

// DATA is processed before execution by the interpreter; at runtime it's a no-op
func (ds *DataStatement) Execute(ops InterpreterOperations) error { return nil }

// ReadStatement represents a READ statement to read values from DATA
type ReadStatement struct {
	Variables []string // Variable names to fill
	Line      int      // Source line number
}

func (rs *ReadStatement) GetLineNumber() int { return rs.Line }

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
	Line int // Source line number
}

func (rs *RemStatement) GetLineNumber() int { return rs.Line }

func (rs *RemStatement) Execute(ops InterpreterOperations) error { return nil }

// FunctionCall represents a function call expression
type FunctionCall struct {
	FunctionName string       // Function name (LEN, LEFT$, RIGHT$, etc.)
	Arguments    []Expression // Function arguments
	Line         int          // Source line number
}

func (fc *FunctionCall) GetLineNumber() int { return fc.Line }

func (fc *FunctionCall) Evaluate(ops InterpreterOperations) (types.Value, error) {
	return ops.EvaluateFunction(fc.FunctionName, fc.Arguments)
}

// ABOUTME: Abstract Syntax Tree node definitions for BASIC language
// ABOUTME: Defines the structure of parsed BASIC programs as tree nodes

package parser

// Node represents any node in the AST
type Node interface {
	GetLineNumber() int
}

// Statement represents any statement node
type Statement interface {
	Node
	statementNode()
}

// Expression represents any expression node
type Expression interface {
	Node
	expressionNode()
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
	Expression Expression // What to print
	Line       int        // Source line number
}

func (ps *PrintStatement) statementNode()     {}
func (ps *PrintStatement) GetLineNumber() int { return ps.Line }

// StringLiteral represents a string literal expression
type StringLiteral struct {
	Value string // The string value (without quotes)
	Line  int    // Source line number
}

func (sl *StringLiteral) expressionNode()    {}
func (sl *StringLiteral) GetLineNumber() int { return sl.Line }

// EndStatement represents an END statement
type EndStatement struct {
	Line int // Source line number
}

func (es *EndStatement) statementNode()     {}
func (es *EndStatement) GetLineNumber() int { return es.Line }

// LetStatement represents a LET assignment statement
type LetStatement struct {
	Variable   string     // Variable name
	Expression Expression // Value to assign
	Line       int        // Source line number
}

func (ls *LetStatement) statementNode()     {}
func (ls *LetStatement) GetLineNumber() int { return ls.Line }

// VariableReference represents a variable reference in an expression
type VariableReference struct {
	Name string // Variable name
	Line int    // Source line number
}

func (vr *VariableReference) expressionNode()    {}
func (vr *VariableReference) GetLineNumber() int { return vr.Line }

// NumberLiteral represents a numeric literal expression
type NumberLiteral struct {
	Value string // The numeric value as string
	Line  int    // Source line number
}

func (nl *NumberLiteral) expressionNode()    {}
func (nl *NumberLiteral) GetLineNumber() int { return nl.Line }

// BinaryOperation represents a binary arithmetic operation
type BinaryOperation struct {
	Left     Expression // Left operand
	Operator string     // Operator (+, -, *, /, ^)
	Right    Expression // Right operand
	Line     int        // Source line number
}

func (bo *BinaryOperation) expressionNode()    {}
func (bo *BinaryOperation) GetLineNumber() int { return bo.Line }

// RunStatement represents a RUN statement
type RunStatement struct {
	Line int // Source line number
}

func (rs *RunStatement) statementNode()     {}
func (rs *RunStatement) GetLineNumber() int { return rs.Line }

// StopStatement represents a STOP statement
type StopStatement struct {
	Line int // Source line number
}

func (ss *StopStatement) statementNode()     {}
func (ss *StopStatement) GetLineNumber() int { return ss.Line }

// GotoStatement represents a GOTO statement
type GotoStatement struct {
	TargetLine int // Target line number to jump to
	Line       int // Source line number
}

func (gs *GotoStatement) statementNode()     {}
func (gs *GotoStatement) GetLineNumber() int { return gs.Line }

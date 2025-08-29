// ABOUTME: Operator precedence handling utility for expression parsing
// ABOUTME: Centralizes precedence rules and provides clean precedence lookup functionality

package parser

import "basic-interpreter/lexer"

// precedence represents operator precedence levels
type precedence int

const (
	_ precedence = iota
	LOWEST
	COMPARE // =, <>, <, >, <=, >=
	SUM     // +, -
	PRODUCT // *, /
	POWER   // ^
	PREFIX  // -X or +X
	CALL    // functions (future use)
)

// PrecedenceTable manages operator precedence lookup
type PrecedenceTable struct {
	precedences map[lexer.TokenType]precedence
}

// NewPrecedenceTable creates a new precedence table with BASIC operator precedences
func NewPrecedenceTable() *PrecedenceTable {
	return &PrecedenceTable{
		precedences: map[lexer.TokenType]precedence{
			lexer.ASSIGN:   COMPARE,
			lexer.NE:       COMPARE,
			lexer.LT:       COMPARE,
			lexer.GT:       COMPARE,
			lexer.LE:       COMPARE,
			lexer.GE:       COMPARE,
			lexer.PLUS:     SUM,
			lexer.MINUS:    SUM,
			lexer.MULTIPLY: PRODUCT,
			lexer.DIVIDE:   PRODUCT,
			lexer.POWER:    POWER,
		},
	}
}

// GetPrecedence returns the precedence of a token type
func (pt *PrecedenceTable) GetPrecedence(tokenType lexer.TokenType) precedence {
	if prec, ok := pt.precedences[tokenType]; ok {
		return prec
	}
	return LOWEST
}

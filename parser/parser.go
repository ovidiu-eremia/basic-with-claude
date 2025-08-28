// ABOUTME: Parser for BASIC language syntax and Abstract Syntax Tree construction
// ABOUTME: Converts tokens from lexer into structured AST nodes representing the program

package parser

import (
	"fmt"
	"strconv"

	"basic-interpreter/lexer"
)

// Parser represents the parser state
type Parser struct {
	lexer *lexer.Lexer
	
	currentToken lexer.Token
	peekToken    lexer.Token
	
	errors []string
}

// New creates a new parser instance
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}
	
	// Read two tokens, so currentToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	
	return p
}

// nextToken advances both currentToken and peekToken
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// Errors returns parsing errors
func (p *Parser) Errors() []string {
	return p.errors
}

// addError adds an error message
func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}

// skipToNextLineOrEOF advances tokens until reaching newline or EOF for error recovery
func (p *Parser) skipToNextLineOrEOF() {
	for p.currentToken.Type != lexer.NEWLINE && p.currentToken.Type != lexer.EOF {
		p.nextToken()
	}
}

// ParseProgram parses the entire program
func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Lines = []*Line{}

	for p.currentToken.Type != lexer.EOF {
		// Skip newlines
		if p.currentToken.Type == lexer.NEWLINE {
			p.nextToken()
			continue
		}
		
		line := p.parseLine()
		if line != nil {
			program.Lines = append(program.Lines, line)
		}
		
		// parseLine() leaves us at NEWLINE or EOF, no need to advance
	}

	return program
}

// parseLine parses a single BASIC line
func (p *Parser) parseLine() *Line {
	if p.currentToken.Type != lexer.NUMBER {
		p.addError(fmt.Sprintf("expected line number, got %s", p.currentToken.Type))
		p.skipToNextLineOrEOF()
		return nil
	}

	lineNum, err := strconv.Atoi(p.currentToken.Literal)
	if err != nil {
		p.addError(fmt.Sprintf("invalid line number: %s", p.currentToken.Literal))
		p.skipToNextLineOrEOF()
		return nil
	}

	line := &Line{
		Number:     lineNum,
		Statements: []Statement{},
		SourceLine: p.currentToken.Line,
	}

	p.nextToken() // consume line number

	// Parse statements on this line
	for p.currentToken.Type != lexer.NEWLINE && p.currentToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			line.Statements = append(line.Statements, stmt)
		}
		// Advance token after parsing statement
		p.nextToken()
	}

	return line
}

// parseStatement parses a statement
func (p *Parser) parseStatement() Statement {
	switch p.currentToken.Type {
	case lexer.PRINT:
		return p.parsePrintStatement()
	case lexer.LET:
		return p.parseLetStatement()
	case lexer.IDENT:
		// Assignment without LET keyword
		return p.parseAssignmentStatement()
	case lexer.END:
		return p.parseEndStatement()
	case lexer.ILLEGAL:
		p.addError(fmt.Sprintf("illegal token: %s", p.currentToken.Literal))
		return nil
	default:
		p.addError(fmt.Sprintf("unknown statement: %s", p.currentToken.Type))
		return nil
	}
}

// parsePrintStatement parses a PRINT statement
func (p *Parser) parsePrintStatement() *PrintStatement {
	stmt := &PrintStatement{Line: p.currentToken.Line}

	p.nextToken() // consume PRINT

	stmt.Expression = p.parseExpression()
	
	return stmt
}

// Operator precedence levels
type precedence int

const (
	_ precedence = iota
	LOWEST
	SUM     // +, -
	PRODUCT // *, /
	POWER   // ^
	CALL    // functions (future use)
)

// precedences maps token types to their precedence levels
var precedences = map[lexer.TokenType]precedence{
	lexer.PLUS:     SUM,
	lexer.MINUS:    SUM,
	lexer.MULTIPLY: PRODUCT,
	lexer.DIVIDE:   PRODUCT,
	lexer.POWER:    POWER,
}

// parseExpression parses an expression using operator precedence parsing
func (p *Parser) parseExpression() Expression {
	return p.parseExpressionWithPrecedence(LOWEST)
}

// parseExpressionWithPrecedence parses expressions with given minimum precedence
func (p *Parser) parseExpressionWithPrecedence(minPrec precedence) Expression {
	left := p.parsePrimaryExpression()
	if left == nil {
		return nil
	}

	for p.peekToken.Type != lexer.NEWLINE && p.peekToken.Type != lexer.EOF && p.peekTokenPrecedence() > minPrec {
		operator := p.peekToken.Literal
		operatorPrec := p.peekTokenPrecedence()
		
		p.nextToken() // consume the operator
		p.nextToken() // move to right operand
		
		right := p.parseExpressionWithPrecedence(operatorPrec)
		if right == nil {
			return nil
		}
		
		left = &BinaryOperation{
			Left:     left,
			Operator: operator,
			Right:    right,
			Line:     left.GetLineNumber(),
		}
	}

	return left
}

// parsePrimaryExpression parses primary expressions (literals, variables, parentheses)
func (p *Parser) parsePrimaryExpression() Expression {
	switch p.currentToken.Type {
	case lexer.STRING:
		return p.parseStringLiteral()
	case lexer.NUMBER:
		return p.parseNumberLiteral()
	case lexer.IDENT:
		return p.parseVariableReference()
	case lexer.LPAREN:
		return p.parseGroupedExpression()
	case lexer.ILLEGAL:
		p.addError(fmt.Sprintf("illegal token in expression: %s", p.currentToken.Literal))
		return nil
	default:
		p.addError(fmt.Sprintf("unexpected token in expression: %s", p.currentToken.Type))
		return nil
	}
}

// parseGroupedExpression parses expressions in parentheses
func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken() // consume '('
	
	expr := p.parseExpression()
	if expr == nil {
		return nil
	}
	
	if p.peekToken.Type != lexer.RPAREN {
		p.addError("expected ')' after grouped expression")
		return nil
	}
	
	p.nextToken() // consume ')'
	return expr
}

// peekTokenPrecedence returns the precedence of the peek token
func (p *Parser) peekTokenPrecedence() precedence {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
}

// parseEndStatement parses an END statement
func (p *Parser) parseEndStatement() *EndStatement {
	return &EndStatement{
		Line: p.currentToken.Line,
	}
}

// parseStringLiteral parses a string literal
func (p *Parser) parseStringLiteral() *StringLiteral {
	return &StringLiteral{
		Value: p.currentToken.Literal,
		Line:  p.currentToken.Line,
	}
}

// parseNumberLiteral parses a number literal
func (p *Parser) parseNumberLiteral() *NumberLiteral {
	return &NumberLiteral{
		Value: p.currentToken.Literal,
		Line:  p.currentToken.Line,
	}
}

// parseVariableReference parses a variable reference
func (p *Parser) parseVariableReference() *VariableReference {
	return &VariableReference{
		Name: p.currentToken.Literal,
		Line: p.currentToken.Line,
	}
}

// parseLetStatement parses a LET statement
func (p *Parser) parseLetStatement() *LetStatement {
	p.nextToken() // consume LET
	return p.parseAssignment()
}

// parseAssignmentStatement parses an assignment statement without LET
func (p *Parser) parseAssignmentStatement() *LetStatement {
	return p.parseAssignment()
}

// parseAssignment parses variable assignment (with or without LET keyword)
func (p *Parser) parseAssignment() *LetStatement {
	stmt := &LetStatement{Line: p.currentToken.Line}
	
	if p.currentToken.Type != lexer.IDENT {
		p.addError(fmt.Sprintf("expected variable name, got %s", p.currentToken.Type))
		return nil
	}
	
	stmt.Variable = p.currentToken.Literal
	p.nextToken() // consume variable name
	
	if p.currentToken.Type != lexer.ASSIGN {
		p.addError(fmt.Sprintf("expected '=' after variable name, got %s", p.currentToken.Type))
		return nil
	}
	
	p.nextToken() // consume '='
	
	stmt.Expression = p.parseExpression()
	
	return stmt
}
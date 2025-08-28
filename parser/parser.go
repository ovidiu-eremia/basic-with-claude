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
	lexer      *lexer.Lexer
	precedence *PrecedenceTable
	
	currentToken lexer.Token
	peekToken    lexer.Token
	
	errors []string
}

// New creates a new parser instance
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:      l,
		precedence: NewPrecedenceTable(),
		errors:     []string{},
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

// addErrorf adds a formatted error message with current token context
func (p *Parser) addErrorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	p.addError(msg)
}

// addTokenError adds an error message with token type context
func (p *Parser) addTokenError(expected string, got lexer.TokenType) {
	p.addErrorf("expected %s, got %s", expected, got)
}

// addLiteralError adds an error message with token literal context
func (p *Parser) addLiteralError(prefix string, literal string) {
	p.addErrorf("%s: %s", prefix, literal)
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
		p.addTokenError("line number", p.currentToken.Type)
		p.skipToNextLineOrEOF()
		return nil
	}

	lineNum, err := strconv.Atoi(p.currentToken.Literal)
	if err != nil {
		p.addLiteralError("invalid line number", p.currentToken.Literal)
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
		return p.parseAssignmentStatement(true) // LET assignment
	case lexer.IDENT:
		return p.parseAssignmentStatement(false) // Direct assignment
	case lexer.END:
		return p.parseEndStatement()
	case lexer.ILLEGAL:
		p.addLiteralError("illegal token", p.currentToken.Literal)
		return nil
	default:
		p.addTokenError("valid statement", p.currentToken.Type)
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

	for p.peekToken.Type != lexer.NEWLINE && p.peekToken.Type != lexer.EOF && p.precedence.GetPrecedence(p.peekToken.Type) > minPrec {
		operator := p.peekToken.Literal
		operatorPrec := p.precedence.GetPrecedence(p.peekToken.Type)
		
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
		p.addLiteralError("illegal token in expression", p.currentToken.Literal)
		return nil
	default:
		p.addTokenError("valid expression", p.currentToken.Type)
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

// parseAssignmentStatement parses variable assignment (with or without LET keyword)
func (p *Parser) parseAssignmentStatement(hasLet bool) *LetStatement {
	if hasLet {
		p.nextToken() // consume LET token
	}
	return p.parseAssignment()
}

// parseAssignment parses variable assignment core logic
func (p *Parser) parseAssignment() *LetStatement {
	stmt := &LetStatement{Line: p.currentToken.Line}
	
	if p.currentToken.Type != lexer.IDENT {
		p.addTokenError("variable name", p.currentToken.Type)
		return nil
	}
	
	stmt.Variable = p.currentToken.Literal
	p.nextToken() // consume variable name
	
	if p.currentToken.Type != lexer.ASSIGN {
		p.addTokenError("'=' after variable name", p.currentToken.Type)
		return nil
	}
	
	p.nextToken() // consume '='
	
	stmt.Expression = p.parseExpression()
	
	return stmt
}
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
	
	curToken  lexer.Token
	peekToken lexer.Token
	
	errors []string
}

// New creates a new parser instance
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}
	
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	
	return p
}

// nextToken advances both curToken and peekToken
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
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
	for p.curToken.Type != lexer.NEWLINE && p.curToken.Type != lexer.EOF {
		p.nextToken()
	}
}

// ParseProgram parses the entire program
func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Lines = []*Line{}

	for p.curToken.Type != lexer.EOF {
		// Skip newlines
		if p.curToken.Type == lexer.NEWLINE {
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
	if p.curToken.Type != lexer.NUMBER {
		p.addError(fmt.Sprintf("expected line number, got %s", p.curToken.Type))
		p.skipToNextLineOrEOF()
		return nil
	}

	lineNum, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		p.addError(fmt.Sprintf("invalid line number: %s", p.curToken.Literal))
		p.skipToNextLineOrEOF()
		return nil
	}

	line := &Line{
		Number:     lineNum,
		Statements: []Statement{},
		SourceLine: p.curToken.Line,
	}

	p.nextToken() // consume line number

	// Parse statements on this line
	for p.curToken.Type != lexer.NEWLINE && p.curToken.Type != lexer.EOF {
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
	switch p.curToken.Type {
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
		p.addError(fmt.Sprintf("illegal token: %s", p.curToken.Literal))
		return nil
	default:
		p.addError(fmt.Sprintf("unknown statement: %s", p.curToken.Type))
		return nil
	}
}

// parsePrintStatement parses a PRINT statement
func (p *Parser) parsePrintStatement() *PrintStatement {
	stmt := &PrintStatement{Line: p.curToken.Line}

	p.nextToken() // consume PRINT

	stmt.Expression = p.parseExpression()
	
	return stmt
}

// parseExpression parses an expression
func (p *Parser) parseExpression() Expression {
	switch p.curToken.Type {
	case lexer.STRING:
		return p.parseStringLiteral()
	case lexer.NUMBER:
		return p.parseNumberLiteral()
	case lexer.IDENT:
		return p.parseVariableReference()
	case lexer.ILLEGAL:
		p.addError(fmt.Sprintf("illegal token in expression: %s", p.curToken.Literal))
		return nil
	default:
		p.addError(fmt.Sprintf("unexpected token in expression: %s", p.curToken.Type))
		return nil
	}
}

// parseEndStatement parses an END statement
func (p *Parser) parseEndStatement() *EndStatement {
	return &EndStatement{
		Line: p.curToken.Line,
	}
}

// parseStringLiteral parses a string literal
func (p *Parser) parseStringLiteral() *StringLiteral {
	return &StringLiteral{
		Value: p.curToken.Literal,
		Line:  p.curToken.Line,
	}
}

// parseNumberLiteral parses a number literal
func (p *Parser) parseNumberLiteral() *NumberLiteral {
	return &NumberLiteral{
		Value: p.curToken.Literal,
		Line:  p.curToken.Line,
	}
}

// parseVariableReference parses a variable reference
func (p *Parser) parseVariableReference() *VariableReference {
	return &VariableReference{
		Name: p.curToken.Literal,
		Line: p.curToken.Line,
	}
}

// parseLetStatement parses a LET statement
func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{Line: p.curToken.Line}
	
	p.nextToken() // consume LET
	
	if p.curToken.Type != lexer.IDENT {
		p.addError(fmt.Sprintf("expected variable name after LET, got %s", p.curToken.Type))
		return nil
	}
	
	stmt.Variable = p.curToken.Literal
	p.nextToken() // consume variable name
	
	if p.curToken.Type != lexer.ASSIGN {
		p.addError(fmt.Sprintf("expected '=' after variable name, got %s", p.curToken.Type))
		return nil
	}
	
	p.nextToken() // consume '='
	
	stmt.Expression = p.parseExpression()
	
	return stmt
}

// parseAssignmentStatement parses an assignment statement without LET
func (p *Parser) parseAssignmentStatement() *LetStatement {
	stmt := &LetStatement{Line: p.curToken.Line}
	
	stmt.Variable = p.curToken.Literal
	p.nextToken() // consume variable name
	
	if p.curToken.Type != lexer.ASSIGN {
		p.addError(fmt.Sprintf("expected '=' after variable name, got %s", p.curToken.Type))
		return nil
	}
	
	p.nextToken() // consume '='
	
	stmt.Expression = p.parseExpression()
	
	return stmt
}
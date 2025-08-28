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
	l *lexer.Lexer
	
	curToken  lexer.Token
	peekToken lexer.Token
	
	errors []string
}

// New creates a new parser instance
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
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
	p.peekToken = p.l.NextToken()
}

// Errors returns parsing errors
func (p *Parser) Errors() []string {
	return p.errors
}

// addError adds an error message
func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
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
		// Skip to next line or EOF to recover from error
		for p.curToken.Type != lexer.NEWLINE && p.curToken.Type != lexer.EOF {
			p.nextToken()
		}
		return nil
	}

	lineNum, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		p.addError(fmt.Sprintf("invalid line number: %s", p.curToken.Literal))
		// Skip to next line or EOF to recover from error
		for p.curToken.Type != lexer.NEWLINE && p.curToken.Type != lexer.EOF {
			p.nextToken()
		}
		return nil
	}

	line := &Line{
		Number:     lineNum,
		Statements: []Statement{},
		Line:       p.curToken.Line,
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
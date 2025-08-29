// ABOUTME: Parser for BASIC language syntax and Abstract Syntax Tree construction
// ABOUTME: Converts tokens from lexer into structured AST nodes representing the program

package parser

import (
	"fmt"
	"strconv"

	"basic-interpreter/lexer"
)

// ParseError represents an error that occurred during parsing
type ParseError struct {
	Message  string
	Position lexer.Position
}

// Error implements the error interface
func (pe *ParseError) Error() string {
	return fmt.Sprintf("parse error at line %d, column %d: %s", pe.Position.Line, pe.Position.Column, pe.Message)
}

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
	case lexer.RUN:
		return p.parseRunStatement()
	case lexer.STOP:
		return p.parseStopStatement()
	case lexer.GOTO:
		return p.parseGotoStatement()
	case lexer.IF:
		return p.parseIfStatement()
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

	// Check if there's an expression to print, or if we're at end of line/file
	if p.currentToken.Type != lexer.NEWLINE && p.currentToken.Type != lexer.EOF {
		stmt.Expression = p.parseExpression()
	} else {
		// No expression means print empty line
		stmt.Expression = &StringLiteral{Value: "", Line: stmt.Line}
	}

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

		// Create appropriate node type based on operator
		if p.isComparisonOperatorString(operator) {
			left = &ComparisonExpression{
				Left:     left,
				Operator: operator,
				Right:    right,
				Line:     left.GetLineNumber(),
			}
		} else {
			left = &BinaryOperation{
				Left:     left,
				Operator: operator,
				Right:    right,
				Line:     left.GetLineNumber(),
			}
		}
	}

	return left
}

// isComparisonOperatorString checks if an operator string is a comparison operator
func (p *Parser) isComparisonOperatorString(operator string) bool {
	switch operator {
	case "=", "<>", "<", ">", "<=", ">=":
		return true
	default:
		return false
	}
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
	case lexer.MINUS:
		return p.parseUnaryOperation()
	case lexer.ILLEGAL:
		p.addLiteralError("illegal token in expression", p.currentToken.Literal)
		return nil
	default:
		p.addTokenError("valid expression", p.currentToken.Type)
		return nil
	}
}

// parseUnaryOperation parses a unary operation
func (p *Parser) parseUnaryOperation() Expression {
	stmt := &UnaryOperation{
		Operator: p.currentToken.Literal,
		Line:     p.currentToken.Line,
	}
	p.nextToken() // consume operator
	stmt.Right = p.parseExpressionWithPrecedence(PREFIX)
	return stmt
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

// parseRunStatement parses a RUN statement
func (p *Parser) parseRunStatement() *RunStatement {
	return &RunStatement{
		Line: p.currentToken.Line,
	}
}

// parseStopStatement parses a STOP statement
func (p *Parser) parseStopStatement() *StopStatement {
	return &StopStatement{
		Line: p.currentToken.Line,
	}
}

// parseGotoStatement parses a GOTO statement
func (p *Parser) parseGotoStatement() *GotoStatement {
	stmt := &GotoStatement{Line: p.currentToken.Line}

	p.nextToken() // consume GOTO

	// Expect a number (target line)
	if p.currentToken.Type != lexer.NUMBER {
		p.addTokenError("line number", p.currentToken.Type)
		return nil
	}

	// Parse the target line number
	targetLine, err := strconv.Atoi(p.currentToken.Literal)
	if err != nil {
		p.addErrorf("invalid line number: %s", p.currentToken.Literal)
		return nil
	}

	stmt.TargetLine = targetLine
	return stmt
}

// parseIfStatement parses an IF...THEN statement
func (p *Parser) parseIfStatement() *IfStatement {
	stmt := &IfStatement{Line: p.currentToken.Line}

	p.nextToken() // consume IF

	// Parse the condition expression
	stmt.Condition = p.parseExpression()
	if stmt.Condition == nil {
		return nil
	}

	// For simple expressions without operators, we need to advance past the primary expression
	if p.currentToken.Type != lexer.THEN && p.peekToken.Type == lexer.THEN {
		p.nextToken()
	}

	// Expect THEN
	if p.currentToken.Type != lexer.THEN {
		p.addTokenError("THEN", p.currentToken.Type)
		return nil
	}

	p.nextToken() // consume THEN

	// Parse the statement to execute when condition is true
	stmt.ThenStmt = p.parseStatement()
	if stmt.ThenStmt == nil {
		return nil
	}

	return stmt
}

// parseComparisonExpression parses a comparison expression (left op right)
func (p *Parser) parseComparisonExpression() Expression {
	// Parse left side - for comparisons, we want to allow arithmetic expressions
	left := p.parseExpressionWithPrecedence(LOWEST) // Parse full expressions
	if left == nil {
		return nil
	}

	// Check if current token is a comparison operator
	operator := ""
	switch p.currentToken.Type {
	case lexer.ASSIGN: // = for comparison
		operator = "="
	case lexer.NE: // <>
		operator = "<>"
	case lexer.LT: // <
		operator = "<"
	case lexer.GT: // >
		operator = ">"
	case lexer.LE: // <=
		operator = "<="
	case lexer.GE: // >=
		operator = ">="
	default:
		p.addTokenError("comparison operator", p.currentToken.Type)
		return nil
	}

	line := p.currentToken.Line
	p.nextToken() // consume comparison operator

	// Parse right side as an expression
	right := p.parseExpressionWithPrecedence(LOWEST) // Parse full expressions
	if right == nil {
		return nil
	}

	return &ComparisonExpression{
		Left:     left,
		Operator: operator,
		Right:    right,
		Line:     line,
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

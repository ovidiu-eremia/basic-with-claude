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

	error *ParseError
}

// New creates a new parser instance
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:      l,
		precedence: NewPrecedenceTable(),
		error:      nil,
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
// Errors returns parsing errors as human-readable strings (kept for compatibility)
func (p *Parser) Errors() []string {
	if p.error == nil {
		return []string{}
	}
	return []string{fmt.Sprintf("line %d: %s", p.error.Position.Line, p.error.Message)}
}

// ParseError returns the parse error if any
func (p *Parser) ParseError() *ParseError {
	return p.error
}

// addErrorf adds a formatted error message with current token context
func (p *Parser) addErrorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	p.addErrorAt(p.currentToken.Line, msg)
}

// addTokenError adds an error message with token type context
func (p *Parser) addTokenError(expected string, got lexer.TokenType) {
	p.addErrorAt(p.currentToken.Line, fmt.Sprintf("expected %s, got %s (token %q)", expected, got, p.currentToken.Literal))
}

// addLiteralError adds an error message with token literal context
func (p *Parser) addLiteralError(prefix string, literal string) {
	p.addErrorAt(p.currentToken.Line, fmt.Sprintf("%s: %s", prefix, literal))
}

// addErrorAt sets a ParseError with an explicit line (only if no error exists yet)
func (p *Parser) addErrorAt(line int, msg string) {
	if p.error == nil {
		p.error = &ParseError{
			Message: msg,
			Position: lexer.Position{
				Line:   line,
				Column: 0, // Column tracking not implemented yet
			},
		}
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

		// Stop parsing if we encountered any error
		if p.error != nil {
			break
		}

		// parseLine() leaves us at NEWLINE or EOF, no need to advance
	}

	return program
}

// parseLine parses a single BASIC line
func (p *Parser) parseLine() *Line {
	if p.currentToken.Type != lexer.NUMBER {
		p.addTokenError("line number", p.currentToken.Type)
		return nil
	}

	lineNum, err := strconv.Atoi(p.currentToken.Literal)
	if err != nil {
		p.addLiteralError("invalid line number", p.currentToken.Literal)
		return nil
	}

	line := &Line{
		Number:     lineNum,
		Statements: []Statement{},
		SourceLine: p.currentToken.Line,
	}

	p.nextToken() // consume line number

	// Parse statements on this line. On first error, skip rest of the line.
	for p.currentToken.Type != lexer.NEWLINE && p.currentToken.Type != lexer.EOF {
		// Support colon-separated statements
		if p.currentToken.Type == lexer.COLON {
			p.nextToken()
			continue
		}
		stmt := p.parseStatement()
		if stmt == nil {
			// An error occurred; stop here
			break
		}
		line.Statements = append(line.Statements, stmt)
		// Advance token after parsing a successful statement
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
	case lexer.INPUT:
		return p.parseInputStatement()
	case lexer.END:
		return p.parseEndStatement()
	case lexer.RUN:
		return p.parseRunStatement()
	case lexer.STOP:
		return p.parseStopStatement()
	case lexer.GOTO:
		return p.parseGotoStatement()
	case lexer.GOSUB:
		return p.parseGosubStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.FOR:
		return p.parseForStatement()
	case lexer.NEXT:
		return p.parseNextStatement()
	case lexer.DATA:
		return p.parseDataStatement()
	case lexer.READ:
		return p.parseReadStatement()
	case lexer.REM:
		return p.parseRemStatement()
	case lexer.ILLEGAL:
		p.addLiteralError("illegal token", p.currentToken.Literal)
		return nil
	default:
		p.addTokenError("unrecognized statement", p.currentToken.Type)
		return nil
	}
}

// parseRemStatement parses a REM statement which consumes the rest of the line
func (p *Parser) parseRemStatement() *RemStatement {
	stmt := &RemStatement{Line: p.currentToken.Line}
	// Consume REM token
	p.nextToken()
	// Skip tokens until end of line or EOF, but leave currentToken on last non-NEWLINE token
	for p.peekToken.Type != lexer.NEWLINE && p.peekToken.Type != lexer.EOF {
		p.nextToken()
	}
	// Leave currentToken at NEWLINE/EOF so caller can advance appropriately
	return stmt
}

// parseDataStatement parses a DATA statement: DATA <const>[, <const>...]
func (p *Parser) parseDataStatement() *DataStatement {
	stmt := &DataStatement{Line: p.currentToken.Line}
	p.nextToken() // consume DATA

	// Parse zero or more constants until end of line/EOF
	for p.currentToken.Type != lexer.NEWLINE && p.currentToken.Type != lexer.EOF {
		var expr Expression
		switch p.currentToken.Type {
		case lexer.STRING:
			expr = p.parseStringLiteral()
		case lexer.NUMBER:
			expr = p.parseNumberLiteral()
		default:
			p.addTokenError("constant (number or string)", p.currentToken.Type)
			return nil
		}
		stmt.Values = append(stmt.Values, expr)

		// If next token is comma, consume it and continue
		if p.peekToken.Type == lexer.COMMA {
			p.nextToken() // move to COMMA
			p.nextToken() // move past COMMA to next value
			continue
		}

		// Otherwise, break and let outer loop advance
		break
	}
	return stmt
}

// parseReadStatement parses a READ statement: READ <var>[, <var>...]
func (p *Parser) parseReadStatement() *ReadStatement {
	stmt := &ReadStatement{Line: p.currentToken.Line}
	p.nextToken() // consume READ

	// Expect at least one identifier
	if p.currentToken.Type != lexer.IDENT {
		p.addTokenError("variable name", p.currentToken.Type)
		return nil
	}

	// Parse variables separated by commas until end of line/EOF
	for p.currentToken.Type == lexer.IDENT {
		stmt.Variables = append(stmt.Variables, p.currentToken.Literal)
		if p.peekToken.Type == lexer.COMMA {
			p.nextToken() // move to COMMA
			p.nextToken() // move to next IDENT
			continue
		}
		break
	}
	return stmt
}

// parsePrintStatement parses a PRINT statement
func (p *Parser) parsePrintStatement() *PrintStatement {
	stmt := &PrintStatement{Line: p.currentToken.Line}

	// Look ahead: if next token ends the statement, this is an empty PRINT
	if p.peekToken.Type == lexer.NEWLINE || p.peekToken.Type == lexer.EOF || p.peekToken.Type == lexer.COLON {
		// Empty PRINT -> outputs blank line
		stmt.Expression = &StringLiteral{Value: "", Line: stmt.Line}
		return stmt
	}

	// Consume PRINT and parse first expression
	p.nextToken()
	first := p.parseExpression()
	if first == nil {
		return nil
	}

	// Collect additional items separated by ';' or ','
	items := []Expression{first}
	noNewline := false
	for {
		// If next token is a separator, handle it
		if p.peekToken.Type == lexer.SEMICOLON || p.peekToken.Type == lexer.COMMA {
			p.nextToken() // move to separator
			// If the separator is the last token before end-of-statement, suppress newline
			if p.peekToken.Type == lexer.NEWLINE || p.peekToken.Type == lexer.EOF || p.peekToken.Type == lexer.COLON {
				noNewline = true
				break
			}
			// Otherwise parse another expression
			p.nextToken()
			nextExpr := p.parseExpression()
			if nextExpr == nil {
				return nil
			}
			items = append(items, nextExpr)
			continue
		}
		break
	}

	// If only one item and no special flags, keep legacy field for compatibility
	if len(items) == 1 && !noNewline {
		stmt.Expression = items[0]
	} else {
		stmt.Items = items
		stmt.NoNewline = noNewline
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
		operatorType := p.peekToken.Type
		operatorPrec := p.precedence.GetPrecedence(p.peekToken.Type)

		p.nextToken() // consume the operator
		p.nextToken() // move to right operand

		// Handle right associativity for power operator
		rightPrec := operatorPrec
		if operatorType == lexer.POWER {
			rightPrec = operatorPrec - 1 // Right associative
		}
		right := p.parseExpressionWithPrecedence(rightPrec)
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
		// Check if this is a function call (identifier followed by left parenthesis)
		if p.peekToken.Type == lexer.LPAREN {
			return p.parseFunctionCall()
		}
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
		p.addErrorAt(p.currentToken.Line, "expected ')' after grouped expression")
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

// parseGosubStatement parses a GOSUB statement
func (p *Parser) parseGosubStatement() *GosubStatement {
	stmt := &GosubStatement{Line: p.currentToken.Line}

	p.nextToken() // consume GOSUB

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

// parseReturnStatement parses a RETURN statement
func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Line: p.currentToken.Line}
	// No need to consume more tokens - RETURN is a simple statement
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

	// Support optional THEN if followed directly by GOTO (e.g., IF A=B GOTO 100)
	if p.peekToken.Type == lexer.GOTO {
		p.nextToken() // move to GOTO
		// Parse the statement to execute (GOTO ...)
		stmt.ThenStmt = p.parseStatement()
		if stmt.ThenStmt == nil {
			return nil
		}
		return stmt
	}

	// For simple expressions without operators, advance past the primary expression to THEN if needed
	if p.currentToken.Type != lexer.THEN && p.peekToken.Type == lexer.THEN {
		p.nextToken()
	}

	// Expect THEN for standard form
	if p.currentToken.Type != lexer.THEN {
		p.addTokenError("THEN", p.currentToken.Type)
		return nil
	}

	p.nextToken() // consume THEN

	// Support short form: THEN <lineNumber> meaning GOTO <lineNumber>
	if p.currentToken.Type == lexer.NUMBER {
		targetLine, err := strconv.Atoi(p.currentToken.Literal)
		if err != nil {
			p.addErrorf("invalid line number: %s", p.currentToken.Literal)
			return nil
		}
		stmt.ThenStmt = &GotoStatement{TargetLine: targetLine, Line: stmt.Line}
		return stmt
	}

	// Parse the statement to execute when condition is true (regular form)
	stmt.ThenStmt = p.parseStatement()
	if stmt.ThenStmt == nil {
		return nil
	}

	return stmt
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

// parseFunctionCall parses a function call (identifier followed by parentheses)
func (p *Parser) parseFunctionCall() *FunctionCall {
	functionCall := &FunctionCall{
		FunctionName: p.currentToken.Literal,
		Arguments:    []Expression{}, // Initialize empty slice
		Line:         p.currentToken.Line,
	}

	p.nextToken() // consume function name
	p.nextToken() // consume '('

	// Parse arguments
	if p.currentToken.Type != lexer.RPAREN {
		// First argument
		arg := p.parseExpression()
		if arg == nil {
			return nil
		}
		functionCall.Arguments = append(functionCall.Arguments, arg)

		// Additional arguments separated by commas
		for p.peekToken.Type == lexer.COMMA {
			p.nextToken() // consume current argument token
			p.nextToken() // consume ','
			arg = p.parseExpression()
			if arg == nil {
				return nil
			}
			functionCall.Arguments = append(functionCall.Arguments, arg)
		}

		// Move to the closing parenthesis for this call.
		// After parsing an argument, currentToken may still be on the last
		// token of the argument (or on a nested function's ')'). In both
		// simple and nested cases, if the next token is a ')', advance so
		// that currentToken points at this call's closing parenthesis.
		if p.peekToken.Type == lexer.RPAREN {
			p.nextToken()
		}
	}

	// We should now be at the closing parenthesis
	if p.currentToken.Type != lexer.RPAREN {
		p.addTokenError("')'", p.currentToken.Type)
		return nil
	}

	// Don't consume the closing parenthesis here - let the caller handle token advancement
	return functionCall
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

// parseInputStatement parses an INPUT statement
func (p *Parser) parseInputStatement() *InputStatement {
	stmt := &InputStatement{Line: p.currentToken.Line}
	p.nextToken() // consume INPUT

	// Check if we have a prompt string
	if p.currentToken.Type == lexer.STRING {
		stmt.Prompt = p.currentToken.Literal
		p.nextToken() // consume prompt string

		// Expect semicolon after prompt
		if p.currentToken.Type != lexer.SEMICOLON {
			p.addTokenError("semicolon", p.currentToken.Type)
			return nil
		}
		p.nextToken() // consume semicolon
	}

	// Expect variable name
	if p.currentToken.Type != lexer.IDENT {
		p.addTokenError("variable name", p.currentToken.Type)
		return nil
	}
	stmt.Variable = p.currentToken.Literal
	return stmt
}

// parseForStatement parses a FOR statement: FOR I = 1 TO 5 [STEP X]
func (p *Parser) parseForStatement() *ForStatement {
	stmt := &ForStatement{Line: p.currentToken.Line}

	p.nextToken() // consume FOR

	// Expect variable name
	if p.currentToken.Type != lexer.IDENT {
		p.addTokenError("variable name", p.currentToken.Type)
		return nil
	}
	stmt.Variable = p.currentToken.Literal

	p.nextToken() // consume variable name

	// Expect equals sign
	if p.currentToken.Type != lexer.ASSIGN {
		p.addTokenError("'=' after variable name", p.currentToken.Type)
		return nil
	}

	p.nextToken() // consume '='

	// Parse start value expression
	stmt.StartValue = p.parseExpression()
	if stmt.StartValue == nil {
		return nil
	}

	// For simple expressions without operators, we need to advance past the primary expression
	if p.currentToken.Type != lexer.TO && p.peekToken.Type == lexer.TO {
		p.nextToken()
	}

	// Expect TO keyword
	if p.currentToken.Type != lexer.TO {
		p.addTokenError("TO keyword", p.currentToken.Type)
		return nil
	}

	p.nextToken() // consume TO

	// Parse end value expression
	stmt.EndValue = p.parseExpression()
	if stmt.EndValue == nil {
		return nil
	}

	// For simple expressions without operators, we may need to advance to STEP if present
	if p.currentToken.Type != lexer.STEP && p.peekToken.Type == lexer.STEP {
		p.nextToken()
	}

	// Optional STEP clause
	if p.currentToken.Type == lexer.STEP {
		p.nextToken() // consume STEP
		stmt.StepValue = p.parseExpression()
		if stmt.StepValue == nil {
			return nil
		}
	}

	return stmt
}

// parseNextStatement parses a NEXT statement: NEXT I or NEXT
func (p *Parser) parseNextStatement() *NextStatement {
	stmt := &NextStatement{Line: p.currentToken.Line}

	p.nextToken() // consume NEXT

	// Check if there's a variable name (optional in NEXT)
	if p.currentToken.Type == lexer.IDENT {
		stmt.Variable = p.currentToken.Literal
		// Token will be consumed by the main parser loop
	}

	return stmt
}

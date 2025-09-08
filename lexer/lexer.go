// ABOUTME: Lexical analyzer for BASIC language tokens and keywords
// ABOUTME: Converts BASIC source code into a stream of tokens for parsing

package lexer

import "strings"

// TokenType represents the type of a token
type TokenType string

// Token types for BASIC language
const (
	ILLEGAL   TokenType = "ILLEGAL"
	EOF       TokenType = "EOF"
	NUMBER    TokenType = "NUMBER"
	STRING    TokenType = "STRING"
	IDENT     TokenType = "IDENT"
	ASSIGN    TokenType = "="
	PRINT     TokenType = "PRINT"
	LET       TokenType = "LET"
	END       TokenType = "END"
	RUN       TokenType = "RUN"
	STOP      TokenType = "STOP"
	GOTO      TokenType = "GOTO"
	INPUT     TokenType = "INPUT"
	DATA      TokenType = "DATA"
	READ      TokenType = "READ"
	NEWLINE   TokenType = "NEWLINE"
	PLUS      TokenType = "+"
	MINUS     TokenType = "-"
	MULTIPLY  TokenType = "*"
	DIVIDE    TokenType = "/"
	POWER     TokenType = "^"
	COLON     TokenType = ":"
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	COMMA     TokenType = ","
	IF        TokenType = "IF"
	THEN      TokenType = "THEN"
	GT        TokenType = ">"
	LT        TokenType = "<"
	NE        TokenType = "<>"
	GE        TokenType = ">="
	LE        TokenType = "<="
	SEMICOLON TokenType = ";"
	FOR       TokenType = "FOR"
	TO        TokenType = "TO"
	NEXT      TokenType = "NEXT"
	STEP      TokenType = "STEP"
	GOSUB     TokenType = "GOSUB"
	RETURN    TokenType = "RETURN"
	REM       TokenType = "REM"
)

// keywords maps BASIC keywords to their token types
var keywords = map[string]TokenType{
	"PRINT":  PRINT,
	"LET":    LET,
	"END":    END,
	"RUN":    RUN,
	"STOP":   STOP,
	"GOTO":   GOTO,
	"INPUT":  INPUT,
	"DATA":   DATA,
	"READ":   READ,
	"IF":     IF,
	"THEN":   THEN,
	"FOR":    FOR,
	"TO":     TO,
	"NEXT":   NEXT,
	"STEP":   STEP,
	"GOSUB":  GOSUB,
	"RETURN": RETURN,
	"REM":    REM,
}

// Position represents a position in the source code
type Position struct {
	Line   int
	Column int
}

// Token represents a single token with its type, literal value, and line number
type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

// Lexer represents the lexical analyzer
type Lexer struct {
	input           string
	currentPosition int  // current position in input (points to current char)
	nextPosition    int  // current reading position in input (after current char)
	currentChar     byte // current char under examination
	line            int  // current line number
}

// New creates a new lexer instance
func New(input string) *Lexer {
	lexer := &Lexer{
		input: input,
		line:  1,
	}
	lexer.readChar()
	return lexer
}

// createToken creates a token with the current line number
func (l *Lexer) createToken(tokenType TokenType, literal string) Token {
	return Token{Type: tokenType, Literal: literal, Line: l.line}
}

// createSingleCharToken creates a token for single-character operators and advances position
func (l *Lexer) createSingleCharToken(tokenType TokenType) Token {
	tok := l.createToken(tokenType, string(l.currentChar))
	l.readChar()
	return tok
}

// readChar reads the next character and advances the position
func (l *Lexer) readChar() {
	if l.nextPosition >= len(l.input) {
		l.currentChar = 0 // ASCII NUL represents "EOF"
	} else {
		l.currentChar = l.input[l.nextPosition]
	}
	l.currentPosition = l.nextPosition
	l.nextPosition++
}

// NextToken scans and returns the next token
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	switch l.currentChar {
	case '=':
		return l.createSingleCharToken(ASSIGN)
	case '+':
		return l.createSingleCharToken(PLUS)
	case '-':
		return l.createSingleCharToken(MINUS)
	case '*':
		return l.createSingleCharToken(MULTIPLY)
	case '/':
		return l.createSingleCharToken(DIVIDE)
	case '^':
		return l.createSingleCharToken(POWER)
	case '(':
		return l.createSingleCharToken(LPAREN)
	case ')':
		return l.createSingleCharToken(RPAREN)
	case ':':
		return l.createSingleCharToken(COLON)
	case ',':
		return l.createSingleCharToken(COMMA)
	case ';':
		return l.createSingleCharToken(SEMICOLON)
	case '<':
		return l.readComparisonOperator('<')
	case '>':
		return l.readComparisonOperator('>')
	case '"':
		if literal, terminated := l.readString(); terminated {
			return l.createToken(STRING, literal)
		} else {
			return l.createToken(ILLEGAL, "unterminated string")
		}
	case '.':
		// Support leading-dot decimals like .3
		if isDigit(l.peekChar()) {
			start := l.currentPosition
			l.readChar() // consume '.'
			for isDigit(l.currentChar) {
				l.readChar()
			}
			return l.createToken(NUMBER, l.input[start:l.currentPosition])
		}
		// Otherwise '.' is illegal in this grammar
		return l.createSingleCharToken(ILLEGAL)
	case '\n':
		tok := l.createToken(NEWLINE, string(l.currentChar))
		l.line++
		l.readChar()
		return tok
	case 0:
		return l.createToken(EOF, "")
	default:
		if isLetter(l.currentChar) {
			literal := l.readIdentifier()
			return l.createToken(lookupIdent(literal), literal)
		} else if isDigit(l.currentChar) {
			literal := l.readNumber()
			return l.createToken(NUMBER, literal)
		} else {
			tok := l.createToken(ILLEGAL, string(l.currentChar))
			l.readChar()
			return tok
		}
	}
}

// skipWhitespace skips whitespace characters except newlines
func (l *Lexer) skipWhitespace() {
	for l.currentChar == ' ' || l.currentChar == '\t' || l.currentChar == '\r' {
		l.readChar()
	}
}

// readString reads a string literal, returns (content, terminated)
func (l *Lexer) readString() (content string, terminated bool) {
	position := l.currentPosition + 1
	for {
		l.readChar()
		if l.currentChar == '"' || l.currentChar == 0 {
			break
		}
	}

	if l.currentChar == 0 {
		return "", false // Unterminated string
	}

	result := l.input[position:l.currentPosition]
	l.readChar() // Skip closing quote
	return result, true
}

// readIdentifier reads an identifier/keyword
func (l *Lexer) readIdentifier() string {
	position := l.currentPosition
	for isLetter(l.currentChar) || isDigit(l.currentChar) {
		l.readChar()
	}
	// Handle string variables ending with $
	if l.currentChar == '$' {
		l.readChar()
	}
	return l.input[position:l.currentPosition]
}

// readNumber reads a numeric literal
func (l *Lexer) readNumber() string {
	position := l.currentPosition
	for isDigit(l.currentChar) {
		l.readChar()
	}
	if l.currentChar == '.' {
		l.readChar()
		for isDigit(l.currentChar) {
			l.readChar()
		}
	}
	return l.input[position:l.currentPosition]
}

// readComparisonOperator reads comparison operators (< <= <> > >=)
func (l *Lexer) readComparisonOperator(firstChar byte) Token {
	switch firstChar {
	case '<':
		if l.peekChar() == '=' {
			l.readChar() // consume '<'
			l.readChar() // consume '='
			return l.createToken(LE, "<=")
		} else if l.peekChar() == '>' {
			l.readChar() // consume '<'
			l.readChar() // consume '>'
			return l.createToken(NE, "<>")
		} else {
			return l.createSingleCharToken(LT)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar() // consume '>'
			l.readChar() // consume '='
			return l.createToken(GE, ">=")
		} else {
			return l.createSingleCharToken(GT)
		}
	default:
		tok := l.createToken(ILLEGAL, string(firstChar))
		l.readChar()
		return tok
	}
}

// peekChar returns the next character without advancing position
func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	}
	return l.input[l.nextPosition]
}

// isLetter checks if character is a letter
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

// isDigit checks if character is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// lookupIdent checks if identifier is a keyword
func lookupIdent(ident string) TokenType {
	// Convert to uppercase for case-insensitive keyword matching
	upperIdent := strings.ToUpper(ident)
	if tok, ok := keywords[upperIdent]; ok {
		return tok
	}
	return IDENT // Non-keyword identifiers are now valid variable names
}

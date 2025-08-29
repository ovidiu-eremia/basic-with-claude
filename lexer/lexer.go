// ABOUTME: Lexical analyzer for BASIC language tokens and keywords
// ABOUTME: Converts BASIC source code into a stream of tokens for parsing

package lexer

// TokenType represents the type of a token
type TokenType string

// Token types for BASIC language
const (
	ILLEGAL  TokenType = "ILLEGAL"
	EOF      TokenType = "EOF"
	NUMBER   TokenType = "NUMBER"
	STRING   TokenType = "STRING"
	IDENT    TokenType = "IDENT"
	ASSIGN   TokenType = "="
	PRINT    TokenType = "PRINT"
	LET      TokenType = "LET"
	END      TokenType = "END"
	RUN      TokenType = "RUN"
	STOP     TokenType = "STOP"
	GOTO     TokenType = "GOTO"
	NEWLINE  TokenType = "NEWLINE"
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	MULTIPLY TokenType = "*"
	DIVIDE   TokenType = "/"
	POWER    TokenType = "^"
	LPAREN   TokenType = "("
	RPAREN   TokenType = ")"
)

// keywords maps BASIC keywords to their token types
var keywords = map[string]TokenType{
	"PRINT": PRINT,
	"LET":   LET,
	"END":   END,
	"RUN":   RUN,
	"STOP":  STOP,
	"GOTO":  GOTO,
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
	case '"':
		if literal, terminated := l.readString(); terminated {
			return l.createToken(STRING, literal)
		} else {
			return l.createToken(ILLEGAL, "unterminated string")
		}
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
	return l.input[position:l.currentPosition]
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
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT // Non-keyword identifiers are now valid variable names
}

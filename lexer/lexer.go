// ABOUTME: Lexical analyzer for BASIC language tokens and keywords
// ABOUTME: Converts BASIC source code into a stream of tokens for parsing

package lexer

// TokenType represents the type of a token
type TokenType string

// Token types for BASIC language
const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"
	NUMBER  TokenType = "NUMBER"
	STRING  TokenType = "STRING"
	IDENT   TokenType = "IDENT"
	ASSIGN  TokenType = "="
	PRINT   TokenType = "PRINT"
	LET     TokenType = "LET"
	END     TokenType = "END"
	NEWLINE TokenType = "NEWLINE"
)

// Token represents a single token with its type, literal value, and line number
type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

// Lexer represents the lexical analyzer
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line number
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

// readChar reads the next character and advances the position
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII NUL represents "EOF"
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// NextToken scans and returns the next token
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = Token{Type: ASSIGN, Literal: string(l.ch), Line: l.line}
		l.readChar()
	case '"':
		tok.Line = l.line
		if literal, ok := l.readString(); ok {
			tok.Type = STRING
			tok.Literal = literal
		} else {
			tok.Type = ILLEGAL
			tok.Literal = "unterminated string"
		}
	case '\n':
		tok = Token{Type: NEWLINE, Literal: string(l.ch), Line: l.line}
		l.line++
		l.readChar()
	case 0:
		tok.Literal = ""
		tok.Type = EOF
		tok.Line = l.line
	default:
		if isLetter(l.ch) {
			tok.Line = l.line
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = NUMBER
			tok.Literal = l.readNumber()
			tok.Line = l.line
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch), Line: l.line}
		}
		l.readChar()
	}

	return tok
}

// skipWhitespace skips whitespace characters except newlines
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

// readString reads a string literal, returns (content, success)
func (l *Lexer) readString() (string, bool) {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	
	if l.ch == 0 {
		return "", false // Unterminated string
	}
	
	result := l.input[position:l.position]
	l.readChar() // Skip closing quote
	return result, true
}

// readIdentifier reads an identifier/keyword
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a numeric literal
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
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
	keywords := map[string]TokenType{
		"PRINT": PRINT,
		"LET":   LET,
		"END":   END,
	}
	
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT // Non-keyword identifiers are now valid variable names
}
package lexer

import (
	"monkey-interpreter/token"
)

// Lexer is simply two pointers
type Lexer struct {
	input   string
	pos     int
	readPos int
	ch      byte
}

// New returns a new lexer
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// Now this only supports ASCII because we use bytes to represent chars
// TODO: use `rune` to support UTF-8
// readChar is called at NextToken to make sure l.ch is always up-to-date (an invariant!)
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		// reached end, return NUL char
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos // inv: l.pos always points to l.ch
	l.readPos++
}

// peekChar is different from readChar in that it returns instead of updates l.ch
// And also peekChar doesn't advance pointers
func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.pos
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.pos]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	position := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.pos]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// NextToken returns the next token and advances char
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			firstCh := l.ch
			l.readChar()
			tok = newTokenString(token.EQ, string(firstCh)+string(l.ch))
		} else {
			tok = newTokenByte(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newTokenByte(token.SEMICOLON, l.ch)
	case '(':
		tok = newTokenByte(token.LPAREN, l.ch)
	case ')':
		tok = newTokenByte(token.RPAREN, l.ch)
	case ',':
		tok = newTokenByte(token.COMMA, l.ch)
	case '+':
		tok = newTokenByte(token.PLUS, l.ch)
	case '-':
		tok = newTokenByte(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			firstCh := l.ch
			l.readChar()
			tok = newTokenString(token.NEQ, string(firstCh)+string(l.ch))
		} else {
			tok = newTokenByte(token.BANG, l.ch)
		}
	case '/':
		tok = newTokenByte(token.SLASH, l.ch)
	case '*':
		tok = newTokenByte(token.STAR, l.ch)
	case '<':
		tok = newTokenByte(token.LT, l.ch)
	case '>':
		tok = newTokenByte(token.GT, l.ch)
	case '{':
		tok = newTokenByte(token.LBRACE, l.ch)
	case '}':
		tok = newTokenByte(token.RBRACE, l.ch)
	case 0:
		tok = newTokenString(token.EOF, "")
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			return newTokenString(token.LookupIdent(literal), literal)
		} else if isDigit(l.ch) {
			return newTokenString(token.INT, l.readNumber())
		} else {
			tok = newTokenByte(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

// helper functions to construct new tokens
func newTokenByte(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
func newTokenString(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal}
}

package lexer

import (
	"kstmc.com/gosha/internal/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readCh()
	return l
}

func (l *Lexer) readCh() {
	l.ch = 0
	if l.readPosition < len(l.input) {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := string(l.ch)
			l.readCh()
			tok = token.Token{Type: token.EQ, Literal: ch + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := string(l.ch)
			l.readCh()
			tok = token.Token{Type: token.NEQ, Literal: ch + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '"':
		tok.Literal = l.readString()
		tok.Type = token.STRING
	case ':':
		if l.peekChar() == '=' {
			ch := string(l.ch)
			l.readCh()
			tok = token.Token{Type: token.INITASSIGN, Literal: ch + string(l.ch)}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case '>':
		tok = newToken(token.GT, l.ch)
	case '%':
		tok = newToken(token.PERCENT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '\n':
		tok = newToken(token.NLINE, l.ch)
	case '#':
		tok = newToken(token.HASH, l.ch)
	case '$':
		if l.peekChar() == '(' {
			l.readCh()
			tok.Literal = l.readBash()
			tok.Type = token.BASH
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '&':
		if l.peekChar() == '&' {
			ch := string(l.ch)
			l.readCh()
			tok = token.Token{Type: token.AND, Literal: ch + string(l.ch)}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := string(l.ch)
			l.readCh()
			tok = token.Token{Type: token.OR, Literal: ch + string(l.ch)}
		} else {
			tok = newToken(token.PIPE, l.ch)
		}
	case '-':
		if l.peekChar() == 'f' {
			ch := string(l.ch)
			l.readCh()
			tok = token.Token{Type: token.FOPER, Literal: ch + string(l.ch)}
		} else {
			tok = newToken(token.MINUS, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.FindIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readCh()
	return tok
}

func (l *Lexer) readBash() string {
	position := l.position + 1
	l.readCh()
	braceCount := 1
	for braceCount > 0 && l.ch != 0 {
		switch l.ch {
		case '(':
			braceCount++
		case ')':
			braceCount--
		}

		if braceCount > 0 && l.ch != 0 {
			l.readCh()
		} else {
			break
		}
	}

	return l.input[position:l.position]
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readString() string {
	position := l.position + 1
	l.readCh()
	for l.ch != '"' && l.ch != 0 {
		l.readCh()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readCh()
	}

	return l.input[position:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readCh()
	}

	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) skipWhitespace() {
	for {
		switch l.ch {
		case ' ':
			l.readCh()
		case '\t':
			l.readCh()
		//case '\n':
		//	l.readCh()
		case '\r':
			l.readCh()
		default:
			return
		}
	}
}

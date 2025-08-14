package token

import (
	"os"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT    = "IDENT"
	INT      = "INT"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	HASH     = "#"

	LT  = "<"
	GT  = ">"
	EQ  = "=="
	NEQ = "!="
	AND = "&&"
	OR  = "||"

	STRING = `"`

	ASSIGN     = "="
	INITASSIGN = ":="
	FOPER      = "-f"
	PLUS       = "+"
	MINUS      = "-"
	PERCENT    = "%"
	COMMA      = ","

	LPAREN = "("
	RPAREN = ")"

	LBRACE = "{"
	RBRACE = "}"
	PIPE   = "|"

	NLINE = "\n"

	BASH = "$"

	// Keywords

	FUNCTION = "FUNC"
	VAR      = "VAR"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	DTYPE    = "DTYPE"
	CALL     = "CALL"
	BUILDIN  = "BUILDIN"
	FOR      = "FOR"
)

var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"var":    VAR,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"string": DTYPE,
	"int":    DTYPE,
	"bool":   DTYPE,
	"any":    DTYPE,
	"for":    FOR,
}

func SetupBashCalls() error {
	files, err := os.ReadDir("/usr/bin/")
	if err != nil {
		return err
	}

	for _, file := range files {
		keywords[file.Name()] = BASH
	}

	return nil
}

func FindIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}

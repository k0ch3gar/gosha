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

	REF = "&"

	STRING = `"`

	ASSIGN       = "="
	INITASSIGN   = ":="
	FOPER        = "-f"
	PLUS         = "+"
	MINUS        = "-"
	PERCENT      = "%"
	COMMA        = ","
	CHANOPERATOR = "<-"

	LBRACKET = "["
	RBRACKET = "]"

	LPAREN = "("
	RPAREN = ")"

	LBRACE = "{"
	RBRACE = "}"
	PIPE   = "|"

	NLINE = "\n"

	BASHEXPR = "$()"
	BASHVAR  = "$"

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
	GO       = "GO"
	CHAN     = "CHAN"
	BREAK    = "BREAK"
)

var keywords = map[string]TokenType{
	"go":     GO,
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
	"chan":   CHAN,
	"break":  BREAK,
}

func SetupBashCalls() error {
	files, err := os.ReadDir("/usr/bin/")
	if err != nil {
		return err
	}

	for _, file := range files {
		keywords[file.Name()] = BASHEXPR
	}

	return nil
}

func FindIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}

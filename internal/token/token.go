package token

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

	ASSIGN     = "="
	INITASSIGN = ":="
	PLUS       = "+"
	MINUS      = "-"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"

	LBRACE = "{"
	RBRACE = "}"
	PIPE   = "|"

	NLINE = "\n"

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
	"echo":   CALL,
}

func FindIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}

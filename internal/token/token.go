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

	LT  = "<"
	GT  = ">"
	EQ  = "=="
	NEQ = "!="

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

	// Keywords

	FUNCTION = "FUNC"
	VAR      = "VAR"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"var":    VAR,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func FindIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}

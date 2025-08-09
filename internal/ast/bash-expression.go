package ast

import (
	"bytes"
	"strings"

	"kstmc.com/gosha/internal/token"
)

type BashExpression struct {
	Token token.Token

	Value []string
}

func (be *BashExpression) expressionNode() {

}

func (be *BashExpression) TokenLiteral() string {
	return be.Token.Literal
}

func (be *BashExpression) String() string {
	var out bytes.Buffer

	out.WriteString(be.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(be.Value, " "))
	out.WriteString(")")

	return out.String()
}

package ast

import (
	"bytes"
	"kstmc.com/gosha/internal/token"
)

type BashExpression struct {
	Token token.Token
	Value string
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
	out.WriteString(be.Value)
	out.WriteString(")")

	return out.String()
}

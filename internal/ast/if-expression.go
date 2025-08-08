package ast

import (
	"bytes"

	"kstmc.com/gosha/internal/token"
)

type IfStatement struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfStatement) statementNode() {

}

func (ie *IfStatement) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ie.Token.Literal)
	out.WriteString(" ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

package ast

import "kstmc.com/gosha/internal/token"

type StringLiteral struct {
	Token token.Token
	Value string
}

func (se *StringLiteral) expressionNode() {

}

func (se *StringLiteral) TokenLiteral() string {
	return se.Token.Literal
}

func (se *StringLiteral) String() string {
	return se.Value
}

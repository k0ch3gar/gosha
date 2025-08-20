package ast

import "kstmc.com/gosha/internal/token"

type BashVarExpression struct {
	Token token.Token

	Value string
}

func (bve *BashVarExpression) expressionNode() {

}

func (bve *BashVarExpression) TokenLiteral() string {
	return bve.Token.Literal
}

func (bve *BashVarExpression) String() string {
	return "$" + bve.Value
}

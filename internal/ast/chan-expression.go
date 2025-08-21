package ast

import (
	"kstmc.com/gosha/internal/token"
)

type SendChanStatement struct {
	Token       token.Token
	Source      Expression
	Destination *Identifier
}

func (sce *SendChanStatement) statementNode() {

}

func (sce *SendChanStatement) TokenLiteral() string {
	return sce.Token.Literal
}

func (sce *SendChanStatement) String() string {
	return sce.Destination.String() + " <- " + sce.Source.String()
}

type ReadChanExpression struct {
	Token  token.Token
	Source *Identifier
}

func (rce *ReadChanExpression) expressionNode() {

}

func (rce *ReadChanExpression) TokenLiteral() string {
	return rce.Token.Literal
}

func (rce *ReadChanExpression) String() string {
	return "<- " + rce.Source.String()
}

package ast

import "kstmc.com/gosha/internal/token"

type AssignStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (as *AssignStatement) statementNode() {

}

func (as *AssignStatement) TokenLiteral() string {
	return as.Token.Literal
}

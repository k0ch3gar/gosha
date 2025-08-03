package ast

import (
	"bytes"
	"kstmc.com/gosha/internal/token"
)

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

func (as *AssignStatement) String() string {
	var out bytes.Buffer

	out.WriteString(as.Name.String() + " ")
	out.WriteString(as.TokenLiteral() + " ")

	if as.Value != nil {
		out.WriteString(as.Value.String())
	}

	return out.String()
}

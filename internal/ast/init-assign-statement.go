package ast

import (
	"bytes"
	"kstmc.com/gosha/internal/token"
)

type InitAssignStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (is *InitAssignStatement) statementNode() {

}

func (is *InitAssignStatement) TokenLiteral() string {
	return is.Token.Literal
}

func (is *InitAssignStatement) String() string {
	var out bytes.Buffer

	out.WriteString(is.Name.String() + " ")
	out.WriteString(is.TokenLiteral() + " ")

	if is.Value != nil {
		out.WriteString(is.Value.String())
	}

	return out.String()
}

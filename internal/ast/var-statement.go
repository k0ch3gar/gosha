package ast

import (
	"bytes"

	"kstmc.com/gosha/internal/token"
)

type VarStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (vs *VarStatement) statementNode() {

}

func (vs *VarStatement) TokenLiteral() string {
	return vs.Token.Literal
}

func (vs *VarStatement) String() string {
	var out bytes.Buffer

	out.WriteString(vs.Token.Literal)
	out.WriteString(" ")
	out.WriteString(vs.Name.String())
	out.WriteString(" = ")
	if vs.Value != nil {
		out.WriteString(vs.Value.String())
	}

	return out.String()
}

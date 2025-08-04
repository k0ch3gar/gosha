package ast

import (
	"bytes"
	"kstmc.com/gosha/internal/token"
	"strings"
)

type FunctionLiteral struct {
	Token      token.Token
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
	ReturnType *DataType
}

func (fl *FunctionLiteral) expressionNode() {

}

func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	if fl.Name != nil {
		out.WriteString(" " + fl.Name.Value)
	}

	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.ReturnType.Token.Literal)
	out.WriteString(" ")
	out.WriteString(fl.Body.String())

	return out.String()
}

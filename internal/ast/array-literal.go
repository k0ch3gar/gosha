package ast

import (
	"bytes"
	"strings"

	"kstmc.com/gosha/internal/token"
)

type SliceLiteral struct {
	Token  token.Token
	Type   DataType
	Values []Expression
}

func (sl *SliceLiteral) expressionNode() {

}

func (sl *SliceLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *SliceLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("[]")
	out.WriteString(sl.Type.Name())
	out.WriteString("{")

	var values []string
	for _, value := range sl.Values {
		values = append(values, value.String())
	}

	out.WriteString(strings.Join(values, ", "))
	out.WriteString("}")

	return out.String()
}

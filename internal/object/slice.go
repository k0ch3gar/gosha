package object

import (
	"bytes"
	"strings"

	"kstmc.com/gosha/internal/ast"
)

type SliceObject struct {
	Values    []Object
	ValueType ast.DataType
}

func (so *SliceObject) Type() ast.DataType {
	return so.ValueType
}

func (so *SliceObject) Inspect() string {
	var out bytes.Buffer

	out.WriteString("[")
	var values []string
	for _, value := range so.Values {
		values = append(values, value.Inspect())
	}

	out.WriteString(strings.Join(values, ", "))
	out.WriteString("]")

	return out.String()
}

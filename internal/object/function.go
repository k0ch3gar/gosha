package object

import (
	"bytes"
	"strings"

	"kstmc.com/gosha/internal/ast"
)

type Function struct {
	Name       *ast.Identifier
	Parameters []*ast.Identifier
	ReturnType ast.DataType
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ast.DataType {
	return &ast.FunctionDataType{
		Parameters: func() []ast.DataType {
			var params []ast.DataType
			for _, param := range f.Parameters {
				params = append(params, *param.DataType)
			}

			return params
		}(),
		ReturnValue: f.ReturnType,
	}
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("func")
	if f.Name != nil {
		out.WriteString(" " + f.Name.Value)
	}

	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.ReturnType.Name())
	out.WriteString(" {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

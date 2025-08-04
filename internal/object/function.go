package object

import (
	"bytes"
	"kstmc.com/gosha/internal/ast"
	"strings"
)

type Function struct {
	Name       *ast.Identifier
	Parameters []*ast.Identifier
	ReturnType *ast.DataType
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
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
	out.WriteString(f.ReturnType.Name)
	out.WriteString(" {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

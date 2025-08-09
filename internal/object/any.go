package object

import (
	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/parser"
)

type Any struct {
	Value string
}

func (a *Any) Inspect() string {
	return a.Value
}

func (a *Any) Type() ast.DataType {
	return parser.ANY
}

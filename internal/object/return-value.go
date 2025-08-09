package object

import (
	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/parser"
)

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ast.DataType {
	return parser.RETURN
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

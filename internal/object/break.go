package object

import (
	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/parser"
)

type BreakObject struct {
}

func (bo *BreakObject) Inspect() string {
	return "break"
}

func (bo *BreakObject) Type() ast.DataType {
	return parser.BREAK
}

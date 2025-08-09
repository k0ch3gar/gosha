package object

import (
	"fmt"

	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/parser"
)

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ast.DataType {
	return parser.BOOLEAN
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

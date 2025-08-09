package object

import (
	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/parser"
)

type Nil struct {
}

func (n *Nil) Type() ast.DataType {
	return parser.NIL
}

func (n *Nil) Inspect() string {
	return "nil"
}

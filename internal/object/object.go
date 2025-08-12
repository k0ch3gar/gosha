package object

import "kstmc.com/gosha/internal/ast"

type Object interface {
	Type() ast.DataType
	Inspect() string
}

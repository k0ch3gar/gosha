package object

import (
	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/parser"
)

type Error struct {
	Message string
}

func (e *Error) Type() ast.DataType {
	return parser.ERROR
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

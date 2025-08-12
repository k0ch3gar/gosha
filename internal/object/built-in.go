package object

import (
	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/parser"
)

//type BuiltinFunction func(env *Environment, args ...Object) Object

type Builtin struct {
	FnName string
}

func (bi *Builtin) Inspect() string {
	return bi.FnName
}

func (bi *Builtin) Type() ast.DataType {
	return parser.BUILTIN
}

var Builtins = map[string]*Builtin{
	"print": {
		FnName: "print",
	},
	"read": {
		FnName: "read",
	},
}

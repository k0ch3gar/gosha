package object

import (
	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/parser"
)

//type BuiltinFunction func(env *Environment, args ...Object) Object

type Builtin struct {
	FnName     string
	Parameters []*ast.Identifier
	ReturnType ast.DataType
}

func (bi *Builtin) Inspect() string {
	return bi.FnName
}

func (bi *Builtin) Type() ast.DataType {
	dType := &ast.BuiltinDataType{
		ReturnType: bi.ReturnType,
	}

	for _, ident := range bi.Parameters {
		dType.Parameters = append(dType.Parameters, *ident.DataType)
	}

	return dType
}

var Builtins = map[string]*Builtin{
	"print": {
		FnName:     "print",
		ReturnType: parser.NIL,
	},
	"read": {
		FnName:     "read",
		ReturnType: parser.NIL,
	},
	"append": {
		FnName:     "append",
		ReturnType: &ast.SliceDataType{},
	},
	"len": {
		FnName:     "len",
		ReturnType: parser.INT,
	},
}

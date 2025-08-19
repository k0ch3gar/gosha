package object

import "kstmc.com/gosha/internal/ast"

type ReferenceObject struct {
	Value *Object
}

func (ro *ReferenceObject) Inspect() string {
	return "&" + (*ro.Value).Inspect()
}

func (ro *ReferenceObject) Type() ast.DataType {
	return &ast.ReferenceDataType{
		ValueType: (*ro.Value).Type(),
	}
}

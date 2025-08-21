package object

import "kstmc.com/gosha/internal/ast"

type ChanObject struct {
	Chan     chan Object
	ChanType ast.DataType
}

func (co *ChanObject) Inspect() string {
	return ""
}

func (co *ChanObject) Type() ast.DataType {
	return &ast.ChanDataType{ValueType: co.ChanType}
}

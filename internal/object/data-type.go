package object

import "kstmc.com/gosha/internal/ast"

type DataTypeObject struct {
	DataType ast.DataType
}

func (dto *DataTypeObject) Inspect() string {
	return dto.DataType.Name()
}

func (dto *DataTypeObject) Type() ast.DataType {
	return dto.DataType
}

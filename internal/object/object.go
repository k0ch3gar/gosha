package object

import "kstmc.com/gosha/internal/ast"

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NIL_OBJ          = "NIL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	ANY_OBJ          = "ANY"
	STRING_OBJ       = "STRING"
)

type Object interface {
	Type() ast.DataType
	Inspect() string
}

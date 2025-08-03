package object

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NIL_OBJ     = "NIL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

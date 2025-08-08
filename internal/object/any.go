package object

type Any struct {
	Value string
}

func (a *Any) Inspect() string {
	return a.Value
}

func (a *Any) Type() ObjectType {
	return ANY_OBJ
}

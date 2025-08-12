package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func UnwrapEnvironment(env *Environment) *Environment {
	return env.outer
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

func (e *Environment) Update(name string, value Object) Object {
	if _, ok := e.store[name]; ok {
		return e.Set(name, value)
	} else if e.outer != nil {
		return e.outer.Update(name, value)
	}

	return value
}

func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}

func (e *Environment) Contains(name string) bool {
	_, ok := e.store[name]
	if !ok && e.outer != nil {
		_, ok = e.outer.Get(name)
	}
	return ok
}

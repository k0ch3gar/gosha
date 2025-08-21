package object

import (
	"fmt"
	"strconv"

	"kstmc.com/gosha/internal/ast"
	"kstmc.com/gosha/internal/parser"
)

type BuiltinFunc func(...Object) Object

type Builtin struct {
	Name string
	Fn   BuiltinFunc
	//FnName     string
	//Parameters []*ast.Identifier
	//ReturnType ast.DataType
}

func (bi *Builtin) Inspect() string {
	return bi.Name
}

func (bi *Builtin) Type() ast.DataType {
	return parser.BUILTIN
}

var Builtins = map[string]*Builtin{
	"print": {
		Name: "print",
		Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect() + " ")
			}

			fmt.Println()
			return &Nil{}
		},
	},
	"len": {
		Name: "len",
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Nil{}
			}

			switch arg := args[0].(type) {
			case *SliceObject:
				return &Integer{Value: int64(len(arg.Values))}
			default:
				return &Nil{}
			}
		},
	},
	"append": {
		Name: "append",
		Fn: func(args ...Object) Object {
			if len(args) < 1 {
				return &Nil{}
			}

			slice, ok := args[0].(*SliceObject)
			if !ok {
				return &Nil{}
			}

			for _, arg := range args[1:] {
				if arg.Type().Name() != slice.Type().(*ast.SliceDataType).Type.Name() {
					return &Nil{}
				}
			}

			newSlice := &SliceObject{
				ValueType: slice.ValueType,
				Values:    append(slice.Values, args[1:]...),
			}

			return newSlice
		},
	},
	"read": {
		Name: "read",
		Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return &Nil{}
			}

			ref, ok := args[0].(*ReferenceObject)
			if !ok {
				return &Nil{}
			}

			arg := ref.Value
			switch arg := (*arg).(type) {
			case *Integer:
				fmt.Scan(&arg.Value)
			case *String:
				fmt.Scan(&arg.Value)
			case *Boolean:
				fmt.Scan(&arg.Value)
			}

			return &Nil{}
		},
	},
	"make": {
		Name: "make",
		Fn: func(args ...Object) Object {
			arg, ok := args[0].(*DataTypeObject)
			if !ok {
				return &Error{Message: fmt.Sprintf("cannot make %s", arg)}
			}
			switch arg := arg.DataType.(type) {
			case *ast.IntegerDataType:
				return makeIntegerObject(args[1:])
			case *ast.StringDataType:
				return makeStringObject(args[1:])
			case *ast.ChanDataType:
				return makeChanObject(args)
			default:
				return &Error{Message: fmt.Sprintf("cannot make %s", arg)}
			}
		},
	},
}

func makeChanObject(objects []Object) Object {
	if len(objects) != 2 {
		return &Error{Message: fmt.Sprintf("unexpected amount of arguments. expected 2, provided %d", len(objects)+1)}
	}

	val, ok := objects[1].(*Integer)
	if !ok {
		return &Error{Message: fmt.Sprintf("expected Integer argument to make chan, got=%T", objects[0])}
	}

	ch := make(chan Object, val.Value)
	return &ChanObject{Chan: ch, ChanType: objects[0].(*DataTypeObject).DataType.(*ast.ChanDataType).ValueType}
}

func makeStringObject(objects []Object) Object {
	if len(objects) != 1 {
		return &Error{Message: fmt.Sprintf("unexpected amount of arguments. expected 2, provided %d", len(objects)+1)}
	}

	switch obj := objects[0].(type) {
	case *Integer:
		return &String{Value: strconv.FormatInt(obj.Value, 10)}
	case *String:
		return &String{Value: obj.Value}
	default:
		return &Error{Message: fmt.Sprintf("cannot make Integer from %s", obj.Inspect())}
	}
}

func makeIntegerObject(objects []Object) Object {
	if len(objects) != 1 {
		return &Error{Message: fmt.Sprintf("unexpected amount of arguments. expected 2, provided %d", len(objects)+1)}
	}

	switch obj := objects[0].(type) {
	case *Integer:
		return &Integer{Value: obj.Value}
	case *String:
		val, err := strconv.Atoi(obj.Value)
		if err != nil {
			return &Error{Message: err.Error()}
		}

		return &Integer{Value: int64(val)}
	default:
		return &Error{Message: fmt.Sprintf("cannot make Integer from %s", obj.Inspect())}
	}
}

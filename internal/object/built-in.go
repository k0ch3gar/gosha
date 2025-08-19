package object

import (
	"fmt"

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
			}

			return &Nil{}
		},
	},
}

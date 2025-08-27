package evaluator

import (
	"github.com/lancelote/writing-an-interpreter-in-go/object"
	"os"
)

var builtins = map[string]*object.Builtin{
	"exit": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError("`exit()` doesn't accept arguments")
			}

			os.Exit(0)
			return nil
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("`first` accepts 1 argument, got %d", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` should be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, want 1, got %d", len(args))
			}

			switch arg := args[0].(type) {

			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}

			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("`last` accepts 1 argument, got %d", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `last` should be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[len(arr.Elements)-1]
			}

			return NULL
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("`push` accepts 2 arguments, got %d", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("first argument to `push` should be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("`rest` accepts 1 argument, got %d", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `rest` should be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}

			return NULL
		},
	},
}

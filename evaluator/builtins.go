package evaluator

import (
	"github.com/lancelote/writing-an-interpreter-in-go/object"
	"os"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, want 1, got %d", len(args))
			}

			switch arg := args[0].(type) {

			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"exit": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError("`exit()` doesn't accept arguments")
			}

			os.Exit(0)
			return nil
		},
	},
}

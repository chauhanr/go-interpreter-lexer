package evaluator

import (
	"go-interpreter-lexer/object"
	"fmt"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{Fn: func(args ...object.Object) object.Object{
			if len(args) != 1{
				return newError("wrong number of arguments. got %d, want=1", len(args))
			}
			switch arg := args[0].(type){
				case *object.Array:
					return &object.Integer{Value: int64(len(arg.Elements))}
				case *object.String:
					return &object.Integer{Value: int64(len(arg.Value))}
				default :
					return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first" : &object.Builtin{ Fn: func(args ...object.Object) object.Object{
		if len(args) != 1{
			return newError("wrong number of arguments to first function got %d, wanted 1", len(args))
		}
		if args[0].Type() != object.ARRAY_OBJ{
			return newError("argument for `first` function is suppose to be an object.Array but got %T", args[0])
		}
		arr := args[0].(*object.Array)
		if len(arr.Elements) > 0{
			return arr.Elements[0]
		}
		return NULL
	  },
	},
	"last" : &object.Builtin{ Fn: func(args ...object.Object) object.Object{
		if len(args) != 1{
			return newError("wrong number of arguments to `last` function got %d, wanted 1", len(args))
		}
		if args[0].Type() != object.ARRAY_OBJ{
			return newError("argument for `last` function is suppose to be an object.Array but got %T", args[0])
		}
		arr := args[0].(*object.Array)
		if len(arr.Elements) >0 {
			return arr.Elements[len(arr.Elements) -1]
		}
		return NULL
	  },
	},
	"puts": &object.Builtin{Fn: func(args ...object.Object) object.Object{
			for _, arg := range args{
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}

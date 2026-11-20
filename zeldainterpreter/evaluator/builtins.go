package evaluator

import (
	"fmt"
	"io"
	"os"

	"zeldainterpreter/object"
)

// Output es donde escribe la función integrada puts. Por defecto es la
// consola; el runner y las pruebas la cambian para capturar la salida.
var Output io.Writer = os.Stdout

// builtins son funciones escritas en Go que se pueden llamar desde Monkey.
// Existen porque hay cosas que el lenguaje no puede hacer por sí mismo,
// como medir un string o imprimir en pantalla.
var builtins = map[string]*object.Builtin{

	// len(x): longitud de un string (en bytes) o de un array.
	"len": {Fn: func(args ...object.Object) object.Object {
		if err := checkArgCount(args, 1); err != nil {
			return err
		}
		switch arg := args[0].(type) {
		case *object.String:
			return &object.Integer{Value: int64(len(arg.Value))}
		case *object.Array:
			return &object.Integer{Value: int64(len(arg.Elements))}
		default:
			return newError("argument to `len` not supported, got %s", args[0].Type())
		}
	}},

	// first(arr): el primer elemento, o null si está vacío.
	"first": {Fn: func(args ...object.Object) object.Object {
		arr, err := arrayArg("first", args, 1)
		if err != nil {
			return err
		}
		if len(arr.Elements) == 0 {
			return NULL
		}
		return arr.Elements[0]
	}},

	// last(arr): el último elemento, o null si está vacío.
	"last": {Fn: func(args ...object.Object) object.Object {
		arr, err := arrayArg("last", args, 1)
		if err != nil {
			return err
		}
		if len(arr.Elements) == 0 {
			return NULL
		}
		return arr.Elements[len(arr.Elements)-1]
	}},

	// rest(arr): un array NUEVO con todos los elementos menos el primero.
	"rest": {Fn: func(args ...object.Object) object.Object {
		arr, err := arrayArg("rest", args, 1)
		if err != nil {
			return err
		}
		if len(arr.Elements) == 0 {
			return NULL
		}
		elems := make([]object.Object, len(arr.Elements)-1)
		copy(elems, arr.Elements[1:]) // copiamos: el original no se modifica
		return &object.Array{Elements: elems}
	}},

	// push(arr, x): un array NUEVO con x agregado al final.
	"push": {Fn: func(args ...object.Object) object.Object {
		arr, err := arrayArg("push", args, 2)
		if err != nil {
			return err
		}
		elems := make([]object.Object, len(arr.Elements), len(arr.Elements)+1)
		copy(elems, arr.Elements)
		return &object.Array{Elements: append(elems, args[1])}
	}},

	// puts(a, b, ...): imprime cada argumento en una línea. Devuelve null.
	"puts": {Fn: func(args ...object.Object) object.Object {
		for _, arg := range args {
			fmt.Fprintln(Output, arg.Inspect())
		}
		return NULL
	}},
}

// checkArgCount valida la cantidad de argumentos de un built-in.
func checkArgCount(args []object.Object, want int) *object.Error {
	if len(args) != want {
		return newError("wrong number of arguments. got=%d, want=%d", len(args), want)
	}
	return nil
}

// arrayArg valida la cantidad de argumentos y que el primero sea un array.
func arrayArg(name string, args []object.Object, want int) (*object.Array, object.Object) {
	if err := checkArgCount(args, want); err != nil {
		return nil, err
	}
	arr, ok := args[0].(*object.Array)
	if !ok {
		return nil, newError("argument to `%s` must be ARRAY, got %s", name, args[0].Type())
	}
	return arr, nil
}

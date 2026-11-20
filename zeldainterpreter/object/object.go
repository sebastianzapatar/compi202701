// Package object define los VALORES que existen mientras un programa se
// ejecuta: enteros, booleanos, strings, funciones, arrays, hashes...
//
// Todo valor implementa la interfaz Object. Así el evaluador puede devolver
// siempre el mismo tipo (object.Object) sin importar qué produjo el código.
package object

import (
	"fmt"
	"hash/fnv"
	"strings"

	"zeldainterpreter/ast"
)

// ObjectType es el nombre del tipo de un valor. Aparece en los mensajes de
// error, por ejemplo: "type mismatch: INTEGER + BOOLEAN".
type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	STRING_OBJ       = "STRING"
	NULL_OBJ         = "NULL"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	RETURN_VALUE_OBJ = "RETURN_VALUE" // envoltorio interno, no visible para quien programa
	ERROR_OBJ        = "ERROR"        // error en tiempo de ejecución
)

// Object es el contrato que cumplen todos los valores.
type Object interface {
	Type() ObjectType // el nombre del tipo: "INTEGER", "STRING"...
	Inspect() string  // cómo se muestra el valor en la consola
}

// ============================================================================
// Valores simples
// ============================================================================

// Integer es un número entero de 64 bits.
type Integer struct{ Value int64 }

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// Boolean es true o false. El evaluador crea UNA sola instancia de cada uno
// (singletons TRUE y FALSE) y las reutiliza.
type Boolean struct{ Value bool }

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// String es un texto.
type String struct{ Value string }

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// Null representa "ausencia de valor", por ejemplo el resultado de un if
// falso sin else. También existe una única instancia (NULL).
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// ============================================================================
// Objetos de control (los usa el evaluador internamente)
// ============================================================================

// ReturnValue envuelve el valor de un `return` para que "suba" por los
// bloques anidados sin que se ejecute nada más. Se desenvuelve al llegar al
// nivel del programa o de la función.
type ReturnValue struct{ Value Object }

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// Error representa un error en tiempo de ejecución. Igual que ReturnValue,
// detiene la evaluación y sube hasta el nivel más alto.
type Error struct{ Message string }

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// ============================================================================
// Funciones
// ============================================================================

// Function es una función definida en Monkey con fn(...) { ... }.
//
// El campo Env es la clave de los CLOSURES: guarda el entorno donde la
// función fue CREADA, así puede seguir usando esas variables aunque se
// llame desde otro lugar.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	params := make([]string, len(f.Parameters))
	for i, p := range f.Parameters {
		params[i] = p.String()
	}
	return "fn(" + strings.Join(params, ", ") + ") {\n" + f.Body.String() + "\n}"
}

// BuiltinFunction es la firma de las funciones integradas escritas en Go.
// Reciben cualquier cantidad de argumentos (función variádica).
type BuiltinFunction func(args ...Object) Object

// Builtin envuelve una función integrada (len, puts, push...).
type Builtin struct{ Fn BuiltinFunction }

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

// ============================================================================
// Colecciones
// ============================================================================

// Array es una lista ordenada de valores. En Monkey los arrays son
// INMUTABLES: push y rest devuelven arrays nuevos.
type Array struct{ Elements []Object }

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	parts := make([]string, len(a.Elements))
	for i, e := range a.Elements {
		parts[i] = e.Inspect()
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// HashKey es la llave que usamos internamente en el map de Go.
//
// ¿Por qué no usar el Object directamente como llave? Porque Object es un
// puntero: dos &String{"nombre"} distintos tienen direcciones distintas y no
// coincidirían. HashKey es un struct de VALORES, que Go compara campo a campo.
// Incluimos el Type para que el entero 1 y el booleano true no choquen.
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// Hashable lo implementan los tipos que pueden usarse como llave de un hash.
type Hashable interface {
	Object
	HashKey() HashKey
}

func (i *Integer) HashKey() HashKey { return HashKey{Type: i.Type(), Value: uint64(i.Value)} }

func (b *Boolean) HashKey() HashKey {
	var v uint64
	if b.Value {
		v = 1
	}
	return HashKey{Type: b.Type(), Value: v}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a() // función hash FNV-1a de la biblioteca estándar
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// HashPair guarda la llave ORIGINAL (para poder imprimirla) junto al valor.
type HashPair struct {
	Key   Object
	Value Object
}

// Hash es un diccionario llave → valor.
// Keys conserva el orden de inserción para que Inspect sea siempre igual.
type Hash struct {
	Pairs map[HashKey]HashPair
	Keys  []HashKey
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	parts := make([]string, 0, len(h.Keys))
	for _, k := range h.Keys {
		pair := h.Pairs[k]
		parts = append(parts, pair.Key.Inspect()+": "+pair.Value.Inspect())
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

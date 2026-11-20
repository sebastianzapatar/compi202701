// Polimorfismo con interfaces: muchos tipos, un mismo contrato.
package main

import (
	"fmt"
	"strings"
)

type Object interface {
	Type() string
	Inspect() string
}

type Integer struct{ Value int64 }
type Boolean struct{ Value bool }
type Array struct{ Elements []Object }

func (i *Integer) Type() string    { return "INTEGER" }
func (i *Integer) Inspect() string { return fmt.Sprint(i.Value) }
func (b *Boolean) Type() string    { return "BOOLEAN" }
func (b *Boolean) Inspect() string { return fmt.Sprint(b.Value) }
func (a *Array) Type() string      { return "ARRAY" }
func (a *Array) Inspect() string {
	parts := make([]string, len(a.Elements))
	for i, e := range a.Elements {
		parts[i] = e.Inspect() // cada elemento sabe imprimirse a sí mismo
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// Verificación en tiempo de compilación: si un tipo deja de cumplir Object, no compila.
var (
	_ Object = (*Integer)(nil)
	_ Object = (*Boolean)(nil)
	_ Object = (*Array)(nil)
)

func truthy(o Object) bool {
	switch v := o.(type) { // type switch: decidir según el tipo concreto
	case *Boolean:
		return v.Value
	default:
		return true
	}
}

func main() {
	objs := []Object{
		&Integer{Value: 5},
		&Boolean{Value: false},
		&Array{Elements: []Object{&Integer{Value: 1}, &Boolean{Value: true}}},
	}
	for _, o := range objs {
		fmt.Printf("%-7s %-12s truthy=%v\n", o.Type(), o.Inspect(), truthy(o))
	}
}

// Prototype: crear objetos copiando (clonando) uno existente.
// Ejemplo: push en Monkey no modifica el array original: clona y luego agrega.
package main

import "fmt"

type Object interface {
	Inspect() string
	Clone() Object
}

type Integer struct{ Value int64 }

func (i *Integer) Inspect() string { return fmt.Sprint(i.Value) }
func (i *Integer) Clone() Object {
	c := *i
	return &c
}

type Array struct{ Elements []Object }

func (a *Array) Inspect() string { return fmt.Sprintf("%d elementos", len(a.Elements)) }

// Clone hace una copia PROFUNDA: también clona cada elemento.
func (a *Array) Clone() Object {
	elems := make([]Object, len(a.Elements))
	for i, e := range a.Elements {
		elems[i] = e.Clone()
	}
	return &Array{Elements: elems}
}

func push(arr *Array, v Object) *Array {
	c := arr.Clone().(*Array) // el original queda intacto
	c.Elements = append(c.Elements, v)
	return c
}

func main() {
	original := &Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 2}}}
	nuevo := push(original, &Integer{Value: 3})

	nuevo.Elements[0].(*Integer).Value = 99 // modificar la copia...
	fmt.Println("original:", original.Inspect(), "primero =", original.Elements[0].Inspect())
	fmt.Println("nuevo:   ", nuevo.Inspect(), "primero =", nuevo.Elements[0].Inspect())
}

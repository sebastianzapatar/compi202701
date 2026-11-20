// "Herencia" en Go: NO existe. Se usa EMBEDDING (composición con promoción
// de métodos) para reutilizar código, e interfaces para el polimorfismo.
package main

import "fmt"

// BaseNode tiene comportamiento reutilizable.
type BaseNode struct {
	Literal string
}

func (b BaseNode) TokenLiteral() string { return b.Literal }
func (b BaseNode) Kind() string         { return "genérico" }
func (b BaseNode) Describe() string     { return "nodo " + b.Kind() } // llama a SU Kind

// Integer embebe BaseNode: sus métodos se PROMUEVEN a Integer.
type Integer struct {
	BaseNode
	Value int64
}

// Kind oculta (shadowing) al método promovido. No es un override virtual.
func (i Integer) Kind() string { return "entero" }

type Node interface {
	TokenLiteral() string
	Kind() string
}

func main() {
	n := Integer{BaseNode: BaseNode{Literal: "5"}, Value: 5}

	fmt.Println(n.TokenLiteral())  // promovido desde BaseNode
	fmt.Println(n.Kind())          // método propio de Integer
	fmt.Println(n.BaseNode.Kind()) // acceso explícito al embebido
	fmt.Println(n.Describe())      // ¡"nodo genérico"! No hay despacho virtual

	var node Node = n // Integer satisface Node gracias a los métodos promovidos
	fmt.Println(node.Kind())
}

// Composite: tratar igual a objetos individuales y a composiciones de objetos.
// ¡El AST es un Composite! Un bloque contiene sentencias que contienen expresiones.
package main

import (
	"fmt"
	"strings"
)

type Node interface {
	String() string
	Count() int // cuántos nodos tiene este subárbol
}

// Hojas
type Integer struct{ Value int }
type Ident struct{ Name string }

func (i Integer) String() string { return fmt.Sprint(i.Value) }
func (i Integer) Count() int     { return 1 }
func (i Ident) String() string   { return i.Name }
func (i Ident) Count() int       { return 1 }

// Compuestos
type Infix struct {
	Left  Node
	Op    string
	Right Node
}

func (n Infix) String() string {
	return "(" + n.Left.String() + " " + n.Op + " " + n.Right.String() + ")"
}
func (n Infix) Count() int { return 1 + n.Left.Count() + n.Right.Count() }

type Block struct{ Stmts []Node }

func (b Block) String() string {
	parts := make([]string, len(b.Stmts))
	for i, s := range b.Stmts {
		parts[i] = s.String()
	}
	return "{ " + strings.Join(parts, "; ") + " }"
}

func (b Block) Count() int {
	total := 1
	for _, s := range b.Stmts {
		total += s.Count() // no importa si s es hoja o compuesto
	}
	return total
}

func main() {
	tree := Block{Stmts: []Node{
		Infix{Left: Integer{1}, Op: "+", Right: Infix{Left: Integer{2}, Op: "*", Right: Ident{"x"}}},
		Ident{"y"},
	}}
	fmt.Println(tree, "->", tree.Count(), "nodos")
}

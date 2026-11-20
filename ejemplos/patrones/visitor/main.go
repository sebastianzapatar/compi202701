// Visitor: agregar operaciones nuevas a una jerarquía de nodos SIN modificar los nodos.
// El libro usa un type switch en Eval; el Visitor es la alternativa clásica de POO.
package main

import (
	"fmt"
	"strconv"
)

type Visitor interface {
	VisitInt(n *IntLit) any
	VisitInfix(n *Infix) any
}

type Node interface{ Accept(v Visitor) any }

type IntLit struct{ Value int }

type Infix struct {
	Left  Node
	Op    string
	Right Node
}

func (n *IntLit) Accept(v Visitor) any { return v.VisitInt(n) }
func (n *Infix) Accept(v Visitor) any  { return v.VisitInfix(n) }

// Operación 1: evaluar.
type Evaluator struct{}

func (Evaluator) VisitInt(n *IntLit) any { return n.Value }
func (e Evaluator) VisitInfix(n *Infix) any {
	l, r := n.Left.Accept(e).(int), n.Right.Accept(e).(int)
	switch n.Op {
	case "+":
		return l + r
	case "*":
		return l * r
	}
	return 0
}

// Operación 2: imprimir. Es nueva y no requirió tocar IntLit ni Infix.
type Printer struct{}

func (Printer) VisitInt(n *IntLit) any { return strconv.Itoa(n.Value) }
func (p Printer) VisitInfix(n *Infix) any {
	return "(" + n.Left.Accept(p).(string) + " " + n.Op + " " + n.Right.Accept(p).(string) + ")"
}

func main() {
	ast := &Infix{Left: &IntLit{1}, Op: "+", Right: &Infix{Left: &IntLit{2}, Op: "*", Right: &IntLit{3}}}
	fmt.Println(ast.Accept(Printer{}), "=", ast.Accept(Evaluator{}))
}

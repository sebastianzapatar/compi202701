package ast

import (
	"fmt"
	"strings"
)

// Tree dibuja el AST como un árbol de texto, con el rol de cada hijo.
// Lo usa el modo "arbol" de la consola. Ejemplo para  let x = 1 + 2;
//
//	Program
//	└── LetStatement
//	    ├── nombre: Identifier x
//	    └── valor: InfixExpression +
//	        ├── izq: IntegerLiteral 1
//	        └── der: IntegerLiteral 2
func Tree(n Node) string {
	var b strings.Builder
	label, children := describe(n)
	b.WriteString(label + "\n")
	writeChildren(&b, children, "")
	return b.String()
}

// child es un hijo de un nodo junto con su rol ("izq", "condición"...).
type child struct {
	role string
	node Node
}

// writeChildren dibuja los hijos con las ramas ├── y └──.
// `prefix` acumula las líneas verticales de los niveles superiores.
func writeChildren(b *strings.Builder, children []child, prefix string) {
	for i, c := range children {
		isLast := i == len(children)-1
		branch, nextPrefix := "├── ", prefix+"│   "
		if isLast {
			branch, nextPrefix = "└── ", prefix+"    "
		}
		label, grandChildren := describe(c.node)
		if c.role != "" {
			label = c.role + ": " + label
		}
		b.WriteString(prefix + branch + label + "\n")
		writeChildren(b, grandChildren, nextPrefix)
	}
}

// describe devuelve la etiqueta de un nodo y la lista de sus hijos.
// Es un "type switch": decide qué hacer según el tipo concreto del nodo,
// la misma técnica que usa el evaluador.
func describe(n Node) (string, []child) {
	switch n := n.(type) {
	case *Program:
		return "Program", statements(n.Statements, "")
	case *LetStatement:
		return "LetStatement", []child{{"nombre", n.Name}, {"valor", n.Value}}
	case *ReturnStatement:
		return "ReturnStatement", []child{{"valor", n.ReturnValue}}
	case *ExpressionStatement:
		return "ExpressionStatement", []child{{"", n.Expression}}
	case *BlockStatement:
		return "BlockStatement", statements(n.Statements, "")
	case *Identifier:
		return "Identifier " + n.Value, nil
	case *IntegerLiteral:
		return fmt.Sprintf("IntegerLiteral %d", n.Value), nil
	case *StringLiteral:
		return fmt.Sprintf("StringLiteral %q", n.Value), nil
	case *Boolean:
		return fmt.Sprintf("Boolean %t", n.Value), nil
	case *PrefixExpression:
		return "PrefixExpression " + n.Operator, []child{{"operando", n.Right}}
	case *InfixExpression:
		return "InfixExpression " + n.Operator, []child{{"izq", n.Left}, {"der", n.Right}}
	case *IfExpression:
		kids := []child{{"condición", n.Condition}, {"entonces", n.Consequence}}
		if n.Alternative != nil {
			kids = append(kids, child{"sino", n.Alternative})
		}
		return "IfExpression", kids
	case *FunctionLiteral:
		var kids []child
		for _, p := range n.Parameters {
			kids = append(kids, child{"parámetro", p})
		}
		return "FunctionLiteral", append(kids, child{"cuerpo", n.Body})
	case *CallExpression:
		kids := []child{{"función", n.Function}}
		for _, a := range n.Arguments {
			kids = append(kids, child{"argumento", a})
		}
		return "CallExpression", kids
	case *ArrayLiteral:
		var kids []child
		for i, e := range n.Elements {
			kids = append(kids, child{fmt.Sprintf("[%d]", i), e})
		}
		return "ArrayLiteral", kids
	case *IndexExpression:
		return "IndexExpression", []child{{"colección", n.Left}, {"índice", n.Index}}
	case *HashLiteral:
		var kids []child
		for _, p := range n.Pairs {
			kids = append(kids, child{"llave", p.Key}, child{"valor", p.Value})
		}
		return "HashLiteral", kids
	case nil:
		return "(vacío)", nil
	default:
		return fmt.Sprintf("%T", n), nil
	}
}

func statements(stmts []Statement, role string) []child {
	kids := make([]child, len(stmts))
	for i, s := range stmts {
		kids[i] = child{role, s}
	}
	return kids
}

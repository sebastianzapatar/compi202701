package ast

import (
	"testing"

	"zeldainterpreter/token"
)

// TestStringAndTree construye a mano el AST de  let x = -y;  y verifica
// sus dos representaciones de texto.
func TestStringAndTree(t *testing.T) {
	program := &Program{Statements: []Statement{
		&LetStatement{
			Token: token.Token{Type: token.LET, Literal: "let"},
			Name:  &Identifier{Token: token.Token{Type: token.IDENT, Literal: "x"}, Value: "x"},
			Value: &PrefixExpression{
				Token:    token.Token{Type: token.MINUS, Literal: "-"},
				Operator: "-",
				Right:    &Identifier{Token: token.Token{Type: token.IDENT, Literal: "y"}, Value: "y"},
			},
		},
	}}

	if got := program.String(); got != "let x = (-y);" {
		t.Errorf("String() = %q", got)
	}

	want := "Program\n" +
		"└── LetStatement\n" +
		"    ├── nombre: Identifier x\n" +
		"    └── valor: PrefixExpression -\n" +
		"        └── operando: Identifier y\n"
	if got := Tree(program); got != want {
		t.Errorf("Tree() =\n%s\nse esperaba:\n%s", got, want)
	}
}

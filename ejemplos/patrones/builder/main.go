// Builder: construir un objeto complejo paso a paso.
// Ejemplo: armar programas/AST de prueba con una API fluida en los tests del parser.
package main

import (
	"errors"
	"fmt"
	"strings"
)

type Program struct{ Statements []string }

func (p Program) String() string { return strings.Join(p.Statements, "\n") }

type ProgramBuilder struct {
	stmts []string
	err   error
}

func NewProgram() *ProgramBuilder { return &ProgramBuilder{} }

func (b *ProgramBuilder) Let(name, value string) *ProgramBuilder {
	if name == "" {
		b.err = errors.New("let sin nombre")
		return b
	}
	b.stmts = append(b.stmts, fmt.Sprintf("let %s = %s;", name, value))
	return b
}

func (b *ProgramBuilder) Return(value string) *ProgramBuilder {
	b.stmts = append(b.stmts, fmt.Sprintf("return %s;", value))
	return b
}

// Build valida y entrega el producto terminado.
func (b *ProgramBuilder) Build() (Program, error) {
	if b.err != nil {
		return Program{}, b.err
	}
	return Program{Statements: b.stmts}, nil
}

func main() {
	p, err := NewProgram().Let("x", "5").Let("y", "x * 2").Return("y").Build()
	fmt.Println(p)
	fmt.Println("err:", err)

	_, err = NewProgram().Let("", "1").Return("1").Build()
	fmt.Println("err:", err)
}

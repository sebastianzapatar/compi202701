// Facade: una interfaz simple frente a un subsistema complejo.
// Ejemplo: Interpreter.Run esconde lexer, parser y evaluador detrás de un solo método.
package main

import (
	"fmt"
	"strconv"
	"strings"
)

// --- Subsistema (muy simplificado) ---
type lexer struct{}
type parser struct{}
type evaluator struct{ env map[string]int }

func (lexer) Tokens(src string) []string { return strings.Fields(src) }

func (parser) ParseLet(toks []string) (name, value string, err error) {
	if len(toks) != 4 || toks[0] != "let" || toks[2] != "=" {
		return "", "", fmt.Errorf("se esperaba: let <nombre> = <valor>")
	}
	return toks[1], toks[3], nil
}

func (e *evaluator) Bind(name, value string) (int, error) {
	n, err := strconv.Atoi(value)
	if err != nil {
		v, ok := e.env[value]
		if !ok {
			return 0, fmt.Errorf("identifier not found: %s", value)
		}
		n = v
	}
	e.env[name] = n
	return n, nil
}

// --- Fachada ---
type Interpreter struct {
	l lexer
	p parser
	e *evaluator
}

func New() *Interpreter { return &Interpreter{e: &evaluator{env: map[string]int{}}} }

func (in *Interpreter) Run(src string) (string, error) {
	name, value, err := in.p.ParseLet(in.l.Tokens(src))
	if err != nil {
		return "", err
	}
	n, err := in.e.Bind(name, value)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s = %d", name, n), nil
}

func main() {
	z := New() // el cliente solo conoce New y Run
	for _, src := range []string{"let x = 5", "let y = x", "let z = w", "x + 1"} {
		out, err := z.Run(src)
		if err != nil {
			fmt.Printf("%-10s -> error: %v\n", src, err)
			continue
		}
		fmt.Printf("%-10s -> %s\n", src, out)
	}
}

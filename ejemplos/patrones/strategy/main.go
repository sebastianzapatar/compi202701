// Strategy: una familia de algoritmos intercambiables detrás de una misma interfaz.
// El Pratt parser lo usa: cada token tiene su "estrategia" de parsing (prefixParseFns).
package main

import (
	"errors"
	"fmt"
)

// Estrategia como tipo función (la forma más idiomática en Go).
type BinaryOp func(a, b int64) (int64, error)

var operators = map[string]BinaryOp{
	"+": func(a, b int64) (int64, error) { return a + b, nil },
	"-": func(a, b int64) (int64, error) { return a - b, nil },
	"*": func(a, b int64) (int64, error) { return a * b, nil },
	"/": func(a, b int64) (int64, error) {
		if b == 0 {
			return 0, errors.New("division by zero")
		}
		return a / b, nil
	},
}

func evalInfix(op string, a, b int64) (int64, error) {
	strategy, ok := operators[op]
	if !ok {
		return 0, fmt.Errorf("unknown operator: INTEGER %s INTEGER", op)
	}
	return strategy(a, b)
}

// Estrategia como interfaz: motores de ejecución intercambiables.
type Engine interface{ Name() string }
type TreeWalker struct{}
type BytecodeVM struct{}

func (TreeWalker) Name() string { return "tree-walking (este curso)" }
func (BytecodeVM) Name() string { return "bytecode + VM (la secuela del libro)" }

func main() {
	for _, op := range []string{"+", "*", "/", "%"} {
		v, err := evalInfix(op, 12, 0)
		if err != nil {
			fmt.Printf("12 %s 0 -> ERROR: %v\n", op, err)
			continue
		}
		fmt.Printf("12 %s 0 -> %d\n", op, v)
	}
	for _, e := range []Engine{TreeWalker{}, BytecodeVM{}} {
		fmt.Println("motor:", e.Name())
	}
}

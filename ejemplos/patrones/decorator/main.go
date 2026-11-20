// Decorator: agregar responsabilidades envolviendo un objeto, sin modificar su código.
// Ejemplo: agregar traza y medición de tiempo a un evaluador.
package main

import (
	"fmt"
	"strings"
	"time"
)

type Evaluator interface {
	Eval(expr string) (int, error)
}

// Componente base: suma enteros "a + b + c".
type sumEvaluator struct{}

func (sumEvaluator) Eval(expr string) (int, error) {
	total := 0
	for _, p := range strings.Split(expr, "+") {
		var n int
		if _, err := fmt.Sscan(strings.TrimSpace(p), &n); err != nil {
			return 0, fmt.Errorf("no es entero: %q", strings.TrimSpace(p))
		}
		total += n
	}
	return total, nil
}

// Decorador 1: traza cada llamada.
type tracing struct{ inner Evaluator }

func (t tracing) Eval(expr string) (int, error) {
	fmt.Printf("-> Eval(%q)\n", expr)
	v, err := t.inner.Eval(expr)
	fmt.Printf("<- %d, err=%v\n", v, err)
	return v, err
}

// Decorador 2: mide el tiempo.
type timing struct{ inner Evaluator }

func (t timing) Eval(expr string) (int, error) {
	start := time.Now()
	defer func() { fmt.Println("   (medido:", time.Since(start) >= 0, ")") }()
	return t.inner.Eval(expr)
}

func main() {
	// Los decoradores se apilan como capas: tracing(timing(base)).
	var ev Evaluator = tracing{inner: timing{inner: sumEvaluator{}}}
	ev.Eval("1 + 2 + 3")
	ev.Eval("1 + x")
}

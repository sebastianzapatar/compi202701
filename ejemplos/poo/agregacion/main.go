// Agregación vs composición.
//   - Agregación: el todo USA una parte que vive fuera de él (se la entregan).
//     Ej.: el REPL y su io.Writer, que puede compartirse con otros.
//   - Composición: el todo CREA y es dueño de sus partes; mueren con él.
//     Ej.: el REPL y su Environment.
package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Environment struct{ store map[string]int }

func NewEnvironment() *Environment { return &Environment{store: map[string]int{}} }

type REPL struct {
	out io.Writer    // agregación: compartido, no le pertenece
	env *Environment // composición: lo crea y es su dueño
}

func NewREPL(out io.Writer) *REPL {
	return &REPL{out: out, env: NewEnvironment()}
}

func (r *REPL) Let(name string, v int) {
	r.env.store[name] = v
	fmt.Fprintf(r.out, "%s = %d\n", name, v)
}

func main() {
	var log strings.Builder
	shared := io.MultiWriter(os.Stdout, &log)

	a := NewREPL(shared)
	b := NewREPL(shared) // ambos REPL agregan el MISMO writer
	a.Let("x", 1)
	b.Let("x", 2)

	// Cada REPL compuso su propio entorno: no se mezclan.
	fmt.Println("x en a:", a.env.store["x"], "| x en b:", b.env.store["x"])

	a, b = nil, nil // los REPL y sus entornos quedan sin referencias (basura)...
	fmt.Print("...pero el writer agregado sigue vivo:\n", log.String())
}

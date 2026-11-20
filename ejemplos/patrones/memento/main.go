// Memento: guardar y restaurar el estado de un objeto sin romper su encapsulación.
// Ejemplo: si un programa falla a la mitad, el REPL restaura el entorno anterior.
package main

import (
	"errors"
	"fmt"
	"maps"
)

type Environment struct{ store map[string]int }

// Snapshot es el memento: una instantánea opaca (campo privado).
type Snapshot struct{ store map[string]int }

func (e *Environment) Save() Snapshot     { return Snapshot{store: maps.Clone(e.store)} }
func (e *Environment) Restore(s Snapshot) { e.store = maps.Clone(s.store) }

func runProgram(env *Environment, stmts []string) error {
	for _, s := range stmts {
		if s == "ERROR" {
			return errors.New("type mismatch: INTEGER + BOOLEAN")
		}
		env.store[s] = len(env.store) + 1
	}
	return nil
}

func main() {
	env := &Environment{store: map[string]int{"x": 1}}

	snap := env.Save() // el REPL (caretaker) guarda antes de ejecutar
	err := runProgram(env, []string{"a", "b", "ERROR", "c"})
	fmt.Println("tras el error:", env.store, "|", err)

	env.Restore(snap)
	fmt.Println("restaurado:   ", env.store)
}

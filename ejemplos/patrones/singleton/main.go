// Singleton: una única instancia compartida. En Go: variables de paquete o sync.Once.
package main

import (
	"fmt"
	"sync"
)

type Boolean struct{ Value bool }

// Instancias únicas, como TRUE/FALSE/NULL en el evaluador del libro.
var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

func nativeBool(b bool) *Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

// Registro de built-ins creado de forma perezosa, una sola vez y seguro entre goroutines.
type Registry struct{ names []string }

var (
	registry *Registry
	once     sync.Once
)

func Builtins() *Registry {
	once.Do(func() {
		fmt.Println("(inicializando registro: esto se imprime una sola vez)")
		registry = &Registry{names: []string{"len", "first", "last", "rest", "push", "puts"}}
	})
	return registry
}

func main() {
	fmt.Println("mismo TRUE:", nativeBool(1 < 2) == nativeBool(3 > 0))

	var wg sync.WaitGroup
	for range 3 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = Builtins()
		}()
	}
	wg.Wait()
	fmt.Println("misma instancia:", Builtins() == Builtins(), Builtins().names)
}

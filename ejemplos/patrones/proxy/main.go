// Proxy: un sustituto que controla el acceso al objeto real.
// Ejemplos: entorno con constantes protegidas y caché (memoización) de una función costosa.
package main

import (
	"errors"
	"fmt"
)

type Env interface {
	Get(name string) (int, bool)
	Set(name string, v int) error
}

type realEnv struct{ store map[string]int }

func (e *realEnv) Get(n string) (int, bool) {
	v, ok := e.store[n]
	return v, ok
}

func (e *realEnv) Set(n string, v int) error {
	e.store[n] = v
	return nil
}

// protectedEnv es un proxy de protección: impide reasignar constantes.
type protectedEnv struct {
	real   Env
	consts map[string]bool
}

func (p protectedEnv) Get(n string) (int, bool) { return p.real.Get(n) }
func (p protectedEnv) Set(n string, v int) error {
	if p.consts[n] {
		return errors.New("no se puede reasignar la constante " + n)
	}
	return p.real.Set(n, v)
}

// cachedFib es un proxy de caché: evita recalcular.
type cachedFib struct {
	calls int
	memo  map[int]int
}

func (c *cachedFib) Fib(n int) int {
	if v, ok := c.memo[n]; ok {
		return v
	}
	c.calls++
	v := n
	if n >= 2 {
		v = c.Fib(n-1) + c.Fib(n-2)
	}
	c.memo[n] = v
	return v
}

func main() {
	var env Env = protectedEnv{
		real:   &realEnv{store: map[string]int{"PI": 3}},
		consts: map[string]bool{"PI": true},
	}
	fmt.Println("Set x: ", env.Set("x", 1))
	fmt.Println("Set PI:", env.Set("PI", 4))

	f := &cachedFib{memo: map[int]int{}}
	fmt.Println("fib(40) =", f.Fib(40), "con", f.calls, "cálculos reales")
}

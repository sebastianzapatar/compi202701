// Flyweight: compartir objetos inmutables cuando hay muchísimos iguales.
// Ejemplo: caché de enteros pequeños (como CPython) e "interning" de identificadores.
package main

import "fmt"

type Integer struct{ Value int64 } // debe tratarse como INMUTABLE

const minCached, maxCached = -5, 256

var smallInts = func() []*Integer {
	cache := make([]*Integer, maxCached-minCached+1)
	for i := range cache {
		cache[i] = &Integer{Value: int64(i + minCached)}
	}
	return cache
}()

// NewInteger devuelve una instancia compartida para valores pequeños.
func NewInteger(v int64) *Integer {
	if v >= minCached && v <= maxCached {
		return smallInts[v-minCached]
	}
	return &Integer{Value: v}
}

// Interner guarda un único string por nombre de identificador.
type Interner struct{ pool map[string]*string }

func (in *Interner) Intern(s string) *string {
	if p, ok := in.pool[s]; ok {
		return p
	}
	p := &s
	in.pool[s] = p
	return p
}

func main() {
	fmt.Println("7 compartido:   ", NewInteger(7) == NewInteger(7))
	fmt.Println("1000 compartido:", NewInteger(1000) == NewInteger(1000))

	in := &Interner{pool: map[string]*string{}}
	a, b := in.Intern("contador"), in.Intern("contador")
	fmt.Println("mismo string:", a == b, "| tamaño del pool:", len(in.pool))
}

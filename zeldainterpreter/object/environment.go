package object

import "sort"

// Environment guarda las variables (los "bindings"): asocia cada nombre con
// su valor.
//
// Los entornos se ENCADENAN con el campo outer. Cuando se llama una función
// se crea un entorno nuevo cuyo outer es el entorno donde la función fue
// definida. Buscar un nombre es mirar en el entorno actual y, si no está,
// preguntar hacia afuera:
//
//	[ llamada: y=3 ] --outer--> [ closure: x=2 ] --outer--> [ global ]
type Environment struct {
	store map[string]Object
	outer *Environment // nil en el entorno global
}

// NewEnvironment crea el entorno global (sin outer).
func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}

// NewEnclosedEnvironment crea un entorno "encerrado" dentro de otro.
// Se usa en cada llamada a función.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get busca un nombre en este entorno y, si no está, en los de afuera.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return obj, ok
}

// Set guarda un valor SIEMPRE en este entorno (nunca en outer).
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// Names devuelve, ordenados, los nombres definidos en ESTE entorno.
// Lo usa el comando :env de la consola.
func (e *Environment) Names() []string {
	names := make([]string, 0, len(e.store))
	for name := range e.store {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

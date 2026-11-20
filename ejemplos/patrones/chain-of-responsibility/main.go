// Chain of Responsibility: pasar una petición por una cadena de manejadores
// hasta que uno la atienda. Ejemplo: comandos del REPL antes de evaluar código.
package main

import (
	"fmt"
	"strings"
)

type Handler interface {
	Handle(line string) string
	SetNext(h Handler) Handler
}

// base implementa el "pasar al siguiente"; los manejadores concretos la embeben.
type base struct{ next Handler }

func (b *base) SetNext(h Handler) Handler {
	b.next = h
	return h
}

func (b *base) pass(line string) string {
	if b.next == nil {
		return "(nadie atendió la línea)"
	}
	return b.next.Handle(line)
}

type emptyLine struct{ base }

type command struct {
	base
	name, reply string
}

type evaluate struct{ base }

func (h *emptyLine) Handle(l string) string {
	if strings.TrimSpace(l) == "" {
		return "(nada que evaluar)"
	}
	return h.pass(l)
}

func (h *command) Handle(l string) string {
	if strings.TrimSpace(l) == h.name {
		return h.reply
	}
	return h.pass(l)
}

func (h *evaluate) Handle(l string) string { return "evaluando: " + l }

func main() {
	chain := &emptyLine{}
	chain.SetNext(&command{name: ":help", reply: "comandos: :help :env"}).
		SetNext(&command{name: ":env", reply: "x = 5"}).
		SetNext(&evaluate{})

	for _, in := range []string{"", ":help", ":env", "1 + 2"} {
		fmt.Printf("%-8q -> %s\n", in, chain.Handle(in))
	}
}

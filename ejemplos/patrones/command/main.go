// Command: encapsular una petición como objeto para guardarla, deshacerla o rehacerla.
// Ejemplo: historial del REPL con :undo.
package main

import "fmt"

type Env map[string]int

type Command interface {
	Execute(env Env)
	Undo(env Env)
}

type LetCommand struct {
	Name    string
	Value   int
	prev    int
	existed bool
}

func (c *LetCommand) Execute(env Env) {
	c.prev, c.existed = env[c.Name]
	env[c.Name] = c.Value
}

func (c *LetCommand) Undo(env Env) {
	if c.existed {
		env[c.Name] = c.prev
	} else {
		delete(env, c.Name)
	}
}

type History struct {
	env  Env
	done []Command
}

func (h *History) Run(c Command) {
	c.Execute(h.env)
	h.done = append(h.done, c)
}

func (h *History) Undo() {
	if len(h.done) == 0 {
		return
	}
	last := h.done[len(h.done)-1]
	h.done = h.done[:len(h.done)-1]
	last.Undo(h.env)
}

func main() {
	h := &History{env: Env{}}
	h.Run(&LetCommand{Name: "x", Value: 1})
	h.Run(&LetCommand{Name: "x", Value: 2})
	h.Run(&LetCommand{Name: "y", Value: 9})
	fmt.Println("antes: ", h.env)
	h.Undo() // deshace let y
	h.Undo() // deshace let x = 2
	fmt.Println("después:", h.env)
}

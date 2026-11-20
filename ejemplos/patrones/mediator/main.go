// Mediator: los componentes no se hablan entre sí; se comunican a través de un mediador.
// Ejemplo: mini IDE con editor, consola y panel de errores coordinados por la sesión.
package main

import "fmt"

type Mediator interface {
	Notify(sender, event, data string)
}

type Editor struct{ m Mediator }
type Console struct{ lines []string }
type ErrorPanel struct{ errors []string }

func (e *Editor) Submit(code string) { e.m.Notify("editor", "submit", code) }
func (c *Console) Print(s string)    { c.lines = append(c.lines, s) }
func (p *ErrorPanel) Show(s string)  { p.errors = append(p.errors, s) }
func (p *ErrorPanel) Clear()         { p.errors = nil }

// Session es el mediador: conoce a todos y decide quién reacciona a qué.
type Session struct {
	editor  *Editor
	console *Console
	panel   *ErrorPanel
}

func NewSession() *Session {
	s := &Session{console: &Console{}, panel: &ErrorPanel{}}
	s.editor = &Editor{m: s}
	return s
}

func (s *Session) Notify(sender, event, data string) {
	if sender == "editor" && event == "submit" {
		s.panel.Clear()
		if data == "" {
			s.panel.Show("código vacío")
			return
		}
		s.console.Print(">> " + data)
	}
}

func main() {
	s := NewSession()
	s.editor.Submit("let x = 5;")
	s.editor.Submit("")
	fmt.Println("consola:", s.console.lines)
	fmt.Println("errores:", s.panel.errors)
}

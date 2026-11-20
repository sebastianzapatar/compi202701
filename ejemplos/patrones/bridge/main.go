// Bridge: separar una abstracción de su implementación para que varíen por separado.
// Abstracción: QUÉ se reporta de un error.  Implementación: DÓNDE se muestra.
package main

import (
	"fmt"
	"html"
)

// Renderer es la implementación (el "puente").
type Renderer interface {
	Title(s string) string
	Body(s string) string
}

type terminal struct{}
type htmlOut struct{}

func (terminal) Title(s string) string { return "== " + s + " ==" }
func (terminal) Body(s string) string  { return "   " + s }
func (htmlOut) Title(s string) string  { return "<h3>" + html.EscapeString(s) + "</h3>" }
func (htmlOut) Body(s string) string   { return "<p>" + html.EscapeString(s) + "</p>" }

// ErrorReport es la abstracción: guarda una referencia a la implementación.
type ErrorReport struct {
	r   Renderer
	msg string
}

func (e ErrorReport) Short() string { return e.r.Title("Error") + "\n" + e.r.Body(e.msg) }

// DetailedReport refina la abstracción sin tocar los renderers.
type DetailedReport struct {
	ErrorReport
	line, col int
}

func (d DetailedReport) Full() string {
	return d.Short() + "\n" + d.r.Body(fmt.Sprintf("línea %d, columna %d", d.line, d.col))
}

func main() {
	for _, r := range []Renderer{terminal{}, htmlOut{}} {
		rep := DetailedReport{ErrorReport: ErrorReport{r: r, msg: "expected =, got <INT>"}, line: 3, col: 7}
		fmt.Println(rep.Full())
	}
}

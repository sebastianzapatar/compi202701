// Abstract Factory: crear FAMILIAS de objetos relacionados que deben usarse juntos.
// Ejemplo: el tema del REPL (prompt + formato de errores) en versión "color" o "plana".
package main

import "fmt"

type Prompt interface{ Render() string }
type ErrorFormatter interface{ Format(msg string) string }

// Theme es la fábrica abstracta.
type Theme interface {
	NewPrompt() Prompt
	NewErrorFormatter() ErrorFormatter
}

// Familia "color" (secuencias ANSI de terminal)
type colorTheme struct{}
type colorPrompt struct{}
type colorErrors struct{}

func (colorTheme) NewPrompt() Prompt                 { return colorPrompt{} }
func (colorTheme) NewErrorFormatter() ErrorFormatter { return colorErrors{} }
func (colorPrompt) Render() string                   { return "\033[36m>>\033[0m " }
func (colorErrors) Format(m string) string           { return "\033[31merror:\033[0m " + m }

// Familia "plana" (archivos de log, terminales sin color)
type plainTheme struct{}
type plainPrompt struct{}
type plainErrors struct{}

func (plainTheme) NewPrompt() Prompt                 { return plainPrompt{} }
func (plainTheme) NewErrorFormatter() ErrorFormatter { return plainErrors{} }
func (plainPrompt) Render() string                   { return ">> " }
func (plainErrors) Format(m string) string           { return "[ERROR] " + m }

// El cliente solo conoce interfaces: nunca mezcla piezas de familias distintas.
func session(t Theme) {
	p, e := t.NewPrompt(), t.NewErrorFormatter()
	fmt.Printf("%q\n", p.Render()+"let x = ;")
	fmt.Printf("%q\n", e.Format("no prefix parse function for ; found"))
}

func main() {
	session(plainTheme{})
	session(colorTheme{})
}

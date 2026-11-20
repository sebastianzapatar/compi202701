// Template Method: el esqueleto de un algoritmo queda fijo y los pasos variables se delegan.
// En Go no hay clases abstractas: el esqueleto es una función y los pasos, una interfaz.
package main

import (
	"fmt"
	"strings"
)

// Language define los pasos que cada lenguaje debe proveer.
type Language interface {
	Tokenize(src string) []string
	Validate(tokens []string) error
	Execute(tokens []string) string
}

// Run es el método plantilla: el orden de los pasos nunca cambia.
func Run(lang Language, src string) string {
	toks := lang.Tokenize(src)
	if err := lang.Validate(toks); err != nil {
		return "error: " + err.Error()
	}
	return lang.Execute(toks)
}

type Calc struct{}

func (Calc) Tokenize(s string) []string { return strings.Fields(s) }
func (Calc) Validate(t []string) error {
	if len(t)%2 == 0 {
		return fmt.Errorf("expresión incompleta: %v", t)
	}
	return nil
}
func (Calc) Execute(t []string) string { return "calculando " + strings.Join(t, "") }

type Shout struct{}

func (Shout) Tokenize(s string) []string { return strings.Split(s, "") }
func (Shout) Validate([]string) error    { return nil }
func (Shout) Execute(t []string) string  { return strings.ToUpper(strings.Join(t, "")) }

func main() {
	fmt.Println(Run(Calc{}, "1 + 2"))
	fmt.Println(Run(Calc{}, "1 +"))
	fmt.Println(Run(Shout{}, "monkey"))
}

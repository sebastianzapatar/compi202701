// Factory Method: una función decide qué tipo concreto crear y devuelve una interfaz.
// Ejemplo: el REPL elige cómo imprimir resultados según una opción de línea de comandos.
package main

import (
	"fmt"
	"strings"
)

type Printer interface{ Print(result string) string }

type plainPrinter struct{}
type jsonPrinter struct{}
type boxPrinter struct{ width int }

func (plainPrinter) Print(r string) string { return r }
func (jsonPrinter) Print(r string) string  { return fmt.Sprintf(`{"result": %q}`, r) }
func (b boxPrinter) Print(r string) string {
	line := strings.Repeat("-", b.width)
	return fmt.Sprintf("+%s+\n|%-*s|\n+%s+", line, b.width, r, line)
}

// NewPrinter es el método fábrica: el cliente pide por nombre y recibe la abstracción.
func NewPrinter(format string) (Printer, error) {
	switch format {
	case "plain":
		return plainPrinter{}, nil
	case "json":
		return jsonPrinter{}, nil
	case "box":
		return boxPrinter{width: 8}, nil
	}
	return nil, fmt.Errorf("formato desconocido: %s", format)
}

func main() {
	for _, f := range []string{"plain", "json", "box", "xml"} {
		p, err := NewPrinter(f)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		fmt.Println(p.Print("42"))
	}
}

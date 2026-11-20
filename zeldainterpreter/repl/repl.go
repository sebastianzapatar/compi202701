// Package repl implementa la consola interactiva (Read-Eval-Print Loop):
// LEE una línea, la EVALÚA, IMPRIME el resultado y vuelve a empezar.
//
// Además del código Monkey, acepta comandos que empiezan con ":".
package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"zeldainterpreter/object"
	"zeldainterpreter/runner"
)

const (
	PROMPT       = ">> " // esperando una entrada nueva
	CONTINUATION = ".. " // la entrada anterior tiene llaves/paréntesis sin cerrar
)

// session guarda el estado de la consola mientras está abierta.
type session struct {
	out  io.Writer
	env  *object.Environment // persiste entre líneas: así recordamos las variables
	mode runner.Mode
}

// Start abre la consola. Lee de `in` y escribe en `out` (normalmente la
// terminal, pero las pruebas pasan un string y un buffer).
func Start(in io.Reader, out io.Writer, mode runner.Mode) {
	s := &session{out: out, env: object.NewEnvironment(), mode: mode}
	s.printBanner()

	scanner := bufio.NewScanner(in)
	var pending strings.Builder // acumula líneas cuando un bloque no está cerrado

	for {
		if pending.Len() == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONTINUATION)
		}

		if !scanner.Scan() { // Ctrl+D (o fin del input): salimos
			fmt.Fprintln(out)
			return
		}
		line := scanner.Text()

		// Los comandos (:ayuda, :modo...) solo se aceptan al inicio de una entrada.
		if pending.Len() == 0 && strings.HasPrefix(strings.TrimSpace(line), ":") {
			if quit := s.runCommand(strings.TrimSpace(line)); quit {
				return
			}
			continue
		}

		pending.WriteString(line + "\n")

		// ¿Faltan cerrar llaves, paréntesis o corchetes? Seguimos leyendo.
		// Así se pueden escribir funciones de varias líneas en la consola.
		if openDelimiters(pending.String()) > 0 {
			continue
		}

		source := pending.String()
		pending.Reset()
		if strings.TrimSpace(source) == "" {
			continue
		}
		runner.Run(source, runner.Options{Mode: s.mode, ShowResult: true}, s.env, out)
	}
}

func (s *session) printBanner() {
	fmt.Fprintln(s.out, "ZeldaInterpreter · intérprete del lenguaje Monkey")
	fmt.Fprintf(s.out, "Modo actual: %s. Escribe :ayuda para ver los comandos, :salir para terminar.\n", s.mode)
}

// runCommand ejecuta un comando de la consola. Devuelve true si hay que salir.
func (s *session) runCommand(line string) bool {
	fields := strings.Fields(line)
	cmd, args := fields[0], fields[1:]

	switch cmd {
	case ":salir", ":q":
		fmt.Fprintln(s.out, "¡Hasta pronto!")
		return true

	case ":ayuda", ":h":
		fmt.Fprint(s.out, `Comandos:
  :modo <m>        cambia el modo: eval | tokens | ast | arbol
  :env             muestra las variables definidas
  :cargar <ruta>   ejecuta un archivo .mk en esta sesión
  :reiniciar       borra todas las variables
  :salir           termina la consola
Puedes escribir funciones en varias líneas: la consola espera a que cierres { ( [.
`)

	case ":modo":
		if len(args) == 0 {
			fmt.Fprintf(s.out, "Modo actual: %s\n", s.mode)
			break
		}
		mode, err := runner.ParseMode(args[0])
		if err != nil {
			fmt.Fprintln(s.out, err)
			break
		}
		s.mode = mode
		fmt.Fprintf(s.out, "Modo cambiado a: %s\n", s.mode)

	case ":env":
		names := s.env.Names()
		if len(names) == 0 {
			fmt.Fprintln(s.out, "(no hay variables definidas)")
		}
		for _, name := range names {
			val, _ := s.env.Get(name)
			fmt.Fprintf(s.out, "  %-12s %-9s %s\n", name, val.Type(), summarize(val))
		}

	case ":cargar":
		if len(args) == 0 {
			fmt.Fprintln(s.out, "Uso: :cargar <ruta-del-archivo.mk>")
			break
		}
		code, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Fprintln(s.out, "No se pudo leer el archivo:", err)
			break
		}
		runner.Run(string(code), runner.Options{Mode: s.mode, ShowResult: true}, s.env, s.out)

	case ":reiniciar":
		s.env = object.NewEnvironment()
		fmt.Fprintln(s.out, "Entorno reiniciado.")

	default:
		fmt.Fprintf(s.out, "Comando desconocido %q. Escribe :ayuda.\n", cmd)
	}
	return false
}

// summarize muestra un valor en una sola línea corta (para :env).
func summarize(obj object.Object) string {
	if fn, ok := obj.(*object.Function); ok {
		params := make([]string, len(fn.Parameters))
		for i, p := range fn.Parameters {
			params[i] = p.Value
		}
		return "fn(" + strings.Join(params, ", ") + ")"
	}
	text := strings.ReplaceAll(obj.Inspect(), "\n", " ")
	if len(text) > 50 {
		text = text[:47] + "..."
	}
	return text
}

// openDelimiters cuenta cuántos { ( [ quedan sin cerrar, ignorando lo que
// esté dentro de strings. Un resultado > 0 significa "la entrada sigue".
func openDelimiters(src string) int {
	depth, inString := 0, false
	for i := 0; i < len(src); i++ {
		switch ch := src[i]; {
		case ch == '"':
			inString = !inString
		case inString:
			// dentro de un string los símbolos no cuentan
		case ch == '{' || ch == '(' || ch == '[':
			depth++
		case ch == '}' || ch == ')' || ch == ']':
			depth--
		}
	}
	return depth
}

// ZeldaInterpreter: intérprete del lenguaje Monkey escrito en Go.
//
// Modos de uso:
//
//	zelda                          consola interactiva (REPL)
//	zelda programa.mk              ejecuta un archivo
//	zelda -                        ejecuta lo que llegue por la entrada estándar
//	zelda -e 'puts(1 + 2)'         ejecuta código escrito en la línea de comandos
//	zelda -modo=tokens programa.mk muestra los tokens (también: ast, arbol, eval)
//
// Importante: las opciones (-modo, -e...) van ANTES del nombre del archivo.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"zeldainterpreter/object"
	"zeldainterpreter/repl"
	"zeldainterpreter/runner"
)

const version = "1.0.0"

func main() {
	// main solo traduce el resultado de run a un código de salida del proceso.
	// Separarlo así permite probar run() sin abrir una terminal real.
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// run interpreta los argumentos y devuelve el código de salida:
// 0 = todo bien, 1 = el programa tuvo errores, 2 = uso incorrecto.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("zelda", flag.ContinueOnError)
	flags.SetOutput(stderr)

	modeFlag := flags.String("modo", "eval", "qué mostrar: eval | tokens | ast | arbol")
	code := flags.String("e", "", "código Monkey a ejecutar directamente")
	showResult := flags.Bool("resultado", false, "al ejecutar un archivo, mostrar también el valor final")
	showVersion := flags.Bool("version", false, "mostrar la versión y salir")

	flags.Usage = func() {
		fmt.Fprint(stderr, `ZeldaInterpreter — intérprete del lenguaje Monkey

Uso:
  zelda [opciones]                 consola interactiva
  zelda [opciones] archivo.mk      ejecutar un archivo
  zelda [opciones] -               ejecutar desde la entrada estándar
  zelda [opciones] -e "código"     ejecutar código de la línea de comandos

Opciones (van ANTES del archivo):
`)
		flags.PrintDefaults()
		fmt.Fprint(stderr, `
Ejemplos:
  zelda programas/fibonacci.mk
  zelda -modo=arbol -e "1 + 2 * 3"
  zelda -modo=tokens programas/hola.mk
  cat programas/closures.mk | zelda -
`)
	}

	if err := flags.Parse(args); err != nil {
		return 2 // -h / -help también terminan aquí, tras mostrar la ayuda
	}
	if *showVersion {
		fmt.Fprintln(stdout, "ZeldaInterpreter", version)
		return 0
	}

	mode, err := runner.ParseMode(*modeFlag)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	// Elegimos el modo de entrada según los argumentos.
	switch {
	case *code != "": // 1) código en la línea de comandos
		return execute(*code, runner.Options{Mode: mode, ShowResult: true}, stdout)

	case flags.NArg() == 0: // 2) sin archivo: consola interactiva
		repl.Start(stdin, stdout, mode)
		return 0

	case flags.NArg() > 1:
		fmt.Fprintln(stderr, "solo se puede ejecutar un archivo a la vez")
		return 2

	default: // 3) archivo o "-" (entrada estándar)
		source, err := readSource(flags.Arg(0), stdin)
		if err != nil {
			fmt.Fprintln(stderr, "error:", err)
			return 1
		}
		return execute(source, runner.Options{Mode: mode, ShowResult: *showResult}, stdout)
	}
}

// readSource lee el programa desde un archivo, o desde stdin si la ruta es "-".
func readSource(path string, stdin io.Reader) (string, error) {
	var data []byte
	var err error
	if path == "-" {
		data, err = io.ReadAll(stdin)
	} else {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// execute corre un programa completo en un entorno nuevo.
func execute(source string, opts runner.Options, out io.Writer) int {
	if runner.Run(source, opts, object.NewEnvironment(), out) {
		return 0
	}
	return 1
}

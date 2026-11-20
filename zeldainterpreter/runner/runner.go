// Package runner lleva un código fuente por el pipeline del intérprete
// (lexer → parser → evaluador) y muestra el resultado según el MODO elegido.
//
// Lo usan tanto la consola interactiva (REPL) como la ejecución de archivos,
// así ambos se comportan exactamente igual.
package runner

import (
	"fmt"
	"io"
	"strings"

	"zeldainterpreter/ast"
	"zeldainterpreter/evaluator"
	"zeldainterpreter/lexer"
	"zeldainterpreter/object"
	"zeldainterpreter/parser"
)

// Mode indica hasta qué fase del pipeline se ejecuta y qué se muestra.
type Mode string

const (
	ModeTokens Mode = "tokens" // solo el lexer: lista de tokens
	ModeAST    Mode = "ast"    // lexer + parser: el programa con paréntesis explícitos
	ModeTree   Mode = "arbol"  // lexer + parser: el AST dibujado como árbol
	ModeEval   Mode = "eval"   // pipeline completo: ejecutar el programa
)

// Modes lista los modos válidos (en el orden en que se muestran en la ayuda).
var Modes = []Mode{ModeEval, ModeTokens, ModeAST, ModeTree}

// ParseMode convierte un texto (por ejemplo de la línea de comandos) en un Mode.
func ParseMode(s string) (Mode, error) {
	for _, m := range Modes {
		if string(m) == strings.ToLower(strings.TrimSpace(s)) {
			return m, nil
		}
	}
	return "", fmt.Errorf("modo desconocido %q (usa: eval, tokens, ast o arbol)", s)
}

// Options configura una ejecución.
type Options struct {
	Mode       Mode
	ShowResult bool // mostrar el valor final (el REPL lo hace siempre)
}

// Run ejecuta `source` en el entorno `env` y escribe todo en `out`.
// Devuelve false si hubo errores de sintaxis, de ejecución o un panic.
func Run(source string, opts Options, env *object.Environment, out io.Writer) (ok bool) {
	// Red de seguridad: si el evaluador provoca un panic de Go (por ejemplo,
	// al dividir por cero), lo atrapamos para que la consola no se cierre.
	// defer + recover es el mecanismo de Go para recuperarse de un panic.
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(out, "ERROR INTERNO (panic de Go): %v\n", r)
			fmt.Fprintln(out, "  El intérprete del libro no valida este caso (¿división por cero o faltan argumentos?).")
			ok = false
		}
	}()

	// Fase 1: análisis léxico.
	if opts.Mode == ModeTokens {
		printTokens(source, out)
		return true
	}

	// Fase 2: análisis sintáctico.
	p := parser.New(lexer.New(source))
	program := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		PrintParserErrors(out, errs)
		return false
	}

	switch opts.Mode {
	case ModeAST:
		for _, stmt := range program.Statements {
			fmt.Fprintln(out, stmt.String())
		}
		return true
	case ModeTree:
		fmt.Fprint(out, ast.Tree(program))
		return true
	}

	// Fase 3: evaluación.
	evaluator.Output = out // la salida de puts va al mismo lugar
	result := evaluator.Eval(program, env)

	if errObj, isErr := result.(*object.Error); isErr {
		fmt.Fprintln(out, errObj.Inspect())
		return false
	}
	if opts.ShowResult && result != nil {
		fmt.Fprintln(out, result.Inspect())
	}
	return true
}

// printTokens muestra una tabla con cada token producido por el lexer.
func printTokens(source string, out io.Writer) {
	for i, tok := range lexer.Tokenize(source) {
		fmt.Fprintf(out, "%3d  %-10s %q\n", i+1, tok.Type, tok.Literal)
	}
}

// mascot es la mascota de ZeldaInterpreter: un pequeño mago que aparece
// cuando hay errores de sintaxis.
const mascot = `
       /\
      /  \       ZeldaInterpreter
     /_*__\      ¡Ups! El hechizo tiene errores de sintaxis:
     ( o.o )
      > ^ <
`

// PrintParserErrors muestra los errores de sintaxis con la mascota.
func PrintParserErrors(out io.Writer, errs []string) {
	fmt.Fprint(out, mascot)
	for _, msg := range errs {
		fmt.Fprintln(out, "  •", msg)
	}
}

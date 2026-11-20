package main

import (
	"bytes"
	"strings"
	"testing"
)

// runCLI ejecuta el programa como si se llamara desde la terminal.
func runCLI(t *testing.T, stdin string, args ...string) (string, int) {
	t.Helper()
	var out, errOut bytes.Buffer
	code := run(args, strings.NewReader(stdin), &out, &errOut)
	return out.String() + errOut.String(), code
}

func TestModos(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		stdin    string
		want     string
		wantCode int
	}{
		{"eval con -e", []string{"-e", "1 + 2 * 3"}, "", "7\n", 0},
		{"tokens", []string{"-modo=tokens", "-e", "x"}, "", `IDENT      "x"`, 0},
		{"ast", []string{"-modo=ast", "-e", "1 + 2 * 3"}, "", "(1 + (2 * 3))", 0},
		{"arbol", []string{"-modo=arbol", "-e", "-x"}, "", "└── operando: Identifier x", 0},
		{"archivo", []string{"programas/fibonacci.mk"}, "", "fibonacci(15) =\n610", 0},
		{"stdin", []string{"-"}, `puts("desde stdin")`, "desde stdin", 0},
		{"error de sintaxis", []string{"-e", "let = 1"}, "", "expected next token to be IDENT", 1},
		{"error de ejecución", []string{"-e", "1 + true"}, "", "ERROR: type mismatch", 1},
		{"panic atrapado", []string{"-e", "1 / 0"}, "", "ERROR INTERNO", 1},
		{"modo inválido", []string{"-modo=xyz", "-e", "1"}, "", "modo desconocido", 2},
		{"archivo inexistente", []string{"no-existe.mk"}, "", "error:", 1},
		{"versión", []string{"-version"}, "", "ZeldaInterpreter", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, code := runCLI(t, tt.stdin, tt.args...)
			if code != tt.wantCode || !strings.Contains(out, tt.want) {
				t.Errorf("código=%d (esperado %d)\nsalida:\n%s\nse esperaba que contuviera %q", code, tt.wantCode, out, tt.want)
			}
		})
	}
}

// TestREPL simula una sesión completa en la consola interactiva.
func TestREPL(t *testing.T) {
	session := strings.Join([]string{
		"let a = 5;",
		"a * 2",
		"let doble = fn(x) {", // función en varias líneas
		"  x * 2",
		"};",
		"doble(21)",
		":env",
		":modo ast",
		"1 + 2 * 3",
		":modo eval",
		":reiniciar",
		"a",
		":salir",
	}, "\n")

	out, code := runCLI(t, session)
	if code != 0 {
		t.Fatalf("código de salida %d", code)
	}
	for _, want := range []string{
		"10\n",                           // a * 2
		".. ",                            // prompt de continuación
		"42\n",                           // doble(21)
		"doble        FUNCTION  fn(x)",   // :env
		"(1 + (2 * 3))",                  // modo ast
		"ERROR: identifier not found: a", // tras :reiniciar
		"¡Hasta pronto!",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("la sesión no contiene %q\n--- salida ---\n%s", want, out)
		}
	}
}

package evaluator

import (
	"bytes"
	"testing"

	"zeldainterpreter/lexer"
	"zeldainterpreter/object"
	"zeldainterpreter/parser"
)

// eval es un ayudante: lexer → parser → evaluador en un entorno nuevo.
func eval(t *testing.T, input string) object.Object {
	t.Helper()
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		t.Fatalf("errores de sintaxis en %q: %v", input, errs)
	}
	return Eval(program, object.NewEnvironment())
}

// TestEval compara el Inspect() del resultado con lo esperado.
// Está organizada por temas, en el mismo orden del curso.
func TestEval(t *testing.T) {
	tests := []struct{ input, want string }{
		// Enteros y aritmética
		{"5", "5"},
		{"-10", "-10"},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", "50"},
		{"7 / 2", "3"},
		// Booleanos, ! y comparaciones
		{"!!5", "true"},
		{"!0", "false"},
		{"1 < 2 == true", "true"},
		{"true == 1", "false"},
		// Condicionales y truthiness
		{"if (0) { 10 } else { 20 }", "10"},
		{"if (false) { 10 }", "null"},
		// return anidado
		{"9; return 2 * 5; 9;", "10"},
		{"if (10 > 1) { if (10 > 1) { return 10; } return 1; }", "10"},
		// Errores
		{"-true", "ERROR: unknown operator: -BOOLEAN"},
		{"5 + true; 5;", "ERROR: type mismatch: INTEGER + BOOLEAN"},
		{"true + false", "ERROR: unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "ERROR: unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "ERROR: identifier not found: foobar"},
		{"let x = 5; x(1);", "ERROR: not a function: INTEGER"},
		// Variables y entornos
		{"let a = 5; let b = a; let c = a + b + 5; c;", "15"},
		{"let x = 10; let f = fn(x) { x * 2 }; f(3) + x;", "16"},
		// Funciones y closures
		{"fn(x) { x + 2; }", "fn(x) {\n(x + 2)\n}"},
		{"let newAdder = fn(x) { fn(y) { x + y } }; newAdder(2)(3);", "5"},
		{"let x = 10; let f = fn() { x }; let g = fn(x) { f() }; g(99);", "10"}, // alcance léxico
		{"fn() {}()", "null"},
		// Strings
		{`"Hola" + " " + "mundo"`, "Hola mundo"},
		{`"a" - "b"`, "ERROR: unknown operator: STRING - STRING"},
		// Built-ins
		{`len("four")`, "4"},
		{"len(1)", "ERROR: argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "ERROR: wrong number of arguments. got=2, want=1"},
		{"first([])", "null"},
		{"rest(rest([1, 2, 3]))", "[3]"},
		{"push([], 1)", "[1]"},
		{"let a = [1]; let b = push(a, 2); len(a);", "1"}, // inmutabilidad
		// Arrays
		{"[1, 2 * 2, 3 + 3]", "[1, 4, 6]"},
		{"[1, 2, 3][3]", "null"},
		{"[1, 2, 3][-1]", "null"},
		// Hashes
		{`{"one": 10 - 9, "thr" + "ee": 6 / 2, 4: 4, true: 5}`, "{one: 1, three: 3, 4: 4, true: 5}"},
		{`{"foo": 5}["foo"]`, "5"},
		{`{1: "a"}[1 == 1]`, "null"}, // BOOLEAN true no es la llave INTEGER 1
		{`{"name": "Monkey"}[fn(x) { x }]`, "ERROR: unusable as hash key: FUNCTION"},
		// Gran final
		{`let map = fn(arr, f) {
			let iter = fn(arr, acc) {
				if (len(arr) == 0) { acc } else { iter(rest(arr), push(acc, f(first(arr)))) }
			};
			iter(arr, []);
		};
		map([1, 2, 3, 4], fn(x) { x * 2 });`, "[2, 4, 6, 8]"},
		{"let fib = fn(n) { if (n < 2) { n } else { fib(n - 1) + fib(n - 2) } }; fib(15);", "610"},
	}

	for _, tt := range tests {
		got := eval(t, tt.input)
		if got == nil {
			t.Errorf("%q: se obtuvo nil, se esperaba %q", tt.input, tt.want)
			continue
		}
		if got.Inspect() != tt.want {
			t.Errorf("%q\n  se esperaba %q\n  se obtuvo   %q", tt.input, tt.want, got.Inspect())
		}
	}
}

func TestLetDoesNotProduceValue(t *testing.T) {
	if got := eval(t, "let a = 5;"); got != nil {
		t.Errorf("un let no debería producir valor, se obtuvo %s", got.Inspect())
	}
}

func TestSingletons(t *testing.T) {
	if eval(t, "1 < 2") != TRUE || eval(t, "!true") != FALSE || eval(t, "if (false) { 1 }") != NULL {
		t.Error("TRUE, FALSE y NULL deben ser instancias únicas")
	}
}

func TestPuts(t *testing.T) {
	var buf bytes.Buffer
	previous := Output
	Output = &buf
	defer func() { Output = previous }() // restauramos para no afectar otras pruebas

	eval(t, `puts("hola", 1 + 2)`)
	if buf.String() != "hola\n3\n" {
		t.Errorf("puts escribió %q", buf.String())
	}
}

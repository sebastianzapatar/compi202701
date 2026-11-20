package parser

import (
	"testing"

	"zeldainterpreter/ast"
	"zeldainterpreter/lexer"
)

// parse es un ayudante: parsea el input y falla la prueba si hubo errores.
func parse(t *testing.T, input string) *ast.Program {
	t.Helper()
	p := New(lexer.New(input))
	program := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		t.Fatalf("errores inesperados para %q: %v", input, errs)
	}
	return program
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input     string
		wantName  string
		wantValue string
	}{
		{"let x = 5;", "x", "5"},
		{"let y = true;", "y", "true"},
		{"let total = a + b * 2;", "total", "(a + (b * 2))"},
	}
	for _, tt := range tests {
		program := parse(t, tt.input)
		if len(program.Statements) != 1 {
			t.Fatalf("%q: se esperaba 1 sentencia, hay %d", tt.input, len(program.Statements))
		}
		let, ok := program.Statements[0].(*ast.LetStatement) // type assertion
		if !ok {
			t.Fatalf("%q: no es *ast.LetStatement sino %T", tt.input, program.Statements[0])
		}
		if let.Name.Value != tt.wantName || let.Value.String() != tt.wantValue {
			t.Errorf("%q: se obtuvo nombre=%s valor=%s", tt.input, let.Name.Value, let.Value.String())
		}
	}
}

func TestReturnStatements(t *testing.T) {
	program := parse(t, "return 5; return x + 1;")
	if len(program.Statements) != 2 {
		t.Fatalf("se esperaban 2 sentencias, hay %d", len(program.Statements))
	}
	for _, stmt := range program.Statements {
		if _, ok := stmt.(*ast.ReturnStatement); !ok {
			t.Errorf("se esperaba *ast.ReturnStatement, se obtuvo %T", stmt)
		}
	}
}

// TestOperatorPrecedence es la prueba más importante del parser: compara el
// String() (que agrega paréntesis explícitos) con la agrupación esperada.
func TestOperatorPrecedence(t *testing.T) {
	tests := []struct{ input, want string }{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
	}
	for _, tt := range tests {
		if got := parse(t, tt.input).String(); got != tt.want {
			t.Errorf("%q\n  se esperaba %s\n  se obtuvo   %s", tt.input, tt.want, got)
		}
	}
}

func TestCompoundExpressions(t *testing.T) {
	tests := []struct{ input, want string }{
		{"if (x < y) { x } else { y }", "if(x < y) xelse y"},
		{"fn(x, y) { x + y; }", "fn(x, y) (x + y)"},
		{"fn() { 1 }()", "fn() 1()"},
		{`{"uno": 1, "dos": 1 + 1}`, "{uno:1, dos:(1 + 1)}"},
		{"[]", "[]"},
	}
	for _, tt := range tests {
		if got := parse(t, tt.input).String(); got != tt.want {
			t.Errorf("%q: se esperaba %q, se obtuvo %q", tt.input, tt.want, got)
		}
	}
}

// TestParserErrors verifica el PRIMER mensaje de error para entradas inválidas.
func TestParserErrors(t *testing.T) {
	tests := []struct{ input, want string }{
		{"let x 5;", "expected next token to be =, got INT instead"},
		{"let = 10;", "expected next token to be IDENT, got = instead"},
		{"5 + ;", "no prefix parse function for ; found"},
		{"(1 + 2", "expected next token to be ), got EOF instead"},
		{"fn(x y) {}", "expected next token to be ), got IDENT instead"},
		{`{"a" 1}`, "expected next token to be :, got INT instead"},
		{"99999999999999999999", `could not parse "99999999999999999999" as integer`},
	}
	for _, tt := range tests {
		p := New(lexer.New(tt.input))
		p.ParseProgram()
		errs := p.Errors()
		if len(errs) == 0 || errs[0] != tt.want {
			t.Errorf("%q: se esperaba el error %q, se obtuvo %v", tt.input, tt.want, errs)
		}
	}
}

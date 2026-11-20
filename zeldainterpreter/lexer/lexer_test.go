package lexer

import (
	"testing"

	"zeldainterpreter/token"
)

// TestNextToken recorre un programa que usa TODOS los tipos de token y
// verifica, uno por uno, el tipo y el literal producidos.
func TestNextToken(t *testing.T) {
	input := `let cinco = 5;
let suma = fn(x, y) { x + y; };
!-/*5;
5 < 10 > 5;
if (5 < 10) { return true; } else { return false; }
10 == 10; 10 != 9;
"hola mundo"
[1, 2];
{"clave": "valor"}
@`

	tests := []struct {
		wantType    token.TokenType
		wantLiteral string
	}{
		{token.LET, "let"}, {token.IDENT, "cinco"}, {token.ASSIGN, "="}, {token.INT, "5"}, {token.SEMICOLON, ";"},
		{token.LET, "let"}, {token.IDENT, "suma"}, {token.ASSIGN, "="}, {token.FUNCTION, "fn"},
		{token.LPAREN, "("}, {token.IDENT, "x"}, {token.COMMA, ","}, {token.IDENT, "y"}, {token.RPAREN, ")"},
		{token.LBRACE, "{"}, {token.IDENT, "x"}, {token.PLUS, "+"}, {token.IDENT, "y"}, {token.SEMICOLON, ";"},
		{token.RBRACE, "}"}, {token.SEMICOLON, ";"},
		{token.BANG, "!"}, {token.MINUS, "-"}, {token.SLASH, "/"}, {token.ASTERISK, "*"}, {token.INT, "5"}, {token.SEMICOLON, ";"},
		{token.INT, "5"}, {token.LT, "<"}, {token.INT, "10"}, {token.GT, ">"}, {token.INT, "5"}, {token.SEMICOLON, ";"},
		{token.IF, "if"}, {token.LPAREN, "("}, {token.INT, "5"}, {token.LT, "<"}, {token.INT, "10"}, {token.RPAREN, ")"},
		{token.LBRACE, "{"}, {token.RETURN, "return"}, {token.TRUE, "true"}, {token.SEMICOLON, ";"}, {token.RBRACE, "}"},
		{token.ELSE, "else"}, {token.LBRACE, "{"}, {token.RETURN, "return"}, {token.FALSE, "false"}, {token.SEMICOLON, ";"}, {token.RBRACE, "}"},
		{token.INT, "10"}, {token.EQ, "=="}, {token.INT, "10"}, {token.SEMICOLON, ";"},
		{token.INT, "10"}, {token.NOT_EQ, "!="}, {token.INT, "9"}, {token.SEMICOLON, ";"},
		{token.STRING, "hola mundo"},
		{token.LBRACKET, "["}, {token.INT, "1"}, {token.COMMA, ","}, {token.INT, "2"}, {token.RBRACKET, "]"}, {token.SEMICOLON, ";"},
		{token.LBRACE, "{"}, {token.STRING, "clave"}, {token.COLON, ":"}, {token.STRING, "valor"}, {token.RBRACE, "}"},
		{token.ILLEGAL, "@"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.wantType || tok.Literal != tt.wantLiteral {
			t.Fatalf("token %d: se esperaba %s %q, se obtuvo %s %q", i, tt.wantType, tt.wantLiteral, tok.Type, tok.Literal)
		}
	}
}

// TestBordes documenta casos que suelen sorprender a los estudiantes.
func TestBordes(t *testing.T) {
	tests := []struct {
		input string
		want  []token.TokenType
	}{
		{"x5", []token.TokenType{token.IDENT, token.INT, token.EOF}},                // isLetter no acepta dígitos
		{"5x", []token.TokenType{token.INT, token.IDENT, token.EOF}},                // número y luego nombre
		{"!==", []token.TokenType{token.NOT_EQ, token.ASSIGN, token.EOF}},           // el lexer es "codicioso"
		{"=!=", []token.TokenType{token.ASSIGN, token.NOT_EQ, token.EOF}},           // de izquierda a derecha
		{"3.14", []token.TokenType{token.INT, token.ILLEGAL, token.INT, token.EOF}}, // no hay decimales
		{"", []token.TokenType{token.EOF}},
	}
	for _, tt := range tests {
		got := Tokenize(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("%q: se esperaban %d tokens, se obtuvieron %d (%v)", tt.input, len(tt.want), len(got), got)
			continue
		}
		for i := range got {
			if got[i].Type != tt.want[i] {
				t.Errorf("%q: token %d es %s, se esperaba %s", tt.input, i, got[i].Type, tt.want[i])
			}
		}
	}
}

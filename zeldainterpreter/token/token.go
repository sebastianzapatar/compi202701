// Package token define el VOCABULARIO de ZeldaInterpreter: todos los tipos de
// token que el lexer puede producir y que el parser sabe consumir.
//
// Un token es como una "palabra" del programa. Por ejemplo, en
//
//	let x = 5;
//
// hay cinco tokens: LET, IDENT("x"), ASSIGN("="), INT("5") y SEMICOLON(";").
package token

// TokenType es la CATEGORÍA de un token.
//
// Usamos string (y no int) a propósito: así un tipo se puede imprimir
// directamente al depurar. fmt.Println(tok.Type) muestra "IDENT" en lugar
// de un número que habría que buscar en una tabla.
type TokenType string

// Token es la unidad mínima con significado que produce el lexer.
type Token struct {
	Type    TokenType // categoría: IDENT, INT, PLUS, LET...
	Literal string    // texto exacto del código fuente: "x", "5", "+", "let"...
}

// Tipos de token.
//
// Ojo: el VALOR de cada constante es lo que aparece en los mensajes de error.
// Por eso los operadores valen su propio símbolo y el parser puede decir:
// "expected next token to be =, got INT instead".
const (
	// --- Especiales ---
	ILLEGAL = "ILLEGAL" // un carácter que el lenguaje no reconoce (por ejemplo @ o $)
	EOF     = "EOF"     // "End Of File": ya no queda nada por leer

	// --- Identificadores y literales ---
	IDENT  = "IDENT"  // nombres creados por quien programa: x, suma, mi_lista
	INT    = "INT"    // números enteros: 0, 7, 1024
	STRING = "STRING" // texto entre comillas dobles: "hola"

	// --- Operadores ---
	ASSIGN   = "=" // asignación en let
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!" // negación lógica
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	EQ       = "==" // operador de DOS caracteres
	NOT_EQ   = "!=" // operador de DOS caracteres

	// --- Delimitadores ---
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":" // separa llave y valor en un hash: {"a": 1}
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"

	// --- Palabras reservadas (keywords) ---
	FUNCTION = "FUNCTION" // fn
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

// keywords asocia cada palabra reservada con su tipo de token.
// Cualquier palabra que NO esté en este mapa es un identificador.
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent decide si una palabra leída por el lexer es reservada o es un
// identificador del usuario.
//
//	LookupIdent("let")   -> LET
//	LookupIdent("letra") -> IDENT
func LookupIdent(word string) TokenType {
	if tokType, isKeyword := keywords[word]; isKeyword {
		return tokType
	}
	return IDENT
}

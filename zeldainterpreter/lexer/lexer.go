// Package lexer implementa el ANÁLISIS LÉXICO: convierte el código fuente
// (un string) en una secuencia de tokens.
//
// El lexer NO entiende la gramática. Para él, "let let let ;;" es perfectamente
// válido: produce LET, LET, LET, SEMICOLON, SEMICOLON. Decidir si eso tiene
// sentido es trabajo del parser.
package lexer

import "zeldainterpreter/token"

// Lexer recorre el código fuente carácter por carácter.
//
// Usa DOS índices sobre el input:
//
//	input:  l  e  t     x
//	        ^  ^
//	        |  readPosition: el SIGUIENTE carácter a leer
//	        position: el carácter ACTUAL (guardado en ch)
//
// Tener ambos nos permite "mirar adelante" (peekChar) sin avanzar,
// algo necesario para reconocer operadores de dos caracteres como "==".
type Lexer struct {
	input        string // el código fuente completo
	position     int    // índice del carácter actual (el que está en ch)
	readPosition int    // índice del próximo carácter a leer (position + 1)
	ch           byte   // carácter actual; el valor 0 significa "se acabó el input"
}

// New crea un lexer listo para usar sobre el código fuente dado.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // carga el primer carácter en l.ch
	return l
}

// Tokenize es un atajo: recorre TODO el input y devuelve la lista de tokens,
// incluido el EOF final. Lo usan el modo "tokens" de la consola y las pruebas.
func Tokenize(input string) []token.Token {
	l := New(input)
	var tokens []token.Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			return tokens
		}
	}
}

// NextToken es el corazón del lexer: cada llamada devuelve el siguiente token.
// Cuando ya no queda nada, devuelve siempre un token EOF.
func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace() // los espacios no producen tokens: solo separan

	var tok token.Token

	switch l.ch {
	// Operadores que pueden tener uno o dos caracteres.
	case '=':
		tok = l.oneOrTwoChars('=', token.EQ, token.ASSIGN) // "==" o "="
	case '!':
		tok = l.oneOrTwoChars('=', token.NOT_EQ, token.BANG) // "!=" o "!"

	// Operadores y delimitadores de un solo carácter.
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)

	// Texto entre comillas.
	case '"':
		tok = token.Token{Type: token.STRING, Literal: l.readString()}

	// Fin del input.
	case 0:
		return token.Token{Type: token.EOF, Literal: ""}

	default:
		// Palabras: pueden ser palabras reservadas (let, fn...) o identificadores.
		if isLetter(l.ch) {
			word := l.readIdentifier()
			// ¡Return temprano! readIdentifier ya dejó l.ch en el primer
			// carácter que NO es parte de la palabra. Si siguiéramos hasta el
			// readChar() del final, nos saltaríamos ese carácter.
			return token.Token{Type: token.LookupIdent(word), Literal: word}
		}
		// Números enteros. Mismo motivo para el return temprano.
		if isDigit(l.ch) {
			return token.Token{Type: token.INT, Literal: l.readNumber()}
		}
		// Cualquier otra cosa es un carácter ilegal.
		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar() // avanzamos para que la próxima llamada empiece en el carácter siguiente
	return tok
}

// readChar avanza una posición en el input.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // 0 (el carácter NUL de ASCII) marca el final del input
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// peekChar devuelve el siguiente carácter SIN avanzar ("mirar sin tocar").
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// oneOrTwoChars resuelve operadores como "=" / "==" o "!" / "!=".
// Si el próximo carácter es `second`, consume los dos y devuelve `twoCharType`;
// si no, devuelve el token de un solo carácter `oneCharType`.
func (l *Lexer) oneOrTwoChars(second byte, twoCharType, oneCharType token.TokenType) token.Token {
	if l.peekChar() == second {
		first := l.ch
		l.readChar() // consumimos el segundo carácter
		return token.Token{Type: twoCharType, Literal: string(first) + string(l.ch)}
	}
	return newToken(oneCharType, l.ch)
}

// skipWhitespace salta espacios, tabulaciones y saltos de línea.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readIdentifier lee una palabra completa (letras y guion bajo) y la devuelve.
// Usa slicing: recuerda dónde empezó y corta el input hasta donde terminó.
func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

// readNumber lee un entero completo (una secuencia de dígitos).
func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

// readString lee el texto entre comillas y lo devuelve SIN las comillas.
// Se detiene en la comilla de cierre o, si falta, al final del input.
func (l *Lexer) readString() string {
	start := l.position + 1 // +1 para saltar la comilla de apertura
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[start:l.position]
}

// isLetter define qué caracteres pueden formar identificadores.
// Incluimos '_' para permitir nombres como mi_variable.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit indica si el carácter es un dígito decimal.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// newToken construye un token de un solo carácter.
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

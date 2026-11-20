// Package parser implementa el ANÁLISIS SINTÁCTICO: toma los tokens del
// lexer y construye el AST.
//
// Usa dos técnicas:
//   - Descenso recursivo para las sentencias (let, return, bloques).
//   - Pratt parsing ("Top Down Operator Precedence") para las expresiones,
//     que es donde está lo difícil: precedencia y asociatividad.
package parser

import (
	"fmt"
	"strconv"

	"zeldainterpreter/ast"
	"zeldainterpreter/lexer"
	"zeldainterpreter/token"
)

// ============================================================================
// Precedencias
// ============================================================================

// Niveles de precedencia, de MENOR a MAYOR "fuerza de atracción".
// iota les asigna 1, 2, 3...: lo que importa no es el número sino el ORDEN.
const (
	_ int = iota
	LOWEST
	EQUALS      // ==  !=
	LESSGREATER // <   >
	SUM         // +   -
	PRODUCT     // *   /
	PREFIX      // -x  !x
	CALL        // f(x)
	INDEX       // lista[i]
)

// precedences dice qué precedencia tiene cada token cuando aparece como
// operador INFIJO (entre dos expresiones).
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,  // add(1, 2): el "(" después de una expresión es una llamada
	token.LBRACKET: INDEX, // lista[0]:  el "[" después de una expresión es un índice
}

// ============================================================================
// El struct Parser
// ============================================================================

// Las dos clases de funciones de parsing de Pratt:
//   - prefixParseFn: se usa cuando el token aparece AL INICIO de una expresión.
//   - infixParseFn:  se usa cuando el token aparece ENTRE dos expresiones;
//     recibe como argumento la expresión de la izquierda, ya parseada.
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(left ast.Expression) ast.Expression
)

// Parser mira los tokens a través de una "ventana" de dos posiciones:
// curToken (el actual) y peekToken (el siguiente). Con eso decide qué hacer
// sin tener que retroceder nunca.
type Parser struct {
	l      *lexer.Lexer
	errors []string // errores de sintaxis acumulados (no nos detenemos en el primero)

	curToken  token.Token
	peekToken token.Token

	// Tablas de Pratt: qué función usar para cada tipo de token.
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// New crea un parser y registra todas las funciones de parsing.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         []string{},
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}

	// Tokens que pueden INICIAR una expresión.
	p.registerPrefix(token.IDENT, p.parseIdentifier)         // x
	p.registerPrefix(token.INT, p.parseIntegerLiteral)       // 5
	p.registerPrefix(token.STRING, p.parseStringLiteral)     // "hola"
	p.registerPrefix(token.TRUE, p.parseBoolean)             // true
	p.registerPrefix(token.FALSE, p.parseBoolean)            // false
	p.registerPrefix(token.BANG, p.parsePrefixExpression)    // !x
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)   // -x
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression) // (1 + 2)
	p.registerPrefix(token.IF, p.parseIfExpression)          // if (...) {...}
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral) // fn(x) {...}
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)    // [1, 2]
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)       // {"a": 1}

	// Tokens que pueden aparecer ENTRE dos expresiones.
	for _, op := range []token.TokenType{
		token.PLUS, token.MINUS, token.ASTERISK, token.SLASH,
		token.EQ, token.NOT_EQ, token.LT, token.GT,
	} {
		p.registerInfix(op, p.parseInfixExpression) // a + b, a == b...
	}
	p.registerInfix(token.LPAREN, p.parseCallExpression)    // f(x)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression) // lista[0]

	// Leemos dos tokens para llenar la ventana curToken / peekToken.
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) registerPrefix(t token.TokenType, fn prefixParseFn) { p.prefixParseFns[t] = fn }
func (p *Parser) registerInfix(t token.TokenType, fn infixParseFn)   { p.infixParseFns[t] = fn }

// Errors devuelve los errores de sintaxis encontrados.
func (p *Parser) Errors() []string { return p.errors }

// ============================================================================
// Movimiento por los tokens
// ============================================================================

// nextToken desliza la ventana una posición hacia adelante.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool  { return p.curToken.Type == t }
func (p *Parser) peekTokenIs(t token.TokenType) bool { return p.peekToken.Type == t }

// expectPeek es una "aserción": si el siguiente token es del tipo esperado,
// avanza y devuelve true. Si no, registra un error y NO avanza.
// Así se impone el orden correcto de los tokens.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// peekPrecedence devuelve la precedencia del SIGUIENTE token (LOWEST si no es operador).
func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
}

// curPrecedence devuelve la precedencia del token ACTUAL.
func (p *Parser) curPrecedence() int {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return LOWEST
}

// ============================================================================
// Errores
// ============================================================================

func (p *Parser) peekError(expected token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", expected, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	p.errors = append(p.errors, fmt.Sprintf("no prefix parse function for %s found", t))
}

// ============================================================================
// Programa y sentencias
// ============================================================================

// ParseProgram es el punto de entrada: parsea sentencias hasta llegar a EOF.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}

	for !p.curTokenIs(token.EOF) {
		if stmt := p.parseStatement(); stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// parseStatement decide qué tipo de sentencia empieza en curToken.
//
// Nota sobre un error clásico de Go: las funciones parseXxxStatement
// devuelven la INTERFAZ ast.Statement (no *ast.LetStatement). Si devolvieran
// un puntero concreto nil, al convertirlo a interfaz quedaría una interfaz
// NO nil que contiene un puntero nil, y el `stmt != nil` de ParseProgram
// fallaría silenciosamente.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement:  let <identificador> = <expresión>;
func (p *Parser) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{Token: p.curToken} // curToken es LET

	if !p.expectPeek(token.IDENT) { // después de let debe venir un nombre
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) { // luego un =
		return nil
	}

	p.nextToken() // saltamos el = y quedamos en el inicio de la expresión
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) { // el ; final es opcional
		p.nextToken()
	}
	return stmt
}

// parseReturnStatement:  return <expresión>;
func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken() // saltamos "return"
	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// parseExpressionStatement: una expresión sola en una línea, como  x + 1;
func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// parseBlockStatement:  { <sentencias> }
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken, Statements: []ast.Statement{}}
	p.nextToken() // saltamos {

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		if stmt := p.parseStatement(); stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

// ============================================================================
// EL CORAZÓN: parseExpression (algoritmo de Pratt)
// ============================================================================

// parseExpression parsea una expresión completa.
//
// `precedence` es la fuerza con la que la llamada ACTUAL "sujeta" lo que tiene
// a su derecha. El algoritmo:
//
//  1. Busca la función PREFIX del token actual y obtiene un primer resultado
//     (leftExp). Si no hay, es un error de sintaxis.
//  2. Mientras el siguiente operador atraiga MÁS fuerte que `precedence`,
//     le entrega leftExp para que lo use como su lado izquierdo.
//  3. Cuando aparece un operador igual o más débil (o un ;), se detiene:
//     ese operador le pertenece a una llamada de más arriba.
//
// Ejemplo con 1 + 2 * 3: el * (PRODUCT) atrae más que el + (SUM), así que el
// 2 termina dentro de la multiplicación: (1 + (2 * 3)).
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()            // el operador pasa a ser curToken
		leftExp = infix(leftExp) // y "absorbe" lo que llevábamos como lado izquierdo
	}
	return leftExp
}

// ============================================================================
// Funciones PREFIX
// ============================================================================

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	value, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil { // por ejemplo, un número demasiado grande para int64
		p.errors = append(p.errors, fmt.Sprintf("could not parse %q as integer", p.curToken.Literal))
		return nil
	}
	return &ast.IntegerLiteral{Token: p.curToken, Value: value}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// parsePrefixExpression:  -x  o  !x
func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	p.nextToken()
	// PREFIX es alto: en -a * b, el - se queda solo con la a → ((-a) * b).
	expr.Right = p.parseExpression(PREFIX)
	return expr
}

// parseGroupedExpression:  ( <expresión> )
// Los paréntesis NO crean un nodo: solo hacen que el interior se parsee con
// LOWEST, como una expresión nueva. Su efecto queda en la forma del árbol.
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken() // saltamos (
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

// parseIfExpression:  if (<condición>) { ... } else { ... }
func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expr.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expr.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) { // el else es opcional
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expr.Alternative = p.parseBlockStatement()
	}
	return expr
}

// parseFunctionLiteral:  fn(<parámetros>) { <cuerpo> }
func (p *Parser) parseFunctionLiteral() ast.Expression {
	fn := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	fn.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	fn.Body = p.parseBlockStatement()
	return fn
}

// parseFunctionParameters lee  x, y, z)  y deja curToken en el ")".
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	params := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) { // fn() sin parámetros
		p.nextToken()
		return params
	}

	p.nextToken()
	params = append(params, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // saltamos el nombre anterior...
		p.nextToken() // ...y la coma
		params = append(params, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return params
}

// parseArrayLiteral:  [1, 2, 3]
func (p *Parser) parseArrayLiteral() ast.Expression {
	return &ast.ArrayLiteral{Token: p.curToken, Elements: p.parseExpressionList(token.RBRACKET)}
}

// parseHashLiteral:  {"llave": valor, ...}
func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)
		hash.Pairs = append(hash.Pairs, ast.HashPair{Key: key, Value: value})

		// Después de una pareja debe venir "}" o ","
		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	return hash
}

// ============================================================================
// Funciones INFIX (reciben la expresión de la izquierda)
// ============================================================================

// parseInfixExpression:  <left> <operador> <right>
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	// Pasar la MISMA precedencia da asociatividad a la izquierda:
	// 1 + 2 + 3  se agrupa como  ((1 + 2) + 3).
	expr.Right = p.parseExpression(precedence)
	return expr
}

// parseCallExpression:  <función>(<argumentos>)
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	return &ast.CallExpression{
		Token:     p.curToken,
		Function:  function,
		Arguments: p.parseExpressionList(token.RPAREN),
	}
}

// parseIndexExpression:  <colección>[<índice>]
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	expr.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return expr
}

// parseExpressionList lee expresiones separadas por comas hasta el token `end`.
// La usan las llamadas (end = ")") y los arrays (end = "]").
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) { // lista vacía: f() o []
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}
	return list
}

// Package ast define el ÁRBOL DE SINTAXIS ABSTRACTA (AST): la estructura que
// construye el parser y que luego recorre el evaluador.
//
// "Abstracta" significa que el árbol guarda solo lo importante: los puntos y
// coma, los paréntesis y los espacios ya cumplieron su función (guiar al
// parser) y no aparecen en el árbol.
//
// Todos los nodos tienen un método String() que imprime la expresión con
// paréntesis explícitos, por ejemplo "(1 + (2 * 3))". Eso hace muy fácil
// verificar en las pruebas que el parser respetó la precedencia.
package ast

import (
	"strings"

	"zeldainterpreter/token"
)

// ============================================================================
// Interfaces base
// ============================================================================

// Node es lo que tienen en común TODOS los nodos del árbol.
type Node interface {
	TokenLiteral() string // literal del token que originó el nodo (útil para depurar)
	String() string       // representación del nodo como código fuente
}

// Statement es una SENTENCIA: no produce un valor (let, return...).
type Statement interface {
	Node
	statementNode() // método "marcador": solo sirve para distinguir los tipos
}

// Expression es una EXPRESIÓN: produce un valor (5, x + 1, add(2, 3)...).
type Expression interface {
	Node
	expressionNode() // método "marcador"
}

// ============================================================================
// Raíz del árbol
// ============================================================================

// Program es la raíz de todo AST: un programa es una lista de sentencias.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out strings.Builder
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// ============================================================================
// Sentencias
// ============================================================================

// LetStatement representa:  let <Name> = <Value>;
type LetStatement struct {
	Token token.Token // el token LET
	Name  *Identifier // nombre de la variable
	Value Expression  // expresión cuyo valor se guarda
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	return "let " + ls.Name.String() + " = " + stringOrEmpty(ls.Value) + ";"
}

// ReturnStatement representa:  return <ReturnValue>;
type ReturnStatement struct {
	Token       token.Token // el token RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	return "return " + stringOrEmpty(rs.ReturnValue) + ";"
}

// ExpressionStatement envuelve una expresión usada como sentencia.
// En Monkey una línea puede ser solo una expresión:  x + 10;
type ExpressionStatement struct {
	Token      token.Token // el primer token de la expresión
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string       { return stringOrEmpty(es.Expression) }

// BlockStatement es una lista de sentencias entre llaves: { ... }.
// Aparece en el cuerpo de las funciones y en las ramas de un if.
type BlockStatement struct {
	Token      token.Token // el token {
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out strings.Builder
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// ============================================================================
// Expresiones simples (hojas del árbol)
// ============================================================================

// Identifier es un nombre: x, suma, mi_lista.
type Identifier struct {
	Token token.Token // el token IDENT
	Value string      // el nombre
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral es un número entero escrito en el código: 42.
type IntegerLiteral struct {
	Token token.Token
	Value int64 // el valor ya convertido de texto a número
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// StringLiteral es un texto entre comillas: "hola".
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// Boolean es true o false.
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// ============================================================================
// Expresiones con operadores
// ============================================================================

// PrefixExpression: un operador ANTES de su operando:  -x   !ok
type PrefixExpression struct {
	Token    token.Token // el operador
	Operator string      // "-" o "!"
	Right    Expression  // el operando
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + stringOrEmpty(pe.Right) + ")"
}

// InfixExpression: un operador ENTRE dos operandos:  a + b   x == y
type InfixExpression struct {
	Token    token.Token // el operador
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return "(" + stringOrEmpty(ie.Left) + " " + ie.Operator + " " + stringOrEmpty(ie.Right) + ")"
}

// ============================================================================
// Expresiones compuestas
// ============================================================================

// IfExpression:  if (<Condition>) { <Consequence> } else { <Alternative> }
// En Monkey el if es una EXPRESIÓN: produce un valor.
type IfExpression struct {
	Token       token.Token // el token IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement // nil si no hay else
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	s := "if" + stringOrEmpty(ie.Condition) + " " + ie.Consequence.String()
	if ie.Alternative != nil {
		s += "else " + ie.Alternative.String()
	}
	return s
}

// FunctionLiteral:  fn(<Parameters>) { <Body> }
// Las funciones son valores: se pueden guardar en variables y pasar como argumento.
type FunctionLiteral struct {
	Token      token.Token // el token FUNCTION (fn)
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	params := make([]string, len(fl.Parameters))
	for i, p := range fl.Parameters {
		params[i] = p.String()
	}
	return fl.TokenLiteral() + "(" + strings.Join(params, ", ") + ") " + fl.Body.String()
}

// CallExpression:  <Function>(<Arguments>)
// Function puede ser un identificador (add) o un literal (fn(x){x}).
type CallExpression struct {
	Token     token.Token // el token (
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	return stringOrEmpty(ce.Function) + "(" + joinExpressions(ce.Arguments) + ")"
}

// ArrayLiteral:  [1, 2 * 2, "tres"]
type ArrayLiteral struct {
	Token    token.Token // el token [
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string       { return "[" + joinExpressions(al.Elements) + "]" }

// IndexExpression:  <Left>[<Index>]   por ejemplo  lista[0]  o  persona["nombre"]
type IndexExpression struct {
	Token token.Token // el token [
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	return "(" + stringOrEmpty(ie.Left) + "[" + stringOrEmpty(ie.Index) + "])"
}

// HashPair es una pareja llave: valor dentro de un HashLiteral.
type HashPair struct {
	Key   Expression
	Value Expression
}

// HashLiteral:  {"nombre": "Ada", "edad": 36}
//
// Decisión de diseño: guardamos las parejas en un SLICE (y no en un map de Go)
// para conservar el orden en que se escribieron. Así String() siempre imprime
// lo mismo y las pruebas son deterministas.
type HashLiteral struct {
	Token token.Token // el token {
	Pairs []HashPair
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	parts := make([]string, len(hl.Pairs))
	for i, pair := range hl.Pairs {
		parts[i] = stringOrEmpty(pair.Key) + ":" + stringOrEmpty(pair.Value)
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

// ============================================================================
// Ayudantes
// ============================================================================

// stringOrEmpty evita un panic si un nodo quedó incompleto por un error de
// sintaxis (por ejemplo, un let sin valor).
func stringOrEmpty(n Node) string {
	if n == nil {
		return ""
	}
	return n.String()
}

func joinExpressions(exps []Expression) string {
	parts := make([]string, len(exps))
	for i, e := range exps {
		parts[i] = stringOrEmpty(e)
	}
	return strings.Join(parts, ", ")
}

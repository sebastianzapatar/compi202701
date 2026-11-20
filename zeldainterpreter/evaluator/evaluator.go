// Package evaluator le da SIGNIFICADO al AST: lo recorre nodo por nodo
// ("tree-walking") y produce valores (object.Object).
//
// Toda la lógica cabe en una idea: Eval mira qué tipo de nodo recibe
// (con un type switch) y decide qué hacer; para los nodos con hijos, se
// llama a sí misma recursivamente.
package evaluator

import (
	"fmt"

	"zeldainterpreter/ast"
	"zeldainterpreter/object"
)

// Valores únicos (singletons). Solo existe un true, un false y un null:
//   - ahorramos memoria (no creamos miles de copias iguales), y
//   - podemos comparar booleanos por puntero:  left == right.
var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval evalúa cualquier nodo del AST dentro de un entorno de variables.
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// ---------------------------- Sentencias ----------------------------
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return nil // un let es una sentencia: no produce valor

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val} // lo envolvemos para que "suba"

	// ---------------------------- Literales ----------------------------
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.FunctionLiteral:
		// Una función es un valor. Guardamos el entorno ACTUAL: eso la vuelve un closure.
		return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}

	// ---------------------------- Expresiones ----------------------------
	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	}

	return nil
}

// ============================================================================
// Programa y bloques: el manejo de return y de errores
// ============================================================================

// evalProgram evalúa las sentencias de nivel superior.
// Aquí sí DESENVOLVEMOS el ReturnValue: ya no hay más niveles hacia arriba.
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range program.Statements {
		result = Eval(stmt, env)
		switch r := result.(type) {
		case *object.ReturnValue:
			return r.Value // desenvolver y terminar
		case *object.Error:
			return r // un error detiene todo
		}
	}
	return result
}

// evalBlockStatement evalúa un bloque { ... }.
// Aquí NO desenvolvemos el ReturnValue: lo devolvemos tal cual para que siga
// subiendo. Si lo desenvolviéramos, en
//
//	if (true) { if (true) { return 10; } return 1; }
//
// el bloque externo no se enteraría del return interno y devolvería 1.
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range block.Statements {
		result = Eval(stmt, env)
		if result != nil {
			if t := result.Type(); t == object.RETURN_VALUE_OBJ || t == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

// ============================================================================
// Operadores prefijos:  !x   -x
// ============================================================================

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

// evalBangOperatorExpression: "!" niega la "verdad" del valor.
// Solo false y null son falsos; todo lo demás (incluso 0) es verdadero.
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE, NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	integer, ok := right.(*object.Integer)
	if !ok {
		return newError("unknown operator: -%s", right.Type())
	}
	return &object.Integer{Value: -integer.Value}
}

// ============================================================================
// Operadores infijos:  a + b   a == b ...
// ============================================================================

// evalInfixExpression elige qué hacer según los TIPOS de ambos lados.
// ¡El orden de los casos importa! Los enteros se revisan primero porque se
// comparan por VALOR; el caso "==" de más abajo compara PUNTEROS, lo que
// solo sirve para los singletons TRUE, FALSE y NULL.
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left.(*object.Integer), right.(*object.Integer))

	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left.(*object.String), right.(*object.String))

	case operator == "==":
		return nativeBoolToBooleanObject(left == right)

	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right *object.Integer) object.Object {
	a, b := left.Value, right.Value
	switch operator {
	case "+":
		return &object.Integer{Value: a + b}
	case "-":
		return &object.Integer{Value: a - b}
	case "*":
		return &object.Integer{Value: a * b}
	case "/":
		// Igual que en el libro: dividir por cero provoca un panic de Go.
		// El runner lo atrapa con recover(). (Controlarlo aquí es un reto
		// de la guía de la Semana 6.)
		return &object.Integer{Value: a / b}
	case "<":
		return nativeBoolToBooleanObject(a < b)
	case ">":
		return nativeBoolToBooleanObject(a > b)
	case "==":
		return nativeBoolToBooleanObject(a == b)
	case "!=":
		return nativeBoolToBooleanObject(a != b)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalStringInfixExpression: los strings solo admiten concatenación con +.
func evalStringInfixExpression(operator string, left, right *object.String) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
	return &object.String{Value: left.Value + right.Value}
}

// ============================================================================
// Condicionales
// ============================================================================

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	switch {
	case isTruthy(condition):
		return orNull(Eval(ie.Consequence, env))
	case ie.Alternative != nil:
		return orNull(Eval(ie.Alternative, env))
	default:
		return NULL // if falso sin else: el valor de la expresión es null
	}
}

// isTruthy define qué valores cuentan como "verdaderos" en un if.
func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL, FALSE:
		return false
	default:
		return true // TRUE, cualquier entero (incluso 0), strings, arrays...
	}
}

// ============================================================================
// Identificadores
// ============================================================================

// evalIdentifier busca primero en el entorno y luego en las funciones
// integradas. Por eso un programa puede "tapar" un built-in: let len = ...;
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found: " + node.Value)
}

// evalExpressions evalúa una lista de expresiones de izquierda a derecha.
// Si alguna falla, devuelve una lista con ese único error.
func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	result := make([]object.Object, 0, len(exps))
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

// ============================================================================
// Llamadas a funciones y closures
// ============================================================================

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return orNull(unwrapReturnValue(evaluated))

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

// extendFunctionEnv crea el entorno de una llamada.
//
// Su outer es fn.Env (donde la función fue DEFINIDA), no el entorno de quien
// la llama. Esa sola decisión produce el ALCANCE LÉXICO y hace funcionar los
// closures.
func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for i, param := range fn.Parameters {
		// Igual que en el libro: si faltan argumentos, args[i] provoca un panic
		// (el runner lo atrapa). Validarlo es un reto de la guía de la Semana 7.
		env.Set(param.Value, args[i])
	}
	return env
}

// unwrapReturnValue quita el envoltorio del return al salir de la función,
// para que el return no "escape" y detenga también a quien llamó.
func unwrapReturnValue(obj object.Object) object.Object {
	if rv, ok := obj.(*object.ReturnValue); ok {
		return rv.Value
	}
	return obj
}

// ============================================================================
// Índices:  lista[0]   persona["nombre"]
// ============================================================================

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left.(*object.Array), index.(*object.Integer))
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left.(*object.Hash), index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

// evalArrayIndexExpression: un índice fuera de rango devuelve null (no error).
func evalArrayIndexExpression(array *object.Array, index *object.Integer) object.Object {
	i := index.Value
	if i < 0 || i >= int64(len(array.Elements)) {
		return NULL
	}
	return array.Elements[i]
}

// ============================================================================
// Hashes
// ============================================================================

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	hash := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}

	for _, pair := range node.Pairs {
		key := Eval(pair.Key, env)
		if isError(key) {
			return key
		}
		hashable, ok := key.(object.Hashable) // ¿este valor puede ser llave?
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(pair.Value, env)
		if isError(value) {
			return value
		}

		hk := hashable.HashKey()
		if _, exists := hash.Pairs[hk]; !exists {
			hash.Keys = append(hash.Keys, hk) // recordamos el orden de inserción
		}
		hash.Pairs[hk] = object.HashPair{Key: key, Value: value}
	}
	return hash
}

func evalHashIndexExpression(hash *object.Hash, index object.Object) object.Object {
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}
	pair, ok := hash.Pairs[key.HashKey()]
	if !ok {
		return NULL // llave inexistente
	}
	return pair.Value
}

// ============================================================================
// Ayudantes
// ============================================================================

func nativeBoolToBooleanObject(b bool) *object.Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR_OBJ
}

// orNull convierte "sin valor" (nil de Go) en NULL de Monkey.
// Ocurre, por ejemplo, con un bloque vacío o uno que termina en un let:
// fn() { let x = 1; }() produce null en vez de dejar un nil suelto que
// podría provocar un panic más adelante.
func orNull(obj object.Object) object.Object {
	if obj == nil {
		return NULL
	}
	return obj
}

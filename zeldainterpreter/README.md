# ZeldaInterpreter

Intérprete completo del lenguaje **Monkey**, escrito en Go para el curso *Lenguajes y Compiladores 2027-1*.
Sigue la arquitectura de *Writing An Interpreter In Go* (Thorsten Ball): lexer → parser Pratt → evaluador *tree-walking*.
Cada archivo está comentado en español para que se pueda leer como material de estudio.

```
código fuente ──► lexer ──► tokens ──► parser ──► AST ──► evaluador ──► resultado
   "1 + 2"        lexer/               parser/    ast/    evaluator/    object/
```

---

## 1. Requisitos y primeros pasos

Necesitas **Go 1.23 o superior** (`go version`).

```bash
cd zeldainterpreter
go test ./...           # corre todas las pruebas (deben salir en "ok")
go run .                # abre la consola interactiva
go build -o zelda .     # compila el ejecutable "zelda" (en Windows: zelda.exe)
```

---

## 2. Modos de uso

| Cómo se llama | Qué hace |
|---|---|
| `zelda` | Abre la **consola interactiva** (REPL). |
| `zelda programas/fibonacci.mk` | **Ejecuta un archivo** `.mk`. |
| `zelda -` | Ejecuta lo que llegue por la **entrada estándar**: `cat prog.mk \| zelda -` |
| `zelda -e "puts(1 + 2)"` | Ejecuta **código escrito en la línea de comandos**. |

Y en cualquiera de ellos puedes elegir **qué fase del intérprete ver** con `-modo`:

| `-modo=` | Fases que se ejecutan | Qué muestra |
|---|---|---|
| `eval` *(por defecto)* | lexer + parser + evaluador | el resultado del programa |
| `tokens` | lexer | la tabla de tokens |
| `ast` | lexer + parser | el programa con paréntesis explícitos: `(1 + (2 * 3))` |
| `arbol` | lexer + parser | el AST dibujado como árbol |

Otras opciones: `-resultado` (al ejecutar un archivo, muestra también el valor final), `-version` y `-h` (ayuda).

> Las opciones van **antes** del nombre del archivo: `zelda -modo=arbol programa.mk`.

Ejemplos:

```bash
$ zelda -modo=arbol -e "let x = 1 + 2 * 3;"
Program
└── LetStatement
    ├── nombre: Identifier x
    └── valor: InfixExpression +
        ├── izq: IntegerLiteral 1
        └── der: InfixExpression *
            ├── izq: IntegerLiteral 2
            └── der: IntegerLiteral 3

$ zelda -modo=tokens -e "x != 10"
  1  IDENT      "x"
  2  !=         "!="
  3  INT        "10"
  4  EOF        ""
```

**Códigos de salida:** `0` todo bien · `1` el programa tuvo errores (sintaxis o ejecución) · `2` uso incorrecto del comando.

---

## 3. La consola interactiva

```
$ go run .
ZeldaInterpreter · intérprete del lenguaje Monkey
Modo actual: eval. Escribe :ayuda para ver los comandos, :salir para terminar.
>> let doble = fn(x) {
..   x * 2
.. };
>> doble(21)
42
>> :modo arbol
Modo cambiado a: arbol
>> -x
Program
└── ExpressionStatement
    └── PrefixExpression -
        └── operando: Identifier x
```

| Comando | Qué hace |
|---|---|
| `:ayuda` | lista los comandos |
| `:modo <m>` | cambia el modo: `eval`, `tokens`, `ast` o `arbol` |
| `:env` | muestra las variables definidas en la sesión |
| `:cargar <ruta>` | ejecuta un archivo dentro de la sesión actual (sus variables quedan disponibles) |
| `:reiniciar` | borra todas las variables |
| `:salir` | termina (también `Ctrl+D`) |

Si una línea deja `{`, `(` o `[` sin cerrar, la consola muestra `..` y sigue leyendo: así puedes escribir funciones en varias líneas.

---

## 4. Estructura del proyecto (y en qué semana se ve cada parte)

Te recomendamos **leer el código en este orden**: es el mismo orden en que los datos atraviesan el intérprete.

| # | Carpeta | Qué contiene | Clase |
|---|---|---|---|
| 1 | `token/` | Los tipos de token y las palabras reservadas. | Semana 3 |
| 2 | `lexer/` | El análisis léxico: texto → tokens (`NextToken`, `peekChar`). | Semana 3 |
| 3 | `ast/` | Los nodos del árbol (`ast.go`) y el dibujante de árboles (`tree.go`). | Semanas 4–5 |
| 4 | `parser/` | Descenso recursivo para sentencias + Pratt parsing para expresiones. | Semanas 4–5 |
| 5 | `object/` | Los valores en tiempo de ejecución (`object.go`) y las variables (`environment.go`). | Semanas 6–8 |
| 6 | `evaluator/` | El evaluador tree-walking (`evaluator.go`) y las funciones integradas (`builtins.go`). | Semanas 6–8 |
| 7 | `runner/` | Lleva el código por el pipeline según el modo; atrapa los `panic`. | — |
| 8 | `repl/` | La consola interactiva y sus comandos. | Semanas 3 y 5 |
| 9 | `main.go` | La línea de comandos: archivo, stdin, `-e`, `-modo`. | — |
| — | `programas/` | Programas `.mk` de ejemplo. | — |

Cada paquete tiene su archivo `_test.go`. Las pruebas también son documentación: muestran con ejemplos concretos qué debe hacer cada parte.

```
zeldainterpreter/
├── go.mod
├── main.go            main_test.go
├── token/token.go
├── lexer/lexer.go     lexer_test.go
├── ast/ast.go         tree.go          tree_test.go
├── parser/parser.go   parser_test.go
├── object/object.go   environment.go   environment_test.go
├── evaluator/evaluator.go   builtins.go   evaluator_test.go
├── runner/runner.go
├── repl/repl.go
└── programas/*.mk
```

---

## 5. El lenguaje en una página

```js
let edad = 20;                                  // variables (enteros, booleanos, strings)
let mensaje = if (edad > 17) { "adulto" } else { "menor" };   // if es una expresión

let suma = fn(a, b) { a + b };                  // las funciones son valores
let newAdder = fn(x) { fn(y) { x + y } };       // closures
newAdder(2)(40);                                // => 42

let lista = [1, 2, 3];                          // arrays (inmutables)
let persona = {"nombre": "Ada", "edad": 36};    // hashes
persona["nombre"];                              // => Ada
```

**Funciones integradas:** `len`, `first`, `last`, `rest`, `push`, `puts`.

> Monkey no tiene comentarios. Las líneas `// ...` de arriba son solo explicación para este README.

---

## 6. Programas de ejemplo

```bash
go run . programas/hola.mk             # strings y concatenación
go run . programas/aritmetica.mk       # precedencia, booleanos, if como expresión
go run . programas/fibonacci.mk        # recursión
go run . programas/closures.mk         # closures y composición de funciones
go run . programas/map_reduce.mk       # funciones de orden superior (el "gran final" del libro)
go run . programas/hashes.mk           # arrays de hashes
go run . programas/error_ejecucion.mk  # cómo se reporta un error en tiempo de ejecución
go run . programas/error_sintaxis.mk   # cómo se reportan los errores de sintaxis
```

---

## 7. Decisiones de diseño

- **Mismos mensajes de error que el libro** (`type mismatch: INTEGER + BOOLEAN`, `expected next token to be =, got INT instead`…), para que coincidan con las clases y las guías de ejercicios.
- **Los hashes conservan el orden de inserción** (`ast.HashLiteral` usa un slice y `object.Hash` guarda `Keys`), así `Inspect()` siempre imprime lo mismo.
- **Las funciones de parsing de sentencias devuelven la interfaz `ast.Statement`**, para evitar la trampa de Go de una interfaz no-nil que contiene un puntero nil (explicada en `parser.go`).
- **Los `panic` no cierran la consola:** dividir por cero o llamar una función con menos argumentos provoca un `panic` (igual que en el libro), pero `runner.Run` lo atrapa con `defer`/`recover` y lo informa.

### Lo que NO incluye (a propósito)

Estas extensiones son retos de las guías de ejercicios, así que quedan para que las implementen los estudiantes: comentarios `//`, operadores `<=` y `>=`, números de línea y columna en los errores, `ERROR: division by zero`, validación de la cantidad de argumentos, `%`, `**`, `&&` y `||`, asignación sin `let` y ciclos.

Una recursión infinita (`let f = fn() { f() }; f();`) agota la pila de Go y termina el proceso: es un error fatal que `recover` no puede atrapar. Detenerla con `context` es el reto del módulo Go IV.

---

## 8. Relación con el código del libro

Este intérprete implementa el mismo lenguaje y la misma arquitectura que *Writing An Interpreter In Go*. Conserva los nombres de los tipos y funciones principales (`NextToken`, `ParseProgram`, `parseExpression`, `Eval`…) para que puedas seguirlo junto con las clases. El módulo se llama `zeldainterpreter`; en las clases, siguiendo el libro, los imports son `monkey/...`.

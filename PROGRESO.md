# 📋 Planeación y Seguimiento — ZeldaInterpreter · Lenguajes y Compiladores 2027-1

> **Libro base:** *Writing An Interpreter In Go* — Thorsten Ball (Ed. 1.7, 2020)  
> **Lenguaje a construir:** **Monkey** (escrito desde cero en **Go**, sin librerías externas)  
> **Duración del curso:** **8 semanas** (2 semanas de Go + 6 semanas de construcción del intérprete)  
> **Formato de material:** Portal web estático híbrido (**Modo Guía + Modo Presentación / Slides**) listo para GitHub Pages, usando la paleta y estilo visual de `compi20262`.  
> **Sistema de diseño (ui-ux-pro-max):** estilo documentación técnica oscura (Minimalism/Swiss), iconos SVG Lucide (sin emojis), foco visible, skip-link, `prefers-reduced-motion`, menú móvil, sin ligaduras en código.

---

## 📊 Estado General del Proyecto

| Componente | Archivo | Estado | Notas |
| :--- | :--- | :---: | :--- |
| Diseño visual y modo slides | `assets/styles.css` | ✅ Hecho | Paleta `#0f172a` / `#38bdf8` / `#f472b6`, responsive + modo presentación |
| Controlador híbrido JS | `assets/presentation.js` | ✅ Hecho | Toggle Guía ↔ Slides (`P`), teclado (`←`/`→`/`Esc`), copiar código, scroll-spy |
| Portal principal del curso | `index.html` | ✅ Hecho | Hero, pipeline, preview de Monkey, roadmap de 8 semanas |
| Documento de seguimiento | `PROGRESO.md` | ✅ Hecho | Plan completo de 8 semanas + checklist de avance |
| **Semana 1:** Fundamentos de Go | `clases/semana-01-go-fundamentos.html` | ✅ Hecho | Tipos, control de flujo, slices, maps, primer proyecto Go + taller |
| **Semana 2:** Go para Intérpretes | `clases/semana-02-go-intermedio.html` | ✅ Hecho | Structs, punteros, interfaces, type switches, closures, TDD (`go test`) |
| **Go III:** POO en Go | `clases/go-03-poo-composicion.html` | ✅ Hecho | Objetos, receivers/method sets, agregación vs composición, embedding (“herencia”), polimorfismo, SOLID + 2 simuladores |
| **Go IV:** Concurrencia | `clases/go-04-concurrencia.html` | ✅ Hecho | Goroutines, Mutex, canales, select, context, pipeline/worker pool + simuladores de carrera y de canal |
| **Go V:** Patrones I | `clases/go-05-patrones-creacionales-estructurales.html` | ✅ Hecho | 12 patrones creacionales y estructurales + explorador interactivo |
| **Go VI:** Patrones II | `clases/go-06-patrones-comportamiento.html` | ✅ Hecho | 10 patrones de comportamiento + explorador interactivo |
| **Intérprete completo** | `zeldainterpreter/` | ✅ Hecho | token, lexer, ast (+árbol), parser Pratt, object, evaluator, runner, repl, main; modos eval/tokens/ast/arbol; archivo, stdin, `-e` y REPL con comandos; pruebas en cada paquete |
| Ejemplos Go ejecutables (29) | `ejemplos/` | ✅ Hecho | POO (4), concurrencia (3), patrones (22); compilados, con `go vet` y su salida real en `salida.txt` |
| **Semana 3:** Lexing (Cap. 0–1) | `clases/semana-03-lexing.html` | ✅ Hecho | Tokens, Lexer, `peekChar()`, REPL + simulador interactivo de Lexer |
| **Semana 4:** Parsing I (Cap. 2.1–2.5) | `clases/semana-04-parsing-sentencias.html` | ✅ Hecho | AST, `let`/`return`, `expectPeek`, errores + simulador paso a paso de `curToken`/`peekToken` |
| **Semana 5:** Parsing II (Cap. 2.6–2.9) | `clases/semana-05-pratt-parsing.html` | ✅ Hecho | Pratt, precedencias, `if`, `fn`, llamadas, RPPL + visualizador interactivo de traza y AST |
| **Semana 6:** Evaluación I (Cap. 3.1–3.8) | `clases/semana-06-evaluacion-base.html` | ✅ Hecho | Objetos, singletons, prefijos/infijos, truthiness, `return`, errores + REPL en vivo |
| **Semana 7:** Evaluación II (Cap. 3.9–3.11) | `clases/semana-07-entornos-closures.html` | ✅ Hecho | `Environment`, `outer`, funciones, closures (diagramas), GC + REPL con visor de entorno |
| **Semana 8:** Extensiones (Cap. 4 + Cierre) | `clases/semana-08-extensiones.html` | ✅ Hecho | Strings, built-ins, arrays, hashes (`HashKey`), Grand Finale + playground completo |
| Navegación compartida del curso | `assets/course.js` | ✅ Hecho | **Fuente única** de semanas/rutas/PDFs: sidebar, anterior/siguiente, roadmap e iconos SVG |
| Simulador ZeldaInterpreter (todas las páginas) | `assets/simulator.js` | ✅ Hecho | Pestañas Tokens / AST / Evaluación, estado por semana, ejemplos y entorno global |
| Intérprete Monkey en JS | `assets/monkey.js` + `assets/repl-widget.js` | ✅ Hecho | Réplica del intérprete del libro (lexer, Pratt con traza, evaluador) para los demos |
| Guías de ejercicios (8) | `ejercicios/pdf/semana-0N-ejercicios.pdf` | ✅ Hecho | Partes A (conceptos), B (trazas), C (Go), reto y rúbrica; generadas con `scripts/build_ejercicios.py` |
| Solucionarios (8, docente) | `ejercicios/solucionario/` | ✅ Hecho | Mismas guías con respuestas verificadas; **no publicar** a estudiantes |
| **Córrelo en ZeldaInterpreter** (Sem 3–8 + portal) | `scripts/build_correr.py` | ✅ Hecho | Sección en cada clase con comandos `go run .` (archivo, `-e`, stdin, `-modo`, REPL) y su **salida real**; el portal tiene la sección “7. Cómo correr ZeldaInterpreter” |

---

## 🗓️ Planeación Detallada de las 8 Semanas

### Bloque I — Dominando Go para Construir Compiladores (Semanas 1–2)

#### ✅ Semana 1: Fundamentos de Go (Parte 1)
- **Objetivo:** Que el estudiante configure su entorno Go y domine los tipos y colecciones básicas que usará el Lexer y el Parser.
- **Temas cubiertos:**
  - ¿Por qué Go para escribir un intérprete? (Intro del libro, pp. 8–9).
  - Instalación, `go mod init monkey`, estructura de paquetes y visibilidad (Mayúscula = exportado).
  - Variables (`var`, `:=`), valores cero (*zero values*) y constantes (`const`).
  - Tipos clave para compiladores: `byte` (ASCII), `string` (slicing `input[pos:end]`), `int64`, `bool`.
  - Control de flujo: `if`, `for` (único ciclo en Go), `switch` (sin `break` obligatorio).
  - Funciones, múltiples valores de retorno y manejo idiomático de errores.
  - Colecciones esenciales: **Slices** (`[]T`, `append`) y **Maps** (`map[K]V`, patrón `val, ok := m[k]`).
- **Entregable / Taller:** Mini contador de palabras y frecuencias usando `map[string]int` y slicing.

#### ✅ Semana 2: Go Intermedio para Intérpretes (Parte 2)
- **Objetivo:** Dominar los 4 pilares de Go que sostienen la arquitectura de Monkey: `struct` + punteros, interfaces, funciones como valores y pruebas unitarias.
- **Temas cubiertos:**
  - `struct` y funciones constructoras (`New(...) *Tipo`).
  - Métodos con *value receiver* vs *pointer receiver* (por qué `*Lexer` y `*Parser` mutan estado).
  - Punteros en Go (`&`, `*`, `nil`) sin aritmética peligrosa.
  - Interfaces implícitas (`ast.Node`, `object.Object`, `io.Reader`/`io.Writer`).
  - *Type assertions* y *Type switches* (`switch n := node.(type)`) — base del Evaluador.
  - Funciones de primera clase, tipos de función (`type prefixParseFn func() ast.Expression`) y closures.
  - Testing nativo con el paquete `testing` y **Table-Driven Tests** (`go test -v ./...`).
- **Entregable / Taller:** Mini evaluador de operaciones aritméticas con interfaz `Operacion` y pruebas table-driven.

---

### Bloque II — Construcción del Intérprete Monkey (Semanas 3–8)

#### ✅ Semana 3: Introducción a Intérpretes y Análisis Léxico (Intro + Cap. 1, pp. 5–27)
- **Objetivo:** Entender el pipeline de un intérprete y construir completamente los paquetes `token`, `lexer` y el primer `repl`.
- **Temas cubiertos:**
  - Intérprete vs Compilador; tipos de intérpretes (Tree-Walking, Bytecode VM, JIT).
  - El lenguaje Monkey y pipeline de 3 fases.
  - `token/token.go`: `TokenType`, `Token`, constantes (`ILLEGAL`, `EOF`, operadores, delimitadores, keywords) y `LookupIdent()`.
  - `lexer/lexer.go`: estructura con 4 campos (`input`, `position`, `readPosition`, `ch`), `readChar()`, `NextToken()`, `skipWhitespace()`, `readIdentifier()`, `readNumber()`.
  - Tokens de dos caracteres (`==`, `!=`) mediante `peekChar()`.
  - Pruebas unitarias en `lexer/lexer_test.go` y primer REPL tokenizador (`repl/repl.go` + `main.go`).
  - **Extra incluido:** Simulador interactivo de Lexer en el navegador dentro de la presentación.

#### ✅ Semana 4: Análisis Sintáctico I — AST y Sentencias (Cap. 2.1 – 2.5, pp. 28–45)
- **Objetivo:** Comprender qué es un AST y construir la base del Parser recursivo descendente para sentencias `let` y `return`.
- **Temas cubiertos:**
  - Qué es un Parser (analogía con `JSON.parse`) y qué hace "abstracto" al AST.
  - ¿Por qué escribir un Parser a mano en lugar de usar yacc/ANTLR?
  - Paquete `ast`: interfaces `Node`, `Statement`, `Expression` y nodo raíz `Program`.
  - Paquete `parser`: `curToken`, `peekToken`, `nextToken()`, `expectPeek()`, recolección de errores (`Errors()`).
  - Parsing de `LetStatement` y `ReturnStatement` + pruebas en `parser_test.go`.

#### ✅ Semana 5: Análisis Sintáctico II — Pratt Parsing y Expresiones (Cap. 2.6 – 2.9, pp. 46–102)
- **Objetivo:** Dominar el algoritmo de Pratt (*Top-Down Operator Precedence*) y parsear todas las expresiones de Monkey.
- **Temas cubiertos:**
  - Sentencias vs Expresiones en Monkey.
  - Anatomía de Pratt Parsing: precedencias (`LOWEST`..`CALL`), mapas `prefixParseFns` e `infixParseFns`.
  - Traza visual paso a paso de `parseExpression(precedence)` sobre expresiones como `1 + 2 * 3`.
  - Implementación de: `Identifier`, `IntegerLiteral`, `PrefixExpression` (`!`, `-`), `InfixExpression`.
  - Extensión: `Boolean`, expresiones agrupadas `( )`, `IfExpression`, `FunctionLiteral`, `CallExpression`.
  - Actualización del REPL a Read-Parse-Print-Loop (RPPL) + cara de Monkey en errores de sintaxis.
  - **Extra incluido:** Visualizador interactivo paso a paso de precedencia / AST.

#### ✅ Semana 6: Evaluación I — Tree-Walking, Objetos y Control de Flujo (Cap. 3.1 – 3.8, pp. 103–135)
- **Objetivo:** Darle vida al AST construyendo el sistema de objetos y el núcleo del evaluador.
- **Temas cubiertos:**
  - Estrategias de evaluación (Tree-Walking vs Bytecode vs JIT).
  - Paquete `object`: interfaz `Object`, `Integer`, `Boolean` (singletons `TRUE`/`FALSE`), `Null` (`NULL`).
  - Paquete `evaluator`: `Eval(node ast.Node)` y recorrido recursivo del AST.
  - Evaluación de prefijos (`!`, `-`), infijos (enteros y booleanos), condicionales (`if/else`, *truthiness*).
  - Sentencias `return` (`ReturnValue` y propagación en bloques anidados).
  - Manejo de errores en tiempo de ejecución (`object.Error` y cortocircuito).

#### ✅ Semana 7: Evaluación II — Entornos, Funciones y Closures (Cap. 3.9 – 3.11, pp. 136–151)
- **Objetivo:** Implementar memoria léxica (`Environment`), llamadas a función y closures reales.
- **Temas cubiertos:**
  - Bindings (`let`) y `object/environment.go` (`Get`, `Set`, `NewEnclosedEnvironment`).
  - Evaluación de identificadores y propagación del entorno en `Eval(node, env)`.
  - `object.Function`: parámetros, cuerpo y entorno capturado.
  - Evaluación de llamadas: `evalExpressions`, `applyFunction`, `extendFunctionEnv`, `unwrapReturnValue`.
  - Visualización de Closures ("hojas encadenadas de entorno") y funciones de orden superior.
  - *Who's taking the trash out?*: recolección de basura delegada al runtime de Go.

#### ✅ Semana 8: Extendiendo el Intérprete — Strings, Built-ins, Arrays y Hashes (Cap. 4, pp. 152–201)
- **Objetivo:** Practicar el ciclo completo de agregar nuevas características de extremo a extremo y cerrar el curso con un intérprete completo.
- **Temas cubiertos:**
  - **Strings:** `token.STRING`, `readString()`, `StringLiteral`, `object.String`, concatenación `+`.
  - **Built-in Functions:** `object.Builtin`, `len`, `first`, `last`, `rest`, `push`, `puts`.
  - **Arrays:** `[ ... ]`, `ArrayLiteral`, `IndexExpression` (`arr[i]`), precedencia `INDEX`.
  - **Hashes:** `{ k: v }`, interfaz `Hashable` y `HashKey` (por qué comparar por valor y no por puntero), `HashLiteral`, indexación.
  - **The Grand Finale:** `map`, `reduce` y `fibonacci` corriendo en vivo en nuestro REPL.
  - Cierre: qué sigue (*The Lost Chapter* — macros, y *Writing A Compiler In Go* — Bytecode + VM).

---

## ✅ Checklist de Pendientes Próximos

- [x] Verificar las Semanas 4–8 (2026-11-20): 13 páginas sin errores de JS, enlaces y PDFs correctos, `go test ./...` en verde, `monkey.js` da los mismos resultados que el intérprete Go, y todas las respuestas de las tablas de las guías 4–8 coinciden con `zeldainterpreter`.
- [ ] Revisar con el docente el contenido de las Semanas 4–8 y las 8 guías de ejercicios (puntajes, dificultad, fechas de entrega).
- [ ] Decidir cómo se entregan los solucionarios (están en `ejercicios/solucionario/` y excluidos en `.gitignore`).
- [ ] Decidir cómo encajan en el calendario los módulos Go III–VI (¿semanas adicionales, clases extra o autoestudio?).
- [ ] (Opcional) Agregar fechas reales del semestre 2027-1 a cada semana en `assets/course.js`.
- [ ] (Opcional) Publicar en GitHub Pages y probar en el proyector del salón (modo presentación + pantalla completa).
- [ ] (Opcional) Parciales / quices a partir de las Partes A y B de las guías.

## 🔁 Cómo mantener el curso

- **Agregar o renombrar una clase:** editar solo el arreglo `WEEKS` en `assets/course.js` (sidebar, anterior/siguiente y roadmap se actualizan solos).
- **Editar ejemplos Go:** modificar `ejemplos/**/main.go` y ejecutar `python3 scripts/embed_codigo.py --run` (gofmt + vet + run + incrusta código y salida en las clases).
- **Editar patrones:** los textos están en `scripts/build_patrones.py`; regenerar con `python3 scripts/build_patrones.py && python3 scripts/embed_codigo.py`.
- **Editar los ejemplos de “Córrelo en ZeldaInterpreter”:** modificar `PAGES` en `scripts/build_correr.py` y ejecutar `python3 scripts/build_correr.py` (compila el intérprete, corre cada comando y reescribe la sección entre `<!-- CORRER:INICIO -->` y `<!-- CORRER:FIN -->`). Vuelve a correrlo si cambias `zeldainterpreter/` o `programas/`.
- **Editar ejercicios:** modificar `scripts/ejercicios_contenido.py` y ejecutar `python3 scripts/build_ejercicios.py` (regenera HTML, PDF y solucionario).

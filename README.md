# ZeldaInterpreter — Lenguajes y Compiladores 2027-1

Intérprete del lenguaje **Monkey** escrito en Go durante el curso.

Material del curso basado en *Writing An Interpreter In Go* (Thorsten Ball): 8 clases híbridas (guía de estudio + diapositivas) y una guía de ejercicios en PDF por semana.

## Cómo usarlo

- Abre `index.html` en el navegador (o publícalo en GitHub Pages).
- En cualquier clase, pulsa **`P`** para el modo presentación. Navega con `←` / `→`, `Esc` para salir y el botón de pantalla completa para proyectar.
- Cada clase tiene un botón **Ejercicios de la semana (PDF)** en el menú lateral y al final, en la diapositiva del taller.

## Estructura

```
index.html                     Portal: visión general, pipeline, roadmap y guías
clases/
  semana-01-go-fundamentos.html     Go I
  semana-02-go-intermedio.html      Go II
  go-03-poo-composicion.html        Go III   · objetos, agregación, composición, embedding
  go-04-concurrencia.html           Go IV    · goroutines, canales, select, context
  go-05-patrones-creacionales-estructurales.html  Go V · 12 patrones
  go-06-patrones-comportamiento.html               Go VI · 10 patrones
  semana-03-lexing.html             Cap. 1   · simulador de Lexer
  semana-04-parsing-sentencias.html Cap. 2A  · simulador de expectPeek
  semana-05-pratt-parsing.html      Cap. 2B  · visualizador de Pratt + AST
  semana-06-evaluacion-base.html    Cap. 3A  · REPL en vivo
  semana-07-entornos-closures.html  Cap. 3B  · REPL con visor de entorno
  semana-08-extensiones.html        Cap. 4   · playground completo
assets/
  styles.css        Sistema visual (tokens, componentes, modo slides, móvil, impresión)
  course.js         FUENTE ÚNICA del curso: semanas, rutas, PDFs, iconos SVG, marca y logo pixel-art
  presentation.js   Modo guía ↔ modo presentación, menú móvil, copiar código
  monkey.js         Réplica en JS del intérprete del libro (lexer, Pratt, evaluador)
  repl-widget.js    REPL reutilizable de las semanas 6–8
  simulator.js      Simulador ZeldaInterpreter (Tokens → AST → Evaluación), presente en todas las páginas
zeldainterpreter/    EL INTÉRPRETE COMPLETO en Go (token, lexer, ast, parser, object, evaluator, runner, repl)
ejemplos/            Programas Go ejecutables (poo/, concurrencia/, patrones/) + salida.txt
ejercicios/
  pdf/              Guías para estudiantes (semana-0N-ejercicios.pdf)
  solucionario/     Guías con respuestas — SOLO DOCENTE (excluido en .gitignore)
  fuente/           HTML imprimible de cada guía + worksheet.css
scripts/
  ejercicios_contenido.py   Preguntas, trazas, problemas, respuestas y rúbricas
  build_ejercicios.py       Genera HTML + PDF + solucionario (usa Chrome headless)
  build_patrones.py         Genera las clases Go V/VI y assets/patrones.js (explorador)
  embed_codigo.py           Incrusta el código de ejemplos/ y su salida real en las clases (--run para ejecutarlos)
  build_correr.py           Genera «Córrelo en ZeldaInterpreter» (Sem 3–8 y portal) con comandos y salida real
PROGRESO.md                 Planeación de 8 semanas y seguimiento
```

## Mantenimiento

**Agregar o renombrar una clase:** edita solo `WEEKS` en `assets/course.js`. El menú lateral, los botones anterior/siguiente, el roadmap del portal y los enlaces a los PDF se generan a partir de ahí. Cada página de clase declara su número con `<body data-week="N">`.

**Editar o regenerar los ejercicios:**

```bash
python3 scripts/build_ejercicios.py        # todas las semanas
python3 scripts/build_ejercicios.py 5      # solo la semana 5
python3 scripts/build_ejercicios.py --html # solo HTML, sin PDF
```

Requiere Python 3 y Google Chrome instalado. Las salidas esperadas de los programas Monkey se verificaron con `assets/monkey.js`.

**Estructura de cada guía:** Parte A · Conceptos (25 pts), Parte B · Trazas a mano (35 pts), Parte C · Programación en Go (40 pts), un reto opcional (+10) y la rúbrica. La semana 8 incluye el enunciado y la rúbrica del proyecto final.

**Simulador ZeldaInterpreter:** cualquier página puede incluirlo con `<div data-simulator data-stage="tokens|ast|eval" data-code="..."></div>` (más `monkey.js` y `simulator.js`). Las pestañas indican qué fase ya construyeron los estudiantes según `<body data-week>`.

**Nombre del módulo Go:** el código de las clases conserva `go mod init monkey` e imports `monkey/token`, etc., para coincidir con el libro. ZeldaInterpreter es el nombre del proyecto/intérprete.

**Cómo correr el intérprete en las clases:** cada clase de las semanas 3–8 tiene la sección *Córrelo en ZeldaInterpreter* (y el portal, *7. Cómo correr ZeldaInterpreter*) con comandos `go run .` y su salida real. Se generan con `python3 scripts/build_correr.py` (edita los ejemplos en `PAGES` dentro del script; requiere Go). Regenera después de cambiar `zeldainterpreter/`.

**Ejemplos Go:** `cd ejemplos && go run ./patrones/builder` (o cualquier carpeta). Tras editar un ejemplo: `python3 scripts/embed_codigo.py --run`. Los patrones siguen el catálogo GoF; como referencia se enlaza https://refactoring.guru/es/design-patterns/catalog (textos y código del curso son propios).

## El intérprete: `zeldainterpreter/`

Implementación completa y comentada en español, con pruebas en cada paquete. Ver `zeldainterpreter/README.md`.

```bash
cd zeldainterpreter
go test ./...                          # pruebas
go run .                               # consola interactiva (:ayuda para ver comandos)
go run . programas/fibonacci.mk        # ejecutar un archivo
go run . -modo=arbol -e "1 + 2 * 3"    # ver el AST (modos: eval, tokens, ast, arbol)
cat programas/hola.mk | go run . -     # leer desde la entrada estándar
```

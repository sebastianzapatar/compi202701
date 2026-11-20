#!/usr/bin/env python3
"""
Genera las clases de patrones de diseño (Go V y Go VI) y assets/patrones.js
a partir de los datos de este archivo. El código Go de cada patrón vive en
ejemplos/patrones/<slug>/main.go y se incrusta con scripts/embed_codigo.py.

Uso:  python3 scripts/build_patrones.py && python3 scripts/embed_codigo.py
"""
import json
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
REF = "https://refactoring.guru/es/design-patterns/catalog"

# slug, nombre, categoría, intención, problema, usar, evitar, zelda, idea Go
P = [
    # ---------------- CREACIONALES ----------------
    ("factory-method", "Factory Method", "creacional",
     "Delegar en una función o método la decisión de qué tipo concreto crear, devolviendo una abstracción.",
     "El REPL debe imprimir resultados en texto plano, JSON o en un recuadro según una opción, sin llenar el código de <code>if</code>.",
     "Cuando el cliente no debe conocer los tipos concretos o cuando la elección depende de configuración.",
     "Si solo existe un tipo concreto: un constructor <code>NewX</code> basta.",
     "<code>lexer.New</code>, <code>parser.New</code> y <code>object.NewEnvironment</code> son fábricas simples; <code>NewPrinter(format)</code> lo sería con varios formatos.",
     "En Go suele ser una función <code>NewX(...) (Interfaz, error)</code>, sin jerarquías de creadores."),
    ("abstract-factory", "Abstract Factory", "creacional",
     "Crear familias completas de objetos relacionados que están diseñados para usarse juntos.",
     "El REPL tiene prompt y formato de errores; con colores ANSI o sin ellos. Mezclar un prompt con color y errores sin color sería inconsistente.",
     "Cuando hay varias familias intercambiables (temas, plataformas, backends) y hay que garantizar coherencia.",
     "Si solo hay un producto por familia: un Factory Method es suficiente.",
     "Un <code>Theme</code> del REPL que produce <code>Prompt</code> y <code>ErrorFormatter</code> coherentes entre sí.",
     "Una interfaz con varios métodos <code>NewX()</code>; cada familia es un struct vacío que la implementa."),
    ("builder", "Builder", "creacional",
     "Construir un objeto complejo paso a paso, separando la construcción de su representación final.",
     "Las pruebas del parser necesitan programas y ASTs de ejemplo; armarlos con structs anidados es largo y frágil.",
     "Objetos con muchas partes opcionales, validación al final o construcción en varios pasos.",
     "Objetos simples con dos o tres campos: un literal de struct es más claro.",
     "Un <code>ProgramBuilder</code> para tests y <code>strings.Builder</code>/<code>bytes.Buffer</code> en todos los <code>String()</code> del AST.",
     "API fluida que devuelve <code>*Builder</code>; alternativa idiomática: <em>functional options</em> (<code>New(WithTrace(true))</code>)."),
    ("prototype", "Prototype", "creacional",
     "Crear nuevos objetos clonando uno existente en lugar de construirlos desde cero.",
     "<code>push</code> de Monkey no debe modificar el array original: necesita una copia del array antes de agregar.",
     "Cuando copiar es más simple que reconstruir, o cuando se necesitan copias independientes de estructuras anidadas.",
     "Si los objetos son inmutables: se pueden compartir sin copiar.",
     "<code>push</code> y <code>rest</code> del Cap. 4 crean arrays nuevos; un método <code>Clone()</code> generaliza la idea.",
     "Un método <code>Clone() T</code>; cuidado con la diferencia entre copia superficial y profunda de slices y maps."),
    ("singleton", "Singleton", "creacional",
     "Garantizar que exista una única instancia de un tipo y ofrecer un punto de acceso global a ella.",
     "Crear un <code>&amp;Boolean{true}</code> nuevo en cada comparación desperdicia memoria e impide comparar por puntero.",
     "Valores únicos por naturaleza (TRUE, FALSE, NULL) o registros globales de solo lectura.",
     "Para estado mutable compartido: dificulta las pruebas y oculta dependencias.",
     "<code>TRUE</code>, <code>FALSE</code> y <code>NULL</code> del evaluador (Semana 6) y el mapa <code>builtins</code> (Semana 8).",
     "Variables de paquete o <code>sync.Once</code> para inicialización perezosa y segura entre goroutines."),
    # ---------------- ESTRUCTURALES ----------------
    ("adapter", "Adapter", "estructural",
     "Convertir la interfaz de un tipo existente en otra que el cliente espera, sin modificar ninguno de los dos.",
     "El lexer lee de un <code>string</code>, pero queremos alimentarlo desde archivos, la red o stdin (<code>io.Reader</code>).",
     "Integrar código de terceros o de la biblioteca estándar con interfaces propias.",
     "Si puedes cambiar directamente una de las dos interfaces sin romper a nadie.",
     "Un <code>Source</code> del lexer adaptado desde <code>io.Reader</code>; el REPL ya recibe <code>io.Reader</code>/<code>io.Writer</code>.",
     "Un struct pequeño que envuelve al adaptado e implementa la interfaz destino."),
    ("bridge", "Bridge", "estructural",
     "Separar una abstracción de su implementación para que ambas puedan evolucionar de forma independiente.",
     "Los reportes de error tienen variantes (corto, detallado) y destinos (terminal, HTML). Combinarlos por herencia daría 2 × 2 tipos.",
     "Cuando hay dos dimensiones de variación que se multiplican.",
     "Si solo hay una dimensión: una interfaz simple es suficiente.",
     "Un reporte de errores del parser que puede mostrarse en terminal o en la web del curso.",
     "La abstracción guarda una interfaz (el “puente”) como campo; se refina con embedding."),
    ("composite", "Composite", "estructural",
     "Componer objetos en estructuras de árbol y tratar igual a las hojas y a los compuestos.",
     "Un programa tiene sentencias que contienen expresiones que contienen otras expresiones; queremos recorrer todo de forma uniforme.",
     "Cualquier estructura jerárquica: ASTs, sistemas de archivos, menús, escenas gráficas.",
     "Si los elementos no forman realmente una jerarquía.",
     "¡El AST completo! <code>Program</code> y <code>BlockStatement</code> contienen <code>Statement</code>s; <code>InfixExpression</code> contiene <code>Expression</code>s.",
     "Una interfaz común (<code>Node</code>) implementada por hojas y por nodos que guardan slices de <code>Node</code>."),
    ("decorator", "Decorator", "estructural",
     "Agregar responsabilidades a un objeto envolviéndolo en otro que implementa la misma interfaz.",
     "Queremos trazar o cronometrar el evaluador sin tocar su código ni crear subclases.",
     "Funcionalidad opcional y combinable: logs, métricas, caché, reintentos, autorización.",
     "Si el orden de las capas genera confusión o se apilan demasiadas.",
     "El tracing del parser (<code>trace</code>/<code>untrace</code>) envuelto como decorador de <code>Eval</code>.",
     "Un struct con un campo <code>inner Interfaz</code>; el patrón de <em>middleware</em> HTTP es el mismo."),
    ("facade", "Facade", "estructural",
     "Ofrecer una interfaz simple a un subsistema complejo.",
     "Ejecutar Monkey requiere crear lexer, parser, entorno, evaluar y revisar errores: demasiados pasos para quien solo quiere un resultado.",
     "Para dar un punto de entrada sencillo a una librería o a un conjunto de paquetes.",
     "Si la fachada empieza a crecer sin control y se vuelve un “objeto dios”.",
     "Una función <code>zelda.Run(src) (string, error)</code> para incrustar el intérprete en otras aplicaciones.",
     "Un struct o función de alto nivel que coordina los paquetes internos."),
    ("flyweight", "Flyweight", "estructural",
     "Compartir objetos inmutables en lugar de crear muchas copias idénticas.",
     "Un programa que itera mil veces crea miles de <code>Integer</code> con valores pequeños repetidos.",
     "Cuando hay muchísimos objetos iguales e inmutables y la memoria importa.",
     "Si los objetos son mutables: un cambio afectaría a todos los que los comparten.",
     "<code>TRUE</code>/<code>FALSE</code>/<code>NULL</code> ya son flyweights; un caché de enteros pequeños sería el siguiente paso.",
     "Un slice o mapa de instancias precreadas y una función fábrica que las reutiliza."),
    ("proxy", "Proxy", "estructural",
     "Proveer un sustituto que controla el acceso al objeto real (protección, caché, carga perezosa, registro).",
     "Queremos impedir que se reasignen constantes, o evitar recalcular resultados costosos.",
     "Control de acceso, memoización, objetos remotos o costosos de crear.",
     "Si solo agrega una indirección sin beneficio real.",
     "Un <code>Environment</code> de solo lectura para constantes o un caché de funciones puras.",
     "Un struct que implementa la misma interfaz y delega en el objeto real."),
    # ---------------- COMPORTAMIENTO ----------------
    ("chain-of-responsibility", "Chain of Responsibility", "comportamiento",
     "Pasar una petición por una cadena de manejadores; cada uno decide si la atiende o la pasa al siguiente.",
     "El REPL debe reconocer comandos (<code>:help</code>, <code>:env</code>, <code>:quit</code>) antes de tratar la línea como código Monkey.",
     "Cuando varios objetos pueden atender una petición y el orden importa.",
     "Si siempre la atiende el mismo manejador.",
     "Comandos del REPL y, en general, cualquier cadena de <em>middlewares</em>.",
     "Una interfaz <code>Handler</code> y un struct base embebido con el puntero al siguiente."),
    ("command", "Command", "comportamiento",
     "Encapsular una petición como un objeto para guardarla, encolarla, deshacerla o rehacerla.",
     "Queremos que el REPL permita <code>:undo</code> de la última sentencia <code>let</code>.",
     "Historial, deshacer/rehacer, colas de tareas, macros.",
     "Para acciones simples que nunca se guardan ni se deshacen.",
     "Historial del REPL; la secuela del libro compila a instrucciones (<em>opcodes</em>) que son comandos.",
     "Una interfaz <code>Execute/Undo</code> o simplemente valores de tipo <code>func()</code>."),
    ("iterator", "Iterator", "comportamiento",
     "Recorrer los elementos de una colección sin exponer su estructura interna.",
     "El parser necesita tokens uno a uno sin saber cómo el lexer los obtiene del texto.",
     "Colecciones con recorridos personalizados o perezosos.",
     "Si basta con un slice y <code>for range</code>.",
     "¡El lexer! <code>NextToken()</code> es un iterador; con Go 1.23 puede exponerse como <code>iter.Seq[Token]</code>.",
     "Método <code>Next()</code> o funciones <code>iter.Seq</code> que funcionan con <code>for range</code>."),
    ("mediator", "Mediator", "comportamiento",
     "Centralizar la comunicación entre componentes para que no dependan unos de otros directamente.",
     "Un editor, una consola y un panel de errores deben coordinarse; si se llaman entre sí, todo queda acoplado.",
     "Interfaces gráficas, salas de chat, sistemas con muchas interacciones cruzadas.",
     "Si el mediador termina concentrando toda la lógica del sistema.",
     "Una “mini IDE” web para Monkey que coordina editor, salida y errores (como el simulador del curso).",
     "Un struct que conoce a todos los componentes y un método <code>Notify</code>."),
    ("memento", "Memento", "comportamiento",
     "Capturar y restaurar el estado interno de un objeto sin exponer sus detalles.",
     "Si un programa falla a mitad de ejecución, el REPL debería volver al entorno anterior.",
     "Deshacer, puntos de control (checkpoints) y transacciones.",
     "Si el estado es enorme y copiarlo es demasiado costoso.",
     "Guardar el <code>Environment</code> antes de evaluar y restaurarlo si aparece un <code>*object.Error</code>.",
     "Un tipo con campos privados que solo el dueño sabe crear y leer; <code>maps.Clone</code> ayuda."),
    ("observer", "Observer", "comportamiento",
     "Notificar automáticamente a varios suscriptores cuando ocurre un evento.",
     "Queremos un tracer, un profiler y un depurador escuchando lo que hace el evaluador sin modificarlo cada vez.",
     "Eventos, notificaciones, interfaces reactivas.",
     "Si el orden de las notificaciones importa o hay demasiados observadores encadenados.",
     "Eventos <code>call</code>/<code>error</code> del evaluador para depurar programas Monkey.",
     "Un slice de interfaces <code>Observer</code>, o canales cuando los suscriptores son goroutines."),
    ("state", "State", "comportamiento",
     "Permitir que un objeto cambie su comportamiento cuando cambia su estado interno.",
     "Un escáner se comporta distinto dentro de un string, de un comentario o en modo normal.",
     "Máquinas de estados: lexers, protocolos, flujos de interfaz.",
     "Si solo hay dos estados y un <code>if</code> es suficiente.",
     "Un lexer con modos (string, comentario). El REPL podría tener modo normal y modo multilínea.",
     "Una interfaz <code>State</code> o funciones de estado <code>type stateFn func(*Lexer) stateFn</code>."),
    ("strategy", "Strategy", "comportamiento",
     "Definir una familia de algoritmos intercambiables detrás de una misma interfaz.",
     "Cada operador o tipo de token necesita su propio algoritmo, y queremos agregar nuevos sin tocar el núcleo.",
     "Algoritmos que varían según el contexto y se eligen en tiempo de ejecución.",
     "Si hay un solo algoritmo y no se espera que cambie.",
     "¡El Pratt parser! <code>prefixParseFns</code> e <code>infixParseFns</code> son mapas de estrategias por token (Semana 5).",
     "Tipos función guardados en un <code>map</code>; o una interfaz si la estrategia tiene estado."),
    ("template-method", "Template Method", "comportamiento",
     "Fijar el esqueleto de un algoritmo y dejar que ciertos pasos los implementen otros.",
     "Todo intérprete sigue tokenizar → validar → ejecutar, pero cada lenguaje lo hace distinto.",
     "Algoritmos con pasos fijos y puntos de extensión bien definidos.",
     "Si los pasos cambian de orden según el caso.",
     "El pipeline lexer → parser → evaluador; también el patrón de pruebas <em>table-driven</em> (mismo esqueleto, distintos casos).",
     "Sin clases abstractas: una función que recibe una interfaz con los pasos."),
    ("visitor", "Visitor", "comportamiento",
     "Agregar operaciones nuevas a una jerarquía de tipos sin modificar esos tipos.",
     "Sobre el mismo AST queremos evaluar, imprimir, contar nodos, verificar tipos… sin llenar cada nodo de métodos.",
     "Estructuras estables (como un AST) con muchas operaciones distintas.",
     "Si se agregan tipos de nodo con frecuencia: cada visitor debe actualizarse.",
     "Alternativa al <code>switch node := node.(type)</code> de <code>Eval</code>; compiladores reales usan visitors para varias pasadas.",
     "Métodos <code>Accept(v Visitor)</code> en cada nodo; en Go el <em>type switch</em> suele ser más simple."),
]

CATS = {"creacional": "Creacional", "estructural": "Estructural", "comportamiento": "Comportamiento"}
PAGES = {
    "go5": ("go-05-patrones-creacionales-estructurales.html", ("creacional", "estructural")),
    "go6": ("go-06-patrones-comportamiento.html", ("comportamiento",)),
}


def page_of(cat):
    return next(k for k, (_, cats) in PAGES.items() if cat in cats)


def slide(i, p):
    slug, name, cat, intent, problem, use, avoid, zelda, gonote = p
    return f'''
        <section id="p-{slug}" class="section slide" data-pattern="{slug}">
            <h2>{i}. {name}</h2>
            <div class="pat-meta">
                <span class="meta-pill">{CATS[cat]}</span>
                <span class="meta-pill"><span data-icon="code"></span><code>ejemplos/patrones/{slug}</code></span>
                <a class="meta-pill" href="{REF}" target="_blank" rel="noopener"><span data-icon="link"></span>Catálogo de referencia</a>
            </div>
            <div class="grid-2">
                <div>
                    <div class="card card-primary">
                        <h3 style="color: var(--primary);">Intención</h3>
                        <p>{intent}</p>
                        <h4>Problema en ZeldaInterpreter</h4>
                        <p>{problem}</p>
                    </div>
                    <div class="grid-2">
                        <div class="card card-success"><h3 style="color: var(--success);">Úsalo cuando…</h3><p>{use}</p></div>
                        <div class="card card-warning"><h3 style="color: var(--warning);">Evítalo si…</h3><p>{avoid}</p></div>
                    </div>
                    <div class="info-box"><strong>En Go:</strong> {gonote}</div>
                    <div class="tip-box"><strong>En el intérprete:</strong> {zelda}</div>
                </div>
                <div>
                    <div class="code-block">
                        <div class="code-header"><span>ejemplos/patrones/{slug}/main.go</span></div>
<pre><code class="language-go" data-src="ejemplos/patrones/{slug}/main.go"></code></pre>
                    </div>
                    <div class="code-block">
                        <div class="code-header"><span>$ go run ./patrones/{slug}</span></div>
<pre><code class="language-plaintext" data-out="ejemplos/patrones/{slug}"></code></pre>
                    </div>
                </div>
            </div>
        </section>
'''


HEAD = '''<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{title} | ZeldaInterpreter · Lenguajes y Compiladores 2027-1</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800&family=JetBrains+Mono:wght@400;700&display=swap" rel="stylesheet">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/styles/atom-one-dark.min.css">
    <link rel="stylesheet" href="../assets/styles.css">
</head>
<body data-week="{key}" data-week-num="2">
    <a class="skip-link" href="#contenido">Saltar al contenido</a>

    <aside class="sidebar">
        <div class="logo">
            <a href="../index.html">
                <h1 class="brand"><span class="brand-logo" data-logo="56"></span><span class="brand-name">ZeldaInterpreter</span></h1>
                <p class="logo-subtitle">{subtitle}</p>
            </a>
        </div>

        <a href="../index.html" class="back-link"><span data-icon="home"></span>Portal del Curso</a>

        <button type="button" class="present-btn" id="toggle-slides-btn" aria-pressed="false">
            <span data-icon="presentation"></span>Modo Presentación <kbd>P</kbd>
        </button>
        <a class="pdf-btn" data-pdf-link href="#"><span data-icon="file-down"></span>Ejercicios del módulo (PDF)</a>

        <nav>
            <div class="nav-group-title">Contenido de la Clase</div>
            <ul>
{nav}
            </ul>
            <div data-course-nav></div>
        </nav>

        <div class="sidebar-footer">Bloque I ampliado · Módulo {footer}</div>
    </aside>

    <main class="content" id="contenido">
'''

TAIL = '''
        <nav class="class-footer-nav" data-course-footer aria-label="Navegación entre clases"></nav>
    </main>

    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/highlight.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/languages/go.min.js"></script>
    <script src="../assets/course.js"></script>
    <script src="../assets/presentation.js"></script>
    <script src="../assets/monkey.js"></script>
    <script src="../assets/simulator.js"></script>
    <script src="../assets/patrones.js"></script>
</body>
</html>
'''


def build_page(key):
    fname, cats = PAGES[key]
    pats = [p for p in P if p[2] in cats]
    first = key == "go5"
    title = "Go V: Patrones Creacionales y Estructurales" if first else "Go VI: Patrones de Comportamiento"
    footer = "Go V" if first else "Go VI"
    sub = "Go V · Patrones creacionales y estructurales" if first else "Go VI · Patrones de comportamiento"
    n_intro = 3
    nav = ['<li><a href="#portada" class="active">Portada y objetivos</a></li>',
           '<li><a href="#que-es">1. ¿Qué es un patrón?</a></li>',
           '<li><a href="#explorador">2. Explorador de patrones</a></li>']
    for i, p in enumerate(pats, start=n_intro):
        nav.append(f'<li><a href="#p-{p[0]}">{i}. {p[1]}</a></li>')
    last = n_intro + len(pats)
    nav += [f'<li><a href="#comparativa">{last}. Comparativa</a></li>',
            '<li><a href="#simulador-zelda">Simulador ZeldaInterpreter</a></li>',
            f'<li><a href="#resumen">{last + 1}. Ideas clave</a></li>',
            f'<li><a href="#taller">{last + 2}. Taller</a></li>']
    navh = "\n".join("                " + x for x in nav)

    if first:
        hero = f'''
        <header id="portada" class="hero slide">
            <div class="hero-badge accent"><span data-icon="layers"></span>Go V · Bloque I ampliado</div>
            <h1>Patrones de Diseño I: <span class="gradient-text">Creacionales y Estructurales</span></h1>
            <p>Los patrones son soluciones probadas a problemas de diseño que se repiten. Veremos los 12 patrones creacionales y estructurales del catálogo clásico (GoF), cada uno con un programa Go ejecutable y su lugar dentro de ZeldaInterpreter.</p>
            <div class="hero-meta">
                <span class="meta-pill"><span data-icon="code"></span>12 programas en <code>ejemplos/patrones/</code></span>
                <a class="meta-pill" href="{REF}" target="_blank" rel="noopener"><span data-icon="link"></span>Referencia: catálogo de Refactoring.Guru</a>
            </div>
            <div class="objectives" aria-label="Objetivos de la clase">
                <div class="objective"><span class="num">1</span><span>Reconocer el problema que resuelve cada patrón.</span></div>
                <div class="objective"><span class="num">2</span><span>Implementarlo de forma idiomática en Go (interfaces, funciones, embedding).</span></div>
                <div class="objective"><span class="num">3</span><span>Identificar dónde ya aparece en el intérprete.</span></div>
                <div class="objective"><span class="num">4</span><span>Saber cuándo <strong>no</strong> usarlo.</span></div>
            </div>
        </header>
'''
        intro = f'''
        <section id="que-es" class="section slide">
            <h2>1. ¿Qué es un patrón de diseño?</h2>
            <p>Un patrón describe un <strong>problema recurrente</strong>, la <strong>idea</strong> de su solución y sus <strong>consecuencias</strong>. No es código para copiar: es un vocabulario común para diseñar. El catálogo clásico viene del libro <em>Design Patterns</em> (1994) de Gamma, Helm, Johnson y Vlissides, la “Gang of Four”.</p>
            <div class="grid-3">
                <div class="card card-primary"><h3 style="color: var(--primary);">Creacionales (5)</h3><p>Cómo se <strong>crean</strong> los objetos sin acoplar al cliente con tipos concretos.</p></div>
                <div class="card card-accent"><h3 style="color: var(--accent);">Estructurales (7)</h3><p>Cómo se <strong>componen</strong> objetos y tipos para formar estructuras mayores.</p></div>
                <div class="card card-success"><h3 style="color: var(--success);">De comportamiento (10)</h3><p>Cómo se <strong>reparten responsabilidades</strong> y se comunican los objetos (en Go VI).</p></div>
            </div>
            <div class="info-box">
                <strong>Patrones en Go:</strong> muchos patrones se simplifican porque Go tiene funciones de primera clase, interfaces implícitas, embedding y canales. Un Strategy puede ser un <code>map[string]func(...)</code>; un Singleton, una variable de paquete con <code>sync.Once</code>. Para profundizar, consulta el <a href="{REF}" target="_blank" rel="noopener">catálogo de Refactoring.Guru</a>, que tiene una explicación ilustrada de cada patrón y ejemplos en Go.
            </div>
        </section>
'''
    else:
        hero = f'''
        <header id="portada" class="hero slide">
            <div class="hero-badge accent"><span data-icon="layers"></span>Go VI · Bloque I ampliado</div>
            <h1>Patrones de Diseño II: <span class="gradient-text">De Comportamiento</span></h1>
            <p>Los 10 patrones de comportamiento organizan <strong>cómo colaboran</strong> los objetos: quién atiende una petición, cómo se notifican eventos, cómo se recorren estructuras y cómo cambian los algoritmos. Varios de ellos ya están escondidos dentro del intérprete del libro.</p>
            <div class="hero-meta">
                <span class="meta-pill"><span data-icon="code"></span>10 programas en <code>ejemplos/patrones/</code></span>
                <a class="meta-pill" href="{REF}" target="_blank" rel="noopener"><span data-icon="link"></span>Referencia: catálogo de Refactoring.Guru</a>
            </div>
            <div class="objectives" aria-label="Objetivos de la clase">
                <div class="objective"><span class="num">1</span><span>Reconocer Strategy, Iterator y Visitor dentro del intérprete.</span></div>
                <div class="objective"><span class="num">2</span><span>Implementar comandos, eventos y estados de forma idiomática en Go.</span></div>
                <div class="objective"><span class="num">3</span><span>Comparar el Visitor con el <em>type switch</em> del evaluador.</span></div>
                <div class="objective"><span class="num">4</span><span>Elegir el patrón adecuado y justificarlo.</span></div>
            </div>
        </header>
'''
        intro = f'''
        <section id="que-es" class="section slide">
            <h2>1. Patrones de comportamiento: repaso</h2>
            <p>En Go V vimos cómo <strong>crear</strong> y <strong>componer</strong> objetos. Los patrones de comportamiento tratan sobre <strong>algoritmos y responsabilidades</strong>. Tres de ellos ya los usamos sin nombrarlos:</p>
            <div class="grid-3">
                <div class="card card-primary"><h3 style="color: var(--primary);">Iterator</h3><p><code>lexer.NextToken()</code> entrega un token a la vez sin exponer el string interno.</p></div>
                <div class="card card-accent"><h3 style="color: var(--accent);">Strategy</h3><p><code>prefixParseFns</code>/<code>infixParseFns</code>: una función de parsing distinta por token.</p></div>
                <div class="card card-success"><h3 style="color: var(--success);">Visitor (alternativa)</h3><p><code>Eval</code> recorre el AST con un <em>type switch</em>; un Visitor haría lo mismo con métodos <code>Accept</code>.</p></div>
            </div>
            <div class="info-box"><strong>Referencia:</strong> el <a href="{REF}" target="_blank" rel="noopener">catálogo de Refactoring.Guru</a> explica cada patrón con diagramas y ejemplos en Go. Aquí los vemos aplicados al mundo de los intérpretes.</div>
        </section>
'''
    explorer = '''
        <section id="explorador" class="section slide">
            <h2>2. Explorador de patrones</h2>
            <p>Filtra por categoría o busca por palabra clave (por ejemplo <em>AST</em>, <em>REPL</em> o <em>lexer</em>). Los 22 patrones del curso están aquí; los de la otra clase abren en su página.</p>
            <div data-pattern-explorer></div>
        </section>
'''
    rows = "".join(
        f"<tr><td><a href=\"#p-{p[0]}\">{p[1]}</a></td><td>{p[3]}</td><td>{p[7]}</td></tr>" for p in pats)
    compar = f'''
        <section id="comparativa" class="section slide">
            <h2>{last}. Comparativa rápida</h2>
            <div class="table-wrap">
                <table>
                    <thead><tr><th>Patrón</th><th>Intención</th><th>En ZeldaInterpreter</th></tr></thead>
                    <tbody>{rows}</tbody>
                </table>
            </div>
            <div class="warn-box"><strong>Cuidado con la “patronitis”:</strong> un patrón agrega indirección. Úsalo cuando el problema que resuelve ya existe, no por si acaso. En Go, la solución más simple (una función, un <code>map</code>, una interfaz de un método) suele ser la correcta.</div>
        </section>
'''
    sim = ('''
        <section id="simulador-zelda" class="section slide">
            <h2>Simulador ZeldaInterpreter</h2>
            <p>Abre la pestaña <strong>AST</strong>: el árbol que ves es un <strong>Composite</strong> (bloques que contienen sentencias que contienen expresiones), y su texto se arma con un <strong>Builder</strong> (<code>bytes.Buffer</code>).</p>
            <div data-simulator data-stage="ast" data-code="let area = fn(b, h) { b * h / 2 };&#10;area(4, 3) + 1;"></div>
        </section>
''' if first else '''
        <section id="simulador-zelda" class="section slide">
            <h2>Simulador ZeldaInterpreter</h2>
            <p>La pestaña <strong>Tokens</strong> es un <strong>Iterator</strong> en acción; la pestaña <strong>AST</strong> muestra el resultado de las <strong>Strategies</strong> del Pratt parser; y la pestaña <strong>Evaluación</strong> recorre el árbol como lo haría un <strong>Visitor</strong>.</p>
            <div data-simulator data-stage="eval" data-code="let compose = fn(f, g) { fn(x) { f(g(x)) } };&#10;let inc = fn(x) { x + 1 };&#10;let dbl = fn(x) { x * 2 };&#10;compose(inc, dbl)(10);"></div>
        </section>
''')
    ideas = ("".join(f'<div class="key-idea"><h3>{p[1]}</h3><p>{p[3]}</p></div>' for p in pats[:6]))
    resumen = f'''
        <section id="resumen" class="section slide">
            <h2>{last + 1}. Ideas clave</h2>
            <div class="key-ideas">{ideas}</div>
            <div class="tip-box"><strong>Regla de oro en Go:</strong> empieza por lo más simple (funciones e interfaces pequeñas) y reconoce el patrón cuando aparezca; no al revés.</div>
        </section>
'''
    if first:
        taller_items = '''
                    <li>Ejecuta los 12 ejemplos: <code>cd ejemplos &amp;&amp; for d in patrones/*/; do go run ./$d; done</code>.</li>
                    <li>Implementa un <strong>Facade</strong> <code>zelda.Run(src string) (string, error)</code> sobre tu intérprete real.</li>
                    <li>Agrega un <strong>Decorator</strong> de traza a tu función <code>Eval</code> sin modificarla.</li>
                    <li>Implementa el <strong>Flyweight</strong> de enteros pequeños en tu paquete <code>object</code> y mide cuántas asignaciones ahorra con <code>testing.B</code>.</li>'''
    else:
        taller_items = '''
                    <li>Ejecuta los 10 ejemplos y explica qué problema resuelve cada uno.</li>
                    <li>Agrega a tu REPL una <strong>Chain of Responsibility</strong> con <code>:help</code>, <code>:env</code> y <code>:quit</code>.</li>
                    <li>Implementa <code>:undo</code> con <strong>Command</strong> o con <strong>Memento</strong> y justifica cuál elegiste.</li>
                    <li>Escribe un <strong>Visitor</strong> que cuente los nodos de tu AST real y compáralo con un <em>type switch</em>.</li>'''
    taller = f'''
        <section id="taller" class="section slide">
            <h2>{last + 2}. Taller — {footer}</h2>
            <div class="card glass">
                <h3><span data-icon="wrench"></span>Patrones aplicados a tu intérprete</h3>
                <ol class="steps">{taller_items}
                </ol>
            </div>
            <div class="exercise-cta">
                <div>
                    <h3>Guía de ejercicios · {footer}</h3>
                    <p>Identificar patrones, elegir el adecuado para un problema y aplicarlos en Go.</p>
                </div>
                <a class="btn-download" data-pdf-link href="#"><span data-icon="file-down"></span>Descargar PDF</a>
            </div>
        </section>
'''
    slides = "".join(slide(i, p) for i, p in enumerate(pats, start=n_intro))
    html = (HEAD.format(title=title, key=key, subtitle=sub, nav=navh, footer=footer)
            + hero + intro + explorer + slides + compar + sim + resumen + taller + TAIL)
    (ROOT / "clases" / fname).write_text(html, encoding="utf-8")
    print("generado:", fname, f"({len(pats)} patrones)")


def build_js():
    data = [{"slug": p[0], "name": p[1], "cat": p[2], "intent": p[3], "zelda": p[7],
             "page": PAGES[page_of(p[2])][0]} for p in P]
    js = """/**
 * Explorador de patrones de diseño (generado por scripts/build_patrones.py).
 * Se monta en <div data-pattern-explorer></div>.
 */
(function () {
    'use strict';
    const PATTERNS = %s;
    const CATS = { todos: 'Todos', creacional: 'Creacionales', estructural: 'Estructurales', comportamiento: 'De comportamiento' };

    function mount(root) {
        const here = location.pathname.split('/').pop();
        let cat = 'todos', q = '';
        root.innerHTML = `
            <div class="pat-filters" role="group" aria-label="Filtrar por categoría">
                ${Object.entries(CATS).map(([k, v]) => `<button type="button" class="btn-pill" data-cat="${k}" aria-pressed="${k === 'todos'}">${v}</button>`).join('')}
                <label class="sr-only" for="pat-q">Buscar patrón</label>
                <input id="pat-q" class="pat-search" type="search" placeholder="Buscar: AST, REPL, lexer, caché…" autocomplete="off">
            </div>
            <div class="viz-label" id="pat-count" aria-live="polite"></div>
            <div class="pat-grid"></div>`;
        const grid = root.querySelector('.pat-grid');
        const count = root.querySelector('#pat-count');
        const strip = (s) => s.replace(/<[^>]+>/g, '');
        const norm = (s) => strip(s).normalize('NFD').replace(/[\\u0300-\\u036f]/g, '').toLowerCase();

        function render() {
            const list = PATTERNS.filter((p) => (cat === 'todos' || p.cat === cat) &&
                (!q || norm(p.name + ' ' + p.intent + ' ' + p.zelda).includes(norm(q))));
            count.textContent = list.length + ' patrón(es)';
            grid.innerHTML = list.length ? list.map((p) => {
                const local = p.page === here;
                const href = (local ? '' : p.page) + '#p-' + p.slug;
                return `<a class="pat-card cat-${p.cat}" href="${href}">
                    <span class="pat-cat">${CATS[p.cat]}</span>
                    <b>${p.name}</b>
                    <p>${p.intent}</p>
                    <span class="pat-zelda">${strip(p.zelda)}</span>
                    ${local ? '' : '<span class="pat-other">Se ve en la otra clase de patrones →</span>'}
                </a>`;
            }).join('') : '<p class="pat-empty">Ningún patrón coincide con la búsqueda.</p>';
        }
        root.querySelectorAll('[data-cat]').forEach((b) => b.addEventListener('click', () => {
            cat = b.dataset.cat;
            root.querySelectorAll('[data-cat]').forEach((x) => x.setAttribute('aria-pressed', String(x === b)));
            render();
        }));
        const input = root.querySelector('#pat-q');
        input.addEventListener('input', () => { q = input.value.trim(); render(); });
        input.addEventListener('keydown', (e) => e.stopPropagation()); // no cambiar de diapositiva al escribir
        render();
    }
    document.addEventListener('DOMContentLoaded', () => document.querySelectorAll('[data-pattern-explorer]').forEach(mount));
})();
""" % json.dumps(data, ensure_ascii=False, indent=8)
    (ROOT / "assets" / "patrones.js").write_text(js, encoding="utf-8")
    print("generado: assets/patrones.js")


if __name__ == "__main__":
    build_page("go5")
    build_page("go6")
    build_js()

/**
 * Explorador de patrones de diseño (generado por scripts/build_patrones.py).
 * Se monta en <div data-pattern-explorer></div>.
 */
(function () {
    'use strict';
    const PATTERNS = [
        {
                "slug": "factory-method",
                "name": "Factory Method",
                "cat": "creacional",
                "intent": "Delegar en una función o método la decisión de qué tipo concreto crear, devolviendo una abstracción.",
                "zelda": "<code>lexer.New</code>, <code>parser.New</code> y <code>object.NewEnvironment</code> son fábricas simples; <code>NewPrinter(format)</code> lo sería con varios formatos.",
                "page": "go-05-patrones-creacionales-estructurales.html"
        },
        {
                "slug": "abstract-factory",
                "name": "Abstract Factory",
                "cat": "creacional",
                "intent": "Crear familias completas de objetos relacionados que están diseñados para usarse juntos.",
                "zelda": "Un <code>Theme</code> del REPL que produce <code>Prompt</code> y <code>ErrorFormatter</code> coherentes entre sí.",
                "page": "go-05-patrones-creacionales-estructurales.html"
        },
        {
                "slug": "builder",
                "name": "Builder",
                "cat": "creacional",
                "intent": "Construir un objeto complejo paso a paso, separando la construcción de su representación final.",
                "zelda": "Un <code>ProgramBuilder</code> para tests y <code>strings.Builder</code>/<code>bytes.Buffer</code> en todos los <code>String()</code> del AST.",
                "page": "go-05-patrones-creacionales-estructurales.html"
        },
        {
                "slug": "prototype",
                "name": "Prototype",
                "cat": "creacional",
                "intent": "Crear nuevos objetos clonando uno existente en lugar de construirlos desde cero.",
                "zelda": "<code>push</code> y <code>rest</code> del Cap. 4 crean arrays nuevos; un método <code>Clone()</code> generaliza la idea.",
                "page": "go-05-patrones-creacionales-estructurales.html"
        },
        {
                "slug": "singleton",
                "name": "Singleton",
                "cat": "creacional",
                "intent": "Garantizar que exista una única instancia de un tipo y ofrecer un punto de acceso global a ella.",
                "zelda": "<code>TRUE</code>, <code>FALSE</code> y <code>NULL</code> del evaluador (Semana 6) y el mapa <code>builtins</code> (Semana 8).",
                "page": "go-05-patrones-creacionales-estructurales.html"
        },
        {
                "slug": "adapter",
                "name": "Adapter",
                "cat": "estructural",
                "intent": "Convertir la interfaz de un tipo existente en otra que el cliente espera, sin modificar ninguno de los dos.",
                "zelda": "Un <code>Source</code> del lexer adaptado desde <code>io.Reader</code>; el REPL ya recibe <code>io.Reader</code>/<code>io.Writer</code>.",
                "page": "go-05-patrones-creacionales-estructurales.html"
        },
        {
                "slug": "bridge",
                "name": "Bridge",
                "cat": "estructural",
                "intent": "Separar una abstracción de su implementación para que ambas puedan evolucionar de forma independiente.",
                "zelda": "Un reporte de errores del parser que puede mostrarse en terminal o en la web del curso.",
                "page": "go-05-patrones-creacionales-estructurales.html"
        },
        {
                "slug": "composite",
                "name": "Composite",
                "cat": "estructural",
                "intent": "Componer objetos en estructuras de árbol y tratar igual a las hojas y a los compuestos.",
                "zelda": "¡El AST completo! <code>Program</code> y <code>BlockStatement</code> contienen <code>Statement</code>s; <code>InfixExpression</code> contiene <code>Expression</code>s.",
                "page": "go-05-patrones-creacionales-estructurales.html"
        },
        {
                "slug": "decorator",
                "name": "Decorator",
                "cat": "estructural",
                "intent": "Agregar responsabilidades a un objeto envolviéndolo en otro que implementa la misma interfaz.",
                "zelda": "El tracing del parser (<code>trace</code>/<code>untrace</code>) envuelto como decorador de <code>Eval</code>.",
                "page": "go-05-patrones-creacionales-estructurales.html"
        },
        {
                "slug": "facade",
                "name": "Facade",
                "cat": "estructural",
                "intent": "Ofrecer una interfaz simple a un subsistema complejo.",
                "zelda": "Una función <code>zelda.Run(src) (string, error)</code> para incrustar el intérprete en otras aplicaciones.",
                "page": "go-05-patrones-creacionales-estructurales.html"
        },
        {
                "slug": "flyweight",
                "name": "Flyweight",
                "cat": "estructural",
                "intent": "Compartir objetos inmutables en lugar de crear muchas copias idénticas.",
                "zelda": "<code>TRUE</code>/<code>FALSE</code>/<code>NULL</code> ya son flyweights; un caché de enteros pequeños sería el siguiente paso.",
                "page": "go-05-patrones-creacionales-estructurales.html"
        },
        {
                "slug": "proxy",
                "name": "Proxy",
                "cat": "estructural",
                "intent": "Proveer un sustituto que controla el acceso al objeto real (protección, caché, carga perezosa, registro).",
                "zelda": "Un <code>Environment</code> de solo lectura para constantes o un caché de funciones puras.",
                "page": "go-05-patrones-creacionales-estructurales.html"
        },
        {
                "slug": "chain-of-responsibility",
                "name": "Chain of Responsibility",
                "cat": "comportamiento",
                "intent": "Pasar una petición por una cadena de manejadores; cada uno decide si la atiende o la pasa al siguiente.",
                "zelda": "Comandos del REPL y, en general, cualquier cadena de <em>middlewares</em>.",
                "page": "go-06-patrones-comportamiento.html"
        },
        {
                "slug": "command",
                "name": "Command",
                "cat": "comportamiento",
                "intent": "Encapsular una petición como un objeto para guardarla, encolarla, deshacerla o rehacerla.",
                "zelda": "Historial del REPL; la secuela del libro compila a instrucciones (<em>opcodes</em>) que son comandos.",
                "page": "go-06-patrones-comportamiento.html"
        },
        {
                "slug": "iterator",
                "name": "Iterator",
                "cat": "comportamiento",
                "intent": "Recorrer los elementos de una colección sin exponer su estructura interna.",
                "zelda": "¡El lexer! <code>NextToken()</code> es un iterador; con Go 1.23 puede exponerse como <code>iter.Seq[Token]</code>.",
                "page": "go-06-patrones-comportamiento.html"
        },
        {
                "slug": "mediator",
                "name": "Mediator",
                "cat": "comportamiento",
                "intent": "Centralizar la comunicación entre componentes para que no dependan unos de otros directamente.",
                "zelda": "Una “mini IDE” web para Monkey que coordina editor, salida y errores (como el simulador del curso).",
                "page": "go-06-patrones-comportamiento.html"
        },
        {
                "slug": "memento",
                "name": "Memento",
                "cat": "comportamiento",
                "intent": "Capturar y restaurar el estado interno de un objeto sin exponer sus detalles.",
                "zelda": "Guardar el <code>Environment</code> antes de evaluar y restaurarlo si aparece un <code>*object.Error</code>.",
                "page": "go-06-patrones-comportamiento.html"
        },
        {
                "slug": "observer",
                "name": "Observer",
                "cat": "comportamiento",
                "intent": "Notificar automáticamente a varios suscriptores cuando ocurre un evento.",
                "zelda": "Eventos <code>call</code>/<code>error</code> del evaluador para depurar programas Monkey.",
                "page": "go-06-patrones-comportamiento.html"
        },
        {
                "slug": "state",
                "name": "State",
                "cat": "comportamiento",
                "intent": "Permitir que un objeto cambie su comportamiento cuando cambia su estado interno.",
                "zelda": "Un lexer con modos (string, comentario). El REPL podría tener modo normal y modo multilínea.",
                "page": "go-06-patrones-comportamiento.html"
        },
        {
                "slug": "strategy",
                "name": "Strategy",
                "cat": "comportamiento",
                "intent": "Definir una familia de algoritmos intercambiables detrás de una misma interfaz.",
                "zelda": "¡El Pratt parser! <code>prefixParseFns</code> e <code>infixParseFns</code> son mapas de estrategias por token (Semana 5).",
                "page": "go-06-patrones-comportamiento.html"
        },
        {
                "slug": "template-method",
                "name": "Template Method",
                "cat": "comportamiento",
                "intent": "Fijar el esqueleto de un algoritmo y dejar que ciertos pasos los implementen otros.",
                "zelda": "El pipeline lexer → parser → evaluador; también el patrón de pruebas <em>table-driven</em> (mismo esqueleto, distintos casos).",
                "page": "go-06-patrones-comportamiento.html"
        },
        {
                "slug": "visitor",
                "name": "Visitor",
                "cat": "comportamiento",
                "intent": "Agregar operaciones nuevas a una jerarquía de tipos sin modificar esos tipos.",
                "zelda": "Alternativa al <code>switch node := node.(type)</code> de <code>Eval</code>; compiladores reales usan visitors para varias pasadas.",
                "page": "go-06-patrones-comportamiento.html"
        }
];
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
        const norm = (s) => strip(s).normalize('NFD').replace(/[\u0300-\u036f]/g, '').toLowerCase();

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

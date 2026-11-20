/**
 * Fuente única de verdad del curso: semanas, rutas, PDFs de ejercicios e iconos.
 * Todas las páginas leen este archivo para construir la navegación entre clases,
 * así que agregar o renombrar una clase se hace SOLO aquí.
 * Lenguajes y Compiladores 2027-1 — Intérprete Monkey en Go
 */
(function () {
    'use strict';

    const WEEKS = [
        {
            n: 1, key: '1', num: '01', footer: 'Semana 1', block: 'go', tag: 'SEMANA 01 · GO',
            file: 'semana-01-go-fundamentos.html',
            pdf: 'semana-01-ejercicios.pdf',
            short: 'Fundamentos de Go',
            title: 'Fundamentos de Go para Compiladores',
            ref: 'Intro del libro + Base Go',
            desc: 'Por qué Go, módulos, tipos básicos (<code>byte</code>, <code>int64</code>, <code>string</code>), control de flujo, funciones, slices y maps.'
        },
        {
            n: 2, key: '2', num: '02', footer: 'Semana 2', block: 'go', tag: 'SEMANA 02 · GO',
            file: 'semana-02-go-intermedio.html',
            pdf: 'semana-02-ejercicios.pdf',
            short: 'Go para Intérpretes',
            title: 'Go Intermedio: Structs, Interfaces y Testing',
            ref: 'Patrones clave del libro',
            desc: 'Punteros y receivers, interfaces implícitas, type switches, funciones como valores (closures) y <em>Table-Driven Tests</em>.'
        },
        {
            n: 2, key: 'go3', num: 'G3', footer: 'Go III', block: 'go', tag: 'GO III · POO',
            file: 'go-03-poo-composicion.html',
            pdf: 'go-03-ejercicios.pdf',
            short: 'POO en Go',
            title: 'Objetos, Agregación, Composición y “Herencia” en Go',
            ref: 'Structs, métodos, embedding, interfaces',
            desc: 'Encapsulación con paquetes, constructores, <strong>agregación vs composición</strong>, embedding en lugar de herencia y polimorfismo con interfaces.'
        },
        {
            n: 2, key: 'go4', num: 'G4', footer: 'Go IV', block: 'go', tag: 'GO IV · CONCURRENCIA',
            file: 'go-04-concurrencia.html',
            pdf: 'go-04-ejercicios.pdf',
            short: 'Concurrencia',
            title: 'Concurrencia: Goroutines, Canales y Sincronización',
            ref: 'goroutines, channels, select, sync, context',
            desc: 'Goroutines, canales con y sin buffer, <code>select</code>, <code>WaitGroup</code>, <code>Mutex</code>, <code>context</code>, pipelines y <em>worker pools</em>.'
        },
        {
            n: 2, key: 'go5', num: 'G5', footer: 'Go V', block: 'go', tag: 'GO V · PATRONES I',
            file: 'go-05-patrones-creacionales-estructurales.html',
            pdf: 'go-05-ejercicios.pdf',
            short: 'Patrones I',
            title: 'Patrones de Diseño I: Creacionales y Estructurales',
            ref: 'Catálogo GoF · 12 patrones',
            desc: 'Factory Method, Abstract Factory, Builder, Prototype, Singleton, Adapter, Bridge, Composite, Decorator, Facade, Flyweight y Proxy en Go.'
        },
        {
            n: 2, key: 'go6', num: 'G6', footer: 'Go VI', block: 'go', tag: 'GO VI · PATRONES II',
            file: 'go-06-patrones-comportamiento.html',
            pdf: 'go-06-ejercicios.pdf',
            short: 'Patrones II',
            title: 'Patrones de Diseño II: De Comportamiento',
            ref: 'Catálogo GoF · 10 patrones',
            desc: 'Chain of Responsibility, Command, Iterator, Mediator, Memento, Observer, State, Strategy, Template Method y Visitor en Go.'
        },
        {
            n: 3, key: '3', num: '03', footer: 'Semana 3', block: 'monkey', tag: 'SEMANA 03 · CAP. 1',
            file: 'semana-03-lexing.html',
            pdf: 'semana-03-ejercicios.pdf',
            short: 'Tokens y Lexer',
            title: 'Análisis Léxico: Tokens, Lexer y REPL',
            ref: 'Cap. 1 (pp. 5–27)',
            desc: 'Definición de <code>Token</code>, arquitectura de doble puntero del <code>Lexer</code>, <code>peekChar()</code> para <code>==</code>/<code>!=</code>, pruebas y primer REPL.'
        },
        {
            n: 4, key: '4', num: '04', footer: 'Semana 4', block: 'monkey', tag: 'SEMANA 04 · CAP. 2A',
            file: 'semana-04-parsing-sentencias.html',
            pdf: 'semana-04-ejercicios.pdf',
            short: 'AST y Sentencias',
            title: 'Análisis Sintáctico I: AST y Sentencias',
            ref: 'Cap. 2.1 – 2.5 (pp. 28–45)',
            desc: '¿Por qué un parser manual? Diseño de <code>Node</code>, <code>Statement</code> y <code>Expression</code>, parsing de <code>let</code> y <code>return</code>, manejo de errores.'
        },
        {
            n: 5, key: '5', num: '05', footer: 'Semana 5', block: 'monkey', tag: 'SEMANA 05 · CAP. 2B',
            file: 'semana-05-pratt-parsing.html',
            pdf: 'semana-05-ejercicios.pdf',
            short: 'Pratt Parsing',
            title: 'Análisis Sintáctico II: Pratt Parsing',
            ref: 'Cap. 2.6 – 2.9 (pp. 46–102)',
            desc: 'Top-Down Operator Precedence paso a paso: prefijos, infijos, precedencias, booleanos, <code>if/else</code>, <code>fn</code> y llamadas.'
        },
        {
            n: 6, key: '6', num: '06', footer: 'Semana 6', block: 'monkey', tag: 'SEMANA 06 · CAP. 3A',
            file: 'semana-06-evaluacion-base.html',
            pdf: 'semana-06-ejercicios.pdf',
            short: 'Evaluación y Objetos',
            title: 'Evaluación I: Objetos, Expresiones y Control',
            ref: 'Cap. 3.1 – 3.8 (pp. 103–135)',
            desc: 'Tree-Walking, sistema de objetos (<code>Integer</code>, <code>Boolean</code>, <code>Null</code>), operadores, condicionales, <code>return</code> y errores.'
        },
        {
            n: 7, key: '7', num: '07', footer: 'Semana 7', block: 'monkey', tag: 'SEMANA 07 · CAP. 3B',
            file: 'semana-07-entornos-closures.html',
            pdf: 'semana-07-ejercicios.pdf',
            short: 'Entornos y Closures',
            title: 'Evaluación II: Entornos, Funciones y Closures',
            ref: 'Cap. 3.9 – 3.11 (pp. 136–151)',
            desc: 'Bindings con <code>Environment</code>, ámbitos encadenados (<code>outer</code>), funciones de primera clase, closures y Garbage Collection.'
        },
        {
            n: 8, key: '8', num: '08', footer: 'Semana 8', block: 'monkey', tag: 'SEMANA 08 · CAP. 4',
            file: 'semana-08-extensiones.html',
            pdf: 'semana-08-ejercicios.pdf',
            short: 'Strings, Arrays y Hashes',
            title: 'Extensiones: Strings, Built-ins, Arrays y Hashes',
            ref: 'Cap. 4 (pp. 152–201)',
            desc: 'Ciclo completo de evolución del lenguaje: concatenación, <code>len/push</code>, indexación <code>a[i]</code>, <code>HashKey</code> y gran final.'
        }
    ];

    /* Iconos Lucide (MIT) — trazos 24x24, stroke=currentColor */
    const ICONS = {
        'book': '<path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"/><path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"/>',
        'presentation': '<path d="M2 3h20"/><path d="M21 3v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V3"/><path d="m7 21 5-5 5 5"/>',
        'arrow-left': '<path d="m12 19-7-7 7-7"/><path d="M19 12H5"/>',
        'arrow-right': '<path d="M5 12h14"/><path d="m12 5 7 7-7 7"/>',
        'file-down': '<path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"/><path d="M14 2v4a2 2 0 0 0 2 2h4"/><path d="M12 18v-6"/><path d="m9 15 3 3 3-3"/>',
        'lightbulb': '<path d="M15 14c.2-1 .7-1.7 1.5-2.5 1-.9 1.5-2.2 1.5-3.5A6 6 0 0 0 6 8c0 1 .2 2.2 1.5 3.5.7.7 1.3 1.5 1.5 2.5"/><path d="M9 18h6"/><path d="M10 22h4"/>',
        'alert': '<path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3"/><path d="M12 9v4"/><path d="M12 17h.01"/>',
        'info': '<circle cx="12" cy="12" r="10"/><path d="M12 16v-4"/><path d="M12 8h.01"/>',
        'check': '<circle cx="12" cy="12" r="10"/><path d="m9 12 2 2 4-4"/>',
        'menu': '<line x1="4" x2="20" y1="12" y2="12"/><line x1="4" x2="20" y1="6" y2="6"/><line x1="4" x2="20" y1="18" y2="18"/>',
        'x': '<path d="M18 6 6 18"/><path d="m6 6 12 12"/>',
        'maximize': '<path d="M8 3H5a2 2 0 0 0-2 2v3"/><path d="M21 8V5a2 2 0 0 0-2-2h-3"/><path d="M3 16v3a2 2 0 0 0 2 2h3"/><path d="M16 21h3a2 2 0 0 0 2-2v-3"/>',
        'home': '<path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/>',
        'code': '<polyline points="16 18 22 12 16 6"/><polyline points="8 6 2 12 8 18"/>',
        'flask': '<path d="M10 2v7.527a2 2 0 0 1-.211.896L4.72 20.55a1 1 0 0 0 .9 1.45h12.76a1 1 0 0 0 .9-1.45l-5.069-10.127A2 2 0 0 1 14 9.527V2"/><path d="M8.5 2h7"/><path d="M7 16h10"/>',
        'folder': '<path d="M20 20a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.9a2 2 0 0 1-1.69-.9L9.6 3.9A2 2 0 0 0 7.93 3H4a2 2 0 0 0-2 2v13a2 2 0 0 0 2 2Z"/>',
        'clock': '<circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>',
        'target': '<circle cx="12" cy="12" r="10"/><circle cx="12" cy="12" r="6"/><circle cx="12" cy="12" r="2"/>',
        'wrench': '<path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>',
        'tree': '<rect x="16" y="16" width="6" height="6" rx="1"/><rect x="2" y="16" width="6" height="6" rx="1"/><rect x="9" y="2" width="6" height="6" rx="1"/><path d="M5 16v-3a1 1 0 0 1 1-1h12a1 1 0 0 1 1 1v3"/><path d="M12 12V8"/>',
        'play': '<polygon points="6 3 20 12 6 21 6 3"/>',
        'terminal': '<polyline points="4 17 10 11 4 5"/><line x1="12" x2="20" y1="19" y2="19"/>',
        'layers': '<path d="m12.83 2.18a2 2 0 0 0-1.66 0L2.6 6.08a1 1 0 0 0 0 1.83l8.58 3.91a2 2 0 0 0 1.66 0l8.58-3.9a1 1 0 0 0 0-1.83Z"/><path d="m22 17.65-9.17 4.16a2 2 0 0 1-1.66 0L2 17.65"/><path d="m22 12.65-9.17 4.16a2 2 0 0 1-1.66 0L2 12.65"/>',
        'zap': '<path d="M4 14a1 1 0 0 1-.78-1.63l9.9-10.2a.5.5 0 0 1 .86.46l-1.92 6.02A1 1 0 0 0 13 10h7a1 1 0 0 1 .78 1.63l-9.9 10.2a.5.5 0 0 1-.86-.46l1.92-6.02A1 1 0 0 0 11 14z"/>',
        'calendar': '<rect width="18" height="18" x="3" y="4" rx="2"/><path d="M16 2v4"/><path d="M8 2v4"/><path d="M3 10h18"/>',
        'link': '<path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/>'
    };

    /* ===================== Marca: ZeldaInterpreter ===================== */
    const BRAND = 'ZeldaInterpreter';

    /* Mascota original (16x16 pixel art): aprendiz de mago con capucha azul,
       broche rosa y bastón con gema — colores de la paleta del curso. */
    const LOGO_ROWS = [
        '.......hh.......',
        '......hHHh....P.',
        '.....hHHHHh..PPP',
        '....hHHHHHHh..P.',
        '...hHHHHHHHHh.s.',
        '...hHrrrrrrHh.s.',
        '...hrSSSSSSrh.s.',
        '...hSSESSESSh.s.',
        '....SSSSSSSS..s.',
        '.....SSmmSS...s.',
        '....cCCCCCCcSSs.',
        '...cCCCLLCCCc.s.',
        '...cCCCPPCCCc.s.',
        '...cCCLLLLCCc.s.',
        '....cCCCCCCc..s.',
        '....dd....dd..s.'
    ];
    const LOGO_COLORS = {
        h: '#075985', H: '#38bdf8', P: '#f472b6', r: '#1e293b', S: '#fdba74', E: '#0f172a',
        m: '#be185d', c: '#0c4a6e', C: '#0284c7', L: '#e0f2fe', s: '#a16207', d: '#334155'
    };
    function logoSvg(size) {
        let rects = '';
        LOGO_ROWS.forEach((row, y) => {
            [...row].forEach((c, x) => {
                if (LOGO_COLORS[c]) rects += `<rect x="${x}" y="${y}" width="1" height="1" fill="${LOGO_COLORS[c]}"/>`;
            });
        });
        return `<svg class="brand-svg" viewBox="0 0 16 16" width="${size || 56}" height="${size || 56}" shape-rendering="crispEdges" role="img" aria-label="Logo de ${BRAND}: mascota pixel art, aprendiz de mago con capucha azul"><title>${BRAND}</title>${rects}</svg>`;
    }
    function hydrateLogos(scope) {
        (scope || document).querySelectorAll('[data-logo]').forEach((el) => {
            if (el.dataset.logoDone) return;
            el.innerHTML = logoSvg(Number(el.dataset.logo) || 56);
            el.dataset.logoDone = '1';
        });
    }

    function svg(name, extraClass) {
        const body = ICONS[name];
        if (!body) return '';
        return `<svg class="icon${extraClass ? ' ' + extraClass : ''}" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" focusable="false">${body}</svg>`;
    }

    /* Detecta si estamos dentro de /clases/ o en la raíz para resolver rutas relativas */
    function paths() {
        const inClases = /\/clases\//.test(window.location.pathname);
        return {
            root: inClases ? '../' : '',
            clase: (w) => (inClases ? '' : 'clases/') + w.file,
            pdf: (w) => (inClases ? '../' : '') + 'ejercicios/pdf/' + w.pdf
        };
    }

    function currentWeek() {
        const key = String(document.body.dataset.week || '');
        return WEEKS.find((w) => w.key === key) || null;
    }

    /* Reemplaza <span data-icon="nombre"></span> por el SVG correspondiente */
    function hydrateIcons(scope) {
        (scope || document).querySelectorAll('[data-icon]').forEach((el) => {
            if (el.dataset.iconDone) return;
            el.insertAdjacentHTML('afterbegin', svg(el.dataset.icon));
            el.dataset.iconDone = '1';
        });
    }

    /* Grupo "Clases del Curso" en el sidebar */
    function renderCourseNav() {
        const mount = document.querySelector('[data-course-nav]');
        if (!mount) return;
        const p = paths();
        const cur = currentWeek();
        const items = WEEKS.map((w) => {
            const active = cur && cur.key === w.key;
            return `<li><a href="${p.clase(w)}" class="week-link${active ? ' current' : ''}"${active ? ' aria-current="page"' : ''}>
                <span class="week-num ${w.block}">${w.num}</span>
                <span>${w.short}</span></a></li>`;
        }).join('');
        mount.innerHTML = `<hr class="nav-separator"><div class="nav-group-title">Clases del Curso</div><ul class="week-list">${items}</ul>`;
    }

    /* Botón de descarga del PDF de ejercicios de la semana actual */
    function renderPdfButtons() {
        const p = paths();
        const cur = currentWeek();
        document.querySelectorAll('[data-pdf-link]').forEach((el) => {
            const key = String(el.dataset.pdfLink || (cur && cur.key) || '');
            const w = WEEKS.find((x) => x.key === key);
            if (!w) return;
            el.setAttribute('href', p.pdf(w));
            el.setAttribute('download', '');
        });
    }

    /* Navegación inferior Anterior / Siguiente */
    function renderFooterNav() {
        const mount = document.querySelector('[data-course-footer]');
        if (!mount) return;
        const p = paths();
        const cur = currentWeek();
        if (!cur) return;
        const i = WEEKS.indexOf(cur);
        const prev = WEEKS[i - 1];
        const next = WEEKS[i + 1];
        const prevHtml = prev
            ? `<a href="${p.clase(prev)}" rel="prev">${svg('arrow-left')}<span><small>${prev.footer}</small>${prev.short}</span></a>`
            : `<a href="${p.root}index.html">${svg('home')}<span><small>Inicio</small>Portal del Curso</span></a>`;
        const nextHtml = next
            ? `<a href="${p.clase(next)}" rel="next" class="next"><span><small>${next.footer}</small>${next.short}</span>${svg('arrow-right')}</a>`
            : `<a href="${p.root}index.html" class="next"><span><small>Fin del curso</small>Portal del Curso</span>${svg('home')}</a>`;
        mount.innerHTML = prevHtml + nextHtml;
    }

    /* Roadmap del portal (index.html) */
    function renderRoadmap() {
        const mount = document.querySelector('[data-course-roadmap]');
        if (!mount) return;
        const p = paths();
        mount.innerHTML = WEEKS.map((w) => `
            <article class="module-card ready">
                <div class="module-card-top">
                    <span class="week-tag${w.block === 'go' ? ' go-block' : ''}">${w.tag}</span>
                    <span class="module-ref">${w.ref}</span>
                </div>
                <h3><a class="stretched" href="${p.clase(w)}">${w.title}</a></h3>
                <p>${w.desc}</p>
                <div class="module-card-footer">
                    <span class="module-cta">${svg('presentation')} Abrir clase</span>
                    <a class="module-pdf" href="${p.pdf(w)}" download>${svg('file-down')} Ejercicios PDF</a>
                </div>
            </article>`).join('');
    }

    window.COURSE = { WEEKS, ICONS, svg, hydrateIcons, BRAND, logoSvg, hydrateLogos };

    document.addEventListener('DOMContentLoaded', () => {
        renderCourseNav();
        renderFooterNav();
        renderRoadmap();
        renderPdfButtons();
        hydrateIcons();
        hydrateLogos();
    });
})();

/**
 * Simulador ZeldaInterpreter — widget reutilizable en todas las clases.
 * Muestra las tres fases del intérprete sobre el mismo código Monkey:
 *   1. Tokens (Lexer)  ·  2. AST (Parser Pratt)  ·  3. Evaluación (Tree-walking)
 *
 * Uso en HTML:
 *   <div data-simulator data-stage="tokens|ast|eval" data-code="..."></div>
 * Requiere monkey.js y course.js (para iconos/logo). La semana se lee de <body data-week>.
 */
(function () {
    'use strict';

    const STAGES = [
        // start: semana en que se empieza a construir · done: semana en que queda completa
        { id: 'tokens', n: 1, label: 'Tokens', sub: 'Lexer', start: 3, done: 3 },
        { id: 'ast', n: 2, label: 'AST', sub: 'Parser', start: 4, done: 5 },
        { id: 'eval', n: 3, label: 'Evaluación', sub: 'Evaluador', start: 6, done: 7 }
    ];

    const SAMPLES = [
        { label: 'Variables', code: 'let five = 5;\nlet ten = 10;\nfive + ten * 2;' },
        { label: 'Precedencia', code: '-a * b + c / (d - e) == f' },
        { label: 'if / else', code: 'let x = 3;\nif (x > 2) { x * 10 } else { 0 }' },
        { label: 'Closure', code: 'let newAdder = fn(x) { fn(y) { x + y } };\nlet addTwo = newAdder(2);\naddTwo(40);' },
        { label: 'Fibonacci', code: 'let fib = fn(n) {\n  if (n < 2) { n } else { fib(n - 1) + fib(n - 2) }\n};\nfib(15);' },
        { label: 'Arrays y hashes', code: 'let p = {"nombre": "Zelda", "notas": [4, 5, 3]};\nputs(p["nombre"]);\nlen(p["notas"]) + first(p["notas"]);' },
        { label: 'Error', code: 'let x = 5;\nx + true;' }
    ];

    let uid = 0;

    function el(tag, attrs, html) {
        const e = document.createElement(tag);
        Object.entries(attrs || {}).forEach(([k, v]) => e.setAttribute(k, v));
        if (html !== undefined) e.innerHTML = html;
        return e;
    }

    function mount(root) {
        const M = window.Monkey;
        if (!M) return;
        const id = 'zsim' + (++uid);
        const week = Number(document.body.dataset.weekNum || document.body.dataset.week || 0);
        let stage = root.dataset.stage || 'tokens';
        const icon = (n) => (window.COURSE ? window.COURSE.svg(n) : '');
        const brand = window.COURSE ? window.COURSE.BRAND : 'ZeldaInterpreter';

        root.classList.add('zsim');
        root.innerHTML = `
            <div class="zsim-head">
                <span class="zsim-logo" data-logo="44"></span>
                <div>
                    <div class="zsim-title">Simulador ${brand}</div>
                    <div class="zsim-sub">Escribe Monkey y observa las tres fases del intérprete en vivo</div>
                </div>
            </div>
            <label class="viz-label" for="${id}-input">Código Monkey</label>
            <textarea id="${id}-input" class="playground-input" spellcheck="false" rows="4"></textarea>
            <div class="playground-actions zsim-samples" aria-label="Ejemplos"></div>
            <div class="zsim-tabs" role="tablist" aria-label="Fases del intérprete"></div>
            <div class="zsim-panel" role="tabpanel" id="${id}-panel" tabindex="0" aria-live="polite"></div>`;

        const input = root.querySelector('textarea');
        const tabs = root.querySelector('.zsim-tabs');
        const panel = root.querySelector('.zsim-panel');
        const samples = root.querySelector('.zsim-samples');
        input.value = root.dataset.code || SAMPLES[0].code;

        SAMPLES.forEach((s) => {
            const b = el('button', { type: 'button', class: 'btn-pill' }, M.esc(s.label));
            b.addEventListener('click', () => { input.value = s.code; render(); });
            samples.appendChild(b);
        });

        STAGES.forEach((st, i) => {
            const status = !week ? '' : week >= st.done ? 'construido' : week >= st.start ? 'en construcción (hoy)' : `semana ${st.start}`;
            const b = el('button', {
                type: 'button', role: 'tab', id: `${id}-tab-${st.id}`, 'aria-controls': `${id}-panel`, 'data-stage': st.id
            }, `<span class="zsim-step">${st.n}</span><span class="zsim-tab-text"><b>${st.label}</b><small>${st.sub}${status ? ' · ' + status : ''}</small></span>`);
            if (status === 'construido') b.classList.add('is-built');
            if (status.startsWith('en construcción')) b.classList.add('is-building');
            b.addEventListener('click', () => select(st.id, false));
            b.addEventListener('keydown', (e) => {
                if (!['ArrowRight', 'ArrowLeft', 'Home', 'End'].includes(e.key)) return;
                e.preventDefault();
                e.stopPropagation(); // no cambiar de diapositiva
                let j = i;
                if (e.key === 'ArrowRight') j = (i + 1) % STAGES.length;
                if (e.key === 'ArrowLeft') j = (i - 1 + STAGES.length) % STAGES.length;
                if (e.key === 'Home') j = 0;
                if (e.key === 'End') j = STAGES.length - 1;
                select(STAGES[j].id, true);
            });
            tabs.appendChild(b);
            if (i < STAGES.length - 1) tabs.appendChild(el('span', { class: 'zsim-arrow', 'aria-hidden': 'true' }, icon('arrow-right')));
        });

        function select(s, focus) {
            stage = s;
            tabs.querySelectorAll('[role="tab"]').forEach((t) => {
                const on = t.dataset.stage === s;
                t.setAttribute('aria-selected', String(on));
                t.tabIndex = on ? 0 : -1;
                if (on) {
                    panel.setAttribute('aria-labelledby', t.id);
                    if (focus) t.focus();
                }
            });
            render();
        }

        function renderTokens(src) {
            const toks = M.lex(src);
            const chips = toks.map((t) =>
                `<span class="token-chip${t.type === 'ILLEGAL' ? ' tok-illegal' : ''}"><span class="tok-type">${M.esc(t.type)}</span><span class="tok-lit">${M.esc(JSON.stringify(t.literal))}</span></span>`).join('');
            return `<div class="viz-label">${toks.length} tokens · <code>lexer.NextToken()</code> hasta <code>EOF</code></div><div class="token-stream">${chips}</div>`;
        }

        function renderAst(src) {
            const r = M.parse(src);
            if (r.errors.length) {
                return `<div class="viz-label">p.Errors()</div><div class="viz-output error">${M.esc(r.errors.join('\n'))}</div>`;
            }
            const tree = r.program.statements.length === 1 ? M.treeHtml(r.program.statements[0]) : M.treeHtml(r.program);
            return `<div class="viz-label">program.String()</div><div class="viz-output">${M.esc(M.str(r.program) || '(vacío)')}</div>
                    <div class="viz-label">Árbol de sintaxis abstracta</div><div class="viz-panel zsim-tree"><div class="ast-tree">${tree}</div></div>`;
        }

        function renderEval(src) {
            const r = M.run(src);
            if (r.parseErrors) {
                return `<div class="viz-label">El parser encontró errores (no se evalúa)</div><div class="viz-output error">${M.esc(r.parseErrors.join('\n'))}</div>`;
            }
            const isErr = r.result && r.result.type === 'ERROR';
            const res = r.result === undefined ? '(sin valor: la última sentencia es un let)' : M.inspect(r.result);
            const type = r.result && r.result.type ? `<span class="zsim-type">${M.esc(r.result.type)}</span>` : '';
            const binds = [...r.env.store.entries()].map(([k, v]) => {
                const val = v.type === 'FUNCTION' ? `fn(${v.parameters.map((p) => p.value).join(', ')})` : M.inspect(v);
                return `<div class="bind"><span>${M.esc(k)}</span><span>${M.esc(val.length > 48 ? val.slice(0, 45) + '…' : val)}</span></div>`;
            }).join('');
            return `<div class="viz-label">Resultado · Inspect() ${type}</div><div class="viz-output${isErr ? ' error' : ''}">${M.esc(res)}</div>
                <div class="viz-grid">
                    <div class="env-frame"><h4>Entorno global</h4>${binds || '<span style="color:var(--text-muted)">(vacío)</span>'}</div>
                    <div><div class="viz-label">Salida de puts</div><pre class="viz-panel zsim-puts">${M.esc(r.output.length ? r.output.join('\n') : '(nada)')}</pre></div>
                </div>`;
        }

        let timer = null;
        function render() {
            const src = input.value;
            panel.innerHTML = stage === 'tokens' ? renderTokens(src) : stage === 'ast' ? renderAst(src) : renderEval(src);
        }
        input.addEventListener('input', () => { clearTimeout(timer); timer = setTimeout(render, 150); });

        if (window.COURSE) window.COURSE.hydrateLogos(root);
        select(STAGES.some((s) => s.id === stage) ? stage : 'tokens', false);
    }

    document.addEventListener('DOMContentLoaded', () => {
        document.querySelectorAll('[data-simulator]').forEach(mount);
    });
})();

/**
 * Widget REPL reutilizable para las clases 6–8.
 * Requiere monkey.js. Busca: #repl-input, #repl-run, #repl-output,
 * botones [data-repl] y, opcionalmente, #repl-env (bindings globales) y #repl-puts (salida de puts).
 */
(function () {
    'use strict';
    document.addEventListener('DOMContentLoaded', () => {
        const input = document.getElementById('repl-input');
        const out = document.getElementById('repl-output');
        if (!input || !out || !window.Monkey) return;
        const envEl = document.getElementById('repl-env');
        const putsEl = document.getElementById('repl-puts');
        const M = window.Monkey;

        function short(o) {
            const s = M.inspect(o);
            return s.length > 60 ? s.slice(0, 57) + '…' : s;
        }

        function run() {
            const r = M.run(input.value);
            if (r.parseErrors) {
                out.textContent = 'parser errors:\n' + r.parseErrors.map((e) => '  ' + e).join('\n');
                out.classList.add('error');
                input.setAttribute('aria-invalid', 'true');
                if (envEl) envEl.innerHTML = '';
                if (putsEl) putsEl.textContent = '';
                return;
            }
            input.removeAttribute('aria-invalid');
            const isErr = r.result && r.result.type === 'ERROR';
            out.classList.toggle('error', !!isErr);
            out.textContent = r.result === undefined ? '(sin valor: la última sentencia fue un let)' : M.inspect(r.result);
            if (putsEl) putsEl.textContent = r.output.length ? r.output.join('\n') : '(puts no imprimió nada)';
            if (envEl) {
                const rows = [...r.env.store.entries()];
                envEl.innerHTML = rows.length
                    ? '<h4>Entorno global (env.store)</h4>' + rows.map(([k, v]) =>
                        `<div class="bind"><span>${M.esc(k)}</span><span>${M.esc(v.type === 'FUNCTION' ? 'FUNCTION · fn(' + v.parameters.map((p) => p.value).join(', ') + ')' : v.type + ' · ' + short(v))}</span></div>`).join('')
                    : '<h4>Entorno global (env.store)</h4><span style="color:var(--text-muted)">(vacío)</span>';
            }
        }

        const runBtn = document.getElementById('repl-run');
        if (runBtn) runBtn.addEventListener('click', run);
        input.addEventListener('keydown', (e) => {
            if (e.key === 'Enter' && (e.metaKey || e.ctrlKey)) { e.preventDefault(); run(); }
        });
        document.querySelectorAll('[data-repl]').forEach((b) => b.addEventListener('click', () => {
            input.value = b.dataset.repl;
            run();
        }));
        run();
    });
})();

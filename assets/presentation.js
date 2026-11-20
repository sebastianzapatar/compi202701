/**
 * Controlador Híbrido: Modo Guía (Sidebar + Scroll) <-> Modo Presentación (Slides)
 * Lenguajes y Compiladores 2027-1 — Intérprete Monkey en Go
 */
(function () {
    'use strict';

    document.addEventListener('DOMContentLoaded', () => {
        initSyntaxHighlighting();
        initCopyButtons();
        initScrollSpy();
        initMobileMenu();
        initPresentationMode();
    });

    const icon = (name) => (window.COURSE ? window.COURSE.svg(name) : '');

    /* Menú lateral colapsable en pantallas pequeñas (< 960px) */
    function initMobileMenu() {
        const sidebar = document.querySelector('.sidebar');
        if (!sidebar) return;
        if (!sidebar.id) sidebar.id = 'sidebar';

        const bar = document.createElement('div');
        bar.className = 'mobile-bar';
        const C = window.COURSE;
        const brandName = C ? C.BRAND : 'ZeldaInterpreter';
        bar.innerHTML = `
            <a class="mobile-brand" href="${sidebar.querySelector('.logo a')?.getAttribute('href') || '#'}">${C ? C.logoSvg(30) : ''}<span class="brand-name" style="font-size:1rem">${brandName}</span></a>
            <button type="button" class="mobile-menu-btn" aria-controls="${sidebar.id}" aria-expanded="false" aria-label="Abrir menú de navegación">${icon('menu')}<span>Menú</span></button>`;
        document.body.insertBefore(bar, document.body.firstChild);

        const btn = bar.querySelector('.mobile-menu-btn');
        const setOpen = (open) => {
            document.body.classList.toggle('menu-open', open);
            btn.setAttribute('aria-expanded', String(open));
            btn.setAttribute('aria-label', open ? 'Cerrar menú de navegación' : 'Abrir menú de navegación');
            btn.innerHTML = `${icon(open ? 'x' : 'menu')}<span>${open ? 'Cerrar' : 'Menú'}</span>`;
        };
        btn.addEventListener('click', () => setOpen(!document.body.classList.contains('menu-open')));
        sidebar.querySelectorAll('a').forEach((a) => a.addEventListener('click', () => setOpen(false)));
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape' && document.body.classList.contains('menu-open')) setOpen(false);
        });
    }

    function initSyntaxHighlighting() {
        if (window.hljs && typeof window.hljs.highlightAll === 'function') {
            window.hljs.highlightAll();
        }
    }

    function initCopyButtons() {
        document.querySelectorAll('.code-block').forEach((block) => {
            const codeEl = block.querySelector('pre code');
            if (!codeEl) return;

            let header = block.querySelector('.code-header');
            if (!header) {
                header = document.createElement('div');
                header.className = 'code-header';
                const langClass = Array.from(codeEl.classList).find((c) => c.startsWith('language-'));
                const label = document.createElement('span');
                label.textContent = langClass ? langClass.replace('language-', '').toUpperCase() : 'CODE';
                header.appendChild(label);
                block.insertBefore(header, block.firstChild);
            }

            if (header.querySelector('.copy-btn')) return;

            const btn = document.createElement('button');
            btn.type = 'button';
            btn.className = 'copy-btn';
            btn.textContent = 'Copiar';
            btn.addEventListener('click', async () => {
                try {
                    await navigator.clipboard.writeText(codeEl.innerText);
                    btn.textContent = '¡Copiado!';
                    setTimeout(() => {
                        btn.textContent = 'Copiar';
                    }, 1600);
                } catch {
                    btn.textContent = 'Error';
                    setTimeout(() => {
                        btn.textContent = 'Copiar';
                    }, 1600);
                }
            });
            header.appendChild(btn);
        });
    }

    function initScrollSpy() {
        const navLinks = Array.from(document.querySelectorAll('.sidebar nav a[href^="#"]'));
        if (!navLinks.length) return;

        const targets = navLinks
            .map((a) => document.querySelector(a.getAttribute('href')))
            .filter(Boolean);

        const updateActive = () => {
            if (document.body.classList.contains('slides-mode')) return;
            const scrollPos = window.scrollY + 140;
            let currentId = targets[0]?.id || '';
            for (const sec of targets) {
                if (sec.offsetTop <= scrollPos) {
                    currentId = sec.id;
                }
            }
            navLinks.forEach((link) => {
                const active = link.getAttribute('href') === `#${currentId}`;
                link.classList.toggle('active', active);
            });
        };

        window.addEventListener('scroll', updateActive, { passive: true });
        updateActive();
    }

    function initPresentationMode() {
        const slides = Array.from(document.querySelectorAll('.slide'));
        const toggleBtn = document.getElementById('toggle-slides-btn');
        if (!slides.length || !toggleBtn) return;

        let currentIndex = 0;

        // Build Presentation HUD dynamically
        const hud = document.createElement('div');
        hud.className = 'slide-hud';
        hud.setAttribute('aria-label', 'Controles de presentación');
        hud.innerHTML = `
            <div class="slide-progress-track"><div class="slide-progress-bar" id="hud-progress"></div></div>
            <div class="hud-left">
                <button type="button" class="hud-btn exit" id="hud-exit">${icon('x')}<span>Salir <kbd>Esc</kbd></span></button>
                <select class="hud-select" id="hud-jump" aria-label="Ir a diapositiva"></select>
            </div>
            <div class="hud-center">
                <button type="button" class="hud-btn" id="hud-prev" aria-label="Diapositiva anterior">${icon('arrow-left')}<span>Anterior</span></button>
                <span class="hud-counter" id="hud-counter" aria-live="polite">1 / ${slides.length}</span>
                <button type="button" class="hud-btn" id="hud-next" aria-label="Diapositiva siguiente"><span>Siguiente</span>${icon('arrow-right')}</button>
            </div>
            <div class="hud-right">
                <button type="button" class="hud-btn" id="hud-fullscreen" aria-label="Pantalla completa">${icon('maximize')}<span>Pantalla completa</span></button>
            </div>
        `;
        document.body.appendChild(hud);

        const progressEl = document.getElementById('hud-progress');
        const counterEl = document.getElementById('hud-counter');
        const prevBtn = document.getElementById('hud-prev');
        const nextBtn = document.getElementById('hud-next');
        const exitBtn = document.getElementById('hud-exit');
        const jumpSelect = document.getElementById('hud-jump');
        const fsBtn = document.getElementById('hud-fullscreen');

        slides.forEach((slide, idx) => {
            const titleEl = slide.querySelector('h1, h2');
            const title = titleEl ? titleEl.textContent.trim() : `Diapositiva ${idx + 1}`;
            const opt = document.createElement('option');
            opt.value = String(idx);
            opt.textContent = `${idx + 1}. ${title}`;
            jumpSelect.appendChild(opt);
        });

        function renderSlide(index) {
            currentIndex = Math.max(0, Math.min(slides.length - 1, index));
            slides.forEach((s, i) => {
                s.classList.toggle('is-active', i === currentIndex);
                if (i === currentIndex) s.scrollTop = 0;
            });
            counterEl.textContent = `${currentIndex + 1} / ${slides.length}`;
            progressEl.style.width = `${((currentIndex + 1) / slides.length) * 100}%`;
            prevBtn.disabled = currentIndex === 0;
            nextBtn.disabled = currentIndex === slides.length - 1;
            jumpSelect.value = String(currentIndex);
            if (document.body.classList.contains('slides-mode')) {
                history.replaceState(null, '', `#slide-${currentIndex + 1}`);
            }
        }

        function enterSlides(startIndex = 0) {
            document.body.classList.add('slides-mode');
            toggleBtn.setAttribute('aria-pressed', 'true');
            renderSlide(startIndex);
            nextBtn.focus({ preventScroll: true });
        }

        function exitSlides() {
            document.body.classList.remove('slides-mode');
            toggleBtn.setAttribute('aria-pressed', 'false');
            slides.forEach((s) => s.classList.remove('is-active'));
            if (document.fullscreenElement) {
                document.exitFullscreen().catch(() => {});
            }
            const target = slides[currentIndex];
            if (target && target.id) {
                history.replaceState(null, '', `#${target.id}`);
                target.scrollIntoView({ behavior: 'auto', block: 'start' });
            } else {
                history.replaceState(null, '', window.location.pathname);
            }
        }

        toggleBtn.addEventListener('click', () => {
            if (document.body.classList.contains('slides-mode')) {
                exitSlides();
            } else {
                enterSlides(0);
            }
        });

        exitBtn.addEventListener('click', exitSlides);
        prevBtn.addEventListener('click', () => renderSlide(currentIndex - 1));
        nextBtn.addEventListener('click', () => renderSlide(currentIndex + 1));
        jumpSelect.addEventListener('change', (e) => renderSlide(Number(e.target.value)));

        fsBtn.addEventListener('click', () => {
            if (!document.fullscreenElement) {
                document.documentElement.requestFullscreen().catch(() => {});
            } else {
                document.exitFullscreen().catch(() => {});
            }
        });

        document.addEventListener('keydown', (e) => {
            const tag = (e.target && e.target.tagName) ? e.target.tagName.toLowerCase() : '';
            if (tag === 'input' || tag === 'textarea' || tag === 'select') return;

            const inSlides = document.body.classList.contains('slides-mode');

            if ((e.key === 'p' || e.key === 'P') && !e.metaKey && !e.ctrlKey && !e.altKey) {
                e.preventDefault();
                if (inSlides) exitSlides();
                else enterSlides(currentIndex);
                return;
            }

            if (!inSlides) return;

            if (e.key === 'ArrowRight' || e.key === 'PageDown' || e.key === ' ') {
                e.preventDefault();
                renderSlide(currentIndex + 1);
            } else if (e.key === 'ArrowLeft' || e.key === 'PageUp') {
                e.preventDefault();
                renderSlide(currentIndex - 1);
            } else if (e.key === 'Home') {
                e.preventDefault();
                renderSlide(0);
            } else if (e.key === 'End') {
                e.preventDefault();
                renderSlide(slides.length - 1);
            } else if (e.key === 'Escape') {
                e.preventDefault();
                exitSlides();
            }
        });

        // En modo diapositivas, un enlace interno (#id) salta a la diapositiva que contiene ese id
        document.addEventListener('click', (e) => {
            if (!document.body.classList.contains('slides-mode')) return;
            const a = e.target.closest('a[href^="#"]');
            if (!a) return;
            const target = document.getElementById(a.getAttribute('href').slice(1));
            const slide = target && target.closest('.slide');
            if (!slide) return;
            e.preventDefault();
            renderSlide(slides.indexOf(slide));
        });

        // Support direct links like ?mode=slides or #slide-3
        const params = new URLSearchParams(window.location.search);
        const hashMatch = window.location.hash.match(/^#slide-(\d+)$/);
        if (params.get('mode') === 'slides' || hashMatch) {
            const initial = hashMatch ? Number(hashMatch[1]) - 1 : 0;
            enterSlides(initial);
        }
    }
})();

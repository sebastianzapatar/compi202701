#!/usr/bin/env python3
"""
Genera las guías de ejercicios del curso (HTML imprimible + PDF) a partir de
scripts/ejercicios_contenido.py.

Salidas:
  ejercicios/fuente/semana-0N.html               (versión estudiante, editable/imprimible)
  ejercicios/pdf/semana-0N-ejercicios.pdf        (PDF para estudiantes)
  ejercicios/solucionario/semana-0N-solucionario.pdf  (PDF con respuestas, SOLO docente)

Uso:
  python3 scripts/build_ejercicios.py            # todas las semanas
  python3 scripts/build_ejercicios.py 4 5        # solo semanas 4 y 5
  python3 scripts/build_ejercicios.py go3 go4    # solo módulos Go III y Go IV
  python3 scripts/build_ejercicios.py --html     # solo HTML, sin PDF

Requiere Google Chrome instalado (se usa en modo headless para imprimir a PDF).
"""
import html
import os
import shutil
import subprocess
import sys
import tempfile
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
sys.path.insert(0, str(Path(__file__).resolve().parent))
from ejercicios_contenido import WEEKS, COURSE  # noqa: E402

OUT_SRC = ROOT / "ejercicios" / "fuente"
OUT_PDF = ROOT / "ejercicios" / "pdf"
OUT_SOL = ROOT / "ejercicios" / "solucionario"

CHROME_CANDIDATES = [
    "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
    "/Applications/Chromium.app/Contents/MacOS/Chromium",
    shutil.which("google-chrome") or "",
    shutil.which("chromium") or "",
    shutil.which("chromium-browser") or "",
]

esc = html.escape


def code_block(src, label=None):
    head = f'<div class="code-label">{esc(label)}</div>' if label else ""
    return f'<div class="code">{head}<pre>{esc(src.strip(chr(10)))}</pre></div>'


def answer_lines(n):
    return f'<div class="lines" style="--n:{n}" aria-hidden="true"></div>' if n else ""


def render_item(item, num, solution):
    parts = [f'<div class="item"><div class="item-head"><span class="qnum">{num}</span><div class="qtext">{item["q"]}']
    if item.get("pts"):
        parts.append(f' <span class="pts">({item["pts"]} pts)</span>')
    parts.append("</div></div>")
    body = []
    if item.get("code"):
        body.append(code_block(item["code"], item.get("code_label")))
    if item.get("options"):
        opts = "".join(
            f'<li><span class="box{" checked" if solution and i in item.get("correct", []) else ""}"></span>{o}</li>'
            for i, o in enumerate(item["options"])
        )
        body.append(f'<ul class="options">{opts}</ul>')
    if item.get("table"):
        headers, rows = item["table"]
        ans_rows = item.get("table_ans")
        th = "".join(f"<th>{h}</th>" for h in headers)
        trs = []
        for r_i, row in enumerate(rows):
            tds = []
            for c_i, cell in enumerate(row):
                if cell == "" and solution and ans_rows:
                    tds.append(f'<td class="fill sol">{ans_rows[r_i][c_i]}</td>')
                elif cell == "":
                    tds.append('<td class="fill"></td>')
                else:
                    tds.append(f"<td>{cell}</td>")
            trs.append("<tr>" + "".join(tds) + "</tr>")
        body.append(f'<table class="fill-table"><thead><tr>{th}</tr></thead><tbody>{"".join(trs)}</tbody></table>')
    if item.get("draw"):
        if solution and item.get("ans"):
            pass
        else:
            body.append(f'<div class="draw" style="height:{item["draw"]}mm"><span>Espacio para dibujar</span></div>')
    if item.get("body"):
        body.append(item["body"])
    if solution and item.get("ans"):
        body.append(f'<div class="ans"><span class="ans-tag">Respuesta</span>{item["ans"]}</div>')
    elif not item.get("draw") and not item.get("table") and not item.get("options"):
        body.append(answer_lines(item.get("lines", 3)))
    elif item.get("options") and item.get("lines"):
        body.append(answer_lines(item["lines"]))
    parts.append(f'<div class="item-body">{"".join(body)}</div></div>')
    return "".join(parts)


def render_problem(p, num, solution):
    out = [f'<div class="problem"><div class="problem-head"><span class="qnum">{num}</span><h3>{p["title"]}</h3>']
    if p.get("pts"):
        out.append(f'<span class="pts">{p["pts"]} pts</span>')
    out.append(f'</div><div class="problem-body">{p["body"]}')
    if p.get("code"):
        out.append(code_block(p["code"], p.get("code_label")))
    if p.get("after"):
        out.append(p["after"])
    if solution and p.get("ans"):
        out.append(f'<div class="ans"><span class="ans-tag">Criterios / solución de referencia</span>{p["ans"]}</div>')
    out.append("</div></div>")
    return "".join(out)


def meta(week):
    """Etiquetas de una semana numerada (1..8) o de un módulo ("go3".."go6")."""
    w = WEEKS[week]
    if isinstance(week, int):
        return {"slug": f"semana-{week:02d}", "label": f"Semana {week:02d}",
                "kicker": f"Semana {week:02d} de 08", "badge": f"{week:02d}"}
    return {"slug": w["slug"], "label": w["label"], "kicker": w["label"] + " · Bloque I ampliado", "badge": w["badge"]}


def render(week, solution=False):
    w = WEEKS[week]
    m = meta(week)
    n = 0
    sections = []
    for key, label, cls in (("A", "Parte A · Conceptos", "pa"), ("B", "Parte B · Trazas a mano", "pb")):
        items = []
        for it in w[key]:
            n += 1
            items.append(render_item(it, n, solution))
        pts = sum(i.get("pts", 0) for i in w[key])
        sections.append(
            f'<section class="part {cls}"><h2><span>{label}</span><span class="part-pts">{pts} pts</span></h2>'
            f'<p class="part-desc">{w[key + "_desc"]}</p>{"".join(items)}</section>'
        )
    probs = []
    for p in w["C"]:
        n += 1
        probs.append(render_problem(p, n, solution))
    pts_c = sum(p.get("pts", 0) for p in w["C"])
    sections.append(
        f'<section class="part pc"><h2><span>Parte C · Programación en Go</span><span class="part-pts">{pts_c} pts</span></h2>'
        f'<p class="part-desc">{w["C_desc"]}</p>{"".join(probs)}</section>'
    )
    if w.get("reto"):
        r = w["reto"]
        sol = f'<div class="ans"><span class="ans-tag">Pista para el docente</span>{r["ans"]}</div>' if solution and r.get("ans") else ""
        sections.append(
            f'<section class="reto"><h2><span>Reto opcional · {r["title"]}</span><span class="part-pts">+{r.get("pts", 10)} pts</span></h2>{r["body"]}{sol}</section>'
        )

    rub = "".join(f"<tr><td>{c}</td><td>{d}</td><td class='num'>{p}</td></tr>" for c, d, p in w.get("rubric", COURSE["rubric"]))
    goals = "".join(f"<li>{g}</li>" for g in w["goals"])
    sol_banner = '<div class="sol-banner">SOLUCIONARIO · USO EXCLUSIVO DEL DOCENTE · NO DISTRIBUIR</div>' if solution else ""
    title_suffix = " · Solucionario" if solution else ""

    return f"""<!DOCTYPE html>
<html lang="es">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{m['label']} · {esc(w['short'])}{title_suffix} | Guía de ejercicios</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&family=JetBrains+Mono:wght@400;600&display=swap" rel="stylesheet">
<link rel="stylesheet" href="worksheet.css">
</head>
<body class="{'solution' if solution else 'student'}">
{sol_banner}
<header class="sheet-header">
  <div class="brand">
    <span class="course">ZeldaInterpreter · {COURSE['name']} · {COURSE['term']}</span>
    <span class="book">Libro guía: <em>Writing An Interpreter In Go</em> — Thorsten Ball</span>
  </div>
  <div class="title-row">
    <div>
      <div class="kicker">Guía de ejercicios · {m['kicker']}{title_suffix}</div>
      <h1>{w['title']}</h1>
      <div class="ref">{w['ref']}</div>
    </div>
    <div class="week-badge" aria-hidden="true">{m['badge']}</div>
  </div>
</header>

<div class="student-fields">
  <div class="field"><label>Nombre</label><span></span></div>
  <div class="field"><label>Código</label><span></span></div>
  <div class="field small"><label>Fecha</label><span></span></div>
  <div class="field small"><label>Nota</label><span></span></div>
</div>

<div class="intro">
  <div class="goals">
    <h2>Al terminar esta guía podrás…</h2>
    <ul>{goals}</ul>
  </div>
  <div class="instructions">
    <h2>Instrucciones</h2>
    <ul>
      <li>Partes A y B: a mano, <strong>sin computador</strong>. Razona con el código del libro.</li>
      <li>Parte C: entrega en tu repositorio <code>monkey</code>; <code>go test ./...</code> debe pasar.</li>
      <li>Puntaje total: <strong>100 pts</strong> (+ reto opcional).</li>
      <li>Clase de referencia: <code>clases/{COURSE['files'][week]}</code></li>
    </ul>
  </div>
</div>

{''.join(sections)}

<section class="rubric">
  <h2><span>Rúbrica de evaluación</span></h2>
  <table>
    <thead><tr><th>Criterio</th><th>Descripción</th><th class="num">Pts</th></tr></thead>
    <tbody>{rub}</tbody>
  </table>
</section>

<footer class="sheet-footer">ZeldaInterpreter · {COURSE['name']} {COURSE['term']} · {m['label']} · {esc(w['short'])}</footer>
</body>
</html>
"""


def find_chrome():
    for c in CHROME_CANDIDATES:
        if c and os.path.exists(c):
            return c
    return None


def print_pdf(chrome, html_path, pdf_path):
    with tempfile.TemporaryDirectory() as profile:
        cmd = [
            chrome, "--headless=new", "--disable-gpu", "--no-first-run",
            f"--user-data-dir={profile}", "--no-pdf-header-footer",
            "--virtual-time-budget=6000", "--run-all-compositor-stages-before-draw",
            f"--print-to-pdf={pdf_path}", html_path.as_uri(),
        ]
        subprocess.run(cmd, check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, timeout=120)


def main(argv):
    only_html = "--html" in argv
    order = [k for k in WEEKS if not isinstance(k, int)] + sorted(k for k in WEEKS if isinstance(k, int))
    wanted = [int(a) if a.isdigit() else a for a in argv if not a.startswith("--")]
    weeks = [k for k in order if not wanted or k in wanted]
    for d in (OUT_SRC, OUT_PDF, OUT_SOL):
        d.mkdir(parents=True, exist_ok=True)
    chrome = None if only_html else find_chrome()
    if not only_html and not chrome:
        print("No se encontró Chrome/Chromium: se generan solo los HTML.")
    sol_tmp = Path(tempfile.mkdtemp())
    shutil.copy(OUT_SRC / "worksheet.css", sol_tmp / "worksheet.css")
    for wk in weeks:
        slug = meta(wk)["slug"]
        src = OUT_SRC / f"{slug}.html"
        src.write_text(render(wk, solution=False), encoding="utf-8")
        sol = sol_tmp / f"{slug}-solucionario.html"
        sol.write_text(render(wk, solution=True), encoding="utf-8")
        msg = f"{meta(wk)['label']}: {src.relative_to(ROOT)}"
        if chrome:
            print_pdf(chrome, src, OUT_PDF / f"{slug}-ejercicios.pdf")
            print_pdf(chrome, sol, OUT_SOL / f"{slug}-solucionario.pdf")
            msg += " → PDF estudiante + solucionario"
        print(msg)
    shutil.rmtree(sol_tmp, ignore_errors=True)


if __name__ == "__main__":
    main(sys.argv[1:])

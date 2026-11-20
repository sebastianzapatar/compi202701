#!/usr/bin/env python3
"""
Genera la sección «Córrelo en ZeldaInterpreter» de las clases (semanas 3–8) y
del portal: comandos de ejemplo para correr el intérprete de zeldainterpreter/
junto con su salida REAL.

Cada ejemplo se ejecuta de verdad con el intérprete en Go, así que la salida que
ven los estudiantes es exactamente la que obtendrán en su terminal. La sección
queda entre los marcadores <!-- CORRER:INICIO --> y <!-- CORRER:FIN --> (se
reemplaza en cada ejecución; no la edites a mano: edita EXAMPLES aquí abajo).

Uso:
  python3 scripts/build_correr.py          # regenera todas las páginas
  python3 scripts/build_correr.py 5 8      # solo las semanas 5 y 8
  python3 scripts/build_correr.py index    # solo el portal

Requiere Go (compila zeldainterpreter/ en una carpeta temporal).
"""
import html
import re
import subprocess
import sys
import tempfile
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
INTERP = ROOT / "zeldainterpreter"
START, END = "<!-- CORRER:INICIO -->", "<!-- CORRER:FIN -->"
NAV_LINK = '<li><a href="#correr">Córrelo en ZeldaInterpreter</a></li>'

esc = html.escape

# --------------------------------------------------------------------------
# Contenido. Cada ejemplo tiene:
#   title  título corto            note  explicación (HTML)
#   args   argumentos de `zelda`   (se muestran como `go run . <args>`)
#   repl   líneas que se escriben en la consola (si es una sesión interactiva)
#   stdin  archivo que se envía por la entrada estándar (`cat x | go run . -`)
#   source archivo .mk que se muestra antes del comando
#   exit   True para mostrar también `echo $?` (código de salida)
# --------------------------------------------------------------------------
PAGES = {
    3: {
        "file": "clases/semana-03-lexing.html",
        "intro": "Tu lexer ya existe en <code>zeldainterpreter/</code>. Con <code>-modo=tokens</code> el intérprete se detiene después del análisis léxico y muestra la tabla de tokens: compárala con lo que produce tu <code>NextToken()</code>.",
        "test": "go test ./lexer/",
        "examples": [
            {"title": "Tokens de una sentencia", "args": ["-modo=tokens", "-e", "let cinco = 5;"],
             "note": "Las palabras reservadas salen de <code>LookupIdent()</code>: <code>let</code> es <code>LET</code>, <code>cinco</code> es <code>IDENT</code>."},
            {"title": "Operadores de dos caracteres", "args": ["-modo=tokens", "-e", "10 == 10; 10 != 9;"],
             "note": "<code>==</code> y <code>!=</code> salen como un solo token gracias a <code>peekChar()</code>."},
            {"title": "La consola en modo tokens", "args": ["-modo=tokens"], "repl": ["x != 10", ":salir"],
             "note": "El REPL del Capítulo 1: lee una línea y la tokeniza."},
        ],
    },
    4: {
        "file": "clases/semana-04-parsing-sentencias.html",
        "intro": "Con <code>-modo=ast</code> el intérprete se detiene después del parser e imprime el <code>String()</code> del programa; con <code>-modo=arbol</code> dibuja el AST. Así puedes ver las sentencias <code>let</code> y <code>return</code> que construye tu parser.",
        "test": "go test ./ast/ ./parser/",
        "examples": [
            {"title": "let y return como texto", "args": ["-modo=ast", "-e", "let x = 5; return x;"],
             "note": "Una línea por sentencia de <code>program.Statements</code>."},
            {"title": "El AST dibujado", "args": ["-modo=arbol", "-e", "let edad = 20;"],
             "note": "<code>LetStatement</code> con su <code>Name</code> (Identifier) y su <code>Value</code>."},
            {"title": "Errores de sintaxis", "source": "programas/error_sintaxis.mk", "args": ["programas/error_sintaxis.mk"], "exit": True,
             "note": "Cada mensaje viene de <code>peekError()</code> o de <code>noPrefixParseFnError()</code>; el parser <strong>no se detiene</strong> en el primer error, los acumula en <code>Errors()</code>."},
        ],
    },
    5: {
        "file": "clases/semana-05-pratt-parsing.html",
        "intro": "Comprueba la precedencia de tu Pratt parser: <code>-modo=ast</code> muestra la expresión con paréntesis explícitos y <code>-modo=arbol</code> la dibuja.",
        "test": "go test ./parser/",
        "examples": [
            {"title": "Precedencia con paréntesis", "args": ["-modo=ast", "-e", "-a * b + c / d"],
             "note": "Prefijo (<code>PREFIX</code>) antes que producto (<code>PRODUCT</code>) antes que suma (<code>SUM</code>)."},
            {"title": "1 + 2 * 3 como árbol", "args": ["-modo=arbol", "-e", "1 + 2 * 3"],
             "note": "El <code>*</code> queda <strong>más abajo</strong> en el árbol: se evalúa primero."},
            {"title": "if, fn y llamadas", "args": ["-modo=arbol", "-e", "if (x < y) { add(x, 2) } else { y }"],
             "note": "<code>IfExpression</code>, <code>BlockStatement</code> y <code>CallExpression</code>, todos producidos por funciones <em>prefix</em>/<em>infix</em>."},
            {"title": "La consola cambia de modo", "args": ["-modo=ast"], "repl": ["a + b * c", ":modo arbol", "-x", ":salir"],
             "note": "Con <code>:modo</code> pasas de <code>ast</code> a <code>arbol</code> sin salir del REPL."},
        ],
    },
    6: {
        "file": "clases/semana-06-evaluacion-base.html",
        "intro": "Ahora sí se evalúa: el modo por defecto (<code>eval</code>) recorre el AST con <code>Eval</code> e imprime el <code>Inspect()</code> del resultado.",
        "test": "go test ./object/ ./evaluator/",
        "examples": [
            {"title": "Expresiones e if", "args": ["-e", "(5 + 10 * 2 + 15 / 3) * 2 + -10"],
             "note": "Con <code>-e</code> se muestra el valor final del programa."},
            {"title": "return dentro de bloques", "args": ["-e", "if (10 > 1) { if (10 > 1) { return 10; } return 1; }"],
             "note": "El <code>ReturnValue</code> burbujea hasta <code>evalProgram</code>, que lo desenvuelve."},
            {"title": "Un programa completo", "source": "programas/aritmetica.mk", "args": ["programas/aritmetica.mk"],
             "note": "Al ejecutar un archivo solo se ve lo que imprime <code>puts</code> (agrega <code>-resultado</code> para ver también el valor final)."},
            {"title": "Errores en tiempo de ejecución", "source": "programas/error_ejecucion.mk", "args": ["programas/error_ejecucion.mk"], "exit": True,
             "note": "El <code>*object.Error</code> corta la ejecución: la última línea nunca se imprime. El código de salida es <code>1</code>."},
        ],
    },
    7: {
        "file": "clases/semana-07-entornos-closures.html",
        "intro": "Con <code>Environment</code> el intérprete ya recuerda variables y funciones. La consola conserva el entorno entre líneas, y <code>:env</code> te muestra qué hay guardado en él.",
        "test": "go test ./object/ ./evaluator/",
        "examples": [
            {"title": "Closures", "source": "programas/closures.mk", "args": ["programas/closures.mk"],
             "note": "<code>sumarDos</code> y <code>sumarDiez</code> recuerdan cada una su propio <code>x</code>: el entorno capturado en <code>object.Function.Env</code>."},
            {"title": "Recursión", "args": ["programas/fibonacci.mk"],
             "note": "<code>fibonacci</code> se encuentra a sí misma porque el <code>let</code> la guarda en el entorno global antes de que se llame."},
            {"title": "El entorno en la consola", "args": [], "repl": ["let newAdder = fn(x) { fn(y) { x + y } };", "let sumarDos = newAdder(2);", "sumarDos(40)", ":env", ":salir"],
             "note": "<code>:env</code> lista lo que guardó <code>env.Set</code>. El <code>x = 2</code> no aparece: vive en el entorno encerrado de <code>sumarDos</code>."},
        ],
    },
    8: {
        "file": "clases/semana-08-extensiones.html",
        "intro": "El intérprete completo: strings, funciones integradas, arrays y hashes. Estos son los programas de <code>zeldainterpreter/programas/</code> y las otras formas de ejecutar código.",
        "test": "go test ./...",
        "examples": [
            {"title": "Strings y concatenación", "source": "programas/hola.mk", "args": ["programas/hola.mk"],
             "note": "<code>+</code> entre dos <code>STRING</code> lo resuelve <code>evalStringInfixExpression</code>."},
            {"title": "Arrays, hashes y built-ins", "args": ["-e", 'let p = {"nombre": "Zelda", "notas": [4, 5, 3]}; puts("Hola, " + p["nombre"]); len(p["notas"]) + first(p["notas"])'],
             "note": "Indexación de hash con <code>HashKey</code>, de array con <code>[i]</code>, y las funciones <code>len</code>, <code>first</code> y <code>puts</code>."},
            {"title": "The Grand Finale por la entrada estándar", "stdin": "programas/map_reduce.mk",
             "note": "Con <code>-</code> el programa se lee de la entrada estándar: sirve para encadenar con otros comandos."},
            {"title": "Cargar un archivo en la consola", "args": [], "repl": [":cargar programas/map_reduce.mk", "map([1, 2, 3], fn(x) { x * 10 })", "reduce([1, 2, 3], 1, fn(a, b) { a * b })", ":salir"],
             "note": "<code>:cargar</code> ejecuta el archivo <strong>en la sesión actual</strong>: sus <code>map</code> y <code>reduce</code> quedan disponibles. El <code>null</code> es el valor de la última línea del archivo: <code>puts</code> devuelve <code>null</code>."},
        ],
    },
    "index": {
        "file": "index.html",
        "title": "7. Cómo correr ZeldaInterpreter",
        "intro": "El intérprete de referencia está en la carpeta <code>zeldainterpreter/</code>. Necesitas <strong>Go 1.23 o superior</strong>. Todos los comandos se escriben <strong>dentro de esa carpeta</strong>; las opciones (<code>-modo</code>, <code>-e</code>…) van <strong>antes</strong> del nombre del archivo.",
        "test": "go test ./...",
        "modes": True,
        "examples": [
            {"title": "Ejecutar un archivo", "args": ["programas/fibonacci.mk"],
             "note": "Los programas de ejemplo están en <code>programas/*.mk</code>."},
            {"title": "Ejecutar código directo", "args": ["-e", "let doble = fn(x) { x * 2 }; doble(21)"],
             "note": "Con <code>-e</code> se ejecuta el texto y se muestra el valor final."},
            {"title": "Ver el AST", "args": ["-modo=arbol", "-e", "let x = 1 + 2 * 3;"],
             "note": "Cambia <code>-modo</code> para detenerte en cualquier fase: <code>tokens</code>, <code>ast</code>, <code>arbol</code> o <code>eval</code>."},
            {"title": "Consola interactiva", "args": [], "repl": ["let saludo = \"Hola, \" + \"Monkey\";", "len(saludo)", ":salir"],
             "note": "Sin argumentos abre el REPL. Escribe <code>:ayuda</code> para ver sus comandos."},
        ],
    },
}

# --------------------------------------------------------------------------


def shell_quote(arg):
    """Cita un argumento como lo escribiría una persona en bash/zsh."""
    if re.fullmatch(r"[A-Za-z0-9_./=:-]+", arg):
        return arg
    if "'" not in arg:
        return f"'{arg}'"
    return '"' + arg.replace("\\", "\\\\").replace('"', '\\"').replace("$", "\\$") + '"'


def run(binary, args, stdin_text=None):
    p = subprocess.run([str(binary), *args], cwd=INTERP, input=stdin_text or "",
                       capture_output=True, text=True, timeout=30)
    return (p.stdout + p.stderr), p.returncode


def repl_transcript(raw, lines):
    """Intercala lo que se escribió en la consola con lo que respondió.

    El REPL imprime ">> " (o ".. ") antes de leer cada línea, así que el texto
    entre dos prompts es la respuesta a la línea anterior."""
    parts = re.split(r"(>> |\.\. )", raw)
    out = [parts[0]]
    i = 0
    for j in range(1, len(parts), 2):
        prompt, answer = parts[j], parts[j + 1]
        typed = lines[i] if i < len(lines) else ""
        i += 1
        out.append(f"{prompt}{typed}\n{answer}")
    return "".join(out).rstrip("\n")


def block(header, lang, text):
    return (f'<div class="code-block"><div class="code-header"><span>{esc(header)}</span></div>\n'
            f'<pre><code class="language-{lang}">{esc(text)}</code></pre></div>')


def example_html(binary, ex):
    pieces = [f'<h3>{esc(ex["title"])}</h3>', f'<p>{ex["note"]}</p>']
    if ex.get("source"):
        src = (INTERP / ex["source"]).read_text(encoding="utf-8").rstrip("\n")
        pieces.append(block(ex["source"], "javascript", src))

    if ex.get("stdin"):
        cmd = f"cat {ex['stdin']} | go run . -"
        out, code = run(binary, ["-"], (INTERP / ex["stdin"]).read_text(encoding="utf-8"))
    elif "repl" in ex:
        cmd = " ".join(["go run .", *map(shell_quote, ex["args"])]).strip()
        raw, code = run(binary, ex["args"], "\n".join(ex["repl"]) + "\n")
        out = repl_transcript(raw, ex["repl"])
    else:
        cmd = " ".join(["go run .", *map(shell_quote, ex["args"])])
        out, code = run(binary, ex["args"])

    out = out.strip("\n")
    if ex.get("exit"):
        cmd += "\necho $?        # código de salida"
        out += f"\n{code}"
    elif code != 0:
        raise SystemExit(f"El ejemplo «{ex['title']}» terminó con código {code}:\n{out}")

    pieces.append(block("terminal", "bash", cmd))
    pieces.append(block("consola interactiva" if "repl" in ex else "salida", "plaintext", out))
    return '<div class="card">\n' + "\n".join(pieces) + "\n</div>"


MODES_TABLE = """<div class="table-wrap"><table>
<thead><tr><th>Cómo se llama</th><th>Qué hace</th></tr></thead>
<tbody>
<tr><td><code>go run .</code></td><td>Abre la <strong>consola interactiva</strong> (REPL). <code>:ayuda</code> lista sus comandos.</td></tr>
<tr><td><code>go run . programa.mk</code></td><td>Ejecuta un <strong>archivo</strong> <code>.mk</code>.</td></tr>
<tr><td><code>go run . -e "código"</code></td><td>Ejecuta <strong>código escrito en la línea de comandos</strong>.</td></tr>
<tr><td><code>cat programa.mk | go run . -</code></td><td>Ejecuta lo que llegue por la <strong>entrada estándar</strong>.</td></tr>
<tr><td><code>-modo=tokens | ast | arbol | eval</code></td><td>Elige <strong>hasta qué fase</strong> llega: lexer, parser (texto o árbol) o evaluación completa.</td></tr>
<tr><td><code>go build -o zelda .</code></td><td>Compila el ejecutable <code>zelda</code> (en Windows <code>zelda.exe</code>); luego <code>./zelda programa.mk</code>.</td></tr>
</tbody></table></div>"""


def section_html(binary, page, key):
    out, _ = run_cmd(["go", "test", *page["test"].split()[2:]])
    test_out = re.sub(r"\t[0-9.]+s$", "", out.strip(), flags=re.M)
    test_out = test_out.replace("\t(cached)", "")
    cards = [example_html(binary, ex) for ex in page["examples"]]
    grid = '<div class="grid-2 correr-grid">\n' + "\n".join(cards) + "\n</div>"
    title = page.get("title", "Córrelo en ZeldaInterpreter")
    extra = MODES_TABLE if page.get("modes") else ""
    where = ("" if key == "index" else
             '<div class="info-box"><strong>¿Y mi propio intérprete?</strong> Si construyes el tuyo siguiendo el libro, los mismos comandos funcionan sobre tu código a medida que avanzas. <code>zeldainterpreter/</code> es la versión de referencia completa para comparar resultados.</div>')
    return f"""{START}
        <!-- Generado por scripts/build_correr.py: no editar a mano. -->
        <section id="correr" class="section slide">
            <h2>{title}</h2>
            <p>{page["intro"]}</p>
            {extra}
            <div class="grid-2 correr-grid">
<div>{block("terminal · primero, las pruebas (deben salir en ok)", "bash", "cd zeldainterpreter\n" + page["test"])}</div>
<div>{block("salida", "plaintext", test_out)}</div>
            </div>
            {grid}
            {where}
        </section>
        {END}"""


def run_cmd(cmd):
    p = subprocess.run(cmd, cwd=INTERP, capture_output=True, text=True, timeout=300)
    if p.returncode != 0:
        raise SystemExit(f"Falló {' '.join(cmd)}:\n{p.stdout}{p.stderr}")
    return p.stdout + p.stderr, p.returncode


def inject(path, section, before_id):
    text = path.read_text(encoding="utf-8")
    if START in text:
        text = re.sub(re.escape(START) + r".*?" + re.escape(END), lambda _: section, text, flags=re.S)
    else:
        anchor = re.search(r'\n(\s*)(<!-- =+ [^\n]*-->\s*\n\s*)?<section id="' + before_id + '"', text)
        if not anchor:
            raise SystemExit(f"{path}: no encontré la sección #{before_id}")
        indent = anchor.group(1)
        text = text[:anchor.start()] + "\n" + indent + section + "\n" + text[anchor.start():]
    if 'href="#correr"' not in text:  # enlace en el menú lateral, antes del simulador
        m = re.search(r'(\s*)<li><a href="#simulador-zelda"', text)
        if not m:
            raise SystemExit(f"{path}: no encontré el enlace del simulador en el menú")
        text = text[:m.start()] + m.group(1) + NAV_LINK + text[m.start():]
    path.write_text(text, encoding="utf-8")


def main(argv):
    keys = [int(a) if a.isdigit() else a for a in argv] or list(PAGES)
    with tempfile.TemporaryDirectory() as tmp:
        binary = Path(tmp) / "zelda"
        run_cmd(["go", "build", "-o", str(binary), "."])
        for key in keys:
            page = PAGES[key]
            path = ROOT / page["file"]
            # En las clases va antes del simulador; en el portal, al final (después de Metodología).
            before = "simulador-zelda" if key != "index" else None
            section = section_html(binary, page, key)
            if before:
                inject(path, section, before)
            else:
                inject_end(path, section)
            print(f"{key}: {page['file']}")


def inject_end(path, section):
    text = path.read_text(encoding="utf-8")
    if START in text:
        text = re.sub(re.escape(START) + r".*?" + re.escape(END), lambda _: section, text, flags=re.S)
    else:
        i = text.rindex("</main>")
        text = text[:i] + "    " + section + "\n\n    " + text[i:]
    if 'href="#correr"' not in text:
        m = re.search(r'(\s*)(<li><a href="#methodology"[^\n]*</li>)', text)
        if not m:
            raise SystemExit(f"{path}: no encontré el enlace de Metodología en el menú")
        text = text[:m.end()] + m.group(1) + '<li><a href="#correr">7. Cómo correr ZeldaInterpreter</a></li>' + text[m.end():]
    path.write_text(text, encoding="utf-8")


if __name__ == "__main__":
    main(sys.argv[1:])

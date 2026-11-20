#!/usr/bin/env python3
"""
Incrusta el código Go de ejemplos/ (y su salida real) dentro de las clases HTML.

En el HTML se marcan los bloques así:
  <pre><code class="language-go" data-src="ejemplos/patrones/builder/main.go"></code></pre>
  <pre><code class="language-plaintext" data-out="ejemplos/patrones/builder"></code></pre>

Uso:
  python3 scripts/embed_codigo.py          # solo incrusta (usa los salida.txt existentes)
  python3 scripts/embed_codigo.py --run    # gofmt + go vet + go run de cada ejemplo, luego incrusta
"""
import html
import re
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
EJ = ROOT / "ejemplos"


def run_examples():
    subprocess.run(["gofmt", "-w", "."], cwd=EJ, check=True)
    subprocess.run(["go", "vet", "./..."], cwd=EJ, check=True)
    for main in sorted(EJ.rglob("main.go")):
        d = main.parent
        rel = "./" + str(d.relative_to(EJ))
        r = subprocess.run(["go", "run", rel], cwd=EJ, capture_output=True, text=True, timeout=60)
        (d / "salida.txt").write_text(r.stdout + r.stderr, encoding="utf-8")
        print(("ok   " if r.returncode == 0 else "FALLÓ"), rel)


SRC_RE = re.compile(r'(<code class="language-go" data-src="([^"]+)">)(.*?)(</code>)', re.S)
OUT_RE = re.compile(r'(<code class="language-plaintext" data-out="([^"]+)">)(.*?)(</code>)', re.S)


def embed(path: Path):
    s = path.read_text(encoding="utf-8")

    def src(m):
        code = (ROOT / m.group(2)).read_text(encoding="utf-8").rstrip("\n")
        return m.group(1) + html.escape(code, quote=False) + m.group(4)

    def out(m):
        txt = (ROOT / m.group(2) / "salida.txt").read_text(encoding="utf-8").rstrip("\n")
        txt = re.sub(r"\x1b\[[0-9;]*m", "", txt)  # quitar colores ANSI
        return m.group(1) + html.escape(txt, quote=False) + m.group(4)

    new = OUT_RE.sub(out, SRC_RE.sub(src, s))
    if new != s:
        path.write_text(new, encoding="utf-8")
        print("incrustado:", path.relative_to(ROOT))


if __name__ == "__main__":
    if "--run" in sys.argv:
        run_examples()
    for f in sorted((ROOT / "clases").glob("*.html")):
        embed(f)

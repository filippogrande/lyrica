#!/usr/bin/env python3
"""Scrive il campo dell'immagine nel front-matter di un file di contenuto.

    python3 scripts/declare-image.py <file.md> cover|band <nome-file.webp>

- `cover` aggiorna il campo `cover:` (album), `band` aggiorna `image:` (band).
- La riga si aggiorna se c'è già, si aggiunge se manca: la si lancia dopo
  aver scaricato l'immagine, invece di aprire l'editor e sbagliare il nome.
- Non tocca il corpo del file. Se manca un front-matter valido esce con un
  errore esplicito: nessun fallback silenzioso (docs/GUIDELINES.md §4).

La usa il workflow "Immagini dagli archivi" e si può lanciare anche a mano.
"""

import sys
from pathlib import Path

# Il campo dipende da cosa rappresenta l'immagine, non dal gusto: le cover
# stanno in `cover` (album.md), le foto in `image` (band.md).
FIELD_BY_KIND = {"cover": "cover", "band": "image"}


def field_for(kind: str) -> str:
    if kind not in FIELD_BY_KIND:
        raise SystemExit(f"errore: tipo sconosciuto {kind!r}: usa cover o band")
    return FIELD_BY_KIND[kind]


def split_frontmatter(text: str, path: Path) -> tuple[list[str], list[str]]:
    """Restituisce le righe del front-matter e quelle da lì in poi."""
    lines = text.splitlines()
    if not lines or lines[0].strip() != "---":
        raise SystemExit(f"errore: {path} non inizia con il front-matter YAML")
    for index in range(1, len(lines)):
        if lines[index].strip() == "---":
            return lines[1:index], lines[index:]
    raise SystemExit(f"errore: front-matter senza chiusura in {path}")


def upsert(front: list[str], field: str, value: str) -> str:
    """Aggiorna o aggiunge la riga del campo, e dice quale delle due cose ha fatto."""
    wanted = f'{field}: "{value}"'
    for index, line in enumerate(front):
        if line.split(":", 1)[0].strip() == field:
            front[index] = wanted
            return "aggiornata"
    front.append(wanted)
    return "aggiunta"


def main(argv: list[str]) -> int:
    if len(argv) != 4:
        raise SystemExit("uso: declare-image.py <file.md> cover|band <nome-file.webp>")
    path = Path(argv[1])
    field = field_for(argv[2])
    value = argv[3]
    if not path.is_file():
        raise SystemExit(f"errore: {path} non esiste")

    front, rest = split_frontmatter(path.read_text(encoding="utf-8"), path)
    outcome = upsert(front, field, value)
    path.write_text("\n".join(["---", *front, *rest]) + "\n", encoding="utf-8")
    print(f'{path}: {field}: "{value}" ({outcome})')
    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv))

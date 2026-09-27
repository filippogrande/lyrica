# CONTENT_SCHEMA — schema dei contenuti

Come sono fatti i file di contenuto. È la specifica di riferimento per la CLI (`lyrica new`) e per il validatore in CI.

## Struttura delle cartelle

```
content/
└── bands/
    └── <band-slug>/
        ├── band.md                     # anagrafica della band
        └── <album-slug>/
            ├── album.md                # anagrafica dell'album + tracklist
            └── tracks/
                ├── <track-slug>.md     # testo originale + traduzioni
                └── ...
covers/
└── <album-slug>.webp                   # 600x600, quadrata
```

Lo slug è la chiave di tutto: **cartella, URL e nome della cover devono coincidere**.

## Regole generali

- Ogni file di contenuto = **front-matter YAML** + corpo Markdown.
- Le lingue usano codici **BCP-47** a due lettere dove esiste (`it`, `en`, `de`, `ru`, `ja`).
- Gli slug sono **minuscoli, ASCII, con trattini**; se due band omonime collidono, si disambigua (`nirvana-us`, `nirvana-uk`).
- Le date sono `YYYY-MM-DD` in ora locale italiana.
- Nessun campo opzionale viene inventato: se manca, il validatore decide se è errore o se si usa un default.

## `band.md`

```yaml
---
name: "Tocotronic"
slug: "tocotronic"
country: "DE"
tags: ["indie", "hamburger schule"]
original_langs: ["de"]
formed_year: 1993
members: ["Dirk von Lowtzow", "Jan Müller", "Arne Zank"]
description: >
  Band di Amburgo, una delle voci principali della Hamburger Schule.
---
```

| Campo | Obbl. | Note |
|---|---|---|
| `name` | sì | come si legge in pagina |
| `slug` | sì | deve coincidere con il nome della cartella |
| `country` | no | sigla ISO a 2 lettere |
| `tags` | no | **tag liberi** di genere; generano le pagine `/tag/{tag}` |
| `original_langs` | sì | una o più lingue originali della band |
| `formed_year` | no | numero |
| `members` | no | lista di stringhe |
| `description` | no | Markdown breve, mostrato nella pagina band |

## `album.md`

```yaml
---
title: "Digital ist besser"
slug: "digital-ist-besser"
year: 1995
cover: "digital-ist-besser.webp"   # opzionale
tracks:
  - slug: "digital-ist-besser"
    title: "Digital ist besser"
  - slug: "ich-moechte-dich"
    title: "Ich möchte dich"
  - slug: "instrumental-01"
    title: "(strumentale)"
    instrumental: true
---
```

- `tracks` definisce **l'ordine** della tracklist così com'è scritta.
- `instrumental: true` segna il brano strumentale: sta in tracklist ma **non ha pagina**.
- `cover` è il **nome del file** in `covers/`; se assente, il layout va senza immagine.

## `tracks/<slug>.md`

```yaml
---
title: "Digital ist besser"
slug: "digital-ist-besser"
added_date: 2026-09-27
featured: false
instrumental: false
original_langs: ["de"]
singers: ["Dirk von Lowtzow"]      # opzionale: solo se il brano ha più voci
blocks:
  - lang: de
    role: original
    stanzas:
      - singer: "Dirk von Lowtzow"   # opzionale, solo nei brani multi-voce
        lines:
          - "Ich bin ein Digital Native"
          - "und du bist auch einer"
  - lang: it
    role: translation
    translator: "Filippo"
    stanzas:
      - lines:
          - "Sono un nativo digitale"
          - "e anche tu lo sei"
---

Eventuali note redazionali **fuori dal testo**: non si stampano nel corpo del testo.
```

### Campi del brano

| Campo | Obbl. | Note |
|---|---|---|
| `title` | sì | titolo mostrato |
| `slug` | sì | coincide col nome del file |
| `added_date` | sì | ordina i "recenti" e il feed RSS |
| `featured` | no | `true` = può entrare nella sezione "in evidenza" |
| `instrumental` | no | `true` = nessuna pagina |
| `original_langs` | sì | una o più |
| `singers` | no | elenco delle voci del brano |
| `blocks` | sì | almeno un blocco `role: original` |

### Campi di un blocco

| Campo | Obbl. | Note |
|---|---|---|
| `lang` | sì | BCP-47 |
| `role` | sì | `original` \| `translation` |
| `translator` | no | chi ha tradotto (mostrato in pagina) |
| `stanzas` | sì | lista ordinata di strofe |
| `live` | no | `true` = versione live/acustica: resta nello stesso brano, mostrato con nota |

### Campi di una strofa

| Campo | Obbl. | Note |
|---|---|---|
| `singer` | no | mostrato **sopra la strofa**; nei brani multi-voce va su **ogni** strofa |
| `lines` | sì | una riga di testo per elemento, nessuna riga vuota |

Il tipo di strofa (verse/chorus/bridge) **non esiste** nello schema: è stato deciso che non si stampa.

## Regole di validazione (applicate in CI)

Un contenuto è **invalido** se:

1. `slug` del front-matter ≠ nome del file o della cartella.
2. Due entità della stessa collezione hanno lo **stesso slug**.
3. `blocks` non contiene almeno un `role: original`, o contiene una `lang` duplicata.
4. Un blocco `role: translation` è **incompleto** (una strofa ha meno righe dell'originale, o manca una strofa).
5. `added_date` manca o non è una data valida.
6. `cover` punta a un file inesistente in `covers/`.
7. `original_langs` contiene una lingua che non compare in nessun blocco.
8. Un `title` in `tracks` di `album.md` non ha il file corrispondente in `tracks/` (o viceversa).
9. Un link interno a un brano/band inesistente.

La regola 4 è la più importante: **una traduzione a metà non si pubblica** (D29).

## Esempio minimo completo

Band → album con un brano tradotto, senza cover e con un brano strumentale: è il contenuto da usare per il primo test della pipeline in FASE 2.

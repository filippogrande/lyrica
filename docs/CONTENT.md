# CONTENT — schema dei contenuti e come si aggiungono

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
├── <album-slug>.webp                   # cover dell'album, 600x600
└── <band-slug>.webp                    # foto della band, 600x600 (opzionale)
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
image: "tocotronic.webp"            # opzionale: foto della band in covers/
description: >
  Band di Amburgo, una delle voci principali della Hamburger Schule.
---
```

| Campo | Obbl. | Note |
|---|---|---|
| `name` | sì | come si legge in pagina |
| `slug` | sì | deve coincidere con il nome della cartella |
| `country` | no | sigla ISO a 2 lettere |
| `tags` | no | **tag liberi** di genere; generano le pagine `/it/tag/{tag}` |
| `original_langs` | sì | una o più lingue originali della band |
| `formed_year` | no | numero |
| `members` | no | lista di stringhe |
| `image` | no | **nome del file** in `covers/` con la foto della band; se manca, la card usa il segnaposto con l'iniziale |
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
- `cover` è il **nome del file** in `covers/`; se assente, la pagina album va senza immagine.

## Immagini (`covers/`)

| Uso | File | Dove si dichiara |
|---|---|---|
| Cover di un album | `covers/<album-slug>.webp` | `cover` in `album.md` |
| Foto di una band | `covers/<band-slug>.webp` | `image` in `band.md` |

- **600x600, quadrata, webp**: una sola dimensione per tutte le immagini, così il formato non si sceglie ogni volta.
- Un'immagine **dichiarata e assente è un errore** (regola 6 per le cover, regola 10 per le band): il sito non deve avere link rotti.
- Un'immagine **non dichiarata non è un errore**: la pagina album va senza immagine, le card della home usano il **segnaposto con l'iniziale** (una lettera su fondo del tema, non un'immagine finta).

### Da dove arrivano le immagini

Si prendono da **archivi pubblici**, per MBID, non dal primo risultato di una ricerca per immagini.

| Cosa | Fonte | Identificativo | Chiave API |
|---|---|---|---|
| Cover di un album | [Cover Art Archive](https://coverartarchive.org/) (MusicBrainz + Internet Archive) | release-group MBID | no |
| Foto di una band | [fanart.tv](https://fanart.tv/) | artist MBID | sì, personale |

Il **MBID** si legge nell'URL della pagina MusicBrainz: `musicbrainz.org/release-group/<mbid>` per un album, `musicbrainz.org/artist/<mbid>` per una band.

**Le immagini restano nel repo**: si scaricano una volta e il sito le serve da `covers/`. Nessuna chiamata alle API a runtime, nessuna a ogni build — l'API si interroga solo quando aggiungi una band o un album (D82). Il footer attribuisce le immagini con i link ai due archivi (D84): gli archivi non sono una fonte di licenza, le immagini restano **dei rispettivi proprietari** (vale il disclaimer in footer, D07).

#### Dal browser, un click per immagine (consigliato)

1. GitHub → **Actions** → **Immagini dagli archivi** → **Run workflow**.
2. Tre campi: `tipo` (`cover` o `band`), `mbid`, `slug`. In più, se vuoi, il **file di contenuto** da aggiornare (`content/bands/<band>/band.md` o `.../<album>/album.md`).
3. Il workflow scarica l'immagine, la mette in `covers/<slug>.webp` a 600x600, scrive il campo nel file se gliel'hai indicato e apre da sé una PR `immagini/<slug>`.

Resta da fare solo: guardare l'immagine e mergiare. La chiave di fanart.tv va nei **secret del repo** (`Settings → Secrets and variables → Actions`): `FANART_API_KEY` (progetto) o `FANART_CLIENT_KEY` (personale). Serve solo al workflow (D85).

#### Dal proprio computer (stessi passi, senza browser)

```
# dalla root del repo
bash scripts/fetch-covers.sh cover <release-group-mbid> <album-slug>
bash scripts/fetch-covers.sh band  <artist-mbid>        <band-slug>

# e, per scrivere il campo senza aprire l'editor:
python3 scripts/declare-image.py content/bands/<band>/band.md band <band-slug>.webp
```

- Lo script scarica, converte in **600x600 webp** (ritaglio centrale) e **verifica il risultato**; se il file esiste già si ferma, a meno di `--force`. Non tocca i file di contenuto: stampa la riga da aggiungere (oppure la scrive `declare-image.py`).
- Per le band stampa il **nome dell'artista** trovato: serve a confermare di aver preso il MBID giusto. Le foto delle band **non esistono su MusicBrainz**: le indicizza fanart.tv (di solito `artistthumb`).
- Chiavi sul proprio computer: `FANART_API_KEY` / `FANART_CLIENT_KEY`, oppure una riga in `~/.fanart_api_key` / `~/.fanart_client_key`. Basta una delle due. **Mai nel repo**, mai stampate, mai nell'URL.
- Serve **ImageMagick** (`brew install imagemagick`), `curl` e `python3`: se manca qualcosa, lo script esce con un errore esplicito e non scrive niente.
- Si scarica **un'immagine per volta**, senza script che martellano l'API.

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
| `featured` | no | `true` = entra nella corsia "In evidenza" della home |
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
9. Un link interno a un band/album/brano inesistente.
10. `image` di una band punta a un file inesistente in `covers/` (o a una cartella invece che a un file).

La regola 4 è la più importante: **una traduzione a metà non si pubblica**.

### Warning (build verde, da sistemare o giustificare in PR)

- brano senza nessuna traduzione (da segnare "solo originale");
- album senza cover;
- **band senza immagine**: la card usa il segnaposto con l'iniziale (non è un errore: una band può restare tipografica);
- tag usato una volta sola (possibile refuso: `metal` vs `metalcore`);
- descrizione band mancante;
- chiavi di locale mancanti rispetto a `it.yaml`.

## Come si aggiunge un contenuto

Il flusso normale: si generano i file con la CLI e si riempiono. Non si scrive a mano partendo da zero.

### 1. Nuova band

```
lyrica new band "Nome Band"
```

Crea `content/bands/<slug>/band.md` con il front-matter precompilato. Da compilare: `tags` (liberi), `original_langs`, `description` (1-3 righe). La **foto della band** è opzionale: si scarica con il workflow o con lo script (vedi sopra), finisce in `covers/<band-slug>.webp` e si dichiara con `image`.

### 2. Nuovo album

```
lyrica new album "Nome Band" "Titolo Album" --year 1995
```

Crea `album.md` con la tracklist vuota da riempire **nell'ordine dell'album**. I brani strumentali si segnano subito con `instrumental: true`. La cover si scarica in `covers/<album-slug>.webp`, **600x600 quadrata**: se non c'è, la pagina album va senza immagine e le card della home mostrano il **segnaposto con l'iniziale** dell'album — una lettera su fondo del tema, non un'immagine inventata (D78).

### 3. Nuovo brano

```
lyrica new brano "Nome Band" "album-slug" "Titolo Brano"
```

Crea `tracks/<track-slug>.md` con lo scheletro dei blocchi lingua. Poi, a mano:

- scrivi il testo originale **strofa per strofa**, un verso per riga;
- se il brano ha più voci, aggiungi `singer:` a **ogni** strofa;
- scrivi la traduzione **completa**: strofe della stessa lunghezza dell'originale;
- `added_date` = il giorno in cui la pubblichi;
- `featured: true` solo se il brano deve stare nella corsia "In evidenza" della home;
- niente annotazioni per verso.

### 4. Verifica in locale

```
lyrica build      # valida i contenuti e genera public/
lyrica serve      # serve il risultato su una porta locale
```

### 5. Pubblicazione

- Commit su un **branch** (`content/<band-slug>`).
- **Pull Request** verso `main`.
- La CI **valida i contenuti**; il merge in `main` fa deploy.

### Checklist prima di aprire la PR

- [ ] La traduzione è **completa** strofa per strofa.
- [ ] `added_date` è la data di pubblicazione reale.
- [ ] Slug, cartelle e nome della cover **coincidono**.
- [ ] Le immagini dichiarate (`cover`, `image`) **esistono** in `covers/` e sono 600x600 webp.
- [ ] I nomi dei cantanti sono solo nei brani multi-voce, e su ogni strofa.
- [ ] Nessuna annotazione per verso.
- [ ] Il contenuto si legge bene **da telefono**, non solo da desktop.
- [ ] Se un brano non ha traduzioni, è segnato come "solo originale" e non è linkato.

### Aggiungere una lingua di traduzione a un brano

Aggiungi al file un blocco con `role: translation`, `lang:`, `translator:` e le strofe **complete**. Il validatore controlla che combacino con l'originale: se non combaciano, la PR è rossa. Il selettore lingua in pagina si aggiorna da solo.

### Aggiungere una lingua all'interfaccia

Copia `locales/it.yaml` in `locales/<lang>.yaml` e traduci i **valori**, non le chiavi. Il build genera `/<lang>/...` con hreflang. Una chiave mancante è un errore in CI.

### Rinominare uno slug (cambio URL)

Rinomina cartella o file, aggiungi la **regola di redirect 301** nel file dei redirect (vecchio URL → nuovo URL) e verifica il redirect nella build locale.

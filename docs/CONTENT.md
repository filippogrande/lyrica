# CONTENT — schema dei contenuti e come si aggiungono

## Stato: due formati, uno in arrivo

Oggi i contenuti in repo sono nel **formato a file unico** (`tracks/<slug>.md`): un file per brano con
front-matter e blocchi di testo inline. È questo il formato descritto nella sezione
`tracks/<slug>.md` più sotto, ed è quello che il loader legge oggi per tutto il catalogo.

Il **formato a cartella** (`tracks/<slug>/` con `song.md` e un file per lingua) è la
destinazione: lo crea `lyrica new brano`, lo legge il loader, e `lyrica migrate` lo produce a
partire dal formato precedente. Per questo la sezione `tracks/<slug>/` è marcata come
**destinazione**: i file su disco sono ancora quasi tutti nel formato a file unico.

Il passaggio si fa **con la migrazione, non a mano**:

```
go run ./cmd/lyrica migrate [--band <slug>]          # prova: dice cosa farebbe
go run ./cmd/lyrica migrate --apply                  # scrive davvero
```

oppure da browser: **Actions → Migrazione brani (cartella per brano) → Run workflow**
(`.github/workflows/migrazione-brani.yml`), che apre una PR.

## Struttura delle cartelle

```
content/
└── bands/
    └── <band-slug>/
        ├── band.md                     # anagrafica della band
        └── <album-slug>/
            ├── album.md                # anagrafica dell'album + tracklist
            └── tracks/
                ├── <track-slug>.md     # FORMATO ATTUALE: testo originale + traduzioni inline
                └── <track-slug>/       # DESTINAZIONE: una cartella per brano
                    ├── song.md
                    ├── original-<lang>.md
                    └── <lang>.md
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

### Lingua di un brano (una o più lingue originali)

Un brano ha **un solo blocco `original`**, nella lingua **dominante** — quella in cui è cantata la maggior parte del testo — e quel blocco contiene il testo **come è cantato**, anche quando mescola più lingue. `original_langs` dichiara **tutte** le lingue originali, **in ordine**: la prima è la dominante, cioè la lingua del blocco `original` (D28, D96, D97).

- **Parole o frasi isolate** in altre lingue **non** rendono il brano multilingue: es. *ADIEU* (Rammstein) è tedesco anche se contiene parole non-tedesche. Una lingua sola, nessuna traduzione extra.
- Se una **seconda lingua copre una parte sostanziale** del testo (es. *Kinglayer* di Bring Me the Horizon / Babymetal con 2-3 frasi in giapponese) è un caso **al limite**: si decide **caso per caso**.
- L'array accetta **N lingue**, non solo due: la soglia è umana, la decide chi scrive il contenuto.
- Le traduzioni sono **una per lingua d'arrivo**, più **una versione completa per ogni lingua dichiarata in `original_langs`**: in un brano `["it", "fr", "de"]` si possono avere una `translation it`, una `translation fr` e una `translation de`, ciascuna con **tutte** le strofe tradotte in quella lingua. Serve anche ai **dialetti e alle varianti** (la versione in lingua regionale accanto a quella standard).
- La stessa lingua su **ruoli diversi** non è un doppione: `original de` e `translation de` convivono (regola 3, D96).

È un criterio umano, non automatico: il validatore non misura percentuali, decide chi scrive il contenuto.

Due cose restano **regole di scrittura**, verificate in review e non dalla CI: il blocco `original` è **uno solo**, e la sua lingua è la **prima** di `original_langs`. Il validatore controlla la **coppia** `role` + `lang` (regola 3), non l'ordine dell'array.

### Varianti di scrittura (BCP-47 script subtag)

Per le lingue con più alfabeti (giapponese kanji/romaji, cinese semplificato/tradizionale, serbo cirillico/latino, coreano hangul/latino), il codice BCP-47 accetta un **subtag di scrittura** dopo il trattino (D99):

| Codice | Variante | Esempio |
|--------|----------|---------|
| `ja-Jpan` | Giapponese in Kanji/Kana | 日本語 |
| `ja-Latn` | Giapponese in Romaji | Nihongo |
| `zh-Hans` | Cinese semplificato | 简体中文 |
| `zh-Hant` | Cinese tradizionale | 繁體中文 |
| `sr-Cyrl` | Serbo cirillico | српски |
| `sr-Latn` | Serbo latino | srpski |

La lingua base determina la bandiera (`ja` → 🇯🇵, `zh` → 🇨🇳); il subtag mostra il nome descrittivo nell'etichetta (`🇯🇵 JA Kanji`, `🇨🇳 ZH Semplificato`).

Entrambe le varianti di scrittura sono **originali** — stesso testo, alfabeto diverso — quindi usano entrambe `role: original` e stanno **entrambe a sinistra** nel selettore dell'originale:

```yaml
original_langs: ["ja-Jpan", "ja-Latn"]
blocks:
  - role: original
    lang: ja-Jpan      # Kanji/Kana (dominante, prima in original_langs)
    stanzas:
      - lines: ["日本語のテキスト"]
  - role: original
    lang: ja-Latn      # Romaji = stesso testo, alfabeto latino
    stanzas:
      - lines: ["Nihongo no tekisuto"]
  - role: translation
    lang: it           # Traduzione italiana (lato destro)
    stanzas:
      - lines: ["Testo in giapponese"]
```

Il lato sinistro mostra il dropdown con `🇯🇵 JA Kanji` · `🇯🇵 JA Romaji`; il lato destro ha `🇮🇹 IT`. (D99)

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
| `name` | sì | come si legge in pagina; è anche il nome con cui si cerca la foto della band |
| `slug` | sì | deve coincidere con il nome della cartella |
| `country` | no | sigla ISO a 2 lettere |
| `tags` | no | **tag liberi** di genere; generano le pagine `/it/tag/{tag}` (**previste, non esistenti**) |
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
  - slug: "digital-ist-besser"     # ha il file in tracks/: in pagina è un link
    title: "Digital ist besser"
  - slug: "ich-moechte-dich"
    title: "Ich möchte dich"
  - slug: "spieluhr"               # testo non ancora scritto
    title: "Spieluhr"
    status: pending
  - slug: "instrumental-01"        # strumentale senza file
    title: "(strumentale)"
    status: instrumental
---
```

- `tracks` definisce **l'ordine** della tracklist così com'è scritta: è l'ordine del disco e quello che si stampa.
- `status` dichiara una traccia **senza file**: `pending` (testo non ancora scritto) o `instrumental`.
- `cover` è il **nome del file** in `covers/`; se assente, la pagina album va senza immagine.
- `title` e `year` sono anche i dati con cui si cerca la cover: vanno scritti come sono sull'album, non come viene più comodo.

### Tracklist completa: anche i brani che non ci sono ancora

La tracklist è **il disco intero**, non l'elenco dei brani già scritti: si mettono **tutte** le tracce, nell'ordine dell'album, e per quelle che non hanno il file si dichiara **perché** mancano.

| `status` | Significato | In pagina |
|---|---|---|
| assente | il brano ha il suo file in `tracks/` | **link** se ha almeno una traduzione, altrimenti "Solo originale" (o "Strumentale") |
| `pending` | il testo non è ancora scritto, il file **non esiste** | numero e titolo, **senza link**, con "In arrivo" |
| `instrumental` | brano strumentale, senza file | numero e titolo, **senza link**, con "Strumentale" |

- Una voce **con `status` non può avere il file**: sarebbe una voce che dice di non avere il testo mentre il testo c'è. Il validatore la segnala come errore (si toglie lo `status`).
- `pending` e `instrumental` sono gli unici valori ammessi: un valore diverso è un errore, non un'etichetta libera.
- La **pagina album** si genera quando c'è almeno una traduzione (D35), ma la tracklist mostra **tutte** le tracce dichiarate, comprese quelle in arrivo.
- La **pagina brano** esiste solo per i brani che hanno il testo (e almeno una traduzione): le voci in arrivo non sono link e non sono cliccabili.
- Chi vuole la traduzione di un brano "in arrivo" potrà chiederla dal form di segnalazione (**previsto, non esistente**): la call to action in pagina arriva **insieme al form** (FASE 6, D90).

## Immagini (`covers/`)

| Uso | File | Dove si dichiara |
|---|---|---|
| Cover di un album | `covers/<album-slug>.webp` | `cover` in `album.md` |
| Foto di una band | `covers/<band-slug>.webp` | `image` in `band.md` |

- **600x600, quadrata, webp**: una sola dimensione per tutte le immagini, così il formato non si sceglie ogni volta.
- Un'immagine **dichiarata e assente è un errore** (regola 6 per le cover, regola 10 per le band): il sito non deve avere link rotti.
- Un'immagine **non dichiarata non è un errore**: la pagina album va senza immagine, le card della home usano il **segnaposto con l'iniziale** (una lettera su fondo del tema, non un'immagine finta).

### Da dove arrivano le immagini

Si prendono da **archivi pubblici**, per identificativo, non dal primo risultato di una ricerca per immagini.

| Cosa | Fonte | Identificativo | Chiave API |
|---|---|---|---|
| Cover di un album | [Cover Art Archive](https://coverartarchive.org/) (MusicBrainz + Internet Archive) | release-group, cercato da titolo + anno | no |
| Foto di una band | [fanart.tv](https://fanart.tv/) | artist, cercato dal nome della band | sì, personale |

**Le immagini restano nel repo**: si scaricano una volta e il sito le serve da `covers/`. Nessuna chiamata alle API a runtime, nessuna a ogni build — l'API si interroga solo quando si aggiunge una band o un album (D82). Il footer attribuisce le immagini con i link ai due archivi (D84): gli archivi non sono una fonte di licenza, le immagini restano **dei rispettivi proprietari** (vale il disclaimer in footer, D07).

#### Dal browser, un click (il modo normale)

1. GitHub → **Actions** → **Immagini mancanti** → **Run workflow**. Il campo `band` si può lasciare vuoto: così guarda **tutte** le band.
2. Il workflow cerca da solo quello che manca — band senza `image`, album senza `cover` (o con un `cover` che punta a un file che non c'è) — e lo scarica: la band si riconosce dal `name` in `band.md`, l'album dal `title` e dall'`year` in `album.md`. **Nessun MBID, nessuno slug da digitare.**
3. Scrive i campi mancanti nel front-matter e apre **una PR sola** con tutte le immagini.

Nel **riepilogo del run** (`Actions` → il run → `Summary`) c'è il report: per ogni immagine se è stata scaricata o perché è stata saltata (già presente, nessuna corrispondenza su MusicBrainz, archivio che non ce l'ha).

Prima di mergiare: **guardare le immagini**. La corrispondenza è automatica e prudente — nome esatto e punteggio alto per la band; per l'album solo `Album` con anno a un anno di distanza e punteggio alto — ma un titolo omonimo (live, compilation, edizione estera) può portare alla release sbagliata.

Quello che l'archivio non ha **non è un errore**: l'album resta senza cover (D36) e la card della home usa il segnaposto con l'iniziale.

#### Una immagine sola, quando serve la certezza

Se la cover scelta è quella sbagliata, si corregge indicando il MBID a mano: **Actions → Immagini dagli archivi** (tipo, MBID, slug, e il file da aggiornare). Il MBID si legge nell'URL della pagina MusicBrainz: `musicbrainz.org/release-group/<mbid>` per un album, `musicbrainz.org/artist/<mbid>` per una band.

La chiave di fanart.tv va nei **secret del repo** (`Settings → Secrets and variables → Actions`): `FANART_API_KEY` (progetto) o `FANART_CLIENT_KEY` (personale). Serve solo alle foto delle band (D85).

#### Dal proprio computer (gli stessi passi, senza browser)

```
bash scripts/fetch-images.sh              # riempie tutte le immagini mancanti
bash scripts/fetch-images.sh rammstein    # solo una band

bash scripts/fetch-covers.sh cover <release-group-mbid> <album-slug>
bash scripts/fetch-covers.sh band  <artist-mbid>        <band-slug>

python3 scripts/declare-image.py content/bands/<band>/band.md band <band-slug>.webp
```

- `fetch-images.sh` è quello che fa il lavoro da solo: per ogni band e ogni album legge i dati già scritti nei file di contenuto, cerca l'identificativo, scarica e scrive il campo. Salta ciò che è a posto e stampa il report (anche nel riepilogo del run, se `GITHUB_STEP_SUMMARY` è impostato).
- `fetch-covers.sh` scarica **una** immagine dal MBID, converte in **600x600 webp** con ritaglio centrale e **verifica le dimensioni**; `--force` per riscrivere, `--allow-missing` per uscire con 3 invece che con errore quando l'archivio non ce l'ha. Non tocca i file di contenuto.
- `declare-image.py` scrive `cover:` / `image:` nel front-matter (aggiorna la riga se c'è, la aggiunge se manca, errore esplicito se non c'è front-matter).
- Le foto delle band **non esistono su MusicBrainz**: le indicizza fanart.tv (di solito `artistthumb`, si prende quella con più like).
- Chiavi sul proprio computer: `FANART_API_KEY` / `FANART_CLIENT_KEY`, oppure una riga in `~/.fanart_api_key` / `~/.fanart_client_key`. Basta una delle due. **Mai nel repo**, mai stampate, mai nell'URL.
- Serve **ImageMagick** (`brew install imagemagick`), `curl` e `python3`: se manca qualcosa, gli script escono con un errore esplicito e non scrivono nulla.
- MusicBrainz accetta **una richiesta al secondo**: gli script aspettano da soli fra una ricerca e l'altra. Duecento album richiedono qualche minuto, non duecento click.

## `tracks/<slug>.md` (formato attuale)

È il formato dei file **che ci sono oggi** in `content/`: un file per brano, front-matter con
l'anagrafica **e i blocchi di testo inline**, corpo dopo i `---` per le note redazionali.

```yaml
---
title: "Digital ist besser"
slug: "digital-ist-besser"
added_date: 2026-09-27
featured: false
instrumental: false
original_langs: ["de"]
singers: ["Dirk von Lowtzow"]      # opzionale: solo se il brano ha più voci
artists:                            # opzionale: collaborazioni e featuring
  - slug: "tocotronic"
    role: "primary"
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
| `original_langs` | sì | una o più, **in ordine**: la prima è la lingua dominante, quella del blocco `original` |
| `singers` | no | elenco delle voci del brano |
| `artists` | no | lista di `ArtistCredit` per collaborazioni: `slug` + `role` (`primary` \| `featuring` \| `equal`) |
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

### Artisti e collaborazioni

Il campo `artists` gestisce collaborazioni e featuring. Ogni voce è un oggetto con `slug` (deve esistere in `content/bands/`) e `role`:

| Ruolo | Significato | Separatore UI |
|---|---|---|
| `primary` | artista principale (ospita il file) | `·` |
| `featuring` | artista in featuring | `feat.` |
| `equal` | collaborazione paritaria (50/50) | `×` |

```yaml
# Ratatata — collaborazione paritaria
title: "RATATATA"
artists:
  - slug: "electric-callboy"
    role: "primary"
  - slug: "babymetal"
    role: "equal"

# From Me To U — featuring
title: "From Me To U"
artists:
  - slug: "babymetal"
    role: "primary"
  - slug: "poppy"
    role: "featuring"
```

- Il brano è **linkato sotto ogni artista** nella pagina band (discoverability bidirezionale).
- Il titolo resta **pulito** (nessun "feat." nel titolo): il ruolo si legge dalla lista artisti.
- Se `artists` manca, il validatore emette un **warning** (non errore): il brano è collegato solo alla band dell'album.

#### Un brano in due album (collaborazione paritaria cross-album)

Quando lo **stesso brano** compare in **due album di due artisti diversi** (es. *RATATATA*, in *TANZNEID* degli Electric Callboy e in *METAL FORTH* delle Babymetal) **esiste un solo file**, nell'album dell'artista `primary`; il brano si dichiara con `artists` (`primary` + `equal`, D100, D107) e **compare linkato nella pagina di entrambe le band**.

- **Un file solo**: il testo si scrive una volta, nell'album dell'artista `primary` (qui Electric Callboy). Non si duplica mai.
- **L'altro album** tiene la voce di tracklist **senza file**, con `status: pending` (la tracklist è il disco intero, D90): lì `pending` non significa "testo non scritto" ma "questo testo vive nell'album dell'altro artista".
- **Entrambe le pagine band** linkano il brano (`track.HasArtist`), quindi la traduzione è raggiungibile anche dall'artista che non ospita il file.

## `tracks/<slug>/` (destinazione: una cartella per brano)

È il formato che `lyrica new brano` crea e che `lyrica migrate` produce. Il loader lo legge
già (`internal/content/trackfolder.go`), ma i contenuti in repo sono ancora quasi tutti nel
formato a file unico.

```
tracks/<slug>/
  song.md             anagrafica del brano (stessi campi del front-matter sopra, senza blocks)
  original-<lang>.md un file per ogni blocco role: original
  <lang>.md           un file per ogni blocco role: translation
```

- Il nome del file **è** la chiave del blocco: `original-<lang>.md` è un originale, `<lang>.md`
  è una traduzione. Un brano che ha la stessa lingua come originale e come versione completa
  (D96) non collide: `original-de.md` e `de.md` sono due file diversi, non un doppione.
- `song.md` contiene l'anagrafica **senza** `blocks` e il corpo dopo i `---` come note redazionali.
- Ogni file di blocco è **autocontenuto**: si legge e si valida da solo, e una traduzione nuova è
  un file nuovo, non una modifica a un file che cresce.
- Il loader verifica che **nome del file e front-matter dicano la stessa cosa**: un file chiamato
  `it.md` con `lang: de` è un errore, non una sorpresa.
- I blocchi si leggono **in ordine di nome**, così l'elenco delle lingue in pagina non dipende
  dall'ordine del filesystem.

## Regole di validazione (applicate in CI)

Un contenuto è **invalido** se:

1. `slug` del front-matter ≠ nome del file o della cartella.
2. Due entità della stessa collezione hanno lo **stesso slug**.
3. `blocks` non contiene almeno un `role: original`, o contiene **due blocchi con la stessa coppia `role` + `lang`** (la stessa lingua può comparire una volta come originale e una come traduzione: D96, D97).
4. Un blocco `role: translation` è **incompleta** (una strofa ha meno righe dell'originale, o manca una strofa).
5. `added_date` manca o non è una data valida.
6. `cover` punta a un file inesistente in `covers/`.
7. `original_langs` contiene una lingua che non compare in nessun blocco.
8. Una voce di `tracks` in `album.md` **senza `status`** non ha il file corrispondente in `tracks/` (o viceversa: un file assente dalla tracklist). Una voce **con** `status` che ha anche il file è un errore: o ha il testo (e lo `status` si toglie) o non ce l'ha.
9. Un link interno a un band/album/brano inesistente.
10. `image` di una band punta a un file inesistente in `covers/` (o a una cartella invece che a un file).
11. `status` di una voce di `tracks` con un valore diverso da `pending` e `instrumental`.
12. `artists` contiene uno slug che non esiste in `content/bands/`.
13. Una riga di testo **vuota** dentro una strofa: diventa una riga fantasma in pagina.

La regola 4 è la più importante: **una traduzione a metà non si pubblica**.

### Warning (build verde, da sistemare o giustificare in PR)

- brano senza nessuna traduzione (da segnare "solo originale");
- album senza cover;
- **band senza immagine**: la card usa il segnaposto con l'iniziale (non è un errore: una band può restare tipografica);
- tag usato una volta sola (possibile refuso: `metal` vs `metalcore`);
- descrizione band mancante;
- `title` o `instrumental` di un file diversi da quelli dichiarati nella tracklist di `album.md`;
- tracklist vuota, o con **tutte** le tracce "in arrivo";
- brano senza `artists` (non collegato a band esterne);
- album con una voce "in arrivo" ma nessun brano con il testo: l'album non ha nemmeno un lavoro in corso.

**Non è un warning né un errore**: una chiave di locale mancante rispetto a `it.yaml`. Il
validatore **non** confronta le chiavi fra lingue (`ARCHITECTURE.md`, "Limiti noti").

## Come si aggiunge un contenuto

Il flusso normale: si generano i file con la CLI e si riempiono. Non si scrive a mano partendo da zero.

### 1. Nuova band

```
lyrica new band "Tocotronic" --lang de
```

Crea `content/bands/tocotronic/band.md` con il front-matter precompilato (nome, slug, lingua da `--lang`). Da compilare a mano: `country`, `tags`, `formed_year`, `members`, `description` (1-3 righe). La **foto della band** è opzionale: la prende il workflow "Immagini mancanti" o `scripts/fetch-images.sh` (vedi sopra), finisce in `covers/<band-slug>.webp` e si dichiara con `image`.

### 2. Nuovo album

```
lyrica new album tocotronic "Digital ist besser" --year 1995
```

Il primo argomento è lo **slug della band** (non il nome). Crea `content/bands/tocotronic/digital-ist-besser/album.md` con la tracklist vuota da riempire **nell'ordine dell'album**: la tracklist è il disco intero, quindi si mettono tutte le tracce, segnando `status: pending` quelle di cui non c'è ancora il testo e `status: instrumental` gli strumentali. La cover la scarica il workflow (o `fetch-images.sh`), **600x600 quadrata**: se l'archivio non ce l'ha, la pagina album va senza immagine e le card della home mostrano il **segnaposto con l'iniziale** dell'album — una lettera su fondo del tema, non un'immagine inventata (D78).

### 3. Nuovo brano

```
lyrica new brano tocotronic digital-ist-besser "Ich möchte dich" --lang de
```

Crea la **cartella** `content/bands/tocotronic/digital-ist-besser/tracks/ich-moechte-dich/` con
`song.md` (anagrafica) e `original-de.md` (lo scheletro del blocco originale, una strofa vuota
da riempire), e **aggiunge la voce alla tracklist** di `album.md`. Poi, a mano:

- scrivi il testo originale **strofa per strofa**, un verso per riga;
- se il brano ha più voci, aggiungi `singer:` a **ogni** strofa;
- aggiungi un file `<lang>.md` per ogni traduzione, **completa**: strofe della stessa lunghezza dell'originale;
- `added_date` = il giorno in cui la pubblichi (la CLI lo mette già a oggi);
- `featured: true` solo se il brano deve stare nella corsia "In evidenza" della home;
- niente annotazioni per verso.

Se il brano era già dichiarato `status: pending` in `album.md`, quando scrivi il testo **togli lo `status`**: la voce torna a puntare al file (una voce con `status` e il file è un errore).

Se invece scrivi il brano nel **formato a file unico** (`tracks/<slug>.md` con i blocchi inline), vale
lo stesso schema della sezione `tracks/<slug>.md`: è il formato che `lyrica migrate` trasforma.

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
- [ ] La **tracklist è completa** (tutte le tracce del disco) e le voci senza testo hanno lo `status` giusto.
- [ ] Slug, cartelle e nome della cover **coincidono**.
- [ ] Le immagini dichiarate (`cover`, `image`) **esistono** in `covers/` e sono 600x600 webp.
- [ ] Se la cover l'ha scelta il workflow, **l'hai guardata**: titolo e anno possono combaciare con la release sbagliata.
- [ ] I nomi dei cantanti sono solo nei brani multi-voce, e su ogni strofa.
- [ ] Nessuna annotazione per verso.
- [ ] Il contenuto si legge bene **da telefono**, non solo da desktop.
- [ ] Se un brano non ha traduzioni, è segnato come "solo originale" e non è linkato.
- [ ] Se il brano è una collaborazione, `artists` è presente con i ruoli corretti.
- [ ] Se hai scritto il brano nel formato a file unico, `lyrica migrate` lo porta in cartella (o resta così: il loader legge entrambi).

### Aggiungere una lingua di traduzione a un brano

Aggiungi al brano un blocco con `role: translation`, `lang:`, `translator:` e le strofe
**complete** — nel formato a file unico dentro `blocks:`, in quello a cartella come file
`<lang>.md`. Il validatore controlla che combacino con l'originale: se non combaciano, la PR è
rossa. Il selettore lingua in pagina si aggiorna da solo.

Se la traduzione è la **versione completa** in una lingua che il brano ha già come originale — il tedesco intero di un brano bilingue — vale la stessa regola: un blocco `role: translation` con la stessa `lang` dell'originale e **tutte** le strofe allineate. Non si duplica il blocco `original` (D96). Nel formato a cartella è il caso in cui esistono sia `original-de.md` sia `de.md`.

La versione completa si può scrivere per **ogni** lingua di `original_langs`, non solo per le due di un brano bilingue (D97): un brano `["it", "fr", "de"]` può avere una `translation` completa in italiano, una in francese e una in tedesco, più le traduzioni di arrivo. Ogni blocco resta allineato all'originale, strofa per strofa (regola 4).

### Aggiungere una lingua all'interfaccia

Copia `locales/it.yaml` in `locales/<lang>.yaml` e traduci i **valori**, non le chiavi. Il build genera `/<lang>/...` con hreflang. Una lingua di traduzione senza il suo file di locale è un **errore della build** (`interfaceLangs` in `internal/build/build.go`).

Una traduzione verso una lingua che il catalogo ha **già come lingua originale** non fa nascere un'interfaccia in quella lingua (D96, confermata in D97): chi legge quella lingua legge già l'originale, quindi il pubblico non cambia. Il tedesco "completo" di un brano bilingue non chieda `locales/de.yaml`.

Nota: le **chiavi mancanti** in un locale non sono un errore di validazione. Il validatore non
confronta le chiavi fra lingue; una chiave assente si vede solo a runtime (fallback italiano, o
errore di `MustT` se manca anche in italiano).

### Rinominare uno slug (cambio URL)

Rinomina cartella o file. Il **redirect 301 non esiste**: oggi il vecchio URL dà 404, quindi il
redirect è una decisione da prendere (e da implementare) prima di rinominare. Dove sarà
implementato, la regola è: vecchio URL → nuovo URL in un file di redirect letto dal build e
verificato nella build locale.

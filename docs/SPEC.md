# SPEC — cosa fa Lyrica

## Stato di implementazione

Questo documento descrive il **progetto**, non solo il codice di oggi: FASE 0 (documentazione)
parte prima del codice, quindi alcune cose descritte qui **non esistono ancora**. Dove manca, la
riga lo dice esplicitamente. Per la lista completa di ciò che non esiste, vedi
[`REFERENCE.md`](REFERENCE.md#non-esiste-ancora-dove-non-cercare).

| Esiste oggi | Previsto, non implementato |
|---|---|
| Home, elenco band, pagina band, pagina album, pagina brano, 404 | pagine `/it/tag/{tag}`, `/it/segnala`, `/it/contatti`, `/it/legali` |
| Feed RSS (una per lingua), sitemap, robots, contatori pubblici, ricerca full-text (indice JSON + JS vanilla, D110) | — |
| Redirect della radice verso la lingua negoziata, `/healthz` | form segnalazione/contatti (`/it/api/segnala`, `/it/api/contatti`) |
| Tema chiaro/scuro, dimensione testo, filtro alfabetico, selettore lingua | ads, analytics (Umami), PWA, redirect 301 per slug rinominati |

## In una frase

Un archivio pubblico di **traduzioni di testi musicali**, consultabile **senza account**, pensato per leggere un testo originale e la sua traduzione **a fronte**.

## Chi lo usa

1. Chi ascolta una band straniera e vuole capire cosa dice il testo.
2. Chi ricorda **una frase** di una canzone e cerca di ritrovare il brano (ricerca per verso).
3. Chi segue il sito e vuole sapere quando esce una traduzione nuova (feed RSS).

Nessuno di questi utenti ha un account, né lo avrà.

## Mappa delle pagine

| URL | Pagina | Contenuto principale |
|---|---|---|
| `/it/` | Home | 3-5 "in evidenza", contatori, elenco recenti |
| `/it/bands` | Elenco band | elenco alfabetico + filtro per tag |
| `/it/tag/{tag}` | Band per genere | **prevista, non esistente** |
| `/it/band/{band}` | Band | descrizione, tag, album (solo con traduzioni) |
| `/it/band/{band}/album/{album}` | Album | anno, cover, tracklist |
| `/it/band/{band}/album/{album}/brano/{brano}` | **Brano** | vista a fronte, selettore lingua, "chiedi altre canzoni" |
| `/it/segnala` | Segnala un brano | **prevista (FASE 6), non esistente** |
| `/it/contatti` | Contatti | **prevista (FASE 6), non esistente** |
| `/it/legali` | Note legali | **prevista (FASE 5), non esistente** |
| `/it/404` | Non trovato | messaggio + ricerca (**la ricerca è prevista: la 404 non ha ancora il campo**) |

Endpoint non-pagina: la ricerca **non ha endpoint** — il dropdown filtra l'indice JSON generato a build time per lingua (`internal/build/search.go`, D110). Restano previsti solo `POST /it/api/segnala` e `POST /it/api/contatti` (form, FASE 6) — **non esistono oggi**.

File generati dal build, senza prefisso di lingua: `/sitemap.xml`, `/robots.txt` — **esistono**.
File generati dal build, **uno per lingua dell'interfaccia**: `/it/rss.xml`, `/en/rss.xml` — **non** c'è un feed alla root e non ci sono feed per band (D102). Vale per `/it/404.html` e per ogni pagina.

## Dettaglio delle pagine

### Home
- **In evidenza**: da 3 a 5 elementi scelti a mano, in cima.
- **Contatori**: brani tradotti / band / lingue, stessa riga del footer (D104, D105).
- **Recenti**: elenco ordinato per data di aggiunta, il più nuovo prima.
- La ricerca sta nell'header: campo + dropdown live (D110).

### Elenco band
- Ordine **alfabetico**.
- **Filtro per tag** di genere — **previsto: non esiste la pagina `/it/tag/{tag}`**.
- Ogni voce: nome band + tag.

### Band
- Nome, descrizione breve, tag, paese, lingue originali, anno di formazione.
- Elenco dei **soli album con almeno una traduzione pubblicata**.
- Album senza cover: nessuna immagine, layout senza copertina.

### Album
- Titolo, anno, cover (600x600 webp) se presente.
- **Tracklist completa** in ordine.
- Brani strumentali: in tracklist con la nota "strumentale", **senza link**.
- Brani senza alcuna traduzione: in tracklist con la nota "solo originale", **senza link**.
- Brani tradotti: link alla pagina + etichetta delle lingue disponibili.

### Brano (la pagina centrale)

Struttura verticale, dall'alto in basso:

```
[ banner ads ]                        (previsto, non esistente)
  titolo brano + band + lingue disponibili
  selettore lingua
  ┌─────┬──────────────────┬──────────────────┬─────┐
  │ ads │  ORIGINALE       │  TRADUZIONE      │ ads │   (colonne laterali: previste, non esistenti)
  └─────┴──────────────────┴──────────────────┴─────┘
[ chiedi altre canzoni → form segnalazione ]   (link: previsto FASE 6)
[ banner ads ]                        (previsto, non esistente)
```

- **Sinistra**: testo originale, diviso in strofe. **Destra**: traduzione nella lingua selezionata.
- Se il brano ha **più voci**, sopra ogni strofa compare **il nome del cantante**; nessuna etichetta strutturale (mai "ritornello", "bridge").
- **Selettore lingua**: solo le lingue realmente presenti per quel brano (originali multipli inclusi).
- L'originale è sempre disponibile: non esistono stati vuoti.
- **Nessuna annotazione per verso**: il testo resta pulito.
- **"Chiedi altre canzoni"**: fascia in fondo al testo con il link al form di segnalazione — **prevista: il form non esiste e la fascia non è implementata**.
- Gli **slot pubblicitari** sono **previsti, non implementati**: oggi nessuno slot è nel DOM e il testo prende tutta la larghezza.

### Segnala / Contatti
- **Previste (FASE 6), non esistenti.** Quando esisteranno: form semplici, senza account, con i dati **inoltrati al bot Telegram** dell'autore.
- Protezione: rate-limit, honeypot, informativa privacy inline.
- Al lancio non parte nessuna email.

## Funzionalità trasversali

- **Tema** chiaro/scuro: automatico dal sistema + toggle manuale.
- **Dimensione testo** A- / A+ (accessibilità).
- **Lingua dell'interfaccia**: negoziata dal browser, fallback italiano; le lingue disponibili seguono quelle delle traduzioni presenti.
- **Breadcrumb** su tutte le pagine interne.
- **Contatori** pubblici: in homepage e nel footer, stessa riga e stessi numeri.
- **Feed RSS** per ogni lingua dell'interfaccia, linkato dal `<head>` (rel alternate) e dal footer — **esiste**.
- **`sitemap.xml` e `robots.txt`** generati dal build, per l'indicizzazione (D101) — **esistono**.
- **Header** sticky su mobile con nav compressa.
- **Ricerca full-text**: **esistente** — indice JSON generato a build time + JS vanilla nel dropdown (D110).

## Cosa il sito NON fa (non è un TODO, è il progetto)

- Nessun account, nessun login, nessun profilo, **nessun pannello admin**.
- Nessun commento, nessun upload, nessun contenuto generato dagli utenti.
- Nessun "preferito" o ricerca salvata.
- Nessuna newsletter, nessuna email automatica (al lancio).
- Nessuna API pubblica documentata: gli unici endpoint dinamici previsti sono i form di segnalazione e contatto; la ricerca non usa endpoint (D110).
- Nessun dato personale conservato sul server (i form inoltrano e dimenticano).

## Casi limite e come si comportano

| Caso | Comportamento |
|---|---|
| Album senza cover | layout senza immagine, nessun placeholder |
| Brano strumentale | in tracklist con nota, nessuna pagina |
| Brano senza traduzioni | in tracklist con nota "solo originale", nessuna pagina |
| Brano con più lingue originali | sono opzioni del selettore lingua come le traduzioni |
| Traduzione a metà | non si pubblica, non si costruisce la pagina |
| Lingua UI non tradotta | fallback italiano |
| Slug cambiato | **previsto**: redirect 301 dal file dei redirect. **Oggi non esiste**: lo slug va cambiato a mano e il vecchio URL dà 404 |
| Band omonima | slug disambiguato (`nirvana-us`) |
| Ads spente | **non applicabile**: non esiste il codice degli slot; il testo prende tutta la larghezza |
| Ricerca senza risultati | **esistente**: messaggio "Nessun risultato" nel dropdown (D110) |

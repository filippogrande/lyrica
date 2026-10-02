# SPEC — cosa fa Lyrica

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
| `/it/tag/{tag}` | Band per genere | stesso elenco della pagina band, filtrato |
| `/it/band/{band}` | Band | descrizione, tag, album (solo con traduzioni) |
| `/it/band/{band}/album/{album}` | Album | anno, cover, tracklist |
| `/it/band/{band}/album/{album}/brano/{brano}` | **Brano** | vista a fronte, selettore lingua, "chiedi altre canzoni" |
| `/it/segnala` | Segnala un brano | form → Telegram |
| `/it/contatti` | Contatti | form → Telegram + info |
| `/it/legali` | Note legali | disclaimer, licenze, privacy |
| `/it/404` | Non trovato | messaggio + ricerca |

Endpoint non-pagina: `/it/api/cerca` (frammento HTML per HTMX), `POST /it/api/segnala`, `POST /it/api/contatti`.

File generati dal build, senza prefisso di lingua: `/sitemap.xml`, `/robots.txt`.
File generati dal build, **uno per lingua dell'interfaccia**: `/it/rss.xml`, `/en/rss.xml` — **non** c'è un feed alla root e non ci sono feed per band (D100). Vale per `/it/404.html` e per ogni pagina.

## Dettaglio delle pagine

### Home
- **In evidenza**: da 3 a 5 elementi scelti a mano, in cima.
- **Contatori**: brani tradotti / band / lingue, stessa riga del footer (D102, D103).
- **Recenti**: elenco ordinato per data di aggiunta, il più nuovo prima.
- La ricerca sta nell'header (dropdown live), non ha una sezione dedicata.

### Elenco band
- Ordine **alfabetico**.
- **Filtro per tag** di genere (tag liberi, decisi dall'autore).
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
[ banner ads ]                        (solo se le ads sono attive)
  titolo brano + band + lingue disponibili
  selettore lingua
  ┌─────┬──────────────────┬──────────────────┬─────┐
  │ ads │  ORIGINALE       │  TRADUZIONE      │ ads │   (colonne laterali solo da xl in su)
  └─────┴──────────────────┴──────────────────┴─────┘
[ chiedi altre canzoni → form segnalazione ]
[ banner ads ]                        (solo se le ads sono attive)
```

- **Sinistra**: testo originale, diviso in strofe. **Destra**: traduzione nella lingua selezionata.
- Se il brano ha **più voci**, sopra ogni strofa compare **il nome del cantante**; nessuna etichetta strutturale (mai "ritornello", "bridge").
- **Selettore lingua**: solo le lingue realmente presenti per quel brano (originali multipli inclusi).
- L'originale è sempre disponibile: non esistono stati vuoti.
- **Nessuna annotazione per verso**: il testo resta pulito.
- **"Chiedi altre canzoni"**: fascia in fondo al testo con il link al form di segnalazione.
- Gli **slot pubblicitari** compaiono solo se attivati: a ads spente il DOM non li contiene e il testo prende tutta la larghezza.

### Segnala / Contatti
- Form semplici, senza account: i dati vengono **inoltrati al bot Telegram** dell'autore.
- Protezione: rate-limit, honeypot, informativa privacy inline.
- Al lancio non parte nessuna email.

## Funzionalità trasversali

- **Tema** chiaro/scuro: automatico dal sistema + toggle manuale.
- **Dimensione testo** A- / A+ (accessibilità).
- **Lingua dell'interfaccia**: negoziata dal browser, fallback italiano; le lingue disponibili seguono quelle delle traduzioni presenti.
- **Breadcrumb** su tutte le pagine interne.
- **Contatori** pubblici: in homepage e nel footer, stessa riga e stessi numeri.
- **Feed RSS** per ogni lingua dell'interfaccia, linkato dal `<head>` (rel alternate) e dal footer.
- **`sitemap.xml` e `robots.txt`** generati dal build, per l'indicizzazione (D99).
- **Header** sticky su mobile con nav compressa.

## Cosa il sito NON fa (non è un TODO, è il progetto)

- Nessun account, nessun login, nessun profilo, **nessun pannello admin**.
- Nessun commento, nessun upload, nessun contenuto generato dagli utenti.
- Nessun "preferito" o ricerca salvata.
- Nessuna newsletter, nessuna email automatica (al lancio).
- Nessuna API pubblica documentata: gli unici endpoint dinamici sono ricerca e form.
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
| Slug cambiato | redirect 301 dal file dei redirect |
| Band omonima | slug disambiguato (`nirvana-us`) |
| Ads spente | nessun contenitore nel DOM, nessuno spazio vuoto |
| Ricerca senza risultati | messaggio breve, non un errore |

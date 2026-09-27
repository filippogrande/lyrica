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
| `/it/band/{band}/album/{album}/brano/{brano}` | **Brano** | vista a fronte + selettore lingua |
| `/it/segnala` | Segnala un brano | form → Telegram |
| `/it/contatti` | Contatti | form → Telegram + info |
| `/it/legali` | Note legali | disclaimer, licenze, privacy |
| `/it/404` | Non trovato | messaggio + ricerca |

Endpoint non-pagina: `/it/api/cerca` (frammento HTML per HTMX), `POST /it/api/segnala`, `POST /it/api/contatti`, `/rss.xml`, `/sitemap.xml`, `/robots.txt`, `/manifest.webmanifest`, `/sw.js`.

## Dettaglio delle pagine

### Home
- **In evidenza**: da 3 a 5 elementi scelti a mano (band o brani), in cima.
- **Contatori**: brani tradotti / band / lingue.
- **Recenti**: elenco ordinato per data di aggiunta, il più nuovo prima.
- Il campo di ricerca è raggiungibile dall'header (dropdown live), non serve una sezione dedicata.

### Elenco band
- Ordine **alfabetico**.
- **Filtro per tag** di genere (tag liberi, decisi dall'autore).
- Ogni voce: nome band + tag.

### Band
- Nome, descrizione breve, tag, paese/lingua originale, anno di formazione.
- Elenco degli **album che hanno almeno una traduzione pubblicata**. Gli altri non compaiono.
- Album senza cover: nessuna immagine, layout senza copertina.

### Album
- Titolo, anno, cover (600x600 webp) se presente.
- **Tracklist completa** in ordine.
- Brani strumentali: presenti in tracklist con la nota "strumentale", **senza link**.
- Brani senza alcuna traduzione: presenti con la nota "solo originale", **senza link**.
- Ogni brano tradotto: link alla pagina brano + etichetta delle lingue disponibili.

### Brano (la pagina centrale)
- **Sinistra**: testo originale, diviso in strofe.
- **Destra**: traduzione nella lingua selezionata.
- Se il brano ha **più voci**, sopra ogni strofa compare **il nome del cantante**; nessuna etichetta strutturale (mai "ritornello", "bridge").
- **Selettore lingua** in alto: solo le lingue realmente presenti per quel brano (originali multipli inclusi).
- L'originale è sempre disponibile; non esistono stati vuoti.
- Nessuna annotazione per verso: il testo resta pulito.

### Segnala / Contatti
- Form semplici, senza account: i dati vengono **inoltrabili al bot Telegram** dell'autore.
- Protezione: rate-limit, honeypot, informativa privacy inline.
- Al lancio non parte nessuna email.

## Funzionalità trasversali

- **Tema** chiaro/scuro: automatico dal sistema + toggle manuale.
- **Dimensione testo** A- / A+ (accessibilità).
- **Lingua dell'interfaccia**: negoziata dal browser, fallback italiano.
- **Breadcrumb** su tutte le pagine interne.
- **Contatori** pubblici nel footer.
- **Header** sticky su mobile con nav compressa.

## Cosa il sito NON fa (non è un TODO, è il progetto)

- Nessun account, nessun login, nessun profilo, **nessun pannello admin**.
- Nessun commento, nessun upload, nessun contenuto generato dagli utenti.
- Nessun "preferito" o ricerca salvata.
- Nessuna newsletter, nessuna email automatica (al lancio).
- Nessuna API pubblica documentata: l'unico endpoint dinamico è la ricerca.
- Nessun dato personale conservato sul server (i form inoltrano e dimenticano).

## Casi limite e come si comportano

| Caso | Comportamento |
|---|---|
| Album senza cover | layout senza immagine, nessun placeholder |
| Brano strumentale | in tracklist con nota, nessuna pagina |
| Brano senza traduzioni | in tracklist con nota "solo originale", nessuna pagina |
| Brano con più lingue originali | sono opzioni del selettore lingua come le traduzioni |
| Traduzione a metà | non si pubblica (o meglio: non si costruisce la pagina) |
| Lingua UI non tradotta | fallback italiano |
| Slug cambiato | redirect 301 dal file dei redirect |
| Band omonima | slug disambiguato (`nirvana-us`) |

# FRONTEND — interfaccia, pagina brano, lingue, ricerca, PWA

## Principi

1. **Il testo è il contenuto**: la pagina brano è progettata attorno alla leggibilità del testo, non attorno a decorazioni.
2. **Mobile-first**: si progetta sul telefono e si allarga al desktop, mai il contrario.
3. **Pulito e senza logo**: tipografia e spaziatura fanno il look, non un marchio disegnato.
4. **Niente movimento inutile**: nessuna animazione che non comunichi stato.

## Base tecnica

- **Bootstrap 5, solo CSS**, servito **dal repo** (`assets/css/vendor/`), non da CDN: serve all'offline della PWA e a non dipendere da terzi a runtime.
- Personalizzazioni in `assets/css/lyrica.css`, **solo custom properties** dove possibile.
- **JS custom ridotto al minimo** (tema, dimensione testo, menu, selettore lingua, filtro alfabetico delle band): poche decine di righe, nessun bundle.
- **Nessuno script inline**: la CSP non ammette `unsafe-inline`, quindi anche il tema è un file (`assets/js/theme.js`) caricato nel `<head>` **prima** del paint.

## Griglia e breakpoints

Breakpoints Bootstrap standard: `sm 576`, `md 768`, `lg 992`, `xl 1200`.

| Contesto | Regola |
|---|---|
| Vista a fronte (pagina brano) | `col-12` sotto `lg`, `col-lg-6` da `lg` in su |
| Elenco band (pagina Bands) | griglia di card `auto-fill`: 1 colonna su mobile, 2-3 da `md`, 3-4 da `lg` |
| Hero della pagina band | foto `col-12 col-md-4`, anagrafica `col-12 col-md-8`; sotto `md` si impilano |
| Hero della pagina album | copertina `col-12 col-md-4`, anagrafica `col-12 col-md-8`; sotto `md` si impilano |
| Card album (corsie delle pagine band e album) | stessa misura delle card di corsia (`--rail-card`) |
| Corsie | card larghe `--rail-card`: **10rem** su mobile, **12rem** da `md` |
| Header | sticky sempre, nav compressa sotto `lg` con offcanvas |

## Tipografia

- **Font web multialfabeto**: un font che copre latino, cirillico e greco; per CJK si usa uno stack di fallback di sistema, perché un font CJK completo pesa troppo per il budget.
- Dimensioni relative alla **root**: tutto in `rem`, così il controllo A-/A+ scala l'intera pagina.
- Testo del brano: misura di lettura comoda, interlinea generosa, `max-width` contenuto.
- Numeri e date con `font-variant-numeric: tabular-nums` dove allineati (incluse le **traccie della tracklist**).

## Colori e tema

- Tema gestito con **`data-bs-theme`** su `<html>` (`light`/`dark`), più custom properties per i colori del testo.
- **Auto di default** (`prefers-color-scheme`), **toggle manuale** che salva in `localStorage`.
- **Nessun flash di tema sbagliato**: `assets/js/theme.js` applica il tema salvato prima del primo paint; essendo un file esterno, la CSP resta senza `unsafe-inline`.
- Contrasto **almeno AA** in entrambi i temi.

## Componenti

| Componente | Note |
|---|---|
| **Header** | sticky; nome del sito, nav (Home / Bands / Ricerca), toggle tema, A-/A+; sotto `lg` nav in **offcanvas** |
| **Ricerca** | dropdown live nell'header |
| **Breadcrumb** | su tutte le pagine interne |
| **Card band** | foto se dichiarata (`image` in `band.md`), altrimenti **segnaposto con l'iniziale**; nome, paese, generi, numero di album |
| **Corsia (rail)** | elenco orizzontale: titolo (`.rail-heading`), link "tutte le band" nella corsia delle band, card scorrevoli (vedi "Home a corsie"); usata in home, nella **pagina Bands** e nelle pagine **band** e **album** |
| **Card album** | cover se esiste, altrimenti **segnaposto con l'iniziale**; titolo e anno di pubblicazione (corsie della pagina band e della pagina album) |
| **Anagrafica band** | blocco `dl` con paese, attiva dal, membri, generi, più la descrizione; nella **pagina band** sta a **destra** della foto |
| **Anagrafica album** | blocco `dl` con band (link), anno, numero di brani, lingue tradotte; nella **pagina album** sta a **destra** della copertina |
| **Badge lingua** | lingue disponibili di un brano (`DE`, `IT`, ...) |
| **Strofa** | blocco di testo con eventuale nome del cantante sopra |
| **Avvisi** | "strumentale", "solo originale": testo semplice, non allarmi colorati |
| **Footer** | disclaimer, licenza CC, contatti, RSS, contatori |

## Accessibilità (requisiti, non buone intenzioni)

- Tutto raggiungibile **da tastiera**, con **focus visibile** (mai `outline: none`).
- **Skip link** "vai al testo" come primo elemento focalizzabile.
- Controllo **dimensione testo A-/A+**: scala la root `font-size`, salva la preferenza, si ferma a un massimo ragionevole.
- Immagini decorative con `alt=""`; cover con `alt` descrittivo.
- Ogni blocco di testo ha l'attributo **`lang` corretto**: uno screen reader deve leggere il tedesco come tedesco.
- Contrasto verificato in entrambi i temi prima di chiudere la PR che tocca i colori.

## Budget di performance

- **~300 KB per pagina** di HTML + CSS; le immagini si contano a parte e sono ottimizzate.
- Cover 600x600 webp, `loading="lazy"` e dimensioni dichiarate (`width`/`height`) per evitare spostamenti di layout.
- Se una pagina sfonda il budget, **la pagina si semplifica**, non si alza il budget.

## Cosa non si usa

- Nessun **bundle JS** di Bootstrap, nessun framework JS.
- Nessun **font icon** pesante: le poche icone sono SVG inline.
- Nessuna immagine di sfondo, nessun carosello, nessun popup.
- Nessuno **slider con scorrimento automatico** e nessuna libreria di caroselli: le corsie sono elenchi in overflow (vedi "Home a corsie").

---

# Home a corsie

La home non è un elenco verticale: è una **sequenza di corsie orizzontali**, ognuna con un titolo e le card delle novità. Le corsie nascono in home e vengono riusate fuori dalla home: la **pagina Bands** (corsia delle band, vedi "Pagina Bands"), la **pagina band** (album e brani di quella band, vedi "Pagina band") e la **pagina album** (altri album della band, vedi "Pagina album"). Resta un elenco verticale solo la **tracklist dell'album**.

## Le corsie

| Ordine | Corsia | Cosa contiene | Dove porta |
|---|---|---|---|
| 1 | "In evidenza" (`home.featured`) | brani con `featured: true`, nell'ordine già ricevuto | pagina brano |
| 2 | "Ultimi brani tradotti" (`home.recent`) | brani pubblicati, dalla data di aggiunta più recente | pagina brano |
| 3 | "Ultime band aggiunte" (`home.latest_bands`), con link "Tutte le band" (`home.all_bands`) | band pubblicate, dalla più attiva | pagina band / elenco band |

- **Massimo 12 card per corsia**: oltre, la home smette di essere "le novità" e diventa un archivio da scaricare tutto.
- Una corsia **senza contenuti non esiste nel DOM**; se non c'è nessuna corsia, la home mostra lo **stato vuoto dichiarato** (`home.no_content`), non un layout finto.
- **Pubblicato** significa "con almeno una traduzione": un brano solo originale resta in tracklist e non entra nelle corsie.

## Ordine

- "Ultimi brani tradotti": **data di aggiunta** (`added_date`) decrescente (D37).
- "Ultime band aggiunte": **contributo più recente** della band — la data del suo brano pubblicato più recente, non la data in cui è nata la cartella; a parità di data, **nome alfabetico**. Una band senza brani pubblicati non ha una data e non compare (D79).
- "Ultimi album pubblicati" (pagina band) e "Altri album della band" (pagina album): **anno di pubblicazione** decrescente; a parità di anno, titolo. Un album senza anno dichiarato finisce in coda (D88, D89).

## Anatomia di una card

```
┌──────────────┐
│              │   immagine 1:1 — cover dell'album (corsia brani e album)
│   immagine   │   o foto della band (corsia band)
│              │   oppure segnaposto con l'iniziale
└──────────────┘
Titolo del brano            ← riga 1 (o nome della band, o titolo dell'album)
Band — Album                ← riga 2, in secondo piano
```

- **Immagine**: cover dell'album nelle corsie dei brani e degli album, foto della band in quella delle band; il file sta in `covers/` (`docs/CONTENT.md`).
- **Segnaposto**: quando l'immagine non c'è si mostra **l'iniziale** (dell'album o del nome) su fondo del tema. È un segnaposto **dichiarato**: nessuna immagine inventata, nessuna richiesta in rete, e la card non cambia dimensioni.
- **Seconda riga unica**: per i brani `Band — Album`, per le band `Paese · Generi · N album`, per gli album l'**anno di pubblicazione** (vuota se l'anno non è dichiarato). Le parti vuote si saltano, il separatore resta uno solo.
- Sotto la card non c'è altro: niente descrizione, niente badge lingua (quelli stanno in pagina).

## Scorrimento senza JavaScript

- La corsia è un **elenco `<ul>` che scorre lateralmente** (`overflow-x: auto` + `scroll-snap-type: x proximity`): swipe sul telefono, rotella o frecce da tastiera. **Nessun carosello**, nessuno scorrimento automatico, nessun pallino di paginazione: le card restano nel flusso del documento e sono leggibili da uno screen reader.
- Ogni card è **un link normale**: si apre in una nuova scheda, si copia, si mette nei preferiti. Nessun click intercettato da JS.
- Larghezza della card: custom property **`--rail-card`** (10rem su mobile, 12rem da `md`), immagine sempre **1:1** (`aspect-ratio`).
- L'unica scorciatoia è **"Tutte le band"** nella corsia delle band della home: l'elenco dei brani non esiste come pagina, quindi non c'è un "vedi tutti".

## Immagini nelle corsie

- `loading="lazy"`, `decoding="async"` e `width`/`height` **600x600** dichiarati: le corsie non spostano il layout mentre caricano e non pesano sul primo paint.
- `alt=""`: l'immagine è **decorativa** perché il titolo è nella stessa card — ripeterlo farebbe rumore a uno screen reader. Il segnaposto è `aria-hidden`.
- Nessun formato diverso dal webp e nessun `srcset` per ora: la dimensione è una sola (600x600) e si aggiunge complessità solo se il budget lo chiede.

## Budget immagini

Il budget delle immagini si misura **sul numero di immagini che una pagina carica**, non solo sul peso della singola.

| Contesto | Immagini | Caricamento |
|---|---|---|
| Corsia della home | 1 per card, fino a 12 per corsia | `lazy`: contano quando l'utente scorre |
| Pagina Bands | 1 per band, più la corsia in cima | `lazy`: contano quando l'utente scorre |
| Pagina band | 1 foto + fino a 24 nelle due corsie (12 + 12) | `lazy`: contano quando l'utente scorre |
| Pagina album | 1 copertina + fino a 12 nella corsia degli altri album | `lazy` |
| Copertina di un album in pagina album | 1 | `lazy` |

- Una home con tre corsie piene dichiara fino a **36 immagini**: è il caso peggiore, ed è il motivo per cui tutte sono `lazy` e con dimensioni dichiarate.
- La **pagina Bands** mostra le band come card con foto (stesso formato delle corsie, D87): una pagina con N band dichiara N immagini, tutte `lazy` e con dimensioni dichiarate, quindi contano solo quando si scorre.
- La **pagina band** ha due corsie (album e brani della band) più la foto: anche qui tutto è `lazy` e con dimensioni dichiarate. La **pagina album** aggiunge la corsia degli altri album.
- Le immagini caricate dall'autore sono **600x600 webp, una sola dimensione** (D78): nessuna miniatura separata, nessun ridimensionamento a runtime.
- Se una pagina sfonda il budget, **si riduce il numero di card** (o si esclude una corsia), non si alza il budget.

---

# Pagina Bands

La pagina delle band ha due parti, in quest'ordine (D87):

1. **In cima** la corsia **"Ultime band aggiunte"** (`home.latest_bands`), identica a quella della home: stesse card, stesso scorrimento senza JS. È l'unica corsia della pagina Bands; qui il link "Tutte le band" è **assente**, perché sarebbe un link alla pagina stessa.
2. **Sotto** l'**indice alfabetico completo**: una barra di iniziali con **solo le lettere presenti** più "Tutte", e per ogni lettera un gruppo con una **griglia di card**.

## Le card

Stesso formato delle card delle corsie: immagine 1:1 (foto della band, o **segnaposto con l'iniziale** se manca), nome, e la riga di informazioni `Paese · Generi · N album`. L'immagine ha `alt=""`: è decorativa perché il nome è nella stessa card.

## Il filtro per lettera

- **Tutte le band sono già nel DOM**, raggruppate per iniziale; il filtro **nasconde e mostra i gruppi** con l'attributo `hidden` (`assets/js/lyrica.js`), come il selettore di lingua della pagina brano: nessuna chiamata al server, nessuna pagina separata, funziona **offline** (D87).
- **Senza JavaScript** i gruppi restano **tutti visibili**: il filtro è un di più, non un requisito. Nessun contenuto è nascosto dietro lo script.
- La barra è una `<nav>` con `aria-label` tradotto (`bands.letters_label`); il pulsante attivo è marcato con `aria-pressed` e la classe `is-active`. "Tutte" (`bands.all`) mostra di nuovo tutti i gruppi.

## Griglia e breakpoints

`.band-grid` usa `grid-template-columns: repeat(auto-fill, minmax(var(--rail-card), 1fr))`: le colonne si adattano da sole alla larghezza, con la stessa misura delle card delle corsie (10rem su mobile, 12rem da `md`). Nessuna media query dedicata.

---

# Pagina band

La pagina di una band mostra, in quest'ordine (D88):

1. **Titolo** (nome della band) su tutta la larghezza, poi l'**hero a due colonne**: **foto a sinistra** (`col-12 col-md-4`) e **anagrafica a destra** (`col-12 col-md-8`: paese, attiva dal, membri, generi, descrizione). Sotto `md` le colonne **si impilano**. Se la band non ha foto (`image` in `band.md`), l'anagrafica prende **tutta la larghezza**: nessuna colonna vuota e nessun segnaposto grande su una pagina intera.
2. **Corsia "Ultimi album pubblicati"** (`band.latest_albums`): gli album pubblicati della band, **anno di pubblicazione** decrescente (a parità di anno, titolo), massimo 12 card. Card = copertina o **segnaposto con l'iniziale**, titolo dell'album, **anno** come riga di informazioni; la card è il link alla pagina dell'album. Un album senza anno dichiarato finisce in coda.
3. **Corsia "Ultimi brani tradotti"** (`home.recent`, la stessa etichetta della home): i brani pubblicati **di quella band**, per **data di aggiunta** decrescente. "Pubblicato" significa "con almeno una traduzione": i brani solo originali restano in tracklist nella pagina dell'album e non entrano in corsia.

- Le due corsie usano `railBlock` come la home: massimo 12 card, scorrimento in overflow senza JS, immagini `lazy` con dimensioni dichiarate.
- Una corsia **senza card non esiste nel DOM**: una band con un solo album pubblicato non mostra una corsia di brani vuota.
- La pagina band **non ha** più la lista verticale degli album: gli album si vedono in corsia, con la copertina.

---

# Pagina album

La pagina di un album mostra, in quest'ordine (D89):

1. **Titolo** (titolo dell'album) su tutta la larghezza, poi l'**hero a due colonne**: **copertina a sinistra** (`col-12 col-md-4`) e **anagrafica a destra** (`col-12 col-md-8`): **band** (link alla pagina band), **anno**, **numero di brani**, **lingue tradotte** (le lingue delle traduzioni presenti nell'album, senza ripetizioni e in ordine alfabetico). Sotto `md` le colonne **si impilano**. Se la copertina manca, l'anagrafica prende **tutta la larghezza**. La riga delle lingue è **assente** finché nell'album non c'è nessuna traduzione: nessun campo vuoto stampato.
2. **Tracklist numerata e completa** (`album.tracks`): **tutti** i brani dichiarati in `tracks:` in `album.md`, nell'ordine del disco, ognuno con il **numero di traccia**.
   - Il numero è la **posizione nell'album**, non quella fra i brani tradotti: non cambia quando si pubblica una traduzione in più, quindi il numero resta un riferimento stabile al disco.
   - I brani **pubblicati** (con almeno una traduzione) sono **link** alla pagina del brano, con le lingue disponibili accanto.
   - I brani **non pubblicati** restano in elenco **senza link e non cliccabili**, con l'etichetta *"Solo originale"* (`track.only_original`) o *"Strumentale"* (`track.instrumental`). Un brano strumentale sta in tracklist ma non ha pagina (D31).
   - La tracklist **non si accorcia** mai: mostrare solo i brani tradotti nasconderebbe che manca ancora del lavoro.
3. **Corsia "Altri album della band"** (`album.other_albums`): gli altri album **pubblicati** della stessa band, **anno di pubblicazione** decrescente, **escluso quello che si sta guardando**. Stesse card della pagina band (copertina o segnaposto con l'iniziale, titolo, anno). Se l'album è l'unico pubblicato della band, la corsia **non esiste nel DOM**.

- La corsia degli album è **la stessa** della pagina band: il costruttore sta in `internal/render/album_rail.go`, cambia solo il titolo.

---

# La vista a fronte (pagina brano)

È la pagina per cui esiste il sito: originale a sinistra, traduzione a destra.

## Struttura

```
┌─ header (sticky) ───────────────────────────────────────┐
├─ breadcrumb: Home > Band > Album > Brano ───────────────┤
├─ [ ads: banner largo e basso ] (se attive) ─────────────┤
├─ titolo brano + band + badge lingue ────────────────────┤
├─ [ selettore lingua:  DE | IT | EN ] ───────────────────┤
│         │                        │                      │
│  [ads]  │  ORIGINALE             │  TRADUZIONE         │  [ads]
│  later. │  (nome cantante)       │  (nome cantante)    │  later.
│         │  strofa 1              │  strofa 1           │
│  solo   │  (nome cantante)       │  (nome cantante)    │  solo
│  xl+    │  strofa 2              │  strofa 2           │  xl+
├─ [ chiedi altre canzoni → form segnalazione ] ──────────┤
├─ [ ads: banner largo e basso ] (se attive) ─────────────┤
├─ footer ────────────────────────────────────────────────┤
```

## Desktop (`xl` e oltre)

- Due colonne **affiancate** con il testo, più **due colonne sottili ai lati** per gli ads: `2 + 4 + 4 + 2` sulla griglia da 12.
- Tra `lg` e `xl`: **solo le due colonne di testo**, nessuna colonna laterale — lo spazio non basta senza comprimere la lettura.
- Le strofe corrispondenti stanno **alla stessa altezza**: la griglia si costruisce **per strofa** (riga per riga), non con due colonne indipendenti. Se una strofa ha meno versi dell'originale, la riga cresce e resta allineata in alto.
- Nessuna linea verticale di separazione pesante: lo spazio fa il lavoro.

## Mobile (sotto `lg`)

- **Una colonna**: prima l'originale, poi la traduzione, nella stessa pagina (nessun tab nascosto: il testo si scorre).
- Il **selettore lingua resta sticky** sotto l'header.
- Blocchi separati da titoletti "Originale" / "Traduzione" in maiuscoletto.
- **Nessuna colonna laterale di ads**: solo il banner in alto e quello dopo "chiedi altre canzoni".

## Strofe e cantanti

- Ogni strofa è un blocco `<div>` con i suoi versi.
- Il **tipo di strofa non si stampa mai** (niente "Ritornello", "Bridge").
- Se il brano ha **più voci**, sopra ogni strofa compare **il nome del cantante**, su **ogni** strofa. Se il brano ha una voce sola, nessuna etichetta (sarebbe rumore).
- Il nome del cantante compare **anche nella colonna della traduzione**, così le due colonne restano leggibili in parallelo.

## Selettore lingua

- Mostra **solo le lingue realmente presenti** per quel brano, comprese le lingue originali quando sono più di una.
- L'**originale** è sempre selezionabile; non esistono stati vuoti.
- Implementazione: **tutte le lingue presenti nel DOM** e commutate lato client con poche righe di JS (attributo `hidden`), **nessuna chiamata al server**: cambio istantaneo e funzionante **offline**.
- La lingua scelta si riflette nell'URL come parametro (`?lang=de`) via `history.replaceState`: link condivisibile, nessuna pagina separata, nessun URL duplicato.
- Usabile da tastiera, annuncia il cambio (`aria-live`).

## Sezione "chiedi altre canzoni"

- Sta **sotto il testo**, prima del banner finale: chi ha appena letto una traduzione è la persona con più probabilità di proporne un'altra.
- Fascia breve con una frase e un link al **form di segnalazione**: nessun form inline, nessun popup.
- Su **ogni** pagina brano.

## Note dentro la pagina

- **Nessuna annotazione per verso.** Il testo resta pulito.
- Una versione **live/acustica** è un blocco dello stesso testo con la nota "versione live": resta nella stessa pagina.
- Il blocco "informazioni" (traduttore, lingue, data di aggiunta) sta **in fondo**, non sopra il testo.

## Casi limite

| Caso | Resa |
|---|---|
| Brano con una sola lingua (nessuna traduzione) | **non ha pagina**: in tracklist come "solo originale" |
| Brano strumentale | in tracklist con nota, nessuna pagina |
| Traduzione con strofe più corte | colonna allineata in alto, nessun riempimento artificiale |
| Numero di voci mancante | nessuna etichetta cantante |
| Lingua con alfabeto diverso | il font multialfabeto deve coprirla, altrimenti è un bug |
| Nessun advertiser | gli slot non esistono nel DOM, il testo prende tutta la larghezza |

---

# Lingue dell'interfaccia

## Due concetti distinti

| Concetto | Cos'è | Dove si sceglie |
|---|---|---|
| **Lingua dell'interfaccia** | la lingua delle parti fisse del sito (Home, Bands, Ricerca, breadcrumb, messaggi) | prefisso URL: `/it/...`, `/en/...` |
| **Lingua del testo** | la lingua di una traduzione di un brano | selettore nella pagina brano (`?lang=xx`) |

`/it/` non significa "traduzione italiana": significa "sito in italiano". Un utente può leggere il sito in italiano e il testo in tedesco con traduzione in inglese.

## Da dove nascono le lingue dell'interfaccia

**Dalle lingue in cui gli utenti leggono le traduzioni, non dalle lingue degli originali.**

Il ragionamento è quello dell'utente: se un brano tedesco ha traduzioni in italiano e inglese, chi arriva su quella pagina **parla italiano o inglese** — quasi mai tedesco. Regola operativa:

- **Nessuna interfaccia in una lingua di un originale**: che esista un testo tedesco non implica un'interfaccia tedesca.
- **Una traduzione in una lingua nuova crea pubblico in quella lingua**: quando compaiono traduzioni in francese, serve l'interfaccia in francese.
- Non si anticipa una lingua "per simmetria": una lingua senza contenuti in quella lingua non serve a nessuno.

Caso iniziale: **italiano + inglese**, con le successive quando la prima traduzione in quella lingua entra nel sito.

## Come si implementa

- Italiano = default e **fallback di tutto**.
- Ogni lingua = un file `locales/<lang>.yaml`; il build genera le pagine per **ogni lingua con un locale presente**.
- Aggiungere una lingua = copiare `locales/it.yaml`, tradurre i **valori**, aprire una PR.
- Nessuna lingua viene servita se il locale è incompleto.

## URL

- Ogni pagina esiste con **path prefix lingua**: `/it/band/...`, `/en/band/...`.
- **La radice `/` non è una pagina**: risponde con un **redirect** verso la lingua negoziata.
- Ogni pagina dichiara `<link rel="alternate" hreflang="...">` per tutte le lingue disponibili + `x-default`.
- `canonical` sempre verso l'URL corrente.
- I link interni passano **sempre** dal prefisso di lingua (helper di percorso), dalle corsie della home in poi: un link costruito a mano senza prefisso è un 404.

## Negoziazione

1. Richiesta a `/` → il server legge `Accept-Language` e redirige alla prima lingua supportata.
2. Se nessuna lingua corrisponde → **italiano**.
3. Il cambio lingua manuale nell'header è un **link** alla stessa pagina in un'altra lingua (nessun JS, nessun cookie).

## Stringhe

- Tutte le stringhe visibili stanno in `locales/<lang>.yaml`: **nessun testo hard-coded nei template**.
- Chiavi piatte e parlanti (`nav.home`, `track.translator`, `search.no_results`).
- **Chiave mancante = fallback italiano** + warning in build (non blocca, ma la PR deve sistemarlo o giustificarlo).
- I nomi propri (band, brani, cantanti) **non si traducono mai**.

## Cosa si adatta e cosa no

Si adattano: formato delle date, etichette dei contatori, meta description e `og:` tags (per lingua).

Non cambiano: i contenuti dei brani (sono dato, non interfaccia) e gli slug (stesso slug in tutte le lingue).

## Validazione

Il validatore controlla che tutte le lingue di `locales/` abbiano **le stesse chiavi** del locale italiano. Segnala anche (warning) quando una lingua di traduzione presente nei contenuti **non ha** un locale di interfaccia: è il promemoria che quella lingua ha ormai un pubblico.

---

# Ricerca full-text

## Cosa si cerca

1. **Titoli** dei brani
2. **Nomi delle band**
3. **Versi dei testi** (originali e traduzioni)

La ricerca per verso è la ragione per cui la ricerca esiste: chi ricorda "quella frase" non ricorda il titolo.

## Indice

- Costruito **a build time** e salvato in `public/` come **JSON** unico: nessun database, nessun servizio esterno.
- Ogni voce: tipo (`track` \| `band`), titolo, band, lingue disponibili, URL, versi normalizzati.
- **Normalizzazione** a build time: minuscole, rimozione dei segni diacritici, spazi collassati, punteggiatura ignorata.
- L'indice è caricato in memoria dal server all'avvio: ricerca istantanea, niente I/O per richiesta.

## Endpoint

`GET /it/api/cerca?q=<query>`

- Restituisce un **frammento HTML** già renderizzato (per HTMX), non JSON: il risultato è nella lingua dell'interfaccia senza JS che lo dipinge.
- Massimo **10 risultati**; se ce ne sono altri, l'ultima voce è "affina la ricerca", non una paginazione.
- Query sotto i **2 caratteri**: risposta vuota, nessuna richiesta inutile.
- Rate-limit per IP, per non trasformare la ricerca in un endpoint da abusare.

## Interfaccia

- Campo nell'header; digitando, HTMX aggiorna il dropdown con **debounce ~200 ms**.
- Ogni risultato mostra **titolo del brano + band**, mai uno snippet del verso.
- Navigazione da **tastiera** (frecce, invio) e ruoli ARIA corretti (`combobox`/`listbox`).
- **Esc** chiude e svuota; click fuori chiude.
- Senza risultati: messaggio breve, non un errore.

## Ordinamento (rilevanza)

1. Match nel **titolo**
2. Match nel nome della **band**
3. Match nel **verso**

A parità di categoria: prima i brani con più lingue, poi ordine alfabetico. Nessun punteggio opaco: l'ordine è spiegabile a voce.

## Cosa NON fa

- Nessuna pagina di risultati dedicata.
- Nessuna ricerca semantica, nessun sinonimo, nessuna tolleranza ai refusi: se servirà, sarà una decisione nuova.
- Nessuna cronologia, nessun suggerimento personalizzato.

La ricerca sta nell'header di tutte le pagine e nella **404**: chi sbaglia URL deve poter cercare subito.

---

# PWA — installabilità e lettura offline

## Obiettivo

Un testo tradotto serve spesso **senza rete** (in metro, in aereo, in un posto senza campo): il brano che stai leggendo resta leggibile anche se il telefono perde la connessione.

## Manifest

`manifest.webmanifest`: `name` e `short_name` Lyrica, `display: standalone`, `theme_color` coerente col tema scuro, `background_color` coerente col tema, icone (al lancio **placeholder**, dimensioni 192, 512, maskable). Nessun prompt di installazione forzato: si usa quello del browser.

## Service worker

- Scritto a mano, poche decine di righe, **nessun tool di build**.
- **In cache**: CSS, font, JS minimo (shell) e le **pagine dei brani visitate**.
- Strategie: shell **cache-first** (aggiornata con la versione dell'app); pagine brano **stale-while-revalidate**; cover cache-first.
- **Non in cache**: risposte di ricerca e form — meglio un errore chiaro che un comportamento ambiguo.

## Cosa funziona offline

| Funzione | Offline |
|---|---|
| Pagine brano già visitate (tutte le lingue, cambio lingua compreso) | sì |
| Navigazione tra pagine già visitate | sì |
| Cambio tema e dimensione testo | sì (è locale) |
| Filtro alfabetico delle band | sì (è locale) |
| Ricerca | no (avviso chiaro) |
| Segnalazioni e contatti | no (il form lo dice prima dell'invio) |
| Nuove traduzioni non ancora aperte | no |

## Aggiornamento

Quando esce una nuova versione, il service worker aggiorna la shell **al reload successivo**: nessun banner "nuova versione disponibile" che copre il testo. La versione della cache è legata al commit di build.

## Vincoli

- Il sito resta **perfettamente usabile senza installare nulla**: la PWA è un extra.
- Nessuna notifica push al lancio.
- La cache è limitata dalle regole sopra: non si scarica il sito intero sul telefono dell'utente.

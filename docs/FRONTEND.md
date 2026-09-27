# FRONTEND — interfaccia, pagina brano, lingue, ricerca, PWA

## Principi

1. **Il testo è il contenuto**: la pagina brano è progettata attorno alla leggibilità del testo, non attorno a decorazioni.
2. **Mobile-first**: si progetta sul telefono e si allarga al desktop, mai il contrario.
3. **Pulito e senza logo**: tipografia e spaziatura fanno il look, non un marchio disegnato.
4. **Niente movimento inutile**: nessuna animazione che non comunichi stato.

## Base tecnica

- **Bootstrap 5, solo CSS**, servito **dal repo** (`assets/css/vendor/`), non da CDN: serve all'offline della PWA e a non dipendere da terzi a runtime.
- Personalizzazioni in `assets/css/lyrica.css`, **solo custom properties** dove possibile.
- **JS custom ridotto al minimo** (tema, dimensione testo, menu, selettore lingua): poche decine di righe, nessun bundle.
- **Nessuno script inline**: la CSP non ammette `unsafe-inline`, quindi anche il tema è un file (`assets/js/theme.js`) caricato nel `<head>` **prima** del paint.

## Griglia e breakpoints

Breakpoints Bootstrap standard: `sm 576`, `md 768`, `lg 992`, `xl 1200`.

| Contesto | Regola |
|---|---|
| Vista a fronte (pagina brano) | `col-12` sotto `lg`, `col-lg-6` da `lg` in su |
| Elenco band | 1 colonna su mobile, 2 da `md`, 3 da `lg` |
| Card album | 1 su mobile, 2-3 da `md` |
| Header | sticky sempre, nav compressa sotto `lg` con offcanvas |

## Tipografia

- **Font web multialfabeto**: un font che copre latino, cirillico e greco; per CJK si usa uno stack di fallback di sistema, perché un font CJK completo pesa troppo per il budget.
- Dimensioni relative alla **root**: tutto in `rem`, così il controllo A-/A+ scala l'intera pagina.
- Testo del brano: misura di lettura comoda, interlinea generosa, `max-width` contenuto.
- Numeri e date con `font-variant-numeric: tabular-nums` dove allineati.

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
| **Card band** | nome + tag; nessuna immagine |
| **Card album** | cover se esiste, altrimenti layout testuale (titolo + anno + n° brani) |
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
| Ricerca | no (avviso chiaro) |
| Segnalazioni e contatti | no (il form lo dice prima dell'invio) |
| Nuove traduzioni non ancora aperte | no |

## Aggiornamento

Quando esce una nuova versione, il service worker aggiorna la shell **al reload successivo**: nessun banner "nuova versione disponibile" che copre il testo. La versione della cache è legata al commit di build.

## Vincoli

- Il sito resta **perfettamente usabile senza installare nulla**: la PWA è un extra.
- Nessuna notifica push al lancio.
- La cache è limitata dalle regole sopra: non si scarica il sito intero sul telefono dell'utente.

# DESIGN_SYSTEM — layout, tema, accessibilità

## Principi

1. **Il testo è il contenuto**: la pagina brano è progettata attorno alla leggibilità del testo, non attorno a decorazioni.
2. **Mobile-first**: si progetta sul telefono e si allarga al desktop, mai il contrario.
3. **Pulito e senza logo** (D04): tipografia e spaziatura fanno il look, non un marchio disegnato.
4. **Niente movimento inutile**: nessuna animazione che non comunichi stato.

## Base tecnica

- **Bootstrap 5, solo CSS**, servito **dal repo** (`assets/css/vendor/`), non da CDN: serve per l'offline della PWA e per non dipendere da terzi.
- Personalizzazioni in `assets/css/lyrica.css`, **solo custom properties** dove possibile.
- **JS custom ridotto al minimo** (tema, dimensione testo, apertura menu, selettore lingua): poche decine di righe, nessun bundle.

## Griglia e breakpoints

Breakpoints Bootstrap standard: `sm 576`, `md 768`, `lg 992`, `xl 1200`.

| Contesto | Regola |
|---|---|
| Vista a fronte (pagina brano) | `col-12` sotto `lg`, `col-lg-6` da `lg` in su |
| Elenco band | 1 colonna su mobile, 2 da `md`, 3 da `lg` |
| Card album | 1 su mobile, 2-3 da `md` |
| Header | sticky sempre, nav compressa sotto `lg` con offcanvas |

## Tipografia

- **Font web multialfabeto** (D19): un font che copre latino, cirillico e greco; per CJK si usa uno stack di fallback di sistema, perché un font CJK completo pesa troppo per il budget.
- Dimensioni relative alla **root**: tutto in `rem` così il controllo A-/A+ scala l'intera pagina.
- Testo del brano: misura di lettura comoda, interlinea generosa, `max-width` contenuto (niente righe di 150 caratteri su desktop).
- Numeri e date con `font-variant-numeric: tabular-nums` dove allineati.

## Colori e tema

- Tema gestito con **`data-bs-theme`** su `<html>` (`light`/`dark`), più custom properties per i colori del testo.
- **Auto di default** (`prefers-color-scheme`), **toggle manuale** che salva in `localStorage` (D50).
- **Nessun flash di tema sbagliato**: uno script inline nel `<head>` applica il tema salvato prima del primo paint.
- Contrasto **almeno AA** in entrambi i temi: il testo tradotto è il caso d'uso principale, deve essere leggibile.

## Componenti

| Componente | Note |
|---|---|
| **Header** | sticky, contiene nome del sito, nav (Home / Bands / Ricerca), toggle tema, A-/A+; sotto `lg` nav in **offcanvas** |
| **Ricerca** | dropdown live nell'header (vedi `SEARCH_DESIGN.md`) |
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
- Controllo **dimensione testo A-/A+** (D51): scala la root `font-size`, salva la preferenza, si ferma a un massimo ragionevole.
- Immagini decorative con `alt=""`; cover con `alt` descrittivo.
- Ogni blocco di testo ha l'attributo **`lang` corretto**: uno screen reader deve leggere il tedesco come tedesco.
- Contrasto verificato in entrambi i temi prima di chiudere la PR che tocca i colori.

## Budget di performance

- **~300 KB per pagina** di HTML + CSS (D18). Le immagini si contano a parte e sono ottimizzate.
- Le cover sono 600x600 webp (D21), caricate con `loading="lazy"` e dimensioni dichiarate (`width`/`height`) per evitare spostamenti di layout.
- Come si misura: dimensione del file HTML + CSS totale della pagina in `public/`, controllo a mano prima di chiudere la PR. Se una pagina sfonda il budget, la pagina si semplifica, non si alza il budget.

## Cosa non si usa

- Nessun **bundle JS** di Bootstrap, nessun framework JS.
- Nessun **font icon** pesante: le poche icone necessarie sono SVG inline.
- Nessuna immagine di sfondo, nessun carosello, nessun popup.

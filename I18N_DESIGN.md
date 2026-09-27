# I18N_DESIGN — lingue dell'interfaccia

## Due concetti distinti (importante)

| Concetto | Cos'è | Dove si sceglie |
|---|---|---|
| **Lingua dell'interfaccia** | la lingua delle parti fisse del sito (Home, Bands, Ricerca, breadcrumb, messaggi) | prefisso URL: `/it/...`, `/en/...` |
| **Lingua del testo** | la lingua di una traduzione di un brano | selettore nella pagina brano (`?lang=xx`) |

Non vanno confusi: `/it/` non significa "traduzione italiana", significa "sito in italiano". Un utente può leggere il sito in italiano e il testo in tedesco con traduzione in inglese.

## Lingue dell'interfaccia

- Italiano = default (fallback di tutto).
- Altre lingue si aggiungono creando un file in `locales/<lang>.yaml`; il build genera le pagine per **ogni lingua con un locale presente**.
- Nessuna lingua viene servita se non ha un locale completo: meglio meno lingue, tutte fatte bene.

## URL

- Ogni pagina esiste con **path prefix lingua** (D40): `/it/band/...`, `/en/band/...`.
- **La radice `/` non è una pagina**: risponde con un **redirect** verso la lingua negoziata.
- Ogni pagina dichiara `<link rel="alternate" hreflang="...">` per tutte le lingue disponibili + `x-default`.
- `canonical` sempre verso l'URL corrente.

## Negoziazione della lingua

1. Richiesta a `/` → il server legge `Accept-Language` e redirige alla prima lingua supportata (D52).
2. Se nessuna lingua corrisponde → **italiano**.
3. Il cambio lingua manuale nell'header è un **link** alla stessa pagina in un'altra lingua (nessun JS, nessun cookie, coerente con "zero account").

## Stringhe dell'interfaccia

- Tutte le stringhe visibili stanno in `locales/<lang>.yaml`, **nessun testo hard-coded nei template**.
- Chiavi piatte e parlanti (`nav.home`, `track.translator`, `search.no_results`).
- **Chiave mancante = fallback italiano** e il validatore di build emette **warning** (non blocca la build, ma la PR deve sistemarlo o giustificarlo).
- I nomi propri (band, brani, cantanti) **non si traducono mai**.

## Cosa si adatta alla lingua

- Formato delle date ("27 settembre 2026" vs "27 September 2026").
- Etichette dei contatori e delle azioni.
- Meta description e `og:` tags (generate per lingua).

## Cosa NON cambia

- I contenuti dei brani (testi e traduzioni: sono dato, non interfaccia).
- Gli slug: le band e i brani hanno **lo stesso slug in tutte le lingue** (se il titolo tradotto differisse, resta il titolo originale nello slug).

## Validazione

Il validatore controlla che tutte le lingue di `locales/` abbiano **le stesse chiavi** del locale italiano: una chiave presente in `it` e assente in `en` è un errore in CI.

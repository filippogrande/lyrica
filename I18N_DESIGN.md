# I18N_DESIGN — lingue dell'interfaccia

## Due concetti distinti (importante)

| Concetto | Cos'è | Dove si sceglie |
|---|---|---|
| **Lingua dell'interfaccia** | la lingua delle parti fisse del sito (Home, Bands, Ricerca, breadcrumb, messaggi) | prefisso URL: `/it/...`, `/en/...` |
| **Lingua del testo** | la lingua di una traduzione di un brano | selettore nella pagina brano (`?lang=xx`) |

Non vanno confusi: `/it/` non significa "traduzione italiana", significa "sito in italiano". Un utente può leggere il sito in italiano e il testo in tedesco con traduzione in inglese.

## Da dove nascono le lingue dell'interfaccia

**Dalle lingue in cui gli utenti leggono le traduzioni, non dalle lingue degli originali.**

Il ragionamento è quello dell'utente: se un brano tedesco ha traduzioni in italiano e inglese, chi arriva su quella pagina **parla italiano o inglese** — quasi mai tedesco. Quindi l'interfaccia in italiano e inglese copre il pubblico reale.

Regola operativa:

- **Nessuna interfaccia in una lingua di un originale**: che esista un testo tedesco non implica un'interfaccia tedesca.
- **Una traduzione in una lingua nuova crea pubblico in quella lingua**: quando compaiono traduzioni in francese, serve l'interfaccia in francese.
- Quindi le lingue di interfaccia seguono l'insieme delle **lingue di traduzione presenti nei contenuti** (più l'italiano, che è il default), e si aggiungono quando i contenuti crescono in quella lingua.
- Non si anticipa una lingua di interfaccia "per simmetria": una lingua senza contenuti in quella lingua non serve a nessuno.

Il caso iniziale è quindi: **italiano + inglese** (le due lingue in cui l'autore traduce), con l'aggiunta delle successive quando la prima traduzione in quella lingua entra nel sito.

## Come si implementa una lingua di interfaccia

- Italiano = default e **fallback di tutto**.
- Ogni lingua = un file `locales/<lang>.yaml`; il build genera le pagine per **ogni lingua con un locale presente**.
- Aggiungere una lingua = copiare `locales/it.yaml`, tradurre i **valori** (non le chiavi), aprire una PR: il build genera `/<lang>/...`, hreflang compresi.
- Nessuna lingua viene servita se il locale è incompleto: meglio meno lingue, tutte fatte bene.

## URL

- Ogni pagina esiste con **path prefix lingua** (D40): `/it/band/...`, `/en/band/...`.
- **La radice `/` non è una pagina**: risponde con un **redirect** verso la lingua negoziata.
- Ogni pagina dichiara `<link rel="alternate" hreflang="...">` per tutte le lingue disponibili + `x-default`.
- `canonical` sempre verso l'URL corrente.

## Negoziazione della lingua

1. Richiesta a `/` → il server legge `Accept-Language` e redirige alla prima lingua supportata (D52). **È il comportamento che l'utente si aspetta**: arriva e trova il sito nella sua lingua, senza scegliere nulla.
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

Controlla anche che **ogni lingua di traduzione presente nei contenuti abbia un locale di interfaccia** (warning, non errore): è il promemoria che quella lingua ha ormai un pubblico.

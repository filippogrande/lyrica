# BUILD_PIPELINE — CLI, validazione, build

## La CLI `lyrica`

Un solo binario Go, tre famiglie di comandi:

| Comando | Cosa fa |
|---|---|
| `lyrica new band \| album \| brano` | genera i template di contenuto (D20, vedi `CONTENT_AUTHORING.md`) |
| `lyrica build` | valida i contenuti e genera `public/` |
| `lyrica serve` | serve `public/` in locale (con gli endpoint dinamici), per verificare il risultato |

`lyrica build` è **deterministico**: stesso contenuto → stesso output. Non esiste build "incrementale" che produca risultati diversi tra due esecuzioni.

## Cosa fa `lyrica build`, in ordine

1. **Parse** di tutto `content/` (front-matter YAML + strofe) in modelli in memoria.
2. **Validazione** (vedi sotto). Se fallisce: output `public/` **non viene prodotto**, exit code ≠ 0.
3. **Indice di ricerca** su titoli, band e versi normalizzati → `public/search-index.json`.
4. **Render** delle pagine con Templ: home, elenco band, pagine tag, band, album, brani, legali, segnala, contatti, 404.
5. **Asset**: CSS (Bootstrap vendored + custom), font, JS minimo, cover.
6. **Output di distribuzione**: `rss.xml`, `sitemap.xml`, `robots.txt`, `manifest.webmanifest`, `sw.js`, file dei redirect.

## Regole di validazione

Sono l'unica forma di test richiesta al lancio (D17) e la loro assenza è un errore, non una mancanza accettabile.

**Bloccanti** (build rossa):
- slug incoerenti con file/cartelle;
- slug duplicati nella stessa collezione;
- **traduzioni incomplete** (strofe o versi che non corrispondono all'originale);
- blocco senza `role: original`, o stessa `lang` due volte nello stesso brano;
- `added_date` mancante o invalida;
- `cover` che punta a un file inesistente;
- chiavi di locale mancanti rispetto a `it.yaml`;
- link interni a band, album o brani inesistenti.

**Warning** (build verde, da sistemare o giustificare in PR):
- brano senza nessuna traduzione (ricordarsi di segnarlo "solo originale");
- album senza cover;
- tag usato una volta sola (possibile refuso: `metal` vs `metalcore`);
- descrizione band mancante.

## Cosa NON si fa nella pipeline

- **Nessuno unit test del codice Go** al lancio.
- **Nessun test end-to-end** (niente Playwright): la verifica è `lyrica build` + sito aperto a mano.
- **Nessun test isolato o usa-e-getta** che controlli solo la sintassi: mascherano i bug invece di rivelarli. Se una cosa va verificata, si verifica sul **sito reale generato**.

## CI (GitHub Actions)

| Evento | Job |
|---|---|
| **Pull Request** | `validate`: esegue `lyrica build` sui contenuti della PR. Rosso = non si mergia. |
| **Push su `main`** | `deploy`: build dell'immagine Docker, push nel registry, deploy sul home-lab (vedi `DEPLOY.md`) |

La validazione è veloce perché non compila tutto il sito due volte: un job unico che fallisce presto sui contenuti.

## Cache della build Docker

Nel `Dockerfile` gli asset e i contenuti si copiano **dopo** il download delle dipendenze Go e **prima** della compilazione, nell'ordine che non invalida la cache a ogni modifica di un testo:

1. `go.mod` / `go.sum` → `go mod download`
2. codice Go → compilazione
3. `content/`, `covers/`, `assets/`, `locales/` → render finale

Regola: una modifica a un **testo** non deve invalidare la cache delle dipendenze.

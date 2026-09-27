# BUILD_PIPELINE — CLI, validazione, build

## La CLI `lyrica`

Un solo binario Go, tre famiglie di comandi:

| Comando | Cosa fa |
|---|---|
| `lyrica new band \| album \| brano` | genera i template di contenuto (D20, vedi `CONTENT_AUTHORING.md`) — **FASE 2** |
| `lyrica build` | valida i contenuti e genera `public/` — **FASE 2** |
| `lyrica serve` | serve il sito in locale — **disponibile** |
| `lyrica help` | mostra i comandi — **disponibile** |

I comandi non ancora implementati **restituiscono un errore esplicito**, non fingono di aver funzionato (nessun fallback silenzioso, `DEVELOPMENT_GUIDELINES.md` §5).

## Dipendenze dall'ambiente

| Serve | Perché |
|---|---|
| **Go 1.27.1** | compilare il binario |
| **templ v0.3.1020** | i file `.templ` vanno convertiti in Go con `templ generate`. I file generati (`*_templ.go`) **non si committano** (sono in `.gitignore`) |
| **assets/css/vendor/bootstrap.min.css** | CSS di terze parti servito dal repo: si scarica con `./scripts/fetch-assets.sh` (Bootstrap 5.3.8) |
| **bash** | `scripts/fetch-assets.sh` ha shebang bash e usa `set -euo pipefail` |

Comandi locali, nell'ordine:

```
./scripts/fetch-assets.sh                                  # una volta (o quando cambi versione)
go mod tidy
go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate   # dopo ogni modifica ai .templ
go run ./cmd/lyrica serve
```

Senza `templ generate` il progetto **non compila**: la generazione non è opzionale.

## Trappole già pagate

Errori che hanno **già** rotto la CI o la build dell'immagine: non si ripetono.

1. **Nei file `.templ` non si importa `github.com/a-h/templ`.** Il generatore lo importa da sé, e un import esplicito produce `templ redeclared in this block` + `"github.com/a-h/templ" imported and not used` (e, a cascata, `undefined:` su componenti dello stesso package). Se serve un valore tipizzato da templ (es. `templ.SafeURL` per un `href`), lo si costruisce in un **file Go** (`internal/render/model.go`) ed esposto come metodo di `PageData`: il `.templ` lo consuma senza importare nulla.
2. **`golang:*-alpine` non ha bash.** Lo stage di build dell'immagine deve fare `apk add --no-cache curl bash`, altrimenti `RUN ./scripts/fetch-assets.sh` muore con `env: bash: No such file or directory`.
3. **`*_templ.go` non si committano mai** (`.gitignore`): sono artefatti di generazione e committarli fa divergere codice e template.

## Cosa farà `lyrica build`, in ordine (FASE 2)

1. **Parse** di tutto `content/` (front-matter YAML + strofe) in modelli in memoria.
2. **Validazione** (vedi sotto). Se fallisce: `public/` **non viene prodotto**, exit code ≠ 0.
3. **Indice di ricerca** su titoli, band e versi normalizzati → `public/search-index.json`.
4. **Render** delle pagine con Templ: home, elenco band, pagine tag, band, album, brani, legali, segnala, contatti, 404.
5. **Asset**: CSS, font, JS minimo, cover.
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
- **Nessun test end-to-end** (niente Playwright): la verifica è la build + il sito aperto a mano.
- **Nessun test isolato o usa-e-getta** che controlli solo la sintassi: mascherano i bug invece di rivelarli. Se una cosa va verificata, si verifica sul **sito reale generato**.

## CI (GitHub Actions)

| Workflow | Evento | Cosa fa |
|---|---|---|
| `ci.yml` | **Pull Request** e **push su `main`** | job `Compila`: `go mod tidy` → `templ generate` → `go vet` → `go build` |
| `ci.yml` | **Pull Request** (FASE 2) | job `valida`: esegue `lyrica build` sui contenuti; se i contenuti sono rotti, la PR è rossa |
| `docker-build-push.yml` | **push su `main`** | build dell'immagine e push su Docker Hub (`:latest` e `:<sha>`) |

Perché esiste il job `Compila`: l'ambiente in cui scrive l'agente **non ha un compilatore Go**, quindi la compilazione — e quindi gli errori di sintassi — si verificano qui. Non è un test unitario: è la verifica dell'artefatto.

Il workflow dell'immagine **non gira sulle PR**: una PR non deve poter pubblicare un tag `latest`. Il deploy sul home-lab resta manuale (`DEPLOY.md`).

## Cache della build Docker

Nel `Dockerfile` i passi sono ordinati perché una modifica a un **testo** non deve invalidare la cache delle dipendenze:

1. `go.mod` → `go mod download`
2. codice e asset (`COPY . .`) → `go mod tidy`
3. `templ generate` → `fetch-assets.sh` → `go build`

L'immagine finale **non** contiene il toolchain Go: solo il binario statico, `assets/` e `locales/`.

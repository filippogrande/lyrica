# ARCHITECTURE — come è fatto, come si compila, come si mette online

## In una riga

Sito **statico generato a build time** da contenuti Markdown, servito da un binario Go che espone solo tre cose dinamiche: la ricerca, i due form e il feed.

## Stato di implementazione

Questo documento descrive il **progetto**; non tutto quello che descrive esiste ancora. Le due
cose da non dimenticare leggendo il resto:

- **`internal/search/` e `internal/notify/` non esistono**: la ricerca full-text (FASE 4) e i
  form con invio a Telegram (FASE 6) sono progettati, non scritti. Il container serve **solo**
  file statici, `/healthz` e il redirect della radice.
- **`ads.yaml` non esiste** e non esiste il codice degli slot pubblicitari.

Per la lista completa, vedi [`REFERENCE.md`](REFERENCE.md#non-esiste-ancora-dove-non-cercare).

## Flusso

```
content/  covers/  assets/  locales/
        |
        |  lyrica build   (parse + validazione + render Templ)
        v
     public/            HTML, CSS, JS minimo, cover, indice di ricerca,
                        rss.xml (per lingua), sitemap.xml, robots.txt, manifest, sw.js
        |
        |  docker build (multi-stage: build Go -> immagine finale minimale)
        v
   immagine Docker  ->  docker compose  ->  Cloudflare Tunnel  ->  lyrica.filippomoscatelli.com
```

L'indice di ricerca, `manifest.webmanifest`, `sw.js` e il file dei redirect sono **previsti, non
prodotti** oggi: il build scrive HTML, `rss.xml`, `sitemap.xml`, `robots.txt`, `assets/` e
`covers/`.

## Componenti

| Componente | Responsabilità | Dove |
|---|---|---|
| **Content loader** | legge `content/`, parse di front-matter + strofe, costruisce i modelli | `internal/content/` |
| **Validator** | regole di integrità: front-matter, slug, link, traduzioni incomplete, cover | `internal/content/` |
| **Renderer** | pagine Templ (home, band, album, brano, 404, legali) | `internal/render/` |
| **Search index** | costruisce l'indice a build time e lo interroga | **non esiste** (FASE 4) |
| **i18n** | stringhe UI e negoziazione lingua | `internal/i18n/`, `locales/` |
| **Notifier** | invio messaggi al bot Telegram (segnalazioni, contatti) | **non esiste** (FASE 6) |
| **Web** | server, routing, header di sicurezza, 404 | `internal/web/` |
| **CLI** | `lyrica new ...`, `lyrica valida`, `lyrica build`, `lyrica serve`, `lyrica migrate` | `cmd/lyrica/` |

Il codice sta in inglese, i contenuti e i doc in italiano (`GUIDELINES.md` §4).

## Cosa gira a runtime

Il container serve soltanto:

1. **File statici**: tutte le pagine già generate (`public/`).
2. **`GET /healthz`**: risposta per gli healthcheck.
3. **Redirect della radice**: `/` verso la lingua negoziata da `Accept-Language`.

**Non esistono ancora** `GET /it/api/cerca?q=` (ricerca nell'indice, frammento HTML per HTMX) e
`POST /it/api/segnala` / `POST /it/api/contatti` (validazione, rate-limit, inoltro al bot Telegram):
sono FASE 4 e FASE 6.

Non esiste database, non esiste sessione, non esiste stato utente. Conseguenze volute: nessuna migrazione, nessun backup di dati applicativi, nessun dato personale a riposo.

## Scelte architetturali e alternative scartate

| Scelta | Alternativa scartata | Perché |
|---|---|---|
| Build-time static | SSR a runtime con database | niente DB da mantenere, niente GDPR sui dati a riposo, sito sempre identico e cacheabile |
| Go + Templ | Node/Next, PHP | build veloce, un binario, nessun runtime JS da servire |
| HTMX | SPA con framework JS | un solo endpoint dinamico non giustifica un framework |
| Bootstrap 5 CSS-only | Tailwind con build npm | nessun passaggio npm per servire il sito |
| Docker + Cloudflare Tunnel | k3s | k3s in dismissione nel home-lab |
| Bot Telegram | pannello admin / email | nessun login, nessuna deliverability da gestire |

## Dipendenze dall'ambiente

| Serve | Perché |
|---|---|
| **Go 1.27.1** | compilare il binario |
| **templ v0.3.1020** | i file `.templ` vanno convertiti in Go con `templ generate`; i generati (`*_templ.go`) **non si committano** |
| **bash** | lo script degli asset ha shebang bash e usa `set -euo pipefail` |
| **curl** | download di Bootstrap |
| **assets/css/vendor/bootstrap.min.css** | CSS di terze parti **servito dal repo**, non da CDN (offline della PWA + nessuna dipendenza a runtime) |

Comandi locali, nell'ordine:

```
bash scripts/fetch-assets.sh                               # una volta (o quando cambi versione)
go mod tidy
go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate   # dopo ogni modifica ai .templ
go run ./cmd/lyrica serve
```

Lo script si invoca con `bash` esplicito, non con `./`: nel repo non ha il bit di esecuzione.

## La CLI `lyrica`

| Comando | Cosa fa |
|---|---|
| `lyrica valida` | controlla i contenuti e stampa errori e avvisi, **senza generare nulla** — **disponibile** |
| `lyrica serve` | serve il sito in locale — **disponibile** |
| `lyrica build` | valida i contenuti e genera `public/` — **disponibile** |
| `lyrica new band \| album \| brano` | genera i template di contenuto — **disponibile** |
| `lyrica migrate [--band <slug>] [--apply]` | migra i brani dal file unico alla cartella (prova senza `--apply`) — **disponibile** |
| `lyrica help` | mostra i comandi — **disponibile** |

I comandi non ancora implementati **restituiscono un errore esplicito**, non fingono di aver funzionato.

`lyrica valida` esce con codice ≠ 0 se c'è anche un solo problema bloccante, e stampa **tutti** i problemi trovati (non si ferma al primo). Gli avvisi non fanno fallire il comando: restano visibili su stderr.

### Cosa fa `lyrica build`, in ordine

1. **Parse** di tutto `content/` in modelli in memoria.
2. **Validazione**. Se fallisce: `public/` **non viene prodotto**, exit code ≠ 0.
3. **Indice di ricerca** → `public/search-index.json` — **non ancora implementato** (FASE 4).
4. **Render** delle pagine con Templ.
5. **Asset**: CSS, font, JS minimo, cover.
6. **Output di distribuzione**: `rss.xml` (per ogni lingua dell'interfaccia, in `public/<lang>/rss.xml`), `sitemap.xml` (uno solo, alla root), `robots.txt` (alla root). `manifest.webmanifest`, `sw.js` e il file dei redirect sono **previsti, non prodotti**.

I file XML e `robots.txt` non sono componenti Templ: si scrivono con `encoding/xml` e con `writeRaw` (`internal/build/raw.go`), perché sono dati e non markup. Un errore nella loro scrittura **fa fallire la build**: un feed a metà è un feed rotto e nessuno se ne accorgerebbe.

## CI (GitHub Actions)

| Workflow | Evento | Cosa fa |
|---|---|---|
| `ci.yml` → job `compila` | **PR** e **push su `main`** | `go mod tidy` → `templ generate` → `go vet` → `go build` → `go test ./...` (D98) |
| `ci.yml` → job `valida` | **PR** e **push su `main`** | `go run ./cmd/lyrica valida`, poi `lyrica build` e controllo che le pagine attese esistano |
| `docker-build-push.yml` | **push su `main`** | build dell'immagine e push su Docker Hub (`:latest` e `:<sha>`) |
| `immagini-mancanti.yml` | **manuale** (`workflow_dispatch`) | cerca le immagini che mancano, le scarica dagli archivi pubblici e apre una PR |
| `immagini.yml` | **manuale** (`workflow_dispatch`) | scarica **una** immagine da un MBID e apre una PR |
| `migrazione-brani.yml` | **manuale** (`workflow_dispatch`) | esegue `lyrica migrate` e apre una PR |

Perché esiste il job `compila`: l'ambiente in cui scrive l'agente **non ha un compilatore Go**, quindi gli errori di sintassi si vedono solo qui. Non è un test unitario: è la verifica dell'artefatto.

Il workflow dell'immagine **non gira sulle PR**: una PR non deve poter pubblicare un tag `latest`.

La sicurezza ha un workflow suo, `.github/workflows/security.yml`, separato da `ci.yml`, con tre job. `dipendenze` esegue `govulncheck`: segnala solo le vulnerabilità che il codice chiama davvero, non ogni CVE presente nei moduli. `sast` esegue `gosec` sul codice nostro con `-exclude-generated` (i `*_templ.go` sono artefatti di generazione) e carica i rilievi in SARIF nella scheda code scanning. `immagine` costruisce l'immagine vera dal `Dockerfile` e la passa a Trivy, poi fa una seconda scansione con `scan-type: config` su `Dockerfile` e workflow: l'OS di base e i suoi pacchetti cambiano anche quando il `go.mod` non cambia. Gira sulle **PR**, sui **push su `main`** e una volta a settimana con un cron. Gli strumenti sono installati a `@latest` perché non sono dipendenze del progetto e un database di CVE vecchio è peggio di un tool non pinnato.

## Docker

### Immagine

- Registry **Docker Hub**, `filippogrande/lyrica`, tag `latest` + `<sha del commit>` (il tag con la SHA serve al rollback).
- `Dockerfile` **multi-stage**:
  - stage `build`: `golang:1.27.1-alpine3.24` → `apk add curl bash` → dipendenze → `templ generate` → script asset → `lyrica build` → binario statico `CGO_ENABLED=0`;
  - stage finale: `alpine:3.24.2` + `ca-certificates`, utente non privilegiato (`lyrica`, uid 10001), **solo** binario, `public/`, `locales/`.
- `HEALTHCHECK` su `/healthz`.

Ordine dei passi studiato perché una modifica a un **testo** non invalidi la cache delle dipendenze: `go.mod` → `go mod download` → `COPY . .` → `go mod tidy` → `templ generate` → asset → `lyrica build` → build del binario.

### Deploy sul home-lab

Compose in **`deploy/docker-compose.yml`** (nel repo), copiato sul master in **`/mnt/applicazioni/yml/docker/lyrica/`** insieme alla sua `.env`:

```
cp .env.example .env
docker compose pull
docker compose up -d
```

- Il container ascolta sulla **8080 interna**.
- La **porta host è locale alla macchina**, non un valore del progetto: si sceglie quella libera. Sul master è la **8097** (8090-8096 risultano occupate). Se è occupata, si cambia **solo quel numero** nel compose — non è un deploy, non si tocca né il repo né la `.env`.
- Il file sul server **non si riscarica a ogni modifica del repo**: si aggiorna a mano quando serve. La copia nel repo è la versione di riferimento.

### `.env`

- Il compose la legge con `env_file` (`required: false`, quindi il container parte anche senza).
- Il `.env` **non sta nel repo** (`.gitignore` + `.dockerignore`); nel repo c'è solo il template `deploy/.env.example`.
- I segreti stanno **solo** lì: mai nel compose, mai nell'immagine, mai nella storia di git.
- Una variabile non ancora letta dall'app resta **commentata** nel template.

Variabili: `LYRICA_ADDR` (attiva), `UMAMI_URL`/`UMAMI_SITE_ID` (FASE 5), `TELEGRAM_BOT_TOKEN`/`TELEGRAM_CHAT_ID` (FASE 6).

**Il dominio pubblico non è una variabile**: `siteBaseURL` è una **costante in `internal/build/urls.go`** (`https://lyrica.filippomoscatelli.com`), usata dal feed, dalla sitemap, da `robots.txt` e dai link hreflang (D103). Quindi `LYRICA_BASE_URL` **non va introdotta** in `.env.example`: sarebbe una configurazione che nessuno legge, e il dominio resterebbe in due posti.

## Sicurezza

### Il vantaggio di partenza

Il sito **non ha login, non ha account, non ha database e non ha un pannello admin**. Non esiste una superficie di autenticazione da attaccare, non esistono credenziali utente da rubare, non esiste un archivio da esfiltrare. È la misura di sicurezza più efficace del progetto e va **difesa**, non erosa: ogni proposta di aggiungere autenticazione o storage è una decisione esplicita da registrare in `DECISIONS.md`.

### Superficie reale

1. **Contenuti** (scritti dall'autore, non da utenti): se un testo contenesse markup, il rischio è XSS verso i lettori.
2. **I due form**: abuso, spam, flood — **previsti, non implementati**.
3. **Token del bot Telegram**: se trapelato, qualcuno può scrivere nella chat dell'autore — **previsto, non implementato**.
4. **Endpoint di ricerca**: abuso leggero (scraping dell'indice) — **previsto, non implementato**.
5. **Dipendenze**: Bootstrap vendored, moduli Go, immagini base Docker.

### Contromisure

| Area | Misura |
|---|---|
| Output | **Escaping di default** in Templ; mai HTML grezzo dai contenuti; Markdown renderizzato con sanitizzazione |
| Header | HSTS, `X-Content-Type-Options: nosniff`, `Referrer-Policy: strict-origin-when-cross-origin`, `Permissions-Policy` minimale, **CSP senza `unsafe-inline`** |
| Tema | lo script che applica il tema è un **file esterno** (`assets/js/theme.js`), non inline: è così che la CSP resta senza `unsafe-inline` |
| Form | **Rate-limit per IP** in memoria + **honeypot** + tempo minimo di compilazione; nessun salvataggio dei dati — **previsto, non implementato** |
| Segreti | Token solo in `.env`/variabili d'ambiente, **mai** nel repo, nei log o nei messaggi di errore |
| Dipendenze | `.github/dependabot.yml` obbligatorio, aggiornato nella stessa PR che aggiunge un manifest |
| Esposizione | Dietro **Cloudflare Tunnel**: nessuna porta aperta in casa, nessun IP del home-lab esposto |
| Log | Log di accesso minimali, **nessun contenuto dei form**, nessuna retention lunga degli IP |

Nota: la CSP attuale (`internal/web/headers.go` → `securityHeaders`) è
`default-src 'self'; img-src 'self' data:; style-src 'self'; script-src 'self'; base-uri 'none';
form-action 'self'; frame-ancestors 'none'`. **Non contiene HSTS**: si aggiunge quando il sito è
realmente online dietro Cloudflare.

### Cosa NON si fa

- Non si aggiunge autenticazione "per sicurezza": qui l'autenticazione **crea** superficie, non la riduce.
- Non si salvano i dati dei form "tanto per averli".
- Non si aprono porte del home-lab per far passare un servizio: si usa il tunnel.
- Non si commetta un `.env`, nemmeno "temporaneamente".

### In caso di abuso

1. Il rate-limit risponde **429** senza spiegare la soglia esatta — **previsto, non implementato**.
2. Se lo spam è sistematico: si disattiva temporaneamente l'endpoint e si valuta un captcha leggero **come decisione**, non come patch silenziosa.
3. Se il token Telegram è sospetto: si rigenera dal BotFather e si aggiorna la `.env` sul server.

## Trappole già pagate

1. **Nei file `.templ` non si importa `github.com/a-h/templ`**: il generatore lo importa da sé, un import esplicito produce `templ redeclared in this block`. Se serve un tipo di templ (es. `templ.SafeURL`), lo si costruisce in un file Go (`internal/render/model.go`) come metodo di `PageData`.
2. **`golang:*-alpine` non ha bash**: lo stage di build deve fare `apk add --no-cache curl bash`, altrimenti lo script asset muore con `env: bash: No such file or directory`.
3. **Lo script asset si invoca come `bash scripts/fetch-assets.sh`**: senza il bit di esecuzione, `./scripts/...` fallisce con **exit 126** e fa fallire la build dell'immagine.
4. **`*_templ.go` non si committano mai**: sono artefatti di generazione.
5. **La porta host va scelta guardando cosa è già occupato** (`docker ps --format '{{.Names}} - {{.Ports}}'` o `ss -tlnp`), non copiando un default: sul master la **8097** è la porta scelta perché 8090-8096 sono occupate.
6. **Anche i comandi che non generano nulla richiedono `templ generate`**: `lyrica valida` compila `internal/web` → `internal/render`, quindi un job CI che lo invoca deve generare i template prima di eseguirlo.
7. **I file di output non-HTML non passano da `writeFile`**: `writeFile` accetta solo componenti templ e va avanti a `index.html`. Feed, sitemap e robots passano da `writeRaw` (`internal/build/raw.go`).
8. **La sitemap elenca le pagine che esistono davvero**: una voce di tracklist con `status: pending` non ha pagina, quindi non va in sitemap. Lo stesso vale per i contatori: usano `publishedAlbums` e `HasTranslations`, gli stessi criteri con cui il build decide se generare una pagina.
9. **`PageData.Description` è calcolato ma non stampato**: `internal/render/layout.templ` non emette `<meta name="description">` né tag `og:`. Se servono, vanno aggiunti nel layout.

## Limiti noti

- La ricerca (quando arriverà) caricherà l'indice in memoria: accettabile fino a decine di migliaia di versi.
- La ricerca è **full-text**, deliberatamente: niente semantica, niente tolleranza ai refusi.
- I contenuti richiedono un **nuovo build + deploy**: non si pubblica nulla "a caldo".
- Nessun ambiente di staging: `main` è produzione.
- Il validatore **non** confronta le chiavi dei locale fra lingue: il solo controllo è che una lingua di traduzione abbia il suo file `locales/<lang>.yaml` (`interfaceLangs`, in `internal/build/build.go`). Le chiavi mancanti si vedono solo a runtime, come fallback italiano o errore di `MustT`.
- Il dominio è una **costante**: il sito va pubblicato allo stesso indirizzo, e un dominio diverso richiede una modifica al codice, non una configurazione.
- **La CSP non ha HSTS** (vedi "Contromisure"): si aggiunge con la messa online.

## Dove vive cosa

| Area | Doc | Codice |
|---|---|---|
| La mappa richiesta → file | `REFERENCE.md` | — |
| Regole di lavoro | `GUIDELINES.md` | — |
| Decisioni | `DECISIONS.md` | — |
| Pagine e funzioni | `SPEC.md` | `internal/render/` |
| Schema contenuti e authoring | `CONTENT.md` | `internal/content/`, `cmd/lyrica/` |
| Interfaccia, brano, lingue | `FRONTEND.md` | `internal/render/`, `internal/i18n/`, `assets/` |
| Ads, analytics, legali, form | `FEATURES.md` | **non implementati**: nessun `ads.yaml`, nessun `internal/notify/` |
| Output generati dal build (feed, sitemap, robots, URL base) | `FEATURES.md` (cap. "Feed RSS e contatori") | `internal/build/rss.go`, `sitemap.go`, `robots.go`, `stats.go`, `urls.go`, `raw.go` |
| Server e sicurezza | `ARCHITECTURE.md` (cap. "Sicurezza") | `internal/web/headers.go`, `routes.go`, `handlers.go` |

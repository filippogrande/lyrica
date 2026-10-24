# Lyrica

Sito web che raccoglie **traduzioni di testi musicali** (lyric) di qualsiasi band, con **vista a fronte**: testo originale a sinistra, traduzione selezionabile a destra.

- **Produzione**: https://lyrica.filippomoscatelli.com (non ancora online)
- **Licenza codice**: MIT — **Licenza traduzioni**: CC BY-NC 4.0
- **Nessun login, nessun account, nessun database**: il sito è **statico**, generato dai contenuti che vivono in questo repo.

## Stato

**FASE 0 (documentazione) e FASE 1 (fondamenta) chiuse.** Esiste lo scheletro: server Go in
container, routing di base, layout con header/footer, tema chiaro/scuro, controllo dimensione
testo. L'immagine è pubblicata su Docker Hub da GitHub Actions.

**I contenuti ci sono**: 18 band con album e brani in `content/`, e `lyrica build` genera il
sito da quei contenuti. Le pagine home, bands, band, album e brano esistono; feed RSS, sitemap
e robots.txt sono generati dal build.

Piano e criteri di chiusura → [`docs/ROADMAP.md`](docs/ROADMAP.md). **Attenzione**: alcune fasi
sono dichiarate aperte in `ROADMAP.md` mentre l'artefatto esiste già — vedi i trovati in
[`docs/REFERENCE.md`](docs/REFERENCE.md#trovati-e-limiti-noti).

## Avvio in locale

Serve **Go 1.27.1** e la CLI **templ**. Lo script degli asset si invoca con `bash`, non con `./`
(nel repo non ha il bit di esecuzione).

```
bash scripts/fetch-assets.sh                               # scarica Bootstrap 5.3.8 (una volta)
go mod tidy
go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate   # dopo ogni modifica ai file .templ
go run ./cmd/lyrica build                                 # valida i contenuti e genera public/
go run ./cmd/lyrica serve                                 # http://localhost:8080
```

Senza `templ generate` il progetto **non compila**. I file generati (`*_templ.go`) non si
committano. Comandi, vincoli e paletti → [`AGENTS.md`](AGENTS.md).

## Deploy con Docker

L'immagine è su Docker Hub (`filippogrande/lyrica`); compose e template di configurazione stanno
in **`deploy/`**.

```
docker compose pull
docker compose up -d
```

Il sito risponde sulla **porta host 8097** (8080 è la porta interna al container). Se è occupata,
si cambia **solo quel numero** nel compose.

Sul home-lab il compose e la sua `.env` vivono in `/mnt/applicazioni/yml/docker/lyrica/`. Dettagli,
verifica e rollback → [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md).

L'immagine viene buildata e pubblicata **automaticamente** a ogni push su `main` (tag `latest` e
`<sha>`).

## Stack in una riga

Go + Templ per l'HTML, HTMX per le interazioni, Bootstrap 5 CSS-only per il layout, contenuti in
Markdown/YAML nel repo, **build-time static site**, container Docker dietro Cloudflare Tunnel.

## Struttura del repo

```
cmd/lyrica/     entrypoint: CLI (build, serve, valida, new, migrate)
internal/       build, content, i18n, render, web
content/        band, album, brani e traduzioni (Markdown + YAML)
covers/         copertine 600x600 webp, nome = slug album o band
locales/        stringhe dell'interfaccia per lingua
assets/         css (custom + vendor) e js
scripts/        fetch-assets.sh, fetch-images.sh, fetch-covers.sh, declare-image.py
deploy/         docker-compose.yml e .env.example
docs/           documentazione (indice sotto)
Dockerfile      immagine multi-stage
```

## Indice della documentazione

**Parti da [`AGENTS.md`](AGENTS.md)** — il file letto sempre — e da
**[`docs/REFERENCE.md`](docs/REFERENCE.md)** — la mappa richiesta → file → simbolo.

| File | Tipo | Cosa contiene |
|---|---|---|
| [`AGENTS.md`](AGENTS.md) | how-to | file letto sempre: identità, comandi, paletti, indice |
| [`docs/REFERENCE.md`](docs/REFERENCE.md) | reference | la mappa: richiesta → file → entry point |
| [`docs/GUIDELINES.md`](docs/GUIDELINES.md) | how-to | regole di lavoro, workflow git, paletti non negoziabili |
| [`docs/DECISIONS.md`](docs/DECISIONS.md) | reference | registro delle decisioni (D01→), con il perché |
| [`docs/ROADMAP.md`](docs/ROADMAP.md) | reference | fasi, artefatti, criteri di chiusura |
| [`docs/SPEC.md`](docs/SPEC.md) | reference | cosa fa il sito: pagine, funzionalità, casi limite |
| [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) | explanation | come è fatto: componenti, build, Docker, deploy, sicurezza |
| [`docs/CONTENT.md`](docs/CONTENT.md) | how-to | schema dei contenuti e come si aggiungono |
| [`docs/FRONTEND.md`](docs/FRONTEND.md) | explanation | interfaccia, vista brano, lingue, ricerca, PWA |
| [`docs/FEATURES.md`](docs/FEATURES.md) | explanation | ads, analytics, legali, form, feed e contatori |

Un documento che non descrive più la realtà va **corretto o cancellato nella stessa PR**; a fase
chiusa, il progetto della fase si assorbe o si elimina.

## Segnalazioni

Chi vuole proporre un brano da tradurre userà il form sul sito: la segnalazione arriva via **bot
Telegram** a Filippo. Non serve account. Il form è **previsto in FASE 6** (non ancora
implementato): vedi [`docs/FEATURES.md`](docs/FEATURES.md).

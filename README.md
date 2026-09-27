# Lyrica

Sito web che raccoglie **traduzioni di testi musicali** (lyric) di qualsiasi band, con **vista a fronte**: testo originale a sinistra, traduzione selezionabile a destra.

- **Produzione**: https://lyrica.filippomoscatellis.com (non ancora online)
- **Licenza codice**: MIT — **Licenza traduzioni**: CC BY-NC 4.0
- **Nessun login, nessun account, nessun database**: il sito è **statico**, generato dai file di contenuto che vivono in questo repo.

## Stato

**FASE 0 (documentazione) e FASE 1 (fondamenta) chiuse.**

Esiste lo scheletro: server Go in container, routing di base, layout con header/footer, tema chiaro/scuro, controllo dimensione testo. L'immagine è pubblicata su Docker Hub da GitHub Actions.

**Non esistono ancora contenuti**: nessuna band, nessun brano, nessuna pagina pubblica oltre a home e 404. La FASE 2 (schema, CLI, validazione, prime band) è la prossima.

Piano e criteri di chiusura → [`docs/ROADMAP.md`](docs/ROADMAP.md)

## Avvio con Docker

L'immagine è su Docker Hub (`filippogrande/lyrica`); compose e template di configurazione stanno in **`deploy/`**.

```
docker compose pull
docker compose up -d
```

Il sito risponde su **http://localhost:8085** (8080 è la porta interna al container).

Sul home-lab il compose e la sua `.env` vivono in `/mnt/applicazioni/yml/docker/lyrica/`. Dettagli, verifica e rollback → [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md).

L'immagine viene buildata e pubblicata **automaticamente** a ogni push su `main` (tag `latest` e `<sha>`).

## Avvio in locale senza Docker

Serve **Go 1.27.1** e la CLI **templ**. Lo script degli asset si invoca con `bash`, non con `./` (nel repo non ha il bit di esecuzione).

```
bash scripts/fetch-assets.sh                               # scarica Bootstrap 5.3.8 (una volta)
go mod tidy
go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate   # dopo ogni modifica ai file .templ
go run ./cmd/lyrica serve                                  # http://localhost:8080
```

Senza `templ generate` il progetto **non compila**. I file generati (`*_templ.go`) non si committano.

## Stack in una riga

Go + Templ per l'HTML, HTMX per le interazioni, Bootstrap 5 CSS-only per il layout, contenuti in Markdown/YAML nel repo, **build-time static site**, container Docker dietro Cloudflare Tunnel.

## Struttura del repo

```
.
├── cmd/lyrica/          # entrypoint: CLI e server
├── internal/            # i18n, render (Templ), web (server e routing)
├── locales/             # stringhe dell'interfaccia per lingua
├── assets/              # css (custom + vendor) e js
├── scripts/             # fetch-assets.sh
├── content/             # (FASE 2) band, album, brani e traduzioni
├── covers/              # (FASE 2) copertine 600x600 webp, nome = slug album
├── deploy/              # docker-compose.yml e .env.example
├── docs/                # documentazione (vedi sotto)
├── Dockerfile           # immagine multi-stage
└── ads.yaml             # (FASE 6) configurazione slot pubblicitari
```

## Documentazione

Otto documenti in [`docs/`](docs/), scritti perché **un agente AI possa capire e modificare il sito leggendoli**:

| File | Cosa contiene |
|---|---|
| [`GUIDELINES.md`](docs/GUIDELINES.md) | Regole di lavoro, workflow git, paletti non negoziabili |
| [`DECISIONS.md`](docs/DECISIONS.md) | Registro delle decisioni (D01→), con il perché |
| [`ROADMAP.md`](docs/ROADMAP.md) | Fasi di sviluppo, artefatti e criteri di chiusura |
| [`SPEC.md`](docs/SPEC.md) | Cosa fa il sito: pagine, funzionalità, casi limite |
| [`ARCHITECTURE.md`](docs/ARCHITECTURE.md) | Componenti, build, Docker, deploy, sicurezza |
| [`CONTENT.md`](docs/CONTENT.md) | Schema dei contenuti e come si aggiungono |
| [`FRONTEND.md`](docs/FRONTEND.md) | Design system, pagina brano, lingue, ricerca, PWA |
| [`FEATURES.md`](docs/FEATURES.md) | Ads, analytics, legali e privacy, form, feed e contatori |

Un documento che non descrive più la realtà va **corretto o cancellato nella stessa PR**; a fase chiusa, il progetto della fase si assorbe o si elimina.

## Segnalazioni

Chi vuole proporre un brano da tradurre userà il form sul sito: la segnalazione arriva via **bot Telegram** a Filippo. Non serve account.

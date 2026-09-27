# Lyrica

Sito web che raccoglie **traduzioni di testi musicali** (lyric) di qualsiasi band, con **vista a fronte**: testo originale a sinistra, traduzione selezionabile a destra.

- **Produzione**: https://lyrica.filippomoscatelli.com (non ancora online)
- **Licenza codice**: MIT — **Licenza traduzioni**: CC BY-NC 4.0
- **Nessun login, nessun account, nessun database**: il sito è **statico**, generato dai file di contenuto che vivono in questo repo.

## Stato del progetto

**FASE 0 (documentazione) chiusa. FASE 1 (fondamenta) in corso.**

Esiste lo scheletro: server Go in container, routing di base, layout con header/footer, tema chiaro/scuro e controllo dimensione testo. **Non esistono ancora contenuti**: nessuna band, nessun brano, nessuna pagina pubblica oltre a home e 404.

- Piano e criteri di chiusura per fase → [`ROADMAP.md`](ROADMAP.md)
- Decisioni prese (D01–D73) → [`DECISION.md`](DECISION.md)
- Cosa fa il sito → [`SPEC.md`](SPEC.md)
- Come è fatto → [`ARCHITECTURE.md`](ARCHITECTURE.md)

## Avvio rapido con Docker

L'immagine è su Docker Hub (`filippogrande/lyrica`) e il `docker-compose.yml` è in questo repo.

```
docker compose pull
docker compose up -d
```

Il sito risponde su **http://localhost:8085** (porta host; la 8080 è quella interna al container).

Compose usato sul home-lab: `/mnt/applicazioni/yml/docker/lyrica/docker-compose.yml` — dettagli, verifica e rollback in [`DEPLOY.md`](DEPLOY.md).

L'immagine viene buildata e pubblicata **automaticamente** da GitHub Actions a ogni push su `main` (tag `latest` e `<sha>`).

## Avvio in locale senza Docker

Serve **Go 1.27.1** e la CLI **templ**. Nota: lo script degli asset si invoca con `bash`, non con `./` (nel repo non ha il bit di esecuzione).

```
bash scripts/fetch-assets.sh                               # scarica Bootstrap 5.3.8 (una volta)
go mod tidy
go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate   # dopo ogni modifica ai file .templ
go run ./cmd/lyrica serve                                  # http://localhost:8080
```

Senza `templ generate` il progetto **non compila**. I file generati (`*_templ.go`) non si committano. Dettagli in [`BUILD_PIPELINE.md`](BUILD_PIPELINE.md).

Altri comandi: `lyrica help`, `lyrica build` e `lyrica new` (arrivano in FASE 2; per ora restituiscono un errore esplicito).

## Stack in una riga

Go + Templ per l'HTML, HTMX per le interazioni, Bootstrap 5 CSS-only per il layout, contenuti in Markdown/YAML nel repo, **build-time static site**, container Docker dietro Cloudflare Tunnel.

## Struttura del repo

```
.
├── cmd/lyrica/          # entrypoint: CLI e server
├── internal/
│   ├── i18n/            # lingue dell'interfaccia, negoziazione, stringhe
│   ├── render/          # componenti Templ (layout, header, footer, pagine)
│   └── web/             # server HTTP: routing, handler, header di sicurezza
├── locales/             # stringhe UI per lingua (it.yaml, ...)
├── assets/
│   ├── css/             # CSS custom + vendor/ (Bootstrap, scaricato, non committato)
│   └── js/              # theme.js (tema prima del paint), lyrica.js (interazioni)
├── scripts/             # fetch-assets.sh
├── content/             # (FASE 2) band, album, brani e traduzioni
├── covers/              # (FASE 2) copertine 600x600 webp, nome = slug album
├── Dockerfile           # immagine multi-stage (build + runtime minimale)
├── docker-compose.yml   # avvio sul home-lab
├── ads.yaml             # (FASE 6) configurazione slot pubblicitari
└── *.md                 # questa documentazione
```

## Documentazione

L'idea è che **un agente AI possa capire e modificare il sito leggendo solo questi file**, senza dover ricostruire il contesto dal codice.

| File | Cosa contiene |
|---|---|
| [`DEVELOPMENT_GUIDELINES.md`](DEVELOPMENT_GUIDELINES.md) | Regole di sviluppo, workflow git, paletti non negoziabili |
| [`DECISION.md`](DECISION.md) | Registro delle decisioni prese (con il perché) |
| [`SPEC.md`](SPEC.md) | Cosa fa il sito: pagine, funzionalità, contenuti |
| [`ARCHITECTURE.md`](ARCHITECTURE.md) | Componenti, flusso di build, struttura cartelle |
| [`CONTENT_SCHEMA.md`](CONTENT_SCHEMA.md) | Schema dei contenuti (band/album/brano/traduzioni) |
| [`CONTENT_AUTHORING.md`](CONTENT_AUTHORING.md) | Come si aggiungono contenuti, passo per passo |
| [`BUILD_PIPELINE.md`](BUILD_PIPELINE.md) | CLI, validazione dei contenuti, build statica |
| [`DESIGN_SYSTEM.md`](DESIGN_SYSTEM.md) | Layout, tema, tipografia, accessibilità |
| [`TRACK_VIEW_DESIGN.md`](TRACK_VIEW_DESIGN.md) | La vista a fronte: strofe, cantante, selezione lingua |
| [`I18N_DESIGN.md`](I18N_DESIGN.md) | Lingue dell'interfaccia, URL, hreflang |
| [`SEARCH_DESIGN.md`](SEARCH_DESIGN.md) | Ricerca full-text e dropdown live |
| [`PWA_DESIGN.md`](PWA_DESIGN.md) | Installabilità e lettura offline |
| [`FEED_AND_STATS.md`](FEED_AND_STATS.md) | Feed RSS e contatori pubblici |
| [`ADS_DESIGN.md`](ADS_DESIGN.md) | Slot pubblicitari (spenti al lancio) |
| [`ANALYTICS_DESIGN.md`](ANALYTICS_DESIGN.md) | Umami: eventi e privacy |
| [`LEGAL_DESIGN.md`](LEGAL_DESIGN.md) | Disclaimer, licenze, GDPR, pagine legali |
| [`SUGGESTIONS_DESIGN.md`](SUGGESTIONS_DESIGN.md) | Segnalazioni brani + form contatto + notifica Telegram |
| [`SECURITY.md`](SECURITY.md) | Modello di minaccia di un sito senza login |
| [`DEPLOY.md`](DEPLOY.md) | Docker, Cloudflare Tunnel, GitHub Actions |
| [`ROADMAP.md`](ROADMAP.md) | Piano di sviluppo a fasi |

## Segnalazioni

Chi vuole proporre un brano da tradurre userà il form sul sito: la segnalazione arriva via **bot Telegram** a Filippo. Non serve account.

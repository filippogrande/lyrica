# Lyrica

Sito web che raccoglie **traduzioni di testi musicali** (lyric) di qualsiasi band, con **vista a fronte**: testo originale a sinistra, traduzione selezionabile a destra.

- **Produzione**: https://lyrica.filippomoscatelli.com
- **Licenza codice**: MIT — **Licenza traduzioni**: CC BY-NC 4.0
- **Nessun login, nessun account, nessun database**: il sito è **statico**, generato dai file di contenuto che vivono in questo repo.

## Stato del progetto

**Progettazione chiusa, sviluppo non iniziato.** Nessuna riga di codice è stata scritta: questo repo contiene per ora solo la documentazione che definisce regole e funzionamento.

- Decisioni chiuse → [`DECISION.md`](DECISION.md)
- Cosa fa il sito → [`SPEC.md`](SPEC.md)
- Come è fatto → [`ARCHITECTURE.md`](ARCHITECTURE.md)
- In che ordine si sviluppa → [`ROADMAP.md`](ROADMAP.md)

## Stack in una riga

Go + Templ per l'HTML, HTMX per le interazioni, Bootstrap 5 CSS-only per il layout, contenuti in Markdown/YAML nel repo, **build-time static site**, deploy su Docker nel home-lab dietro Cloudflare Tunnel.

## Struttura del repo (prevista)

```
.
├── content/            # contenuti: band, album, brani e traduzioni (MD + front-matter YAML)
├── assets/            # CSS custom, font, icone
│   └── css/
├── covers/            # copertine album (600x600 webp, nome = slug album)
├── cmd/               # entrypoint del binario Go (server + CLI `lyrica`)
├── internal/          # moduli Go (content loader, render, search, i18n, notify)
├── locales/           # stringhe UI per lingua
├── ads.yaml           # configurazione slot pubblicitari (spenti al lancio)
├── Dockerfile
├── docker-compose.yml
└── *.md               # questa documentazione
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

Chi vuole proporre un brano da tradurre usa il form sul sito: la segnalazione arriva via **bot Telegram** a Filippo. Non serve account.

# ARCHITECTURE — come è fatto Lyrica

## In una riga

Sito **statico generato a build time** da contenuti Markdown, servito da un binario Go che espone solo tre cose dinamiche: la ricerca, i due form e il feed.

## Flusso

```
content/  covers/  assets/  locales/
        |
        |  lyrica build   (parse + validazione + render Templ)
        v
     public/            HTML, CSS, JS minimo, cover, indice di ricerca,
                        rss.xml, sitemap.xml, robots.txt, manifest, sw.js
        |
        |  docker build (multi-stage: build Go -> immagine finale minimale)
        v
   immagine Docker  ->  docker compose  ->  Cloudflare Tunnel  ->  lyrica.filippomoscatelli.com
```

## Componenti

| Componente | Responsabilità | Dove |
|---|---|---|
| **Content loader** | legge `content/`, parse di front-matter + strofe, costruisce i modelli | `internal/content/` |
| **Validator** | regole di integrità: front-matter, slug, link, traduzioni incomplete, cover | `internal/content/` |
| **Renderer** | pagine Templ (home, band, album, brano, 404, legali) | `internal/render/` |
| **Search index** | costruisce l'indice (versi, titoli, band) a build time e lo interroga | `internal/search/` |
| **i18n** | stringhe UI e negoziazione lingua | `internal/i18n/`, `locales/` |
| **Notifier** | invio messaggi al bot Telegram (segnalazioni, contatti) | `internal/notify/` |
| **Web** | server, routing, endpoint HTMX, header di sicurezza, rate-limit | `internal/web/` |
| **CLI** | `lyrica new ...`, `lyrica build`, `lyrica serve` | `cmd/lyrica/` |

Il codice sta in inglese, i contenuti in italiano (vedi `DEVELOPMENT_GUIDELINES.md`).

## Cosa gira a runtime

Il container serve soltanto:

1. **File statici**: tutte le pagine già generate (`public/`).
2. **`GET /it/api/cerca?q=`**: ricerca nell'indice, restituisce un **frammento HTML** per HTMX.
3. **`POST /it/api/segnala`** e **`POST /it/api/contatti`**: validano, rate-limitano e inoltrano al bot Telegram; **non salvano nulla**.

Non esiste database, non esiste sessione, non esiste stato utente. Conseguenze volute: nessuna migrazione, nessun backup di dati applicativi, nessun dato personale a riposo.

## Scelte architetturali e alternative scartate

| Scelta | Alternativa scartata | Perché |
|---|---|---|
| Build-time static | SSR a runtime con database | niente DB da mantenere, niente GDPR sui dati a riposo, il sito è sempre identico e cacheabile |
| Go + Templ | Node/Next, PHP | build veloce, un binario, nessun runtime JS da servire |
| HTMX | SPA con framework JS | un solo endpoint dinamico non giustifica un framework |
| Bootstrap 5 CSS-only | Tailwind con build npm | nessun passaggio npm per servire il sito |
| Docker + Cloudflare Tunnel | k3s | k3s in dismissione nel home-lab |
| Bot Telegram | pannello admin / email | nessun login, nessuna deliverability da gestire |

## Dove vive cosa

| Area | Doc | Codice |
|---|---|---|
| Schema dei contenuti | `CONTENT_SCHEMA.md` | `internal/content/` |
| Come si aggiunge contenuto | `CONTENT_AUTHORING.md` | `cmd/lyrica/` (CLI) |
| Build e validazione | `BUILD_PIPELINE.md` | `internal/content/`, `cmd/lyrica/` |
| Interfaccia e tema | `DESIGN_SYSTEM.md` | `internal/render/`, `assets/css/` |
| Pagina brano | `TRACK_VIEW_DESIGN.md` | `internal/render/track` |
| Lingue | `I18N_DESIGN.md` | `internal/i18n/`, `locales/` |
| Ricerca | `SEARCH_DESIGN.md` | `internal/search/` |
| Feed e contatori | `FEED_AND_STATS.md` | `internal/render/feed` |
| PWA | `PWA_DESIGN.md` | `assets/`, `public/` |
| Pubblicità | `ADS_DESIGN.md` | `ads.yaml`, `internal/render/ads` |
| Analytics | `ANALYTICS_DESIGN.md` | template head |
| Form e notifiche | `SUGGESTIONS_DESIGN.md` | `internal/notify/`, `internal/web/` |
| Sicurezza | `SECURITY.md` | `internal/web/` |
| Legale | `LEGAL_DESIGN.md` | `internal/render/legal` |
| Deploy | `DEPLOY.md` | `Dockerfile`, `docker-compose.yml`, `.github/workflows/` |

## Dipendenze esterne

- **Cloudflare Tunnel** (già attivo nel home-lab) per l'esposizione pubblica.
- **Umami** self-hosted (istanza esistente) per le statistiche.
- **Bot Telegram** dedicato per le notifiche.
- **GitHub Actions** per validazione, build dell'immagine e deploy su merge in `main`.

## Vincoli e limiti noti

- Il sito **non scala su contenuti enormi per pagina**: l'indice di ricerca è caricato in memoria a runtime (accettabile fino a decine di migliaia di versi).
- Non c'è ricerca semantica: la ricerca è **full-text**, deliberatamente.
- I contenuti richiedono un **nuovo build+deploy**: non si pubblica nulla "a caldo".
- Nessun ambiente di staging al lancio: `main` è produzione (vedi `DEVELOPMENT_GUIDELINES.md`).

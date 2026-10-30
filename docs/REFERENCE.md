# REFERENCE — la mappa del repo

Questo file è **la mappa**: una richiesta ("dove vive X?") → il doc → il file → il simbolo.
Non è un doc tematico: i fatti stanno nei doc di competenza (indice in [`AGENTS.md`](../AGENTS.md)).

**Come si usa.** Leggi `AGENTS.md`, poi questa mappa, poi il file: massimo **2 salti** dalla radice.
Ogni riga punta a un file che esiste e a un simbolo che esiste, verificati uno a uno sul codice
(lettura del file + `search_code`), **non dedotti dalla doc**.

**Regole.** Una riga per richiesta. **Mai numeri di riga**: le righe si muovono, i simboli no.
**Mai elenco di directory**: un elenco non risponde a una richiesta. Se sposti un file o rinomini
un simbolo, aggiorni la mappa **nella stessa PR**.

## La mappa

### Contenuti: schema, lettura, validazione

| Richiesta | Doc | File | Simbolo / entry point |
|---|---|---|---|
| Dove vive lo schema dei contenuti (band, album, brano, blocchi, strofe)? | `docs/CONTENT.md` | `internal/content/model.go` | `Track`, `Block`, `Stanza`, `Band`, `Album`, `Catalog` |
| Dove si legge il front-matter YAML e il corpo di un file? | `docs/CONTENT.md` | `internal/content/frontmatter.go` | `readFrontMatterFile`, `splitFrontMatter` |
| Dove si leggono band, album e brani dal disco? | `docs/CONTENT.md` | `internal/content/load.go` | `Load`, `loadBand`, `loadAlbum`, `loadTracks` |
| Dove si legge un brano nel formato a cartella (`song.md` + un file per lingua)? | `docs/CONTENT.md` | `internal/content/trackfolder.go` | `loadTrackFolder`, `blockFileRole` |
| Dove si costruisce la tracklist stampata (ordine del disco, voci senza pagina)? | `docs/CONTENT.md` | `internal/content/tracklist.go` | `Album.Entries`, `TracklistEntry.Published` |
| Dove si converte un brano dal file unico alla cartella? | `docs/CONTENT.md` | `cmd/lyrica/migrate.go` | `runMigrate`, `migrateTrackFile` |
| Dove si genera lo slug di un contenuto? | `docs/CONTENT.md` | `cmd/lyrica/slug.go` | `slugify`, `uniqueSlug`; in `internal/content/slug.go` `validSlug` |
| **Dove vive la validazione dei contenuti?** | `docs/CONTENT.md` ("Regole di validazione") | `internal/content/validate.go` | `Validate` |
| Dove si applica la regola più importante (traduzione incompleta)? | `docs/CONTENT.md` | `internal/content/validate_track.go` | `checkTranslations` |
| Dove si raccolgono tutti gli errori e gli avvisi? | `docs/CONTENT.md` | `internal/content/report.go` | `Report`, `Issue`, `Report.Err` |
| Dove si lancia la validazione da riga di comando? | `docs/CONTENT.md` | `cmd/lyrica/validate.go` | `validateContent` |

### Build e output generato

| Richiesta | Doc | File | Simbolo / entry point |
|---|---|---|---|
| Dove parte la build del sito? | `docs/ARCHITECTURE.md` | `internal/build/build.go` | `Run`, `DefaultOptions` |
| Dove si decide se una pagina esiste (e quindi se un album è pubblicato)? | `docs/ARCHITECTURE.md` | `internal/build/pages.go` | `publishedAlbums`, `writeAlbum` |
| Dove si generano le pagine di una lingua (home, bands, band, album, brano, 404, feed)? | `docs/ARCHITECTURE.md` | `internal/build/pages.go` | `pageWriter.writeLang` |
| Dove si scrive un file che non è HTML (XML, `robots.txt`)? | `docs/ARCHITECTURE.md` | `internal/build/raw.go` | `pageWriter.writeRaw` |
| Dove si genera il feed RSS? | `docs/FEATURES.md` | `internal/build/rss.go` | `pageWriter.writeFeed`, `feedLimit` |
| Dove si genera la sitemap? | `docs/FEATURES.md` | `internal/build/sitemap.go` | `pageWriter.writeSitemap`, `sitemapPages` |
| Dove si genera `robots.txt`? | `docs/FEATURES.md` | `internal/build/robots.go` | `pageWriter.writeRobots`, `robotsTxt` |
| Dove si calcolano i contatori pubblici (brani, band, lingue)? | `docs/FEATURES.md` | `internal/build/stats.go` | `siteStats` |
| Dove vive il dominio pubblico (gli URL assoluti)? | `docs/ARCHITECTURE.md` | `internal/build/urls.go` | `siteBaseURL`, `absoluteURL` |
| Dove si decide quali lingue d'interfaccia vengono generate? | `docs/FRONTEND.md` | `internal/build/build.go` | `interfaceLangs` |
| Dove si copiano asset e cover dentro `public/`? | `docs/ARCHITECTURE.md` | `internal/build/assets.go` | `copyAssets`, `copyCovers` |
| Dove si genera l'indice di ricerca (per lingua)? | docs/FRONTEND.md | internal/build/search.go | writeSearchIndex |

### Rendering e interfaccia

| Richiesta | Doc | File | Simbolo / entry point |
|---|---|---|---|
| Dove si costruiscono gli URL delle pagine? | `docs/FRONTEND.md` | `internal/render/model.go` | `PageData.TrackPath`, `PageData.BandPath`, `PageData.AlbumPath` |
| Dove si costruiscono i link hreflang? | `docs/FRONTEND.md` | `internal/render/model.go` | `PageData.AlternateLinks`; lato sitemap `internal/build/sitemap.go` → `sitemapEntry` |
| Dove si scelgono bandiera ed etichetta di una lingua? | `docs/FRONTEND.md` | `internal/render/lang_flags.go` | `FlagOf`, `LangLabel` |
| Dove si costruisce la vista della pagina brano? | `docs/FRONTEND.md` | `internal/render/track_page.go` | `BuildTrackView` |
| Dove sta il markup della pagina brano (vista a fronte)? | `docs/FRONTEND.md` | `internal/render/track.templ` | `Track`, `textSide` |
| Dove si costruisce la home a corsie? | `docs/FRONTEND.md` | `internal/render/view.go` | `BuildHomeView`, `newTrackRail`, `latestBands` |
| Dove sta il markup della home? | `docs/FRONTEND.md` | `internal/render/home.templ` | `Home`; la corsia `railBlock` sta in `internal/render/rail.templ` |
| Dove si costruisce la pagina delle band e il filtro per lettera? | `docs/FRONTEND.md` | `internal/render/view.go` | `BuildBandsView`, `groupByInitial`; markup `Bands` e `bandIndex` in `internal/render/bands.templ` |
| Dove si costruisce la pagina album (tracklist completa)? | `docs/FRONTEND.md` | `internal/render/album_page.go` | `BuildAlbumView`; markup `Album` in `internal/render/album.templ` |
| Dove si costruisce la corsia degli album (riusata in pagina band e album)? | `docs/FRONTEND.md` | `internal/render/album_rail.go` | `newAlbumRail` |
| Dove sta la struttura comune di ogni pagina (head, header, footer)? | `docs/FRONTEND.md` | `internal/render/layout.templ` | `Layout` |
| Dove si scrive la riga dei contatori? | `docs/FRONTEND.md` | `internal/render/stats.templ` | `StatsLine`; i numeri in `internal/render/stats.go` → `SiteStats.Line` |
| Dove sta il campo di ricerca e il dropdown? | docs/FRONTEND.md | internal/render/searchbox.templ | SearchBox; assets/js/search.js |

### Server, runtime e i18n

| Richiesta | Doc | File | Simbolo / entry point |
|---|---|---|---|
| Dove vive il server HTTP e il routing? | `docs/ARCHITECTURE.md` | `internal/web/routes.go` | `Handler` |
| Dove si avvia il server in locale? | `docs/ARCHITECTURE.md` | `internal/web/server.go` | `Serve`, `defaultAddr` |
| Dove si applicano header di sicurezza e CSP? | `docs/ARCHITECTURE.md` | `internal/web/headers.go` | `securityHeaders` |
| Dove si risponde al 404 e al redirect della radice? | `docs/ARCHITECTURE.md` | `internal/web/handlers.go` | `handleRootRedirect`, `serveNotFound` |
| Dove si decide la lingua dell'interfaccia da `Accept-Language`? | `docs/FRONTEND.md` | `internal/i18n/i18n.go` | `Bundle.Negotiate`, `DefaultLang` |

### CLI, script, deploy, CI

| Richiesta | Doc | File | Simbolo / entry point |
|---|---|---|---|
| Dove si smista il comando `lyrica`? | `docs/ARCHITECTURE.md` | `cmd/lyrica/main.go` | `run` |
| Dove si crea un nuovo contenuto (band, album, brano)? | `docs/CONTENT.md` | `cmd/lyrica/new_band.go` | `newBand`; anche `newAlbum` (`new_album.go`), `newBrano` (`new_brano.go`) |
| Dove si serializza un front-matter YAML? | `docs/CONTENT.md` | `cmd/lyrica/new.go` | `frontMatter`, `parseArgs` |
| Dove si scarica Bootstrap (asset di terze parti)? | `docs/ARCHITECTURE.md` | `scripts/fetch-assets.sh` | `BOOTSTRAP_VERSION` |
| Dove si scaricano le immagini mancanti, o una sola? | `docs/CONTENT.md` ("Immagini") | `scripts/fetch-images.sh` | script bash; anche `scripts/fetch-covers.sh` e `scripts/declare-image.py` |
| Dove si configura il container e la porta host? | `docs/ARCHITECTURE.md` | `deploy/docker-compose.yml` | servizio `web`, `ports` (8097:8080) |
| Quali variabili d'ambiente esistono e quali sono già lette? | `docs/ARCHITECTURE.md` | `deploy/.env.example` | solo `LYRICA_ADDR` è attiva |
| Dove si costruisce l'immagine Docker? | `docs/ARCHITECTURE.md` | `Dockerfile` | stage `build`, stage runtime |
| Dove gira la CI e cosa fa? | `docs/ARCHITECTURE.md` | `.github/workflows/ci.yml` | job `compila`, job `valida` |
| Dove si pubblica l'immagine su Docker Hub? | `docs/ARCHITECTURE.md` | `.github/workflows/docker-build-push.yml` | job `build-and-push` |
| Dove stanno i workflow manuali (immagini, migrazione)? | `docs/CONTENT.md` | `.github/workflows/immagini-mancanti.yml` | `Immagini mancanti`; `immagini.yml` → `Immagini dagli archivi`; `migrazione-brani.yml` → `Migrazione brani` |
| Dove si configura Dependabot? | `docs/GUIDELINES.md` | `.github/dependabot.yml` | `gomod`, `docker`, `github-actions` |
| Dove vivono le stringhe dell'interfaccia? | `docs/FRONTEND.md` | `locales/it.yaml` | chiavi piatte (`nav.home`, `track.added`); lette da `internal/i18n/i18n.go` → `Bundle.MustT` |
| Dove gira il controllo di sicurezza (CVE dipendenze, analisi statica, immagine)? | `docs/ARCHITECTURE.md` | `.github/workflows/security.yml` | job `dipendenze`, `sast`, `immagine` |

### Regole, decisioni, piano

| Richiesta | Doc | File | Simbolo / entry point |
|---|---|---|---|
| Quali sono le regole non negoziabili del repo? | `docs/GUIDELINES.md` | `docs/GUIDELINES.md` | sezioni §1–§9 |
| Perché è stato scelto X? Dove sono le decisioni? | `docs/DECISIONS.md` | `docs/DECISIONS.md` | registro D01→ |
| Qual è il piano di fase e i criteri di chiusura? | `docs/ROADMAP.md` | `docs/ROADMAP.md` | FASE 0 → FASE 6 |
| Cosa fa il sito, pagina per pagina? | `docs/SPEC.md` | `docs/SPEC.md` | "Mappa delle pagine" |
| Come è fatto il sito (componenti, build, Docker, sicurezza)? | `docs/ARCHITECTURE.md` | `docs/ARCHITECTURE.md` | "Componenti", "Sicurezza" |

## Non esiste (ancora): dove NON cercare

Richieste a cui la doc risponde come se la cosa fosse già fatta, ma **il codice non c'è**.
Cercarle nel repo porta a conclusioni sbagliate: la risposta giusta è "previsto, non implementato".

| Richiesta | Cosa dice la doc | Realtà |
|---|---|---|
| Ricerca full-text: `GET /it/api/cerca`, `public/search-index.json` | `docs/SPEC.md`, `docs/FRONTEND.md` ("Ricerca full-text") | nessun `internal/search/`; `internal/web/routes.go` registra solo `/`, `/healthz` e il file server |
| Form segnalazione/contatti: `POST /it/api/segnala`, `POST /it/api/contatti`, `internal/notify/` | `docs/FEATURES.md`, `docs/SPEC.md` | nessun pacchetto `internal/notify/`, nessun handler di form |
| Ads: `ads.yaml`, slot `track_top` / `side_left` / `side_right` / `track_bottom` | `docs/FEATURES.md`, `docs/FRONTEND.md` | nessun `ads.yaml` nel repo, nessuno slot in `internal/render/track.templ` |
| PWA: `manifest.webmanifest`, `sw.js` | `docs/FRONTEND.md` ("PWA") | il build non li scrive: `internal/build/` non li genera |
| Redirect 301 per uno slug rinominato | `docs/SPEC.md`, `docs/CONTENT.md` | nessun file di redirect generato né letto |
| Pagine `/it/tag/{tag}`, `/it/segnala`, `/it/contatti`, `/it/legali` | `docs/SPEC.md` ("Mappa delle pagine") | il build scrive solo home, bands, band, album, brano e 404 (`pageWriter.writeLang`) |
| `meta description` e tag `og:` per pagina | `docs/FRONTEND.md` ("Cosa sta nel `<head>`") | `internal/render/layout.templ` non li stampa: `PageData.Description` è calcolato ma non renderizzato |

## Stato di implementazione

La doc descrive il **progetto** (FASE 0 parte prima del codice, `docs/ROADMAP.md`). Non tutto
quello che descrive esiste: ecco cosa c'è **davvero** oggi.

| Area | Stato | Dove |
|---|---|---|
| Contenuti in `content/` | **esiste**: 18 band, album, brani (formato a file unico) | `internal/content/load.go` |
| Validazione in CI | **esiste** | `internal/content/validate.go`, job `valida` |
| Build statico → `public/` | **esiste** | `internal/build/build.go` |
| Pagine home, bands, band, album, brano, 404 | **esistono** | `internal/render/*` |
| Feed RSS, sitemap, robots, contatori | **esistono** | `internal/build/rss.go`, `sitemap.go`, `robots.go`, `stats.go` |
| Redirect radice → lingua negoziata, `/healthz` | **esistono** | `internal/web/handlers.go` |
| Formato contenuti a cartella (`song.md` + un file per lingua) | **esiste come destinazione**: la CLI e la migrazione lo producono, il loader lo legge; i contenuti in repo sono ancora nel formato precedente | `internal/content/trackfolder.go`, `cmd/lyrica/migrate.go` |
| Ricerca full-text, form/Telegram, ads, PWA, redirect 301, pagine tag/segnala/contatti/legali | **previsti, non implementati** | vedi la sezione precedente |
| Umami, pagine legali pubblicate, messa online | **previsti (FASE 5), non implementati** | `deploy/.env.example` (variabili commentate) |

## Trovati e limiti noti

Cose che la doc dice e il codice non conferma. **Non sono state "sistemate" nel codice**:
questa PR è solo documentazione; le correzioni ai doc sono in questa stessa PR, quelle al
codice sono task a parte.

1. **`internal/search/` e `internal/notify/` non esistono**, ma `docs/ARCHITECTURE.md` li elenca
   come componenti. Corretto in questa PR.
2. **`ads.yaml` non esiste**, ma `README.md`, `docs/FEATURES.md` e `docs/ARCHITECTURE.md` lo danno
   per presente. Corretto in questa PR.
3. **`lyrica build` e `lyrica new ...` sono disponibili**, ma `docs/ARCHITECTURE.md` li marca
   "FASE 2". Corretto in questa PR; il comando `lyrica migrate` non era elencato ed è stato aggiunto.
4. **La porta host è 8097** (`deploy/docker-compose.yml`), non 8085: `README.md` la diceva 8085 e
   `docs/ARCHITECTURE.md` ("trappole già pagate") parlava di 8085. Corretto in questa PR.
5. **I contenuti esistono**, ma `README.md` diceva "non esistono ancora contenuti" e `docs/ROADMAP.md`
   dichiara aperta la FASE 2. `README.md` corretto in questa PR; **`ROADMAP.md` non si tocca** (vedi punto 6).
6. **`docs/ROADMAP.md` è in parte morto**: dichiara aperte FASE 2, 3 e 4 mentre gli artefatti
   esistono (contenuti, pagine, feed/sitemap/robots). **Non lo si cancella né lo si riscrive qui**:
   va discusso in questa PR e poi aggiornato nella PR che chiude la fase.
7. **I commenti nel codice citano doc che non esistono**: `SECURITY.md`, `DECISION.md`,
   `I18N_DESIGN.md`, `SEARCH_DESIGN.md`, `DESIGN_SYSTEM.md`, `PWA_DESIGN.md`. È codice, non doc:
   **non toccato in questa PR**; segnalato per una PR successiva.
8. **Il validatore non confronta le chiavi dei locale fra lingue**: `docs/FRONTEND.md` diceva il
   contrario. Corretto in questa PR (l'unico controllo è che una lingua di traduzione abbia un
   locale, in `internal/build/build.go` → `interfaceLangs`).
9. **`content/` contiene centinaia di file `.md`, non letti tutti**: ne sono stati letti 3
   (`band.md`, `album.md`, `tracks/` di tre band) per capire lo schema. Lo schema è stato preso
   da `docs/CONTENT.md` e verificato contro `internal/content/model.go`.
10. **Simboli verificati per lettura del file**, non con `search_code`: `publishedAlbums`,
    `interfaceLangs`, `pageWriter.writeLang`, `BuildTrackView`, `securityHeaders`,
    `handleRootRedirect`, `Serve`, `siteStats`, `writeRaw`, `FlagOf`/`LangLabel`.
    **Simboli verificati con `search_code`**: `validSlug`, `newBrano`, `lyrica migrate`,
    `song.md` (formato a cartella), `ads.yaml` (assente), `search-index.json` (assente),
    `api/cerca` (assente), `serviceWorker` (assente), `notify` in `internal/` (assente),
    `SECURITY.md`/`DECISION.md`/`I18N_DESIGN.md`/`SEARCH_DESIGN.md` (assenti).
    `search_code` non ha mai risposto "rate limit exceeded" durante questa verifica.

## Manutenzione della mappa

- **Una riga per richiesta, non per area.** Se una riga risponde a due domande diverse, si scinde.
- **Mai numeri di riga.** Se sposti una funzione, la mappa resta vera; se cambi il nome, la riga no.
- **Mai elenchi di directory.** Se una richiesta si risolve con "guarda `internal/`", la riga è inutile.
- **Ogni riga si verifica**: il file esiste e il simbolo esiste. Una riga non verificata si toglie.
- **Si aggiorna nella stessa PR** che sposta un file o rinomina un simbolo. Mappa e codice
  divergono solo per un commit.

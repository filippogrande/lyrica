# ROADMAP — piano di sviluppo

Ordine di sviluppo a fasi. Ogni fase si chiude con un **artefatto verificabile**: se l'artefatto non esiste, la fase non è chiusa.

Le task vivono su **Vikunja, progetto `Lyrica`**; qui c'è l'ordine e il criterio di chiusura.

## FASE 0 — Documentazione

**Artefatto**: il set di documenti nel repo. La documentazione **precede** il codice: è la sorgente da cui si sviluppa.

**Stato**: ✅ chiusa.

## FASE 1 — Fondamenta

**Artefatto**: il binario compila e gira **in container**; il sito mostra una home navigabile con header, footer, tema chiaro/scuro, controllo dimensione testo e 404; l'immagine è pubblicata su Docker Hub da GitHub Actions.

**Criterio di chiusura**: `docker compose pull && docker compose up -d` tira su il sito e risponde su `/healthz`; la CI che compila è verde sulla PR.

**Stato**: ✅ chiusa — codice mergiato, CI `Compila` verde, immagine pubblicata su Docker Hub (`latest` + tag per commit). Primo `up -d` sul home-lab da fare.

## FASE 2 — Contenuti

**Artefatto**: la CLI genera i template; 2-3 band reali con un album e alcuni brani sono nel repo; la validazione in CI passa; il sito mostra una band e un album veri.

**Criterio di chiusura**: `lyrica build` produce le pagine di band/album dai contenuti reali, e la validazione **fallisce** su un contenuto rotto di prova (traduzione incompleta, slug duplicato).

## FASE 3 — Pagine pubbliche

**Artefatto**: un visitatore arriva in homepage, trova una band, apre un album e legge un brano con originale a sinistra e traduzione a destra.

**Criterio di chiusura**: Home → Bands → Band → Album → Brano funziona da telefono, con breadcrumb e cambio lingua.

## FASE 4 — Funzioni trasversali

**Artefatto**: la ricerca live trova un brano cercando un verso; il sito risponde in più lingue di interfaccia; esistono `rss.xml`, `sitemap.xml`, `robots.txt` e la 404 con ricerca.

**Criterio di chiusura**: cercando un frammento di verso si arriva al brano giusto dal dropdown; il feed è valido.

## FASE 5 — Messa online e osservabilità

**Artefatto**: il sito risponde su `lyrica.filippomoscatellis.com` dietro Cloudflare Tunnel; Umami registra le visite; le pagine legali esistono e sono linkate dal footer.

Il container esiste già dalla FASE 1: qui si aggiunge ciò che manca per essere **pubblici**.

**Criterio di chiusura**: l'URL pubblico risponde e Umami mostra gli eventi.

## FASE 6 — Extra

**Artefatto**: il sito è installabile e legge i brani offline; il form di segnalazione manda un messaggio al bot Telegram.

**Criterio di chiusura**: una segnalazione inviata dal sito arriva come messaggio Telegram.

## Regole del piano

- **Non si salta la FASE 0**: il codice segue i doc, non il contrario.
- Non si inizia una fase senza che la precedente abbia il suo artefatto verificabile.
- **Ogni fase si verifica su un container reale**, non solo in locale.
- Se durante lo sviluppo una decisione cambia, si aggiorna `DECISIONS.md` **nella stessa PR**.
- **A fase chiusa si ripuliscono i documenti**: il progetto della fase o si assorbe in una sezione "come funziona" o si elimina (`GUIDELINES.md` §9, D74).

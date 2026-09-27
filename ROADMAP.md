# ROADMAP — piano di sviluppo di Lyrica

Ordine di sviluppo a fasi. Ogni fase si chiude con un artefatto verificabile: se l'artefatto non esiste, la fase non è chiusa.

Le task vivono su **Vikunja, progetto `Lyrica`** (`#1`–`#18`); qui c'è solo l'ordine e il criterio di chiusura.

## FASE 0 — Documentazione

**Task**: #13.

**Artefatto**: il set di MD nel repo (20 file). La documentazione **precede** il codice: è la sorgente da cui si sviluppa.

**Stato**: ✅ chiusa (PR #1 mergiata).

## FASE 1 — Fondamenta

**Task**: #3 (scheletro Go + Templ), #7 (design system), #4 (routing), più **Docker e CI dell'immagine** (D72).

**Artefatto**: il binario compila e gira **in container**; il sito mostra una home navigabile con header, footer, tema chiaro/scuro, controllo dimensione testo e 404; l'immagine è pubblicata su Docker Hub da GitHub Actions.

**Criterio di chiusura**: `docker compose pull && docker compose up -d` sul home-lab tira su il sito e risponde su `/healthz`; la CI che compila è verde sulla PR.

**Stato**: in corso (PR #2).

## FASE 2 — Contenuti

**Task**: #2 (schema), #18 (CLI), #9 (pipeline + validazione), #1 (primi contenuti reali).

**Artefatto**: la CLI genera i template; 2-3 band reali con un album e alcuni brani sono nel repo; la validazione in CI passa; il sito mostra una band e un album veri.

**Criterio di chiusura**: `lyrica build` produce le pagine di band/album dai contenuti reali, e la validazione **fallisce** su un contenuto rotto di prova (traduzione incompleta, slug duplicato).

## FASE 3 — Pagine pubbliche

**Task**: #10 (vista a fronte) + le pagine Home / Bands / Band / Album.

**Artefatto**: un visitatore arriva in homepage, trova una band, apre un album e legge un brano con originale a sinistra e traduzione a destra, con il selettore lingua.

**Criterio di chiusura**: la navigazione completa Home → Bands → Band → Album → Brano funziona da telefono, con breadcrumb e cambio lingua.

## FASE 4 — Funzioni trasversali

**Task**: #5 (ricerca), #6 (i18n completo), #17 (feed + contatori), più 404 e redirect.

**Artefatto**: la ricerca live trova un brano cercando un verso; il sito risponde in più lingue di interfaccia; esistono `rss.xml`, `sitemap.xml`, `robots.txt` e la 404 con ricerca.

**Criterio di chiusura**: cercando un frammento di verso si arriva al brano giusto dal dropdown; il feed è valido.

## FASE 5 — Messa online e osservabilità

**Task**: #8 (Cloudflare Tunnel + dominio), #14 (Umami), #15 (pagine legali).

Il container esiste già dalla FASE 1: qui si aggiunge ciò che manca per essere **pubblici**.

**Artefatto**: il sito risponde su `lyrica.filippomoscatelli.com` dietro Cloudflare Tunnel; Umami registra le visite; le pagine legali esistono e sono linkate dal footer.

**Criterio di chiusura**: l'URL pubblico risponde, e Umami mostra gli eventi.

## FASE 6 — Extra

**Task**: #16 (PWA offline), #12 (segnalazioni + contatto + bot Telegram), #11 (ads, quando ci saranno advertiser).

**Artefatto**: il sito è installabile e legge i brani offline; il form di segnalazione manda un messaggio al bot Telegram.

**Criterio di chiusura**: una segnalazione inviata dal sito arriva come messaggio Telegram.

## Regole del piano

- **Non si salta la FASE 0**: il codice segue i doc, non il contrario.
- Non si inizia una fase senza che la precedente abbia il suo artefatto verificabile.
- **Ogni fase si verifica su un container reale**, non solo in locale: `docker compose up -d` è la prova che la fase tiene.
- Le task dentro una fase possono procedere in parallelo solo dove non ci sono dipendenze: #2 e #3 sbloccano tutto il resto; #6, #7 e #4 sono paralleli dopo #3.
- Se durante lo sviluppo una decisione cambia, si aggiorna `DECISION.md` **nella stessa PR**, non dopo.

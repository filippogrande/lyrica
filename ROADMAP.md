# ROADMAP — piano di sviluppo di Lyrica

Ordine di sviluppo a fasi. Ogni fase si chiude con un artefatto verificabile: se l'artefatto non esiste, la fase non è chiusa.

Le task vivono su **Vikunja, progetto `Lyrica`** (`#1`–`#18`); qui c'è solo l'ordine e il criterio di chiusura.

## FASE 0 — Documentazione

**Task**: #13 (priorità massima).

**Artefatto**: questo set di MD completo nel repo. La documentazione viene **prima** del codice: è la sorgente da cui si sviluppa.

**Criterio di chiusura**: ogni decisione presa è scritta in `DECISION.md`; ogni area ha il suo doc; un agente AI può leggere `README.md` e `ARCHITECTURE.md` e sapere cosa fare senza chiedere.

## FASE 1 — Fondamenta

**Task**: #3 (scheletro Go + Templ), #7 (design system), #4 (routing).

**Artefatto**: il binario compila, serve una homepage vuota ma **navigabile** con header, footer, tema chiaro/scuro e toggle dimensione testo. In locale, non online.

**Criterio di chiusura**: `docker compose up` in locale risponde su una porta e si vedono header, footer, breadcrumb e switch tema funzionanti.

## FASE 2 — Contenuti

**Task**: #2 (schema), #18 (CLI), #9 (pipeline + validazione), #1 (primi contenuti reali).

**Artefatto**: la CLI genera i template; 2-3 band reali con un album e alcuni brani sono nel repo; la validazione in CI passa; il sito mostra una band e un album veri.

**Criterio di chiusura**: `lyrica build` produce le pagine di band/album dai contenuti reali, e la validazione **fallisce** su un contenuto rotto di prova (trad. incompleta, slug duplicato).

## FASE 3 — Pagine pubbliche

**Task**: #10 (vista a fronte) + le pagine Home / Bands / Band / Album.

**Artefatto**: un visitatore arriva in homepage, trova una band, apre un album e legge un brano con originale a sinistra e traduzione a destra, con il selettore lingua.

**Criterio di chiusura**: la navigazione completa Home → Bands → Band → Album → Brano funziona da telefono, con breadcrumb e cambio lingua.

## FASE 4 — Funzioni trasversali

**Task**: #5 (ricerca), #6 (i18n completo), #17 (feed + contatori), più 404 e redirect.

**Artefatto**: la ricerca live trova un brano cercando un verso; il sito risponde in più lingue di interfaccia; esistono `rss.xml`, `sitemap.xml`, `robots.txt` e la 404 con ricerca.

**Criterio di chiusura**: cercando un frammento di verso si arriva al brano giusto dal dropdown; il feed è valido.

## FASE 5 — Deploy e osservabilità

**Task**: #8 (deploy), #14 (Umami), #15 (pagine legali).

**Artefatto**: il sito è **online** su `lyrica.filippomoscatelli.com` dietro Cloudflare Tunnel; merge su `main` fa deploy da solo; Umami registra le visite; le pagine legali esistono e sono linkate dal footer.

**Criterio di chiusura**: un merge su `main` arriva in produzione senza intervento manuale, e Umami mostra gli eventi.

## FASE 6 — Extra

**Task**: #16 (PWA offline), #12 (segnalazioni + contatto + bot Telegram), #11 (ads, quando ci saranno advertiser).

**Artefatto**: il sito è installabile e legge i brani offline; il form di segnalazione manda un messaggio al bot Telegram.

**Criterio di chiusura**: una segnalazione inviata dal sito arriva come messaggio Telegram.

## Regole del piano

- **Non si salta la FASE 0.** Il codice segue i doc, non il contrario.
- Non si inizia una fase senza che la precedente abbia il suo artefatto verificabile.
- Le task dentro una fase possono procedere in parallelo solo dove non ci sono dipendenze: #2 e #3 sbloccano tutto il resto; #6, #7 e #4 sono paralleli dopo #3.
- Se durante lo sviluppo una decisione cambia, si aggiorna `DECISION.md` **nella stessa PR**, non dopo.

# Development Guidelines — Lyrica

Regole non negoziabili del progetto. Se una PR contraddice questo file, la PR è sbagliata (o va prima cambiato questo file, con una decisione registrata in `DECISION.md`).

## 1. Workflow git

- Ogni modifica passa da **branch + Pull Request**. **Mai** push diretto su `main`.
- Nomi branch: `feat/<slug>`, `fix/<slug>`, `docs/<slug>`, `ci/<slug>`.
- **Merge su `main` = deploy.** Dalla FASE 5 in avanti `main` è produzione: quello che mergi va online.
- PR piccole e su un solo tema. Nel body: cosa cambia, come verificarlo, cosa potrebbe rompersi.
- **La documentazione si aggiorna nella stessa PR** che la invalida. Mai "aggiorno i doc dopo".
- Ogni modifica corrisponde a **una task tracciata** (Vikunja, progetto `Lyrica`).

## 2. Paletti di prodotto

- **Nessun login, nessun account, nessun ruolo, nemmeno per l'autore.** Non esiste pannello admin. I test si fanno direttamente in produzione.
- Il sito è **statico**: niente database, niente sessione, nessuno stato per utente. Le uniche scritture sono i form (segnalazione, contatto), che **inoltrano a Telegram e non persistono nulla** sul sito.
- Niente commenti utente, niente upload, niente contenuto generato dagli utenti.
- **Niente scraping di testi**: i contenuti si scrivono a mano.
- Nessun contributo esterno: le traduzioni le scrive solo l'autore.

## 3. Stack e vincoli tecnici

- **Go** + **Templ** per l'HTML, servito **già renderizzato**.
- **HTMX** per le interazioni (ricerca live).
- **Bootstrap 5 solo CSS** (nessun bundle JS). JS custom ridotto al minimo indispensabile: tema, service worker, apertura menu.
- CSS custom tramite custom properties in `assets/css/`.
- Nessun framework JS, nessun passaggio npm come prerequisito per servire il sito.

## 4. Struttura del codice

- **Nessun file oltre ~1000 righe.**
- Si divide per **coesione semantica**, non a caso per numero di righe: un file per funzione o per gruppo di funzioni affini. Molti moduli piccoli che fanno una cosa sola battono pochi file grossi.
- Codice in **inglese** (nomi, commenti, log). Contenuti e documentazione in **italiano**.

## 5. Test e verifica

- **La CI valida i contenuti, non il codice**: front-matter, slug duplicati, link interni, traduzioni incomplete, cover mancanti.
- Non si scrivono test isolati o usa-e-getta che verificano solo la sintassi: mascherano i bug invece di rivelarli.
- La verifica si fa sulla **build reale** (`lyrica build` e poi il sito aperto), non su un frammento di codice.
- Una cosa è "fatta" quando è **verificata sull'artefatto** (la pagina generata, l'output del comando), non quando il codice è stato scritto.

## 6. Performance e accessibilità

- Budget: **~300 KB per pagina** (HTML + CSS; immagini escluse perché ottimizzate a parte).
- Immagini: **webp** + `loading="lazy"`.
- Ogni pagina usabile da tastiera, con contrasto sufficiente e controllo della dimensione del testo.
- Design **mobile-first**: si guarda prima da telefono.

## 7. Sicurezza e segreti

- Segreti (token del bot Telegram) **solo** via variabili d'ambiente o `.env` non committato.
- Header di sicurezza e CSP sempre presenti (vedi `SECURITY.md`).
- Nessun dato personale raccolto senza necessità e senza informativa (vedi `LEGAL_DESIGN.md`).

## 8. Dipendenze

- Ogni manifest presente nel repo (Go modules, `Dockerfile`, `docker-compose.yml`) deve essere coperto da `.github/dependabot.yml`, aggiornato **nella stessa PR** che introduce il manifest.

## 9. Documentazione

- Gli MD di questo repo sono scritti per essere letti **da un agente AI** che deve capire e modificare il sito: ogni file dichiara scope, decisioni e dove vive il codice.
- Un doc che non descrive più la realtà va **corretto o cancellato** nella stessa PR. Mai lasciare stub o note "in realtà ora è diverso".
- Nessun doc contiene segreti, credenziali o dati personali.

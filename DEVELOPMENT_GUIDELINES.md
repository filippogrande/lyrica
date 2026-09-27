# Development Guidelines — Lyrica

Regole non negoziabili del progetto. Se una PR contraddice questo file, la PR è sbagliata (o va prima cambiato questo file, con una decisione registrata in `DECISION.md`).

## 1. Workflow git

- Ogni modifica passa da **branch + Pull Request**. **Mai** push diretto su `main`.
- Nomi branch: `feat/<slug>`, `fix/<slug>`, `docs/<slug>`, `ci/<slug>`, `content/<band-slug>`.
- **Merge su `main` = deploy.** Dalla FASE 5 in avanti `main` è produzione: quello che mergi va online.
- PR piccole e su **un solo tema**: cosa cambia, come verificarlo, cosa potrebbe rompersi.
- **Una PR può toccare più file** quando fanno parte della stessa modifica (es. modulo + rotta + template + doc che lo descrive). La regola "una PR = un file" usata in altri progetti **qui non si applica**: il criterio è il tema della PR, non il numero di file — purché resti piccola e leggibile.
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

## 4. Dimensione del codice — limiti VINCOLANTI

- **Massimo 500 righe per file.** Limite vincolante, non un'indicazione.
- **Massimo 50 righe per funzione.** Limite vincolante.
- Un file che sfora va **diviso per responsabilità** (coesione semantica): non si taglia a metà per fare numero, si separa per ciò che il codice fa. Due file da 500 righe spezzati a caso sono peggio di uno da 1000 ben coeso.
- Preferire molti moduli piccoli che fanno **una cosa sola** a pochi file grossi: i file grossi si rompono e sono difficili da revisionare.
- Codice in **inglese** (nomi, commenti, log). Contenuti e documentazione in **italiano**.

## 5. Nessun fallback silenzioso

- **Vietato** restituire valori fittizi, vuoti o di default per mascherare un errore: niente `return []`, `return {}`, `return nil` "tanto per non far crashare", niente dati mock al posto dei dati reali.
- Se un contenuto è invalido, un file manca o una chiamata al bot Telegram fallisce, il codice **restituisce l'errore e lo rende visibile** (build rossa in locale/CI, messaggio d'errore in pagina, log esplicito).
- Una pagina che finge che tutto vada bene è peggio di una pagina con un errore: il fallback silenzioso nasconde i bug invece di rivelarli.
- L'unico fallback ammesso è quello **dichiarato e progettato**: la lingua dell'interfaccia che non esiste ricade sull'italiano (`I18N_DESIGN.md`) — ed è scritto nei doc, non improvvisato nel codice.

## 6. Test e verifica

- **La CI valida i contenuti, non il codice**: front-matter, slug duplicati, link interni, traduzioni incomplete, cover mancanti.
- **Vietati i file di test temporanei o usa-e-getta.** Un mini-test scritto per l'occasione carica solo una parte del codice, quindi segnala al massimo errori di battitura e **maschera i bug veri** invece di rivelarli.
- La verifica si fa sulla **build reale** (`lyrica build` e poi il sito aperto), con i moduli veri caricati, non su un frammento isolato.
- Una cosa è "fatta" quando è **verificata sull'artefatto** (la pagina generata, l'output del comando), non quando il codice è stato scritto. Se una verifica non è stata eseguita, va detto esplicitamente.

## 7. Performance e accessibilità

- Budget: **~300 KB per pagina** (HTML + CSS; immagini escluse perché ottimizzate a parte).
- Immagini: **webp** + `loading="lazy"`.
- Ogni pagina usabile da tastiera, con contrasto sufficiente e controllo della dimensione del testo.
- Design **mobile-first**: si guarda prima da telefono.

## 8. Sicurezza e segreti

- Segreti (token del bot Telegram) **solo** via variabili d'ambiente o `.env` non committato.
- Header di sicurezza e CSP sempre presenti (vedi `SECURITY.md`).
- Nessun dato personale raccolto senza necessità e senza informativa (vedi `LEGAL_DESIGN.md`).

## 9. Dipendenze

- Ogni manifest presente nel repo (Go modules, `Dockerfile`, `docker-compose.yml`) deve essere coperto da `.github/dependabot.yml`, aggiornato **nella stessa PR** che introduce il manifest.

## 10. Documentazione

- Gli MD di questo repo sono scritti per essere letti **da un agente AI** che deve capire e modificare il sito: ogni file dichiara scope, decisioni e dove vive il codice.
- Un doc che non descrive più la realtà va **corretto o cancellato** nella stessa PR. Mai lasciare stub o note "in realtà ora è diverso".
- Nessun doc contiene segreti, credenziali o dati personali.

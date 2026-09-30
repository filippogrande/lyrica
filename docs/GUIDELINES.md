# GUIDELINES — regole di Lyrica

Regole non negoziabili. Se una PR contraddice questo file, la PR è sbagliata — o va prima cambiato questo file, con una decisione registrata in `DECISIONS.md`.

## 1. Workflow git

- Ogni modifica passa da **branch + Pull Request**. **Mai** push diretto su `main`.
- Nomi branch: `feat/<slug>`, `fix/<slug>`, `docs/<slug>`, `ci/<slug>`, `content/<band-slug>`.
- **Merge su `main` = deploy**: quello che mergi finisce nell'immagine pubblicata. Il `docker compose up -d` sul home-lab resta un'azione manuale.
- PR piccole e **su un solo tema**, anche se toccano più file (D76). Nel body: cosa cambia, come verificarlo, cosa potrebbe rompersi.
- **La documentazione si aggiorna nella stessa PR** che la invalida. Mai "aggiorno i doc dopo".
- Ogni modifica corrisponde a **una task tracciata** (Vikunja, progetto `Lyrica`).

## 2. Paletti di prodotto

- **Nessun login, nessun account, nessun ruolo, nemmeno per l'autore.** Non esiste pannello admin. I test si fanno direttamente in produzione.
- Il sito è **statico**: niente database, niente sessione, nessuno stato per utente. Le uniche scritture sono i form (segnalazione, contatto), che **inoltrano a Telegram e non persistono nulla**.
- Niente commenti utente, niente upload, niente contenuto generato dagli utenti.
- **Niente scraping di testi**: i contenuti si scrivono a mano.
- Nessun contributo esterno: le traduzioni le scrive solo l'autore.

## 3. Stack e vincoli tecnici

- **Go** + **Templ** per l'HTML, servito **già renderizzato**.
- **HTMX** per le interazioni (ricerca live).
- **Bootstrap 5 solo CSS** (nessun bundle JS), servito **dal repo** e non da CDN.
- JS custom ridotto al minimo indispensabile: tema, dimensione testo, service worker, apertura menu, selettore lingua. Nessun framework JS, nessun passaggio npm.

## 4. Struttura del codice (vincolante)

- **Massimo 500 righe per file.** Nessuna eccezione.
- **Massimo 50 righe per funzione.** Se una funzione non ci sta, si scompone.
- Si divide per **coesione semantica**, mai a metà per fare numero: un file per funzione o per gruppo di funzioni affini. Molti moduli piccoli che fanno una cosa sola battono pochi file grossi.
- **Vietati i fallback silenziosi**: niente `return []`, `nil`, valore zero o dato finto per mascherare un errore. Un errore si **propaga e si vede**. L'unica eccezione ammessa è un fallback **progettato e dichiarato** (esempio: lingua dell'interfaccia non disponibile → italiano).
- Codice in **inglese** (nomi, commenti, log). Contenuti e documentazione in **italiano**.

## 5. Test e verifica

- **La CI valida i contenuti, compila il progetto ed esegue i test unitari del codice** (`go test ./...`, D98).
- **Vietati i file di test temporanei o usa-e-getta**: verificavano solo la sintassi e mascheravano i bug veri invece di rivelarli. I test che entrano nel repo provano una **regola** (il validatore, le lingue in pagina) e restano veri anche fra sei mesi.
- La verifica si fa sulla **build reale** (il binario compilato dalla CI, il sito generato e aperto), non su un frammento di codice: i test la affiancano, non la sostituiscono.
- Una cosa è "fatta" quando è **verificata sull'artefatto** (la pagina generata, l'output del comando, il container che risponde), non quando il codice è stato scritto.

## 6. Performance e accessibilità

- Budget: **~300 KB per pagina** (HTML + CSS; immagini a parte perché ottimizzate).
- Immagini: **webp** + `loading="lazy"` + dimensioni dichiarate.
- Ogni pagina usabile da tastiera, con contrasto sufficiente e controllo della dimensione del testo.
- Design **mobile-first**: si guarda prima da telefono.

## 7. Sicurezza e segreti

- Segreti (token del bot Telegram) **solo** via variabili d'ambiente o `.env` non committato.
- Header di sicurezza e CSP sempre presenti (vedi `ARCHITECTURE.md`).
- Nessun dato personale raccolto senza necessità e senza informativa (vedi `FEATURES.md`).
- **Vietato** committare un `.env` reale, anche "temporaneamente".

## 8. Dipendenze

- Ogni manifest nel repo (Go modules, `Dockerfile`, `docker-compose.yml`) deve essere coperto da `.github/dependabot.yml`, aggiornato **nella stessa PR** che introduce il manifest.

## 9. Documentazione

I documenti stanno in **`docs/`** e sono **8**, con questo scope:

| File | Scope |
|---|---|
| `GUIDELINES.md` | questo file: regole di lavoro |
| `DECISIONS.md` | registro delle decisioni (D01→), con il perché |
| `ROADMAP.md` | fasi, artefatti, criteri di chiusura |
| `SPEC.md` | cosa fa il sito: pagine, funzionalità, casi limite |
| `ARCHITECTURE.md` | come è fatto: componenti, build, Docker, deploy, sicurezza |
| `CONTENT.md` | schema dei contenuti e come si aggiungono |
| `FRONTEND.md` | interfaccia: design system, vista brano, lingue, ricerca, PWA |
| `FEATURES.md` | ads, analytics, legali/privacy, form, feed e contatori |

Regole:

- Sono scritti per essere letti **da un agente AI** che deve capire e modificare il sito: ogni file dichiara scope, decisioni e dove vive il codice.
- **Un doc che non descrive più la realtà va corretto o cancellato nella stessa PR.** Mai stub, mai note "in realtà ora è diverso".
- **Chiusura di fase**: quando una fase è chiusa, il progetto in essa contenuto o si **assorbe** in una sezione "come funziona" del doc di competenza, o si **elimina**. La verità è il codice; una doc che mente è peggio di nessuna doc. Il taglio si fa nella PR che chiude la fase.
- Nessun doc contiene segreti, credenziali o dati personali.
- Le decisioni di prodotto restano in `DECISIONS.md` anche quando il dettaglio implementativo sparisce: sono la memoria del perché.

# PWA_DESIGN — installabilità e lettura offline

## Obiettivo

Un testo tradotto serve spesso **senza rete** (in metro, in aereo, in un posto senza campo). La PWA esiste per questo: il brano che stai leggendo resta leggibile anche se il telefono perde la connessione.

## Manifest

- `manifest.webmanifest`: `name` Lyrica, `short_name` Lyrica, `start_url` la home nella lingua negoziata, `display: standalone`, `theme_color` coerente con il tema scuro, `background_color` coerente con il tema.
- **Icone**: al lancio **placeholder** (D23). Sostituzione con icone da pacchetto quando si affronta la grafica. Il manifest dichiara le dimensioni standard (192, 512, maskable).
- Nessun prompt di installazione forzato: si usa quello del browser. Nessun popup "installa la nostra app".

## Service worker

- Scritto a mano, poche decine di righe, **nessun tool di build** (coerente con `DEVELOPMENT_GUIDELINES.md`).
- **Cosa mette in cache**: CSS, font, JS minimo (shell dell'app) e le **pagine dei brani visitate**, così la lettura offline funziona sulle canzoni già aperte.
- Strategia:
  - shell (CSS/font/JS): **cache-first**, aggiornata con la versione dell'app;
  - pagine dei brani: **stale-while-revalidate** — si legge subito la copia locale e in background si aggiorna;
  - cover: cache-first, dimensione controllata.
- **Cosa NON mette in cache**: le risposte di ricerca (endpoint dinamico) e i form: meglio un errore chiaro che un comportamento ambiguo.

## Cosa funziona offline

| Funzione | Offline |
|---|---|
| Pagine dei brani già visitate (tutte le lingue presenti, cambio lingua compreso) | sì |
| Navigazione tra pagine già visitate | sì |
| Cambio tema e dimensione testo | sì (è locale) |
| Ricerca | no (avviso chiaro) |
| Segnalazioni e contatti | no (il form lo dice prima dell'invio) |
| Nuove traduzioni non ancora aperte | no |

## Aggiornamento

- Quando esce una nuova versione, il service worker aggiorna la shell **al reload successivo**: nessun banner "nuova versione disponibile" che copre il testo.
- La versione della cache è legata al commit di build: nessun contenuto vecchio servito a lungo.

## Vincoli

- Il sito resta **perfettamente usabile senza installare nulla**: la PWA è un extra, non un requisito.
- Nessuna notifica push al lancio (non c'è motivo di interrompere nessuno).
- La cache è limitata dalle regole sopra: non si scarica il sito intero sul telefono dell'utente.

# SUGGESTIONS_DESIGN — segnalazioni brani e contatti

## Perché esiste

Chi legge il sito spesso pensa a un brano che manca. Deve poterlo dire **senza account** (D08) e senza che l'autore debba gestire un pannello o una casella email (D66).

## Flusso

```
form sul sito  →  POST /api/segnala  →  validazione + rate-limit  →  messaggio al bot Telegram  →  grazie (nessuna attesa, nessun account)
```

**Nessun dato resta sul server**: il messaggio Telegram **è** il record (D16, vedi `LEGAL_DESIGN.md`).

## Form segnalazione brano (`/it/segnala`)

| Campo | Obbligatorio | Note |
|---|---|---|
| Nome band | sì | testo libero |
| Album | no | se lo sa |
| Titolo brano | sì | testo libero |
| Link di riferimento | no | utile per trovare il testo |
| Perché lo proponi | no | una riga |
| Email | no | **solo** se vuole essere ricontattato. Al lancio non parte nessuna email (D67) |

## Form contatti (`/it/contatti`)

Messaggio libero + email facoltativa + oggetto. Usato anche per le **richieste di rimozione** (vedi `LEGAL_DESIGN.md`).

## Notifica all'autore

- **Bot Telegram dedicato** (D66), separato da altri bot: le segnalazioni non si mescolano con le notifiche di casa o di lavoro.
- Ogni invio produce **un messaggio**: band, album, titolo, link, testo libero, email se presente, data/ora italiane.
- Se l'invio a Telegram **fallisce**: la pagina risponde con un errore onesto ("non è stato possibile inviare, riprova") e **non finge** che sia andata a buon fine. Non c'è coda né retry automatico al lancio.

## Anti-abuso

| Misura | Dettaglio |
|---|---|
| Honeypot | campo nascosto che un umano non compila |
| Tempo minimo | invio rifiutato se compilato in meno di pochi secondi |
| Rate-limit | per IP, in memoria; oltre soglia → 429 |
| Dimensione | limiti di lunghezza su ogni campo |
| Nessun allegato | niente upload: superficie in meno |

## GDPR

- Informativa **inline**, breve, prima del pulsante di invio: chi riceve i dati, perché, per quanto.
- Dichiarato che il messaggio passa da **Telegram** (servizio terzo).
- Nessun cookie, nessun tracciamento del form oltre l'evento aggregato `suggestion_sent` (vedi `ANALYTICS_DESIGN.md`).

## Cosa NON fa

- Nessuna coda di moderazione, nessun pannello, nessuna risposta automatica.
- Nessuna email di conferma all'utente (D67, in backlog come `B04`).
- Nessuna pubblicazione automatica del suggerimento: la valutazione è umana.

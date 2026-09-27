# LEGAL_DESIGN — disclaimer, licenze, privacy

Il sito pubblica **traduzioni di testi di cui non detiene i diritti**: è la prima cosa da gestire con chiarezza.

## Copyright dei testi originali

- I testi originali appartengono ai rispettivi autori, compositori ed editori. Lyrica **non li rivendica** e non ne trae profitto (D07).
- Nel footer, in ogni pagina, un disclaimer esplicito con il rimando alla pagina legale.
- Le traduzioni sono presentate come **opere amatoriali** dell'autore del sito, non come edizioni ufficiali.
- **Procedura di rimozione**: chi detiene i diritti e ritiene un contenuto lesivo può chiedere la rimozione tramite il form contatti; il contenuto viene rimosso e l'eventuale URL rediretto o dismesso.

## Licenze

| Cosa | Licenza |
|---|---|
| Codice del sito | **MIT** (D05) |
| Traduzioni | **CC BY-NC 4.0** (D06) |

La licenza delle traduzioni è dichiarata **in ogni pagina brano** (in fondo, insieme ai metadati) e nel footer.

## Privacy

Il sito **non ha account, non ha login e non salva dati personali** sul server (D16). Questo semplifica l'informativa ma non la elimina: i due form raccolgono dati e li trasmettono.

### Form segnalazione / contatto — cosa succede ai dati

1. L'utente compila il form (dati: quello che scrive, più email **solo se la fornisce**).
2. Il server valida e **inoltra il messaggio al bot Telegram** dell'autore.
3. **Il server non conserva nulla**: nessun database, nessun file di log con i contenuti del form.
4. Il messaggio resta nella chat Telegram dell'autore (servizio terzo, Telegram) — questo va **detto esplicitamente** nell'informativa.

### Informativa minima (contenuto richiesto)

- Titolare: Filippo (contatto nel footer).
- Finalità: rispondere alla segnalazione / al messaggio; valutare l'aggiunta di un brano.
- Base giuridica: consenso implicito nell'invio + interesse legittimo a rispondere.
- Conservazione: il messaggio sulla chat Telegram finché è utile; **sul sito, nessuna conservazione**.
- Destinatari: Telegram (come piattaforma di recapito).
- Diritti: cancellazione su richiesta, via lo stesso form.

### Cookie

- **Nessun cookie.** Le preferenze (tema, dimensione testo) sono in `localStorage` e restano sul dispositivo dell'utente.
- Umami è cookieless (vedi `ANALYTICS_DESIGN.md`) → **nessun banner cookie**.

## Pagine legali da pubblicare

Un'unica pagina `/it/legali` con sezioni chiare e ancorabili:

1. Copyright dei testi e natura amatoriale delle traduzioni
2. Licenze (MIT per il codice, CC BY-NC 4.0 per le traduzioni)
3. Privacy (form, Telegram, nessuna conservazione sul sito)
4. Cookie (nessuno) e statistiche (Umami)
5. Richieste di rimozione

Linkata dal **footer** di ogni pagina e dai due form.

## Se un giorno si attivassero le ads

Vanno rispettate le condizioni di `ADS_DESIGN.md`: pubblicità senza tracking e senza profilazione. Una pubblicità con tracciamento romperebbe questa pagina e richiederebbe banner e consenso: è una decisione, non un dettaglio implementativo.

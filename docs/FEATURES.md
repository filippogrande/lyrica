# FEATURES — ads, analytics, legali, form, feed

Le funzioni che stanno **attorno** al testo: pubblicità, statistiche, obblighi legali, form e distribuzione.

---

# Pubblicità

## Stato al lancio

**Spente.** Gli slot esistono nel codice ma `enabled: false` in `ads.yaml`: nessuno spazio pubblicitario viene renderizzato finché non è una decisione esplicita.

Quando sono spente, **il DOM non contiene nemmeno i contenitori**: il testo prende tutta la larghezza e il layout non ha buchi.

## Principi (se e quando si attivano)

1. **Non invasive**: piccole, statiche, mai popup, mai interstitial, mai sticky.
2. **Nessun tracking**: banner serviti dallo stesso dominio, nessun script di terze parti, nessun cookie di profilazione.
3. **Coerenti con la privacy**: se una pubblicità introducesse tracciamento servirebbero consenso e banner — quindi non si introduce.
4. **Contano nel budget**: lo spazio occupato fa parte dei ~300 KB per pagina.

## Posizionamento (deciso)

### Pagina brano

```
ads            <- banner largo e basso, sopra il titolo
 titolo
 canzone       <- con due colonne laterali di ads ai lati (solo schermi larghi)
 chiedi altre canzoni
ads            <- banner largo e basso, dopo la sezione "chiedi altre canzoni"
```

| Slot | Dove | Formato | Visibilità |
|---|---|---|---|
| `track_top` | sopra il titolo del brano | banner largo e basso | sempre |
| `side_left` / `side_right` | ai **lati** delle due colonne di testo | verticale stretta | **solo da `xl` in su** |
| `track_bottom` | dopo "chiedi altre canzoni" | banner largo e basso | sempre |

Regole delle colonne laterali:

- Il **testo non si stringe** per far posto agli ads: sotto `xl` le colonne laterali **non esistono** e il testo usa tutto lo spazio.
- **Non sono sticky**: scorrono con la pagina. Un banner che segue lo scroll è la definizione di invasivo.
- Non si mettono **tra** originale e traduzione: spezzerebbe la lettura a fronte, che è il motivo per cui il sito esiste.
- Niente ads **dentro** il blocco delle strofe.

### Altre pagine

| Slot | Dove |
|---|---|
| `top` | sotto l'header nelle pagine di elenco |
| `footer` | sopra il footer |

Gli slot **non compaiono** su 404, pagine legali e form.

## Configurazione

`ads.yaml` contiene l'interruttore globale `enabled: false` e, per ogni slot, se è attivo, il frammento HTML o l'immagine locale, il link di destinazione, l'etichetta "sponsorizzato".

Le immagini pubblicitarie vivono nel repo, come le cover: il sito non chiama domini terzi.

## Requisiti prima di attivare

- [ ] Ricontrollare la privacy (informativa e cookie).
- [ ] Etichettatura "sponsorizzato" visibile e leggibile.
- [ ] Verifica del peso aggiunto rispetto al budget.
- [ ] Controllo da telefono **e** da schermo largo.
- [ ] Verifica che nessuno slot copra o interrompa il testo.

## Cosa non si farà mai

- Nessun annuncio **sticky** o che segue lo scroll.
- Nessun annuncio tra originale e traduzione, o dentro le strofe.
- Nessun annuncio con audio o video in autoplay.
- Nessuna raccolta di dati per profilazione.

---

# Analytics — Umami

## Perché Umami

- **Cookieless**: nessun cookie, nessun identificatore persistente → **nessun cookie banner**.
- **Riuso dell'istanza esistente**: nessun container nuovo da mantenere.
- Misura ciò che serve (cosa viene letto, cosa si cerca) **senza** profilare nessuno.

## Setup

1. Nell'istanza Umami già attiva si crea un **nuovo website** dedicato a Lyrica.
2. Si ottiene lo **script tag** e lo si inserisce nel `<head>` di tutte le pagine, in **un solo punto** dei template.
3. Lo script è **asincrono** e non blocca il render.

## Custom event

| Evento | Quando | A cosa serve |
|---|---|---|
| `search` | una ricerca produce risultati | capire cosa cerca la gente e non trova |
| `search_empty` | una ricerca non produce risultati | materia prima per le prossime traduzioni |
| `lang_switch` | cambio della lingua del testo | quali traduzioni servono davvero |
| `ui_lang_switch` | cambio della lingua dell'interfaccia | quali lingue di interfaccia valgono la pena |
| `theme_toggle` | cambio tema | se il tema manuale viene usato |
| `text_size` | uso di A-/A+ | se l'accessibilità serve davvero |
| `pwa_install` | installazione della PWA | quanti la usano come app |
| `suggestion_sent` | segnalazione inviata | volume reale delle proposte |

Nessun evento contiene **dati personali** né il testo cercato: si traccia il fatto che una ricerca c'è stata, non chi l'ha fatta né cosa ha scritto.

## Cosa NON si traccia

- Nessun tracciamento cross-sito, nessun pixel di terze parti.
- Nessun salvataggio del testo delle ricerche nella dashboard.
- Nessun profilo utente, nessuna cronologia.

## Se un giorno servissero più dati

È una **nuova decisione** da registrare in `DECISIONS.md`, non un'estensione silenziosa. Il vincolo di partenza (niente cookie, niente banner) non si rompe in una PR di routine.

---

# Legali, licenze e privacy

Il sito pubblica **traduzioni di testi di cui non detiene i diritti**: è la prima cosa da gestire con chiarezza.

## Copyright dei testi originali

- I testi appartengono ai rispettivi autori, compositori ed editori. Lyrica **non li rivendica** e non ne trae profitto.
- In **ogni pagina**, nel footer, un disclaimer esplicito con rimando alla pagina legale.
- Le traduzioni sono presentate come **opere amatoriali** dell'autore del sito, non come edizioni ufficiali.
- **Procedura di rimozione**: chi detiene i diritti può chiedere la rimozione tramite il form contatti; il contenuto viene rimosso e l'eventuale URL rediretto o dismesso.

## Licenze

| Cosa | Licenza |
|---|---|
| Codice del sito | **MIT** |
| Traduzioni | **CC BY-NC 4.0** |

La licenza delle traduzioni è dichiarata in ogni pagina brano (in fondo, con i metadati) e nel footer.

## Privacy

Il sito **non ha account e non salva dati personali** sul server. Questo semplifica l'informativa ma non la elimina: i due form raccolgono dati e li trasmettono.

### Form: cosa succede ai dati

1. L'utente compila il form (quello che scrive, più email **solo se la fornisce**).
2. Il server valida e **inoltra il messaggio al bot Telegram** dell'autore.
3. **Il server non conserva nulla**: nessun database, nessun file di log con i contenuti del form.
4. Il messaggio resta nella chat Telegram dell'autore — servizio terzo, e va **detto esplicitamente** nell'informativa.

### Informativa minima (contenuto richiesto)

- Titolare: Filippo (contatto nel footer).
- Finalità: rispondere alla segnalazione / al messaggio; valutare l'aggiunta di un brano.
- Base giuridica: consenso implicito nell'invio + interesse legittimo a rispondere.
- Conservazione: il messaggio sulla chat Telegram finché è utile; **sul sito, nessuna conservazione**.
- Destinatari: Telegram (come piattaforma di recapito).
- Diritti: cancellazione su richiesta, via stesso form.

### Cookie

- **Nessun cookie.** Le preferenze (tema, dimensione testo) sono in `localStorage` e restano sul dispositivo.
- Umami è cookieless → **nessun banner cookie**.

## Pagina legale da pubblicare

Un'unica pagina `/it/legali`, con sezioni ancorabili:

1. Copyright dei testi e natura amatoriale delle traduzioni
2. Licenze (MIT per il codice, CC BY-NC 4.0 per le traduzioni)
3. Privacy (form, Telegram, nessuna conservazione sul sito)
4. Cookie (nessuno) e statistiche (Umami)
5. Richieste di rimozione

Linkata dal **footer** di ogni pagina e dai due form.

## Se un giorno si attivassero le ads

Vanno rispettate le condizioni della sezione Pubblicità: niente tracking, niente profilazione. Una pubblicità con tracciamento romperebbe questa pagina e richiederebbe banner e consenso: è una decisione, non un dettaglio implementativo.

---

# Form: segnalazioni e contatti

## Perché esistono

Chi legge il sito spesso pensa a un brano che manca. Deve poterlo dire **senza account** e senza che l'autore debba gestire un pannello o una casella email.

## Flusso

```
form sul sito  →  POST /api/segnala  →  validazione + rate-limit  →  messaggio al bot Telegram  →  grazie
```

**Nessun dato resta sul server**: il messaggio Telegram **è** il record.

## Form segnalazione brano (`/it/segnala`)

| Campo | Obbligatorio | Note |
|---|---|---|
| Nome band | sì | testo libero |
| Album | no | se lo sa |
| Titolo brano | sì | testo libero |
| Link di riferimento | no | utile per trovare il testo |
| Perché lo proponi | no | una riga |
| Email | no | **solo** se vuole essere ricontattato; al lancio non parte nessuna email |

## Form contatti (`/it/contatti`)

Messaggio libero + email facoltativa + oggetto. Usato anche per le **richieste di rimozione**.

## Notifica all'autore

- **Bot Telegram dedicato** (`@LyricaNotifyBot`), separato dagli altri bot: le segnalazioni non si mescolano con notifiche di casa o di lavoro.
- Ogni invio produce **un messaggio**: band, album, titolo, link, testo libero, email se presente, data/ora italiane.
- Se l'invio a Telegram **fallisce**: la pagina risponde con un errore onesto ("non è stato possibile inviare, riprova") e **non finge** che sia andata a buon fine. Nessun coda, nessun retry automatico al lancio.

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
- Dichiarato che il messaggio passa da **Telegram**.
- Nessun cookie, nessun tracciamento del form oltre all'evento aggregato `suggestion_sent`.

## Cosa NON fanno

- Nessuna coda di moderazione, nessun pannello, nessuna risposta automatica.
- Nessuna email di conferma all'utente.
- Nessuna pubblicazione automatica del suggerimento: la valutazione è umana.

---

# Feed RSS e contatori

Tutti e tre i file qui sotto (`rss.xml`, `sitemap.xml`, `robots.txt`) sono **generati dal build**, statici come il resto del sito: nessun servizio esterno, nessuna dipendenza a runtime. Il dominio assoluto è una **costante in `internal/build`** (`https://lyrica.filippomoscatelli.com`): i tre file contengono URL assoluti e il dominio vive in un posto solo, senza variabili d'ambiente.

## Feed RSS

- **Un feed per lingua dell'interfaccia**: `/it/rss.xml` e `/en/rss.xml`. Nessun feed a root e **nessun feed per band al lancio**.
- Contiene le **ultime traduzioni pubblicate**, ordinate per `added_date` decrescente, **limite 30 voci**.
- Solo brani con **almeno una traduzione pubblicata**: una voce compare quando la traduzione c'è, non quando il brano viene creato a metà.
- Ogni voce: titolo del brano, **link assoluto** alla pagina del brano in quella lingua, `guid` stabile (permalink assoluto), `pubDate` in RFC1123 UTC ricavato da `added_date`, e descrizione generata = la meta description (titolo, band, album) più un'anteprima breve delle prime righe della traduzione in italiano (o nella lingua del feed, se presente).
- RSS 2.0 con **struct tipizzate e `encoding/xml`**: i testi contengono virgolette, apostrofi, `&` e il minore, e solo l'encodificatore li scappa come si deve.
- Nel `<head>` di ogni pagina: `<link rel="alternate" type="application/rss+xml">` con `href` **assoluto** verso il feed della lingua corrente.
- Link al feed nel **footer**.

## sitemap.xml

- **Un unico file alla root**, non uno per lingua (D99): le versioni linguistiche della stessa pagina sono un'unica URL con degli alternati.
- Namespace `http://www.sitemaps.org/schemas/sitemap/0.9`, con il namespace `xhtml` dichiarato sull'elemento radice (senza, gli hreflang vengono ignorati).
- Elenca **tutte le pagine reali prodotte dal build**, per ogni lingua dell'interfaccia: home, elenco band, ogni band, ogni album pubblicato, ogni brano tradotto. Le voci di tracklist con `status: pending` non hanno pagina, quindi non compaiono: la sitemap promette solo URL che rispondono.
- Ogni URL: `loc` assoluto, `lastmod` in ISO 8601 (per un brano la sua `added_date`; per album e band la data più recente dei brani pubblicati sotto, omessa quando non ci sono), `changefreq` e `priority` coerenti (home 1.0, elenco band 0.8, band 0.7, album 0.6, brano 0.5).
- **hreflang reciproci**: ogni pagina si dichiara in ogni lingua dell'interfaccia più `x-default` verso l'italiano. Gli stessi legami sono nel `<head>` di ogni pagina: le due liste non possono divergere.

## robots.txt

- Generato dal build, **non statico nel repo**: le righe `Sitemap` contengono gli URL assoluti di `sitemap.xml` e dei `rss.xml` di ogni lingua, e una lista di linghe scritta a mano diventerebbe subito falsa (D69).
- Permesso a tutti (`User-agent: *`, `Allow: /`): il sito è pubblico e l'indicizzazione è il punto.

## Contatori pubblici

Mostrano il valore del sito a colpo d'occhio:

| Contatore | Definizione esatta |
|---|---|
| **Brani tradotti** | brani con almeno una traduzione pubblicata (`HasTranslations`) |
| **Band** | band con almeno un album che ha almeno un brano tradotto (`publishedAlbums`) |
| **Lingue** | le lingue di **traduzione di arrivo** (`Catalog.TranslationLangs`, D69 e D96) |

La semantica delle **lingue** è una decisione esplicita (D104): si mostra ciò che un lettore può trovare — oggi italiano e inglese, quindi **2** — e **non** tutte le lingue originali tradotte, che darebbero 3 per via del tedesco "completo" di un brano bilingue (D96). Il tedesco è una versione completa in pagina, non un'interfaccia in tedesco.

- Sono **calcolati a build time** dai contenuti: nessun numero scritto a mano che possa divergere dalla realtà.
- **Comprimono**: un brano con 4 lingue conta 1 brano, non 4. L'obiettivo è mostrare che il sito cresce, non gonfiare i numeri.
- Dove: **homepage** (sotto la presentazione) e **footer**, stessa riga e stessi numeri.
- Le etichette sono chiavi di locale (`stats.tracks`, `stats.bands`, `stats.langs`) in `locales/it.yaml` e `locales/en.yaml`.
- Nessun contatore di visite: le statistiche sono di Umami, non un numero da esibire.

## Perché questi tre insieme

Sono le tre facce della stessa idea: il sito si aggiorna nel tempo e chi lo segue non deve andarlo a controllare. Il feed serve a chi lo segue, i contatori a chi arriva per la prima volta, sitemap e robots a chi arriva da un motore di ricerca.
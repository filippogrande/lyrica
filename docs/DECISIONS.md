# DECISIONS — registro delle decisioni

Ogni decisione presa ha un numero. Se una decisione viene cambiata, **non si cancella**: si aggiunge una nuova voce che la sostituisce.

## A. Identità, legale, processo

| # | Decisione | Perché |
|---|---|---|
| D01 | Nome del sito: **Lyrica** | doveva suonare musicale, corto, senza richiamare una band specifica |
| D02 | Repo GitHub `filippogrande/lyrica` | convenzione degli altri progetti |
| D03 | URL di produzione `lyrica.filippomoscatelli.com` | sottodominio dedicato, non path su un sito esistente |
| D04 | **Nessun logo**, look pulito | un logo disegnato male peggiora più di quanto un logo aiuti |
| D05 | Licenza del codice: **MIT** | standard per i progetti personali |
| D06 | Licenza delle traduzioni: **CC BY-NC 4.0** | le traduzioni sono opera dell'autore, non commerciali |
| D07 | Disclaimer copyright nel footer + dichiarazione "traduzioni amatoriali, sito no profit" | mitiga il rischio con gli editori dei testi |
| D08 | **Zero login e zero account**, nemmeno per l'autore. Nessun pannello admin. | il sito non ha stato utente; l'autore testa in produzione |
| D09 | Workflow git: **sempre branch + PR**, anche per l'autore | storico leggibile e ogni cambiamento annullabile |
| D10 | **Merge su `main` = deploy** (immagine su Docker Hub; il `up -d` sul home-lab resta manuale) | una sola strada per mettere online, senza automatismi che sfuggono di mano |

## B. Stack e infrastruttura

| # | Decisione | Perché |
|---|---|---|
| D11 | Backend in **Go** (scartato Rust) | build veloce, semplicità, un binario solo |
| D12 | HTML con **Templ**, servito già renderizzato | type-safe, niente template runtime |
| D13 | **HTMX** per le interazioni, **Bootstrap 5 CSS-only** per il layout | niente framework JS, poche dipendenze |
| D14 | Contenuti in **Markdown + front-matter YAML nel repo**, **build-time static site** | niente database, contenuti versionati con il codice |
| D15 | Deploy su **Docker nel home-lab dietro Cloudflare Tunnel** (non k3s) | k3s in dismissione; tunnel già in uso |
| D16 | **Nessun database, nessuno stato server-side** | elimina backup, migrazioni, GDPR sui dati a riposo |
| D17 | Test: **solo validazione dei contenuti in CI** (+ compilazione, vedi D73) | i test di codice non sono richiesti ora |
| D18 | Budget di performance **~300 KB/pagina**; immagini **webp + lazy** | sito leggero, buono anche da mobile |
| D19 | **Font web multialfabeto** (latino, cirillico, greco, CJK di sistema) | le band possono avere testi in alfabeti non latini |
| D20 | **CLI `lyrica new band\|album\|brano`** che genera i template di contenuto | evita errori di struttura nei front-matter |
| D21 | Cover: **600x600 webp quadrata**, nome file = slug dell'album | uniformità di resa senza dover pensare al formato |
| D22 | Umami: **riuso dell'istanza self-hosted esistente** | nessun container in più da mantenere |
| D23 | Favicon e icone PWA: **placeholder** ora, icone da pacchetti esistenti poi | non blocca lo sviluppo per una cosa grafica |

## C. Contenuti e URL

| # | Decisione | Perché |
|---|---|---|
| D24 | Band: **qualsiasi** (TBS è solo il primo caso d'uso) | il sito non è legato a una band |
| D25 | **Tag di genere liberi** per band, con pagine di esplorazione per tag | l'autore scrive il tag che serve, senza lista chiusa da mantenere |
| D26 | Le traduzioni **non sono obbligatorie**: si mostrano le lingue realmente disponibili | si pubblica quello che c'è, quando c'è |
| D27 | **Qualsiasi lingua**, non solo IT/EN/DE | l'architettura non deve porre limiti |
| D28 | Un brano può avere **più lingue originali** (band bilingui) | esistono davvero, non ha senso escluderle |
| D29 | Una traduzione si pubblica **solo se completa** | un testo a metà non serve a nessuno |
| D30 | L'**originale è sempre attivo**; il selettore mostra solo le lingue presenti per quel brano | zero vicoli ciechi |
| D31 | Brani **strumentali**: in tracklist con nota, **senza pagina** | non c'è testo da mostrare |
| D32 | Versione **live/acustica**: blocco con nota dentro lo stesso brano, non pagina separata | è lo stesso testo |
| D33 | **Nessuna nota/annotazione per verso. Mai.** Testo pulito | l'autore non le vuole |
| D34 | Testi **scritti a mano**, niente scraping | legale e qualità |
| D35 | Pagina band: si mostrano **solo gli album con almeno una traduzione pubblicata** | niente pagine vuote |
| D36 | Album senza cover: **nessuna immagine**, layout senza copertina | meglio del placeholder brutto |
| D37 | Elenchi "recenti" ordinati per **data di aggiunta** | riflette l'attività reale del sito |
| D38 | **Slug disambiguati** quando due band hanno lo stesso nome (`nirvana-us`) | gli URL restano univoci |
| D39 | **Redirect 301** gestiti da un file nel repo | i vecchi URL non muoiono mai |
| D40 | URL con **path prefix lingua**: `/it/band/{slug}/album/{slug}/brano/{slug}` | SEO e chiarezza |

## D. Pagine e interfaccia

| # | Decisione | Perché |
|---|---|---|
| D41 | Gerarchia: **Home → Bands → Band → Album → Brano** | percorso naturale per chi cerca un testo |
| D42 | **Homepage**: 3-5 "in evidenza" in cima, poi contatori e recenti | mostra subito il meglio, non solo l'ultimo arrivato |
| D43 | "In evidenza" **scelto a mano** | scelta editoriale, non algoritmica |
| D44 | Pagina **Bands**: elenco alfabetico + filtro per tag | trovare una band senza cercarla per nome |
| D45 | Etichetta della strofa = **solo nome del cantante**; niente Verse/Chorus/Bridge. Nei brani multi-voce il nome sta su **ogni strofa** | è l'informazione utile; l'autore non vuole etichette strutturali |
| D46 | **Breadcrumb** sulle pagine interne | orientamento |
| D47 | Mobile: **header sticky con nav compressa + offcanvas** | il testo deve restare protagonista |
| D48 | **Header**: Home / Bands / Ricerca | minimo indispensabile |
| D49 | **Footer**: disclaimer + licenza CC + contatti + RSS + contatori | tutto ciò che è obbligatorio o utile, fuori dai piedi |
| D50 | Tema **auto** (`prefers-color-scheme`) + **toggle manuale** in `localStorage` | rispetto delle preferenze di sistema, con override |
| D51 | Controllo **dimensione testo A-/A+** | accessibilità reale, non solo dichiarata |
| D52 | Lingua dell'interfaccia **negoziata dal browser** (`Accept-Language`), fallback italiano | l'utente non deve sceglierla se non serve |
| D53 | **404 personalizzata con ricerca** | chi sbaglia URL deve poter cercare subito |

## E. Ricerca, SEO, distribuzione

| # | Decisione | Perché |
|---|---|---|
| D54 | Ricerca **full-text su versi + titoli + band** | chi cerca una frase ricorda il verso, non il titolo |
| D55 | Risultati come **dropdown live as-you-type** (HTMX), **nessuna pagina di ricerca** | meno click, meno pagine da progettare |
| D56 | Risultato = **titolo brano + band**, senza snippet del verso | risultato compatto e leggibile |
| D57 | Ordine: **rilevanza** (match nel titolo → band → verso) | il titolo è il segnale più forte |
| D58 | **Meta description automatica** per brano | costante, senza lavoro manuale per pagina |
| D59 | Dati strutturati **schema.org** (`MusicGroup`, `MusicRecording`) | SEO |
| D60 | **sitemap.xml + hreflang + robots.txt**, sito indicizzabile | il sito deve essere trovato |
| D61 | **Feed RSS** delle nuove traduzioni | chi segue il sito non deve controllarlo |
| D62 | **Contatori pubblici** (brani / band / lingue) | misura il valore del sito a colpo d'occhio |
| D63 | **PWA installabile + lettura offline** dei brani | un testo tradotto serve spesso senza rete |

## F. Pubblicità, analytics, form

| # | Decisione | Perché |
|---|---|---|
| D64 | Ads **statiche, senza tracking, GDPR-compliant**, slot **nascosti al lancio** (`enabled: false`) | non si parte con pubblicità; quando ci sarà, non deve diventare un problema legale |
| D65 | Analytics **Umami** (cookieless) su istanza esistente, con custom event | metriche utili senza cookie banner |
| D66 | Segnalazioni brani e form di contatto arrivano via **bot Telegram dedicato (`@LyricaNotifyBot`)** | nessun pannello admin, nessuna email da gestire |
| D67 | **Nessuna email automatica al lancio** (l'avviso "traduzione pronta" è in backlog) | doppio opt-in e deliverability non valgono il costo ora |
| D68 | Segnalazioni: **rate-limit + honeypot + informativa GDPR**, niente account | il form resta aperto a tutti senza diventare uno spam gateway |

## G. Revisioni successive

| # | Decisione | Perché |
|---|---|---|
| D69 | Le **lingue dell'interfaccia derivano dalle lingue di traduzione** presenti nei contenuti (non dalle lingue degli originali) + italiano di default | se un brano ha traduzioni in IT ed EN, il pubblico parla IT o EN e quasi mai tedesco; quando arriveranno traduzioni in francese, servirà anche l'interfaccia francese |
| D70 | **Ads anche nella pagina brano**, con collocazione precisa: banner largo e basso **sopra il titolo**, **due colonne laterali** ai lati del testo (solo da `xl` in su, mai sticky, mai tra originale e traduzione), banner largo e basso **dopo la sezione "chiedi altre canzoni"** | collocazione voluta dall'autore: gli ads ai lati possono convivere con la lettura a fronte, purché il testo non si stringa sugli schermi piccoli |
| D71 | **Sezione "chiedi altre canzoni"** in fondo a ogni pagina brano, con link al form di segnalazione | chi ha appena letto una traduzione è la persona con più probabilità di proporne un'altra |
| D72 | **Docker da subito** (non in FASE 5): immagine multi-stage su **Docker Hub** `filippogrande/lyrica` con tag `latest` e `<sha>`, build e push automatici da GitHub Actions su push in `main`, compose nel repo, **deploy manuale** con `docker compose pull && docker compose up -d` sul home-lab | poter tirare su il sito con docker compose appena c'è un'immagine, come per gli altri servizi; il deploy resta una decisione esplicita e il runner non ha accesso al home-lab |
| D73 | La CI **compila il progetto** a ogni PR e push, oltre alla validazione dei contenuti che arriverà in FASE 2 | l'ambiente dell'agente non ha un compilatore Go: la compilazione è l'unica verifica disponibile dell'artefatto |
| D74 | **Struttura della documentazione**: 8 file dentro `docs/` (GUIDELINES, DECISIONS, ROADMAP, SPEC, ARCHITECTURE, CONTENT, FRONTEND, FEATURES). I doc di design di una fase, a fase chiusa, si **assorbono** in una sezione "come funziona" o si **eliminano** | 20 file in root rendevano il repo illeggibile; i doc-progetto che sopravvivono al codice iniziano a descrivere qualcosa che non esiste più |
| D75 | I file di deploy (`docker-compose.yml`, `.env.example`) stanno in **`deploy/`**, non nella root | la root del repo deve mostrare il progetto, non l'infrastruttura |
| D76 | Regole di codice **vincolanti**: max **500 righe per file**, max **50 righe per funzione**, **vietati i fallback silenziosi** (eccezione: fallback progettato e dichiarato, es. lingua UI → italiano), **vietati i file di test temporanei/usa-e-getta**, **PR multi-file tematiche ammesse** | sono le regole degli altri progetti dell'autore: un file grosso nasconde il disordine, un fallback silenzioso nasconde un bug, un test usa-e-getta maschera i bug veri |
| D77 | La home è una **sequenza di corsie orizzontali** ("In evidenza", "Ultimi brani tradotti", "Ultime band aggiunte"), **massimo 12 card per corsia**, realizzate come elenco in overflow con `scroll-snap` — **nessun carosello, nessuno scorrimento automatico, nessun JS**; le corsie esistono **solo in home** | la home deve mostrare le novità, non essere un archivio da scaricare tutto; un elenco scorrevole resta nel flusso del documento, si legge da tastiera e da screen reader, e non costa una libreria di slider da mantenere |
| D78 | **Immagine della band**: campo `image` in `band.md`, file in `covers/` (600x600 webp) come le cover; quando l'immagine non c'è (band senza foto o album senza cover) la **card mostra il segnaposto con l'iniziale** del nome, su fondo del tema. Sostituisce D36 **solo per le card delle corsie**: la pagina album continua a non mostrare immagine se la cover manca | una home di sole scritte è povera, ma un'immagine finta o un placeholder grafico sarebbe peggio della lettera: l'iniziale è dichiarata, non costa una richiesta in rete e non sposta il layout |
| D79 | L'ordine di "Ultime band aggiunte" è il **contributo più recente** della band (data di aggiunta del suo brano pubblicato più recente), a parità di data **ordine alfabetico**; una band senza brani pubblicati **non compare** | l'ordine deve riflettere il lavoro reale (D37), non la data in cui è nata la cartella: una band "nuova" perché ha appena pubblicato una traduzione è la novità che interessa a chi apre la home |
| D80 | La **foto della band** compare anche nella **pagina band** (immagine singola sotto il titolo, larghezza massima 20rem come la copertina dell'album, `alt` tradotto); l'**elenco band resta testuale**, senza immagine per riga, e se `image` manca la pagina band non mostra nessuna immagine (il segnaposto con l'iniziale resta solo nelle card delle corsie) | una foto in cima alla pagina band è contesto utile quando c'è; ripeterla su ogni riga di un elenco la renderebbe pesante e illeggibile, e un segnaposto grande su una pagina intera sarebbe peggio di nessuna immagine |
| D81 | Le regole di stile delle pagine interne (breadcrumb, elenchi, anagrafica band, album, **vista a fronte a due colonne da `lg`**) restano in `assets/css/lyrica.css` e sono parte del sito: le classi nel CSS e quelle usate nei template devono coincidere, perché la CI compila i template ma **non vede il CSS** | con le corsie della home sono sparite le regole delle pagine interne e il CSS non è coperto da nessun controllo: la pagina brano è tornata a una colonna senza che niente la segnalasse. Chi tocca gli stili verifica le pagine interne a mano prima di chiudere la PR (stessa logica del contrasto, §Accessibilità) |

## H. Backlog esplicito (deciso di NON fare ora)

| # | Rimandato |
|---|---|
| B01 | Interviste / articoli sulle band |
| B02 | Export PDF / EPUB delle traduzioni |
| B03 | Commenti degli utenti |
| B04 | Email "la traduzione che hai chiesto è pronta" |
| B05 | Area admin (esclusa dall'architettura, non solo rimandata) |

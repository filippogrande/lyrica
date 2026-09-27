# SEARCH_DESIGN — ricerca full-text

## In una riga

Chi cerca un **verso**, un **titolo** o una **band** trova il brano da un **dropdown live nell'header**, senza mai caricare una pagina di ricerca.

## Cosa si cerca

1. **Titoli** dei brani
2. **Nomi delle band**
3. **Versi dei testi** (originali e traduzioni)

La ricerca per verso è la ragione per cui la ricerca esiste: chi ricorda "quella frase" non ricorda il titolo.

## Indice

- Costruito **a build time** e salvato in `public/` come **JSON** unico: nessun database, nessun servizio esterno.
- Ogni voce contiene: tipo (`track` \| `band`), titolo, band, lingue disponibili, URL, e i versi normalizzati.
- **Normalizzazione**: minuscole, rimozione dei segni diacritici, spazi collassati, punteggiatura dei versi ignorata. La normalizzazione si fa **a build time**, non a ogni ricerca.
- L'indice è caricato in memoria dal server all'avvio: ricerca istantanea, niente I/O per richiesta.

## Endpoint

`GET /it/api/cerca?q=<query>`

- Restituisce un **frammento HTML** già renderizzato (per HTMX), non JSON: così il risultato è nella stessa lingua dell'interfaccia e non serve JS per dipingerlo.
- Massimo **10 risultati** per richiesta. Se ce ne sono altri, l'ultima voce è "affina la ricerca", non una paginazione.
- Query sotto i **2 caratteri**: risposta vuota, nessuna richiesta inutile.
- Rate-limit per IP per evitare di trasformare la ricerca in un endpoint da abusare.

## Interfaccia

- Campo nell'header; digitando, HTMX aggiorna il dropdown con **debounce ~200 ms**.
- Ogni risultato mostra **titolo del brano + band** (D56), mai uno snippet del verso: il risultato deve restare compatto e non rivelare mezze frasi fuori contesto.
- Navigazione da **tastiera** (frecce, invio) e ruoli ARIA corretti (`combobox`/`listbox`): è un dropdown, non un div cliccabile.
- **Esc** chiude e svuota; click fuori chiude.
- Senza risultati: messaggio breve, non un errore.

## Ordinamento (rilevanza)

1. Match nel **titolo**
2. Match nel nome della **band**
3. Match nel **verso**

A parità di categoria, prima i brani con più lingue disponibili, poi l'ordine alfabetico. Nessun punteggio opaco: l'ordine è spiegabile a voce.

## Cosa la ricerca NON fa

- Nessuna pagina di risultati dedicata (D55).
- Nessuna ricerca semantica, nessun sinonimo, nessuna tolleranza agli errori di battitura: se serve, si aggiunge dopo e sarà una decisione nuova (vedi `DECISION.md`).
- Nessuna cronologia di ricerca, nessun suggerimento personalizzato: non c'è account e non si traccia l'utente.

## Dove si usa la ricerca

- Header, su tutte le pagine.
- **Pagina 404** (D53): chi sbaglia URL trova subito il campo attivo.

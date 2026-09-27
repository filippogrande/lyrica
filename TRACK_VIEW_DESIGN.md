# TRACK_VIEW_DESIGN — la vista a fronte

È la pagina per cui esiste il sito. Originale a sinistra, traduzione a destra.

## Struttura della pagina

```
┌─ header (sticky) ───────────────────────────────┐
├─ breadcrumb: Home > Band > Album > Brano ───────┤
├─ titolo brano + band + badge lingue ────────────┤
├─ [ selettore lingua:  DE | IT | EN ] ───────────┤
├──────────────────────┬──────────────────────────┤
│  ORIGINALE           │  TRADUZIONE              │
│  (nome cantante)     │  (nome cantante)         │
│  strofa 1            │  strofa 1                │
│                      │                          │
│  (nome cantante)     │  (nome cantante)         │
│  strofa 2            │  strofa 2                │
└──────────────────────┴──────────────────────────┘
```

## Desktop (`lg` e oltre)

- Due colonne **affiancate** (`col-lg-6` ciascuna), allineate in alto.
- Le strofe corrispondenti stanno **alla stessa altezza**: la griglia si costruisce **per strofa** (riga per riga), non con due colonne di testo indipendenti. Se una strofa ha un numero di versi diverso dall'originale, la riga cresce e resta allineata in alto.
- Nessuna linea verticale di separazione pesante: lo spazio fa il lavoro.

## Mobile (sotto `lg`)

- **Una colonna**: prima l'originale, poi la traduzione, sempre nella stessa pagina (nessun tab nascosto: il testo si scorre).
- Il **selettore lingua resta sticky** sotto l'header: cambiare lingua non deve richiedere di risalire la pagina.
- Blocchi separati da titoletti "Originale" / "Traduzione" in maiuscoletto.

## Strofe e cantanti

- Ogni strofa è un blocco `<div>` con i suoi versi.
- Il **tipo di strofa non si stampa mai** (niente "Ritornello", "Bridge"): decisione D45.
- Se il brano ha **più voci**, sopra ogni strofa compare **il nome del cantante** — su **ogni** strofa, anche quella dopo un ritornello. Se il brano ha una voce sola, nessuna etichetta compare (sarebbe rumore).
- Il nome del cantante si mostra **anche nella colonna della traduzione**, così le due colonne restano leggibili in parallelo.

## Selettore lingua

- Mostra **solo le lingue realmente presenti** per quel brano (D30), comprese le lingue originali quando sono più di una (D28).
- L'**originale** è sempre selezionabile; non esistono stati vuoti.
- Implementazione: **tutte le lingue presenti nel DOM** e commutate lato client con poche righe di JS (attributo `hidden`), **nessuna chiamata al server**: il cambio è istantaneo e funziona **offline** (la PWA copia l'intera pagina).
- La lingua scelta si riflette nell'URL come parametro (`?lang=de`) via `history.replaceState`, così il link è condivisibile: nessuna pagina separata per lingua, nessun URL duplicato.
- Il selettore è usabile da tastiera e annuncia il cambio (aria-live) agli screen reader.

## Note dentro la pagina

- **Nessuna annotazione per verso** (D33). Il testo resta pulito.
- Una versione **live/acustica** dello stesso brano è un blocco dello stesso testo con la nota "versione live": resta nella stessa pagina (D32).
- Il blocco "informazioni" (traduttore, lingue, data di aggiunta) sta **in fondo**, non sopra il testo: chi è qui vuole leggere, non leggere metadati.

## Casi limite

| Caso | Resa |
|---|---|
| Brano con una sola lingua (solo originale, nessuna traduzione) | **non ha pagina**: compare in tracklist come "solo originale" |
| Brano strumentale | in tracklist con nota, nessuna pagina |
| Traduzione con strofe più corte | la colonna resta allineata in alto, nessun riempimento artificiale |
| Numero di voci mancante | nessuna etichetta cantante |
| Lingua con alfabeto diverso | il font multialfabeto deve coprirla, altrimenti si segnala come bug |

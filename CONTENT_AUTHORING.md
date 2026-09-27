# CONTENT_AUTHORING — come si aggiunge un contenuto

Il flusso normale, in ordine. Non si scrivono file a mano partendo da zero: si generano con la CLI e si riempiono (D20).

## 1. Nuova band

```
lyrica new band "Nome Band"
```

Crea `content/bands/<slug>/band.md` con il front-matter precompilato (nome, slug, `original_langs`, `tags` vuoti, `description` vuota) e apre il file nell'editor.

Da compilare:
- `tags`: **liberi** (D25) — scrivi il genere che descrive davvero la band; ogni tag genera la pagina `/it/tag/<tag>`.
- `original_langs`: una o più lingue originali (D28).
- `description`: 1-3 righe, Markdown semplice.

## 2. Nuovo album

```
lyrica new album "Nome Band" "Titolo Album" --year 1995
```

Crea `content/bands/<band>/<album>/album.md` con la tracklist vuota da riempire.

- Aggiungi i brani in `tracks:` **nell'ordine dell'album**.
- I brani strumentali si segnano subito con `instrumental: true` (D31).
- La cover va messa in `covers/<album-slug>.webp`, **600x600 quadrata** (D21). Se non ce l'hai, lascia perdere: il layout va senza immagine (D36), non mettere placeholder.

## 3. Nuovo brano

```
lyrica new brano "Nome Band" "album-slug" "Titolo Brano"
```

Crea `tracks/<track-slug>.md` con lo scheletro dei blocchi lingua.

Poi, a mano e con calma:
- scrivi il testo originale, **strofa per strofa**, un verso per riga;
- se il brano ha più voci, aggiungi `singer:` a **ogni** strofa;
- scrivi la traduzione **completa**: strofe della stessa lunghezza dell'originale (una traduzione parziale non si pubblica, D29);
- `added_date` = il giorno in cui la pubblichi;
- niente annotazioni per verso (D33).

## 4. Verifica in locale

```
lyrica build      # valida i contenuti e genera public/
lyrica serve      # serve il risultato su una porta locale
```

Se la validazione passa e la pagina si apre come deve, il contenuto è pronto.

## 5. Pubblicazione

- Commit su un **branch** (`content/<band-slug>` è il nome consigliato).
- **Pull Request** verso `main` (D09).
- La CI **valida i contenuti** (vedi `BUILD_PIPELINE.md`); il merge in `main` fa deploy (D10).

## Checklist prima di aprire la PR

- [ ] La traduzione è **completa** strofa per strofa.
- [ ] `added_date` è la data di pubblicazione reale.
- [ ] Slug, cartelle e nome della cover **coincidono**.
- [ ] I nomi dei cantanti sono solo nei brani multi-voce, e su ogni strofa.
- [ ] Nessuna annotazione per verso.
- [ ] Il contenuto si legge bene **da telefono**, non solo da desktop.
- [ ] Se un brano non ha traduzioni, è segnato come "solo originale" e non è linkato (vedi `SPEC.md`).

## Aggiungere una lingua di traduzione a un brano esistente

1. Apri `tracks/<slug>.md`.
2. Aggiungi un blocco con `role: translation`, `lang: <codice>`, `translator:`, e le strofe **complete**.
3. Il validatore controlla che le strofe combacino con l'originale: se non combaciano, la PR è rossa.

Il selettore lingua in pagina si aggiorna da solo: non c'è niente da configurare.

## Aggiungere una lingua all'interfaccia

1. Copia `locales/it.yaml` in `locales/<lang>.yaml` e traduci i valori (non le chiavi).
2. Il build genera le pagine `/<lang>/...` e le dichiara con hreflang.
3. Se una chiave manca, la CI segnala l'errore (vedi `I18N_DESIGN.md`).

## Rinominare uno slug (cambio URL)

1. Rinomina cartella o file.
2. Aggiungi la **regola di redirect 301** nel file dei redirect (D39): vecchio URL → nuovo URL.
3. Verifica che il vecchio URL risponda con redirect nella build locale.

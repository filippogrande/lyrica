# CSS di terze parti (vendored)

Questa cartella contiene il CSS di terze parti **servito dal repo** (non da CDN),
secondo `DESIGN_SYSTEM.md`.

## `bootstrap.min.css`

- Versione dichiarata: **Bootstrap 5.3.8** (solo CSS, nessun bundle JS).
- **Non è ancora nel repo**: il file viene scaricato con il comando qui sotto.
  Finché manca, la pagina si apre senza stili di Bootstrap (il CSS custom
  `assets/css/lyrica.css` è comunque presente).

```
./scripts/fetch-assets.sh
```

Il file va poi committato: così la build e l'offline della PWA non dipendono da
una rete esterna. Quando si cambia versione si aggiornano insieme
`scripts/fetch-assets.sh` e questa nota.

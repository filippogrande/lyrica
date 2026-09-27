# ADS_DESIGN — pubblicità

## Stato al lancio

**Spente.** Gli slot esistono nel codice ma `enabled: false` in `ads.yaml`: nessuno spazio pubblicitario viene renderizzato finché non è una decisione esplicita (D64).

Quando sono spente, **il DOM non contiene nemmeno i contenitori**: il testo prende tutta la larghezza disponibile e il layout non ha buchi.

## Principi (se e quando si attivano)

1. **Non invasive**: piccole, statiche, mai sopra il testo del brano, mai popup, mai interstitial, mai sticky.
2. **Nessun tracking**: banner serviti dallo stesso dominio, nessuno script di terze parti, nessun cookie di profilazione.
3. **Coerenti con `LEGAL_DESIGN.md`**: se una pubblicità introducesse tracciamento, servirebbero consenso e banner — quindi non si introduce.
4. **Includono il budget**: lo spazio occupato conta nel budget di ~300 KB per pagina.

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
| `side_left` / `side_right` | ai **lati** delle due colonne di testo | verticale stretto | **solo da `xl` in su** |
| `track_bottom` | dopo "chiedi altre canzoni" | banner largo e basso | sempre |

Regole delle colonne laterali:

- Il **testo non si stringe** per far posto agli ads: sotto `xl` le colonne laterali **non esistono** e il testo usa tutto lo spazio.
- Le colonne laterali **non sono sticky**: scorrono con la pagina. Un banner che segue lo scroll è la definizione di invasivo.
- Non si mettono **tra** le due colonne del testo (in mezzo all'originale e alla traduzione): spezzerebbe la lettura a fronte, che è il motivo per cui il sito esiste.
- Niente ads **dentro** il blocco delle strofe.

### Altre pagine

| Slot | Dove |
|---|---|
| `top` | sotto l'header nelle pagine di elenco |
| `footer` | sopra il footer |

Gli slot pubblicitari **non compaiono** su 404, pagine legali e form.

## Configurazione

`ads.yaml` contiene:

- `enabled: false` — interruttore globale;
- per ogni slot: se attivo, il frammento HTML o l'immagine locale, il link di destinazione, l'etichetta "sponsorizzato".

Le immagini pubblicitarie vivono nel repo, come le cover: il sito non chiama domini terzi.

## Requisiti prima di attivare

- [ ] Ricontrollare `LEGAL_DESIGN.md` (informativa e cookie).
- [ ] Etichettatura "sponsorizzato" visibile e leggibile.
- [ ] Verifica del peso aggiunto rispetto al budget.
- [ ] Controllo da telefono **e** da schermo largo: il testo resta larghezza piena sotto `xl` e non si comprime sopra.
- [ ] Verifica che nessuno slot copra o interrompa il testo.

## Cosa non si farà mai

- Nessun annuncio **sticky** o che segue lo scroll.
- Nessun annuncio tra originale e traduzione, o dentro le strofe.
- Nessun annuncio con audio o video con autoplay.
- Nessuna raccolta di dati per profilazione.

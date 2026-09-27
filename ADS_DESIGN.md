# ADS_DESIGN — pubblicità

## Stato al lancio

**Spente.** Gli slot esistono nel codice ma `enabled: false` in `ads.yaml`: nessuno spazio pubblicitario viene renderizzato finché non è una decisione esplicita (D64).

## Principi (se e quando si attivano)

1. **Non invasive**: piccole, statiche, mai sopra il testo del brano, mai popup, mai interstitial.
2. **Nessun tracking**: banner serviti dallo stesso dominio, nessuno script di terze parti, nessun cookie di profilazione.
3. **Coerenti con `LEGAL_DESIGN.md`**: se una pubblicità introducesse tracciamento, servirebbero consenso e banner — quindi non si introduce.
4. **Includono il budget**: lo spazio occupato conta nel budget di ~300 KB per pagina.

## Posizionamento previsto

| Slot | Dove | Formato |
|---|---|---|
| `top` | sotto l'header delle pagine di elenco | banner orizzontale piccolo |
| `side` | colonna laterale, solo da `xl` in su | verticale, piccolo |
| `footer` | sopra il footer | orizzontale piccolo |

**Mai nella pagina brano accanto al testo**: la vista a fronte è il motivo per cui il sito esiste e non va disturbata.

## Configurazione

`ads.yaml` contiene:

- `enabled: false` — interruttore globale;
- per ogni slot: se attivo, il frammento HTML o l'immagine locale, il link di destinazione, l'etichetta "sponsorizzato".

Le immagini pubblicitarie vivono nel repo, come le cover: il sito non chiama domini terzi.

## Requisiti prima di attivare

- [ ] Ricontrollare `LEGAL_DESIGN.md` (informativa e cookie).
- [ ] Etichettatura "sponsorizzato" visibile e leggibile.
- [ ] Verifica del peso aggiunto rispetto al budget.
- [ ] Controllo che nessuno slot copra il testo su mobile.
- [ ] Gli slot restano **assenti** su 404, legali e form (non ha senso pubblicizzare lì).

## Cosa non si farà mai

- Nessun annuncio che si sovrappone al testo, nessun annuncio con audio, nessuna raccolta di dati per profilazione.

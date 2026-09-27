# FEED_AND_STATS — RSS e contatori

## Feed RSS

- **Un feed per lingua dell'interfaccia**: `/it/rss.xml` (e `/en/rss.xml` quando ci sarà l'inglese), coerente con il path prefix di `I18N_DESIGN.md`.
- Contiene le **ultime traduzioni pubblicate**, ordinate per `added_date` decrescente, con un limite ragionevole (20 voci: chi vuole di più ha il sito).
- Ogni voce: titolo del brano, band, link alla pagina, data, breve descrizione generata (lingue disponibili), e autore della traduzione se presente.
- Nel `<head>` di ogni pagina: `<link rel="alternate" type="application/rss+xml">`.
- Link al feed nel **footer** (D49).
- Il feed è **statico**, generato a build time: nessun servizio di feed.

### Criteri

- Una voce compare **quando la traduzione è pubblicata**, non quando il brano viene creato a metà.
- Niente feed per singola band al lancio: se servirà, sarà una decisione nuova.

## Contatori pubblici

Mostrano il valore reale del sito a colpo d'occhio (D62):

| Contatore | Definizione esatta |
|---|---|
| **Brani tradotti** | brani con almeno una traduzione pubblicata |
| **Band** | band con almeno un album visibile |
| **Lingue** | lingue di traduzione presenti nei contenuti |

- Sono **calcolati a build time** dai contenuti: nessun numero scritto a mano che possa divergere dalla realtà.
- Comprimono **bene**: l'obiettivo è mostrare che il sito cresce, non gonfiare i numeri. Un brano con 4 lingue conta 1 brano, non 4.
- Dove: **homepage** e **footer** (D49).
- Nessun contatore di visite: le statistiche sono di Umami (`ANALYTICS_DESIGN.md`), non un numero da esibire.

## Perché questi due insieme

Sono le due facce della stessa idea: il sito si aggiorna nel tempo e chi lo segue non deve andarlo a controllare. Il feed serve a chi lo segue, i contatori a chi arriva per la prima volta.

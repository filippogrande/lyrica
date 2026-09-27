# ANALYTICS_DESIGN — Umami

## Perché Umami

- **Cookieless**: nessun cookie, nessun identificatore persistente → **nessun cookie banner** necessario.
- **Riuso dell'istanza esistente** (D22): nessun container nuovo da mantenere.
- Misura ciò che serve sapere (cosa viene letto, cosa si cerca) **senza** profilare nessuno, coerentemente con "zero account" (D08).

## Setup

1. Nell'istanza Umami già attiva si crea un **nuovo website** dedicato a Lyrica (`lyrica.filippomoscatelli.com`).
2. Si ottiene lo **script tag** e lo si inserisce nel `<head>` di tutte le pagine generate (un solo punto nei template).
3. Lo script è **asincrono** e non blocca il render.

## Custom event

| Evento | Quando | A cosa serve |
|---|---|---|
| `search` | una ricerca produce risultati | capire cosa la gente cerca e non trova |
| `search_empty` | una ricerca non produce risultati | materia prima per le prossime traduzioni |
| `lang_switch` | cambio della lingua del testo nella pagina brano | quali traduzioni servono davvero |
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

È una **nuova decisione** da registrare in `DECISION.md`, non un'estensione silenziosa. Il vincolo di partenza (niente cookie, niente banner) non si rompe in una PR di routine.

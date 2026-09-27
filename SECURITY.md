# SECURITY — modello di minaccia

## Il vantaggio di partenza

Il sito **non ha login, non ha account, non ha database e non ha un pannello admin** (D08, D16). Non esiste una superficie di autenticazione da attaccare, non esistono credenziali utente da rubare, non esiste un archivio da esfiltrare. È la misura di sicurezza più efficace del progetto e va **difesa**, non erosa: ogni proposta di aggiungere autenticazione o storage è una decisione esplicita da registrare in `DECISION.md`.

## Superficie reale

1. **Contenuti** (autore, non utenti): se un testo contenesse markup, il rischio è XSS verso i lettori.
2. **I due form** (`/api/segnala`, `/api/contatti`): abuso, spam, flood.
3. **Token del bot Telegram**: se trapelato, qualcuno può scrivere nella chat dell'autore.
4. **Endpoint di ricerca**: abuso leggero (scraping dell'indice).
5. **Dipendenze**: Bootstrap vendored, moduli Go, immagini base Docker.

## Contromisure

| Area | Misura |
|---|---|
| Output | **Escaping di default** in Templ; mai `unsafe`/HTML grezzo dai contenuti; Markdown renderizzato con sanitizzazione |
| Header | HSTS, `X-Content-Type-Options: nosniff`, `Referrer-Policy: strict-origin-when-cross-origin`, `Permissions-Policy` minimale, **CSP** senza `unsafe-inline` (lo script del tema è inline: si autorizza con hash specifico, non con `unsafe-inline`) |
| Form | **Rate-limit per IP** in memoria + **honeypot** + campo tempo minimo di compilazione; nessun salvataggio dei dati |
| Segreti | Token Telegram solo in `.env`/variabili d'ambiente, **mai** nel repo, mai nei log, mai nei messaggi di errore |
| Dipendenze | `.github/dependabot.yml` obbligatorio e aggiornato nella stessa PR che aggiunge un manifest |
| Esposizione | Il servizio sta dietro **Cloudflare Tunnel**: nessuna porta aperta in casa, nessun IP del home-lab esposto |
| Log | Log di accesso minimali, **nessun contenuto dei form** e nessuna retention lunga degli IP |

## Cosa NON si fa

- Non si aggiunge autenticazione "per sicurezza": qui l'autenticazione **crea** superficie, non la riduce.
- Non si salvano i dati dei form "tanto per averli": è un archivio in più da difendere e da cancellare su richiesta.
- Non si aprono porte del home-lab per far passare un servizio: si usa il tunnel.
- Non si committa un `.env`, nemmeno di esempio con valori reali (esempio con valori finti sì, con nomi di variabile chiari).

## Comportamento in caso di abuso

1. Il rate-limit risponde **429** senza spiegare la soglia esatta.
2. Se lo spam diventa sistematico: si disattiva temporaneamente l'endpoint e si valuta un captcha leggero **come decisione**, non come patch silenziosa.
3. Se il token Telegram è sospetto: si rigenera il token dal BotFather e si aggiorna `.env` sul server.

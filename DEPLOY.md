# DEPLOY — Docker, Cloudflare Tunnel, CI

## Dove gira

Home-lab di Filippo, su **Docker**, esposto da **Cloudflare Tunnel**. **Non** k3s (in dismissione, D15).
URL pubblico: `lyrica.filippomoscatelli.com` (D03).

## Immagine

`Dockerfile` **multi-stage**:

- **Stage build**: immagine Go → scarica le dipendenze, compila il binario, esegue `lyrica build` sui contenuti.
- **Stage finale**: immagine minimale (distroless o alpine) con **solo** il binario e `public/`. Niente toolchain, niente sorgenti.
- L'immagine finale non contiene né segreti né contenuti non pubblicati.

## Configurazione

- `docker-compose.yml` con il servizio `lyrica`, `restart: unless-stopped`, **healthcheck** sull'endpoint di salute.
- Il container **non espone porte all'host** se non quelle necessarie al tunnel: l'accesso pubblico passa dal tunnel.
- Variabili d'ambiente (via `.env` non committato):
  - token del **bot Telegram** e chat id di destinazione;
  - porta interna del server;
  - eventuale dominio pubblico per la generazione degli URL assoluti (sitemap, feed).
- **Nessun segreto nel repo** (vedi `SECURITY.md`).

## Flusso di deploy

1. PR aperta → job **validate** sulla CI (vedi `BUILD_PIPELINE.md`).
2. Merge su `main` (D10).
3. GitHub Actions builda l'immagine e la pubblica nel registry.
4. SSH sul home-lab → `docker compose pull && docker compose up -d`.
5. Verifica automatica: healthcheck + risposta 200 sulla home.
6. Se lo smoke test fallisce, il job è rosso: **si fa rollback** con il tag dell'immagine precedente (nessun rollback automatico al lancio, decisione dell'autore).

## Aggiornamenti

- **Nessun Watchtower** al lancio: un deploy deve essere una decisione, non un evento notturno.
- L'aggiornamento di sistema (base image, Bootstrap) segue le PR di Dependabot.

## Rischi noti e cosa si rompe se sbagli

| Azione | Rischio |
|---|---|
| Merge su `main` senza aver aperto il sito in locale | va in produzione un layout rotto |
| `docker compose down -v` | cancella il container e il suo stato (qui non c'è database: il danno è solo il downtime) |
| Rotazione del token Telegram senza aggiornare `.env` | le segnalazioni smettono di arrivare, il sito resta su |
| Modifica al tunnel Cloudflare | il sito sparisce dall'URL pubblico pur essendo su |

## Backup

Non c'è niente da salvare sul server: **i contenuti sono nel repo**.
Il backup del sito è il repository. Da qui la scelta (D14) di non avere database.

## Cosa documentare dopo il primo deploy reale

Questo file descrive il deploy **previsto**. Alla prima messa online (FASE 5) vanno aggiunti qui: nome esatto dell'immagine, porta interna, comando di rollback verificato e la voce del tunnel Cloudflare usata.

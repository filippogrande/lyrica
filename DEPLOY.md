# DEPLOY — Docker, GitHub Actions, home-lab

## Dove gira

Home-lab di Filippo, **Docker**, esposto in rete locale; il tunnel **Cloudflare** e il dominio `lyrica.filippomoscatellis.com` arrivano quando il sito ha contenuti (FASE 5). **Non** k3s: in dismissione (D15).

Esiste un container **dalla FASE 1**: il sito si avvia con `docker compose` fin da subito, appena c'è un'immagine pubblicata.

## Immagine

- Registry: **Docker Hub**, `filippogrande/lyrica` (stessa convenzione degli altri repo).
- Tag: `latest` e `<sha del commit>` — il tag con la SHA serve al rollback.
- `Dockerfile` **multi-stage**:
  - stage `build`: `golang:1.27.1-alpine3.24` → dipendenze, `templ generate`, script asset (Bootstrap 5.3.8 dentro l'immagine), binario statico `CGO_ENABLED=0`;
  - stage finale: `alpine:3.24.2` + `ca-certificates`, utente non privilegiato (`lyrica`, uid 10001), **solo** binario, `assets/` e `locales/`.
- `HEALTHCHECK` sull'endpoint `/healthz`.

## Build e push (automatici)

Workflow `.github/workflows/docker-build-push.yml`, su **push su `main`**:

1. login su Docker Hub (secrets `DOCKERHUB_USERNAME`, `DOCKERHUB_TOKEN`);
2. build dell'immagine;
3. push di `filippogrande/lyrica:latest` e `:sha`.

L'immagine **non** viene buildata sulle PR: una PR non deve poter pubblicare un tag `latest`.

## .env: cosa ci va e cosa non ci va

Il compose legge un file **`.env`** accanto a sé (`env_file`), che **non sta nel repo**: `.gitignore` esclude `.env` e `.env.*`, ma tiene `.env.example`. Nel repo c'è solo il **template** [`.env.example`](.env.example).

Sul home-lab:

```
cp .env.example .env
```

Regole:

- i **segreti** (token del bot Telegram, ecc.) stanno **solo** nel `.env`, mai nel compose, mai nel repo;
- il compose resta valido anche **senza** `.env` (`required: false`): il container parte con i default;
- se una variabile non è ancora letta dall'applicazione, resta **commentata** nel template: nessun valore finto che finge di funzionare.

Variabili previste: `LYRICA_ADDR` (attiva), `LYRICA_BASE_URL` (FASE 4/5), `UMAMI_URL`/`UMAMI_SITE_ID` (FASE 5), `TELEGRAM_BOT_TOKEN`/`TELEGRAM_CHAT_ID` (FASE 6).

## Deploy sul home-lab (manuale)

Compose atteso in **`/mnt/applicazioni/yml/docker/lyrica/docker-compose.yml`**, con il `docker-compose.yml` e il `.env` di questo repo copiati lì:

```
cp .env.example .env
docker compose pull
docker compose up -d
```

Il container ascolta sulla **8080 interna**, mappata sull'host sulla **8085**.

> **Da confermare al primo avvio**: la porta 8085 sull'host è libera. Se è occupata, si cambia la mappa nel compose (le 8090-8095 sono già usate da altri servizi; Piped usa 8090/8091/8093, release-monitor 8095, trip 8092).

## Verifica dopo il deploy

1. `docker compose ps` → stato `running` e healthcheck `healthy`.
2. `curl -I http://localhost:8085/` → deve rispondere **302** verso `/it/`.
3. `curl http://localhost:8085/healthz` → `ok`.
4. `docker compose logs --tail=50 web` → nessun errore di rendering.

## Rollback

Con il tag SHA dell'ultimo commit buono:

```
docker compose down
docker compose up -d --pull always
```

Per fissare una versione precisa si usa l'immagine `filippogrande/lyrica:<sha>` al posto di `latest` nel compose, e si riporta a `latest` quando il problema è risolto.

## Cosa può rompersi (e cosa si perde)

| Azione | Rischio |
|---|---|
| `docker compose down` | il sito sparisce fino al `up -d`: nessun dato perso (non c'è database) |
| Porta 8085 già occupata | il container non parte: `Bind for 0.0.0.0:8085 failed` |
| Immagine `latest` con un bug | il compose riparte con la stessa immagine rotta: si fa rollback al tag SHA |
| Secrets Docker Hub assenti/scaduti nel repo | la CI è rossa e `latest` resta quella vecchia (il sito non si rompe, non si aggiorna) |
| `.env` con una variabile sbagliata | il container parte ma l'app non funziona come previsto: si controlla `docker compose logs` |

## Perché non c'è deploy automatico dal runner

Il deploy è **una decisione**, non un effetto collaterale di un push: il runner non ha accesso al home-lab e non deve averlo. Il workflow si ferma alla pubblicazione dell'immagine.

## Backup

Non c'è niente da salvare sul server: i contenuti sono nel repo. Il backup del sito è il repository (D14).

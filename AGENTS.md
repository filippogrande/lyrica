# AGENTS — Lyrica

File che si legge **sempre per primo** in questo repo. Da qui si arriva a un file di
codice in **massimo 2 salti**: questo file → `docs/REFERENCE.md` → il file.

## Che cos'è

Sito **statico** di traduzioni di testi di canzoni: originale a sinistra, traduzione
selezionabile a destra. Go + Templ, contenuti come Markdown con front-matter YAML in
`content/`, build a build time in `public/`, immagine Docker su Docker Hub a ogni merge su
`main`.

- **Nessun login, nessun account, nessun database, nessun pannello admin.**
- Nessun contributo esterno: le traduzioni le scrive solo l'autore.
- Codice in inglese, contenuti e documentazione in italiano.

## Prima di scrivere codice

1. Leggi [`docs/REFERENCE.md`](docs/REFERENCE.md): è la mappa richiesta → file → simbolo,
   verificata sul codice.
2. Leggi il doc di competenza per l'area che tocchi (indice sotto).
3. Se codice e doc non si accordano, **il doc si aggiorna nella stessa PR**. Un doc che
   mente è peggio di nessuna doc: va corretto o cancellato, non annotato.

## Comandi

```
bash scripts/fetch-assets.sh                              # una volta (Bootstrap 5.3.8)
go mod tidy

go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate   # dopo ogni modifica ai .templ

go run ./cmd/lyrica build        # valida i contenuti e genera public/
go run ./cmd/lyrica serve        # http://localhost:8080 (richiede public/)
go run ./cmd/lyrica valida       # valida i contenuti senza generare nulla
go test ./...
```

Senza `templ generate` il progetto **non compila**; i file `*_templ.go` non si committano.
Le versioni esatte di **Go** e **templ** stanno in `go.mod`, `ci.yml` e `Dockerfile`: se ne
cambi una, si aggiornano **tutti e tre nella stessa PR** (e il Dockerfile con la stessa
patch di Go della CI).

Deploy: `docker compose pull && docker compose up -d` (compose in `deploy/`, porta host
**8097**). Merge su `main` = nuova immagine su Docker Hub.

## Paletti non negoziabili

- **Massimo 500 righe per file, massimo 50 per funzione.** Nessuna eccezione.
- **Vietati i fallback silenziosi**: un errore si propaga e si vede. Unica eccezione
  progettata: lingua dell'interfaccia non disponibile → italiano.
- **Mai push diretto su `main`**: branch + PR. Nomi branch `feat|fix|docs|ci|content/<slug>`.
- **Segreti solo in `.env`/variabili d'ambiente** (token Telegram, Umami), mai nel repo.
- **Vietati i file di test usa-e-getta**: un test entra nel repo solo se prova una regola.
- La verifica si fa sull'**artefatto** (pagina generata, output del comando, container che
  risponde), non su un frammento di codice.
- Ogni modifica corrisponde a una task tracciata (Vikunja, progetto `Lyrica`).

## Indice dei documenti

| File | Tipo | Cosa contiene |
|---|---|---|
| [`AGENTS.md`](AGENTS.md) | how-to | questo file: identità, comandi, paletti, indice |
| [`docs/REFERENCE.md`](docs/REFERENCE.md) | reference | **la mappa**: richiesta → file → entry point |
| [`docs/GUIDELINES.md`](docs/GUIDELINES.md) | how-to | regole di lavoro, workflow git, paletti non negoziabili |
| [`docs/DECISIONS.md`](docs/DECISIONS.md) | reference | registro delle decisioni (D01→), con il perché |
| [`docs/ROADMAP.md`](docs/ROADMAP.md) | reference | fasi, artefatti, criteri di chiusura |
| [`docs/SPEC.md`](docs/SPEC.md) | reference | cosa fa il sito: pagine, funzionalità, casi limite |
| [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) | explanation | come è fatto: componenti, build, Docker, deploy, sicurezza |
| [`docs/CONTENT.md`](docs/CONTENT.md) | how-to | schema dei contenuti e come si aggiungono |
| [`docs/FRONTEND.md`](docs/FRONTEND.md) | explanation | interfaccia, pagina brano, lingue, ricerca, PWA |
| [`docs/FEATURES.md`](docs/FEATURES.md) | explanation | ads, analytics, legali, form, feed e contatori |

`DECISIONS.md` e `ROADMAP.md` **non si toccano**: non si riscrivono né si consolidano.
Se una decisione registrata non corrisponde più al codice, si **segnala nella PR**, non si
corregge nel registro. Se il piano di fase è morto, si segnala prima di proporre il taglio.

## Convenzione di navigazione

Una richiesta ("dove vive la validazione dei contenuti?") si risolve con **indice → 1 file**,
mai aprendo tutto. La mappa si aggiorna **nella stessa PR** che sposta un file o rinomina un
simbolo: un doc che punta a un simbolo inesistente è un doc che mente.

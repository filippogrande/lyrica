# Lyrica — immagine Docker.
#
# Stage di build in tre passaggi: Go + templ, asset di terze parti, SITO
# GENERATO. L'ultimo è quello che rende pubblico il contenuto: senza
# `lyrica build` la cartella public/ non esiste e il server non parte
# (controllo esplicito in internal/web).
#
# Ordine dei COPY studiato per non invalidare la cache delle dipendenze quando
# cambia un contenuto o un asset.

# --- Stage 1: build -----------------------------------------------------
FROM golang:1.27.1-alpine3.24 AS build

WORKDIR /src

# curl e bash servono a scripts/fetch-assets.sh (Bootstrap pinnato, niente CDN
# a runtime). bash NON e' opzionale: lo script ha shebang bash e usa
# "set -euo pipefail", e l'immagine alpine di Go non include bash (fallirebbe
# con "env: bash: No such file or directory").
RUN apk add --no-cache curl bash

# 1) dipendenze: questo layer resta in cache finché non cambia go.mod
COPY go.mod ./
RUN go mod download

# 2) codice, contenuti, locale, asset
COPY . .
RUN go mod tidy

# 3) i file .templ diventano Go: senza questo passo il progetto non compila
RUN go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate

# 4) CSS di terze parti (Bootstrap 5.3.8) dentro l'immagine.
# Invocato come "bash scripts/..." e NON "./scripts/...": il file non ha il bit
# di esecuzione nel repo (i commit passano dalle API), quindi "./" fallisce
# con exit 126 / Permission denied.
RUN bash scripts/fetch-assets.sh

# 5) il sito: valida i contenuti e genera public/. Se un contenuto è invalido
# la build dell'immagine si ferma QUI, prima di pubblicare qualsiasi cosa.
RUN go run ./cmd/lyrica build

# 6) binario statico
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/lyrica ./cmd/lyrica

# --- Stage 2: runtime ---------------------------------------------------
FROM alpine:3.24.2

RUN apk add --no-cache ca-certificates \
    && adduser -D -u 10001 -h /app lyrica

WORKDIR /app

COPY --from=build /out/lyrica /app/lyrica
# Il sito generato: HTML, CSS, JS, cover. Il server serve SOLO questo, più
# /healthz e il redirect della radice.
COPY --from=build /src/public /app/public
COPY --from=build /src/locales /app/locales

USER lyrica

ENV LYRICA_ADDR=:8080
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://127.0.0.1:8080/healthz || exit 1

ENTRYPOINT ["/app/lyrica"]
CMD ["serve"]

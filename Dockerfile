# Lyrica — immagine Docker.
#
# Due stage: build (Go + templ + asset) e runtime minimale (alpine + binario).
# Ordine dei COPY studiato per non invalidare la cache delle dipendenze quando
# cambia un contenuto o un asset (BUILD_PIPELINE.md).

# --- Stage 1: build -----------------------------------------------------
FROM golang:1.27.1-alpine3.24 AS build

WORKDIR /src

# curl serve a scripts/fetch-assets.sh (Bootstrap pinnato, niente CDN a runtime).
# bash NON e' opzionale: lo script ha shebang bash e usa "set -euo pipefail",
# e l'immagine alpine di Go non include bash (fallirebbe con
# "env: bash: No such file or directory").
RUN apk add --no-cache curl bash

# 1) dipendenze: questo layer resta in cache finché non cambia go.mod
COPY go.mod ./
RUN go mod download

# 2) codice e asset
COPY . .
RUN go mod tidy

# 3) i file .templ diventano Go: senza questo passo il progetto non compila
RUN go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate

# 4) CSS di terze parti (Bootstrap 5.3.8) dentro l'immagine
RUN ./scripts/fetch-assets.sh

# 5) binario statico
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/lyrica ./cmd/lyrica

# --- Stage 2: runtime ---------------------------------------------------
FROM alpine:3.24.2

RUN apk add --no-cache ca-certificates \
    && adduser -D -u 10001 -h /app lyrica

WORKDIR /app

COPY --from=build /out/lyrica /app/lyrica
COPY --from=build /src/assets /app/assets
COPY --from=build /src/locales /app/locales

USER lyrica

ENV LYRICA_ADDR=:8080
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://127.0.0.1:8080/healthz || exit 1

ENTRYPOINT ["/app/lyrica"]
CMD ["serve"]

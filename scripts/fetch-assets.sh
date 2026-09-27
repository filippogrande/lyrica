#!/usr/bin/env bash
#
# Scarica gli asset di terze parti, pinnati alla versione dichiarata.
# Bootstrap è servito dal repo e non da una CDN: serve all'offline della PWA
# e a non dipendere da terzi a runtime (DESIGN_SYSTEM.md, PWA_DESIGN.md).
#
# Uso: ./scripts/fetch-assets.sh

set -euo pipefail

BOOTSTRAP_VERSION="5.3.8"
DEST="assets/css/vendor"
URL="https://cdn.jsdelivr.net/npm/bootstrap@${BOOTSTRAP_VERSION}/dist/css/bootstrap.min.css"

mkdir -p "${DEST}"
curl -fsSL "${URL}" -o "${DEST}/bootstrap.min.css"

echo "bootstrap.min.css ${BOOTSTRAP_VERSION} -> ${DEST}/bootstrap.min.css"
echo "Aggiorna assets/css/vendor/README.md se cambi versione."

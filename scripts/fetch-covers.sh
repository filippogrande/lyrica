#!/usr/bin/env bash
# fetch-covers.sh — scarica UNA immagine dagli archivi pubblici, sapendo il MBID.
#
#   bash scripts/fetch-covers.sh cover <release-group-mbid> <album-slug>
#   bash scripts/fetch-covers.sh band  <artist-mbid>        <band-slug>
#
# Le cover degli album arrivano dal Cover Art Archive (MusicBrainz + Internet
# Archive) e le foto delle band da fanart.tv, che le indicizza per MBID di
# MusicBrainz. Lo script scrive covers/<slug>.webp a 600x600, la sola
# dimensione che il sito usa (docs/CONTENT.md), e non tocca i file di
# contenuto: stampa la riga da aggiungere al front-matter.
#
# Le immagini finiscono NEL REPO: il sito le serve da /covers/, quindi
# l'API si chiama una volta per immagine, mai a runtime e mai a ogni build.
#
# Per riempire tutto quello che manca senza copiare MBID a mano c'è
# scripts/fetch-images.sh: li risolve dal nome della band e dal titolo
# dell'album (D86) e chiama questo script.
#
# Codici di uscita:
#   0  immagine scaricata
#   1  errore (strumento mancante, rete, chiave rifiutata, conversione)
#   2  uso sbagliato
#   3  immagine non disponibile nell'archivio (solo con --allow-missing):
#      non è un errore, è un archivio che quella cover o quella foto non ce l'ha
#
# Regole rispettate (docs/GUIDELINES.md, D76):
#   - nessun fallback silenzioso: se manca uno strumento, l'immagine o la
#     chiave, esce con un errore esplicito e non lascia file a metà;
#   - nessun file temporaneo lasciato in giro: si scarica in una cartella
#     temporanea e la si cancella sempre (trap).
#
# Chiavi fanart.tv (servono solo per le foto delle band, le cover no):
#   FANART_API_KEY     chiave di progetto    -> header api-key    (~/.fanart_api_key)
#   FANART_CLIENT_KEY  chiave personale      -> header client-key (~/.fanart_client_key)
# Basta una delle due. Le chiavi NON stanno nel repo e non vengono mai
# stampate, né finite nei log.
#
# Si invoca con `bash` esplicito, come scripts/fetch-assets.sh: nel repo gli
# script non hanno il bit di esecuzione.

set -euo pipefail

readonly SIZE=600
readonly QUALITY=82
readonly CAA_URL="https://coverartarchive.org"
readonly FANART_API="https://webservice.fanart.tv/v3.2/music"
readonly FANART_KEY_FILE="${HOME}/.fanart_api_key"
readonly FANART_CLIENT_KEY_FILE="${HOME}/.fanart_client_key"
readonly COVERS_DIR="covers"
# MusicBrainz chiede un User-Agent che identifichi chi chiama.
readonly USER_AGENT="lyrica-covers/1.0 (+https://github.com/filippogrande/lyrica)"

IM=""
# Su ImageMagick 7 si usa `magick identify`; su ImageMagick 6 `identify` è un
# binario a sé: chiamare `convert identify` non esiste e la conversione non
# arriva mai in fondo.
IM_IDENTIFY=""
TMP_DIR=""
FANART_PROJECT=""
FANART_PERSONAL=""
ALLOW_MISSING="no"
FORCE="no"
ARGS=()

usage() {
	cat >&2 <<'FINE'
uso:
  bash scripts/fetch-covers.sh cover <release-group-mbid> <album-slug> [--force] [--allow-missing]
  bash scripts/fetch-covers.sh band  <artist-mbid>        <band-slug>  [--force] [--allow-missing]

esempi:
  bash scripts/fetch-covers.sh cover c31a5e2b-0bf8-32e0-8aeb-ef4ba9973932 liebe-ist-fuer-alle-da
  bash scripts/fetch-covers.sh band  cf075492-d880-4afc-b87b-d6b03e33dacc electric-callboy

Il MBID si legge nell'URL della pagina MusicBrainz:
  musicbrainz.org/release-group/<release-group-mbid>   -> cover di un album
  musicbrainz.org/artist/<artist-mbid>                 -> foto di una band

--force           riscrive un file già presente in covers/
--allow-missing   se l'archivio non ha quell'immagine esce con 3 invece che con
                  errore: serve a chi riempie in blocco (scripts/fetch-images.sh)

Chiave fanart.tv (solo per le band): in FANART_API_KEY, FANART_CLIENT_KEY o nei
file ~/.fanart_api_key / ~/.fanart_client_key.

Scrive covers/<slug>.webp a 600x600. Non tocca i file di contenuto: la riga da
aggiungere al front-matter te la stampa alla fine.
FINE
}

die() {
	echo "errore: $*" >&2
	exit 1
}

# missing_or_die distingue "l'archivio non ce l'ha" da "qualcosa è andato
# storto": in blocco la prima cosa si salta, la seconda no.
missing_or_die() {
	if [ "$ALLOW_MISSING" = "yes" ]; then
		echo "saltata: $*" >&2
		exit 3
	fi
	die "$*"
}

cleanup() {
	if [ -n "$TMP_DIR" ]; then rm -rf "$TMP_DIR"; fi
	return 0
}

trap cleanup EXIT

# require_tools verifica ciò che serve PRIMA di scaricare qualcosa, e sceglie
# i nomi giusti dei comandi di ImageMagick (7: magick; 6: convert + identify).
require_tools() {
	command -v curl >/dev/null 2>&1 || die "manca curl"
	command -v python3 >/dev/null 2>&1 || die "manca python3 (serve per leggere la risposta di fanart.tv)"
	if command -v magick >/dev/null 2>&1; then
		IM="magick"
		IM_IDENTIFY="magick identify"
	elif command -v convert >/dev/null 2>&1; then
		IM="convert"
		IM_IDENTIFY="identify"
		command -v identify >/dev/null 2>&1 \
			|| die "ImageMagick 6 senza il binario identify: reinstallalo (niente conversioni a metà)"
	else
		die "manca ImageMagick: installalo (macOS: brew install imagemagick) — niente conversioni a metà"
	fi
}

# check_slug rifiuta gli slug che non possono essere nomi di file in covers/.
check_slug() {
	local slug="$1"
	[ -n "$slug" ] || die "slug vuoto"
	case "$slug" in
	*[!a-z0-9-]*) die "slug non valido: \"$slug\" (minuscolo, ASCII, trattini)" ;;
	esac
}

check_mbid() {
	case "$1" in
	????????-????-????-????-????????????) ;;
	*) die "MBID non valido: \"$1\" (formato 8-4-4-4-12)" ;;
	esac
}

prepare_target() {
	local slug="$1" force="$2" target="${COVERS_DIR}/${1}.webp"
	mkdir -p "$COVERS_DIR"
	if [ -f "$target" ] && [ "$force" != "yes" ]; then
		die "$target esiste già: cancellalo a mano o passa --force"
	fi
	printf '%s' "$target"
}

# to_webp converte in 600x600 webp quadrato (ritaglio centrale) e verifica il
# risultato: se non è 600x600 webp, non si scrive niente.
to_webp() {
	local source="$1" target="$2" dims=""
	"$IM" "$source" -resize "${SIZE}x${SIZE}^" -gravity center -extent "${SIZE}x${SIZE}" \
		-quality "$QUALITY" -strip "$target" || die "conversione fallita: $source"
	# IM_IDENTIFY è volutamente senza virgolette: su ImageMagick 7 sono due
	# parole («magick identify»), su ImageMagick 6 un binario solo.
	# shellcheck disable=SC2086
	dims="$($IM_IDENTIFY -format '%wx%h %m' "$target")"
	[ "$dims" = "${SIZE}x${SIZE} WEBP" ] || die "conversione inattesa: $target risulta $dims"
}

write_line() {
	local kind="$1" slug="$2"
	echo "fatto: ${COVERS_DIR}/${slug}.webp (${SIZE}x${SIZE} webp)" >&2
	if [ "$kind" = "cover" ]; then
		echo "aggiungi in album.md:  cover: \"${slug}.webp\"" >&2
	else
		echo "aggiungi in band.md:   image: \"${slug}.webp\"" >&2
	fi
}

# fetch_cover chiede la cover al Cover Art Archive e distingue il 404 (l'album
# non ha cover) da un errore di rete. Il Cover Art Archive risponde 307 verso
# archive.org, quindi il redirect va seguito (-L): senza, ogni cover sembra un
# errore.
fetch_cover() {
	local mbid="$1" target="$2" code=""
	echo "cover di un album dal Cover Art Archive ($mbid)..." >&2
	code="$(curl -sSL -o "${TMP_DIR}/cover" -w '%{http_code}' --max-time 120 -A "$USER_AGENT" \
		"${CAA_URL}/release-group/${mbid}/front-1200")" \
		|| die "errore di rete verso il Cover Art Archive ($mbid)"
	case "$code" in
	200) ;;
	404) missing_or_die "nessuna cover per il release group $mbid sul Cover Art Archive" ;;
	*) die "risposta inattesa dal Cover Art Archive per $mbid (HTTP $code)" ;;
	esac
	to_webp "${TMP_DIR}/cover" "$target"
}

# fanart_keys legge le chiavi dall'ambiente o dai file, senza mai stamparle.
# Ne basta una: la di progetto (api-key) o la personale (client-key).
fanart_keys() {
	FANART_PROJECT="${FANART_API_KEY:-}"
	FANART_PERSONAL="${FANART_CLIENT_KEY:-}"
	if [ -z "$FANART_PROJECT" ] && [ -f "$FANART_KEY_FILE" ]; then
		FANART_PROJECT="$(tr -d '[:space:]' < "$FANART_KEY_FILE")"
	fi
	if [ -z "$FANART_PERSONAL" ] && [ -f "$FANART_CLIENT_KEY_FILE" ]; then
		FANART_PERSONAL="$(tr -d '[:space:]' < "$FANART_CLIENT_KEY_FILE")"
	fi
	if [ -z "$FANART_PROJECT" ] && [ -z "$FANART_PERSONAL" ]; then
		die "manca la chiave fanart.tv: esporta FANART_API_KEY (o FANART_CLIENT_KEY), oppure scrivila in ${FANART_KEY_FILE}"
	fi
}

# artist_thumb legge la risposta e stampa nome dell'artista e URL della foto
# più apprezzata: due righe, nome per primo.
artist_thumb() {
	python3 - "$1" <<'PY'
import json, sys

data = json.load(open(sys.argv[1], encoding="utf-8"))
thumbs = data.get("artistthumb") or []
if not thumbs:
    sys.exit(1)
best = max(thumbs, key=lambda image: int(image.get("likes") or 0))
print(data.get("name", ""))
print(best.get("url", ""))
PY
}

# fetch_band: la chiave rifiutata (401/403) è un errore, l'artista o la foto
# assenti (404) sono "l'archivio non ce l'ha".
fetch_band() {
	local mbid="$1" target="$2" info="" name="" url="" code=""
	local -a headers=()
	fanart_keys
	# Le chiavi viaggiano in header: non finiscono né nell'URL né nei log.
	if [ -n "$FANART_PROJECT" ]; then headers+=(-H "api-key: ${FANART_PROJECT}"); fi
	if [ -n "$FANART_PERSONAL" ]; then headers+=(-H "client-key: ${FANART_PERSONAL}"); fi
	echo "foto della band da fanart.tv ($mbid)..." >&2
	code="$(curl -sS -o "${TMP_DIR}/artist.json" -w '%{http_code}' --max-time 60 "${headers[@]}" \
		"${FANART_API}/${mbid}")" || die "errore di rete verso fanart.tv ($mbid)"
	case "$code" in
	200) ;;
	404) missing_or_die "nessuna scheda per l'artista $mbid su fanart.tv" ;;
	401 | 403) die "chiave fanart.tv rifiutata (HTTP $code): controlla FANART_API_KEY / FANART_CLIENT_KEY" ;;
	*) die "risposta inattesa da fanart.tv per $mbid (HTTP $code)" ;;
	esac
	info="$(artist_thumb "${TMP_DIR}/artist.json")" \
		|| missing_or_die "nessuna foto (artistthumb) per l'artista $mbid su fanart.tv"
	name="${info%%$'\n'*}"
	url="${info#*$'\n'}"
	[ -n "$url" ] || die "foto senza URL per l'artista $mbid"
	echo "band trovata: ${name}" >&2
	curl -fsSL --max-time 120 -o "${TMP_DIR}/artist" "$url" \
		|| die "download della foto fallito: $url"
	to_webp "${TMP_DIR}/artist" "$target"
}

# parse_args separa i flag dagli argomenti posizionali.
parse_args() {
	local arg=""
	ARGS=()
	for arg in "$@"; do
		case "$arg" in
		--force) FORCE="yes" ;;
		--allow-missing) ALLOW_MISSING="yes" ;;
		-*) usage; exit 2 ;;
		*) ARGS+=("$arg") ;;
		esac
	done
}

main() {
	[ "$#" -ge 1 ] || { usage; exit 2; }
	local action="$1" mbid="" slug="" target=""
	shift
	case "$action" in
	-h | --help | help) usage; exit 0 ;;
	esac
	parse_args "$@"
	mbid="${ARGS[0]:-}"
	slug="${ARGS[1]:-}"
	if [ -z "$mbid" ] || [ -z "$slug" ]; then
		usage
		exit 2
	fi

	require_tools
	check_mbid "$mbid"
	check_slug "$slug"
	TMP_DIR="$(mktemp -d)"
	target="$(prepare_target "$slug" "$FORCE")"

	case "$action" in
	cover) fetch_cover "$mbid" "$target" ;;
	band) fetch_band "$mbid" "$target" ;;
	*)
		usage
		exit 2
		;;
	esac
	write_line "$action" "$slug"
}

main "$@"

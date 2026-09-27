#!/usr/bin/env bash
# fetch-covers.sh — scarica le immagini del sito dagli archivi pubblici.
#
#   bash scripts/fetch-covers.sh cover <release-group-mbid> <album-slug>
#   bash scripts/fetch-covers.sh band  <artist-mbid>        <band-slug>
#
# Le cover degli album arrivano dal Cover Art Archive (MusicBrainz + Internet
# Archive) e le foto delle band da fanart.tv, che le indicizza per MBID di
# MusicBrainz. Lo script scrive covers/<slug>.webp a 600x600, la sola
# dimensione che il sito usa (docs/CONTENT.md).
#
# Regole rispettate (docs/GUIDELINES.md, D76):
#   - nessun fallback silenzioso: se manca uno strumento, l'immagine o la
#     chiave, esce con un errore esplicito e non lascia file a metà;
#   - nessun file temporaneo lasciato in giro: si scarica in una cartella
#     temporanea e la si cancella sempre (trap).
#
# La chiave di fanart.tv NON sta nel repo: si passa in FANART_API_KEY oppure si
# scrive in ~/.fanart_api_key (una riga sola). Serve solo per le foto delle
# band: le cover non richiedono chiavi.
#
# Si invoca con `bash` esplicito, come scripts/fetch-assets.sh: nel repo gli
# script non hanno il bit di esecuzione.

set -euo pipefail

readonly SIZE=600
readonly QUALITY=82
readonly CAA_URL="https://coverartarchive.org"
readonly FANART_API="https://webservice.fanart.tv/v3.2/music"
readonly FANART_KEY_FILE="${HOME}/.fanart_api_key"
readonly COVERS_DIR="covers"
# MusicBrainz chiede un User-Agent che identifichi chi chiama.
readonly USER_AGENT="lyrica-covers/1.0 (+https://github.com/filippogrande/lyrica)"

IM=""
TMP_DIR=""

usage() {
	cat >&2 <<'FINE'
uso:
  bash scripts/fetch-covers.sh cover <release-group-mbid> <album-slug> [--force]
  bash scripts/fetch-covers.sh band  <artist-mbid>        <band-slug>  [--force]

esempi:
  bash scripts/fetch-covers.sh cover c31a5e2b-0bf8-32e0-8aeb-ef4ba9973932 liebe-ist-fuer-alle-da
  bash scripts/fetch-covers.sh band  cf075492-d880-4afc-b87b-d6b03e33dacc electric-callboy

Il MBID si legge nell'URL della pagina MusicBrainz:
  musicbrainz.org/release-group/<release-group-mbid>   -> cover di un album
  musicbrainz.org/artist/<artist-mbid>                 -> foto di una band

Scrive covers/<slug>.webp a 600x600. Non tocca i file di contenuto: la riga da
aggiungere al front-matter te la stampa alla fine.
FINE
}

die() {
	echo "errore: $*" >&2
	exit 1
}

cleanup() {
	[ -n "$TMP_DIR" ] && rm -rf "$TMP_DIR"
	return 0
}

trap cleanup EXIT

# require_tools verifica ciò che serve PRIMA di scaricare qualcosa.
require_tools() {
	command -v curl >/dev/null 2>&1 || die "manca curl"
	command -v python3 >/dev/null 2>&1 || die "manca python3 (serve per leggere la risposta di fanart.tv)"
	if command -v magick >/dev/null 2>&1; then
		IM="magick"
	elif command -v convert >/dev/null 2>&1; then
		IM="convert"
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
	dims="$("$IM" identify -format '%wx%h %m' "$target")"
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

fetch_cover() {
	local mbid="$1" target="$2"
	echo "cover di un album dal Cover Art Archive ($mbid)..." >&2
	curl -fsSL --max-time 120 -A "$USER_AGENT" -o "${TMP_DIR}/cover" \
		"${CAA_URL}/release-group/${mbid}/front-1200" \
		|| die "nessuna cover per il release group $mbid (404) o errore di rete"
	to_webp "${TMP_DIR}/cover" "$target"
}

# fanart_key legge la chiave dall'ambiente o dal file, senza mai stamparla.
fanart_key() {
	local key="${FANART_API_KEY:-}"
	if [ -z "$key" ] && [ -f "$FANART_KEY_FILE" ]; then
		key="$(tr -d '[:space:]' < "$FANART_KEY_FILE")"
	fi
	[ -n "$key" ] || die "manca la chiave fanart.tv: esporta FANART_API_KEY o scrivila in ${FANART_KEY_FILE}"
	printf '%s' "$key"
}

# best_artist_thumb sceglie la foto più apprezzata fra quelle disponibili.
best_artist_thumb() {
	python3 - "$1" <<'PY'
import json, sys

data = json.load(open(sys.argv[1], encoding="utf-8"))
thumbs = data.get("artistthumb") or []
if not thumbs:
    sys.exit(1)
best = max(thumbs, key=lambda image: int(image.get("likes") or 0))
print(best.get("url", ""))
PY
}

fetch_band() {
	local mbid="$1" target="$2" key="" url=""
	key="$(fanart_key)"
	echo "foto della band da fanart.tv ($mbid)..." >&2
	# La chiave viaggia in un header: non finisce né nell'URL né nei log.
	curl -fsSL --max-time 60 -H "api-key: ${key}" -o "${TMP_DIR}/artist.json" \
		"${FANART_API}/${mbid}" \
		|| die "nessuna risposta per l'artista $mbid: chiave non valida (401) o artista assente su fanart.tv (404)"
	url="$(best_artist_thumb "${TMP_DIR}/artist.json")" \
		|| die "nessuna foto (artistthumb) per l'artista $mbid su fanart.tv"
	[ -n "$url" ] || die "foto senza URL per l'artista $mbid"
	curl -fsSL --max-time 120 -o "${TMP_DIR}/artist" "$url" \
		|| die "download della foto fallito: $url"
	to_webp "${TMP_DIR}/artist" "$target"
}

main() {
	[ "$#" -ge 1 ] || { usage; exit 2; }
	local action="$1" mbid="" slug="" force="no" target=""
	shift
	case "$action" in
	-h | --help | help) usage; exit 0 ;;
	esac
	for arg in "$@"; do
		[ "$arg" = "--force" ] && force="yes"
	done
	mbid="${1:-}"
	slug="${2:-}"
	[ -n "$mbid" ] && [ "$slug" != "" ] || { usage; exit 2; }
	[ "$slug" != "--force" ] || { usage; exit 2; }

	require_tools
	check_mbid "$mbid"
	check_slug "$slug"
	TMP_DIR="$(mktemp -d)"
	target="$(prepare_target "$slug" "$force")"

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

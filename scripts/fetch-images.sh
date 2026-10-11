#!/usr/bin/env bash
# fetch-images.sh — riempie DA SOLO le immagini che mancano (D86).
#
#   bash scripts/fetch-images.sh              # tutte le band
#   bash scripts/fetch-images.sh rammstein    # una sola band
#
# Niente MBID da copiare a mano: la band si riconosce dal `name` in band.md,
# l'album dal `title` e dall'`year` in album.md, e l'identificativo corrispondente
# lo cerca MusicBrainz.
#   foto della band  -> artist MBID        -> fanart.tv            (campo `image`)
#   cover dell'album -> release-group MBID -> Cover Art Archive    (campo `cover`)
# Il campo viene scritto nel front-matter con scripts/declare-image.py.
#
# Salta quello che è già a posto (immagine dichiarata e file presente in
# covers/). "L'archivio non ha quell'immagine" NON è un errore: si salta e si va
# avanti. Un errore vero (rete, chiave rifiutata, conversione) fa uscire con 1,
# e il report dice quale immagine.
#
# Alla fine ricontrolla che ogni immagine dichiarata nei contenuti esista: è la
# regola che la CI applica al merge, e le PR create dal workflow non fanno partire
# la CI (GitHub non esegue workflow per gli eventi creati da GITHUB_TOKEN), quindi
# il controllo va fatto qui o non si fa.
#
# Stampa il report in Markdown; con GITHUB_STEP_SUMMARY impostato lo copia anche
# nel riepilogo del run, così si legge senza aprire i log. Ogni cosa viene
# stampata una volta sola: il log del run non deve contenere la tabella due volte.
#
# La corrispondenza è automatica e quindi non è infallibile (titolo uguale per
# un album live o un'edizione diversa): il report serve proprio a controllare
# prima di mergiare (D86).
#
# Si invoca con `bash` esplicito, come gli altri script del repo.

set -euo pipefail

readonly CONTENT_DIR="content/bands"
readonly COVERS_DIR="covers"
readonly FETCH_SCRIPT="scripts/fetch-covers.sh"
readonly DECLARE_SCRIPT="scripts/declare-image.py"
readonly MB_API="https://musicbrainz.org/ws/2"
readonly USER_AGENT="lyrica-covers/1.0 (+https://github.com/filippogrande/lyrica)"
# MusicBrainz chiede al massimo una richiesta al secondo.
readonly MB_DELAY=1
# Pausa fra un'immagine e l'altra: gli archivi non amano le raffiche, e in blocco
# si arriva facilmente al "troppe richieste".
readonly IMAGE_DELAY=2
# Retry per le chiamate MusicBrainz: il server risponde 503 quando è occupato.
readonly MB_RETRIES=3
readonly MB_RETRY_WAIT=5

TMP_DIR=""
ROWS=()
ERRORS=0

die() {
	echo "errore: $*" >&2
	exit 1
}

cleanup() {
	if [ -n "$TMP_DIR" ]; then rm -rf "$TMP_DIR"; fi
	return 0
}

trap cleanup EXIT

# row accumula una riga del report: tipo~cosa~esito.
row() {
	ROWS+=("$1~$2~$3")
}

# frontmatter_field legge un campo dal front-matter. Stampa niente se il campo
# non c'è: chi chiama decide se è un errore.
#
# La command substitution è chiusa SUBITO dopo la parentesi, senza spazi dentro
# le virgolette: uno spazio lì (") ") finirebbe dentro il valore, e un valore
# con uno spazio finale fa sembrare inesistente un file che esiste, oltre a
# sporcare le ricerche (artist:"Rammstein ").
frontmatter_field() {
	local file="$1" field="$2" body="" value=""
	body="$(awk 'NR==1 && $0 != "---" { exit } NR>1 && $0 == "---" { exit } NR>1 { print }' "$file")"
	value="$(printf '%s\n' "$body" | grep -m1 "^${field}:" \
		| sed -e "s/^${field}:[[:space:]]*//" -e 's/^"//' -e 's/"[[:space:]]*$//' -e 's/[[:space:]]*$//')" || true
	printf '%s' "$value"
	return 0
}

# http_code_with_retry scarica in un file e stampa il codice HTTP finale,
# ripetendo quando il server chiede di rallentare (429) o sbaglia lui (5xx).
# MusicBrainz risponde 503 quando è occupato: riprovare è la soluzione.
http_code_with_retry() {
	local out="$1" url="$2" tries="$MB_RETRIES" wait="$MB_RETRY_WAIT" code=""
	shift 2
	while :; do
		code="$(curl -sS -o "$out" -w '%{http_code}' "$@" "$url")" || code="000"
		case "$code" in
		429 | 5??)
			tries=$((tries - 1))
			if [ "$tries" -le 0 ]; then break; fi
			echo " MusicBrainz ha risposto $code: riprovo fra ${wait}s" >&2
			sleep "$wait"
			wait=$((wait * 2))
			;;
		*) break ;;
		esac
	done
	printf '%s' "$code"
}

# resolve_artist_mbid cerca l'artist dal nome esatto della band e stampa
# mbid<TAB>nome<TAB>tipo. Stampa niente se non trova una corrispondenza sicura:
# meglio saltare che scaricare la foto di un'altra band.
# Se MusicBrainz è occupato (503 dopo i retry), stampa niente e chiama decide.
resolve_artist_mbid() {
	local name="$1" json="${TMP_DIR}/artist.json" found=""
	local code
	code="$(http_code_with_retry "$json" "${MB_API}/artist" \
		--max-time 60 -A "$USER_AGENT" -G \
		--data-urlencode "query=artist:\"${name}\"" \
		--data-urlencode "fmt=json" --data-urlencode "limit=5")"
	sleep "$MB_DELAY"
	if [ "$code" != "200" ]; then
		printf '%s' ""
		return 0
	fi
	found="$(python3 - "$json" "$name" <<'PY'
import json, sys

data = json.load(open(sys.argv[1], encoding="utf-8"))
wanted = sys.argv[2].strip().casefold()
candidates = []
for artist in data.get("artists", []):
    score = int(artist.get("score") or 0)
    if score < 95:
        continue
    if (artist.get("name") or "").strip().casefold() != wanted:
        continue
    candidates.append((artist.get("type") == "Group", score, artist))
if not candidates:
    sys.exit(0)
best = max(candidates, key=lambda candidate: (candidate[0], candidate[1]))[2]
print("\t".join([best["id"], best.get("name", ""), best.get("type", "")]))
PY
	)" || true
	printf '%s' "$found"
	return 0
}

# resolve_release_group_mbid cerca il release-group dal titolo e dall'anno e
# stampa mbid<TAB>titolo<TAB>data. Solo album (non singoli, non live o
# compilation), punteggio alto e anno vicino a quello scritto in album.md.
# Se MusicBrainz è occupato (503 dopo i retry), stampa niente e chiama decide.
resolve_release_group_mbid() {
	local artist="$1" title="$2" year="$3" json="${TMP_DIR}/release-group.json" found=""
	local code
	code="$(http_code_with_retry "$json" "${MB_API}/release-group" \
		--max-time 60 -A "$USER_AGENT" -G \
		--data-urlencode "query=releasegroup:\"${title}\" AND artist:\"${artist}\"" \
		--data-urlencode "fmt=json" --data-urlencode "limit=25")"
	sleep "$MB_DELAY"
	if [ "$code" != "200" ]; then
		printf '%s' ""
		return 0
	fi
	found="$(python3 - "$json" "$title" "$year" <<'PY'
import json, sys

data = json.load(open(sys.argv[1], encoding="utf-8"))
wanted = sys.argv[2].strip().casefold()
wanted_year = int(sys.argv[3]) if sys.argv[3].isdigit() else 0


def year_of(group):
    date = (group.get("first-release-date") or "")[:4]
    return int(date) if date.isdigit() else 0


candidates = []
for group in data.get("release-groups", []):
    score = int(group.get("score") or 0)
    if score < 90 or group.get("primary-type") != "Album":
        continue
    if (group.get("title") or "").strip().casefold() != wanted:
        continue
    year = year_of(group)
    if wanted_year and year and abs(year - wanted_year) > 1:
        continue
    distance = abs(year - wanted_year) if (wanted_year and year) else 99
    candidates.append((score, -distance, group))
if not candidates:
    sys.exit(0)
best = max(candidates, key=lambda candidate: (candidate[0], candidate[1]))[2]
print("\t".join([best["id"], best.get("title", ""), best.get("first-release-date", "") or "?"]))
PY
	)" || true
	printf '%s' "$found"
	return 0
}

# declare scrive il campo nel front-matter; se fallisce è un errore vero.
declare() {
	local file="$1" kind="$2" slug="$3"
	python3 "$DECLARE_SCRIPT" "$file" "$kind" "${slug}.webp" >/dev/null \
		|| die "non riesco a scrivere il campo in ${file}"
}

do_band() {
	local band_dir="$1" band_file="${1}band.md" slug="" name="" declared="" found="" mbid="" rc=0
	slug="$(basename "$band_dir")"
	if [ ! -f "$band_file" ]; then
		row "foto" "$slug" "errore: manca band.md"
		ERRORS=$((ERRORS + 1))
		return 0
	fi
	name="$(frontmatter_field "$band_file" name)"
	if [ -z "$name" ]; then
		row "foto" "$slug" "errore: band.md senza name"
		ERRORS=$((ERRORS + 1))
		return 0
	fi
	declared="$(frontmatter_field "$band_file" image)"
	if [ -n "$declared" ] && [ -f "${COVERS_DIR}/${declared}" ]; then
		row "foto" "$name" "già presente"
		return 0
	fi
	if [ -f "${COVERS_DIR}/${slug}.webp" ]; then
		declare "$band_file" band "$slug"
		row "foto" "$name" "dichiarata (il file c'era già)"
		return 0
	fi
	found="$(resolve_artist_mbid "$name")"
	if [ -z "$found" ]; then
		row "foto" "$name" "nessuna corrispondenza su MusicBrainz"
		return 0
	fi
	mbid="${found%%$'\t'*}"
	set +e
	bash "$FETCH_SCRIPT" band "$mbid" "$slug" --allow-missing
	rc=$?
	set -e
	sleep "$IMAGE_DELAY"
	case "$rc" in
	0)
		declare "$band_file" band "$slug"
		row "foto" "$name" "scaricata (artist ${mbid})"
		;;
	3) row "foto" "$name" "nessuna foto su fanart.tv" ;;
	*)
		row "foto" "$name" "errore nel download o nella conversione"
		ERRORS=$((ERRORS + 1))
		;;
	esac
	return 0
}

do_album() {
	local band_name="$1" album_dir="$2" album_file="${2}album.md" slug="" title="" year="" declared="" found="" mbid="" date="" rc=0
	slug="$(basename "$album_dir")"
	if [ -z "$band_name" ]; then
		row "cover" "$slug" "errore: band.md senza name"
		ERRORS=$((ERRORS + 1))
		return 0
	fi
	if [ ! -f "$album_file" ]; then
		row "cover" "${band_name}/${slug}" "errore: manca album.md"
		ERRORS=$((ERRORS + 1))
		return 0
	fi
	title="$(frontmatter_field "$album_file" title)"
	year="$(frontmatter_field "$album_file" year)"
	if [ -z "$title" ]; then
		row "cover" "${band_name}/${slug}" "errore: album.md senza title"
		ERRORS=$((ERRORS + 1))
		return 0
	fi
	declared="$(frontmatter_field "$album_file" cover)"
	if [ -n "$declared" ] && [ -f "${COVERS_DIR}/${declared}" ]; then
		row "cover" "${band_name}/${title}" "già presente"
		return 0
	fi
	if [ -f "${COVERS_DIR}/${slug}.webp" ]; then
		declare "$album_file" cover "$slug"
		row "cover" "${band_name}/${title}" "dichiarata (il file c'era già)"
		return 0
	fi
	found="$(resolve_release_group_mbid "$band_name" "$title" "$year")"
	if [ -z "$found" ]; then
		row "cover" "${band_name}/${title}" "nessuna corrispondenza su MusicBrainz"
		return 0
	fi
	mbid="${found%%$'\t'*}"
	date="$(printf '%s' "$found" | cut -f3)"
	set +e
	bash "$FETCH_SCRIPT" cover "$mbid" "$slug" --allow-missing
	rc=$?
	set -e
	sleep "$IMAGE_DELAY"
	case "$rc" in
	0)
		declare "$album_file" cover "$slug"
		row "cover" "${band_name}/${title}" "scaricata (${date})"
		;;
	3) row "cover" "${band_name}/${title}" "nessuna cover sul Cover Art Archive" ;;
	*)
		row "cover" "${band_name}/${title}" "errore nel download o nella conversione"
		ERRORS=$((ERRORS + 1))
		;;
	esac
	return 0
}

# check_declared_images verifica che ogni immagine dichiarata nei contenuti
# esista in covers/: è la stessa cosa che la CI controlla al merge (regola 6 per
# le cover, regola 10 per le band), e va fatta qui perché questa PR non fa
# partire la CI.
check_declared_images() {
	local band_dir="" band_file="" band_slug="" album_dir="" album_file="" album_slug="" declared=""
	for band_dir in "$CONTENT_DIR"/*/; do
		[ -d "$band_dir" ] || continue
		band_slug="$(basename "$band_dir")"
		band_file="${band_dir}band.md"
		declared="$(frontmatter_field "$band_file" image)"
		if [ -n "$declared" ] && [ ! -f "${COVERS_DIR}/${declared}" ]; then
			row "foto" "$band_slug" "errore: dichiara ${declared}, che non esiste in ${COVERS_DIR}/"
			ERRORS=$((ERRORS + 1))
		fi
		for album_dir in "$band_dir"*/; do
			[ -d "$album_dir" ] || continue
			album_slug="$(basename "$album_dir")"
			album_file="${album_dir}album.md"
			declared="$(frontmatter_field "$album_file" cover)"
			if [ -n "$declared" ] && [ ! -f "${COVERS_DIR}/${declared}" ]; then
				row "cover" "${band_slug}/${album_slug}" "errore: dichiara ${declared}, che non esiste in ${COVERS_DIR}/"
				ERRORS=$((ERRORS + 1))
			fi
		done
	done
	return 0
}

report_table() {
	local entry="" kind="" what="" outcome=""
	echo "| Immagine | Band / album | Esito |"
	echo "|---|---|---|"
	for entry in "${ROWS[@]}"; do
		IFS='~' read -r kind what outcome <<<"$entry" || true
		echo "| ${kind} | ${what} | ${outcome} |"
	done
}

print_report() {
	local entry="" kind="" what="" outcome="" tot=0 scaricate=0 saltate=0 riepilogo=""
	for entry in "${ROWS[@]}"; do
		IFS='~' read -r kind what outcome <<<"$entry" || true
		tot=$((tot + 1))
		case "$outcome" in
		scaricata*) scaricate=$((scaricate + 1)) ;;
		*) saltate=$((saltate + 1)) ;;
		esac
	done
	riepilogo="${tot} immagini considerate, ${scaricate} scaricate, ${saltate} saltate, ${ERRORS} errori."
	echo ""
	report_table
	echo "${riepilogo}"
	if [ -n "${GITHUB_STEP_SUMMARY:-}" ]; then
		{
			echo "## Immagini"
			echo ""
			report_table
			echo ""
			echo "${riepilogo}"
		} >>"$GITHUB_STEP_SUMMARY"
	fi
	return 0
}

main() {
	local filter="${1:-}" band_dir="" band_name="" album_dir=""
	[ -d "$CONTENT_DIR" ] || die "manca ${CONTENT_DIR}: lancia lo script dalla root del repo"
	if [ -n "$filter" ] && [ ! -d "${CONTENT_DIR}/${filter}" ]; then
		die "band sconosciuta: ${CONTENT_DIR}/${filter} non esiste"
	fi
	for band_dir in "$CONTENT_DIR"/*/; do
		[ -d "$band_dir" ] || continue
		if [ -n "$filter" ] && [ "$(basename "$band_dir")" != "$filter" ]; then
			continue
		fi
		do_band "$band_dir"
		band_name="$(frontmatter_field "${band_dir}band.md" name)"
		for album_dir in "$band_dir"*/; do
			[ -d "$album_dir" ] || continue
			do_album "$band_name" "$album_dir"
		done
	done
	check_declared_images
	print_report
	if [ "$ERRORS" -ne 0 ]; then
		echo "errore: ${ERRORS} immagini non riuscite" >&2
		exit 1
	fi
}

TMP_DIR="$(mktemp -d)"
main "$@"

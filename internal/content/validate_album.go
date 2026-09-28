package content

import (
	"os"
	"path/filepath"
	"strings"
)

// validateAlbum controlla un album: anagrafica, tracklist e brani.
func validateAlbum(catalog *Catalog, album *Album, opts Options, report *Report) {
	where := album.Directory
	checkSlugAgainst(album.Slug, filepath.Base(album.Directory), "cartella", where, report)
	if strings.TrimSpace(album.Title) == "" {
		report.Errorf(where, "title mancante")
	}
	if album.Year == 0 {
		report.Warnf(where, "year mancante: l'album resta senza anno in pagina")
	}
	checkTracklist(album, report)
	checkCover(album, opts, report)
	for _, track := range album.Tracks {
		validateTrack(catalog, track, report)
	}
	checkPendingTracks(album, report)
}

// checkTracklist verifica la regola 8 nei due sensi: ogni voce della tracklist
// deve avere il suo file, a meno che dichiari uno status ("non ha ancora una
// pagina"), e ogni file deve stare nella tracklist. È la stessa regola che
// tiene in ordine l'album.
func checkTracklist(album *Album, report *Report) {
	files := make(map[string]*Track, len(album.Tracks))
	for _, track := range album.Tracks {
		files[track.Slug] = track
	}
	listed := make(map[string]TrackRef, len(album.Tracklist))
	for _, ref := range album.Tracklist {
		listed[ref.Slug] = ref
		checkTrackRef(album, ref, files, report)
	}
	for _, track := range album.Tracks {
		ref, ok := listed[track.Slug]
		if !ok {
			report.Errorf(track.Path, "brano assente dalla tracklist di album.md (regola 8)")
			continue
		}
		checkTrackRefMatches(ref, track, report)
	}
}

// checkTrackRef controlla una voce della tracklist: senza status il file deve
// esistere, con lo status il file NON deve esistere (sarebbe una voce che dice
// di non avere il testo mentre il testo c'è).
func checkTrackRef(album *Album, ref TrackRef, files map[string]*Track, report *Report) {
	_, hasFile := files[ref.Slug]
	if strings.TrimSpace(ref.Status) == "" {
		if !hasFile {
			report.Errorf(album.Directory, "la tracklist elenca %q ma non esiste tracks/%s.md: aggiungi il file, oppure dichiara status: %s (regola 8)", ref.Slug, ref.Slug, StatusPending)
		}
		return
	}
	switch ref.Status {
	case StatusPending, StatusInstrumental:
	default:
		report.Errorf(album.Directory, "voce %q con status %q sconosciuto: ammessi %q e %q", ref.Slug, ref.Status, StatusPending, StatusInstrumental)
	}
	if hasFile {
		report.Errorf(album.Directory, "voce %q dichiara status: %s ma tracks/%s.md esiste: il testo c'è, togli lo status", ref.Slug, ref.Status, ref.Slug)
	}
}

// checkTrackRefMatches confronta la voce della tracklist con il file del brano:
// sono due copie degli stessi dati e vanno tenute allineate.
func checkTrackRefMatches(ref TrackRef, track *Track, report *Report) {
	if ref.Title != track.Title {
		report.Warnf(track.Path, "title %q diverso da quello in tracklist (%q)", track.Title, ref.Title)
	}
	if ref.Instrumental != track.Instrumental {
		report.Warnf(track.Path, "instrumental=%t diverso da album.md (%t)", track.Instrumental, ref.Instrumental)
	}
}

// checkPendingTracks avvisa quando una voce dichiarata "in arrivo" non ha
// nemmeno una traccia con il testo: l'album sarebbe tutto da scrivere, non un
// lavoro in corso.
func checkPendingTracks(album *Album, report *Report) {
	entries := album.Entries()
	pending := 0
	published := 0
	for _, entry := range entries {
		if entry.Ref.Status == StatusPending {
			pending++
		}
		if entry.Published() {
			published++
		}
	}
	if len(entries) == 0 {
		report.Warnf(album.Directory, "tracklist vuota: la pagina album sarà senza brani")
		return
	}
	if pending == len(entries) {
		report.Warnf(album.Directory, "tutte le %d tracce sono \"in arrivo\": l'album non ha nemmeno un testo", len(entries))
	}
	if published == 0 && pending < len(entries) {
		report.Warnf(album.Directory, "nessun brano pubblicato: la pagina album non verrà generata finché non c'è una traduzione")
	}
}

// checkCover verifica la regola 6: se la cover è dichiarata deve esistere.
// Senza cover il layout va senza immagine: è un avviso, non un errore.
func checkCover(album *Album, opts Options, report *Report) {
	if strings.TrimSpace(album.Cover) == "" {
		report.Warnf(album.Directory, "album senza cover: il layout va senza immagine, nessun placeholder")
		return
	}
	path := filepath.Join(opts.CoversDir, album.Cover)
	info, err := os.Stat(path)
	if err != nil {
		report.Errorf(album.Directory, "cover %q non trovata in %s (regola 6)", album.Cover, opts.CoversDir)
		return
	}
	if info.IsDir() {
		report.Errorf(album.Directory, "cover %q è una cartella, non un file (regola 6)", album.Cover)
	}
}

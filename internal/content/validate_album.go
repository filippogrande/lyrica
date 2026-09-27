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
}

// checkTracklist verifica la regola 8 nei due sensi: ogni voce della tracklist
// deve avere il suo file, e ogni file deve stare nella tracklist. È la stessa
// regola che tiene in ordine l'album.
func checkTracklist(album *Album, report *Report) {
	files := make(map[string]*Track, len(album.Tracks))
	for _, track := range album.Tracks {
		files[track.Slug] = track
	}
	listed := make(map[string]TrackRef, len(album.Tracklist))
	for _, ref := range album.Tracklist {
		listed[ref.Slug] = ref
		if _, ok := files[ref.Slug]; !ok {
			report.Errorf(album.Directory, "la tracklist elenca %q ma non esiste tracks/%s.md (regola 8)", ref.Slug, ref.Slug)
		}
	}
	for _, track := range album.Tracks {
		ref, ok := listed[track.Slug]
		if !ok {
			report.Errorf(track.Path, "brano assente dalla tracklist di album.md (regola 8)")
			continue
		}
		if ref.Title != track.Title {
			report.Warnf(track.Path, "title %q diverso da quello in tracklist (%q)", track.Title, ref.Title)
		}
		if ref.Instrumental != track.Instrumental {
			report.Warnf(track.Path, "instrumental=%t diverso da album.md (%t)", track.Instrumental, ref.Instrumental)
		}
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

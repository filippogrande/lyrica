package content

import (
	"os"
	"path/filepath"
	"strings"
)

// validateBand controlla una band e scende nei suoi album.
func validateBand(catalog *Catalog, band *Band, opts Options, report *Report) {
	where := band.Directory
	checkSlugAgainst(band.Slug, filepath.Base(band.Directory), "cartella", where, report)
	if strings.TrimSpace(band.Name) == "" {
		report.Errorf(where, "name mancante: è il nome che si legge in pagina")
	}
	if len(band.OriginalLangs) == 0 {
		report.Errorf(where, "original_langs mancante: serve almeno una lingua originale")
	}
	if strings.TrimSpace(band.Description) == "" {
		report.Warnf(where, "descrizione band mancante: la pagina band resta senza testo")
	}
	if len(band.Albums) == 0 {
		report.Warnf(where, "band senza album pubblicati: non compare nell'elenco finché non ne ha uno con traduzioni")
	}
	checkBandImage(band, opts, report)
	checkDuplicateAlbums(band, report)
	for _, album := range band.Albums {
		validateAlbum(catalog, album, opts, report)
	}
	checkLinks(catalog, where, band.Description, report)
}

// checkBandImage verifica la regola 10: se la foto della band è dichiarata deve
// esistere in covers/. Senza foto la card usa il segnaposto con le iniziali:
// è un avviso, non un errore (una band può restare tipografica).
func checkBandImage(band *Band, opts Options, report *Report) {
	if strings.TrimSpace(band.Image) == "" {
		report.Warnf(band.Directory, "band senza immagine: la card usa il segnaposto con le iniziali")
		return
	}
	path := filepath.Join(opts.CoversDir, band.Image)
	info, err := os.Stat(path)
	if err != nil {
		report.Errorf(band.Directory, "immagine %q non trovata in %s (regola 10)", band.Image, opts.CoversDir)
		return
	}
	if info.IsDir() {
		report.Errorf(band.Directory, "immagine %q è una cartella, non un file (regola 10)", band.Image)
	}
}

// checkDuplicateAlbums verifica la regola 2 fra album della stessa band.
func checkDuplicateAlbums(band *Band, report *Report) {
	seen := make(map[string]string, len(band.Albums))
	for _, album := range band.Albums {
		if previous, ok := seen[album.Slug]; ok {
			report.Errorf(album.Directory, "slug %q già usato dall'album in %s (regola 2)", album.Slug, previous)
			continue
		}
		seen[album.Slug] = album.Directory
	}
}

// checkSlugAgainst verifica la regola 1 (slug = nome di cartella o file) e il
// formato dello slug.
func checkSlugAgainst(slug, expected, kind, where string, report *Report) {
	if slug != expected {
		report.Errorf(where, "slug %q diverso dal nome della %s %q (regola 1)", slug, kind, expected)
	}
	if !validSlug(slug) {
		report.Errorf(where, "slug %q non valido: minuscolo, ASCII, parole separate da trattini", slug)
	}
}

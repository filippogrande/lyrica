package content

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Load legge la cartella dei contenuti (content/) e costruisce il catalogo.
// Restituisce il primo errore trovato: se un contenuto non si legge, il
// catalogo non viene costruito a metà.
func Load(root string) (*Catalog, error) {
	bandsDir := filepath.Join(root, "bands")
	entries, err := os.ReadDir(bandsDir)
	if err != nil {
		return nil, fmt.Errorf("lettura di %s: %w", bandsDir, err)
	}
	catalog := &Catalog{}
	for _, entry := range entries {
		if isIgnorable(entry.Name()) {
			continue
		}
		if !entry.IsDir() {
			return nil, fmt.Errorf("%s: attese solo cartelle di band, trovato il file %q", bandsDir, entry.Name())
		}
		band, err := loadBand(filepath.Join(bandsDir, entry.Name()))
		if err != nil {
			return nil, err
		}
		catalog.Bands = append(catalog.Bands, band)
	}
	sort.SliceStable(catalog.Bands, func(i, j int) bool {
		return catalog.Bands[i].Name < catalog.Bands[j].Name
	})
	return catalog, nil
}

// isIgnorable riconosce i file di servizio del sistema operativo (.DS_Store,
// file di editor) che non sono contenuti e non devono far fallire il
// caricamento.
func isIgnorable(name string) bool {
	return strings.HasPrefix(name, ".")
}

// loadBand legge una band e i suoi album.
func loadBand(dir string) (*Band, error) {
	band := &Band{Directory: dir}
	if _, err := readFrontMatterFile(filepath.Join(dir, "band.md"), band); err != nil {
		return nil, err
	}
	albums, err := loadAlbums(dir, band)
	if err != nil {
		return nil, err
	}
	band.Albums = albums
	return band, nil
}

// loadAlbums legge le sottocartelle di una band: ognuna è un album. I file
// nella cartella della band (band.md) non sono album e si saltano.
func loadAlbums(bandDir string, band *Band) ([]*Album, error) {
	entries, err := os.ReadDir(bandDir)
	if err != nil {
		return nil, fmt.Errorf("lettura di %s: %w", bandDir, err)
	}
	var albums []*Album
	for _, entry := range entries {
		if isIgnorable(entry.Name()) || !entry.IsDir() {
			continue
		}
		album, err := loadAlbum(filepath.Join(bandDir, entry.Name()), band)
		if err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}
	return albums, nil
}

// loadAlbum legge album.md e i suoi brani.
func loadAlbum(dir string, band *Band) (*Album, error) {
	album := &Album{Directory: dir, Band: band}
	if _, err := readFrontMatterFile(filepath.Join(dir, "album.md"), album); err != nil {
		return nil, err
	}
	tracks, err := loadTracks(filepath.Join(dir, "tracks"), album)
	if err != nil {
		return nil, err
	}
	album.Tracks = orderByTracklist(album, tracks)
	return album, nil
}

// loadTracks legge i brani di un album. La struttura attuale è una cartella per
// brano (tracks/<slug>/); un file <slug>.md nella stessa cartella è il formato
// precedente, che `lyrica migrate` trasforma e che il loader accetta finché
// l'ultimo brano non è migrato.
func loadTracks(dir string, album *Album) ([]*Track, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("lettura di %s: %w", dir, err)
	}
	var tracks []*Track
	for _, entry := range entries {
		if isIgnorable(entry.Name()) {
			continue
		}
		track, err := loadTrackEntry(filepath.Join(dir, entry.Name()), entry, album)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}
	return tracks, nil
}

// loadTrackEntry legge un brano da una voce della cartella tracks/: una cartella
// (formato attuale) o un file .md (formato precedente). Qualsiasi altra cosa è
// un errore esplicito.
func loadTrackEntry(path string, entry os.DirEntry, album *Album) (*Track, error) {
	switch {
	case entry.IsDir():
		return loadTrackFolder(path, album)
	case strings.HasSuffix(entry.Name(), ".md"):
		return loadLegacyTrack(path, album)
	default:
		return nil, fmt.Errorf("%s: attese cartelle di brani o file .md, trovato %q", filepath.Dir(path), entry.Name())
	}
}

// loadLegacyTrack legge un brano nel formato precedente (tracks/<slug>.md):
// front-matter nei modelli, blocchi inline, corpo in Notes. Il brano si segna
// come Legacy, così il validatore lo riconosce e la migrazione sa che c'è.
func loadLegacyTrack(path string, album *Album) (*Track, error) {
	track := &Track{Path: path, Album: album, Legacy: true}
	notes, err := readFrontMatterFile(path, track)
	if err != nil {
		return nil, err
	}
	track.Notes = notes
	return track, nil
}

// orderByTracklist ordina i brani come dichiarato in album.md: l'ordine
// stampato è quello della tracklist. I brani non elencati vanno in fondo e li
// segnala il validatore (album.md e tracks/ devono coincidere).
func orderByTracklist(album *Album, tracks []*Track) []*Track {
	positions := make(map[string]int, len(album.Tracklist))
	for i, ref := range album.Tracklist {
		positions[ref.Slug] = i
	}
	last := len(positions) + 1
	ordered := append([]*Track(nil), tracks...)
	sort.SliceStable(ordered, func(i, j int) bool {
		pi, ok := positions[ordered[i].Slug]
		if !ok {
			pi = last
		}
		pj, ok := positions[ordered[j].Slug]
		if !ok {
			pj = last
		}
		return pi < pj
	})
	return ordered
}

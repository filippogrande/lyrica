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
// editor) che non sono contenuti e non devono far fallire il caricamento.
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

// loadAlbums legge le sottocartelle di una band: ognuna è un album.
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

// loadTracks legge i file dei brani di un album.
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
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			return nil, fmt.Errorf("%s: attesi solo file .md, trovato %q", dir, entry.Name())
		}
		tracks = append(tracks, loadTrackFile(filepath.Join(dir, entry.Name()), album))
	}
	return tracks, nil
}

// loadTrackFile legge un singolo brano.
func loadTrackFile(path string, album *Album) *Track {
	track := &Track{Path: path, Album: album}
	return track // sostituito dal file successivo: vedi loadTrack
}

// orderByTracklist ordina i brani come dichiarato in album.md: l'ordine
// stampato è quello della tracklist. I brani non elencati vanno in fondo e li
// segnala il validatore (album.md e tracks/ devono coincidere).
func orderByTracklist(album *Album, tracks []*Track) []*Track {
	positions := make(map[string]int, len(album.Tracklist))
	for i, ref := range album.Tracklist {
		positions[ref.Slug] = i
	}
	ordered := append([]*Track(nil), tracks...)
	sort.SliceStable(ordered, func(i, j int) bool {
		pi, oki := positions[ordered[i].Slug]
		pj, okj := positions[ordered[j].Slug]
		if !oki {
			pi = len(positions) + 1
		}
		if !okj {
			pj = len(positions) + 1
		}
		return pi < pj
	})
	return ordered
}

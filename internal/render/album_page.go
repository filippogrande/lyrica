package render

import (
	"sort"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/filippogrande/lyrica/internal/content"
)

// BuildAlbumView assembla la pagina di un album: copertina e anagrafica (band,
// anno, numero di brani, lingue tradotte), tracklist completa e numerata,
// corsia degli altri album della band. I costruttori stanno qui e non in view.go
// per non far crescere oltre il limite il file che raccoglie le viste di tutte
// le pagine.
func BuildAlbumView(page PageData, band *content.Band, album *content.Album) AlbumView {
	rows := make([]TrackRow, 0, len(album.Tracks))
	for _, track := range album.Tracks {
		rows = append(rows, newTrackRow(page, band, album, track))
	}
	view := AlbumView{
		Page:        page,
		Crumbs:      bandCrumbs(page, band, album),
		Title:       album.Title,
		YearText:    yearText(album),
		BandName:    band.Name,
		BandURL:     templ.URL(page.BandPath(band.Slug)),
		CoverAlt:    page.T("album.cover_alt"),
		TracksText:  strconv.Itoa(len(album.Tracks)),
		LangsText:   albumLanguages(album),
		Tracklist:   rows,
		OtherAlbums: newAlbumRail(page, page.T("album.other_albums"), band, otherAlbums(band, album)),
	}
	if album.Cover != "" {
		view.ShowCover = true
		view.CoverURL = templ.URL("/covers/" + album.Cover)
	}
	return view
}

// albumLanguages elenca le lingue delle traduzioni presenti nell'album, senza
// ripetizioni e in ordine alfabetico: è la risposta a "in che lingue posso
// leggere questo album". Vuoto se nell'album non c'è ancora nessuna traduzione.
func albumLanguages(album *content.Album) string {
	seen := make(map[string]bool)
	var langs []string
	for _, track := range album.Tracks {
		for _, lang := range track.TranslationLangs() {
			if !seen[lang] {
				seen[lang] = true
				langs = append(langs, lang)
			}
		}
	}
	sort.Strings(langs)
	return strings.Join(langs, " · ")
}

// otherAlbums elenca gli album pubblicati della band tranne quello corrente:
// la corsia in fondo alla pagina album non rimanda alla pagina che si sta
// guardando.
func otherAlbums(band *content.Band, album *content.Album) []*content.Album {
	var others []*content.Album
	for _, candidate := range publishedAlbums(band) {
		if candidate.Slug != album.Slug {
			others = append(others, candidate)
		}
	}
	return others
}

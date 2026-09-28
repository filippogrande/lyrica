package render

import (
	"strconv"

	"github.com/a-h/templ"
	"github.com/filippogrande/lyrica/internal/content"
)

// BuildAlbumView assembla la pagina di un album: copertina e anagrafica (band,
// anno, numero di tracce, lingue tradotte), tracklist completa e numerata,
// corsia degli altri album della band. I costruttori stanno qui e non in view.go
// per non far crescere oltre il limite il file che raccoglie le viste di tutte
// le pagine.
func BuildAlbumView(page PageData, band *content.Band, album *content.Album) AlbumView {
	entries := album.Entries()
	view := AlbumView{
		Page:        page,
		Crumbs:      bandCrumbs(page, band, album),
		Title:       album.Title,
		YearText:    yearText(album),
		BandName:    band.Name,
		BandURL:     templ.URL(page.BandPath(band.Slug)),
		CoverAlt:    page.T("album.cover_alt"),
		TracksText:  strconv.Itoa(len(entries)),
		LangsText:   albumLanguages(album),
		Tracklist:   newTrackRows(page, band, album, entries),
		OtherAlbums: newAlbumRail(page, page.T("album.other_albums"), band, otherAlbums(band, album)),
	}
	if album.Cover != "" {
		view.ShowCover = true
		view.CoverURL = templ.URL("/covers/" + album.Cover)
	}
	return view
}

// newTrackRows prepara le righe della tracklist: una per traccia dichiarata in
// album.md, nell'ordine del disco. La lista è completa: i brani senza testo ci
// sono comunque, senza link.
func newTrackRows(page PageData, band *content.Band, album *content.Album, entries []content.TracklistEntry) []TrackRow {
	rows := make([]TrackRow, 0, len(entries))
	for _, entry := range entries {
		rows = append(rows, newTrackRow(page, band, album, entry))
	}
	return rows
}

// newTrackRow prepara una riga di tracklist: i brani pubblicati sono link (con
// le lingue disponibili accanto, bandiera compresa), ogni altro resta in elenco
// senza link con il motivo per cui non ha una pagina.
func newTrackRow(page PageData, band *content.Band, album *content.Album, entry content.TracklistEntry) TrackRow {
	if !entry.Published() {
		return TrackRow{Title: entry.Title(), ReasonKey: reasonKey(entry)}
	}
	return TrackRow{
		Title:     entry.Title(),
		URL:       templ.URL(page.TrackPath(band.Slug, album.Slug, entry.Slug())),
		Linked:    true,
		LangsText: LangLabels(entry.Track.TranslationLangs()),
	}
}

// reasonKey è la chiave di locale del motivo per cui un brano non ha una
// pagina: il testo non è ancora scritto ("in arrivo"), è strumentale, oppure
// c'è l'originale ma non una traduzione. Il testo del motivo sta nei locale:
// nessuna stringa scritta nel codice.
func reasonKey(entry content.TracklistEntry) string {
	switch {
	case entry.Ref.Status == content.StatusPending:
		return "track.pending"
	case entry.Instrumental():
		return "track.instrumental"
	default:
		return "track.only_original"
	}
}

// albumLanguages elenca le lingue delle traduzioni presenti nell'album, come
// etichette con bandiera ("🇩🇪 DE · 🇮🇹 IT"): è la risposta a "in che lingue
// posso leggere questo album". Vuoto se nell'album non c'è ancora nessuna
// traduzione.
func albumLanguages(album *content.Album) string {
	var langs []string
	for _, track := range album.Tracks {
		langs = append(langs, track.TranslationLangs()...)
	}
	return LangLabels(langs)
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

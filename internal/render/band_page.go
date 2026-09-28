package render

import (
	"sort"
	"strings"

	"github.com/a-h/templ"
	"github.com/filippogrande/lyrica/internal/content"
)

// BuildBandView assembla la pagina di una band: anagrafica, corsia degli album
// pubblicati (dal più recente per anno) e corsia dei brani tradotti della band.
// I costruttori stanno qui e non in view.go per non far crescere oltre il limite
// il file che raccoglie le viste di tutte le pagine.
func BuildBandView(page PageData, band *content.Band) BandView {
	view := BandView{
		Page:        page,
		Crumbs:      bandCrumbs(page, band, nil),
		Name:        band.Name,
		Country:     band.Country,
		FormedText:  formedText(band),
		MembersText: strings.Join(band.Members, ", "),
		TagsText:    strings.Join(band.Tags, " · "),
		Description: band.Description,
		AlbumRail:   newAlbumRail(page, band, publishedAlbums(band)),
		TrackRail:   newTrackRail(page, page.T("home.recent"), bandTrackEntries(band)),
	}
	if band.Image != "" {
		view.HasImage = true
		view.ImageURL = templ.URL("/covers/" + band.Image)
	}
	return view
}

// newAlbumRail costruisce la corsia degli album pubblicati, dal più recente per
// anno di pubblicazione; a parità di anno, ordine alfabetico. Un album senza
// anno dichiarato (Year 0) finisce in coda, non in testa.
func newAlbumRail(page PageData, band *content.Band, albums []*content.Album) Rail {
	sorted := make([]*content.Album, len(albums))
	copy(sorted, albums)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].Year != sorted[j].Year {
			return sorted[i].Year > sorted[j].Year
		}
		return sorted[i].Title < sorted[j].Title
	})
	rail := Rail{Title: page.T("band.latest_albums")}
	for _, album := range sorted {
		if len(rail.Cards) == railLimit {
			break
		}
		rail.Cards = append(rail.Cards, newAlbumRailCard(page, band, album))
	}
	return rail
}

// newAlbumRailCard prepara la card di un album: immagine = copertina (o
// segnaposto con l'iniziale), seconda riga = anno di pubblicazione. L'URL è la
// pagina dell'album, quindi la corsia resta navigabile come quella dei brani.
func newAlbumRailCard(page PageData, band *content.Band, album *content.Album) RailCard {
	card := RailCard{
		URL:     templ.URL(page.AlbumPath(band.Slug, album.Slug)),
		Title:   album.Title,
		Meta:    yearText(album),
		Initial: initialOf(album.Title),
	}
	if album.Cover != "" {
		card.HasImage = true
		card.ImageURL = templ.URL("/covers/" + album.Cover)
	}
	return card
}

// bandTrackEntries raccoglie i brani pubblicati di una band (quelli con almeno
// una traduzione), dal più recente per data di aggiunta. È la corsia dei brani
// della pagina band: stessa forma delle voci della home, filtro sulla band.
func bandTrackEntries(band *content.Band) []trackEntry {
	var entries []trackEntry
	for _, album := range band.Albums {
		for _, track := range album.Tracks {
			if track.HasTranslations() {
				entries = append(entries, trackEntry{band: band, album: album, track: track})
			}
		}
	}
	sortTrackEntries(entries)
	return entries
}

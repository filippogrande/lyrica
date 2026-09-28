package render

import (
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
		AlbumRail:   newAlbumRail(page, page.T("band.latest_albums"), band, publishedAlbums(band)),
		TrackRail:   newTrackRail(page, page.T("home.recent"), bandTrackEntries(band)),
	}
	if band.Image != "" {
		view.HasImage = true
		view.ImageURL = templ.URL("/covers/" + band.Image)
	}
	return view
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

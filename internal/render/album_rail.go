package render

import (
	"sort"

	"github.com/a-h/templ"
	"github.com/filippogrande/lyrica/internal/content"
)

// newAlbumRail costruisce una corsia di album, dal più recente per anno di
// pubblicazione; a parità di anno, ordine alfabetico. Un album senza anno
// dichiarato (Year 0) finisce in coda, non in testa. Il titolo lo passa chi
// chiama: la stessa corsia serve la pagina band ("ultimi album pubblicati") e la
// pagina album ("altri album della band").
func newAlbumRail(page PageData, title string, band *content.Band, albums []*content.Album) Rail {
	sorted := make([]*content.Album, len(albums))
	copy(sorted, albums)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].Year != sorted[j].Year {
			return sorted[i].Year > sorted[j].Year
		}
		return sorted[i].Title < sorted[j].Title
	})
	rail := Rail{Title: title}
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

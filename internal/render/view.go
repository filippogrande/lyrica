package render

import (
	"sort"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/filippogrande/lyrica/internal/content"
)

// Crumb è una voce del breadcrumb.
type Crumb struct {
	Label   string
	URL     templ.SafeURL
	Current bool
}

// TrackCard è un brano in un elenco (home, sezioni "in evidenza"/"recenti").
type TrackCard struct {
	Title     string
	URL       templ.SafeURL
	MetaText  string
	DateText  string
	LangsText string

	// added serve solo a ordinare i recenti: non si stampa.
	added content.Date
}

// BandCard è una band in un elenco.
type BandCard struct {
	Name        string
	URL         templ.SafeURL
	Country     string
	TagsText    string
	Description string
	AlbumsText  string
}

// AlbumCard è un album elencato nella pagina di una band.
type AlbumCard struct {
	Title      string
	URL        templ.SafeURL
	YearText   string
	TracksText string
}

// TrackRow è una riga di tracklist: linkata se il brano è pubblicato, altrimenti
// solo dichiarata (strumentale o ancora senza traduzione).
type TrackRow struct {
	Title        string
	URL          templ.SafeURL
	Linked       bool
	Instrumental bool
	LangsText    string
}

// StanzaView è una strofa pronta da stampare.
type StanzaView struct {
	Singer string
	Lines  []string
}

// BlockView è un blocco di testo: l'originale o una traduzione.
type BlockView struct {
	Lang       string
	Translator string
	Live       bool
	Stanzas    []StanzaView
}

// HomeView sono i dati della home.
type HomeView struct {
	Page     PageData
	Featured []TrackCard
	Recent   []TrackCard
	Bands    []BandCard
}

// HasContent dice se il sito ha qualcosa da mostrare in home.
func (v HomeView) HasContent() bool { return len(v.Recent) > 0 }

// BandsView sono i dati della pagina con l'elenco delle band.
type BandsView struct {
	Page  PageData
	Bands []BandCard
}

// BandView sono i dati della pagina di una band.
type BandView struct {
	Page        PageData
	Crumbs      []Crumb
	Name        string
	Country     string
	FormedText  string
	MembersText string
	TagsText    string
	Description string
	Albums      []AlbumCard
}

// AlbumView sono i dati della pagina di un album.
type AlbumView struct {
	Page      PageData
	Crumbs    []Crumb
	Title     string
	YearText  string
	BandName  string
	BandURL   templ.SafeURL
	ShowCover bool
	CoverURL  templ.SafeURL
	CoverAlt  string
	Tracklist []TrackRow
}

// TrackView sono i dati della pagina di un brano: originale e traduzione.
type TrackView struct {
	Page           PageData
	Crumbs         []Crumb
	Title          string
	BandName       string
	BandURL        templ.SafeURL
	AlbumTitle     string
	AlbumURL       templ.SafeURL
	DateText       string
	SingersText    string
	Original       BlockView
	Translation    BlockView
	HasTranslation bool
}

// BuildHomeView assembla la home: in evidenza, recenti ed elenco delle band.
func BuildHomeView(page PageData, catalog *content.Catalog) HomeView {
	view := HomeView{Page: page, Bands: BuildBandCards(page, catalog.Bands)}
	forEachPublishedTrack(catalog, func(band *content.Band, album *content.Album, track *content.Track) {
		card := newTrackCard(page, band, album, track)
		view.Recent = append(view.Recent, card)
		if track.Featured {
			view.Featured = append(view.Featured, card)
		}
	})
	sortCardsByDate(view.Recent)
	sortCardsByDate(view.Featured)
	return view
}

// BuildBandCards costruisce le schede delle band che compaiono nel sito: solo
// quelle con almeno un album pubblicato (una band senza album pubblicati non
// compare in nessun elenco).
func BuildBandCards(page PageData, bands []*content.Band) []BandCard {
	var cards []BandCard
	for _, band := range bands {
		if len(publishedAlbums(band)) == 0 {
			continue
		}
		cards = append(cards, newBandCard(page, band))
	}
	return cards
}

// BuildBandView assembla la pagina di una band.
func BuildBandView(page PageData, band *content.Band) BandView {
	albums := publishedAlbums(band)
	cards := make([]AlbumCard, 0, len(albums))
	for _, album := range albums {
		cards = append(cards, newAlbumCard(page, band, album))
	}
	return BandView{
		Page:        page,
		Crumbs:      bandCrumbs(page, band, nil),
		Name:        band.Name,
		Country:     band.Country,
		FormedText:  formedText(band),
		MembersText: strings.Join(band.Members, ", "),
		TagsText:    strings.Join(band.Tags, " · "),
		Description: band.Description,
		Albums:      cards,
	}
}

// BuildAlbumView assembla la pagina di un album con la sua tracklist.
func BuildAlbumView(page PageData, band *content.Band, album *content.Album) AlbumView {
	rows := make([]TrackRow, 0, len(album.Tracks))
	for _, track := range album.Tracks {
		rows = append(rows, newTrackRow(page, band, album, track))
	}
	view := AlbumView{
		Page:      page,
		Crumbs:    bandCrumbs(page, band, album),
		Title:     album.Title,
		YearText:  yearText(album),
		BandName:  band.Name,
		BandURL:   templ.URL(page.BandPath(band.Slug)),
		CoverAlt:  page.T("album.cover_alt"),
		Tracklist: rows,
	}
	if album.Cover != "" {
		view.ShowCover = true
		view.CoverURL = templ.URL("/covers/" + album.Cover)
	}
	return view
}

// BuildTrackView assembla la pagina di un brano: originale e traduzione a
// fronte. La traduzione mostrata è quella nella lingua della pagina; se il
// brano non ce l'ha, si mostra la prima disponibile (la lingua del blocco è
// stampata in pagina, quindi non è un fallback nascosto).
func BuildTrackView(page PageData, band *content.Band, album *content.Album, track *content.Track) TrackView {
	original, _ := track.Original()
	translation, found := track.Translation(page.Lang)
	if !found {
		translation, found = firstTranslation(track)
	}
	view := TrackView{
		Page:           page,
		Crumbs:         trackCrumbs(page, band, album, track),
		Title:          track.Title,
		BandName:       band.Name,
		BandURL:        templ.URL(page.BandPath(band.Slug)),
		AlbumTitle:     album.Title,
		AlbumURL:       templ.URL(page.AlbumPath(band.Slug, album.Slug)),
		DateText:       track.AddedDate.String(),
		SingersText:    strings.Join(track.Singers, ", "),
		Original:       newBlockView(original),
		HasTranslation: found,
	}
	if found {
		view.Translation = newBlockView(translation)
	}
	return view
}

// forEachPublishedTrack visita i brani pubblicati: quelli con almeno una
// traduzione. Un brano senza traduzioni resta in tracklist ma non ha pagina.
func forEachPublishedTrack(catalog *content.Catalog, visit func(*content.Band, *content.Album, *content.Track)) {
	for _, band := range catalog.Bands {
		for _, album := range band.Albums {
			for _, track := range album.Tracks {
				if track.HasTranslations() {
					visit(band, album, track)
				}
			}
		}
	}
}

// publishedAlbums elenca gli album pubblicati: quelli con almeno un brano
// tradotto.
func publishedAlbums(band *content.Band) []*content.Album {
	var albums []*content.Album
	for _, album := range band.Albums {
		for _, track := range album.Tracks {
			if track.HasTranslations() {
				albums = append(albums, album)
				break
			}
		}
	}
	return albums
}

// newTrackCard prepara un brano per un elenco.
func newTrackCard(page PageData, band *content.Band, album *content.Album, track *content.Track) TrackCard {
	return TrackCard{
		Title:     track.Title,
		URL:       templ.URL(page.TrackPath(band.Slug, album.Slug, track.Slug)),
		MetaText:  band.Name + " — " + album.Title,
		DateText:  track.AddedDate.String(),
		LangsText: strings.Join(track.TranslationLangs(), " · "),
		added:     track.AddedDate,
	}
}

// newTrackRow prepara una riga di tracklist.
func newTrackRow(page PageData, band *content.Band, album *content.Album, track *content.Track) TrackRow {
	row := TrackRow{
		Title:        track.Title,
		Instrumental: track.Instrumental,
		LangsText:    strings.Join(track.TranslationLangs(), " · "),
	}
	if track.HasTranslations() {
		row.Linked = true
		row.URL = templ.URL(page.TrackPath(band.Slug, album.Slug, track.Slug))
	}
	return row
}

// newBandCard prepara una band per un elenco.
func newBandCard(page PageData, band *content.Band) BandCard {
	return BandCard{
		Name:        band.Name,
		URL:         templ.URL(page.BandPath(band.Slug)),
		Country:     band.Country,
		TagsText:    strings.Join(band.Tags, " · "),
		Description: band.Description,
		AlbumsText:  strconv.Itoa(len(publishedAlbums(band))),
	}
}

// newAlbumCard prepara un album per un elenco.
func newAlbumCard(page PageData, band *content.Band, album *content.Album) AlbumCard {
	return AlbumCard{
		Title:      album.Title,
		URL:        templ.URL(page.AlbumPath(band.Slug, album.Slug)),
		YearText:   yearText(album),
		TracksText: strconv.Itoa(len(album.Tracks)),
	}
}

// sortCardsByDate ordina i brani dal più recente.
func sortCardsByDate(cards []TrackCard) {
	sort.SliceStable(cards, func(i, j int) bool {
		return cards[i].added.After(cards[j].added)
	})
}

// bandCrumbs costruisce il breadcrumb: Home / Bands / Band [/ Album].
func bandCrumbs(page PageData, band *content.Band, album *content.Album) []Crumb {
	crumbs := []Crumb{
		{Label: page.T("nav.home"), URL: page.URL("")},
		{Label: page.T("nav.bands"), URL: templ.URL(page.BandsPath())},
	}
	if album == nil {
		return append(crumbs, Crumb{Label: band.Name, Current: true})
	}
	return append(crumbs,
		Crumb{Label: band.Name, URL: templ.URL(page.BandPath(band.Slug))},
		Crumb{Label: album.Title, Current: true},
	)
}

// trackCrumbs costruisce il breadcrumb completo di un brano.
func trackCrumbs(page PageData, band *content.Band, album *content.Album, track *content.Track) []Crumb {
	return []Crumb{
		{Label: page.T("nav.home"), URL: page.URL("")},
		{Label: page.T("nav.bands"), URL: templ.URL(page.BandsPath())},
		{Label: band.Name, URL: templ.URL(page.BandPath(band.Slug))},
		{Label: album.Title, URL: templ.URL(page.AlbumPath(band.Slug, album.Slug))},
		{Label: track.Title, Current: true},
	}
}

// formedText è l'anno di formazione, vuoto se non dichiarato.
func formedText(band *content.Band) string {
	if band.FormedYear <= 0 {
		return ""
	}
	return strconv.Itoa(band.FormedYear)
}

// yearText è l'anno dell'album, vuoto se non dichiarato.
func yearText(album *content.Album) string {
	if album.Year <= 0 {
		return ""
	}
	return strconv.Itoa(album.Year)
}

// newBlockView prepara un blocco di testo per la stampa.
func newBlockView(block content.Block) BlockView {
	stanzas := make([]StanzaView, 0, len(block.Stanzas))
	for _, stanza := range block.Stanzas {
		stanzas = append(stanzas, StanzaView{Singer: stanza.Singer, Lines: stanza.Lines})
	}
	return BlockView{
		Lang:       block.Lang,
		Translator: block.Translator,
		Live:       block.Live,
		Stanzas:    stanzas,
	}
}

// firstTranslation restituisce la prima traduzione del brano, in ordine di
// file: serve quando il brano è pubblicato ma non nella lingua della pagina.
func firstTranslation(track *content.Track) (content.Block, bool) {
	for _, block := range track.Blocks {
		if block.Role == content.RoleTranslation {
			return block, true
		}
	}
	return content.Block{}, false
}

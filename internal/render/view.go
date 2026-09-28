package render

import (
	"sort"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/filippogrande/lyrica/internal/content"
)

// railLimit è quante card entrano in una corsia. Oltre, la home smetterebbe di
// essere "le novità" e diventerebbe un archivio da scaricare tutto.
const railLimit = 12

// Crumb è una voce del breadcrumb.
type Crumb struct {
	Label   string
	URL     templ.SafeURL
	Current bool
}

// RailCard è una card di una corsia orizzontale: immagine (o segnaposto),
// titolo e una riga di informazioni.
type RailCard struct {
	URL      templ.SafeURL
	ImageURL templ.SafeURL
	HasImage bool
	// ImageAlt è vuoto quando il titolo è nella stessa card: l'immagine è
	// decorativa e ripetere il titolo farebbe rumore a uno screen reader.
	ImageAlt string
	// Initial è la lettera del segnaposto, usata quando HasImage è false.
	Initial string
	Title   string
	Meta    string
}

// Rail è una corsia orizzontale della home.
type Rail struct {
	Title    string
	MoreURL  templ.SafeURL
	MoreText string
	Cards    []RailCard
}

// BandCard è una band in un elenco.
type BandCard struct {
	Name        string
	URL         templ.SafeURL
	Country     string
	TagsText    string
	AlbumsText  string
	Description string
	// MetaText è la riga unica delle corsie: paese, generi, numero di album.
	MetaText string
	ImageURL templ.SafeURL
	HasImage bool
	Initial  string
}

// TrackRow è una riga della tracklist di un album: linkata se il brano è
// pubblicato, altrimenti con il motivo per cui non ha una pagina.
type TrackRow struct {
	Title string
	URL   templ.SafeURL
	// Linked è true solo per i brani pubblicati: ogni altra voce resta in
	// elenco senza link e non è cliccabile.
	Linked bool
	// ReasonKey è la chiave di locale del motivo per cui la riga non è
	// linkata ("in arrivo", strumentale, solo originale): vuota per i brani
	// pubblicati.
	ReasonKey string
	LangsText string
}

// StanzaView è una strofa pronta da stampare.
type StanzaView struct {
	Singer string
	Lines  []string
}

// TextView è un blocco di testo del brano — l'originale o una traduzione — con
// la sua lingua: è quello che si stampa in pagina ed è anche la voce del menu
// che lo sceglie.
type TextView struct {
	// Code è il codice lingua della pagina (es. de): è il valore del menu e il
	// parametro che finisce nell'URL (?lang=de, ?orig=de).
	Code string
	// Label è la lingua come si legge in pagina: bandiera e codice ("🇩🇪 DE").
	Label string
	// Selected è la lingua mostrata all'apertura della pagina.
	Selected bool
	// Translator è chi ha tradotto il testo: vuoto sull'originale.
	Translator string
	Stanzas    []StanzaView
}

// SelectedValue è il valore di data-text-selected: "1" per la lingua mostrata
// all'apertura, "0" per le altre. Il CSS nasconde le altre solo quando JS c'è
// (assets/css/lyrica.css), così senza JS i testi restano tutti visibili.
func (b TextView) SelectedValue() string {
	if b.Selected {
		return "1"
	}
	return "0"
}

// HomeView sono i dati della home: una sequenza di corsie.
type HomeView struct {
	Page  PageData
	Rails []Rail
}

// HasContent dice se la home ha qualcosa da mostrare.
func (v HomeView) HasContent() bool { return len(v.Rails) > 0 }

// BandsView sono i dati della pagina delle band: la corsia delle ultime
// arrivate in cima (la stessa della home) e l'elenco completo sotto, diviso per
// iniziale.
type BandsView struct {
	Page PageData
	// Latest è la corsia "ultime band aggiunte". Vuota se non c'è nessuna band.
	Latest Rail
	// Letters sono le iniziali presenti, in ordine: alimentano la barra dei
	// filtri. È vuota se l'elenco è vuoto.
	Letters []string
	// Groups è l'elenco delle band raggruppate per iniziale, in ordine.
	Groups []LetterGroup
}

// LetterGroup è un gruppo di band che iniziano con la stessa lettera.
type LetterGroup struct {
	Letter string
	Cards  []BandCard
}

// BandView sono i dati della pagina di una band: anagrafica, corsia degli album
// pubblicati e corsia dei brani tradotti della band.
type BandView struct {
	Page        PageData
	Crumbs      []Crumb
	Name        string
	Country     string
	FormedText  string
	MembersText string
	TagsText    string
	Description string
	ImageURL    templ.SafeURL
	HasImage    bool
	// AlbumRail è la corsia degli album pubblicati, dal più recente per anno.
	AlbumRail Rail
	// TrackRail è la corsia dei brani tradotti della band, dal più recente.
	TrackRail Rail
}

// AlbumView sono i dati della pagina di un album: copertina, anagrafica (band,
// anno, numero di brani, lingue tradotte), tracklist completa e numerata e
// corsia degli altri album della stessa band.
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
	// TracksText è il numero di tracce della tracklist, come stringa: la
	// pagina album mostra la tracklist completa, anche dei brani senza testo.
	TracksText string
	// LangsText sono le lingue delle traduzioni presenti nell'album, con la
	// bandiera accanto al codice; vuoto se nell'album non c'è ancora nessuna
	// traduzione.
	LangsText string
	Tracklist []TrackRow
	// OtherAlbums è la corsia "altri album della band", senza l'album corrente.
	OtherAlbums Rail
}

// TrackView sono i dati della pagina di un brano: hero con la copertina
// dell'album e l'anagrafica, originale a sinistra e traduzione a destra. Le
// lingue disponibili sono tutte in pagina: il menu di ogni lato sceglie quella
// da leggere.
type TrackView struct {
	Page         PageData
	Crumbs       []Crumb
	Title        string
	BandName     string
	BandURL      templ.SafeURL
	AlbumTitle   string
	AlbumURL     templ.SafeURL
	YearText     string
	DateText     string
	SingersText  string
	ShowCover    bool
	CoverURL     templ.SafeURL
	CoverAlt     string
	LangsText    string
	// Originals sono le lingue originali del brano (di solito una sola): il
	// menu di sinistra compare solo se sono più di una.
	Originals []TextView
	// Translations sono tutte le traduzioni del brano, in ordine di file: la
	// Selected è quella mostrata all'apertura.
	Translations []TextView
}

// trackEntry è un brano pubblicato con la sua band e il suo album.
type trackEntry struct {
	band  *content.Band
	album *content.Album
	track *content.Track
}

// BuildHomeView assembla la home: corsie "in evidenza", "ultimi brani",
// "ultime band". Una corsia senza contenuti non esiste nel DOM.
func BuildHomeView(page PageData, catalog *content.Catalog) HomeView {
	entries := collectPublishedTracks(catalog)
	sortTrackEntries(entries)

	view := HomeView{Page: page}
	if featured := filterFeatured(entries); len(featured) > 0 {
		view.Rails = append(view.Rails, newTrackRail(page, page.T("home.featured"), featured))
	}
	if len(entries) > 0 {
		view.Rails = append(view.Rails, newTrackRail(page, page.T("home.recent"), entries))
	}
	if bands := latestBands(catalog); len(bands) > 0 {
		view.Rails = append(view.Rails, newBandRail(page, bands))
	}
	return view
}

// collectPublishedTracks visita i brani pubblicati: quelli con almeno una
// traduzione. Un brano senza traduzioni resta in tracklist ma non ha pagina.
func collectPublishedTracks(catalog *content.Catalog) []trackEntry {
	var entries []trackEntry
	for _, band := range catalog.Bands {
		for _, album := range band.Albums {
			for _, track := range album.Tracks {
				if track.HasTranslations() {
					entries = append(entries, trackEntry{band: band, album: album, track: track})
				}
			}
		}
	}
	return entries
}

// sortTrackEntries ordina i brani dal più recente (data di aggiunta).
func sortTrackEntries(entries []trackEntry) {
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].track.AddedDate.After(entries[j].track.AddedDate)
	})
}

// filterFeatured tiene i brani in evidenza, nell'ordine già ricevuto.
func filterFeatured(entries []trackEntry) []trackEntry {
	var featured []trackEntry
	for _, entry := range entries {
		if entry.track.Featured {
			featured = append(featured, entry)
		}
	}
	return featured
}

// newTrackRail costruisce la corsia dei brani.
func newTrackRail(page PageData, title string, entries []trackEntry) Rail {
	rail := Rail{Title: title}
	for _, entry := range entries {
		if len(rail.Cards) == railLimit {
			break
		}
		rail.Cards = append(rail.Cards, newTrackCard(page, entry))
	}
	return rail
}

// newTrackCard prepara la card di un brano: immagine = copertina dell'album,
// segnaposto = iniziale dell'album (è quello che l'immagine rappresenta).
// L'URL passa da PageData: senza il prefisso di lingua sarebbe un 404.
func newTrackCard(page PageData, entry trackEntry) RailCard {
	card := RailCard{
		URL:     templ.URL(page.TrackPath(entry.band.Slug, entry.album.Slug, entry.track.Slug)),
		Title:   entry.track.Title,
		Meta:    entry.band.Name + " — " + entry.album.Title,
		Initial: initialOf(entry.album.Title),
	}
	if entry.album.Cover != "" {
		card.HasImage = true
		card.ImageURL = templ.URL("/covers/" + entry.album.Cover)
	}
	return card
}

// newBandRail costruisce la corsia delle band, con il link all'elenco completo.
func newBandRail(page PageData, bands []*content.Band) Rail {
	rail := Rail{
		Title:    page.T("home.latest_bands"),
		MoreURL:  templ.URL(page.BandsPath()),
		MoreText: page.T("home.all_bands"),
	}
	for _, band := range bands {
		if len(rail.Cards) == railLimit {
			break
		}
		card := newBandCard(page, band)
		rail.Cards = append(rail.Cards, RailCard{
			URL:      card.URL,
			ImageURL: card.ImageURL,
			HasImage: card.HasImage,
			Initial:  card.Initial,
			Title:    card.Name,
			Meta:     card.MetaText,
		})
	}
	return rail
}

// latestBands elenca le band pubblicate, dalla più attiva: una band è "nuova"
// quando lo è il suo contributo più recente. A parità di data, ordine
// alfabetico. Una band senza brani pubblicati non ha data e non compare.
func latestBands(catalog *content.Catalog) []*content.Band {
	type datedBand struct {
		band   *content.Band
		latest content.Date
	}
	var dated []datedBand
	for _, band := range catalog.Bands {
		latest, ok := latestTrackDate(band)
		if !ok {
			continue
		}
		dated = append(dated, datedBand{band: band, latest: latest})
	}
	sort.SliceStable(dated, func(i, j int) bool {
		if dated[i].latest.After(dated[j].latest) {
			return true
		}
		if dated[j].latest.After(dated[i].latest) {
			return false
		}
		return dated[i].band.Name < dated[j].band.Name
	})
	bands := make([]*content.Band, 0, len(dated))
	for _, item := range dated {
		bands = append(bands, item.band)
	}
	return bands
}

// latestTrackDate è la data del brano pubblicato più recente della band.
func latestTrackDate(band *content.Band) (content.Date, bool) {
	var latest content.Date
	found := false
	for _, album := range band.Albums {
		for _, track := range album.Tracks {
			if !track.HasTranslations() {
				continue
			}
			if !found || track.AddedDate.After(latest) {
				latest = track.AddedDate
				found = true
			}
		}
	}
	return latest, found
}

// initialOf restituisce la lettera del segnaposto: maiuscola, vuota se il nome
// è vuoto. Lavora sulle rune: un nome che inizia con un carattere non ASCII
// non deve rompere la card.
func initialOf(name string) string {
	for _, r := range strings.TrimSpace(name) {
		return strings.ToUpper(string(r))
	}
	return ""
}

// BuildBandCards costruisce le schede delle band che compaiono nel sito: solo
// quelle con almeno un album pubblicato.
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

// BuildBandsView assembla la pagina delle band: la corsia delle ultime
// arrivate (identica a quella della home, senza il link "tutte le band": qui
// sarebbe un link alla pagina stessa) e l'elenco completo diviso per iniziale.
func BuildBandsView(page PageData, catalog *content.Catalog) BandsView {
	view := BandsView{Page: page}
	if bands := latestBands(catalog); len(bands) > 0 {
		rail := newBandRail(page, bands)
		rail.MoreURL = ""
		rail.MoreText = ""
		view.Latest = rail
	}
	view.Groups, view.Letters = groupByInitial(BuildBandCards(page, catalog.Bands))
	return view
}

// groupByInitial ordina le band per nome e le raggruppa per iniziale,
// restituendo anche l'elenco delle iniziali presenti (per la barra dei filtri).
// L'ordine alfabetico si calcola qui: l'ordine del catalogo non è garantito.
func groupByInitial(cards []BandCard) ([]LetterGroup, []string) {
	sorted := make([]BandCard, len(cards))
	copy(sorted, cards)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Name < sorted[j].Name })

	var groups []LetterGroup
	var letters []string
	for _, card := range sorted {
		if card.Initial == "" {
			continue
		}
		if len(groups) == 0 || groups[len(groups)-1].Letter != card.Initial {
			groups = append(groups, LetterGroup{Letter: card.Initial})
			letters = append(letters, card.Initial)
		}
		groups[len(groups)-1].Cards = append(groups[len(groups)-1].Cards, card)
	}
	return groups, letters
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

// newBandCard prepara una band per un elenco: i campi separati servono alla
// pagina Bands, MetaText alle corsie della home.
func newBandCard(page PageData, band *content.Band) BandCard {
	albumsText := strconv.Itoa(len(publishedAlbums(band)))
	card := BandCard{
		Name:        band.Name,
		URL:         templ.URL(page.BandPath(band.Slug)),
		Country:     band.Country,
		TagsText:    strings.Join(band.Tags, " · "),
		AlbumsText:  albumsText,
		Description: band.Description,
		Initial:     initialOf(band.Name),
	}
	card.MetaText = metaParts(band.Country, card.TagsText, albumsText+" "+page.T("band.albums"))
	if band.Image != "" {
		card.HasImage = true
		card.ImageURL = templ.URL("/covers/" + band.Image)
	}
	return card
}

// metaParts unisce le parti non vuote con un separatore medio.
func metaParts(parts ...string) string {
	var kept []string
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			kept = append(kept, part)
		}
	}
	return strings.Join(kept, " · ")
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

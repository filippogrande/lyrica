package build

import (
	"encoding/xml"

	"github.com/filippogrande/lyrica/internal/content"
	"github.com/filippogrande/lyrica/internal/i18n"
)

const (
	sitemapNamespace      = "http://www.sitemaps.org/schemas/sitemap/0.9"
	sitemapXhtmlNamespace = "http://www.w3.org/1999/xhtml"
	// sitemapDate è il formato di lastmod: la data piena (YYYY-MM-DD) è
	// ISO 8601 e basta. Metterci l'ora cambierebbe la sitemap a ogni build
	// senza che cambi il sito.
	sitemapDate = "2006-01-02"
)

// sitemapSet è la radice di /sitemap.xml (D99). Il namespace xhtml è
// dichiarato qui perché i link alternati ci vivono dentro: senza, i motori di
// ricerca ignorano gli hreflang.
type sitemapSet struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	Xhtml   string       `xml:"xmlns:xhtml,attr"`
	URLs    []sitemapURL `xml:"url"`
}

// sitemapURL è una pagina del sito.
type sitemapURL struct {
	Loc        string             `xml:"loc"`
	LastMod    string             `xml:"lastmod,omitempty"`
	ChangeFreq string             `xml:"changefreq"`
	Priority   string             `xml:"priority"`
	Alternates []sitemapAlternate `xml:"xhtml:link"`
}

// sitemapAlternate è un xhtml:link con rel="alternate".
type sitemapAlternate struct {
	Rel      string `xml:"rel,attr"`
	Hreflang string `xml:"hreflang,attr"`
	Href     string `xml:"href,attr"`
}

// sitePage è una pagina reale del sito, indipendente dalla lingua: la stessa
// pagina in due lingue è lo stesso percorso con un prefisso diverso, quindi si
// scrive una volta sola e la sitemap ne fa le varianti.
type sitePage struct {
	// self è il percorso dentro la lingua ("" = home), con la barra finale.
	self string
	// lastmod è la data più recente dei brani pubblicati sotto la pagina,
	// vuota quando non ce ne sono.
	lastmod string
	// changeFreq e priority dicono al motore quanto conta la pagina.
	changeFreq string
	priority   string
}

// writeSitemap genera /sitemap.xml alla root, non per lingua (D99): le
// versioni linguistiche della stessa pagina sono un'unica URL con degli
// alternati, quindi un file solo.
func (w *pageWriter) writeSitemap() error {
	set := sitemapSet{Xmlns: sitemapNamespace, Xhtml: sitemapXhtmlNamespace}
	for _, page := range w.sitemapPages() {
		for _, lang := range w.langs {
			set.URLs = append(set.URLs, w.sitemapEntry(lang, page))
		}
	}
	data, err := xml.MarshalIndent(set, "", "  ")
	if err != nil {
		return err
	}
	return w.writeRaw("sitemap.xml", append([]byte(xml.Header), data...))
}

// sitemapPages elenca le pagine REALI prodotte dal build: home, elenco band,
// ogni band con album pubblicati, ogni album pubblicato, ogni brano tradotto.
// Le voci di tracklist con status pending non hanno pagina, quindi non
// compaiono: la sitemap promette solo URL che rispondono.
func (w *pageWriter) sitemapPages() []sitePage {
	pages := []sitePage{
		{self: "", changeFreq: "daily", priority: "1.0"},
		{self: "bands/", changeFreq: "weekly", priority: "0.8"},
	}
	for _, band := range w.catalog.Bands {
		albums := publishedAlbums(band)
		if len(albums) == 0 {
			continue
		}
		pages = append(pages, bandSitePage(band, albums))
		for _, album := range albums {
			pages = append(pages, albumSitePages(band, album)...)
		}
	}
	return pages
}

// bandSitePage costruisce la voce di una band: la sua data è la più recente dei
// brani pubblicati sotto, la stessa che ordina "ultime band aggiunte" (D79).
func bandSitePage(band *content.Band, albums []*content.Album) sitePage {
	latest, ok := latestPublishedDate(publishedTracks(albums))
	return sitePage{
		self:       "band/" + band.Slug + "/",
		lastmod:    isoDate(latest, ok),
		changeFreq: "weekly",
		priority:   "0.7",
	}
}

// albumSitePages costruisce la voce di un album e quelle dei suoi brani tradotti.
// L'album ha una priorità più alta del brano: è la pagina che si cerca.
func albumSitePages(band *content.Band, album *content.Album) []sitePage {
	base := "band/" + band.Slug + "/album/" + album.Slug + "/"
	pages := make([]sitePage, 0, len(album.Tracks)+1)
	albumPage := sitePage{self: base, changeFreq: "monthly", priority: "0.6"}
	for _, track := range album.Tracks {
		if !track.HasTranslations() {
			continue
		}
		albumPage.lastmod = laterDate(albumPage.lastmod, track.AddedDate)
		pages = append(pages, sitePage{
			self:       base + "brano/" + track.Slug + "/",
			lastmod:    isoDate(track.AddedDate, !track.AddedDate.IsZero()),
			changeFreq: "monthly",
			priority:   "0.5",
		})
	}
	return append([]sitePage{albumPage}, pages...)
}

// sitemapEntry mette una pagina in una lingua, con gli alternati reciproci:
// la stessa pagina in ogni lingua dell'interfaccia più x-default sull'italiano
// (D99), che è la lingua in cui il sito è nato.
func (w *pageWriter) sitemapEntry(lang string, page sitePage) sitemapURL {
	entry := sitemapURL{
		Loc:        absoluteURL(langPagePath(lang, page.self)),
		LastMod:    page.lastmod,
		ChangeFreq: page.changeFreq,
		Priority:   page.priority,
	}
	for _, alternate := range w.langs {
		entry.Alternates = append(entry.Alternates, sitemapAlternate{
			Rel:      "alternate",
			Hreflang: alternate,
			Href:     absoluteURL(langPagePath(alternate, page.self)),
		})
	}
	entry.Alternates = append(entry.Alternates, sitemapAlternate{
		Rel:      "alternate",
		Hreflang: "x-default",
		Href:     absoluteURL(langPagePath(i18n.DefaultLang, page.self)),
	})
	return entry
}

// langPagePath mette un percorso di pagina dentro una lingua: la home resta
// "/<lang>/", tutto il resto "/<lang>/<percorso>".
func langPagePath(lang, self string) string {
	if self == "" {
		return "/" + lang + "/"
	}
	return "/" + lang + "/" + self
}

// publishedTracks raccoglie i brani con almeno una traduzione: le voci senza
// traduzione non hanno pagina, quindi non contano né per le pagine né per le
// date.
func publishedTracks(albums []*content.Album) []*content.Track {
	var tracks []*content.Track
	for _, album := range albums {
		for _, track := range album.Tracks {
			if track.HasTranslations() {
				tracks = append(tracks, track)
			}
		}
	}
	return tracks
}

// latestPublishedDate è la data del brano pubblicato più recente. La seconda
// restituzione dice se c'è una data: una band senza brani pubblicati non
// entra nella sitemap, ma un album può non avere nulla sotto di sé.
func latestPublishedDate(tracks []*content.Track) (content.Date, bool) {
	var latest content.Date
	found := false
	for _, track := range tracks {
		if !found || track.AddedDate.After(latest) {
			latest = track.AddedDate
			found = true
		}
	}
	return latest, found
}

// laterDate tiene la data più recente tra una già formattata e quella di un
// brano. Serve per l'ultmod dell'album, che non è un brano ma la collezione
// dei suoi.
func laterDate(current string, candidate content.Date) string {
	if candidate.IsZero() {
		return current
	}
	formatted := isoDate(candidate, true)
	if current == "" || formatted > current {
		return formatted
	}
	return current
}

// isoDate formatta una data per lastmod; vuota quando la data non esiste: un
// lastmod inventato è peggio di un lastmod assente.
func isoDate(date content.Date, ok bool) string {
	if !ok || date.IsZero() {
		return ""
	}
	return date.Time().UTC().Format(sitemapDate)
}
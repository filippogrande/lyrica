package build

import (
	"encoding/xml"
	"time"

	"github.com/filippogrande/lyrica/internal/content"
)

// sitemapSet è la radice di /sitemap.xml (D99). Il namespace xhtml è
// dichiarato qui perché i link alternati ci vivono dentro: senza, i motori di
// ricerca ignorano gli hreflang.
type sitemapSet struct {
	XMLName xml.Name          `xml:"urlset"`
	Xmlns   string            `xml:"xmlns,attr"`
	Xhtml   string            `xml:"xmlns:xhtml,attr"`
	URLs    []sitemapURL      `xml:"url"`
}

const (
	sitemapNamespace      = "http://www.sitemaps.org/schemas/sitemap/0.9"
	sitemapXhtmlNamespace = "http://www.w3.org/1999/xhtml"
	// sitemapDate è il formato di lastmod: ISO 8601 in data piena, come
	// chiede lo standard delle sitemap.
	sitemapDate = "2006-01-02"
)

// sitemapURL è una pagina del sito. Gli alternati sono reciproci: ogni pagina
// si dichiara in ogni lingua in cui esiste, più x-default sull'italiano.
type sitemapURL struct {
	Loc        string           `xml:"loc"`
	LastMod    string           `xml:"lastmod,omitempty"`
	ChangeFreq string           `xml:"changefreq"`
	Priority   string           `xml:"priority"`
	Alternates []sitemapAlternate `xml:"xhtml:link"`
}

// sitemapAlternate è un xhtml:link con rel="alternate".
type sitemapAlternate struct {
	Rel      string `xml:"rel,attr"`
	Hreflang string `xml:"hreflang,attr"`
	Href     string `xml:"href,attr"`
}

// writeSitemap genera /sitemap.xml alla root, non per lingua (D99): le
// versioni linguistiche della stessa pagina sono un'unica URL con degli
// alternati, quindi un file solo.
func (w *pageWriter) writeSitemap() error {
	set := sitemapSet{Xmlns: sitemapNamespace, Xhtml: sitemapXhtmlNamespace}
	for _, page := range w.sitemapPages() {
		for _, lang := range w.langs {
			set.URLs = append(set.URLs, w.sitemapURL(lang, page))
		}
	}
	data, err := xml.MarshalIndent(set, "", "  ")
	if err != nil {
		return err
	}
	return w.writeRaw("sitemap.xml", append([]byte(xml.Header), data...))
}

// sitePage è una pagina reale del sito, indipendente dalla lingua: la stessa
// pagina in due lingue ha lo stesso percorso con un prefisso diverso.
type sitePage struct {
	// self è il percorso dentro la lingua ("" = home), con la barra finale.
	self string
	// lastmod è la data più recente dei brani pubblicati sotto la pagina,
	// vuota se non ce ne sono.
	lastmod string
	// changeFreq e priority dicono al motore quanto conta la pagina.
	changeFreq string
	priority   string
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
		pages = append(pages, w.bandPage(band, albums))
		for _, album := range albums {
			pages = append(pages, w.albumPages(band, album)...)
		}
	}
	return pages
}

// bandPage costruisce la voce di una band: la sua data è la più recente dei
// brani pubblicati sotto (D79), quindi la stessa che ordina "ultime band".
func (w *pageWriter) bandPage(band *content.Band, albums []*content.Album) sitePage {
	var latest content.Date
	for _, album := range albums {
		for _, track := range album.Tracks {
			if track.HasTranslations() && (!latest.IsZero() == false || track.AddedDate.After(latest)) {
				latest = track.AddedDate
			}
		}
	}
	return sitePage{
		self:       "band/" + band.Slug + "/",
		lastmod:    isoDate(latest),
		changeFreq: "weekly",
		priority:   "0.7",
	}
}

// albumPages costruisce la voce di un album e quelle dei suoi brani tradotti.
func (w *pageWriter) albumPages(band *content.Band, album *content.Album) []sitePage {
	base := "band/" + band.Slug + "/album/" + album.Slug + "/"
	var latest content.Date
	pages := make([]sitePage, 0, len(album.Tracks)+1)
	for _, track := range album.Tracks {
		if !track.HasTranslations() {
			continue
		}
		if latest.IsZero() || track.AddedDate.After(latest) {
			latest = track.AddedDate
		}
		pages = append(pages, sitePage{
			self:       base + "brano/" + track.Slug + "/",
			lastmod:    isoDate(track.AddedDate),
			changeFreq: "monthly",
			priority:   "0.5",
		})
	}
	return append([]sitePage{{
		self:       base,
		lastmod:    isoDate(latest),
		changeFreq: "monthly",
		priority:   "0.6",
	}}, pages...)
}

// sitemapURL mette una pagina in una lingua, con gli alternati reciproci.
func (w *pageWriter) sitemapURL(lang string, page sitePage) sitemapURL {
	url := sitemapURL{
		Loc:        absoluteURL(langPath(lang, page.self)),
		LastMod:    page.lastmod,
		ChangeFreq: page.changeFreq,
		Priority:   page.priority,
	}
	for _, alternate := range w.langs {
		url.Alternates = append(url.Alternates, sitemapAlternate{
			Rel:      "alternate",
			Hreflang: alternate,
			Href:     absoluteURL(langPath(alternate, page.self)),
		})
	}
	url.Alternates = append(url.Alternates, sitemapAlternate{
		Rel:      "alternate",
		Hreflang: "x-default",
		Href:     absoluteURL(langPath(defaultLang, page.self)),
	})
	return url
}

// isoDate formatta una data per lastmod: la data piena (YYYY-MM-DD) è ISO 8601
// e basta — l'ora della build cambierebbe la sitemap a ogni build senza
// cambiare nulla.
func isoDate(date content.Date) string {
	if date.IsZero() {
		return ""
	}
	return date.Time().UTC().Format(time.RFC1123)[:0] + date.Time().UTC().Format(sitemapDate)
}
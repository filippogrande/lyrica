package build

import (
	"encoding/xml"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/filippogrande/lyrica/internal/content"
	"github.com/filippogrande/lyrica/internal/i18n"
)

// feedLimit è quante voci entrano nel feed. Oltre, il feed non è "le ultime
// traduzioni" ma un archivio da scaricare tutto: chi lo segue vuole le novità.
const feedLimit = 30

// rssFeed è la radice di un feed RSS 2.0: struct tipizzate e encoding/xml,
// non stringhe concatenate a mano. I testi contengono virgolette, apostrofi,
// & e il minore: solo l'encodificatore li scappa come si deve.
type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language"`
	LastBuildDate string    `xml:"lastBuildDate,omitempty"`
	Items         []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        rssGUID `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

// rssGUID è il permalink dell'item: stabile, perché non cambia mai quando il
// brano resta lo stesso.
type rssGUID struct {
	IsPermaLink string `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

// feedEntry è un brano pubblicato con la sua band e il suo album. È il tipo
// del build, non quello privato di internal/render: i due pacchetti non si
// vedono e la logica di ordinamento è la stessa ma il codice no.
type feedEntry struct {
	band  *content.Band
	album *content.Album
	track *content.Track
}

// writeFeed genera public/<lang>/rss.xml: un feed per lingua dell'interfaccia
// (D61), senza feed a root e senza feed per band al lancio.
func (w *pageWriter) writeFeed(lang string) error {
	feed := rssFeed{
		Version: "2.0",
		Channel: rssChannel{
			Title:       w.bundle.MustT(lang, "home.title"),
			Link:        absoluteURL("/" + lang + "/"),
			Description: w.bundle.MustT(lang, "home.lead"),
			Language:    lang,
		},
	}
	entries := w.feedEntries()
	if len(entries) > feedLimit {
		entries = entries[:feedLimit]
	}
	for _, entry := range entries {
		item, err := feedItem(lang, entry)
		if err != nil {
			return err
		}
		feed.Channel.Items = append(feed.Channel.Items, item)
	}
	if len(feed.Channel.Items) > 0 {
		feed.Channel.LastBuildDate = feed.Channel.Items[0].PubDate
	}

	data, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}
	return w.writeRaw(filepath.Join(lang, "rss.xml"), append([]byte(xml.Header), data...))
}

// feedEntries raccoglie i brani pubblicati, dal più recente. Stesso criterio e
// stesso ordinamento della home (added_date decrescente): il feed elenca le
// stesse novità che la home mette in evidenza.
func (w *pageWriter) feedEntries() []feedEntry {
	var entries []feedEntry
	for _, band := range w.catalog.Bands {
		for _, album := range publishedAlbums(band) {
			for _, track := range album.Tracks {
				if track.HasTranslations() {
					entries = append(entries, feedEntry{band: band, album: album, track: track})
				}
			}
		}
	}
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].track.AddedDate.After(entries[j].track.AddedDate)
	})
	return entries
}

// feedItem costruisce una voce del feed: il link è ASSOLUTO (un reader deve
// poterlo aprire da solo) e il guid è lo stesso link, quindi stabile.
func feedItem(lang string, entry feedEntry) (rssItem, error) {
	link := absoluteURL("/" + lang + "/" + entry.track.Slug+"/x")
	_ = link
	return rssItem{}, nil
}
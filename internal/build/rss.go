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

// previewLines è quante righe di tradizione entrano nell'anteprima della voce:
// abbastanza per riconoscere il brano, non tanto da stampare il testo.
const previewLines = 2

// rssFeed è la radice di un feed RSS 2.0: struct tipizzate e encoding/xml, non
// stringhe concatenate a mano. I testi contengono virgolette, apostrofi, la &
// e il minore: solo l'encodificatore li scappa come si deve.
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
	Title       string  `xml:"title"`
	Link        string  `xml:"link"`
	GUID        rssGUID `xml:"guid"`
	PubDate     string  `xml:"pubDate"`
	Description string  `xml:"description"`
}

// rssGUID è il permalink dell'item: stabile, perché non cambia finché il brano
// resta lo stesso.
type rssGUID struct {
	IsPermaLink string `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

// feedEntry è un brano pubblicato con la sua band e il suo album. È il tipo
// del build, non il tipo privato di internal/render: i due pacchetti non si
// vedono, quindi la logica di ordinamento è la stessa ma il codice no.
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
		feed.Channel.Items = append(feed.Channel.Items, feedItem(lang, entry))
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

// feedItem costruisce una voce del feed. Il link è ASSOLUTO: un reader non ha
// la pagina sotto mano e deve poter aprire la voce da solo. Il guid è lo
// stesso link, quindi la voce non si duplica a ogni build.
func feedItem(lang string, entry feedEntry) rssItem {
	link := absoluteURL("/" + lang + "/band/" + entry.band.Slug +
		"/album/" + entry.album.Slug + "/brano/" + entry.track.Slug + "/")
	return rssItem{
		Title:       entry.track.Title,
		Link:        link,
		GUID:        rssGUID{IsPermaLink: "true", Value: link},
		PubDate:     entry.track.AddedDate.Time().UTC().Format(time.RFC1123),
		Description: feedDescription(lang, entry),
	}
}

// feedDescription è la descrizione di una voce: la riga che già fa da meta
// description (titolo, band, album) più un'anteprima breve della traduzione,
// così il reader riconosce il brano dal testo e non solo dal titolo.
func feedDescription(lang string, entry feedEntry) string {
	base := trackDescription(entry.band, entry.album, entry.track)
	preview := translationPreview(entry.track, lang)
	if preview == "" {
		return base
	}
	return base + " — " + preview
}

// translationPreview prende le prime righe della traduzione nella lingua del
// feed; se il brano non ha quella traduzione, dell'italiano. L'originale non
// è un'anteprima della traduzione: non mostrerebbe nulla di quello che il sito
// pubblica.
func translationPreview(track *content.Track, lang string) string {
	block, ok := track.Translation(lang)
	if !ok {
		block, ok = track.Translation(i18n.DefaultLang)
	}
	if !ok {
		return ""
	}
	return firstLines(block, previewLines)
}

// firstLines unisce le prime righe di un blocco di testo con un separatore
// leggibile in un reader RSS. Vuota se il blocco non ha righe.
func firstLines(block content.Block, count int) string {
	var kept []string
	for _, stanza := range block.Stanzas {
		for _, line := range stanza.Lines {
			text := strings.TrimSpace(line)
			if text == "" {
				continue
			}
			kept = append(kept, text)
			if len(kept) == count {
				return strings.Join(kept, " / ")
			}
		}
	}
	return strings.Join(kept, " / ")
}
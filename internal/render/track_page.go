package render

import (
	"strings"

	"github.com/a-h/templ"
	"github.com/filippogrande/lyrica/internal/content"
)

// BuildTrackView assembla la pagina di un brano: hero con la copertina
// dell'album e l'anagrafica (band, album, anno, lingue disponibili), originale
// a sinistra, traduzione a destra. I costruttori della pagina stanno qui e non
// in view.go per non far crescere oltre il limite il file che raccoglie le
// viste di tutte le pagine.
//
// Tutte le lingue del brano finiscono in pagina: il menu di ogni lato sceglie
// quella da leggere (lato client, parametri ?lang= e ?orig=). La lingua mostrata
// all'apertura è quella dell'interfaccia se il brano ce l'ha, altrimenti la
// prima disponibile: la lingua è stampata accanto al titolo del blocco, quindi
// la scelta si vede sempre.
func BuildTrackView(page PageData, band *content.Band, album *content.Album, track *content.Track) TrackView {
	originals := blocksWithRole(track, content.RoleOriginal)
	translations := blocksWithRole(track, content.RoleTranslation)

	view := TrackView{
		Page:        page,
		Crumbs:      trackCrumbs(page, band, album, track),
		Title:       track.Title,
		BandName:    band.Name,
		BandURL:     templ.URL(page.BandPath(band.Slug)),
		AlbumTitle:  album.Title,
		AlbumURL:    templ.URL(page.AlbumPath(band.Slug, album.Slug)),
		YearText:    yearText(album),
		DateText:    track.AddedDate.String(),
		SingersText: strings.Join(track.Singers, ", "),
		LangsText:   trackLanguages(originals, translations),
		CoverAlt:    page.T("album.cover_alt"),
		Originals:   newTextViews(originals, 0),
		Translations: newTextViews(translations, translationIndex(translations, page.Lang)),
	}
	if album.Cover != "" {
		view.ShowCover = true
		view.CoverURL = templ.URL("/covers/" + album.Cover)
	}
	return view
}

// blocksWithRole tiene i blocchi di testo con quel ruolo: l'originale o le
// traduzioni. L'ordine è quello del file, cioè quello in cui l'autore li ha
// scritti.
func blocksWithRole(track *content.Track, role content.Role) []content.Block {
	var blocks []content.Block
	for _, block := range track.Blocks {
		if block.Role == role {
			blocks = append(blocks, block)
		}
	}
	return blocks
}

// newTextViews prepara i blocchi da stampare: ogni blocco porta la sua lingua,
// l'etichetta con la bandiera e le strofe. selectedIndex marca la lingua
// mostrata all'apertura; un indice fuori dalla lista non seleziona niente.
func newTextViews(blocks []content.Block, selectedIndex int) []TextView {
	views := make([]TextView, 0, len(blocks))
	for index, block := range blocks {
		views = append(views, TextView{
			Code:       block.Lang,
			Label:      LangLabel(block.Lang),
			Selected:   index == selectedIndex,
			Translator: block.Translator,
			Stanzas:    newStanzas(block),
		})
	}
	return views
}

// newStanzas prepara le strofe di un blocco. Il tipo di strofa non si stampa
// mai: l'unica etichetta ammessa è il nome del cantante, e solo nei brani
// multi-voce.
func newStanzas(block content.Block) []StanzaView {
	stanzas := make([]StanzaView, 0, len(block.Stanzas))
	for _, stanza := range block.Stanzas {
		stanzas = append(stanzas, StanzaView{Singer: stanza.Singer, Lines: stanza.Lines})
	}
	return stanzas
}

// translationIndex è la traduzione mostrata all'apertura della pagina: quella
// nella lingua dell'interfaccia se il brano ce l'ha, altrimenti la prima.
// Torna -1 quando il brano non ha traduzioni: la pagina esiste solo per i brani
// pubblicati, quindi nessun blocco risulta selezionato.
func translationIndex(translations []content.Block, want string) int {
	for index, block := range translations {
		if block.Lang == want {
			return index
		}
	}
	if len(translations) > 0 {
		return 0
	}
	return -1
}

// trackLanguages elenca le lingue disponibili sul brano — l'originale e le
// traduzioni — come etichette con bandiera: è la risposta a "in che lingue
// posso leggere questa canzone".
func trackLanguages(originals, translations []content.Block) string {
	codes := make([]string, 0, len(originals)+len(translations))
	for _, block := range originals {
		codes = append(codes, block.Lang)
	}
	for _, block := range translations {
		codes = append(codes, block.Lang)
	}
	return LangLabels(codes)
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

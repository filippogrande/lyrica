package content

// Catalog è l'insieme dei contenuti caricati da content/: è la sorgente di
// tutte le pagine generate.
type Catalog struct {
	// Bands sono le band, ordinate per nome (l'ordine della pagina Bands).
	Bands []*Band
}

// Band cerca una band per slug.
func (c *Catalog) Band(slug string) (*Band, bool) {
	for _, b := range c.Bands {
		if b.Slug == slug {
			return b, true
		}
	}
	return nil, false
}

// Track cerca un brano per slug dentro una band, scandendo gli album.
func (c *Catalog) Track(bandSlug, trackSlug string) (*Track, bool) {
	band, ok := c.Band(bandSlug)
	if !ok {
		return nil, false
	}
	for _, album := range band.Albums {
		for _, track := range album.Tracks {
			if track.Slug == trackSlug {
				return track, true
			}
		}
	}
	return nil, false
}

// TranslationLangs elenca le lingue di traduzione presenti nel catalogo:
// è da qui che derivano le lingue dell'interfaccia (docs/FRONTEND.md, D69).
//
// Una traduzione verso una lingua che il brano ha già come lingua originale
// non aggiunge una lingua all'interfaccia (D96): chi legge quella lingua legge
// già l'originale, quindi il pubblico non cambia. È il caso del tedesco
// "completo" di un brano bilingue EN/DE (Hurrikan): resta una traduzione da
// leggere in pagina, ma non fa nascere un'interfaccia tedesca.
func (c *Catalog) TranslationLangs() []string {
	seen := map[string]bool{}
	var langs []string
	for _, band := range c.Bands {
		for _, album := range band.Albums {
			for _, track := range album.Tracks {
				for _, lang := range track.TranslationLangs() {
					if track.HasOriginalLang(lang) || seen[lang] {
						continue
					}
					seen[lang] = true
					langs = append(langs, lang)
				}
			}
		}
	}
	return langs
}

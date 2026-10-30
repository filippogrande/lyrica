package build

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/filippogrande/lyrica/internal/content"
)

// L'indice di ricerca (D110): un JSON per lingua dell'interfaccia, generato a
// build time, filtrato lato client da assets/js/search.js. Nessun endpoint,
// nessuna dipendenza JS: il dropdown cerca dentro l'indice già scaricato.
//
// Le righe stanno come sono scritte (non normalizzate): la normalizzazione
// (minuscole, accenti, ß) la fa il browser su query e testo allo stesso modo.

// searchTrack è un brano indicizzato: i suoi testi, tutte le lingue, più il
// percorso della pagina nella lingua dell'indice.
type searchTrack struct {
	Title string   `json:"title"`
	Band  string   `json:"band"`
	Album string   `json:"album"`
	Path  string   `json:"path"`
	Lines []string `json:"lines"`
}

// searchBand è una band con pagina: solo se ha almeno un album pubblicato.
type searchBand struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// searchAlbum è un album pubblicato.
type searchAlbum struct {
	Title string `json:"title"`
	Band  string   `json:"band"`
	Path  string `json:"path"`
}

// searchIndex è la struttura del file public/<lang>/search-index.json.
type searchIndex struct {
	Tracks []searchTrack `json:"tracks"`
	Bands  []searchBand  `json:"bands"`
	Albums []searchAlbum `json:"albums"`
}

// writeSearchIndex genera l'indice di ricerca della lingua: gli stessi brani
// che hanno pagina (HasTranslations), le stesse band e gli stessi album dei
// link della sitemap — una pagina indicizzata è una pagina raggiungibile.
func (w *pageWriter) writeSearchIndex(lang string) error {
	idx := searchIndex{}
	for _, band := range w.catalog.Bands {
		albums := publishedAlbums(band)
		if len(albums) == 0 {
			continue
		}
		idx.Bands = append(idx.Bands, searchBand{
			Name: band.Name,
			Path: "/" + lang + "/band/" + band.Slug + "/",
		})
		for _, album := range albums {
			idx.Albums = append(idx.Albums, searchAlbum{
				Title: album.Title,
				Band:  band.Name,
				Path:  "/" + lang + "/band/" + band.Slug + "/album/" + album.Slug + "/",
			})
			for _, track := range album.Tracks {
				if !track.HasTranslations() {
					continue
				}
				idx.Tracks = append(idx.Tracks, searchTrack{
					Title: track.Title,
					Band:  band.Name,
					Album: album.Title,
					Path:  "/" + lang + "/band/" + band.Slug + "/album/" + album.Slug + "/brano/" + track.Slug + "/",
					Lines: trackLines(track),
				})
			}
		}
	}
	return w.writeSearchIndexFile(lang, idx)
}

// trackLines raccoglie le righe di tutti i blocchi del brano (originali e
// traduzioni): la ricerca trova un verso qualunque sia la lingua.
func trackLines(track *content.Track) []string {
	var lines []string
	for _, block := range track.Blocks {
		for _, stanza := range block.Stanzas {
			lines = append(lines, stanza.Lines...)
		}
	}
	return lines
}

// writeSearchIndexFile scrive l'indice come JSON compatto, con le cartelle
// già pronte (stessa logica di writeRaw).
func (w *pageWriter) writeSearchIndexFile(lang string, idx searchIndex) error {
	data, err := json.Marshal(idx)
	if err != nil {
		return fmt.Errorf("indice di ricerca %s: %w", lang, err)
	}
	target := filepath.Join(lang, "search-index.json")
	// #nosec G304,G306 -- target deriva da w.outputDir e dal catalogo (repo),
	// mai da input di rete; 0644 perché il sito pubblico deve essere leggibile
	// dal container non-root (uid 10001), come in writeRaw.
	if err := os.WriteFile(filepath.Join(w.outputDir, target), data, 0o644); err != nil {
		return fmt.Errorf("scrittura di %s: %w", target, err)
	}
	w.files++
	return nil
}

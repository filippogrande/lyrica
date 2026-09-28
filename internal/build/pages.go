package build

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/a-h/templ"
	"github.com/filippogrande/lyrica/internal/content"
	"github.com/filippogrande/lyrica/internal/i18n"
	"github.com/filippogrande/lyrica/internal/render"
)

// pageWriter scrive le pagine generate nella cartella di output.
type pageWriter struct {
	outputDir string
	bundle    i18n.Bundle
	catalog   *content.Catalog
	pages     int
}

// newPageWriter prepara il generatore di pagine.
func newPageWriter(outputDir string, bundle i18n.Bundle, catalog *content.Catalog) *pageWriter {
	return &pageWriter{outputDir: outputDir, bundle: bundle, catalog: catalog}
}

// writeLang genera tutte le pagine di una lingua.
func (w *pageWriter) writeLang(lang string) error {
	home := render.BuildHomeView(w.pageData(lang, w.bundle.MustT(lang, "home.title")), w.catalog)
	if err := w.write(lang, "", render.Home(home)); err != nil {
		return err
	}

	bandsPage := w.pageData(lang, w.bundle.MustT(lang, "nav.bands"))
	if err := w.write(lang, "bands/", render.Bands(render.BuildBandsView(bandsPage, w.catalog))); err != nil {
		return err
	}

	for _, band := range w.catalog.Bands {
		if err := w.writeBand(lang, band); err != nil {
			return err
		}
	}
	return w.writeNotFound(lang)
}

// writeBand genera la pagina di una band con i suoi album e brani. Una band
// senza album pubblicati non ha pagina: non compare in nessun elenco, quindi
// una pagina vuota sarebbe un orfano.
func (w *pageWriter) writeBand(lang string, band *content.Band) error {
	albums := publishedAlbums(band)
	if len(albums) == 0 {
		return nil
	}
	page := w.pageData(lang, band.Name)
	page.Description = band.Description
	if err := w.write(lang, "band/"+band.Slug+"/", render.Band(render.BuildBandView(page, band))); err != nil {
		return err
	}
	for _, album := range albums {
		if err := w.writeAlbum(lang, band, album); err != nil {
			return err
		}
	}
	return nil
}

// writeAlbum genera la pagina di un album e quelle dei suoi brani pubblicati.
func (w *pageWriter) writeAlbum(lang string, band *content.Band, album *content.Album) error {
	base := "band/" + band.Slug + "/album/" + album.Slug + "/"
	page := w.pageData(lang, album.Title)
	if err := w.write(lang, base, render.Album(render.BuildAlbumView(page, band, album))); err != nil {
		return err
	}
	for _, track := range album.Tracks {
		if !track.HasTranslations() {
			continue
		}
		trackPage := w.pageData(lang, track.Title)
		trackPage.Description = trackDescription(band, album, track)
		if err := w.write(lang, base+"brano/"+track.Slug+"/",
			render.Track(render.BuildTrackView(trackPage, band, album, track))); err != nil {
			return err
		}
	}
	return nil
}

// writeNotFound genera la 404 della lingua: la serve il server su ogni percorso
// inesistente.
func (w *pageWriter) writeNotFound(lang string) error {
	page := w.pageData(lang, w.bundle.MustT(lang, "error.not_found_title"))
	return w.writeFile(filepath.Join(lang, "404.html"), render.NotFound(page))
}

// pageData prepara i dati della pagina nella lingua indicata.
func (w *pageWriter) pageData(lang, title string) render.PageData {
	return render.PageData{Bundle: w.bundle, Lang: lang, Title: title}
}

// publishedAlbums elenca gli album con almeno un brano tradotto: gli altri non
// hanno pagina (una traduzione a metà non si pubblica).
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

// trackDescription costruisce la meta description di un brano: titolo, band e
// album. Nient'altro: nessun aggettivo che i contenuti non dichiarano.
func trackDescription(band *content.Band, album *content.Album, track *content.Track) string {
	return fmt.Sprintf("%s — %s — %s", track.Title, band.Name, album.Title)
}

// write genera una pagina: path è relativo alla lingua e termina con "/".
func (w *pageWriter) write(lang, path string, component templ.Component) error {
	return w.writeFile(filepath.Join(lang, filepath.FromSlash(path), "index.html"), component)
}

// writeFile scrive un componente in un file, creando le cartelle necessarie.
func (w *pageWriter) writeFile(relative string, component templ.Component) error {
	target := filepath.Join(w.outputDir, relative)
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("creazione di %s: %w", filepath.Dir(target), err)
	}
	file, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("creazione di %s: %w", target, err)
	}
	// La build non è legata a una richiesta HTTP: il contesto è nuovo. Il
	// timeout evita che un template in errore blocchi la build per sempre.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := component.Render(ctx, file); err != nil {
		file.Close()
		return fmt.Errorf("rendering di %s: %w", target, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("chiusura di %s: %w", target, err)
	}
	w.pages++
	return nil
}

package build

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/filippogrande/lyrica/internal/render"
)

// writeRaw scrive un file che non è HTML (rss.xml, sitemap.xml, robots.txt):
// writeFile accetta solo componenti templ e va avanti a index.html.
//
// Un errore di scrittura o di creazione cartella è un errore della build, non
// un file saltato: un feed a metà è un feed rotto, e nessuno se ne accorgerebbe.
func (w *pageWriter) writeRaw(relative string, data []byte) error {
	target := filepath.Join(w.outputDir, relative)
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("creazione di %s: %w", filepath.Dir(target), err)
	}
	if err := os.WriteFile(target, data, 0o644); err != nil {
		return fmt.Errorf("scrittura di %s: %w", target, err)
	}
	w.files++
	return nil
}

// siteStats calcola i contatori pubblici del sito dai contenuti (D62, D103).
//
// I numeri sono SEMPRE calcolati, mai scritti a mano: un contatore nel repo può
// solo divergere dalla realtà. I tre valori comprimono, perché l'obiettivo è
// mostrare che il sito cresce e non gonfiare i numeri:
//
//   - Tracks: brani con almeno una traduzione pubblicata (HasTranslations);
//   - Bands: band con almeno un album pubblicato (publishedAlbums);
//   - Langs: le lingue di TRADUZIONE DI ARRIVO (Catalog.TranslationLangs, D69
//     e D96) — ciò che un lettore può trovare, non tutte le lingue originali
//     tradotte.
func siteStats(catalog *content.Catalog) render.SiteStats {
	stats := render.SiteStats{Langs: len(catalog.TranslationLangs())}
	for _, band := range catalog.Bands {
		albums := publishedAlbums(band)
		if len(albums) == 0 {
			continue
		}
		stats.Bands++
		for _, album := range albums {
			for _, track := range album.Tracks {
				if track.HasTranslations() {
					stats.Tracks++
				}
			}
		}
	}
	return stats
}
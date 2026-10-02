package build

import (
	"github.com/filippogrande/lyrica/internal/content"
	"github.com/filippogrande/lyrica/internal/render"
)

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
//
// publishedAlbums e HasTranslations sono le STESSE regole con cui il build
// decide se una pagina esiste: un contatore che usa un criterio diverso
// conterebbe pagine che non ci sono.
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
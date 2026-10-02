package render

import (
	"fmt"

	"github.com/filippogrande/lyrica/internal/i18n"
)

// SiteStats sono i contatori pubblici del sito (D62, D103): il suo valore a
// colpo d'occhio. Li calcola il build dai contenuti — nessun numero scritto a
// mano, quindi nessun numero che possa divergere dalla realtà.
//
// I contatori comprimono: un brano con quattro lingue conta un brano, non
// quattro. L'obiettivo è mostrare che il sito cresce, non gonfiare i numeri.
type SiteStats struct {
	// Tracks sono i brani con almeno una traduzione pubblicata.
	Tracks int
	// Bands sono le band con almeno un album pubblicato.
	Bands int
	// Langs sono le lingue di TRADUZIONE DI ARRIVO (D69, D96): ciò che un
	// lettore può trovare sul sito. Non sono tutte le lingue originali
	// tradotte: il tedesco di un brano bilingue è una versione completa in
	// pagina, non un'interfaccia in tedesco (D96).
	Langs int
}

// TracksText è il contatore dei brani tradotti con la sua etichetta nella
// lingua della pagina.
func (s SiteStats) TracksText(bundle i18n.Bundle, lang string) string {
	return fmt.Sprintf("%d %s", s.Tracks, bundle.MustT(lang, "stats.tracks"))
}

// BandsText è il contatore delle band con la sua etichetta.
func (s SiteStats) BandsText(bundle i18n.Bundle, lang string) string {
	return fmt.Sprintf("%d %s", s.Bands, bundle.MustT(lang, "stats.bands"))
}

// LangsText è il contatore delle lingue con la sua etichetta.
func (s SiteStats) LangsText(bundle i18n.Bundle, lang string) string {
	return fmt.Sprintf("%d %s", s.Langs, bundle.MustT(lang, "stats.langs"))
}

// Line unisce i tre contatori in una riga sola, con il separatore usato nelle
// righe informative del sito.
func (s SiteStats) Line(bundle i18n.Bundle, lang string) string {
	return s.TracksText(bundle, lang) + " · " +
		s.BandsText(bundle, lang) + " · " +
		s.LangsText(bundle, lang)
}
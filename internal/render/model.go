// Package render contiene i componenti Templ del sito e le viste che li
// alimentano.
//
// Separazione: i file .templ contengono SOLO markup e campi già pronti; ogni
// valore (testi composti, URL, conteggi) si calcola in Go nei file view*.go.
// Così i .templ non importano nulla — nemmeno il pacchetto templ, che la
// generazione importa da sé (trappola 1 di ARCHITECTURE.md).
package render

import (
	"github.com/a-h/templ"
	"github.com/filippogrande/lyrica/internal/i18n"
)

// PageData sono i dati comuni a ogni pagina.
type PageData struct {
	Bundle i18n.Bundle
	Lang   string
	Title  string
	// Description è la meta description della pagina.
	Description string
}

// T traduce una chiave dell'interfaccia nella lingua della pagina.
func (d PageData) T(key string) string {
	return d.Bundle.MustT(d.Lang, key)
}

// URL costruisce un URL interno nella lingua corrente.
//
// Restituisce templ.SafeURL perché i file .templ NON importano il pacchetto
// templ: la generazione lo importa da sé, e un import esplicito nei .templ
// produce "templ redeclared in this block".
func (d PageData) URL(suffix string) templ.SafeURL {
	return templ.URL(d.Path(suffix))
}

// Path costruisce il percorso interno nella lingua corrente, come stringa.
func (d PageData) Path(suffix string) string {
	if suffix == "" || suffix == "/" {
		return "/" + d.Lang + "/"
	}
	return "/" + d.Lang + "/" + suffix
}

// BandsPath è il percorso della pagina con l'elenco delle band.
func (d PageData) BandsPath() string { return d.Path("bands/") }

// BandPath è il percorso della pagina di una band.
func (d PageData) BandPath(bandSlug string) string {
	return d.Path("band/" + bandSlug + "/")
}

// AlbumPath è il percorso della pagina di un album.
func (d PageData) AlbumPath(bandSlug, albumSlug string) string {
	return d.Path("band/" + bandSlug + "/album/" + albumSlug + "/")
}

// TrackPath è il percorso della pagina di un brano: è la forma di URL decisa
// nel progetto (/<lang>/band/<band>/album/<album>/brano/<brano>).
func (d PageData) TrackPath(bandSlug, albumSlug, trackSlug string) string {
	return d.Path("band/" + bandSlug + "/album/" + albumSlug + "/brano/" + trackSlug + "/")
}

// Package render contiene i componenti Templ del sito.
package render

import "github.com/filippogrande/lyrica/internal/i18n"

// PageData sono i dati comuni a ogni pagina.
type PageData struct {
	Bundle i18n.Bundle
	Lang   string
	Title  string
}

// T traduce una chiave dell'interfaccia nella lingua della pagina.
func (d PageData) T(key string) string {
	return d.Bundle.MustT(d.Lang, key)
}

// Path costruisce un URL interno nella lingua corrente.
func (d PageData) Path(suffix string) string {
	if suffix == "" || suffix == "/" {
		return "/" + d.Lang + "/"
	}
	return "/" + d.Lang + "/" + suffix
}

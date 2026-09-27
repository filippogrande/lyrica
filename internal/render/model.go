// Package render contiene i componenti Templ del sito.
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
	if suffix == "" || suffix == "/" {
		return templ.URL("/" + d.Lang + "/")
	}
	return templ.URL("/" + d.Lang + "/" + suffix)
}

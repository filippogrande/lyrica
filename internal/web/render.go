package web

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/filippogrande/lyrica/internal/i18n"
	"github.com/filippogrande/lyrica/internal/render"
)

// writePage rende un componente Templ nella risposta con lo status indicato.
// Un errore di rendering non viene nascosto: finisce nel log esplicito
// (DEVELOPMENT_GUIDELINES §5).
func writePage(w http.ResponseWriter, r *http.Request, status int, c templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := c.Render(r.Context(), w); err != nil {
		log.Printf("rendering fallito su %s: %v", r.URL.Path, err)
	}
}

// pageData prepara i dati di pagina: il titolo è una chiave di locale, così
// resta tradotto insieme al resto dell'interfaccia.
func pageData(b i18n.Bundle, lang, titleKey string) render.PageData {
	return render.PageData{Bundle: b, Lang: lang, Title: b.MustT(lang, titleKey)}
}

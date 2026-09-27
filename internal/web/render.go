package web

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/filippogrande/lyrica/internal/i18n"
	"github.com/filippogrande/lyrica/internal/render"
)

// writePage rende un componente Templ nella risposta. Un errore di rendering
// non viene nascosto: finisce nel log (DEVELOPMENT_GUIDELINES §5).
func writePage(w http.ResponseWriter, r *http.Request, c templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := c.Render(r.Context(), w); err != nil {
		log.Printf("rendering fallito su %s: %v", r.URL.Path, err)
	}
}

// d costruisce i dati di pagina per una lingua.
func d(b i18n.Bundle, lang string) render.PageData {
	return render.PageData{Bundle: b, Lang: lang}
}

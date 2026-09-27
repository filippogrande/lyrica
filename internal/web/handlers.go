package web

import (
	"log"
	"net/http"

	"github.com/filippogrande/lyrica/internal/i18n"
	"github.com/filippogrande/lyrica/internal/render"
)

// handleRootRedirect manda la radice alla lingua negoziata dal browser.
func handleRootRedirect(b i18n.Bundle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := b.Negotiate(r.Header.Get("Accept-Language"))
		http.Redirect(w, r, "/"+lang+"/", http.StatusFound)
	}
}

// handleHome rende la home italiana.
func handleHome(b i18n.Bundle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d := render.PageData{Bundle: b, Lang: i18n.DefaultLang, Title: "Lyrica"}
		writePage(w, r, render.Home(d))
	}
}

// handleNotFound rende la 404 personalizzata.
func handleNotFound(b i18n.Bundle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := b.Negotiate(r.Header.Get("Accept-Language"))
		d := render.PageData{Bundle: b, Lang: lang, Title: d(b, lang).Title}
		w.WriteHeader(http.StatusNotFound)
		writePage(w, r, render.NotFound(d))
	}
}

// handleHealth risponde agli healthcheck.
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok\n")); err != nil {
		log.Printf("healthz: scrittura della risposta fallita: %v", err)
	}
}

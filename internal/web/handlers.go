package web

import (
	"log"
	"net/http"

	"github.com/filippogrande/lyrica/internal/i18n"
	"github.com/filippogrande/lyrica/internal/render"
)

// handleRootRedirect manda la radice alla lingua negoziata dal browser (D52).
func handleRootRedirect(b i18n.Bundle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := b.Negotiate(r.Header.Get("Accept-Language"))
		http.Redirect(w, r, "/"+lang+"/", http.StatusFound)
	}
}

// handleHome rende la home nella lingua di default.
func handleHome(b i18n.Bundle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d := pageData(b, i18n.DefaultLang, "home.title")
		writePage(w, r, http.StatusOK, render.Home(d))
	}
}

// handleNotFound rende la 404 personalizzata nella lingua negoziata (D53).
func handleNotFound(b i18n.Bundle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := b.Negotiate(r.Header.Get("Accept-Language"))
		d := pageData(b, lang, "error.not_found_title")
		writePage(w, r, http.StatusNotFound, render.NotFound(d))
	}
}

// handleHealth risponde agli healthcheck del servizio.
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok\n")); err != nil {
		log.Printf("healthz: scrittura della risposta fallita: %v", err)
	}
}

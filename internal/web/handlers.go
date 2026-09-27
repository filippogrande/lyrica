package web

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/filippogrande/lyrica/internal/i18n"
)

// handleRootRedirect manda la radice alla lingua negoziata dal browser (D52).
func handleRootRedirect(b i18n.Bundle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := b.Negotiate(r.Header.Get("Accept-Language"))
		http.Redirect(w, r, "/"+lang+"/", http.StatusFound)
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

// serveNotFound serve la 404 già generata nella lingua negoziata. Se il file
// manca, il 404 di testo finisce nel log: non si inventa una pagina al volo.
func serveNotFound(w http.ResponseWriter, r *http.Request, bundle i18n.Bundle) {
	lang := bundle.Negotiate(r.Header.Get("Accept-Language"))
	target := filepath.Join(publicDir, lang, "404.html")
	body, err := os.ReadFile(target)
	if err != nil {
		log.Printf("404 generata non trovata (%s): %v", target, err)
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	if _, err := w.Write(body); err != nil {
		log.Printf("404: scrittura della risposta fallita: %v", err)
	}
}

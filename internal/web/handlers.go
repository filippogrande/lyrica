package web

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/filippogrande/lyrica/internal/i18n"
)

// LangCookie è il nome del cookie che ricorda la lingua scelta con lo switcher
// (D111): la radice del sito ci riporta nella lingua salvata.
const LangCookie = "lyrica.lang"

// handleRootRedirect manda la radice alla lingua: prima la scelta salvata
// nello switcher (cookie lyrica.lang, D111), poi la negoziazione dal browser
// (D52).
func handleRootRedirect(b i18n.Bundle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := b.Negotiate(r.Header.Get("Accept-Language"))
		if cookie, err := r.Cookie(LangCookie); err == nil {
			lang = b.Negotiate(cookie.Value)
		}
		// #nosec G710 -- lang esce da Negotiate su un set fisso di lingue
		// dell'interfaccia (i file locales/), mai dall'header o dal cookie
		// grezzi.
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

// serveNotFound serve la 404 già generata nella lingua negoziata. Il file si
// legge attraverso os.Root: la lingua viene da Negotiate (set fisso) e Root
// blocca comunque qualsiasi percorso che esca da public/.
func serveNotFound(w http.ResponseWriter, r *http.Request, bundle i18n.Bundle, root *os.Root) {
	lang := bundle.Negotiate(r.Header.Get("Accept-Language"))
	body, err := root.ReadFile(filepath.Join(lang, "404.html"))
	if err != nil {
		// #nosec G706 -- lang esce da Negotiate su set fisso di lingue, mai da
		// input grezzo; err è un errore di os, non testo controllabile.
		log.Printf("404 generata non trovata (%s): %v", lang, err)
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	// #nosec G705 -- il body è il 404.html generato a build time dai template
	// del progetto: nessun contenuto utente arriva qui.
	if _, err := w.Write(body); err != nil {
		log.Printf("404: scrittura della risposta fallita: %v", err)
	}
}

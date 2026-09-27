// Package web serve il sito generato: file statici, salute del servizio e
// (dalle fasi successive) gli endpoint dinamici.
package web

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/filippogrande/lyrica/internal/i18n"
)

// publicDir è la cartella generata da `lyrica build`, dalla radice del progetto.
const publicDir = "public"

// localesDir è la cartella dei locale: serve a negoziare la lingua della 404.
const localesDir = "locales"

// Handler costruisce il router: file generati e salute del servizio.
func Handler() (http.Handler, error) {
	if err := requireBuild(); err != nil {
		return nil, err
	}
	bundle, err := i18n.Load(localesDir)
	if err != nil {
		return nil, err
	}

	files := http.FileServer(http.Dir(publicDir))
	mux := http.NewServeMux()

	// La radice non è una pagina: redirige alla lingua negoziata (D52).
	mux.HandleFunc("/{$}", handleRootRedirect(bundle))

	// Salute del servizio, usata dagli healthcheck.
	mux.HandleFunc("/healthz", handleHealth)

	// Tutto il resto: se il file generato esiste lo si serve, altrimenti 404.
	mux.HandleFunc("/", serveGenerated(files, bundle))

	return securityHeaders(mux), nil
}

// requireBuild si assicura che il sito sia stato generato: servire una cartella
// vuota darebbe 404 su tutto, senza spiegare perché.
func requireBuild() error {
	index := filepath.Join(publicDir, "it", "index.html")
	if _, err := os.Stat(index); err != nil {
		return fmt.Errorf("manca %s: esegui `lyrica build` prima di `lyrica serve` (%w)", index, err)
	}
	return nil
}

// serveGenerated serve il file generato richiesto; se non esiste risponde con
// la 404 generata nella lingua negoziata.
func serveGenerated(files http.Handler, bundle i18n.Bundle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if existsGenerated(r.URL.Path) {
			files.ServeHTTP(w, r)
			return
		}
		serveNotFound(w, r, bundle)
	}
}

// existsGenerated dice se il percorso corrisponde a un file generato o a una
// cartella che contiene index.html.
func existsGenerated(urlPath string) bool {
	clean := path.Clean("/" + urlPath)
	full := filepath.Join(publicDir, filepath.FromSlash(clean))
	info, err := os.Stat(full)
	if err != nil {
		return false
	}
	if !info.IsDir() {
		return true
	}
	_, err = os.Stat(filepath.Join(full, "index.html"))
	return err == nil
}

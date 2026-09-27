package web

import (
	"net/http"

	"github.com/filippogrande/lyrica/internal/i18n"
)

// localesDir è la cartella dei locale, dalla radice del progetto.
const localesDir = "locales"

// assetsDir è la cartella degli asset statici, dalla radice del progetto.
const assetsDir = "assets"

// Handler costruisce il router del sito.
func Handler() (http.Handler, error) {
	bundle, err := i18n.Load(localesDir)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	// La radice non è una pagina: redirige alla lingua negoziata (D52, I18N_DESIGN.md).
	mux.HandleFunc("/{$}", handleRootRedirect(bundle))

	// Home nella lingua italiana: le altre lingue arrivano con i locale (D69).
	mux.HandleFunc("/it/{$}", handleHome(bundle))

	// Salute del servizio, usata dagli healthcheck.
	mux.HandleFunc("/healthz", handleHealth)

	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsDir))))

	// Tutto il resto è 404 personalizzata (D53).
	mux.HandleFunc("/", handleNotFound(bundle))

	return securityHeaders(mux), nil
}

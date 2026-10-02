package build

import (
	"fmt"
	"os"
	"path/filepath"
)

// writeRaw scrive un file che non è HTML (rss.xml, sitemap.xml, robots.txt):
// writeFile accetta solo componenti templ e va avanti a index.html.
//
// Un errore di scrittura o di creazione cartella è un errore della build, non
// un file saltato: un feed a metà è un feed rotto, e nessuno se ne accorgerebbe.
func (w *pageWriter) writeRaw(relative string, data []byte) error {
	target := filepath.Join(w.outputDir, relative)
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("creazione di %s: %w", filepath.Dir(target), err)
	}
	if err := os.WriteFile(target, data, 0o644); err != nil {
		return fmt.Errorf("scrittura di %s: %w", target, err)
	}
	w.files++
	return nil
}
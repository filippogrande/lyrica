// Package build genera il sito statico in public/ a partire dai contenuti.
//
// Ordine (docs/ARCHITECTURE.md): parse → validazione → render → asset. Se la
// validazione fallisce NON si scrive nulla: public/ resta com'era, così un
// contenuto rotto non porta online un sito a metà.
package build

import (
	"fmt"
	"os"
	"strings"

	"github.com/filippogrande/lyrica/internal/content"
	"github.com/filippogrande/lyrica/internal/i18n"
)

// Options indica le cartelle di lavoro, dalla radice del progetto.
type Options struct {
	ContentDir string
	CoversDir  string
	AssetsDir  string
	LocalesDir string
	OutputDir  string
}

// DefaultOptions restituisce le cartelle standard del progetto.
func DefaultOptions() Options {
	return Options{
		ContentDir: "content",
		CoversDir:  "covers",
		AssetsDir:  "assets",
		LocalesDir: "locales",
		OutputDir:  "public",
	}
}

// Run valida i contenuti e genera il sito. La build o riesce per intero o non
// tocca l'output: nessun risultato a metà.
func Run(opts Options) error {
	catalog, err := content.Load(opts.ContentDir)
	if err != nil {
		return err
	}
	report := content.Validate(catalog, content.Options{
		ContentDir: opts.ContentDir,
		CoversDir:  opts.CoversDir,
	})
	for _, issue := range report.Warnings() {
		fmt.Fprintf(os.Stderr, "avviso: %s\n", issue.String())
	}
	if err := report.Err(); err != nil {
		return fmt.Errorf("contenuti non validi:\n%s", report.Error())
	}

	bundle, err := i18n.Load(opts.LocalesDir)
	if err != nil {
		return err
	}
	langs, err := interfaceLangs(bundle, catalog)
	if err != nil {
		return err
	}

	if err := prepareOutput(opts.OutputDir); err != nil {
		return err
	}
	writer := newPageWriter(opts.OutputDir, bundle, catalog)
	for _, lang := range langs {
		if err := writer.writeLang(lang); err != nil {
			return err
		}
	}
	if err := copyAssets(opts.AssetsDir, opts.OutputDir); err != nil {
		return err
	}
	if err := copyCovers(opts.CoversDir, opts.OutputDir); err != nil {
		return err
	}
	fmt.Printf("build: %d pagine in %s (%s)\n", writer.pages, opts.OutputDir, strings.Join(langs, ", "))
	return nil
}

// interfaceLangs sceglie le lingue da generare: quelle che hanno un locale E
// almeno una traduzione nei contenuti (D69), con l'italiano sempre presente.
//
// Una traduzione in una lingua senza il suo file di locale è un errore, non una
// lingua in meno: altrimenti il sito mostrerebbe quella lingua con le stringhe
// di un'altra.
func interfaceLangs(bundle i18n.Bundle, catalog *content.Catalog) ([]string, error) {
	available := map[string]bool{}
	for _, lang := range bundle.Langs() {
		available[lang] = true
	}
	present := map[string]bool{i18n.DefaultLang: true}
	for _, lang := range catalog.TranslationLangs() {
		if !available[lang] {
			return nil, fmt.Errorf("manca locales/%s.yaml: i contenuti hanno traduzioni in %q (D69)", lang, lang)
		}
		present[lang] = true
	}
	var langs []string
	for _, lang := range bundle.Langs() {
		if present[lang] {
			langs = append(langs, lang)
		}
	}
	return langs, nil
}

// prepareOutput svuota e ricrea la cartella di output. Rifiuta i percorsi
// pericolosi: è una cancellazione, non si indovina il bersaglio.
func prepareOutput(dir string) error {
	switch dir {
	case "", ".", "..", "/":
		return fmt.Errorf("cartella di output non valida: %q", dir)
	}
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("pulizia di %s: %w", dir, err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creazione di %s: %w", dir, err)
	}
	return nil
}

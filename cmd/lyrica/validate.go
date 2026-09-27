package main

import (
	"fmt"
	"os"

	"github.com/filippogrande/lyrica/internal/content"
)

// contentDir e coversDir sono le cartelle di lavoro, dalla radice del progetto.
const (
	contentDir = "content"
	coversDir  = "covers"
)

// validateContent legge i contenuti, li valida e stampa il report: prima gli
// errori che bloccano la pubblicazione, poi gli avvisi. Restituisce un errore
// se c'è anche un solo problema bloccante — è questo che rende rossa la PR.
func validateContent() error {
	catalog, err := content.Load(contentDir)
	if err != nil {
		return err
	}
	report := content.Validate(catalog, content.Options{
		ContentDir: contentDir,
		CoversDir:  coversDir,
	})
	for _, issue := range report.Warnings() {
		fmt.Fprintf(os.Stderr, "avviso: %s\n", issue.String())
	}
	if err := report.Err(); err != nil {
		return fmt.Errorf("contenuti non validi:\n%s", report.Error())
	}
	fmt.Printf("contenuti validi: %d band, %d album, %d brani\n",
		len(catalog.Bands), countAlbums(catalog), countTracks(catalog))
	return nil
}

// countAlbums conta gli album pubblicati nel catalogo.
func countAlbums(catalog *content.Catalog) int {
	total := 0
	for _, band := range catalog.Bands {
		total += len(band.Albums)
	}
	return total
}

// countTracks conta i brani pubblicati nel catalogo.
func countTracks(catalog *content.Catalog) int {
	total := 0
	for _, band := range catalog.Bands {
		for _, album := range band.Albums {
			total += len(album.Tracks)
		}
	}
	return total
}

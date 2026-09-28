package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// albumFM modella il front-matter di album.md (docs/CONTENT.md).
type albumFM struct {
	Title  string       `yaml:"title"`
	Slug   string       `yaml:"slug"`
	Year   int          `yaml:"year"`
	Cover  string       `yaml:"cover"`
	Tracks []trackRefFM `yaml:"tracks"`
}

// trackRefFM è una voce della tracklist: slug e titolo, strumentale a parte.
type trackRefFM struct {
	Slug         string `yaml:"slug"`
	Title        string `yaml:"title"`
	Instrumental bool   `yaml:"instrumental"`
}

// newAlbum crea content/bands/<band-slug>/<album-slug>/album.md con la
// tracklist vuota da riempire in ordine.
func newAlbum(args []string) error {
	keys, _, pos, err := parseArgs(args)
	if err != nil {
		return err
	}
	if len(pos) < 2 {
		return fmt.Errorf("uso: lyrica new album <band-slug> <titolo> [--year 1995]")
	}
	bandSlug, title := pos[0], strings.Join(pos[1:], " ")
	bandDir := filepath.Join(contentDir, "bands", bandSlug)
	if _, err := os.Stat(bandDir); os.IsNotExist(err) {
		return fmt.Errorf("band %q inesistente: creala prima con 'lyrica new band'", bandSlug)
	}
	slug := slugify(title)
	slug, err = uniqueSlug(slug, bandDir, "")
	if err != nil {
		return err
	}
	year := 0
	if y := keys["year"]; y != "" {
		year, err = strconv.Atoi(y)
		if err != nil {
		return fmt.Errorf("--year non è un numero: %q", y)
		}
	}
	doc, err := frontMatter(albumFM{Title: title, Slug: slug, Year: year})
	if err != nil {
		return err
	}
	dir := filepath.Join(bandDir, slug)
	if err := os.MkdirAll(filepath.Join(dir, "tracks"), 0o755); err != nil {
		return fmt.Errorf("creazione cartella %s: %w", dir, err)
	}
	path := filepath.Join(dir, "album.md")
	if err := os.WriteFile(path, []byte(doc), 0o644); err != nil {
		return fmt.Errorf("scrittura di %s: %w", path, err)
	}
	fmt.Printf("creato %s\n", path)
	return validateContent()
}

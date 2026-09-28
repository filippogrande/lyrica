package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// bandFM modella il front-matter di band.md (docs/CONTENT.md).
type bandFM struct {
	Name          string   `yaml:"name"`
	Slug          string   `yaml:"slug"`
	Country       string   `yaml:"country"`
	Tags          []string `yaml:"tags"`
	OriginalLangs []string `yaml:"original_langs"`
	FormedYear    int      `yaml:"formed_year"`
	Members       []string `yaml:"members"`
	Image         string   `yaml:"image"`
	Description   string   `yaml:"description"`
}

// newBand crea content/bands/<slug>/band.md con il front-matter precompilato
// (nome, slug, una lingua da --lang). I campi opzionali restano vuoti e la
// validazione finale dice cosa manca ancora (docs/CONTENT.md).
func newBand(args []string) error {
	keys, _, pos, err := parseArgs(args)
	if err != nil {
		return err
	}
	if len(pos) == 0 {
		return fmt.Errorf("uso: lyrica new band <nome> [--lang de]")
	}
	name := strings.Join(pos, " ")
	slug := slugify(name)
	slug, err = uniqueSlug(slug, filepath.Join(contentDir, "bands"), "")
	if err != nil {
		return err
	}
	langs := []string{}
	if l := keys["lang"]; l != "" {
		langs = append(langs, l)
	}
	doc, err := frontMatter(bandFM{Name: name, Slug: slug, OriginalLangs: langs})
	if err != nil {
		return err
	}
	dir := filepath.Join(contentDir, "bands", slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creazione cartella %s: %w", dir, err)
	}
	path := filepath.Join(dir, "band.md")
	if err := os.WriteFile(path, []byte(doc), 0o644); err != nil {
		return fmt.Errorf("scrittura di %s: %w", path, err)
	}
	fmt.Printf("creato %s\n", path)
	return validateContent()
}

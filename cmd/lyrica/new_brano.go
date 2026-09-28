package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// blockFM e stanzaFM modellano il front-matter dei blocchi di testo di un
// brano (docs/CONTENT.md): il tipo di strofa non si stampa mai, restano solo
// lingua, ruolo e righe.
type blockFM struct {
	Lang    string    `yaml:"lang"`
	Role    string    `yaml:"role"`
	Stanzas []stanzaFM `yaml:"stanzas"`
}

type stanzaFM struct {
	Singer string   `yaml:"singer"`
	Lines  []string `yaml:"lines"`
}

// trackFM è il front-matter del file di un brano.
type trackFM struct {
	Title         string    `yaml:"title"`
	Slug          string    `yaml:"slug"`
	AddedDate     string    `yaml:"added_date"`
	Featured      bool      `yaml:"featured"`
	Instrumental  bool      `yaml:"instrumental"`
	OriginalLangs []string  `yaml:"original_langs"`
	Singers       []string  `yaml:"singers"`
	Blocks        []blockFM `yaml:"blocks"`
}

// newBrano crea tracks/<slug>.md con la struttura strofe vuota e il blocco
// originale pronto; poi aggiunge la voce alla tracklist di album.md (regola 8:
// ogni brano va elencato nell'album).
func newBrano(args []string) error {
	keys, bools, pos, err := parseArgs(args)
	if err != nil {
		return err
	}
	if len(pos) < 3 {
		return fmt.Errorf("uso: lyrica new brano <band-slug> <album-slug> <titolo> [--lang de] [--instrumental]")
	}
	bandSlug, albumSlug, title := pos[0], pos[1], strings.Join(pos[2:], " ")
	albumDir := filepath.Join(contentDir, "bands", bandSlug, albumSlug)
	albumPath := filepath.Join(albumDir, "album.md")
	if _, err := os.Stat(albumPath); os.IsNotExist(err) {
		return fmt.Errorf("album %q della band %q inesistente: crealo prima con 'lyrica new album'", albumSlug, bandSlug)
	}
	lang := keys["lang"]
	if lang == "" {
		return fmt.Errorf("serve --lang <codice> per la lingua del blocco (es. --lang de)")
	}
	slug := slugify(title)
	slug, err = uniqueSlug(slug, filepath.Join(albumDir, "tracks"), ".md")
	if err != nil {
		return err
	}
	if err := addToTracklist(albumPath, slug, title, bools["instrumental"]); err != nil {
		return err
	}
	doc, err := frontMatter(trackFM{
		Title:         title,
		Slug:          slug,
		AddedDate:     today(),
		Instrumental:  bools["instrumental"],
		OriginalLangs: []string{lang},
		Blocks:        []blockFM{{Lang: lang, Role: "original"}},
	})
	if err != nil {
		return err
	}
	path := filepath.Join(albumDir, "tracks", slug+".md")
	if err := os.WriteFile(path, []byte(doc), 0o644); err != nil {
		return fmt.Errorf("scrittura di %s: %w", path, err)
	}
	fmt.Printf("creato %s\n", path)
	return validateContent()
}

// addToTracklist legge album.md, aggiunge il brano alla tracklist e riscrive
// il file. È una modifica all'album: la fa solo la CLI, mai a mano.
func addToTracklist(albumPath, slug, title string, instrumental bool) error {
	raw, err := os.ReadFile(albumPath)
	if err != nil {
		return fmt.Errorf("lettura di %s: %w", albumPath, err)
	}
	front, _, err := splitFront(raw)
	if err != nil {
		return err
	}
	var album albumFM
	if err := yaml.Unmarshal(front, &album); err != nil {
		return fmt.Errorf("front-matter di %s non valido: %w", albumPath, err)
	}
	album.Tracks = append(album.Tracks, trackRefFM{Slug: slug, Title: title, Instrumental: instrumental})
	doc, err := frontMatter(album)
	if err != nil {
		return err
	}
	if err := os.WriteFile(albumPath, []byte(doc), 0o644); err != nil {
		return fmt.Errorf("scrittura di %s: %w", albumPath, err)
	}
	return nil
}

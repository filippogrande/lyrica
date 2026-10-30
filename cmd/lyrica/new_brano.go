package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// songFM è il front-matter di song.md: l'anagrafica di un brano nella sua
// cartella (docs/CONTENT.md). I blocchi di testo non stanno qui: ognuno è un
// file a parte.
type songFM struct {
	Title         string   `yaml:"title"`
	Slug          string   `yaml:"slug"`
	AddedDate     string   `yaml:"added_date"`
	Featured      bool     `yaml:"featured"`
	Instrumental  bool     `yaml:"instrumental"`
	OriginalLangs []string `yaml:"original_langs"`
}

// blockFM è il front-matter di un file di blocco (original-<lang>.md o
// <lang>.md): la lingua, il ruolo e le strofe.
type blockFM struct {
	Lang       string     `yaml:"lang"`
	Role       string     `yaml:"role"`
	Translator string     `yaml:"translator,omitempty"`
	Stanzas    []stanzaFM `yaml:"stanzas"`
}

type stanzaFM struct {
	Singer string   `yaml:"singer,omitempty"`
	Lines  []string `yaml:"lines"`
}

// newBrano crea la cartella di un brano — tracks/<slug>/ con song.md e il file
// del blocco originale pronto da riempire — e aggiunge la voce alla tracklist
// di album.md (regola 8: ogni brano va elencato nell'album).
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
	tracksDir := filepath.Join(albumDir, "tracks")
	slug := slugify(title)
	slug, err = uniqueSlug(slug, tracksDir, "")
	if err != nil {
		return err
	}
	if err := addToTracklist(albumPath, slug, title, bools["instrumental"]); err != nil {
		return err
	}
	if err := writeTrackSkeleton(tracksDir, slug, title, lang, bools["instrumental"]); err != nil {
		return err
	}
	fmt.Printf("creato %s\n", filepath.Join(tracksDir, slug))
	return validateContent()
}

// writeTrackSkeleton crea la cartella di un brano: song.md con l'anagrafica e
// il file del blocco originale, pronti da riempire.
func writeTrackSkeleton(tracksDir, slug, title, lang string, instrumental bool) error {
	dir := filepath.Join(tracksDir, slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creazione di %s: %w", dir, err)
	}
	song, err := frontMatter(songFM{
		Title:         title,
		Slug:          slug,
		AddedDate:     today(),
		Instrumental:  instrumental,
		OriginalLangs: []string{lang},
	})
	if err != nil {
		return err
	}
	songPath := filepath.Join(dir, songFileName)
	if err := os.WriteFile(songPath, []byte(song), 0o644); err != nil {
		return fmt.Errorf("scrittura di %s: %w", songPath, err)
	}
	block, err := frontMatter(blockFM{
		Lang:    lang,
		Role:    "original",
		Stanzas: []stanzaFM{{Lines: []string{}}},
	})
	if err != nil {
		return err
	}
	blockPath := filepath.Join(dir, "original-"+lang+".md")
	if err := os.WriteFile(blockPath, []byte(block), 0o644); err != nil {
		return fmt.Errorf("scrittura di %s: %w", blockPath, err)
	}
	return nil
}

// addToTracklist legge album.md, aggiunge il brano alla tracklist e riscrive
// il file. È una modifica all'album: la fa solo la CLI, mai a mano.
func addToTracklist(albumPath, slug, title string, instrumental bool) error {
	// #nosec G304 -- albumPath arriva dagli argomenti della CLI dell'autore
	// (contentDir + slug), mai da input esterni.
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

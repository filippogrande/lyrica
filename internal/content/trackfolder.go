package content

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// La cartella di un brano (docs/CONTENT.md):
//
//	tracks/<slug>/
//	  song.md             anagrafica del brano + note redazionali
//	  original-<lang>.md  un file per ogni blocco role: original
//	  <lang>.md           un file per ogni blocco role: translation
//
// Il nome del file deriva dal contenuto: il prefisso "original-" distingue le
// lingue originali dalle traduzioni, così un brano che ha la stessa lingua come
// originale e come versione completa (D96) non collide (original-de.md e de.md
// sono due file diversi, non un doppione).
//
// Il file di un blocco è autocontenuto: si legge e si valida da solo, e una
// traduzione nuova è un file nuovo, non una modifica a un file che cresce.

// songFileName è il file dell'anagrafica del brano dentro la sua cartella.
const songFileName = "song.md"

// originalPrefix è il prefisso dei file dei blocchi originali.
const originalPrefix = "original-"

// isSongFile dice se il nome è quello dell'anagrafica del brano.
func isSongFile(name string) bool { return name == songFileName }

// blockFileRole ricava lingua e ruolo dal nome di un file di blocco:
// original-<lang>.md è un originale, <lang>.md è una traduzione. Il terzo valore
// dice se il nome è un file di blocco valido.
func blockFileRole(name string) (lang, role string, ok bool) {
	if !strings.HasSuffix(name, ".md") || isSongFile(name) {
		return "", "", false
	}
	stem := strings.TrimSuffix(name, ".md")
	if rest, found := strings.CutPrefix(stem, originalPrefix); found {
		if rest == "" {
			return "", "", false
		}
		return rest, RoleOriginal, true
	}
	if stem == "" {
		return "", "", false
	}
	return stem, RoleTranslation, true
}

// songFrontMatter è l'anagrafica di un brano (song.md): gli stessi campi del
// front-matter di un tempo, senza blocks — nella cartella i blocchi sono file
// separati. I campi sconosciuti restano un errore (frontmatter.go).
type songFrontMatter struct {
	Title         string         `yaml:"title"`
	Slug          string         `yaml:"slug"`
	AddedDate     Date           `yaml:"added_date"`
	Featured      bool           `yaml:"featured"`
	Instrumental  bool           `yaml:"instrumental"`
	OriginalLangs []string       `yaml:"original_langs"`
	Singers       []string       `yaml:"singers"`
	Artists       []ArtistCredit `yaml:"artists"`
}

// apply copia l'anagrafica nel brano.
func (s songFrontMatter) apply(track *Track) {
	track.Title = s.Title
	track.Slug = s.Slug
	track.AddedDate = s.AddedDate
	track.Featured = s.Featured
	track.Instrumental = s.Instrumental
	track.OriginalLangs = s.OriginalLangs
	track.Singers = s.Singers
	track.Artists = s.Artists
}

// blockFileFrontMatter è il front-matter di un file di blocco.
type blockFileFrontMatter struct {
	Lang       string   `yaml:"lang"`
	Role       string   `yaml:"role"`
	Translator string   `yaml:"translator"`
	Live       bool     `yaml:"live"`
	Stanzas    []Stanza `yaml:"stanzas"`
}

// loadTrackFolder legge la cartella di un brano: song.md dà l'anagrafica, ogni
// altro .md è un blocco di testo. Senza song.md è un errore: la cartella non è
// un brano.
func loadTrackFolder(dir string, album *Album) (*Track, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("lettura di %s: %w", dir, err)
	}
	track := &Track{Path: dir, Album: album}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || isIgnorable(entry.Name()) {
			continue
		}
		names = append(names, entry.Name())
	}
	// Ordine stabile: i blocchi si leggono in ordine di nome, così l'elenco
	// delle lingue in pagina non dipende dall'ordine del filesystem.
	sort.Strings(names)
	foundSong := false
	for _, name := range names {
		if isSongFile(name) {
			var fm songFrontMatter
			notes, err := readFrontMatterFile(filepath.Join(dir, name), &fm)
			if err != nil {
				return nil, err
			}
			fm.apply(track)
			track.Notes = notes
			foundSong = true
			continue
		}
		block, err := readBlockFile(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		track.Blocks = append(track.Blocks, block)
	}
	if !foundSong {
		return nil, fmt.Errorf("%s: manca %s, l'anagrafica del brano", dir, songFileName)
	}
	return track, nil
}

// readBlockFile legge un file di blocco e verifica che nome del file e
// front-matter dicano la stessa cosa: un file chiamato it.md con lang: de è un
// errore, non una sorpresa.
func readBlockFile(path string) (Block, error) {
	var fm blockFileFrontMatter
	if _, err := readFrontMatterFile(path, &fm); err != nil {
		return Block{}, err
	}
	nameLang, nameRole, ok := blockFileRole(filepath.Base(path))
	if !ok {
		return Block{}, fmt.Errorf("%s: nome file non valido per un blocco di testo", path)
	}
	if fm.Lang != nameLang {
		return Block{}, fmt.Errorf("%s: lang %q diversa dal nome del file %q", path, fm.Lang, nameLang)
	}
	if fm.Role != nameRole {
		return Block{}, fmt.Errorf("%s: role %q diverso dal nome del file (atteso %q)", path, fm.Role, nameRole)
	}
	return Block{
		Lang:       fm.Lang,
		Role:       fm.Role,
		Translator: fm.Translator,
		Live:       fm.Live,
		Stanzas:    fm.Stanzas,
		Path:       path,
	}, nil
}

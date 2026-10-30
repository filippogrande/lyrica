package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// songFileName è il file dell'anagrafica dentro la cartella di un brano.
const songFileName = "song.md"

// yamlDate tiene una data YAML come testo: accetta la forma nuda (2026-09-30,
// che YAML legge come !!timestamp) e quella fra virgolette, e la riscrive
// com'era. Il loader accetta entrambe (internal/content/date.go).
type yamlDate string

// UnmarshalYAML legge il valore grezzo del nodo, senza convertirlo in data.
func (d *yamlDate) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.ScalarNode {
		return fmt.Errorf("data attesa come YYYY-MM-DD")
	}
	*d = yamlDate(strings.TrimSpace(value.Value))
	return nil
}

// migrateTrackFM è il front-matter di un brano nel formato precedente: un file
// unico con l'anagrafica e i blocchi inline.
type migrateTrackFM struct {
	Title         string           `yaml:"title"`
	Slug          string           `yaml:"slug"`
	AddedDate     yamlDate         `yaml:"added_date"`
	Featured      bool             `yaml:"featured"`
	Instrumental  bool             `yaml:"instrumental"`
	OriginalLangs []string         `yaml:"original_langs"`
	Singers       []string         `yaml:"singers"`
	Artists       []artistCreditFM `yaml:"artists"`
	Blocks        []migrateBlockFM `yaml:"blocks"`
}

// artistCreditFM è un artista del brano nel front-matter.
type artistCreditFM struct {
	Slug string `yaml:"slug"`
	Role string `yaml:"role"`
}

// migrateBlockFM è un blocco di testo nel formato precedente.
type migrateBlockFM struct {
	Lang       string     `yaml:"lang"`
	Role       string     `yaml:"role"`
	Translator string     `yaml:"translator"`
	Live       bool       `yaml:"live"`
	Stanzas    []stanzaFM `yaml:"stanzas"`
}

// songOutFM è il front-matter di song.md: l'anagrafica del brano senza i
// blocchi, che diventano file separati.
type songOutFM struct {
	Title         string           `yaml:"title"`
	Slug          string           `yaml:"slug"`
	AddedDate     yamlDate         `yaml:"added_date"`
	Featured      bool             `yaml:"featured"`
	Instrumental  bool             `yaml:"instrumental"`
	OriginalLangs []string         `yaml:"original_langs"`
	Singers       []string         `yaml:"singers,omitempty"`
	Artists       []artistCreditFM `yaml:"artists,omitempty"`
}

// blockOutFM è il front-matter di un file di blocco.
type blockOutFM struct {
	Lang       string     `yaml:"lang"`
	Role       string     `yaml:"role"`
	Translator string     `yaml:"translator,omitempty"`
	Live       bool       `yaml:"live,omitempty"`
	Stanzas    []stanzaFM `yaml:"stanzas"`
}

// runMigrate trasforma i brani dal formato a file unico (tracks/<slug>.md) alla
// cartella per brano (tracks/<slug>/). Senza --apply è una prova: dice cosa
// farebbe e non tocca niente. Con --band limita a una sola band.
func runMigrate(args []string) error {
	keys, bools, _, err := parseArgs(args)
	if err != nil {
		return err
	}
	apply := bools["apply"]
	files, err := legacyTrackFiles(keys["band"])
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Println("niente da migrare: nessun brano nel formato a file unico")
		return nil
	}
	for _, path := range files {
		rel, err := filepath.Rel(contentDir, path)
		if err != nil {
			rel = path
		}
		if !apply {
			fmt.Printf("prova: migrerei %s\n", rel)
			continue
		}
		if err := migrateTrackFile(path); err != nil {
			return err
		}
		fmt.Printf("migrato %s\n", rel)
	}
	if !apply {
		fmt.Printf("\n%d brani da migrare: rilancia con --apply per scrivere.\n", len(files))
		return nil
	}
	fmt.Printf("\n%d brani migrati.\n", len(files))
	return validateContent()
}

// legacyTrackFiles elenca i file dei brani nel formato precedente, ordinati. Un
// brano già migrato è una cartella, quindi non compare qui.
func legacyTrackFiles(onlyBand string) ([]string, error) {
	bandsDir := filepath.Join(contentDir, "bands")
	bands, err := os.ReadDir(bandsDir)
	if err != nil {
		return nil, fmt.Errorf("lettura di %s: %w", bandsDir, err)
	}
	var files []string
	for _, band := range bands {
		if !band.IsDir() || strings.HasPrefix(band.Name(), ".") {
			continue
		}
		if onlyBand != "" && band.Name() != onlyBand {
			continue
		}
		bandFiles, err := legacyFilesOfBand(filepath.Join(bandsDir, band.Name()))
		if err != nil {
			return nil, err
		}
		files = append(files, bandFiles...)
	}
	sort.Strings(files)
	return files, nil
}

// legacyFilesOfBand raccoglie i file dei brani vecchi di una band, album per
// album.
func legacyFilesOfBand(bandDir string) ([]string, error) {
	albums, err := os.ReadDir(bandDir)
	if err != nil {
		return nil, fmt.Errorf("lettura di %s: %w", bandDir, err)
	}
	var files []string
	for _, album := range albums {
		if !album.IsDir() || strings.HasPrefix(album.Name(), ".") {
			continue
		}
		tracksDir := filepath.Join(bandDir, album.Name(), "tracks")
		entries, err := os.ReadDir(tracksDir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("lettura di %s: %w", tracksDir, err)
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			files = append(files, filepath.Join(tracksDir, entry.Name()))
		}
	}
	return files, nil
}

// migrateTrackFile legge un brano a file unico e scrive la sua cartella:
// song.md con l'anagrafica e le note, un file per ogni blocco. Poi toglie il
// file vecchio. Se qualcosa non torna si ferma prima di scrivere.
func migrateTrackFile(path string) error {
	// #nosec G304 -- path arriva da legacyTrackFiles(), che cammina il repo
	// (content/): nessun input esterno.
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("lettura di %s: %w", path, err)
	}
	front, notes, err := splitFront(raw)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	var fm migrateTrackFM
	if err := yaml.Unmarshal(front, &fm); err != nil {
		return fmt.Errorf("front-matter di %s non valido: %w", path, err)
	}
	if strings.TrimSpace(fm.Slug) == "" {
		return fmt.Errorf("%s: slug mancante, non so come chiamare la cartella", path)
	}
	dir := strings.TrimSuffix(path, ".md")
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("%s: esiste già la cartella %s, la salto", path, dir)
	}
	docs, err := trackFolderDocs(fm, notes)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	// #nosec G301 -- contenuti del repo: 0755 voluto, git e i tool devono
	// attraversare le cartelle.
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creazione di %s: %w", dir, err)
	}
	for name, doc := range docs {
		target := filepath.Join(dir, name)
		// #nosec G306 -- contenuti del repo: 0644 voluto, leggibili da git e CI.
		if err := os.WriteFile(target, []byte(doc), 0o644); err != nil {
			return fmt.Errorf("scrittura di %s: %w", target, err)
		}
	}
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("rimozione di %s: %w", path, err)
	}
	return nil
}

// trackFolderDocs costruisce i file della cartella: song.md più un file per
// blocco, con il nome derivato da ruolo e lingua.
func trackFolderDocs(fm migrateTrackFM, notes string) (map[string]string, error) {
	songDoc, err := frontMatter(songOutFM{
		Title:         fm.Title,
		Slug:          fm.Slug,
		AddedDate:     fm.AddedDate,
		Featured:      fm.Featured,
		Instrumental:  fm.Instrumental,
		OriginalLangs: fm.OriginalLangs,
		Singers:       fm.Singers,
		Artists:       fm.Artists,
	})
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(notes) != "" {
		songDoc += notes + "\n"
	}
	docs := map[string]string{songFileName: songDoc}
	for _, block := range fm.Blocks {
		if strings.TrimSpace(block.Lang) == "" {
			return nil, fmt.Errorf("blocco senza lang: non so come chiamare il file")
		}
		name, err := blockFileName(block.Role, block.Lang)
		if err != nil {
			return nil, err
		}
		if _, clash := docs[name]; clash {
			return nil, fmt.Errorf("due blocchi finirebbero nello stesso file %s", name)
		}
		doc, err := frontMatter(blockOutFM{
			Lang:       block.Lang,
			Role:       block.Role,
			Translator: block.Translator,
			Live:       block.Live,
			Stanzas:    block.Stanzas,
		})
		if err != nil {
			return nil, err
		}
		docs[name] = doc
	}
	return docs, nil
}

// blockFileName è il nome del file di un blocco: original-<lang>.md per un
// originale, <lang>.md per una traduzione (la stessa regola del loader).
func blockFileName(role, lang string) (string, error) {
	switch role {
	case "original":
		return "original-" + lang + ".md", nil
	case "translation":
		return lang + ".md", nil
	default:
		return "", fmt.Errorf("role %q sconosciuto: ammessi original e translation", role)
	}
}

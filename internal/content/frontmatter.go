package content

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// delimitatore del front-matter: una riga composta esattamente da tre trattini.
const frontMatterDelimiter = "---"

// splitFrontMatter separa il front-matter YAML dal corpo del file.
// Il file deve iniziare con una riga "---" e contenerne un'altra per chiudere.
func splitFrontMatter(doc []byte) (front []byte, body string, err error) {
	lines := strings.Split(string(doc), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != frontMatterDelimiter {
		return nil, "", fmt.Errorf("front-matter mancante: il file deve iniziare con una riga %s", frontMatterDelimiter)
	}
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == frontMatterDelimiter {
			front = []byte(strings.Join(lines[1:i], "\n"))
			return front, strings.TrimSpace(strings.Join(lines[i+1:], "\n")), nil
		}
	}
	return nil, "", fmt.Errorf("front-matter non chiuso: manca la riga %s finale", frontMatterDelimiter)
}

// decodeFrontMatter decodifica il front-matter in v rifiutando i campi
// sconosciuti: un refuso in una chiave deve farsi sentire, non essere ignorato.
func decodeFrontMatter(front []byte, v any) error {
	decoder := yaml.NewDecoder(bytes.NewReader(front))
	decoder.KnownFields(true)
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("front-matter YAML non valido: %w", err)
	}
	return nil
}

// readFrontMatterFile legge un file di contenuto, decodifica il front-matter
// in v e restituisce il corpo del file.
func readFrontMatterFile(path string, v any) (string, error) {
	doc, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("lettura di %s: %w", path, err)
	}
	front, body, err := splitFrontMatter(doc)
	if err != nil {
		return "", fmt.Errorf("%s: %w", path, err)
	}
	if err := decodeFrontMatter(front, v); err != nil {
		return "", fmt.Errorf("%s: %w", path, err)
	}
	return body, nil
}

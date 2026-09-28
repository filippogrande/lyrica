// Comando 'new': genera i template di contenuto (band, album, brano) e poi
// lancia la validazione per dire subito cosa manca da compilare (task #18).
package main

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// runNew smista il comando new sul tipo di contenuto.
func runNew(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("uso: lyrica new band|album|brano ... (vedi docs/CONTENT.md)")
	}
	switch args[0] {
	case "band":
		return newBand(args[1:])
	case "album":
		return newAlbum(args[1:])
	case "brano":
		return newBrano(args[1:])
	default:
		return fmt.Errorf("tipo di contenuto %q sconosciuto: ammessi band, album, brano", args[0])
	}
}

// frontMatter serializza un front-matter YAML delimitato da ---.
func frontMatter(v any) (string, error) {
	doc, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return "---\n" + string(doc) + "---\n", nil
}

// parseArgs separa i flag (--chiave valore, --booleani) dagli argomenti
// posizionali. Un flag che pretende un valore ma non lo ha restituisce un
// errore esplicito, mai un fallback silenzioso.
func parseArgs(args []string) (keys map[string]string, booleans map[string]bool, pos []string, err error) {
	keys = map[string]string{}
	booleans = map[string]bool{}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if !strings.HasPrefix(a, "--") {
			pos = append(pos, a)
			continue
		}
		name := strings.TrimPrefix(a, "--")
		if k, v, ok := strings.Cut(name, "="); ok {
			keys[k] = v
			continue
		}
		switch name {
		case "instrumental":
			booleans[name] = true
		case "lang", "year":
			if i+1 >= len(args) {
				return nil, nil, nil, fmt.Errorf("flag --%s richiede un valore", name)
			}
			keys[name] = args[i+1]
			i++
		default:
			return nil, nil, nil, fmt.Errorf("flag sconosciuto --%s", name)
		}
	}
	return keys, booleans, pos, nil
}

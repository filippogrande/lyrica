package content

import (
	"regexp"
	"strings"
)

// linkPattern trova i link interni scritti in Markdown: ](/it/band/...).
var linkPattern = regexp.MustCompile(`\]\((/[^)\s]+)\)`)

// checkLinks verifica la regola 9: un link interno deve puntare a qualcosa che
// esiste davvero. Vale per la descrizione di una band e per le note di un
// brano, gli unici testi liberi dello schema.
func checkLinks(catalog *Catalog, where, text string, report *Report) {
	for _, match := range linkPattern.FindAllStringSubmatch(text, -1) {
		target := match[1]
		if !internalTargetExists(catalog, target) {
			report.Errorf(where, "link interno %q non esiste (regola 9)", target)
		}
	}
}

// internalTargetExists interpreta i percorsi del sito:
//
//	/<lang>/band/<band>/
//	/<lang>/band/<band>/album/<album>/
//	/<lang>/band/<band>/album/<album>/brano/<brano>/
//
// Un percorso scritto in un altro modo è un link rotto: il sito non ha altre
// forme di URL.
func internalTargetExists(catalog *Catalog, target string) bool {
	parts := splitPath(target)
	if len(parts) < 3 || parts[1] != "band" {
		return false
	}
	band, ok := catalog.Band(parts[2])
	if !ok {
		return false
	}
	switch {
	case len(parts) == 3:
		return true
	case len(parts) == 5 && parts[3] == "album":
		_, ok := albumOf(band, parts[4])
		return ok
	case len(parts) == 7 && parts[3] == "album" && parts[5] == "brano":
		if _, ok := albumOf(band, parts[4]); !ok {
			return false
		}
		_, ok := catalog.Track(band.Slug, parts[6])
		return ok
	default:
		return false
	}
}

// splitPath spezza un percorso nei suoi segmenti non vuoti.
func splitPath(target string) []string {
	var parts []string
	for _, part := range strings.Split(target, "/") {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

// albumOf cerca un album per slug dentro una band.
func albumOf(band *Band, slug string) (*Album, bool) {
	for _, album := range band.Albums {
		if album.Slug == slug {
			return album, true
		}
	}
	return nil, false
}

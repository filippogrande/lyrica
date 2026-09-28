package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// slugify trasforma un titolo in slug: minuscolo, ASCII, parole separate da
// trattini (docs/CONTENT.md, "Regole generali"). Gli umlaut tedeschi si
// sciolgono come prevede la traslitterazione a mano (Ü30 -> ue30, härte ->
// haerte): chi li vuole a video li vede scritti in chiaro, lo slug resta ASCII.
var slugNonAscii = regexp.MustCompile(`[^a-z0-9]+`)

var slugUmlaut = strings.NewReplacer(
	"ä", "ae", "ö", "oe", "ü", "ue", "ß", "ss",
)

// slugify costruisce lo slug dal titolo.
func slugify(title string) string {
	s := strings.ToLower(title)
	s = slugUmlaut.Replace(s)
	s = slugNonAscii.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// uniqueSlug restituisce uno slug libero nella cartella data: se base esiste
// già, aggiunge -2, -3, ... finché non trova un nome che non collide. Serve a
// disambiguare slug omonimi (nirvana-us / nirvana-uk, docs/SPEC.md).
func uniqueSlug(base, dir, ext string) (string, error) {
	slug := base
	for n := 2; ; n++ {
		if _, err := os.Stat(filepath.Join(dir, slug+ext)); os.IsNotExist(err) {
			return slug, nil
		} else if err != nil {
			return "", err
		}
		slug = fmt.Sprintf("%s-%d", base, n)
	}
}

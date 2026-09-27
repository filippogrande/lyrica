package content

import "regexp"

// slugPattern descrive gli slug ammessi: minuscoli, ASCII, parole separate da
// trattini (docs/CONTENT.md, "Regole generali"). Niente maiuscole, niente
// diacritici, niente underscore: la traslitterazione si fa a mano una volta
// (Ü30 -> ue30, härte -> haerte).
var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// validSlug dice se uno slug rispetta il formato ammesso.
func validSlug(slug string) bool {
	return slugPattern.MatchString(slug)
}

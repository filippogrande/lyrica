package content

import (
	"path/filepath"
	"strings"
)

// validateTrack controlla un brano: anagrafica, blocchi e traduzioni.
func validateTrack(catalog *Catalog, track *Track, report *Report) {
	where := track.Path
	expected := strings.TrimSuffix(filepath.Base(track.Path), ".md")
	checkSlugAgainst(track.Slug, expected, "file", where, report)
	if strings.TrimSpace(track.Title) == "" {
		report.Errorf(where, "title mancante")
	}
	if track.AddedDate.IsZero() {
		report.Errorf(where, "added_date mancante o non valida (regola 5)")
	}
	checkBlockLangs(track, report)
	checkOriginalLangs(track, report)
	checkArtists(catalog, track, report)
	original, hasOriginal := track.Original()
	if !hasOriginal {
		report.Errorf(where, "nessun blocco con role: original (regola 3)")
	}
	for _, block := range track.Blocks {
		checkStanzaLines(track.Path+" ["+block.Lang+"]", block.Stanzas, report)
	}
	if hasOriginal {
		checkTranslations(track, original, report)
	}
	if !track.HasTranslations() {
		report.Warnf(where, "brano senza traduzioni: va segnato come \"solo originale\" e non linkato")
	}
	checkLinks(catalog, where, track.Notes, report)
}

// checkArtists verifica che tutti gli slug in artists esistano nel catalogo.
func checkArtists(catalog *Catalog, track *Track, report *Report) {
	if len(track.Artists) == 0 {
		report.Warnf(track.Path, "artists mancante: il brano non sarà collegato ad alcuna band oltre a quella dell'album")
		return
	}
	hasPrimary := false
	for _, a := range track.Artists {
		if strings.TrimSpace(a.Slug) == "" {
			report.Errorf(track.Path, "artist slug vuoto")
			continue
		}
		if a.Role == ArtistRolePrimary {
			hasPrimary = true
		}
		// Verifica che la band esista
		found := false
		for _, band := range catalog.Bands {
			if band.Slug == a.Slug {
				found = true
				break
			}
		}
		if !found {
			report.Errorf(track.Path, "artist slug %q non trovato in content/bands/", a.Slug)
		}
	}
	if !hasPrimary {
		report.Warnf(track.Path, "nessun artist con role: primary; il primo della lista farà da primary")
	}
}

// checkBlockLangs verifica la regola 3: role noto, una lingua una volta per
// ruolo, almeno una strofa per blocco.
//
// La stessa lingua può comparire due volte se i ruoli sono diversi (D96): in un
// brano bilingue l'originale è in tedesco e la sua "versione completa" è una
// traduzione in tedesco, che è un blocco legittimo, non un doppione.
func checkBlockLangs(track *Track, report *Report) {
	seen := make(map[string]bool, len(track.Blocks))
	for _, block := range track.Blocks {
		if strings.TrimSpace(block.Lang) == "" {
			report.Errorf(track.Path, "blocco senza lang (regola 3)")
			continue
		}
		key := block.Role + "|" + block.Lang
		if seen[key] {
			report.Errorf(track.Path, "due blocchi %s con lang %q: la stessa lingua compare una volta per ruolo (regola 3)", block.Role, block.Lang)
		}
		seen[key] = true
		switch block.Role {
		case RoleOriginal, RoleTranslation:
		default:
			report.Errorf(track.Path, "blocco %s con role %q sconosciuto: ammessi %q e %q (regola 3)", block.Lang, block.Role, RoleOriginal, RoleTranslation)
		}
		if len(block.Stanzas) == 0 {
			report.Errorf(track.Path, "blocco %s senza strofe (regola 3)", block.Lang)
		}
	}
}

// checkOriginalLangs verifica la regola 7: una lingua dichiarata in
// original_langs ma assente da ogni blocco è un errore.
func checkOriginalLangs(track *Track, report *Report) {
	if len(track.OriginalLangs) == 0 {
		report.Errorf(track.Path, "original_langs mancante: serve almeno una lingua originale")
		return
	}
	for _, lang := range track.OriginalLangs {
		if !hasBlockLang(track, lang) {
			report.Errorf(track.Path, "original_langs dichiara %q ma nessun blocco ha quella lingua (regola 7)", lang)
		}
	}
}

// hasBlockLang dice se un brano ha un blocco in quella lingua.
func hasBlockLang(track *Track, lang string) bool {
	for _, block := range track.Blocks {
		if block.Lang == lang {
			return true
		}
	}
	return false
}

// checkTranslations verifica la regola 4, la più importante: una traduzione
// incompleta non si pubblica. Stesso numero di strofe dell'originale e stessa
// lunghezza per ogni strofa.
func checkTranslations(track *Track, original Block, report *Report) {
	for _, block := range track.Blocks {
		if block.Role != RoleTranslation {
			continue
		}
		where := track.Path + " [" + block.Lang + "]"
		if len(block.Stanzas) != len(original.Stanzas) {
			report.Errorf(where, "traduzione incompleta: %d strofe contro %d dell'originale (regola 4)", len(block.Stanzas), len(original.Stanzas))
			continue
		}
		for index, stanza := range block.Stanzas {
			if expected := len(original.Stanzas[index].Lines); len(stanza.Lines) != expected {
				report.Errorf(where, "strofa %d: %d righe contro %d dell'originale (regola 4)", index+1, len(stanza.Lines), expected)
			}
		}
	}
}

// checkStanzaLines verifica che nessuna riga sia vuota: una riga vuota
// diventa una riga fantasma in pagina (docs/CONTENT.md, campi di una strofa).
func checkStanzaLines(where string, stanzas []Stanza, report *Report) {
	for index, stanza := range stanzas {
		for _, line := range stanza.Lines {
			if strings.TrimSpace(line) == "" {
				report.Errorf(where, "strofa %d: riga vuota", index+1)
			}
		}
	}
}

package content

import (
	"path/filepath"
	"strings"
)

// validateTrack controlla un brano: anagrafica, blocchi e traduzioni.
func validateTrack(catalog *Catalog, track *Track, report *Report) {
	where := track.Path
	checkSlugAgainst(track.Slug, trackName(track), trackKind(track), where, report)
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
		checkStanzaLines(blockWhere(track, block), block.Stanzas, report)
	}
	if hasOriginal {
		checkTranslations(track, original, report)
	}
	if !track.HasTranslations() {
		report.Warnf(where, "brano senza traduzioni: va segnato come \"solo originale\" e non linkato")
	}
	checkLinks(catalog, where, track.Notes, report)
}

// trackName è il nome da cui deve derivare lo slug: la cartella (formato a
// cartella) o il nome del file senza .md (formato precedente).
func trackName(track *Track) string {
	base := filepath.Base(track.Path)
	if track.Legacy {
		return strings.TrimSuffix(base, ".md")
	}
	return base
}

// trackKind è come si chiama, in un messaggio d'errore, la cosa di cui lo slug
// deve portare il nome.
func trackKind(track *Track) string {
	if track.Legacy {
		return "file"
	}
	return "cartella"
}

// blockWhere è il riferimento da stampare per un blocco: il suo file quando
// esiste (formato a cartella), altrimenti il brano con la lingua fra parentesi.
func blockWhere(track *Track, block Block) string {
	if block.Path != "" {
		return block.Path
	}
	return track.Path + " [" + block.Lang + "]"
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
// Nel formato a cartella la stessa coppia non può ripetersi: il nome del file
// (original-<lang>.md, <lang>.md) è la chiave del blocco, quindi due blocchi con
// la stessa coppia scriverebbero lo stesso file. Il controllo resta per il
// formato precedente e per le regole che valgono comunque (role noto, strofe).
func checkBlockLangs(track *Track, report *Report) {
	seen := make(map[string]bool, len(track.Blocks))
	for _, block := range track.Blocks {
		if strings.TrimSpace(block.Lang) == "" {
			report.Errorf(track.Path, "blocco senza lang (regola 3)")
			continue
		}
		key := block.Role + "|" + block.Lang
		if seen[key] {
			report.Errorf(blockWhere(track, block), "due blocchi %s con lang %q: la stessa lingua compare una volta per ruolo (regola 3)", block.Role, block.Lang)
		}
		seen[key] = true
		switch block.Role {
		case RoleOriginal, RoleTranslation:
		default:
			report.Errorf(blockWhere(track, block), "blocco %s con role %q sconosciuto: ammessi %q e %q (regola 3)", block.Lang, block.Role, RoleOriginal, RoleTranslation)
		}
		if len(block.Stanzas) == 0 {
			report.Errorf(blockWhere(track, block), "blocco %s senza strofe (regola 3)", block.Lang)
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
		where := blockWhere(track, block)
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

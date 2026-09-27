package content

import (
	"sort"
)

// Options indica dove stanno i contenuti e le cover.
type Options struct {
	ContentDir string
	CoversDir  string
}

// Validate controlla l'integrità dei contenuti e raccoglie tutti i problemi
// trovati (docs/CONTENT.md, "Regole di validazione"). Le regole 1-9 sono
// errori; l'elenco dei warning è nello stesso doc.
func Validate(catalog *Catalog, opts Options) *Report {
	report := &Report{}
	checkDuplicateBands(catalog, report)
	for _, band := range catalog.Bands {
		validateBand(catalog, band, opts, report)
	}
	checkTags(catalog, report)
	return report
}

// checkDuplicateBands verifica la regola 2 fra band diverse.
func checkDuplicateBands(catalog *Catalog, report *Report) {
	seen := make(map[string]string, len(catalog.Bands))
	for _, band := range catalog.Bands {
		if previous, ok := seen[band.Slug]; ok {
			report.Errorf(band.Directory, "slug %q già usato da %s: due band non possono averlo (regola 2)", band.Slug, previous)
			continue
		}
		seen[band.Slug] = band.Directory
	}
}

// checkTags avvisa sui tag usati una volta sola: quasi sempre sono refusi
// ("metal" scritto dove si intendeva "metalcore").
func checkTags(catalog *Catalog, report *Report) {
	counts := map[string]int{}
	for _, band := range catalog.Bands {
		for _, tag := range band.Tags {
			counts[tag]++
		}
	}
	tags := make([]string, 0, len(counts))
	for tag := range counts {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	for _, tag := range tags {
		if counts[tag] == 1 {
			report.Warnf("tag "+tag, "usato una volta sola: possibile refuso")
		}
	}
}

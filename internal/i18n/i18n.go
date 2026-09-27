// Package i18n gestisce le lingue dell'interfaccia e le stringhe tradotte.
//
// Le lingue disponibili derivano dai file locales/<lang>.yaml. Le lingue di
// interfaccia seguono le lingue delle traduzioni presenti nei contenuti
// (DECISION.md, D69).
package i18n

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// DefaultLang è la lingua di fallback dell'interfaccia.
const DefaultLang = "it"

// Bundle contiene le stringhe di tutte le lingue disponibili.
type Bundle struct {
	langs   []string
	strings map[string]map[string]string
}

// Load legge tutti i file locales/<lang>.yaml dalla cartella indicata.
func Load(dir string) (Bundle, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return Bundle{}, fmt.Errorf("lettura dei locale da %q: %w", dir, err)
	}

	b := Bundle{strings: map[string]map[string]string{}}
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".yaml" {
			continue
		}
		lang := strings.TrimSuffix(e.Name(), ".yaml")
		values, err := loadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return Bundle{}, err
		}
		b.strings[lang] = values
		b.langs = append(b.langs, lang)
	}

	if len(b.langs) == 0 {
		return Bundle{}, fmt.Errorf("nessun locale trovato in %q", dir)
	}
	if _, ok := b.strings[DefaultLang]; !ok {
		return Bundle{}, fmt.Errorf("manca il locale di default %q in %q", DefaultLang, dir)
	}

	sort.Strings(b.langs)
	return b, nil
}

// loadFile legge un singolo file di locale.
func loadFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("lettura di %q: %w", path, err)
	}
	values := map[string]string{}
	if err := yaml.Unmarshal(data, &values); err != nil {
		return nil, fmt.Errorf("parsing di %q: %w", path, err)
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("il locale %q è vuoto", path)
	}
	return values, nil
}

// Langs elenca le lingue di interfaccia disponibili, in ordine alfabetico.
func (b Bundle) Langs() []string { return b.langs }

// Negotiate sceglie la lingua leggendo Accept-Language, con fallback italiano.
// È l'unico fallback ammesso ed è progettato (I18N_DESIGN.md).
func (b Bundle) Negotiate(acceptLanguage string) string {
	for _, part := range strings.Split(acceptLanguage, ",") {
		base := baseLang(part)
		if base == "" {
			continue
		}
		if _, ok := b.strings[base]; ok {
			return base
		}
	}
	return DefaultLang
}

// baseLang estrae il codice lingua da una voce di Accept-Language.
func baseLang(part string) string {
	tag := strings.TrimSpace(strings.Split(part, ";")[0])
	tag = strings.ToLower(tag)
	return strings.Split(tag, "-")[0]
}

// MustT traduce una chiave. Se la chiave manca sia nella lingua richiesta sia
// nel fallback italiano il programma si ferma con un errore visibile, invece di
// mostrare un testo vuoto: nessun fallback silenzioso (DEVELOPMENT_GUIDELINES §5).
func (b Bundle) MustT(lang, key string) string {
	if v, ok := b.strings[lang][key]; ok {
		return v
	}
	if v, ok := b.strings[DefaultLang][key]; ok {
		return v
	}
	panic(fmt.Sprintf("i18n: chiave mancante %q (lingua %q, fallback %q)", key, lang, DefaultLang))
}

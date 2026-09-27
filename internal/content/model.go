// Package content legge i contenuti del sito (band, album, brani) dai file
// Markdown con front-matter YAML — schema in docs/CONTENT.md — e ne costruisce
// i modelli in memoria.
//
// Il pacchetto NON giudica i contenuti: legge e segnala con un errore esplicito
// i problemi di formato (front-matter mancante, YAML invalido, campo
// sconosciuto, data non valida). Le regole di integrità fra contenuti (slug
// duplicati, traduzioni incomplete, link rotti) stanno in validate.go e girano
// in CI.
//
// Nessun fallback silenzioso: se un file non si legge, si restituisce l'errore
// invece di saltarlo o di restituire valori vuoti (docs/GUIDELINES.md §4).
package content

const (
	// RoleOriginal è il blocco del testo originale: obbligatorio in ogni brano.
	RoleOriginal = "original"
	// RoleTranslation è un blocco tradotto: si pubblica solo se completo.
	RoleTranslation = "translation"
)

// Band è l'anagrafica di una band (content/bands/<slug>/band.md).
type Band struct {
	Name          string   `yaml:"name"`
	Slug          string   `yaml:"slug"`
	Country       string   `yaml:"country"`
	Tags          []string `yaml:"tags"`
	OriginalLangs []string `yaml:"original_langs"`
	FormedYear    int      `yaml:"formed_year"`
	Members       []string `yaml:"members"`
	Description   string   `yaml:"description"`

	// Directory è la cartella della band: la popola il loader.
	Directory string `yaml:"-"`
	// Albums sono gli album trovati su disco.
	Albums []*Album `yaml:"-"`
}

// Album è un album con la sua tracklist (content/bands/<band>/<album>/album.md).
type Album struct {
	Title string `yaml:"title"`
	Slug  string `yaml:"slug"`
	Year  int    `yaml:"year"`
	Cover string `yaml:"cover"`
	// Tracklist è l'ordine dichiarato in album.md: è quello che si stampa.
	Tracklist []TrackRef `yaml:"tracks"`

	Directory string   `yaml:"-"`
	Band      *Band    `yaml:"-"`
	Tracks    []*Track `yaml:"-"`
}

// TrackRef è un brano nella tracklist dell'album: slug, titolo e se è
// strumentale (gli strumentali stanno in tracklist ma non hanno pagina).
type TrackRef struct {
	Slug         string `yaml:"slug"`
	Title        string `yaml:"title"`
	Instrumental bool   `yaml:"instrumental"`
}

// Track è un brano con i suoi blocchi di testo (tracks/<slug>.md).
type Track struct {
	Title         string   `yaml:"title"`
	Slug          string   `yaml:"slug"`
	AddedDate     Date     `yaml:"added_date"`
	Featured      bool     `yaml:"featured"`
	Instrumental  bool     `yaml:"instrumental"`
	OriginalLangs []string `yaml:"original_langs"`
	Singers       []string `yaml:"singers"`
	Blocks        []Block  `yaml:"blocks"`

	// Notes è il corpo del file dopo il front-matter: note redazionali.
	// Non si stampa mai nella pagina (docs/CONTENT.md).
	Notes string `yaml:"-"`
	// Path è il percorso del file: lo popola il loader.
	Path string `yaml:"-"`
	// Album è l'album di appartenenza: lo popola il loader.
	Album *Album `yaml:"-"`
}

// Block è un blocco di testo in una lingua: originale o traduzione.
type Block struct {
	Lang       string   `yaml:"lang"`
	Role       string   `yaml:"role"`
	Translator string   `yaml:"translator"`
	Live       bool     `yaml:"live"`
	Stanzas    []Stanza `yaml:"stanzas"`
}

// Stanza è una strofa. Singer, se presente, si stampa sopra la strofa: nei
// brani multi-voce va su ogni strofa (docs/CONTENT.md).
type Stanza struct {
	Singer string   `yaml:"singer"`
	Lines  []string `yaml:"lines"`
}

// HasTranslations dice se il brano ha almeno un blocco di traduzione.
func (t *Track) HasTranslations() bool {
	for _, b := range t.Blocks {
		if b.Role == RoleTranslation {
			return true
		}
	}
	return false
}

// Original restituisce il primo blocco originale del brano e se esiste.
// L'assenza è un errore di contenuto: la segnala il validatore.
func (t *Track) Original() (Block, bool) {
	for _, b := range t.Blocks {
		if b.Role == RoleOriginal {
			return b, true
		}
	}
	return Block{}, false
}

// Translation restituisce il blocco di traduzione nella lingua richiesta.
func (t *Track) Translation(lang string) (Block, bool) {
	for _, b := range t.Blocks {
		if b.Role == RoleTranslation && b.Lang == lang {
			return b, true
		}
	}
	return Block{}, false
}

// TranslationLangs elenca le lingue di traduzione presenti nel brano.
func (t *Track) TranslationLangs() []string {
	var langs []string
	for _, b := range t.Blocks {
		if b.Role == RoleTranslation {
			langs = append(langs, b.Lang)
		}
	}
	return langs
}

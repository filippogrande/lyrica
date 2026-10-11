package content

import (
	"os"
	"path/filepath"
	"testing"
)

// writeTestFile scrive un file di contenuto sotto root, creando le cartelle.
func writeTestFile(t *testing.T, root, rel, body string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("creazione di %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("scrittura di %s: %v", path, err)
	}
}

// onlyTrack restituisce l'unico brano del catalogo, o fa fallire il test.
func onlyTrack(t *testing.T, catalog *Catalog) *Track {
	t.Helper()
	var tracks []*Track
	for _, band := range catalog.Bands {
		for _, album := range band.Albums {
			tracks = append(tracks, album.Tracks...)
		}
	}
	if len(tracks) != 1 {
		t.Fatalf("atteso 1 brano, trovati %d", len(tracks))
	}
	return tracks[0]
}

const testBand = `---
name: "Band"
slug: "band"
country: "IT"
original_langs: ["de"]
`

const testAlbum = `---
title: "Album"
slug: "album"
tracks:
  - slug: "prova"
    title: "Prova"
`

// writeBandScaffold scrive band.md e album.md, il minimo perché Load arrivi al
// brano.
func writeBandScaffold(t *testing.T, root string) {
	t.Helper()
	writeTestFile(t, root, "bands/band/band.md", testBand)
	writeTestFile(t, root, "bands/band/album/album.md", testAlbum)
}

// Il formato a cartella: song.md dà l'anagrafica, un file per blocco.
func TestLoadCartellaBrano(t *testing.T) {
	root := t.TempDir()
	writeBandScaffold(t, root)
	writeTestFile(t, root, "bands/band/album/tracks/prova/song.md", `---
title: "Prova"
slug: "prova"
added_date: 2026-01-02
original_langs: ["de"]
`)
	writeTestFile(t, root, "bands/band/album/tracks/prova/original-de.md", `---
lang: de
role: original
stanzas:
  - lines: ["Zeile"]
`)
	writeTestFile(t, root, "bands/band/album/tracks/prova/it.md", `---
lang: it
role: translation
translator: "Filippo"
stanzas:
  - lines: ["Riga"]
`)
	catalog, err := Load(root)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	track := onlyTrack(t, catalog)
	if track.Legacy {
		t.Fatal("un brano in cartella non è nel formato precedente")
	}
	if len(track.Blocks) != 2 {
		t.Fatalf("attesi 2 blocchi, trovati %d", len(track.Blocks))
	}
	original, ok := track.Original()
	if !ok || original.Lang != "de" || original.Role != RoleOriginal {
		t.Fatalf("originale atteso de/original, trovato %+v", original)
	}
	if _, ok := track.Translation("it"); !ok {
		t.Fatal("traduzione it non trovata")
	}
	if langs := catalog.TranslationLangs(); len(langs) != 1 || langs[0] != "it" {
		t.Fatalf("attese lingue [it], trovate %v", langs)
	}
}

// Il formato precedente resta leggibile finché la migrazione non è finita.
func TestLoadFormatoPrecedente(t *testing.T) {
	root := t.TempDir()
	writeBandScaffold(t, root)
	writeTestFile(t, root, "bands/band/album/tracks/prova.md", `---
title: "Prova"
slug: "prova"
added_date: 2026-01-02
original_langs: ["de"]
blocks:
  - lang: de
    role: original
    stanzas:
      - lines: ["Zeile"]
  - lang: it
    role: translation
    stanzas:
      - lines: ["Riga"]
`)
	catalog, err := Load(root)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	track := onlyTrack(t, catalog)
	if !track.Legacy {
		t.Fatal("un brano a file unico è nel formato precedente")
	}
	if len(track.Blocks) != 2 {
		t.Fatalf("attesi 2 blocchi, trovati %d", len(track.Blocks))
	}
}

// Una cartella di brano senza song.md non è un brano.
func TestLoadCartellaSenzaSong(t *testing.T) {
	root := t.TempDir()
	writeBandScaffold(t, root)
	writeTestFile(t, root, "bands/band/album/tracks/prova/it.md", `---
lang: it
role: translation
stanzas:
  - lines: ["Riga"]
`)
	if _, err := Load(root); err == nil {
		t.Fatal("atteso un errore: la cartella del brano non ha song.md")
	}
}

// Lingua dichiarata nel front-matter e diversa dal nome del file: errore.
func TestLoadBloccoNomeIncoerente(t *testing.T) {
	root := t.TempDir()
	writeBandScaffold(t, root)
	writeTestFile(t, root, "bands/band/album/tracks/prova/song.md", `---
title: "Prova"
slug: "prova"
added_date: 2026-01-02
original_langs: ["fr"]
`)
	writeTestFile(t, root, "bands/band/album/tracks/prova/original-de.md", `---
lang: fr
role: original
stanzas:
  - lines: ["Zeile"]
`)
	if _, err := Load(root); err == nil {
		t.Fatal("atteso un errore: lang del front-matter diversa dal nome del file")
	}
}

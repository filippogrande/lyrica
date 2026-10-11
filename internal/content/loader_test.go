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

const testAlbum = `---
title: "Album"
slug: "album"
tracks:
  - slug: "prova"
    title: "Prova"
`

// Il formato a cartella: song.md dà l'anagrafica, un file per blocco.
func TestLoadCartellaBrano(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "bands/band/album/album.md", testAlbum)
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
	writeTestFile(t, root, "bands/band/album/album.md", testAlbum)
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
	writeTestFile(t, root, "bands/band/album/album.md", testAlbum)
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
	writeTestFile(t, root, "bands/band/album/album.md", testAlbum)
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
}ttese lingue [it], trovate %v\", langs)\n\t}\n}\n\n// Il formato precedente resta leggibile finché la migrazione non è finita.\nfunc TestLoadFormatoPrecedente(t *testing.T) {\n\troot := t.TempDir()\n\twriteTestFile(t, root, \"bands/band/album/album.md\", testAlbum)\n\twriteTestFile(t, root, \"bands/band/album/tracks/prova.md\", `---\ntitle: \"Prova\"\nslug: \"prova\"\nadded_date: 2026-01-02\noriginal_langs: [\"de\"]\nblocks:\n  - lang: de\n    role: original\n    stanzas:\n      - lines: [\"Zeile\"]\n  - lang: it\n    role: translation\n    stanzas:\n      - lines: [\"Riga\"]\n`)\n\tcatalog, err := Load(root)\n\tif err != nil {\n\t\tt.Fatalf(\"Load: %v\", err)\n\t}\n\ttrack := onlyTrack(t, catalog)\n\tif !track.Legacy {\n\t\tt.Fatal(\"un brano a file unico è nel formato precedente\")\n\t}\n\tif len(track.Blocks) != 2 {\n\t\tt.Fatalf(\"attesi 2 blocchi, trovati %d\", len(track.Blocks))\n\t}\n}\n\n// Una cartella di brano senza song.md non è un brano.\nfunc TestLoadCartellaSenzaSong(t *testing.T) {\n\troot := t.TempDir()\n\twriteTestFile(t, root, \"bands/band/album/album.md\", testAlbum)\n\twriteTestFile(t, root, \"bands/band/album/tracks/prova/it.md\", `---\nlang: it\nrole: translation\nstanzas:\n  - lines: [\"Riga\"]\n`)\n\tif _, err := Load(root); err == nil {\n\t\tt.Fatal(\"atteso un errore: la cartella del brano non ha song.md\")\n\t}\n}\n\n// Lingua dichiarata nel front-matter e diversa dal nome del file: errore.\nfunc TestLoadBloccoNomeIncoerente(t *testing.T) {\n\troot := t.TempDir()\n\twriteTestFile(t, root, \"bands/band/album/album.md\", testAlbum)\n\twriteTestFile(t, root, \"bands/band/album/tracks/prova/song.md\", `---\ntitle: \"Prova\"\nslug: \"prova\"\nadded_date: 2026-01-02\noriginal_langs: [\"fr\"]\n`)\n\twriteTestFile(t, root, \"bands/band/album/tracks/prova/original-de.md\", `---\nlang: fr\nrole: original\nstanzas:\n  - lines: [\"Zeile\"]\n`)\n\tif _, err := Load(root); err == nil {\n\t\tt.Fatal(\"atteso un errore: lang del front-matter diversa dal nome del file\")\n\t}\n}\n"}]}}
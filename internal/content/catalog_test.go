package content

import "testing"

// testCatalog monta un catalogo minimo con i brani passati: le lingue che il
// catalogo dichiara si provano senza leggere niente da disco.
func testCatalog(tracks ...*Track) *Catalog {
	album := &Album{Slug: "album", Tracks: tracks}
	band := &Band{Slug: "band", Albums: []*Album{album}}
	return &Catalog{Bands: []*Band{band}}
}

func TestTranslationLangsElencaLeLingueDiArrivo(t *testing.T) {
	track := &Track{
		Slug:          "prova",
		OriginalLangs: []string{"de"},
		Blocks: []Block{
			testBlock("de", RoleOriginal),
			testBlock("it", RoleTranslation),
			testBlock("en", RoleTranslation),
		},
	}
	got := testCatalog(track).TranslationLangs()
	if len(got) != 2 || got[0] != "it" || got[1] != "en" {
		t.Fatalf("attese [it en], trovate %v", got)
	}
}

// Una traduzione verso una lingua che il brano ha già come originale non apre
// una lingua d'interfaccia (D96, confermata in D97): la build non deve chiedere
// un locales/<lang>.yaml per quella lingua.
func TestTranslationLangsEscludeLeLingueOriginali(t *testing.T) {
	track := &Track{
		Slug:          "hurrikan",
		OriginalLangs: []string{"de", "en"},
		Blocks: []Block{
			testBlock("de", RoleOriginal),
			testBlock("de", RoleTranslation),
			testBlock("en", RoleTranslation),
			testBlock("it", RoleTranslation),
		},
	}
	got := testCatalog(track).TranslationLangs()
	if len(got) != 1 || got[0] != "it" {
		t.Fatalf("attesa [it], trovate %v", got)
	}
}

func TestTranslationLangsNonRipeteUnaLingua(t *testing.T) {
	first := &Track{
		Slug:          "uno",
		OriginalLangs: []string{"de"},
		Blocks:        []Block{testBlock("de", RoleOriginal), testBlock("it", RoleTranslation)},
	}
	second := &Track{
		Slug:          "due",
		OriginalLangs: []string{"de"},
		Blocks:        []Block{testBlock("de", RoleOriginal), testBlock("it", RoleTranslation)},
	}
	got := testCatalog(first, second).TranslationLangs()
	if len(got) != 1 || got[0] != "it" {
		t.Fatalf("attesa [it] una volta sola, trovate %v", got)
	}
}

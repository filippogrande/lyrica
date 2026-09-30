package content

import "testing"

// testBlock è un blocco minimo valido per la regola 3: una strofa con una riga.
func testBlock(lang, role string) Block {
	return Block{Lang: lang, Role: role, Stanzas: []Stanza{{Lines: []string{"riga"}}}}
}

// checkRegola3 fa girare il controllo dei blocchi su un brano che non esiste su
// disco: la regola 3 guarda solo i blocchi, non il file.
func checkRegola3(blocks []Block) *Report {
	report := &Report{}
	checkBlockLangs(&Track{Path: "tracks/prova.md", Blocks: blocks}, report)
	return report
}

func TestCheckBlockLangsAccettaUnaLinguaPerRuolo(t *testing.T) {
	cases := []struct {
		name   string
		blocks []Block
	}{
		{
			name:   "monolingue",
			blocks: []Block{testBlock("de", RoleOriginal), testBlock("it", RoleTranslation)},
		},
		{
			name: "bilingue",
			blocks: []Block{
				testBlock("de", RoleOriginal),
				testBlock("en", RoleTranslation),
				testBlock("it", RoleTranslation),
			},
		},
		{
			name: "originale e versione completa nella stessa lingua",
			blocks: []Block{
				testBlock("de", RoleOriginal),
				testBlock("de", RoleTranslation),
				testBlock("it", RoleTranslation),
			},
		},
		{
			name: "tre lingue originali con la versione completa in tutte e tre",
			blocks: []Block{
				testBlock("it", RoleOriginal),
				testBlock("it", RoleTranslation),
				testBlock("fr", RoleTranslation),
				testBlock("de", RoleTranslation),
				testBlock("en", RoleTranslation),
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if report := checkRegola3(tc.blocks); report.HasErrors() {
				t.Fatalf("attesi zero errori, trovati: %s", report.Error())
			}
		})
	}
}

func TestCheckBlockLangsRifiutaLaStessaCoppiaRuoloLingua(t *testing.T) {
	cases := []struct {
		name   string
		blocks []Block
	}{
		{
			name:   "due blocchi original con la stessa lingua",
			blocks: []Block{testBlock("de", RoleOriginal), testBlock("de", RoleOriginal)},
		},
		{
			name: "due traduzioni nella stessa lingua",
			blocks: []Block{
				testBlock("de", RoleOriginal),
				testBlock("it", RoleTranslation),
				testBlock("it", RoleTranslation),
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if report := checkRegola3(tc.blocks); !report.HasErrors() {
				t.Fatal("atteso un errore: la stessa coppia role + lang non può comparire due volte")
			}
		})
	}
}

func TestCheckBlockLangsSegnalaBlocchiRotti(t *testing.T) {
	cases := []struct {
		name   string
		blocks []Block
	}{
		{
			name:   "blocco senza lingua",
			blocks: []Block{{Role: RoleOriginal, Stanzas: []Stanza{{Lines: []string{"riga"}}}}},
		},
		{
			name:   "role sconosciuto",
			blocks: []Block{testBlock("de", "coro")},
		},
		{
			name:   "blocco senza strofe",
			blocks: []Block{{Lang: "de", Role: RoleOriginal}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if report := checkRegola3(tc.blocks); !report.HasErrors() {
				t.Fatal("atteso un errore")
			}
		})
	}
}

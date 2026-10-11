package content

import "testing"

// Il nome del file di un blocco dice lingua e ruolo, e distingue l'originale
// dalla traduzione della stessa lingua (D96).
func TestBlockFileRole(t *testing.T) {
	cases := []struct {
		name string
		lang string
		role string
		ok   bool
	}{
		{"original-de.md", "de", RoleOriginal, true},
		{"original-ja-Jpan.md", "ja-Jpan", RoleOriginal, true},
		{"it.md", "it", RoleTranslation, true},
		{"ja-Latn.md", "ja-Latn", RoleTranslation, true},
		{"song.md", "", "", false},
		{"original-.md", "", "", false},
		{"original", "", "", false},
		{"README.txt", "", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lang, role, ok := blockFileRole(tc.name)
			if ok != tc.ok || lang != tc.lang || role != tc.role {
				t.Fatalf("blockFileRole(%q) = (%q, %q, %t), atteso (%q, %q, %t)",
					tc.name, lang, role, ok, tc.lang, tc.role, tc.ok)
			}
		})
	}
}

// La stessa lingua su ruoli diversi è lecita: original it e translation it
// convivono nel brano e nei file (D96, D97).
func TestBlockFileRoleDistingueIruoli(t *testing.T) {
	langOriginal, roleOriginal, ok := blockFileRole("original-it.md")
	if !ok || langOriginal != "it" || roleOriginal != RoleOriginal {
		t.Fatalf("original-it.md non letto come originale it")
	}
	langTranslation, roleTranslation, ok := blockFileRole("it.md")
	if !ok || langTranslation != "it" || roleTranslation != RoleTranslation {
		t.Fatalf("it.md non letto come traduzione it")
	}
}

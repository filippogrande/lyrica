package render

import (
	"testing"

	"github.com/filippogrande/lyrica/internal/content"
)

func testTextBlock(lang, role string) content.Block {
	return content.Block{Lang: lang, Role: role}
}

// L'anagrafica elenca una voce per blocco: in un brano bilingue il tedesco
// compare due volte, come originale e come versione completa (D96, D97).
func TestTrackLanguagesContaUnaVocePerBlocco(t *testing.T) {
	originals := []content.Block{testTextBlock("de", content.RoleOriginal)}
	translations := []content.Block{
		testTextBlock("de", content.RoleTranslation),
		testTextBlock("it", content.RoleTranslation),
	}
	got := trackLanguages(originals, translations)
	want := "🇩🇪 DE · 🇩🇪 DE · 🇮🇹 IT"
	if got != want {
		t.Fatalf("atteso %q, trovato %q", want, got)
	}
}

func TestTrackLanguagesSenzaTraduzioneNellaStessaLingua(t *testing.T) {
	originals := []content.Block{testTextBlock("de", content.RoleOriginal)}
	translations := []content.Block{testTextBlock("it", content.RoleTranslation)}
	got := trackLanguages(originals, translations)
	want := "🇩🇪 DE · 🇮🇹 IT"
	if got != want {
		t.Fatalf("atteso %q, trovato %q", want, got)
	}
}

// LangLabels resta l'elenco delle lingue senza ripetizioni: è quello che usano
// l'anagrafica dell'album e la tracklist.
func TestLangLabelsNonRipeteLeLingue(t *testing.T) {
	got := LangLabels([]string{"it", "de", "it"})
	want := "🇩🇪 DE · 🇮🇹 IT"
	if got != want {
		t.Fatalf("atteso %q, trovato %q", want, got)
	}
}

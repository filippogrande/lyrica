package render

import (
	"sort"
	"strings"
)

// Le bandiere sono l'etichetta delle lingue in tutto il sito (D92): una
// bandiera si legge prima del codice lingua, quindi "🇩🇪 DE" dice più di "DE".
// Il codice resta accanto alla bandiera perché la lingua sia riconoscibile
// anche a chi non distingue le bandiere e a chi le legge con uno screen reader.
//
// L'inglese è 🇬🇧 (Regno Unito): è la scelta dichiarata dall'autore.
var langFlags = map[string]string{
	"ar": "🇸🇦",
	"bg": "🇧🇬",
	"bs": "🇧🇦",
	"cs": "🇨🇿",
	"da": "🇩🇰",
	"de": "🇩🇪",
	"el": "🇬🇷",
	"en": "🇬🇧",
	"es": "🇪🇸",
	"et": "🇪🇪",
	"fa": "🇮🇷",
	"fi": "🇫🇮",
	"fr": "🇫🇷",
	"he": "🇮🇱",
	"hi": "🇮🇳",
	"hr": "🇭🇷",
	"hu": "🇭🇺",
	"is": "🇮🇸",
	"it": "🇮🇹",
	"ja": "🇯🇵",
	"ko": "🇰🇷",
	"lt": "🇱🇹",
	"lv": "🇱🇻",
	"nl": "🇳🇱",
	"no": "🇳🇴",
	"pl": "🇵🇱",
	"pt": "🇵🇹",
	"ro": "🇷🇴",
	"ru": "🇷🇺",
	"sk": "🇸🇰",
	"sl": "🇸🇮",
	"sr": "🇷🇸",
	"sv": "🇸🇪",
	"tr": "🇹🇷",
	"uk": "🇺🇦",
	"zh": "🇨🇳",
}

// FlagOf restituisce la bandiera di una lingua; per una lingua senza bandiera
// (o un codice non previsto) torna il globo: meglio un segno dichiarato che
// nessuna etichetta.
func FlagOf(code string) string {
	if flag, ok := langFlags[strings.ToLower(strings.TrimSpace(code))]; ok {
		return flag
	}
	return "🌐"
}

// LangLabel è la lingua come si legge in pagina: bandiera e codice maiuscolo
// ("🇩🇪 DE"). Il codice non si traduce: è lo stesso in tutte le lingue.
// Una lingua non dichiarata torna vuota: un'etichetta senza codice non dice
// niente.
func LangLabel(code string) string {
	trimmed := strings.ToUpper(strings.TrimSpace(code))
	if trimmed == "" {
		return ""
	}
	return FlagOf(code) + " " + trimmed
}

// LangLabels unisce le etichette di più lingue con il separatore delle liste
// brevi ("🇩🇪 DE · 🇮🇹 IT"): i codici ripetuti contano una volta e l'ordine è
// alfabetico, così la stessa lista si legge uguale in ogni pagina del sito.
// Per l'elenco delle voci di un brano, dove la stessa lingua può valere due
// volte, si usa LangLabelsRepeated.
func LangLabels(codes []string) string {
	return joinLabels(sortedCodes(codes, true))
}

// LangLabelsRepeated unisce le etichette di più lingue tenendo i codici
// ripetuti, in ordine alfabetico: è l'elenco delle VOCI di un brano, dove la
// stessa lingua vale due volte se è sia l'originale sia la sua versione
// completa (D97). Dove una lingua è una voce sola si usa LangLabels.
func LangLabelsRepeated(codes []string) string {
	return joinLabels(sortedCodes(codes, false))
}

// sortedCodes normalizza i codici (minuscolo, senza spazi) e li ordina
// alfabeticamente; con unique i codici ripetuti contano una volta sola.
func sortedCodes(codes []string, unique bool) []string {
	seen := make(map[string]bool, len(codes))
	sorted := make([]string, 0, len(codes))
	for _, code := range codes {
		normalized := strings.ToLower(strings.TrimSpace(code))
		if normalized == "" {
			continue
		}
		if unique && seen[normalized] {
			continue
		}
		seen[normalized] = true
		sorted = append(sorted, normalized)
	}
	sort.Strings(sorted)
	return sorted
}

// joinLabels rende le etichette delle lingue unite dal separatore delle liste
// brevi.
func joinLabels(codes []string) string {
	labels := make([]string, 0, len(codes))
	for _, code := range codes {
		labels = append(labels, LangLabel(code))
	}
	return strings.Join(labels, " · ")
}

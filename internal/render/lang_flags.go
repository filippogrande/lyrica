package render

import (
	"sort"
	"strings"
)

// Le bandiere sono l'etichetta delle lingue in tutto il sito (D92): una
// bandiera si legge prima del codice lingua, quindi "🇩🇪 DE" dice pi� di "DE".
// Il codice resta accanto alla bandiera perch� la lingua sia riconoscibile
// anche a chi non distingue le bandiere e a chi le legge con uno screen reader.
//
// L'inglese � �� (Regno Unito): � la scelta dichiarata dall'autore.
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

// Script labels — nome descrittivo della variante di scrittura, mostrato
// dopo il codice lingua in etichetta quando il BCP-47 ha un subtag di script
// (D99). Un sottotag non riconosciuto mostra il codice completo.
var scriptLabels = map[string]string{
	"jpan": "Kanji",	  // Giapponese in caratteri nativi
	"latn": "Romaji",	 // Giapponese/Cinese in alfabeto latino
	"hans": "Semplificato", // Cinese semplificato (zh-Hans)
	"hant": "Tradizionale", // Cinese tradizionale (zh-Hant)
	"cyrl": "Cirillico",  // Serbo in alfabeto cirillico (sr-Cyrl)
	"kore": "Hangul",	 // Coreano in Hangul (ko-Kore)
	"hira": "Hiragana",   // Giapponese hiragana (raro)
	"kana": "Katakana",   // Giapponese katakana (raro)
}

// FlagOf restituisce la bandiera di una lingua; per un codice con sottotag
// di scrittura (es. ja-Latn) estrae il prefisso prima di "-". Per una lingua
// senza bandiera (o un codice non previsto) torna il globo: meglio un segno
// dichiarato che nessuna etichetta.
func FlagOf(code string) string {
	base := strings.Split(strings.ToLower(strings.TrimSpace(code)), "-")[0]
	if flag, ok := langFlags[base]; ok {
		return flag
	}
	return "🌐"
}

// LangLabel � la lingua come si legge in pagina: bandiera e codice maiuscolo
// ("🇩🇪 DE"). Se il codice BCP-47 ha un subtag di scrittura riconosciuto, lo
// mostra dopo il codice ("🇯🇵 JA Romaji"). Un subtag sconosciuto resta
// visibile come codice completo ("🇯🇵 JA-UNKN").
func LangLabel(code string) string {
	normalized := strings.ToLower(strings.TrimSpace(code))
	if normalized == "" {
		return ""
	}
	parts := strings.Split(normalized, "-")
	base := parts[0]
	label := FlagOf(base) + " " + strings.ToUpper(base)
	if len(parts) > 1 {
		if scriptLabel, ok := scriptLabels[parts[1]]; ok {
			label += " " + scriptLabel
		} else {
			// subtag sconosciuto: mostra il codice completo
			label = FlagOf(base) + " " + strings.ToUpper(base) + "-" + strings.ToUpper(parts[1])
		}
	}
	return label
}

// LangLabels unisce le etichette di pi� lingue con il separatore delle liste
// brevi ("🇩🇪 DE · 🇮🇹 IT"): i codici ripetuti contano una volta e l'ordine �
// alfabetico, cos� la stessa lista si legge uguale in ogni pagina del sito.
// Per l'elenco delle voci di un brano, dove la stessa lingua pu� valere due
// volte, si usa LangLabelsRepeated.
func LangLabels(codes []string) string {
	return joinLabels(sortedCodes(codes, true))
}

// LangLabelsRepeated unisce le etichette di pi� lingue tenendo i codici
// ripetuti, in ordine alfabetico: � l'elenco delle VOCI di un brano, dove la
// stessa lingua vale due volte se � sia l'originale sia la sua versione
// completa (D97). Dove una lingua � una voce sola si usa LangLabels.
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

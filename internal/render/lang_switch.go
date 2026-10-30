package render

import (
	"github.com/a-h/templ"
	"github.com/filippogrande/lyrica/internal/i18n"
)

// LangSwitchData è il link dello switcher di lingua: la stessa pagina
// nell'altra lingua dell'interfaccia.
type LangSwitchData struct {
	Code string
	URL  templ.SafeURL
}

// langSwitchData calcola il link dello switcher (D111): l'altra lingua
// dell'interfaccia (d.Langs), stessa pagina (stesso Self). Sulla 404
// (NoAlternates) la pagina equivalente nell'altra lingua non esiste: si linka
// la radice di quella lingua.
func langSwitchData(d PageData) LangSwitchData {
	other := ""
	for _, lang := range d.Langs {
		if lang != d.Lang {
			other = lang
			break
		}
	}
	if other == "" {
		other = i18n.DefaultLang
	}
	self := d.Self
	if d.NoAlternates {
		self = ""
	}
	return LangSwitchData{Code: other, URL: templ.URL(langPath(other, self))}
}

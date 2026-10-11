// Package render contiene i componenti Templ del sito e le viste che li
// alimentano.
//
// Separazione: i file .templ contengono SOLO markup e campi già pronti; ogni
// valore (testi composti, URL, conteggi) si calcola in Go nei file view*.go.
// Così i .templ non importano nulla — nemmeno il pacchetto templ, che la
// generazione importa da sé (trappola 1 di ARCHITECTURE.md).
package render

import (
	"strings"

	"github.com/a-h/templ"
	"github.com/filippogrande/lyrica/internal/i18n"
)

// PageData sono i dati comuni a ogni pagina.
type PageData struct {
	Bundle i18n.Bundle
	Lang   string
	Title  string
	// Description è la meta description della pagina.
	Description string
	// Base è il dominio assoluto del sito, definito in un posto solo
	// (internal/build): serve ai link che devono uscire dal sito (feed
	// RSS, hreflang). Vuoto solo dove l'indirizzo assoluto non serve.
	Base string
	// Langs sono le lingue dell'interfaccia generate dal build (D69): sono le
	// sole in cui la pagina esiste, quindi le uniche hreflang.
	Langs []string
	// Self è il percorso della pagina dentro la lingua, con la barra finale
	// ("" per la home): serve a ricostruire la pagina nelle altre lingue.
	Self string
	// NoAlternates esclude la pagina dai link hreflang: la 404 non ha una
	// versione per lingua, e un alternato che porta alla home mentirebbe.
	NoAlternates bool
	// Stats sono i contatori pubblici del sito, calcolati a build time.
	Stats SiteStats
}

// T traduce una chiave dell'interfaccia nella lingua della pagina.
func (d PageData) T(key string) string {
	return d.Bundle.MustT(d.Lang, key)
}

// URL costruisce un URL interno nella lingua corrente.
//
// Restituisce templ.SafeURL perché i file .templ NON importano il pacchetto
// templ: la generazione lo importa da sé, e un import esplicito nei .templ
// produce "templ redeclared in this block".
func (d PageData) URL(suffix string) templ.SafeURL {
	return templ.URL(d.Path(suffix))
}

// Path costruisce il percorso interno nella lingua corrente, come stringa.
func (d PageData) Path(suffix string) string {
	if suffix == "" || suffix == "/" {
		return "/" + d.Lang + "/"
	}
	return "/" + d.Lang + "/" + suffix
}

// BandsPath è il percorso della pagina con l'elenco delle band.
func (d PageData) BandsPath() string { return d.Path("bands/") }

// BandPath è il percorso della pagina di una band.
func (d PageData) BandPath(bandSlug string) string {
	return d.Path("band/" + bandSlug + "/")
}

// AlbumPath è il percorso della pagina di un album.
func (d PageData) AlbumPath(bandSlug, albumSlug string) string {
	return d.Path("band/" + bandSlug + "/album/" + albumSlug + "/")
}

// TrackPath è il percorso della pagina di un brano: è la forma di URL decisa
// nel progetto (/<lang>/band/<band>/album/<album>/brano/<brano>).
func (d PageData) TrackPath(bandSlug, albumSlug, trackSlug string) string {
	return d.Path("band/" + bandSlug + "/album/" + albumSlug + "/brano/" + trackSlug + "/")
}

// AbsURL trasforma un percorso interno in URL assoluto, usando il dominio
// base. È l'unico modo con cui un .templ esce dal sito: il link al feed RSS e
// i link hreflang devono essere assoluti, quelli a pagina possono restare
// relativi.
func (d PageData) AbsURL(path string) templ.SafeURL {
	return templ.URL(strings.TrimSuffix(d.Base, "/") + path)
}

// FeedURL è l'indirizzo assoluto del feed RSS della lingua corrente.
func (d PageData) FeedURL() templ.SafeURL {
	return d.AbsURL(d.Path("rss.xml"))
}

// AlternateLinks sono i link hreflang della pagina corrente: una voce per ogni
// lingua dell'interfaccia più x-default sull'italiano, come nella sitemap.
func (d PageData) AlternateLinks() []AlternateLink {
	if d.NoAlternates {
		return nil
	}
	links := make([]AlternateLink, 0, len(d.Langs)+1)
	for _, lang := range d.Langs {
		links = append(links, AlternateLink{Hreflang: lang, URL: d.AbsURL(langPath(lang, d.Self))})
	}
	links = append(links, AlternateLink{
		Hreflang: "x-default",
		URL:      d.AbsURL(langPath(i18n.DefaultLang, d.Self)),
	})
	return links
}

// langPath mette un percorso di pagina dentro una lingua: la home resta
// "/<lang>/", tutto il resto "/<lang>/<percorso>".
func langPath(lang, self string) string {
	if self == "" {
		return "/" + lang + "/"
	}
	return "/" + lang + "/" + self
}
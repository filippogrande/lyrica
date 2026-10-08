package render

import "github.com/a-h/templ"

// AlternateLink è un link alternato hreflang: la stessa pagina in un'altra
// lingua dell'interfaccia, più la voce x-default sull'italiano. Gli stessi
// legami sono nella sitemap: le due liste non possono divergere.
//
// Il tipo dell'URL è templ.SafeURL perché i file .templ non importano il
// pacchetto templ (la generazione lo importa da sé).
type AlternateLink struct {
	Hreflang string
	URL      templ.SafeURL
}
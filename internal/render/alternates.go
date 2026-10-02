package render

// AlternateLink è un link alternato hreflang: la stessa pagina in un'altra
// lingua dell'interfaccia, più la voce x-default sull'italiano. Gli stessi
// legami sono nella sitemap: le due liste non possono divergere.
type AlternateLink struct {
	Hreflang string
	URL      templSafeURL
}
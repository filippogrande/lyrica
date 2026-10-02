package build

import "strings"

// robotsTxt è il contenuto di /robots.txt (D100).
//
// Permesso a tutti: il sito è pubblico e l'indicizzazione è il punto (D60). Il
// file non è statico nel repo perché le righe Sitemap contengono gli URL
// assoluti delle lingue generate: statico, andrebbe cambiato a mano ogni volta
// che nasce una lingua di interfaccia (D69) e diventerebbe subito falso.
func robotsTxt(langs []string) []byte {
	var sb strings.Builder
	sb.WriteString("User-agent: *\n")
	sb.WriteString("Allow: /\n")
	sb.WriteString("\n")
	sb.WriteString("Sitemap: " + absoluteURL("/sitemap.xml") + "\n")
	for _, lang := range langs {
		sb.WriteString("Sitemap: " + absoluteURL("/"+lang+"/rss.xml") + "\n")
	}
	return []byte(sb.String())
}

// writeRobots genera /robots.txt alla root.
func (w *pageWriter) writeRobots() error {
	return w.writeRaw("robots.txt", robotsTxt(w.langs))
}
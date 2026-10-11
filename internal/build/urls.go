package build

import "strings"

// siteBaseURL è il dominio pubblico del sito (D03). Vive in un posto solo: il
// feed RSS, la sitemap e robots.txt hanno bisogno di URL ASSOLUTI, e un
// dominio in tre posti è un dominio che prima o poi diverge.
//
// Nessuna variabile d'ambiente: il dominio è quello di produzione e il sito
// non viene servito da un altro nome.
const siteBaseURL = "https://lyrica.filippomoscatelli.com"

// absoluteURL unisce il dominio base e il percorso di pagina senza doppi
// slash. Il percorso può arrivare con o senza slash iniziale: i percorsi delle
// pagine lo hanno sempre ("/<lang>/band/..."), quelli delle root no.
func absoluteURL(path string) string {
	base := strings.TrimSuffix(siteBaseURL, "/")
	if path == "" || path == "/" {
		return base + "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return base + path
}
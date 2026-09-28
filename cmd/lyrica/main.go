// Comando lyrica: toolchain del sito (server locale, build, authoring contenuti).
package main

import (
	"fmt"
	"os"

	"github.com/filippogrande/lyrica/internal/web"
)

const usage = `lyrica - toolchain del sito Lyrica

Uso:
  lyrica build                     valida i contenuti e genera public/
  lyrica serve                     serve public/ in locale
  lyrica valida                    controlla i contenuti e stampa errori e avvisi
  lyrica new band|album|brano ...  crea i template di contenuto (vedi docs/CONTENT.md)
  lyrica help                      mostra questo messaggio
`

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "errore: %v\n", err)
		os.Exit(1)
	}
}

// run esegue il comando richiesto. Nessun comando "fa finta": quelli non
// ancora implementati restituiscono un errore esplicito (GUIDELINES.md &S5).
func run(args []string) error {
	if len(args) == 0 {
		fmt.Print(usage)
		return nil
	}

	switch args[0] {
	case "build":
		return runBuild()
	case "serve":
		return web.Serve()
	case "valida":
		return validateContent()
	case "new":
		return runNew(args[1:])
	case "help", "-h", "--help":
		fmt.Print(usage)
		return nil
	default:
		return fmt.Errorf("comando sconosciuto %q\n\n%s", args[0], usage)
	}
}

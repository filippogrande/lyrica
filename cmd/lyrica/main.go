// Comando lyrica: toolchain del sito (server locale, build, authoring contenuti).
package main

import (
	"fmt"
	"os"

	"github.com/filippogrande/lyrica/internal/web"
)

const usage = `lyrica - toolchain del sito Lyrica

Uso:
  lyrica serve                     avvia il sito in locale
  lyrica build                     valida i contenuti e genera public/ (FASE 2)
  lyrica new band|album|brano ...  crea i template di contenuto (FASE 2)
  lyrica help                      mostra questo messaggio
`

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "errore: %v\n", err)
		os.Exit(1)
	}
}

// run esegue il comando richiesto. Nessun comando "fa finta": quelli non
// ancora implementati restituiscono un errore esplicito (DEVELOPMENT_GUIDELINES §5).
func run(args []string) error {
	if len(args) == 0 {
		fmt.Print(usage)
		return nil
	}

	switch args[0] {
	case "serve":
		return web.Serve()
	case "build":
		return fmt.Errorf("comando 'build' non ancora implementato: arriva in FASE 2 (vedi ROADMAP.md)")
	case "new":
		return fmt.Errorf("comando 'new' non ancora implementato: arriva in FASE 2 (vedi ROADMAP.md)")
	case "help", "-h", "--help":
		fmt.Print(usage)
		return nil
	default:
		return fmt.Errorf("comando sconosciuto %q\n\n%s", args[0], usage)
	}
}

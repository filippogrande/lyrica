package build

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// copyAssets copia gli asset statici (CSS, JS, vendor) dentro l'output: il sito
// generato è autonomo, nessuna dipendenza a runtime.
func copyAssets(sourceDir, outputDir string) error {
	return copyTree(sourceDir, filepath.Join(outputDir, "assets"))
}

// copyCovers copia le copertine, se la cartella esiste: un sito senza cover è
// valido, il layout va senza immagine.
func copyCovers(sourceDir, outputDir string) error {
	if _, err := os.Stat(sourceDir); err != nil {
		return nil
	}
	return copyTree(sourceDir, filepath.Join(outputDir, "covers"))
}

// copyTree copia un albero di file. Un file illeggibile è un errore: un asset
// mancante in output è un sito rotto, non un dettaglio.
func copyTree(sourceDir, destinationDir string) error {
	return filepath.WalkDir(sourceDir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("lettura di %s: %w", path, err)
		}
		relative, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("percorso relativo di %s: %w", path, err)
		}
		target := filepath.Join(destinationDir, relative)
		if entry.IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return fmt.Errorf("creazione di %s: %w", target, err)
			}
			return nil
		}
		return copyFile(path, target)
	})
}

// copyFile copia un singolo file.
func copyFile(source, target string) error {
	input, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("apertura di %s: %w", source, err)
	}
	defer input.Close()

	output, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("creazione di %s: %w", target, err)
	}
	if _, err := io.Copy(output, input); err != nil {
		output.Close()
		return fmt.Errorf("copia di %s in %s: %w", source, target, err)
	}
	if err := output.Close(); err != nil {
		return fmt.Errorf("chiusura di %s: %w", target, err)
	}
	return nil
}

package main

import (
	"fmt"
	"strings"
	"time"
)

// today restituisce la data odierna in formato YYYY-MM-DD, nel fuso locale
// della macchina che esegue la CLI (l'ora locale è quella di chi scrive).
func today() string {
	return time.Now().Format("2006-01-02")
}

// splitFront separa un file di contenuto nel front-matter YAML (fra le due
// righe ---) e nel corpo successivo. Se il formato non è quello atteso
// restituisce un errore esplicito (docs/CONTENT.md).
func splitFront(doc []byte) (front []byte, body string, err error) {
	lines := strings.Split(string(doc), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return nil, "", fmt.Errorf("front-matter mancante: il file deve iniziare con ---")
	}
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			front = []byte(strings.Join(lines[1:i], "\n"))
			return front, strings.TrimSpace(strings.Join(lines[i+1:], "\n")), nil
		}
	}
	return nil, "", fmt.Errorf("front-matter non chiuso: manca la riga --- finale")
}

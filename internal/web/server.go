// Package web contiene il server HTTP del sito: routing, header di sicurezza e
// rendering delle pagine.
package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// defaultAddr è la porta di ascolto locale quando LYRICA_ADDR non è impostata.
const defaultAddr = ":8080"

// Serve avvia il server HTTP in locale.
func Serve() error {
	handler, err := Handler()
	if err != nil {
		return err
	}

	addr := os.Getenv("LYRICA_ADDR")
	if addr == "" {
		addr = defaultAddr
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Lyrica in ascolto su http://localhost%s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server: %w", err)
	}
	return nil
}

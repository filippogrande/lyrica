package main

import "github.com/filippogrande/lyrica/internal/build"

// runBuild valida i contenuti e genera il sito in public/.
func runBuild() error {
	return build.Run(build.DefaultOptions())
}

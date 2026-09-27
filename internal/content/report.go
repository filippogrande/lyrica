package content

import (
	"fmt"
	"strings"
)

// Severity distingue ciò che blocca la pubblicazione da ciò che va sistemato.
type Severity int

const (
	// SeverityWarning: build verde, da sistemare o giustificare in PR.
	SeverityWarning Severity = iota
	// SeverityError: contenuto invalido, la build non produce nulla.
	SeverityError
)

// String restituisce l'etichetta italiana della gravità.
func (s Severity) String() string {
	if s == SeverityError {
		return "errore"
	}
	return "avviso"
}

// Issue è un singolo problema trovato nei contenuti.
type Issue struct {
	Severity Severity
	// Where è il file o l'entità a cui si riferisce il problema.
	Where   string
	Message string
}

// String descrive il problema in una riga.
func (i Issue) String() string {
	return fmt.Sprintf("%s: %s: %s", i.Severity, i.Where, i.Message)
}

// Report raccoglie TUTTI i problemi di una validazione: non si ferma al primo,
// altrimenti si sistema un errore per volta e si riapre la PR cinque volte.
type Report struct {
	Issues []Issue
}

// Errorf registra un problema che blocca la pubblicazione.
func (r *Report) Errorf(where, format string, args ...any) {
	r.Issues = append(r.Issues, Issue{
		Severity: SeverityError,
		Where:    where,
		Message:  fmt.Sprintf(format, args...),
	})
}

// Warnf registra un problema che NON blocca la pubblicazione.
func (r *Report) Warnf(where, format string, args ...any) {
	r.Issues = append(r.Issues, Issue{
		Severity: SeverityWarning,
		Where:    where,
		Message:  fmt.Sprintf(format, args...),
	})
}

// Errors elenca i problemi bloccanti.
func (r *Report) Errors() []Issue { return r.bySeverity(SeverityError) }

// Warnings elenca gli avvisi.
func (r *Report) Warnings() []Issue { return r.bySeverity(SeverityWarning) }

// HasErrors dice se c'è almeno un problema bloccante.
func (r *Report) HasErrors() bool { return len(r.bySeverity(SeverityError)) > 0 }

// Err restituisce nil se non ci sono problemi bloccanti, altrimenti il report
// stesso come errore.
func (r *Report) Err() error {
	if !r.HasErrors() {
		return nil
	}
	return r
}

// Error elenca prima gli errori e poi gli avvisi.
func (r *Report) Error() string {
	blocking := r.Errors()
	warnings := r.Warnings()
	lines := make([]string, 0, len(blocking)+len(warnings)+1)
	lines = append(lines, fmt.Sprintf("%d errori, %d avvisi", len(blocking), len(warnings)))
	for _, issue := range blocking {
		lines = append(lines, "- "+issue.String())
	}
	for _, issue := range warnings {
		lines = append(lines, "- "+issue.String())
	}
	return strings.Join(lines, "\n")
}

// bySeverity filtra i problemi per gravità.
func (r *Report) bySeverity(severity Severity) []Issue {
	var filtered []Issue
	for _, issue := range r.Issues {
		if issue.Severity == severity {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

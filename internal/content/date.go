package content

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Date è una data nel formato YYYY-MM-DD usata nei front-matter (added_date).
//
// Il parsing è stretto: una data scritta male produce un errore, non un valore
// vuoto. Una data mancante resta "zero" e la segnala il validatore.
type Date struct {
	t time.Time
}

const dateLayout = "2006-01-02"

// UnmarshalYAML accetta YYYY-MM-DD sia scritta con le virgolette (!!str) sia
// nuda (che YAML classifica come !!timestamp: sono la stessa data, e il formato
// documentato è senza virgolette). Qualsiasi altro tipo è un errore.
func (d *Date) UnmarshalYAML(value *yaml.Node) error {
	switch value.Tag {
	case "!!timestamp":
		var parsed time.Time
		if err := value.Decode(&parsed); err != nil {
			return fmt.Errorf("data non valida: %w", err)
		}
		d.t = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC)
		return nil
	case "!!str":
		var raw string
		if err := value.Decode(&raw); err != nil {
			return fmt.Errorf("data non valida: %w", err)
		}
		parsed, err := time.Parse(dateLayout, strings.TrimSpace(raw))
		if err != nil {
			return fmt.Errorf("data %q non valida: il formato e' YYYY-MM-DD", raw)
		}
		d.t = parsed
		return nil
	default:
		return fmt.Errorf("data attesa come YYYY-MM-DD, trovato un valore di tipo %s", value.Tag)
	}
}

// IsZero dice se la data non è stata impostata.
func (d Date) IsZero() bool { return d.t.IsZero() }

// Time restituisce la data come time.Time (mezzanotte UTC).
func (d Date) Time() time.Time { return d.t }

// String restituisce la data in formato YYYY-MM-DD, vuota se non impostata.
func (d Date) String() string {
	if d.IsZero() {
		return ""
	}
	return d.t.Format(dateLayout)
}

// After dice se la data è successiva a un'altra: serve a ordinare "recenti"
// e feed per added_date decrescente.
func (d Date) After(other Date) bool { return d.t.After(other.t) }

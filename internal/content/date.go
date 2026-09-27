package content

import (
	"fmt"
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

// UnmarshalYAML accetta solo stringhe YYYY-MM-DD.
func (d *Date) UnmarshalYAML(value *yaml.Node) error {
	var raw string
	if err := value.Decode(&raw); err != nil {
		return fmt.Errorf("data attesa come stringa YYYY-MM-DD: %w", err)
	}
	parsed, err := time.Parse(dateLayout, raw)
	if err != nil {
		return fmt.Errorf("data %q non valida: il formato e' YYYY-MM-DD", raw)
	}
	d.t = parsed
	return nil
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

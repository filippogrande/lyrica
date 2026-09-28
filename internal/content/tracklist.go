package content

// TracklistEntry è una voce della tracklist di un album, nell'ordine dichiarato
// in album.md: o il brano con il suo file, o una voce dichiarata senza file (il
// motivo è in Ref.Status). È la lista completa del disco: comprende i brani che
// non hanno ancora un testo, che in pagina restano in elenco senza link.
type TracklistEntry struct {
	Ref   TrackRef
	Track *Track
}

// Published dice se il brano è pubblicato: ha il file e almeno una traduzione.
// Solo i brani pubblicati hanno una pagina e sono linkati in tracklist.
func (e TracklistEntry) Published() bool {
	return e.Track != nil && e.Track.HasTranslations()
}

// Instrumental dice se la voce è strumentale: dichiarata nella voce della
// tracklist (status: instrumental o instrumental: true) o nel file del brano.
func (e TracklistEntry) Instrumental() bool {
	if e.Ref.Status == StatusInstrumental || e.Ref.Instrumental {
		return true
	}
	return e.Track != nil && e.Track.Instrumental
}

// Title è il titolo da stampare: quello del file del brano quando c'è,
// altrimenti quello dichiarato in album.md.
func (e TracklistEntry) Title() string {
	if e.Track != nil {
		return e.Track.Title
	}
	return e.Ref.Title
}

// Slug è lo slug della traccia: quello del file e della voce di tracklist
// (il validatore verifica che coincidano).
func (e TracklistEntry) Slug() string { return e.Ref.Slug }

// Entries restituisce la tracklist dell'album nell'ordine dichiarato, con il
// file del brano quando esiste. L'ordine è quello di album.md: i brani che non
// compaiono nella tracklist non entrano qui (e sono un errore di validazione).
func (a *Album) Entries() []TracklistEntry {
	bySlug := make(map[string]*Track, len(a.Tracks))
	for _, track := range a.Tracks {
		bySlug[track.Slug] = track
	}
	entries := make([]TracklistEntry, 0, len(a.Tracklist))
	for _, ref := range a.Tracklist {
		entries = append(entries, TracklistEntry{Ref: ref, Track: bySlug[ref.Slug]})
	}
	return entries
}

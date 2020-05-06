package canarytools

// InputModule interface receives incidents, and feeds them to consumers
type InputModule interface {
	Feed(incidnetsChan chan<- Incident)
	// TODO: add stats
}

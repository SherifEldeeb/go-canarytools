package canarytools

// Feeder interface receives incidents, and feeds them to consumers
type Feeder interface {
	Feed(incidnetsChan chan<- Incident)
	// TODO: add stats
}

// Filter interface drops, or keeps, incidents based on specified criteria
type Filter interface {
	Filter(incidnetsChan <-chan Incident, filteredIncidnetsChan chan<- Incident)
	// TODO: add stats
}

// Mapper interface maps incident fields to specifed schema,
// then serializes the incident to []byte
type Mapper interface {
	Map(filteredIncidnetsChan <-chan Incident, outChan chan<- []byte)
	// TODO: add stats
}

// Forwarder interface sends the incidnets to their destination
type Forwarder interface {
	Forward(outChan <-chan []byte)
	// TODO: add stats
}

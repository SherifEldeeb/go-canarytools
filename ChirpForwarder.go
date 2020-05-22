package canarytools

import log "github.com/sirupsen/logrus"

// ChirpForwarder is the main struct of the forwarder
// it contains the configurations, and various components of the forwarder
type ChirpForwarder struct {
	// configs
	cfg ChirpForwarderConfig

	// Work chans
	incidentsChan         chan Incident
	filteredIncidentsChan chan Incident
	outChan               chan []byte
	incidentAckerChan     chan []byte

	// interfaces
	feeder        Feeder
	incidentAcker IncidentAcker
	filter        Filter
	mapper        Mapper
	forwarder     Forwarder

	// logger
	l *log.Logger
}

// NewChirpForwarder creates a new chirp forwarder
func NewChirpForwarder(cfg ChirpForwarderConfig) (cf *ChirpForwarder, err error) {
	cf = &ChirpForwarder{}

	cf.cfg = cfg

	// create work chans
	cf.incidentsChan = make(chan Incident)
	cf.filteredIncidentsChan = make(chan Incident)
	cf.outChan = make(chan []byte)
	cf.incidentAckerChan = make(chan []byte)

	// set logger
	cf.l = log.New()
	return
}

// Run starts forwarding incidents
func (cf *ChirpForwarder) Run() {

	// All good, let's roll...
	go cf.feeder.Feed(cf.incidentsChan)
	go cf.incidentAcker.AckIncidents(cf.incidentAckerChan)
	go cf.filter.Filter(cf.incidentsChan, cf.filteredIncidentsChan)
	go cf.mapper.Map(cf.filteredIncidentsChan, cf.outChan)
	cf.forwarder.Forward(cf.outChan, cf.incidentAckerChan)
}

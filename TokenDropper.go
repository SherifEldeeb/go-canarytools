package canarytools

import log "github.com/sirupsen/logrus"

// TokenDropper is the main struct of the droper
// it contains the configurations, and various components
type TokenDropper struct {
	// configs
	cfg TokenDropperConfig

	// Work chans
	// incidentsChan         chan Incident
	// filteredIncidentsChan chan Incident
	// outChan               chan []byte
	// incidentAckerChan     chan []byte

	// interfaces
	// feeder        Feeder
	// incidentAcker IncidentAcker
	// filter        Filter
	// mapper        Mapper
	// forwarder     Forwarder

	// logger
	l *log.Logger
}

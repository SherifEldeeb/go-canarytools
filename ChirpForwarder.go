package canarytools

import log "github.com/sirupsen/logrus"

// ChirpForwarder is the main struct of the forwarder
// it contains the configurations, and various components of the forwarder
type ChirpForwarder struct {
	// configs
	cfg ChirpForwarderConfig
	// interfaces
	feeder        Feeder
	incidentAcker IncidentAcker
	filter        Filter
	mapper        Mapper
	forwarder     Forwarder
	// logger
	l *log.Logger
}

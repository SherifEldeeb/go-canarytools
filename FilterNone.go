package canarytools

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

// FilterNone is a filter that does nothing
type FilterNone struct {
	l *log.Logger
}

func NewFilterNone(loglevel string) (filterNone *FilterNone, err error) {
	filterNone = &FilterNone{}
	// logging config
	filterNone.l = log.New()
	switch loglevel {
	case "info":
		filterNone.l.SetLevel(log.InfoLevel)
	case "warning":
		filterNone.l.SetLevel(log.WarnLevel)
	case "debug":
		filterNone.l.SetLevel(log.DebugLevel)
	default:
		return nil, errors.New("unsupported log level (can be 'info', 'warning' or 'debug')")
	}
	return
}

// Filter filters the incidents, in this case it simply passes them through
func (fn FilterNone) Filter(incidnetsChan <-chan Incident, filteredIncidnetsChan chan<- Incident) {
	for v := range incidnetsChan {
		filteredIncidnetsChan <- v
	}
}

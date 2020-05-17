package canarytools

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// FilterDropEvents is a filter that drops "events" from incidents.
// this will reduce incident size keeping only basic information about the
// incident.
type FilterDropEvents struct {
	l *log.Logger
}

func NewFilterDropEvents(l *log.Logger) (fde *FilterDropEvents, err error) {
	fde = &FilterDropEvents{}
	fde.l = l
	return
}

// Filter filters the incidents, in this case it drops "description.events"
func (fn FilterDropEvents) Filter(incidnetsChan <-chan Incident, filteredIncidnetsChan chan<- Incident) {
	fn.l.WithFields(log.Fields{
		"source": "FilterDropEvents",
		"stage":  "filter",
	}).Info("starting FilterDropEvents")

	for v := range incidnetsChan {
		fn.l.WithFields(log.Fields{
			"source":  "FilterDropEvents",
			"stage":   "filter",
			"content": fmt.Sprintf("%#v", v),
		}).Trace("Orignial incident 'before filtering'")
		if _, ok := v.Description["events"]; ok {
			delete(v.Description, "events")
		} else {
			fn.l.WithFields(log.Fields{
				"source": "FilterDropEvents",
				"stage":  "filter",
			}).Debug("Orignial incident didn't have 'events' array; FilterDropEvents didn't do anything")
		}
		fn.l.WithFields(log.Fields{
			"source":  "FilterDropEvents",
			"stage":   "filter",
			"content": fmt.Sprintf("%#v", v),
		}).Trace("Modified incident 'after filtering'")
		fn.l.WithFields(log.Fields{
			"source":   "FilterDropEvents",
			"stage":    "filter",
			"Incidnet": v.Summary,
		}).Debug("Filter Incident")
		filteredIncidnetsChan <- v
	}
}

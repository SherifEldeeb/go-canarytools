package canarytools

import (
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"
)

// MapperJSON encodes incident details into JSON,
// escapeHTML specifies whether problematic HTML characters should be escaped
// inside JSON quoted strings. The default behavior is to escape &, <, and >
// to \u0026, \u003c, and \u003e to avoid certain safety problems that can
// arise when embedding JSON in HTML.
// In non-HTML requirements where the escaping interferes with the readability
// of the output, setting escapeHTML to (false) disables this behavior.
// TODO: Fix!
type MapperJSON struct {
	// encoder    *json.Encoder
	// buf        *bytes.Buffer
	escapeHTML bool
	// scanner *bufio.Scanner
	l *log.Logger
	// TODO: add stats
}

// NewMapperJSON creates a new MapperJSON
func NewMapperJSON(escapeHTML bool, l *log.Logger) (mapperJSON *MapperJSON, err error) {
	mapperJSON = &MapperJSON{}
	mapperJSON.l = l
	mapperJSON.escapeHTML = escapeHTML
	return
}

// Map maps canary incidents to JSON
func (mj MapperJSON) Map(filteredIncidnetsChan <-chan Incident, outChan chan<- []byte) {
	// new method, till i fogure how to do the json.Encoder with buffer :/
	for v := range filteredIncidnetsChan {
		mj.l.WithFields(log.Fields{
			"source":   "MapperJSON",
			"stage":    "map",
			"incident": v.Summary,
		}).Info("JSON Marshaling incident")

		j, err := json.Marshal(v)
		if err != nil {
			mj.l.WithFields(log.Fields{
				"source": "MapperJSON",
				"stage":  "map",
				"err":    err,
			}).Error("error marshaling value")
			continue
		}
		outChan <- []byte(strings.TrimSpace(string(j)) + "\n")

	}
}

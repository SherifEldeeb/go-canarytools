package canarytools

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"

	log "github.com/sirupsen/logrus"
)

// MapperJSON encodes incident details into JSON,
// escapeHTML specifies whether problematic HTML characters should be escaped
// inside JSON quoted strings. The default behavior is to escape &, <, and >
// to \u0026, \u003c, and \u003e to avoid certain safety problems that can
// arise when embedding JSON in HTML.
// In non-HTML requirements where the escaping interferes with the readability
// of the output, setting escapeHTML to (false) disables this behavior.
type MapperJSON struct {
	escapeHTML bool
	encoder    *json.Encoder
	buf        *bytes.Buffer
	scanner    *bufio.Scanner
	l          *log.Logger

	// TODO: add stats
}

// NewMapperJSON creates a new MapperJSON
func NewMapperJSON(escapeHTML bool, loglevel string) (mapperJSON *MapperJSON, err error) {
	mapperJSON = &MapperJSON{}

	// logging config
	mapperJSON.l = log.New()
	switch loglevel {
	case "info":
		mapperJSON.l.SetLevel(log.InfoLevel)
	case "warning":
		mapperJSON.l.SetLevel(log.WarnLevel)
	case "debug":
		mapperJSON.l.SetLevel(log.DebugLevel)
	default:
		return nil, errors.New("unsupported log level (can be 'info', 'warning' or 'debug')")
	}

	// It'll work like this: JSON Encoder -> (byte slice <-> bytes.Buffer) -> bufio.Scanner -> scanner.Scan -> scanner.Text
	// we do this to benefit from EscapeHTML, otherwise json.Marshal would have been easier.
	b := []byte{}
	mapperJSON.buf = bytes.NewBuffer(b)
	mapperJSON.encoder = json.NewEncoder(mapperJSON.buf)
	mapperJSON.encoder.SetEscapeHTML(escapeHTML)
	mapperJSON.escapeHTML = escapeHTML
	mapperJSON.scanner = bufio.NewScanner(mapperJSON.buf)

	return
}

// Map maps canary incidents to JSON
func (mj MapperJSON) Map(filteredIncidnetsChan <-chan Incident, outChan chan<- []byte) {
	go func() {
		for i := range filteredIncidnetsChan {
			mj.encoder.Encode(i)
		}
	}()
	for mj.scanner.Scan() {
		o := mj.scanner.Text()
		outChan <- []byte(o)
	}

}

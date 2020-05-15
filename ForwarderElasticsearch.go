package canarytools

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"

	log "github.com/sirupsen/logrus"
)

// ElasticForwarder sends alerts to elasticsearch
type ElasticForwarder struct {
	index  string
	client *elasticsearch.Client
	l      *log.Logger
	// TODO: TLS!
}

// NewElasticForwarder creates a new ElasticForwarder,
// it verifies configurations and tries to ping the cluster
func NewElasticForwarder(cfg elasticsearch.Config, index string, l *log.Logger) (elasticforwarder *ElasticForwarder, err error) {
	elasticforwarder = &ElasticForwarder{}

	elasticforwarder.index = index
	elasticforwarder.l = l
	elasticforwarder.client, err = elasticsearch.NewClient(cfg)
	if err != nil {
		l.WithFields(log.Fields{
			"source": "NewElasticForwarder",
			"stage":  "forward",
			"err":    err,
		}).Error("NewElasticForwarder error creating client")
		return
	}
	p, err := elasticforwarder.client.Info()
	if err != nil || p.IsError() {
		l.WithFields(log.Fields{
			"source": "NewElasticForwarder",
			"stage":  "forward",
			"err":    err,
			"status": p.Status(),
		}).Error("NewElasticForwarder error getting cluster info")
		return nil, errors.New(p.String())
	}
	defer p.Body.Close() // freakin' leak!

	l.WithFields(log.Fields{
		"source": "NewElasticForwarder",
		"stage":  "forward",
		"mesage": p,
	}).Info("elasticsearch cluster info")

	return
}

func (ef ElasticForwarder) Forward(outChan <-chan []byte, incidentAckerChan chan<- []byte) {
	for i := range outChan {
		var indexname = ef.index + "-"                          // preparing index name canarychirps-yyyy.MM.dd
		var indexsuffix = time.Now().UTC().Format("2006.01.02") // preparing index suffix canarychirps-yyyy.MM.dd
		var incidentTime = time.Now().UTC()
		// we'll have to unmarshal the incident to extract timestamp,
		// then add "@timestamp" to the event, then marshal it again
		// we could've done that earlier in the pipeline, but to maintain
		// consitency.
		var j map[string]interface{} // temp
		err := json.Unmarshal(i, &j)
		if err != nil {
			ef.l.WithFields(log.Fields{
				"source": "ElasticForwarder",
				"stage":  "forward",
				"err":    err,
			}).Error("Forward unmarshaling incident")
		}
		// getting updated_time
		ut, ok := j["updated_time"] // updated_time: "1588805467"
		if ok {                     // we found it
			utstring, ok := ut.(string) // prevent panic if not string
			if ok {                     // it's a string
				utint, err := strconv.Atoi(utstring) // "1588805467" -> 1588805467
				if err == nil {                      // error converting to int?
					incidentTime = time.Unix(int64(utint), 0).UTC()
					indexsuffix = incidentTime.Format("2006.01.02")
				}
				indexname = indexname + indexsuffix // we now have index name
			}
		}
		// add "@timestamp"
		j["@timestamp"] = incidentTime.Format("2006-01-02T15:04:05.999Z")
		b, err := json.Marshal(j) // we got a json back
		if err != nil {
			ef.l.WithFields(log.Fields{
				"source": "ElasticForwarder",
				"stage":  "forward",
				"err":    err,
			}).Error("Forward error marshaling incident")
		}
		buf := bytes.NewReader(b)
		if err != nil {
			ef.l.WithFields(log.Fields{
				"source": "ElasticForwarder",
				"stage":  "forward",
				"err":    err,
			}).Error("Forward error writing to buffer")
		}
		// Set up the request object.
		req := esapi.IndexRequest{
			Index:   indexname,
			Body:    buf,
			Refresh: "true",
		}

		// Perform the request with the client.
		res, err := req.Do(context.Background(), ef.client)
		if err != nil || res.IsError() {
			ef.l.WithFields(log.Fields{
				"source": "ElasticForwarder",
				"stage":  "forward",
				"err":    err,
				"status": res.Status(),
			}).Error("Forward error indexing document")
			continue
		}
		defer res.Body.Close()
		// add to incident acker

		buf.Seek(0, 0)
		i, _ := ioutil.ReadAll(buf)
		incidentAckerChan <- i
	}
}

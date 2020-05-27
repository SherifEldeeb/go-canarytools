package canarytools

import (
	"reflect"
	"testing"

	"github.com/elastic/go-elasticsearch/v7"
	log "github.com/sirupsen/logrus"
)

func TestNewElasticForwarder(t *testing.T) {
	type args struct {
		cfg   elasticsearch.Config
		index string
		l     *log.Logger
	}
	tests := []struct {
		name                 string
		args                 args
		wantElasticforwarder *ElasticForwarder
		wantErr              bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotElasticforwarder, err := NewElasticForwarder(tt.args.cfg, tt.args.index, tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewElasticForwarder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotElasticforwarder, tt.wantElasticforwarder) {
				t.Errorf("NewElasticForwarder() = %v, want %v", gotElasticforwarder, tt.wantElasticforwarder)
			}
		})
	}
}

func TestElasticForwarder_Forward(t *testing.T) {
	type args struct {
		outChan           <-chan []byte
		incidentAckerChan chan<- []byte
	}
	tests := []struct {
		name string
		ef   ElasticForwarder
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ef.Forward(tt.args.outChan, tt.args.incidentAckerChan)
		})
	}
}

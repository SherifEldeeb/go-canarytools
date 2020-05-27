package canarytools

import (
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestNewMapperJSON(t *testing.T) {
	type args struct {
		escapeHTML bool
		l          *log.Logger
	}
	tests := []struct {
		name           string
		args           args
		wantMapperJSON *MapperJSON
		wantErr        bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMapperJSON, err := NewMapperJSON(tt.args.escapeHTML, tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMapperJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMapperJSON, tt.wantMapperJSON) {
				t.Errorf("NewMapperJSON() = %v, want %v", gotMapperJSON, tt.wantMapperJSON)
			}
		})
	}
}

func TestMapperJSON_Map(t *testing.T) {
	type args struct {
		filteredIncidnetsChan <-chan Incident
		outChan               chan<- []byte
	}
	tests := []struct {
		name string
		mj   MapperJSON
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mj.Map(tt.args.filteredIncidnetsChan, tt.args.outChan)
		})
	}
}

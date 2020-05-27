package canarytools

import (
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestNewFilterNone(t *testing.T) {
	type args struct {
		l *log.Logger
	}
	tests := []struct {
		name           string
		args           args
		wantFilterNone *FilterNone
		wantErr        bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFilterNone, err := NewFilterNone(tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFilterNone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFilterNone, tt.wantFilterNone) {
				t.Errorf("NewFilterNone() = %v, want %v", gotFilterNone, tt.wantFilterNone)
			}
		})
	}
}

func TestFilterNone_Filter(t *testing.T) {
	type args struct {
		incidnetsChan         <-chan Incident
		filteredIncidnetsChan chan<- Incident
	}
	tests := []struct {
		name string
		fn   FilterNone
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fn.Filter(tt.args.incidnetsChan, tt.args.filteredIncidnetsChan)
		})
	}
}

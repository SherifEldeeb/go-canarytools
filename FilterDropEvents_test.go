package canarytools

import (
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestNewFilterDropEvents(t *testing.T) {
	type args struct {
		l *log.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantFde *FilterDropEvents
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFde, err := NewFilterDropEvents(tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFilterDropEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFde, tt.wantFde) {
				t.Errorf("NewFilterDropEvents() = %v, want %v", gotFde, tt.wantFde)
			}
		})
	}
}

func TestFilterDropEvents_Filter(t *testing.T) {
	type args struct {
		incidnetsChan         <-chan Incident
		filteredIncidnetsChan chan<- Incident
	}
	tests := []struct {
		name string
		fn   FilterDropEvents
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

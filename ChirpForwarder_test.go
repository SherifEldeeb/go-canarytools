package canarytools

import (
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestNewChirpForwarder(t *testing.T) {
	type args struct {
		cfg ChirpForwarderConfig
		l   *log.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantCf  *ChirpForwarder
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCf, err := NewChirpForwarder(tt.args.cfg, tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChirpForwarder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCf, tt.wantCf) {
				t.Errorf("NewChirpForwarder() = %v, want %v", gotCf, tt.wantCf)
			}
		})
	}
}

func TestChirpForwarder_setFeeder(t *testing.T) {
	tests := []struct {
		name string
		cf   *ChirpForwarder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cf.setFeeder()
		})
	}
}

func TestChirpForwarder_setFilter(t *testing.T) {
	tests := []struct {
		name string
		cf   *ChirpForwarder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cf.setFilter()
		})
	}
}

func TestChirpForwarder_setTLSConfig(t *testing.T) {
	tests := []struct {
		name string
		cf   *ChirpForwarder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cf.setTLSConfig()
		})
	}
}

func TestChirpForwarder_setForwarder(t *testing.T) {
	tests := []struct {
		name string
		cf   *ChirpForwarder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cf.setForwarder()
		})
	}
}

func TestChirpForwarder_setMapper(t *testing.T) {
	tests := []struct {
		name string
		cf   *ChirpForwarder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cf.setMapper()
		})
	}
}

func TestChirpForwarder_Run(t *testing.T) {
	tests := []struct {
		name string
		cf   *ChirpForwarder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cf.Run()
		})
	}
}

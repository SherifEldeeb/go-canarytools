package main

import (
	"reflect"
	"testing"

	"github.com/SherifEldeeb/canarytools"
	log "github.com/sirupsen/logrus"
)

func Test_popultaeVarsFromEnv(t *testing.T) {
	type args struct {
		cfg *canarytools.ChirpForwarderConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			popultaeVarsFromEnv(tt.args.cfg)
		})
	}
}

func Test_populateVarsFromFlags(t *testing.T) {
	type args struct {
		cfg *canarytools.ChirpForwarderConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			populateVarsFromFlags(tt.args.cfg)
		})
	}
}

func Test_setDefaultVarsEmpty(t *testing.T) {
	var cfg = &canarytools.ChirpForwarderConfig{}
	var l = log.New()
	setDefaultVars(cfg, l)
	if reflect.DeepEqual(cfg, &canarytools.ChirpForwarderConfig{}) {
		t.Error("setting default vars produced an empty ChirpForwarderConfig")
	}
}

func Test_setDefaultVarsDefaultValues(t *testing.T) {
	var cfg = &canarytools.ChirpForwarderConfig{}
	var l = log.New()
	setDefaultVars(cfg, l)

	// default values check
	if l.Level != log.InfoLevel {
		t.Error("loglevel not set to 'info'")
	}
	if cfg.ThenWhat != "nothing" {
		t.Error("ThenWhat not set to 'nothing'")
	}
	if cfg.WhichIncidents != "unacknowledged" {
		t.Error("WhichIncidents not set to 'unacknowledged'")
	}
	if cfg.IncidentFilter != "none" {
		t.Error("IncidentFilter not set to 'none'")
	}
	if cfg.OmFileMaxSize != 8 {
		t.Error("OmFileMaxSize not set to '8'")
	}
	if cfg.OmFileMaxBackups != 14 {
		t.Error("OmFileMaxBackups not set to '14'")
	}
	if cfg.OmFileMaxAge != 120 {
		t.Error("OmFileMaxAge not set to '120'")
	}
	if cfg.OmFileName != "canaryChirps.json" {
		t.Error("OmFileName not set to 'canaryChirps.json'")
	}
}

func Test_setDefaultVars(t *testing.T) {
	type args struct {
		cfg *canarytools.ChirpForwarderConfig
		l   *log.Logger
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setDefaultVars(tt.args.cfg, tt.args.l)
		})
	}
}

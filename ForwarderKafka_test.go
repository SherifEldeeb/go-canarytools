package canarytools

import (
	"crypto/tls"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestNewKafkaForwarder(t *testing.T) {
	type args struct {
		brokers   []string
		topic     string
		tlsconfig *tls.Config
		l         *log.Logger
	}
	tests := []struct {
		name               string
		args               args
		wantKafkaforwarder *KafkaForwarder
		wantErr            bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKafkaforwarder, err := NewKafkaForwarder(tt.args.brokers, tt.args.topic, tt.args.tlsconfig, tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKafkaForwarder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotKafkaforwarder, tt.wantKafkaforwarder) {
				t.Errorf("NewKafkaForwarder() = %v, want %v", gotKafkaforwarder, tt.wantKafkaforwarder)
			}
		})
	}
}

func TestKafkaForwarder_Forward(t *testing.T) {
	type args struct {
		outChan           <-chan []byte
		incidentAckerChan chan<- []byte
	}
	tests := []struct {
		name string
		kf   KafkaForwarder
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.kf.Forward(tt.args.outChan, tt.args.incidentAckerChan)
		})
	}
}

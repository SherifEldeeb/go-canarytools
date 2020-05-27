package canarytools

import (
	"crypto/tls"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestNewTCPForwarder(t *testing.T) {
	type args struct {
		host      string
		port      int
		tlsConfig *tls.Config
		sslUseSSL bool
		l         *log.Logger
	}
	tests := []struct {
		name             string
		args             args
		wantTcpforwarder *TCPForwarder
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTcpforwarder, err := NewTCPForwarder(tt.args.host, tt.args.port, tt.args.tlsConfig, tt.args.sslUseSSL, tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTCPForwarder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTcpforwarder, tt.wantTcpforwarder) {
				t.Errorf("NewTCPForwarder() = %v, want %v", gotTcpforwarder, tt.wantTcpforwarder)
			}
		})
	}
}

func TestTCPForwarder_Forward(t *testing.T) {
	type args struct {
		outChan           <-chan []byte
		incidentAckerChan chan<- []byte
	}
	tests := []struct {
		name string
		t    TCPForwarder
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.Forward(tt.args.outChan, tt.args.incidentAckerChan)
		})
	}
}

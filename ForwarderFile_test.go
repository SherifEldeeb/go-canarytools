package canarytools

import (
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestFileForwader_Forward(t *testing.T) {
	type args struct {
		outChan           <-chan []byte
		incidentAckerChan chan<- []byte
	}
	tests := []struct {
		name string
		f    FileForwader
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.Forward(tt.args.outChan, tt.args.incidentAckerChan)
		})
	}
}

func TestNewFileForwader(t *testing.T) {
	type args struct {
		filename   string
		maxsize    int
		maxbackups int
		maxage     int
		compress   bool
		l          *log.Logger
	}
	tests := []struct {
		name             string
		args             args
		wantFileforwader *FileForwader
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFileforwader, err := NewFileForwader(tt.args.filename, tt.args.maxsize, tt.args.maxbackups, tt.args.maxage, tt.args.compress, tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFileForwader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFileforwader, tt.wantFileforwader) {
				t.Errorf("NewFileForwader() = %v, want %v", gotFileforwader, tt.wantFileforwader)
			}
		})
	}
}

package canarytools

import (
	"reflect"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestNewClient(t *testing.T) {
	type args struct {
		domain         string
		apikey         string
		thenWhat       string
		sinceWhen      string
		whichIncidents string
		fetchInterval  int
		l              *log.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantC   *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := NewClient(tt.args.domain, tt.args.apikey, tt.args.thenWhat, tt.args.sinceWhen, tt.args.whichIncidents, tt.args.fetchInterval, tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("NewClient() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func TestClient_Ping(t *testing.T) {
	tests := []struct {
		name    string
		c       Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("Client.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_GetUnacknowledgedIncidents(t *testing.T) {
	type args struct {
		since time.Time
	}
	tests := []struct {
		name          string
		c             Client
		args          args
		wantIncidents []Incident
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIncidents, err := tt.c.GetUnacknowledgedIncidents(tt.args.since)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetUnacknowledgedIncidents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIncidents, tt.wantIncidents) {
				t.Errorf("Client.GetUnacknowledgedIncidents() = %v, want %v", gotIncidents, tt.wantIncidents)
			}
		})
	}
}

func TestClient_GetAllDevices(t *testing.T) {
	tests := []struct {
		name        string
		c           Client
		wantDevices []Device
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDevices, err := tt.c.GetAllDevices()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetAllDevices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDevices, tt.wantDevices) {
				t.Errorf("Client.GetAllDevices() = %v, want %v", gotDevices, tt.wantDevices)
			}
		})
	}
}

func TestClient_GetLiveDevices(t *testing.T) {
	tests := []struct {
		name        string
		c           Client
		wantDevices []Device
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDevices, err := tt.c.GetLiveDevices()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetLiveDevices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDevices, tt.wantDevices) {
				t.Errorf("Client.GetLiveDevices() = %v, want %v", gotDevices, tt.wantDevices)
			}
		})
	}
}

func TestClient_GetDeadDevices(t *testing.T) {
	tests := []struct {
		name        string
		c           Client
		wantDevices []Device
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDevices, err := tt.c.GetDeadDevices()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetDeadDevices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDevices, tt.wantDevices) {
				t.Errorf("Client.GetDeadDevices() = %v, want %v", gotDevices, tt.wantDevices)
			}
		})
	}
}

func TestClient_getIncidents(t *testing.T) {
	type args struct {
		which string
		since time.Time
	}
	tests := []struct {
		name          string
		c             Client
		args          args
		wantIncidents []Incident
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIncidents, err := tt.c.getIncidents(tt.args.which, tt.args.since)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.getIncidents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIncidents, tt.wantIncidents) {
				t.Errorf("Client.getIncidents() = %v, want %v", gotIncidents, tt.wantIncidents)
			}
		})
	}
}

func TestClient_getDevices(t *testing.T) {
	type args struct {
		which string
	}
	tests := []struct {
		name        string
		c           Client
		args        args
		wantDevices []Device
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDevices, err := tt.c.getDevices(tt.args.which)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.getDevices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDevices, tt.wantDevices) {
				t.Errorf("Client.getDevices() = %v, want %v", gotDevices, tt.wantDevices)
			}
		})
	}
}

func TestClient_GetAllIncidents(t *testing.T) {
	type args struct {
		since time.Time
	}
	tests := []struct {
		name          string
		c             Client
		args          args
		wantIncidents []Incident
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIncidents, err := tt.c.GetAllIncidents(tt.args.since)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetAllIncidents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIncidents, tt.wantIncidents) {
				t.Errorf("Client.GetAllIncidents() = %v, want %v", gotIncidents, tt.wantIncidents)
			}
		})
	}
}

func TestClient_Feed(t *testing.T) {
	type args struct {
		incidnetsChan chan<- Incident
	}
	tests := []struct {
		name string
		c    *Client
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Feed(tt.args.incidnetsChan)
		})
	}
}

func TestClient_AckIncidents(t *testing.T) {
	type args struct {
		ackedIncident <-chan []byte
	}
	tests := []struct {
		name string
		c    *Client
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.AckIncidents(tt.args.ackedIncident)
		})
	}
}

func TestClient_AckIncident(t *testing.T) {
	type args struct {
		incident string
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.AckIncident(tt.args.incident); (err != nil) != tt.wantErr {
				t.Errorf("Client.AckIncident() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package canarytools

import "testing"

func TestLoadTokenFile(t *testing.T) {
	type args struct {
		f string
	}
	tests := []struct {
		name          string
		args          args
		wantApikey    string
		wantApidomain string
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotApikey, gotApidomain, err := LoadTokenFile(tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadTokenFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotApikey != tt.wantApikey {
				t.Errorf("LoadTokenFile() gotApikey = %v, want %v", gotApikey, tt.wantApikey)
			}
			if gotApidomain != tt.wantApidomain {
				t.Errorf("LoadTokenFile() gotApidomain = %v, want %v", gotApidomain, tt.wantApidomain)
			}
		})
	}
}

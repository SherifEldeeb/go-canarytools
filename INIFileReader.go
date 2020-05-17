package canarytools

import (
	"errors"

	"gopkg.in/ini.v1"
)

// LoadTokenFile reads the token file '.canarytools.config' that contains the
// API key, and your unique canary domain.
// You can get that file 'after enabling API' from:
// Canary console -> API -> Download Token File.
func LoadTokenFile(f string) (apikey string, apidomain string, err error) {
	token, err := ini.Load(f)
	if err != nil {
		return
	}
	apikey = token.Section("CanaryTools").Key("api_key").String()
	apidomain = token.Section("CanaryTools").Key("domain").String()
	if apikey == "" || apidomain == "" {
		err = errors.New("couldn't get key or domain from file")
	}
	return
}

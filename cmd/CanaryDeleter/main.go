package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	canarytools "github.com/thinkst/go-canarytools"
)

var (
	cfg canarytools.CanaryDeleterConfig
	err error

	// constants for build-time hardcoding of params
	// those could be set at build time to create a self-running executable
	// go build -ldflags "-X main.DOMAIN=$domain_hash  -X main.APIKEY=$api_auth -w -s -linkmode=internal"

	// DOMAIN is the domain hash
	DOMAIN string
	// APIKEY is the main console API auth token
	APIKEY string
	// BUILDTIME is the time the tools was built
	BUILDTIME string
	// SHA1VER is the built sha1
	SHA1VER string
)

func init() {
	populateVarsFromFlags(&cfg)
}

func main() {
	flag.Parse()
	l := log.New()
	l.WithFields(log.Fields{
		"BUILDTIME": BUILDTIME,
		"SHA1VER":   SHA1VER,
	}).Info("starting CanaryDeleter")

	// Set LogLevel
	switch cfg.LogLevel {
	case "info":
		l.SetLevel(log.InfoLevel)
	case "warning":
		l.SetLevel(log.WarnLevel)
	case "debug":
		l.SetLevel(log.DebugLevel)
	case "trace":
		l.SetLevel(log.TraceLevel)
	default:
		l.Warn("unsupported log level, or none specified; will set to 'Debug'")
		l.SetLevel(log.DebugLevel)
	}

	// Set all hardcoded info, if provided
	if DOMAIN != "" && cfg.ConsoleAPIDomain == "" { // command line values always Supersede hardcoded ones
		l.Debug("found pre-configured domain hash value")
		cfg.ConsoleAPIDomain = DOMAIN
	}
	if APIKEY != "" && cfg.ConsoleAPIKey == "" { // command line values always Supersede hardcoded ones
		l.Debug("found pre-configured API auth value")
		cfg.ConsoleAPIKey = APIKEY
	}

	// Finish config logic
	err = finishConfig(&cfg, l)
	if err != nil {
		l.WithField("err", err).Fatal("configuration error")
	}

	l.Info("building an API client")
	c, err := canarytools.NewClient(cfg.ConsoleAPIConfig, l)
	if err != nil {
		l.Fatal(err)
	}
	l.WithField("FlockName", cfg.FlockName).Info("getting flock_id from FlockName")
	// does the flock exist?
	flockID, err := c.GetFlockIDFromName(cfg.FlockName)
	if err != nil {
		l.Fatal(err)
	}
	cfg.FlockID = flockID
	l.WithField("FlockName", cfg.FlockName).WithField("flock_id", cfg.FlockID).Info("got flock_id")

	switch cfg.DeleteWhat {
	case "incidents":
		l.WithField("FlockName", cfg.FlockName).WithField("flock_id", cfg.FlockID).Info("deleting all incidents")
		err = c.DeleteMultipleIncidents("flock_id", flockID, cfg.IncludeUnacknowledged)
		if err != nil {
			l.Fatal(err)
		}
	case "tokens":
		t, err := c.FetchCanarytokenAll()
		if err != nil {
			l.Fatal(err)
		}
		l.WithField("FlockName", cfg.FlockName).WithField("flock_id", cfg.FlockID).Info("deleting all tokens")
		for _, token := range t {
			if token.FlockID == flockID {
				l.Info("deleteing:", token.Canarytoken)
				err = c.DeleteCanarytoken(token.Canarytoken)
				if err != nil {
					l.Error("error:", err)
				}
			}
		}
	default:
		l.Fatal("you have to tell me what to delete using '-what' ... supported values are 'incidents' & 'tokens'")
	}
	l.Info("done!")

}

func finishConfig(cfg *canarytools.CanaryDeleterConfig, l *log.Logger) (err error) {
	// first, we didn't get api key and domain through flags? let's try to load them from file
	if cfg.ConsoleAPIKey == "" && cfg.ConsoleAPIDomain == "" {
		// if we don't have them, we try to load it from same drectory
		if cfg.ConsoleTokenFile == "" { // if not
			cwd, _ := os.Getwd()
			cfg.ConsoleTokenFile = filepath.Join(cwd, "canarytools.config")
		}
		// do we have canarytools.config in same path? get data from it...
		if _, err := os.Stat(cfg.ConsoleTokenFile); os.IsNotExist(err) {
			return fmt.Errorf("canarytools.config does not exist, and we couldn't get domain hash and API key")
		}
		cfg.ConsoleAPIKey, cfg.ConsoleAPIDomain, err = canarytools.LoadTokenFile(cfg.ConsoleTokenFile)
		if err != nil || cfg.ConsoleAPIDomain == "" || cfg.ConsoleAPIKey == "" {
			return fmt.Errorf("error parsing token file: %s", err)
		}
	}
	if cfg.FlockName == "" {
		return fmt.Errorf("You have to specify the flock name using '-flock'")
	}

	return
}

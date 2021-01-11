package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	canarytools "github.com/thinkst/go-canarytools"
)

func populateVarsFromFlags(cfg *canarytools.CanaryDeleterConfig) {
	// General flags
	flag.StringVar(&cfg.ConsoleAPIDomain, "domain", "", "Canary console domain (hash)")
	flag.StringVar(&cfg.LogLevel, "loglevel", "info", "set loglevel, can be one of ('info', 'warning', 'debug' or 'trace')")
	flag.StringVar(&cfg.ConsoleAPIKey, "apikey", "", "API key")
	flag.StringVar(&cfg.ConsoleTokenFile, "tokenfile", "", "the token file 'canarytools.config' which contains auth token and the domain")

	// What to cleanup? valid options are "alerts" and "tokens"
	flag.StringVar(&cfg.DeleteWhat, "what", "incidents", `What to cleanup? valid options are "incidents" and "tokens"`)
	flag.BoolVar(&cfg.IncludeUnacknowledged, "include-unacknowledged-incidents", true, `Include Unacknowledged Incidents?`)
	flag.BoolVar(&cfg.DumpToJson, "dump", true, `dump incidents to a JSON file before deleting them`)
	flag.BoolVar(&cfg.DumpOnly, "dumponly", false, `only dump incidents to a JSON file?`)

	// Flock Specific flags
	flag.StringVar(&cfg.FlockName, "flock", "", "Which flock to target?")

	// Node specific flag
	flag.StringVar(&cfg.NodeID, "node", "", "Which 'Node ID' to target?")
}

func finishConfig(cfg *canarytools.CanaryDeleterConfig, l *log.Logger) (err error) {
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

	// flock or node?
	if cfg.FlockName != "" && cfg.NodeID != "" {
		l.Fatal("you can't provide both '-flock' and '-node' at the same time ... pick one")
	}

	if cfg.FlockName != "" {
		cfg.FilterType = "flock_id"
	}

	if cfg.NodeID != "" {
		cfg.FilterType = "node_id"
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
	if cfg.FlockName == "" && cfg.NodeID == "" {
		return fmt.Errorf("You have to specify either the flock name using '-flock', or node ID using '-node'")
	}

	if cfg.DumpOnly && !cfg.DumpToJson {
		return fmt.Errorf("you have set -dump to false, and -dumponly to true, which doesn't work")
	}

	return
}

package main

import (
	"flag"

	canarytools "github.com/thinkst/go-canarytools"
)

func populateVarsFromFlags(cfg *canarytools.CanaryDeleterConfig) {
	// General flags
	flag.StringVar(&cfg.ConsoleAPIDomain, "domain", "", "Canary console domain (hash)")
	flag.StringVar(&cfg.LogLevel, "loglevel", "info", "set loglevel, can be one of ('info', 'warning', 'debug' or 'trace')")
	flag.StringVar(&cfg.ConsoleAPIKey, "apikey", "", "API key")
	flag.StringVar(&cfg.ConsoleTokenFile, "tokenfile", "", "the token file 'canarytools.config' which contains auth token and the domain")

	// What to cleanup? valid options are "alerts" and "tokens"
	flag.StringVar(&cfg.DeleteWhat, "what", "", `What to cleanup? valid options are "alerts" and "tokens"`)
	flag.BoolVar(&cfg.IncludeUnacknowledged, "include-unacknowledged-incidents", true, `Include Unacknowledged Incidents?`)

	// Flock Specific flags
	flag.StringVar(&cfg.FlockName, "flock", "", "Which flock to target?")
}

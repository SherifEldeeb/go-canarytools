package main

import (
	"flag"

	canarytools "github.com/thinkst/go-canarytools"
)

func populateVarsFromFlags(cfg *canarytools.ConsoleAPIConfig) {
	// General flags
	flag.StringVar(&cfg.ConsoleAPIDomain, "domain", "", "Canary console domain (hash)")
	flag.StringVar(&cfg.OpMode, "opmode", "api", "operate using console API or Factory? valid values are 'api' & 'factory'")
	flag.StringVar(&cfg.ConsoleAPIKey, "apikey", "", "API key (can't be specified with '-factoryauth')")
	flag.StringVar(&cfg.ConsoleTokenFile, "tokenfile", "", "the token file 'canarytools.config' which contains auth token and the domain")
}

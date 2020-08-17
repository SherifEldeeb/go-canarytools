package main

import (
	"flag"
	"fmt"

	"github.com/SherifEldeeb/canarytools"
	log "github.com/sirupsen/logrus"
)

var (
	chirpforwarder canarytools.ChirpForwarder
	cfg            canarytools.ChirpForwarderConfig
	err            error
)

func init() {
	populateVarsFromFlags(&cfg) // first: set vars with flags
	popultaeVarsFromEnv(&cfg)   // then:  populate remaining vars from environment
	// explicit command line flags overrides environment variables; values from
	// environment variables are only set if not already set by flags
}

func main() {
	// Fancy logos are quite important
	fmt.Println(logo)

	// create logger, this will be used throughout!
	// TODO: support file log out.
	l := log.New()

	// parse arguments
	flag.Parse()

	// setting default vars for those that are not set
	setDefaultVars(&cfg, l)

	// create a new chirpforwarder
	cf, err := canarytools.NewChirpForwarder(cfg, l)
	if err != nil {
		l.WithFields(log.Fields{
			"err": err,
		}).Fatal("error setting up ChirpForwarder")
	}

	// All good, let's roll...
	cf.Run()
}

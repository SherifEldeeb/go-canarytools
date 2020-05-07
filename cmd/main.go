package main

import (
	"flag"

	"github.com/SherifEldeeb/canarytools"
	log "github.com/sirupsen/logrus"
)

var (
	// General flags
	feederModule    = flag.String("feeder", "consoleapi", "input module")
	forwarderModule = flag.String("output", "tcp", "output module")
	loglevel        = flag.String("loglevel", "info", "set loglevel, can be one of ('info', 'warning' or 'debug')")
	thenWhat        = flag.String("then", "nothing", "what to do after getting an incident? ")

	// INPUT MODULES
	// Console API input module
	fmConsoleAPIKey           = flag.String("apikey", "", "API Key")
	fmConsoleAPIDomain        = flag.String("domain", "", "canarytools domain")
	fmConsoleAPIFetchInterval = flag.Int("interval", 5, "alert fetch interval 'in seconds'")
	// TODO: webhook
	// TODO: syslog

	// OUTPUT MODULES
	// TCP/UDP output module
	omTCPUDPPort = flag.Int("port", 4455, "[output] TCP/UDP port")
	omTCPUDPHost = flag.String("host", "127.0.0.1", "[output] host")
)

// interface placeholders
var (
	feeder    canarytools.Feeder
	filter    canarytools.Filter
	mapper    canarytools.Mapper
	forwarder canarytools.Forwarder
)

// implemented modules
var (
	validFeederModules = map[string]bool{
		"consoleapi": true,
	}
	validForwarderModules = map[string]bool{
		"tcp":  true,
		"file": true,
	}
)

func main() {
	log.Info("starting canary ChirpForwarder")

	// create chans
	var incidentsChan = make(chan canarytools.Incident)
	var filteredIncidentsChan = make(chan canarytools.Incident)
	var outChan = make(chan []byte)

	// parse arguments
	flag.Parse()

	// setting log level
	switch *loglevel {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.Fatal("unsupported log level (can be 'info', 'warning' or 'debug')")
	}

	// few sanity checks
	// valid input module?
	_, ok := validFeederModules[*feederModule]
	if !ok {
		log.Fatal("invalid input module specifed")
	}
	_, ok = validForwarderModules[*forwarderModule]
	if !ok {
		log.Fatal("invalid output module specifed")
	}

	// Input modules look good?
	switch *feederModule {
	case "consoleapi":
		if len(*fmConsoleAPIKey) != 32 {
			log.Fatal("invalid API Key (length != 32)")
		}
		if *fmConsoleAPIDomain == "" {
			log.Fatal("domain must be provided")
		}
		////////////////////
		// start...
		log.WithFields(log.Fields{
			"domain":          *fmConsoleAPIDomain,
			"fmConsoleAPIKey": (*fmConsoleAPIKey)[0:4] + "..." + (*fmConsoleAPIKey)[len(*fmConsoleAPIKey)-4:len(*fmConsoleAPIKey)],
		}).Info("ChirpForwarder Configs")

		// building a new clint, testing connection...
		log.Debug("building new client and pinging console")
		c, err := canarytools.NewClient(*fmConsoleAPIDomain, *fmConsoleAPIKey, *loglevel, *fmConsoleAPIFetchInterval)
		if err != nil {
			log.WithFields(log.Fields{
				"message": err,
			}).Fatal("error during creating client, or pinging console")
		}
		log.Debug("ping successful! we're good to go")
		feeder = c
	}

	// Output modules look good?
	switch *forwarderModule {
	case "tcp":
		// bulding new TCP out
		t, err := canarytools.NewTCPForwarder(*omTCPUDPHost, *omTCPUDPPort, *loglevel)
		if err != nil {
			log.WithFields(log.Fields{
				"message": err,
			}).Fatal("error during creating TCP Out client")
		}
		forwarder = t
	}

	// filter
	filter = &canarytools.FilterNone{}

	//mapper
	mapper, err := canarytools.NewMapperJSON(false, *loglevel)
	if err != nil {
		log.WithFields(log.Fields{
			"message": err,
		}).Fatal("error creating JON Mapper")
	}

	// All good, let's roll...
	go feeder.Feed(incidentsChan)
	go filter.Filter(incidentsChan, filteredIncidentsChan)
	go mapper.Map(filteredIncidentsChan, outChan)
	forwarder.Forward(outChan)

}

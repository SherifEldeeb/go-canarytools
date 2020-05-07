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

	// create logger & setting log level
	l := log.New()
	switch *loglevel {
	case "info":
		l.SetLevel(log.InfoLevel)
	case "warning":
		l.SetLevel(log.WarnLevel)
	case "debug":
		l.SetLevel(log.DebugLevel)
	default:
		l.Fatal("unsupported log level (can be 'info', 'warning' or 'debug')")
	}

	// few sanity checks
	// valid input module?
	_, ok := validFeederModules[*feederModule]
	if !ok {
		l.Fatal("invalid input module specifed")
	}
	_, ok = validForwarderModules[*forwarderModule]
	if !ok {
		l.Fatal("invalid output module specifed")
	}

	// Input modules look good?
	switch *feederModule {
	case "consoleapi":
		if len(*fmConsoleAPIKey) != 32 {
			l.Fatal("invalid API Key (length != 32)")
		}
		if *fmConsoleAPIDomain == "" {
			l.Fatal("domain must be provided")
		}
		////////////////////
		// start...
		l.WithFields(log.Fields{
			"domain":          *fmConsoleAPIDomain,
			"fmConsoleAPIKey": (*fmConsoleAPIKey)[0:4] + "..." + (*fmConsoleAPIKey)[len(*fmConsoleAPIKey)-4:len(*fmConsoleAPIKey)],
		}).Info("ChirpForwarder Configs")

		// building a new clint, testing connection...
		l.Debug("building new client and pinging console")
		c, err := canarytools.NewClient(*fmConsoleAPIDomain, *fmConsoleAPIKey, *fmConsoleAPIFetchInterval, l)
		if err != nil {
			l.WithFields(log.Fields{
				"message": err,
			}).Fatal("error during creating client, or pinging console")
		}
		l.Debug("ping successful! we're good to go")
		feeder = c
	}

	// Output modules look good?
	switch *forwarderModule {
	case "tcp":
		// bulding new TCP out
		t, err := canarytools.NewTCPForwarder(*omTCPUDPHost, *omTCPUDPPort, l)
		if err != nil {
			l.WithFields(log.Fields{
				"message": err,
			}).Fatal("error during creating TCP Out client")
		}
		forwarder = t
	}

	// filter
	filter, err := canarytools.NewFilterNone(l)
	if err != nil {
		l.WithFields(log.Fields{
			"message": err,
		}).Fatal("error creating None filter")
	}

	//mapper
	mapper, err := canarytools.NewMapperJSON(false, l)
	if err != nil {
		l.WithFields(log.Fields{
			"message": err,
		}).Fatal("error creating JON Mapper")
	}

	// All good, let's roll...
	go feeder.Feed(incidentsChan)
	go filter.Filter(incidentsChan, filteredIncidentsChan)
	go mapper.Map(filteredIncidentsChan, outChan)
	forwarder.Forward(outChan)

}

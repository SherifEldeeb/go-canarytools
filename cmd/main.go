package main

import (
	"flag"
	"sync"

	"github.com/SherifEldeeb/canarytools"
	log "github.com/sirupsen/logrus"
)

var (
	// General flags
	inputModule  = flag.String("input", "consoleapi", "input module")
	outputModule = flag.String("output", "tcp", "output module")
	loglevel     = flag.String("loglevel", "info", "set loglevel, can be one of ('info', 'warning' or 'debug')")

	// INPUT MODULES
	// Console API input module
	imConsoleAPIKey           = flag.String("apikey", "", "API Key")
	imConsoleAPIDomain        = flag.String("domain", "", "canarytools domain")
	imConsoleAPIFetchInterval = flag.Int("interval", 5, "alert fetch interval 'in seconds'")
	// TODO: webhook
	// TODO: syslog

	// OUTPUT MODULES
	// TCP/UDP output module
	omTCPUDPPort = flag.Int("port", 4455, "[output] TCP/UDP port")
	omTCPUDPHost = flag.String("host", "127.0.0.1", "[output] host")
)

// implemented modules
var (
	validInputModules = map[string]bool{
		"consoleapi": true,
	}
	validOutputModules = map[string]bool{
		"tcp": true,
	}
)

func main() {
	log.Info("starting canary ChirpForwarder")

	// create chans
	var incidentsChan = make(chan canarytools.Incident)

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
	_, ok := validInputModules[*inputModule]
	if !ok {
		log.Fatal("invalid input module specifed")
	}
	_, ok = validOutputModules[*outputModule]
	if !ok {
		log.Fatal("invalid output module specifed")
	}

	// Input modules look good?
	// var im interface{}
	switch *inputModule {
	case "consoleapi":
		if len(*imConsoleAPIKey) != 32 {
			log.Fatal("invalid API Key (length != 32)")
		}
		if *imConsoleAPIDomain == "" {
			log.Fatal("domain must be provided")
		}
		// im = canarytools.Client{}
	}

	// Output modules look good?
	// var om interface{}
	switch *outputModule {
	case "tcp":
		if *omTCPUDPPort > 65535 || *omTCPUDPPort < 1 {
			log.Fatal("invalid port number")
		}
		if *omTCPUDPHost == "" {
			log.Fatal("output host can't be empty")
		}
		// om = canarytools.Client{}
	}

	////////////////////
	// start...
	log.WithFields(log.Fields{
		"domain":          *imConsoleAPIDomain,
		"imConsoleAPIKey": (*imConsoleAPIKey)[0:4] + "..." + (*imConsoleAPIKey)[len(*imConsoleAPIKey)-4:len(*imConsoleAPIKey)],
	}).Info("ChirpForwarder Configs")

	// building a new clint, testing connection...
	log.Debug("building new client and pinging console")
	c, err := canarytools.NewClient(*imConsoleAPIDomain, *imConsoleAPIKey, *loglevel, *imConsoleAPIFetchInterval)
	if err != nil {
		log.WithFields(log.Fields{
			"message": err,
		}).Fatal("error during creating client, or pinging console")
	}
	log.Debug("ping successful! we're good to go")

	// bulding new TCP out
	t, err := canarytools.NewTCPOutput(*omTCPUDPHost, *omTCPUDPPort)
	if err != nil {
		log.WithFields(log.Fields{
			"message": err,
		}).Fatal("error during creating TCP Out client")
	}
	// feed'em
	go c.Feed(incidentsChan)
	go t.Out(incidentsChan)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

}

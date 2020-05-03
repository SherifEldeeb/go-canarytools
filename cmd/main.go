package main

import (
	"flag"

	"github.com/SherifEldeeb/canarytools"
	log "github.com/sirupsen/logrus"
)

var (
	apiKey   = flag.String("apikey", "", "API Key")
	domain   = flag.String("domain", "", "canarytools domain")
	loglevel = flag.String("loglevel", "info", "set loglevel, can be one of ('info', 'warning' or 'debug')")
)

func main() {
	log.Info("starting canary ChirpManager")

	// parse arguments
	flag.Parse()

	// few sanity checks
	if len(*apiKey) != 32 {
		log.Fatal("invalid API Key (length != 32)")
	}
	if *domain == "" {
		log.Fatal("domain must be provided")
	}

	// start...
	log.WithFields(log.Fields{
		"domain": *domain,
		"apikey": (*apiKey)[0:4] + "..." + (*apiKey)[len(*apiKey)-4:len(*apiKey)],
	}).Info("ChirpManager Configs")

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

	// building a new clint, testing connection...
	log.Debug("building new client and pinging console")
	c, err := canarytools.NewClient(*domain, *apiKey, *loglevel)
	if err != nil {
		log.WithFields(log.Fields{
			"message": err,
		}).Fatal("error during creating client, or pinging console")
	}
	log.Debug("ping successful! we're good to go")

	// get all devices
	log.Debug("getting all devices")
	dvcs, err := c.GetAllDevices()
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("found total of %d devices", len(dvcs))

	// get all unacked incidents
	log.Debug("getting all unacked incidents")
	unackedInc, maxUpdateID, err := c.GetUnacknowledgedIncidents(0)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("found total of %d unacked incidents; max incident ID %d", len(unackedInc), maxUpdateID)
	for _, v := range unackedInc {
		log.WithFields(log.Fields{
			"UpdatedID": v.UpdatedID,
		}).Debug(v.Summary)
	}

}

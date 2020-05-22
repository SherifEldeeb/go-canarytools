package canarytools

import (
	"os"
	"os/user"
	"path"

	log "github.com/sirupsen/logrus"
)

// ChirpForwarder is the main struct of the forwarder
// it contains the configurations, and various components of the forwarder
type ChirpForwarder struct {
	// configs
	cfg ChirpForwarderConfig

	// Work chans
	incidentsChan         chan Incident
	filteredIncidentsChan chan Incident
	outChan               chan []byte
	incidentAckerChan     chan []byte

	// interfaces
	feeder        Feeder
	incidentAcker IncidentAcker
	filter        Filter
	mapper        Mapper
	forwarder     Forwarder

	// logger
	l *log.Logger
}

// NewChirpForwarder creates a new chirp forwarder
func NewChirpForwarder(cfg ChirpForwarderConfig, l *log.Logger) (cf *ChirpForwarder, err error) {
	cf = &ChirpForwarder{}

	cf.cfg = cfg

	// create work chans
	cf.incidentsChan = make(chan Incident)
	cf.filteredIncidentsChan = make(chan Incident)
	cf.outChan = make(chan []byte)
	cf.incidentAckerChan = make(chan []byte)

	// set logger
	cf.l = l
	return
}

func (cf *ChirpForwarder) setFeeder() {
	var err error
	switch cf.cfg.FeederModule {
	case "consoleapi":
		// did you specify both token file && manually using apikey+domain?
		if cf.cfg.ImConsoleTokenFile != "" && (cf.cfg.ImConsoleAPIDomain != "" || cf.cfg.ImConsoleAPIKey != "") {
			cf.l.Fatal("look, you either use 'tokenfile' or 'apikey+domain', not both")
		}
		// so, what if token file is not specfied, but neither apikey+domain?
		// we'll look for the "canarytools.config" file in user's home directory
		if cf.cfg.ImConsoleTokenFile == "" && cf.cfg.ImConsoleAPIDomain == "" && cf.cfg.ImConsoleAPIKey == "" {
			cf.l.Warn("none of 'tokenfile', 'apikey' & 'domain' has been provided! will look for 'canarytools.config' file in user's home directory")
			u, err := user.Current()
			if err != nil {
				cf.l.WithFields(log.Fields{
					"err": err,
				}).Fatal("error getting current user")
			}
			cf.cfg.ImConsoleTokenFile = path.Join(u.HomeDir, "canarytools.config")
			cf.l.WithField("path", cf.cfg.ImConsoleTokenFile).Warn("automatically looking for canarytools.config")
			if _, err := os.Stat(cf.cfg.ImConsoleTokenFile); os.IsNotExist(err) {
				cf.l.Fatal("couldn't get apikey+domain! provide using environment variables, command line flags, or path to token file")
			}
		}
		// tokenfile specified? get values from there
		if cf.cfg.ImConsoleTokenFile != "" {
			cf.cfg.ImConsoleAPIKey, cf.cfg.ImConsoleAPIDomain, err = LoadTokenFile(cf.cfg.ImConsoleTokenFile)
			if err != nil || cf.cfg.ImConsoleAPIDomain == "" || cf.cfg.ImConsoleAPIKey == "" {
				cf.l.WithFields(log.Fields{
					"err":    err,
					"api":    cf.cfg.ImConsoleAPIKey,
					"domain": cf.cfg.ImConsoleAPIDomain,
				}).Fatal("error parsing token file")
			}
			cf.l.WithFields(log.Fields{
				"path":   cf.cfg.ImConsoleTokenFile,
				"api":    cf.cfg.ImConsoleAPIKey,
				"domain": cf.cfg.ImConsoleAPIDomain,
			}).Info("successfully parsed token file, using values from there")
		}
		// few checks
		if len(cf.cfg.ImConsoleAPIKey) != 32 {
			cf.l.Fatal("invalid API Key (length != 32)")
		}
		if cf.cfg.ImConsoleAPIDomain == "" {
			cf.l.Fatal("domain must be provided")
		}
		////////////////////
		// start...
		cf.l.WithFields(log.Fields{
			"domain":                 cf.cfg.ImConsoleAPIDomain,
			"cf.cfg.ImConsoleAPIKey": (cf.cfg.ImConsoleAPIKey)[0:4] + "..." + (cf.cfg.ImConsoleAPIKey)[len(cf.cfg.ImConsoleAPIKey)-4:len(cf.cfg.ImConsoleAPIKey)],
		}).Info("ChirpForwarder Configs")

		// building a new clint, testing connection...
		cf.l.Debug("building new client and pinging console")
		c, err := NewClient(cf.cfg.ImConsoleAPIDomain, cf.cfg.ImConsoleAPIKey, cf.cfg.ThenWhat, cf.cfg.SinceWhenString, cf.cfg.WhichIncidents, cf.cfg.ImConsoleAPIFetchInterval, cf.l)
		if err != nil {
			cf.l.WithFields(log.Fields{
				"err": err,
			}).Fatal("error during creating client, or pinging console")
		}
		cf.l.Debug("ping successful! we're good to go")
		cf.feeder = c
		cf.incidentAcker = c
	default:
		cf.l.WithField("feeder", cf.cfg.FeederModule).Fatal("unsupported feeder module specified")
	}
}

// Run starts forwarding incidents
func (cf *ChirpForwarder) Run() {

	// All good, let's roll...
	go cf.feeder.Feed(cf.incidentsChan)
	go cf.incidentAcker.AckIncidents(cf.incidentAckerChan)
	go cf.filter.Filter(cf.incidentsChan, cf.filteredIncidentsChan)
	go cf.mapper.Map(cf.filteredIncidentsChan, cf.outChan)
	cf.forwarder.Forward(cf.outChan, cf.incidentAckerChan)
}

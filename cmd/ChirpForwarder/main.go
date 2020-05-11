package main

import (
	"crypto/tls"
	"flag"
	"net/http"
	"time"

	"github.com/SherifEldeeb/canarytools"
	"github.com/elastic/go-elasticsearch/v7"
	log "github.com/sirupsen/logrus"
	"github.com/stackimpact/stackimpact-go"
)

var (
	// General flags
	feederModule    = flag.String("feeder", "consoleapi", "input module")
	forwarderModule = flag.String("output", "tcp", "output module")
	loglevel        = flag.String("loglevel", "info", "set loglevel, can be one of ('info', 'warning' or 'debug')")
	thenWhat        = flag.String("then", "nothing", "what to do after getting an incident? can be one of ('nothing', or 'ack')")
	sinceWhenString = flag.String("since", "", `get events newer than this time.
	format has to be like this: 'yyyy-MM-dd HH:mm:ss'
	if nothing provided, it will check value from '.canary.lastcheck' file,
	if .canary.lastcheck file does not exist, it will default to events from last 7 days`)
	whichIncidents = flag.String("which", "unacknowledged", "which incidents to fetch? can be one of ('all', or 'unacknowledged')")

	// INPUT MODULES
	// Console API input module
	imConsoleAPIKey           = flag.String("apikey", "", "API Key")
	imConsoleAPIDomain        = flag.String("domain", "", "canarytools domain")
	imConsoleAPIFetchInterval = flag.Int("interval", 5, "alert fetch interval 'in seconds'")
	// TODO: webhook
	// TODO: syslog

	// OUTPUT MODULES
	// TCP/UDP output module
	omTCPUDPPort = flag.Int("port", 4455, "[OUT|TCP] TCP/UDP port")
	omTCPUDPHost = flag.String("host", "127.0.0.1", "[OUT|TCP] host")

	// File forward module
	omFileMaxSize    = flag.Int("maxsize", 8, "[OUT|FILE] file max size in megabytes")
	omFileMaxBackups = flag.Int("maxbackups", 14, "[OUT|FILE] file max number of files to keep")
	omFileMaxAge     = flag.Int("maxage", 120, "[OUT|FILE] file max age in days 'older than this will be deleted'")
	omFileCompress   = flag.Bool("compress", false, "[OUT|FILE] file compress log files?")
	omFileName       = flag.String("filename", "canaryChirps.json", "[OUT|FILE] file name")

	// elasticsearch forward module
	omElasticHost        = flag.String("eshost", "http://127.0.0.1:9200", "[OUT|ELASTIC] elasticsearch host")
	omElasticUser        = flag.String("esuser", "elastic", "[OUT|ELASTIC] elasticsearch user 'basic auth'")
	omElasticPass        = flag.String("espass", "elastic", "[OUT|ELASTIC] elasticsearch password 'basic auth'")
	omElasticCloudAPIKey = flag.String("escloudapikey", "", "[OUT|ELASTIC] elasticsearch Base64-encoded token for authorization; if set, overrides username and password")
	omElasticCloudID     = flag.String("escloudid", "", "[OUT|ELASTIC] endpoint for the Elastic Service (https://elastic.co/cloud)")
	omElasticIndex       = flag.String("esindex", "canarychirps", "[OUT|ELASTIC] elasticsearch index")
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
		"tcp":     true,
		"file":    true,
		"elastic": true,
	}
)

func main() {
	// Profiler Start
	agent := stackimpact.Start(stackimpact.Options{
		AgentKey: "aff482334b4e5bf0d9f4fea81dda16fa8068eb32",
		AppName:  "ChirpForwarder",
	})
	span := agent.Profile()
	defer span.Stop()
	//

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
	case "trace":
		l.SetLevel(log.TraceLevel)
	default:
		l.Fatal("unsupported log level (should be 'info', 'warning', 'debug' or 'trace')")
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
		if len(*imConsoleAPIKey) != 32 {
			l.Fatal("invalid API Key (length != 32)")
		}
		if *imConsoleAPIDomain == "" {
			l.Fatal("domain must be provided")
		}
		////////////////////
		// start...
		l.WithFields(log.Fields{
			"domain":          *imConsoleAPIDomain,
			"imConsoleAPIKey": (*imConsoleAPIKey)[0:4] + "..." + (*imConsoleAPIKey)[len(*imConsoleAPIKey)-4:len(*imConsoleAPIKey)],
		}).Info("ChirpForwarder Configs")

		// building a new clint, testing connection...
		l.Debug("building new client and pinging console")
		c, err := canarytools.NewClient(*imConsoleAPIDomain, *imConsoleAPIKey, *thenWhat, *sinceWhenString, *whichIncidents, *imConsoleAPIFetchInterval, l)
		if err != nil {
			l.WithFields(log.Fields{
				"err": err,
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
				"err": err,
			}).Fatal("error during creating TCP Out client")
		}
		forwarder = t
	case "file":
		// bulding new file out
		ff, err := canarytools.NewFileForwader(*omFileName, *omFileMaxSize, *omFileMaxBackups, *omFileMaxAge, *omFileCompress, l)
		if err != nil {
			l.WithFields(log.Fields{
				"err": err,
			}).Fatal("error during creating File Out client")
		}
		forwarder = ff
	case "elastic":
		// bulding new elastic out
		cfg := elasticsearch.Config{
			Addresses: []string{*omElasticHost}, // A list of Elasticsearch nodes to use.
			Username:  *omElasticUser,           // Username for HTTP Basic Authentication.
			Password:  *omElasticPass,           // Password for HTTP Basic Authentication.
			CloudID:   *omElasticCloudID,
			APIKey:    *omElasticCloudAPIKey,
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   10,
				ResponseHeaderTimeout: time.Second,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
		ef, err := canarytools.NewElasticForwarder(cfg, *omElasticIndex, l)
		if err != nil {
			l.WithFields(log.Fields{
				"err": err,
			}).Fatal("error during creating File Out client")
		}
		forwarder = ef
	default:
		l.Fatal("unsupported output module")
	}

	// filter
	filter, err := canarytools.NewFilterNone(l)
	if err != nil {
		l.WithFields(log.Fields{
			"err": err,
		}).Fatal("error creating None filter")
	}

	//mapper
	mapper, err := canarytools.NewMapperJSON(false, l)
	if err != nil {
		l.WithFields(log.Fields{
			"err": err,
		}).Fatal("error creating JON Mapper")
	}

	// All good, let's roll...
	go feeder.Feed(incidentsChan)
	go filter.Filter(incidentsChan, filteredIncidentsChan)
	go mapper.Map(filteredIncidentsChan, outChan)
	forwarder.Forward(outChan)
}

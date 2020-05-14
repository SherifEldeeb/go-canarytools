package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/SherifEldeeb/canarytools"
	"github.com/elastic/go-elasticsearch/v7"
	log "github.com/sirupsen/logrus"
	"github.com/stackimpact/stackimpact-go"
)

var (
	// General flags
	feederModule    string
	forwarderModule string
	loglevel        string
	thenWhat        string
	sinceWhenString string
	whichIncidents  string

	// SSL/TLS Client configs
	// used by TCP & Elastic output
	sslUseSSL       bool
	sslSkipInsecure bool
	sslCA           string
	sslKey          string
	sslCert         string

	// INPUT MODULES
	// Console API input module
	imConsoleAPIKey           string
	imConsoleAPIDomain        string
	imConsoleAPIFetchInterval int

	// OUTPUT MODULES
	// TCP/UDP output module
	omTCPUDPPort int
	omTCPUDPHost string

	// File forward module
	omFileMaxSize    int
	omFileMaxBackups int
	omFileMaxAge     int
	omFileCompress   bool
	omFileName       string

	// elasticsearch forward module
	omElasticHost        string
	omElasticUser        string
	omElasticPass        string
	omElasticCloudAPIKey string
	omElasticCloudID     string
	omElasticIndex       string
)

// setting vars
func init() {
	// General flags
	flag.StringVar(&feederModule, "feeder", "consoleapi", "input module")
	flag.StringVar(&forwarderModule, "output", "tcp", "output module")
	flag.StringVar(&loglevel, "loglevel", "info", "set loglevel, can be one of ('info', 'warning' or 'debug')")
	flag.StringVar(&thenWhat, "then", "nothing", "what to do after getting an incident? can be one of ('nothing', or 'ack')")
	flag.StringVar(&sinceWhenString, "since", "", `get events newer than this time.
	format has to be like this: 'yyyy-MM-dd HH:mm:ss'
	if nothing provided, it will check value from '.canary.lastcheck' file,
	if .canary.lastcheck file does not exist, it will default to events from last 7 days`)
	flag.StringVar(&whichIncidents, "which", "unacknowledged", "which incidents to fetch? can be one of ('all', or 'unacknowledged')")

	// SSL/TLS Client configs
	// used by TCP & Elastic output
	flag.BoolVar(&sslUseSSL, "ssl", false, "[SSL/TLS CLIENT] are we using SSL/TLS? setting this to true enables encrypted clinet configs")
	flag.BoolVar(&sslSkipInsecure, "insecure", false, "[SSL/TLS CLIENT] ignore cert errors")
	flag.StringVar(&sslCA, "sslclientca", "", "[SSL/TLS CLIENT] path to client rusted CA certificate file")
	flag.StringVar(&sslKey, "sslclientkey", "", "[SSL/TLS CLIENT] path to client SSL/TLS Key  file")
	flag.StringVar(&sslCert, "sslclientcert", "", "[SSL/TLS CLIENT] path to client SSL/TLS cert  file")

	// INPUT MODULES
	// Console API input module
	flag.StringVar(&imConsoleAPIKey, "apikey", "", "API Key")
	flag.StringVar(&imConsoleAPIDomain, "domain", "", "canarytools domain")
	flag.IntVar(&imConsoleAPIFetchInterval, "interval", 5, "alert fetch interval 'in seconds'")
	// TODO: webhook
	// TODO: syslog

	// OUTPUT MODULES
	// TCP/UDP output module
	flag.IntVar(&omTCPUDPPort, "port", 4455, "[OUT|TCP] TCP/UDP port")
	flag.StringVar(&omTCPUDPHost, "host", "127.0.0.1", "[OUT|TCP] host")

	// File forward module
	flag.IntVar(&omFileMaxSize, "maxsize", 8, "[OUT|FILE] file max size in megabytes")
	flag.IntVar(&omFileMaxBackups, "maxbackups", 14, "[OUT|FILE] file max number of files to keep")
	flag.IntVar(&omFileMaxAge, "maxage", 120, "[OUT|FILE] file max age in days 'older than this will be deleted'")
	flag.BoolVar(&omFileCompress, "compress", false, "[OUT|FILE] file compress log files?")
	flag.StringVar(&omFileName, "filename", "canaryChirps.json", "[OUT|FILE] file name")

	// elasticsearch forward module
	flag.StringVar(&omElasticHost, "eshost", "http://127.0.0.1:9200", "[OUT|ELASTIC] elasticsearch host")
	flag.StringVar(&omElasticUser, "esuser", "elastic", "[OUT|ELASTIC] elasticsearch user 'basic auth'")
	flag.StringVar(&omElasticPass, "espass", "elastic", "[OUT|ELASTIC] elasticsearch password 'basic auth'")
	flag.StringVar(&omElasticCloudAPIKey, "escloudapikey", "", "[OUT|ELASTIC] elasticsearch Base64-encoded token for authorization; if set, overrides username and password")
	flag.StringVar(&omElasticCloudID, "escloudid", "", "[OUT|ELASTIC] endpoint for the Elastic Cloud Service 'https://elastic.co/cloud'")
	flag.StringVar(&omElasticIndex, "esindex", "canarychirps", "[OUT|ELASTIC] elasticsearch index")
}

// interface placeholders
var (
	feeder        canarytools.Feeder
	incidentAcker canarytools.IncidentAcker
	filter        canarytools.Filter
	mapper        canarytools.Mapper
	forwarder     canarytools.Forwarder
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
	log.Info("starting canary ChirpForwarder")
	// Profiler Start
	agent := stackimpact.Start(stackimpact.Options{
		AgentKey: "aff482334b4e5bf0d9f4fea81dda16fa8068eb32",
		AppName:  "ChirpForwarder",
	})
	span := agent.Profile()
	defer span.Stop()
	// Profiler end

	// parse arguments
	flag.Parse()
	// Check Environment
	// env := os.Environ()

	// create chans
	var incidentsChan = make(chan canarytools.Incident)
	var filteredIncidentsChan = make(chan canarytools.Incident)
	var outChan = make(chan []byte)
	var incidentAckerChan = make(chan []byte)

	// create logger & setting log level
	l := log.New()
	switch loglevel {
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
	_, ok := validFeederModules[feederModule]
	if !ok {
		l.Fatal("invalid input module specifed")
	}
	_, ok = validForwarderModules[forwarderModule]
	if !ok {
		l.Fatal("invalid output module specifed")
	}

	// Input modules look good?
	switch feederModule {
	case "consoleapi":
		if len(imConsoleAPIKey) != 32 {
			l.Fatal("invalid API Key (length != 32)")
		}
		if imConsoleAPIDomain == "" {
			l.Fatal("domain must be provided")
		}
		////////////////////
		// start...
		l.WithFields(log.Fields{
			"domain":          imConsoleAPIDomain,
			"imConsoleAPIKey": (imConsoleAPIKey)[0:4] + "..." + (imConsoleAPIKey)[len(imConsoleAPIKey)-4:len(imConsoleAPIKey)],
		}).Info("ChirpForwarder Configs")

		// building a new clint, testing connection...
		l.Debug("building new client and pinging console")
		c, err := canarytools.NewClient(imConsoleAPIDomain, imConsoleAPIKey, thenWhat, sinceWhenString, whichIncidents, imConsoleAPIFetchInterval, l)
		if err != nil {
			l.WithFields(log.Fields{
				"err": err,
			}).Fatal("error during creating client, or pinging console")
		}
		l.Debug("ping successful! we're good to go")
		feeder = c
		incidentAcker = c
	}

	// Prepping SSL/TLS configs
	var tlsConfig = &tls.Config{}
	if sslUseSSL {
		// ignore cert verification errors?
		tlsConfig.InsecureSkipVerify = sslSkipInsecure
		// custom CA?
		if sslCA != "" {
			// Get the SystemCertPool, continue with an empty pool on error
			rootCAs, _ := x509.SystemCertPool()
			if rootCAs == nil {
				rootCAs = x509.NewCertPool()
			}
			// Read in the cert file
			certs, err := ioutil.ReadFile(sslCA)
			if err != nil {
				l.WithFields(log.Fields{
					"err":    err,
					"cafile": sslCA,
				}).Fatal("Failed to read CA file")
			}
			// Append our cert to the system pool
			if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
				l.Fatal("couldn't add CA cert! (file might be improperly formatted)")
			}
			tlsConfig.RootCAs = rootCAs
		}
		// custom key + cert?
		if sslKey != "" && sslCert != "" {
			// Load client cert
			clientCert, err := tls.LoadX509KeyPair(sslCert, sslKey)
			if err != nil {
				l.Fatal(err)
			}
			tlsConfig.Certificates = []tls.Certificate{clientCert}
		}
	}

	// Output modules look good?
	switch forwarderModule {
	case "tcp":
		// bulding new TCP out
		t, err := canarytools.NewTCPForwarder(omTCPUDPHost, omTCPUDPPort, tlsConfig, sslUseSSL, l)
		if err != nil {
			l.WithFields(log.Fields{
				"err": err,
			}).Fatal("error during creating TCP Out client")
		}
		forwarder = t
	case "file":
		// bulding new file out
		ff, err := canarytools.NewFileForwader(omFileName, omFileMaxSize, omFileMaxBackups, omFileMaxAge, omFileCompress, l)
		if err != nil {
			l.WithFields(log.Fields{
				"err": err,
			}).Fatal("error during creating File Out client")
		}
		forwarder = ff
	case "elastic":
		// bulding new elastic out
		cfg := elasticsearch.Config{
			Addresses: []string{omElasticHost}, // A list of Elasticsearch nodes to use.
			Username:  omElasticUser,           // Username for HTTP Basic Authentication.
			Password:  omElasticPass,           // Password for HTTP Basic Authentication.
			CloudID:   omElasticCloudID,
			APIKey:    omElasticCloudAPIKey,
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   10,
				ResponseHeaderTimeout: time.Duration(10) * time.Second,
				TLSClientConfig:       tlsConfig,
			},
		}
		ef, err := canarytools.NewElasticForwarder(cfg, omElasticIndex, l)
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
	go incidentAcker.AckIncidents(incidentAckerChan)
	go filter.Filter(incidentsChan, filteredIncidentsChan)
	go mapper.Map(filteredIncidentsChan, outChan)
	forwarder.Forward(outChan, incidentAckerChan)
}

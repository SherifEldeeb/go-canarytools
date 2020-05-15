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
	feederModule    string // CANARY_FEEDER
	forwarderModule string // CANARY_OUTPUT
	loglevel        string // CANARY_LOGLEVEL
	thenWhat        string // CANARY_THEN
	sinceWhenString string // CANARY_SINCE
	whichIncidents  string // CANARY_WHICH

	// SSL/TLS Client configs
	// used by TCP & Elastic output
	sslUseSSL       bool   // CANARY_SSL
	sslSkipInsecure bool   // CANARY_INSECURE
	sslCA           string // CANARY_SSLCLIENTCA
	sslKey          string // CANARY_SSLCLIENTKEY
	sslCert         string // CANARY_SSLCLIENTCERT

	// INPUT MODULES
	// Console API input module
	imConsoleAPIKey           string // CANARY_APIKEY
	imConsoleAPIDomain        string // CANARY_DOMAIN
	imConsoleAPIFetchInterval int    // CANARY_INTERVAL

	// OUTPUT MODULES
	// TCP/UDP output module
	omTCPUDPPort int    // CANARY_PORT
	omTCPUDPHost string // CANARY_HOST

	// File forward module
	omFileMaxSize    int    // CANARY_MAXSIZE
	omFileMaxBackups int    // CANARY_MAXBACKUPS
	omFileMaxAge     int    // CANARY_MAXAGE
	omFileCompress   bool   // CANARY_COMPRESS
	omFileName       string // CANARY_FILENAME

	// elasticsearch forward module
	omElasticHost        string // CANARY_ESHOST
	omElasticUser        string // CANARY_ESUSER
	omElasticPass        string // CANARY_ESPASS
	omElasticCloudAPIKey string // CANARY_ESCLOUDAPIKEY
	omElasticCloudID     string // CANARY_ESCLOUDID
	omElasticIndex       string // CANARY_ESINDEX

)

// interface placeholders
var (
	feeder        canarytools.Feeder
	incidentAcker canarytools.IncidentAcker
	filter        canarytools.Filter
	mapper        canarytools.Mapper
	forwarder     canarytools.Forwarder
)

// setting vars
func init() {
	popultaeVarsFromEnv()   // first: get from environment variables
	populateVarsFromFlags() // then: override with flags (if set)
}

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
	default:
		l.WithField("feeder", feederModule).Fatal("unsupported feeder module specified")
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
		l.WithField("outputModule", forwarderModule).Fatal("unsupported output module")
	}

	// filter
	// currently, we haven't build anything here yet
	filter, err := canarytools.NewFilterNone(l)
	if err != nil {
		l.WithFields(log.Fields{
			"err": err,
		}).Fatal("error creating None filter")
	}

	// mapper
	// only JSON mapper is implemented
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

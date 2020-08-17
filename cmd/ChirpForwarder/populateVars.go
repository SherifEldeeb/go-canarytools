package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/SherifEldeeb/canarytools"

	log "github.com/sirupsen/logrus"
)

func popultaeVarsFromEnv(cfg *canarytools.ChirpForwarderConfig) {
	// Artificial Intelligence vars parser!
	// lots of if statements ahead ... :)

	// General flags
	if cfg.FeederModule == "" {
		cfg.FeederModule, _ = os.LookupEnv("CANARY_FEEDER")
	}
	if cfg.ForwarderModule == "" {
		cfg.ForwarderModule, _ = os.LookupEnv("CANARY_OUTPUT")
	}
	if cfg.Loglevel == "" {
		cfg.Loglevel, _ = os.LookupEnv("CANARY_LOGLEVEL")
	}
	if cfg.ThenWhat == "" {
		cfg.ThenWhat, _ = os.LookupEnv("CANARY_THEN")
	}
	if cfg.SinceWhenString == "" {
		cfg.SinceWhenString, _ = os.LookupEnv("CANARY_SINCE")
	}
	if cfg.WhichIncidents == "" {
		cfg.WhichIncidents, _ = os.LookupEnv("CANARY_WHICH")
	}
	if cfg.IncidentFilter == "" {
		cfg.IncidentFilter, _ = os.LookupEnv("CANARY_FILTER")
	}

	// SSL/TLS Client configs
	// used by TCP &cfg. Elastic output
	if cfg.SSLUseSSL == false {
		sslUseSSLBool, _ := os.LookupEnv("CANARY_SSL")
		cfg.SSLUseSSL, _ = strconv.ParseBool(sslUseSSLBool)
	}
	if cfg.SSLSkipInsecure == false {
		sslSkipInsecureBool, _ := os.LookupEnv("CANARY_INSECURE")
		cfg.SSLSkipInsecure, _ = strconv.ParseBool(sslSkipInsecureBool)
	}
	if cfg.SSLCA == "" {
		cfg.SSLCA, _ = os.LookupEnv("CANARY_SSLCLIENTCA")
	}
	if cfg.SSLKey == "" {
		cfg.SSLKey, _ = os.LookupEnv("CANARY_SSLCLIENTKEY")
	}
	if cfg.SSLCert == "" {
		cfg.SSLCert, _ = os.LookupEnv("CANARY_SSLCLIENTCERT")
	}

	// INPUT MODULES
	// Console API input module
	if cfg.ImConsoleAPIKey == "" {
		cfg.ImConsoleAPIKey, _ = os.LookupEnv("CANARY_APIKEY")
	}
	if cfg.ImConsoleAPIDomain == "" {
		cfg.ImConsoleAPIDomain, _ = os.LookupEnv("CANARY_DOMAIN")
	}
	if cfg.ImConsoleTokenFile == "" {
		cfg.ImConsoleTokenFile, _ = os.LookupEnv("CANARY_TOKENFILE")
	}
	if cfg.ImConsoleAPIFetchInterval == 0 {
		imConsoleAPIFetchIntervalInt, _ := os.LookupEnv("CANARY_INTERVAL")
		cfg.ImConsoleAPIFetchInterval, _ = strconv.Atoi(imConsoleAPIFetchIntervalInt)
	}

	// OUTPUT MODULES
	// TCP/UDP output module
	if cfg.OmTCPUDPPort == 0 {
		omTCPUDPPortInt, _ := os.LookupEnv("CANARY_PORT")
		cfg.OmTCPUDPPort, _ = strconv.Atoi(omTCPUDPPortInt)
	}
	if cfg.OmTCPUDPHost == "" {
		cfg.OmTCPUDPHost, _ = os.LookupEnv("CANARY_HOST")
	}

	// File forward module
	if cfg.OmFileMaxSize == 0 {
		omFileMaxSizeInt, _ := os.LookupEnv("CANARY_MAXSIZE")
		cfg.OmFileMaxSize, _ = strconv.Atoi(omFileMaxSizeInt)
	}
	if cfg.OmFileMaxBackups == 0 {
		omFileMaxBackupsInt, _ := os.LookupEnv("CANARY_MAXBACKUPS")
		cfg.OmFileMaxBackups, _ = strconv.Atoi(omFileMaxBackupsInt)
	}
	if cfg.OmFileMaxAge == 0 {
		omFileMaxAgeInt, _ := os.LookupEnv("CANARY_MAXAGE")
		cfg.OmFileMaxAge, _ = strconv.Atoi(omFileMaxAgeInt)
	}
	if cfg.OmFileCompress == false {
		omFileCompressBool, _ := os.LookupEnv("CANARY_COMPRESS")
		cfg.OmFileCompress, _ = strconv.ParseBool(omFileCompressBool)
	}
	if cfg.OmFileName == "" {
		cfg.OmFileName, _ = os.LookupEnv("CANARY_FILENAME")
	}

	// elasticsearch forward module
	if cfg.OmElasticHost == "" {
		cfg.OmElasticHost, _ = os.LookupEnv("CANARY_ESHOST")
	}
	if cfg.OmElasticUser == "" {
		cfg.OmElasticUser, _ = os.LookupEnv("CANARY_ESUSER")
	}
	if cfg.OmElasticPass == "" {
		cfg.OmElasticPass, _ = os.LookupEnv("CANARY_ESPASS")
	}
	if cfg.OmElasticCloudAPIKey == "" {
		cfg.OmElasticCloudAPIKey, _ = os.LookupEnv("CANARY_ESCLOUDAPIKEY")
	}
	if cfg.OmElasticCloudID == "" {
		cfg.OmElasticCloudID, _ = os.LookupEnv("CANARY_ESCLOUDID")
	}
	if cfg.OmElasticIndex == "" {
		cfg.OmElasticIndex, _ = os.LookupEnv("CANARY_ESINDEX")
	}

	// kafka forward module
	if cfg.OmKafkaBrokers == "" {
		cfg.OmKafkaBrokers, _ = os.LookupEnv("CANARY_KAFKABROKERS")
	}
	if cfg.OmKafkaTopic == "" {
		cfg.OmKafkaTopic, _ = os.LookupEnv("CANARY_KAFKATOPIC")
	}
}

func populateVarsFromFlags(cfg *canarytools.ChirpForwarderConfig) {
	// General flags
	flag.StringVar(&cfg.FeederModule, "feeder", "consoleapi", "input module")
	flag.StringVar(&cfg.ForwarderModule, "output", "", "output module ('tcp', 'file', 'elastic' or 'kafka')")
	flag.StringVar(&cfg.Loglevel, "loglevel", "", "set loglevel, can be one of ('info', 'warning' or 'debug')")
	flag.StringVar(&cfg.ThenWhat, "then", "", "what to do after getting an incident? can be one of ('nothing', or 'ack')")
	flag.StringVar(&cfg.SinceWhenString, "since", "", `get events newer than this time.
		format has to be like this: 'yyyy-MM-dd HH:mm:ss'
		if nothing provided, it will check value from '.canary.lastcheck' file,
		if .canary.lastcheck file does not exist, it will default to events from last 7 days`)
	flag.StringVar(&cfg.WhichIncidents, "which", "", "which incidents to fetch? can be one of ('all', or 'unacknowledged')")
	flag.StringVar(&cfg.IncidentFilter, "filter", "", "filter to apply to incident ('none', or 'dropevents')")

	// SSL/TLS Client configs
	// used by TCP &cfg. Elastic output
	flag.BoolVar(&cfg.SSLUseSSL, "ssl", false, "[SSL/TLS CLIENT] are we using SSL/TLS? setting this to true enables encrypted clinet configs")
	flag.BoolVar(&cfg.SSLSkipInsecure, "insecure", false, "[SSL/TLS CLIENT] ignore cert errors")
	flag.StringVar(&cfg.SSLCA, "sslclientca", "", "[SSL/TLS CLIENT] path to client rusted CA certificate file")
	flag.StringVar(&cfg.SSLKey, "sslclientkey", "", "[SSL/TLS CLIENT] path to client SSL/TLS Key  file")
	flag.StringVar(&cfg.SSLCert, "sslclientcert", "", "[SSL/TLS CLIENT] path to client SSL/TLS cert  file")

	// INPUT MODULES
	// Console API input module
	flag.StringVar(&cfg.ImConsoleAPIKey, "apikey", "", "API Key")
	flag.StringVar(&cfg.ImConsoleAPIDomain, "domain", "", "canarytools domain")
	flag.StringVar(&cfg.ImConsoleTokenFile, "tokenfile", "", "the token file 'canarytools.config' which contains api token and the domain")
	flag.IntVar(&cfg.ImConsoleAPIFetchInterval, "interval", 0, "alert fetch interval 'in seconds'")

	// OUTPUT MODULES
	// TCP/UDP output module
	flag.IntVar(&cfg.OmTCPUDPPort, "port", 0, "[OUT|TCP] TCP port")
	flag.StringVar(&cfg.OmTCPUDPHost, "host", "", "[OUT|TCP] host")

	// File forward module
	flag.IntVar(&cfg.OmFileMaxSize, "maxsize", 0, "[OUT|FILE] file max size in megabytes")
	flag.IntVar(&cfg.OmFileMaxBackups, "maxbackups", 0, "[OUT|FILE] file max number of files to keep")
	flag.IntVar(&cfg.OmFileMaxAge, "maxage", 0, "[OUT|FILE] file max age in days 'older than this will be deleted'")
	flag.BoolVar(&cfg.OmFileCompress, "compress", false, "[OUT|FILE] file compress log files?")
	flag.StringVar(&cfg.OmFileName, "filename", "", "[OUT|FILE] file name")

	// elasticsearch forward module
	flag.StringVar(&cfg.OmElasticHost, "eshost", "", "[OUT|ELASTIC] elasticsearch host")
	flag.StringVar(&cfg.OmElasticUser, "esuser", "", "[OUT|ELASTIC] elasticsearch user 'basic auth'")
	flag.StringVar(&cfg.OmElasticPass, "espass", "", "[OUT|ELASTIC] elasticsearch password 'basic auth'")
	flag.StringVar(&cfg.OmElasticCloudAPIKey, "escloudapikey", "", "[OUT|ELASTIC] elasticsearch Base64-encoded token for authorization; if set, overrides username and password")
	flag.StringVar(&cfg.OmElasticCloudID, "escloudid", "", "[OUT|ELASTIC] endpoint for the Elastic Cloud Service 'https://elastic.co/cloud'")
	flag.StringVar(&cfg.OmElasticIndex, "esindex", "canarychirps", "[OUT|ELASTIC] elasticsearch index")

	// kafka forward module
	flag.StringVar(&cfg.OmKafkaBrokers, "kafkabrokers", "", `[OUT|KAFKA] kafka brokers "broker:port"
		for multiple brokers, separate using semicolon "broker1:9092;broker2:9092"`)
	flag.StringVar(&cfg.OmKafkaTopic, "kafkatopic", "", "[OUT|KAFKA] kafka topic 'defaults to canarychirps if not set'")
}

func setDefaultVars(cfg *canarytools.ChirpForwarderConfig, l *log.Logger) {
	if l == nil {
		panic("no logger specififed; will create a new one with default settings")
		// l = log.New()
	}
	switch cfg.Loglevel {
	case "info":
		l.SetLevel(log.InfoLevel)
	case "warning":
		l.SetLevel(log.WarnLevel)
	case "debug":
		l.SetLevel(log.DebugLevel)
	case "trace":
		l.SetLevel(log.TraceLevel)
	default:
		l.Warn("unsupported log level, or none specified; will set to 'info', ")
		l.SetLevel(log.InfoLevel)
	}
	// setting default values for those that doesn't exist
	// had to do it here instead of flag package to support envrionment vars
	switch cfg.ThenWhat {
	case "nothing":
	case "ack":
	default:
		l.Warn("'then' is not valid, or not specified; will set to 'nothing'")
		cfg.ThenWhat = "nothing"
	}

	switch cfg.WhichIncidents {
	case "all":
	case "unacknowledged":
	default:
		l.Warn("'which' is not valid, or not specified; will set to 'unacknowledged'")
		cfg.WhichIncidents = "unacknowledged"
	}

	switch cfg.IncidentFilter {
	case "none":
	case "dropevents":
	default:
		l.Warn("'filter' is not valid, or not specified; will set to 'none'")
		cfg.IncidentFilter = "none"
	}

	if cfg.ImConsoleAPIFetchInterval == 0 {
		l.Warn("'interval' is not valid, or not specified; will set to '60 seconds'")
		cfg.ImConsoleAPIFetchInterval = 60
	}

	// File forward module
	if cfg.OmFileMaxSize == 0 {
		l.Warn("'maxsize' is not valid, or not specified; will set to '8 Megabytes'")
		cfg.OmFileMaxSize = 8
	}
	if cfg.OmFileMaxBackups == 0 {
		l.Warn("'maxbackups' is not valid, or not specified; will set to '14 files'")
		cfg.OmFileMaxBackups = 14
	}
	if cfg.OmFileMaxAge == 0 {
		l.Warn("'maxage' is not valid, or not specified; will set to '120 days'")
		cfg.OmFileMaxAge = 120
	}
	if cfg.OmFileName == "" {
		l.Warn("'filename' is not valid, or not specified; will set to 'canaryChirps.json'")
		cfg.OmFileName = "canaryChirps.json"
	}
}

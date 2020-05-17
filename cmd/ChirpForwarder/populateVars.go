package main

import (
	"flag"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func popultaeVarsFromEnv() {
	// Artificial Intelligence vars parser!
	// lots of if statements ahead ... :)

	// General flags
	if feederModule == "" {
		feederModule, _ = os.LookupEnv("CANARY_FEEDER")
	}
	if forwarderModule == "" {
		forwarderModule, _ = os.LookupEnv("CANARY_OUTPUT")
	}
	if loglevel == "" {
		loglevel, _ = os.LookupEnv("CANARY_LOGLEVEL")
	}
	if thenWhat == "" {
		thenWhat, _ = os.LookupEnv("CANARY_THEN")
	}
	if sinceWhenString == "" {
		sinceWhenString, _ = os.LookupEnv("CANARY_SINCE")
	}
	if whichIncidents == "" {
		whichIncidents, _ = os.LookupEnv("CANARY_WHICH")
	}
	if incidentFilter == "" {
		incidentFilter, _ = os.LookupEnv("CANARY_FILTER")
	}

	// SSL/TLS Client configs
	// used by TCP & Elastic output
	if sslUseSSL == false {
		sslUseSSLBool, _ := os.LookupEnv("CANARY_SSL")
		sslUseSSL, _ = strconv.ParseBool(sslUseSSLBool)
	}
	if sslSkipInsecure == false {
		sslSkipInsecureBool, _ := os.LookupEnv("CANARY_INSECURE")
		sslSkipInsecure, _ = strconv.ParseBool(sslSkipInsecureBool)
	}
	if sslCA == "" {
		sslCA, _ = os.LookupEnv("CANARY_SSLCLIENTCA")
	}
	if sslKey == "" {
		sslKey, _ = os.LookupEnv("CANARY_SSLCLIENTKEY")
	}
	if sslCert == "" {
		sslCert, _ = os.LookupEnv("CANARY_SSLCLIENTCERT")
	}

	// INPUT MODULES
	// Console API input module
	if imConsoleAPIKey == "" {
		imConsoleAPIKey, _ = os.LookupEnv("CANARY_APIKEY")
	}
	if imConsoleAPIDomain == "" {
		imConsoleAPIDomain, _ = os.LookupEnv("CANARY_DOMAIN")
	}
	if imConsoleTokenFile == "" {
		imConsoleTokenFile, _ = os.LookupEnv("CANARY_TOKENFILE")
	}
	if imConsoleAPIFetchInterval == 0 {
		imConsoleAPIFetchIntervalInt, _ := os.LookupEnv("CANARY_INTERVAL")
		imConsoleAPIFetchInterval, _ = strconv.Atoi(imConsoleAPIFetchIntervalInt)
	}

	// OUTPUT MODULES
	// TCP/UDP output module
	if omTCPUDPPort == 0 {
		omTCPUDPPortInt, _ := os.LookupEnv("CANARY_PORT")
		omTCPUDPPort, _ = strconv.Atoi(omTCPUDPPortInt)
	}
	if omTCPUDPHost == "" {
		omTCPUDPHost, _ = os.LookupEnv("CANARY_HOST")
	}

	// File forward module
	if omFileMaxSize == 0 {
		omFileMaxSizeInt, _ := os.LookupEnv("CANARY_MAXSIZE")
		omFileMaxSize, _ = strconv.Atoi(omFileMaxSizeInt)
	}
	if omFileMaxBackups == 0 {
		omFileMaxBackupsInt, _ := os.LookupEnv("CANARY_MAXBACKUPS")
		omFileMaxBackups, _ = strconv.Atoi(omFileMaxBackupsInt)
	}
	if omFileMaxAge == 0 {
		omFileMaxAgeInt, _ := os.LookupEnv("CANARY_MAXAGE")
		omFileMaxAge, _ = strconv.Atoi(omFileMaxAgeInt)
	}
	if omFileCompress == false {
		omFileCompressBool, _ := os.LookupEnv("CANARY_COMPRESS")
		omFileCompress, _ = strconv.ParseBool(omFileCompressBool)
	}
	if omFileName == "" {
		omFileName, _ = os.LookupEnv("CANARY_FILENAME")
	}

	// elasticsearch forward module
	if omElasticHost == "" {
		omElasticHost, _ = os.LookupEnv("CANARY_ESHOST")
	}
	if omElasticUser == "" {
		omElasticUser, _ = os.LookupEnv("CANARY_ESUSER")
	}
	if omElasticPass == "" {
		omElasticPass, _ = os.LookupEnv("CANARY_ESPASS")
	}
	if omElasticCloudAPIKey == "" {
		omElasticCloudAPIKey, _ = os.LookupEnv("CANARY_ESCLOUDAPIKEY")
	}
	if omElasticCloudID == "" {
		omElasticCloudID, _ = os.LookupEnv("CANARY_ESCLOUDID")
	}
	if omElasticIndex == "" {
		omElasticIndex, _ = os.LookupEnv("CANARY_ESINDEX")
	}

	// kafka forward module
	if omKafkaBrokers == "" {
		omKafkaBrokers, _ = os.LookupEnv("CANARY_KAFKABROKERS")
	}
	if omKafkaTopic == "" {
		omKafkaTopic, _ = os.LookupEnv("CANARY_KAFKATOPIC")
	}
}

func populateVarsFromFlags() {
	// General flags
	flag.StringVar(&feederModule, "feeder", "consoleapi", "input module")
	flag.StringVar(&forwarderModule, "output", "", "output module ('tcp', 'file', 'elastic' or 'kafka')")
	flag.StringVar(&loglevel, "loglevel", "", "set loglevel, can be one of ('info', 'warning' or 'debug')")
	flag.StringVar(&thenWhat, "then", "", "what to do after getting an incident? can be one of ('nothing', or 'ack')")
	flag.StringVar(&sinceWhenString, "since", "", `get events newer than this time.
		format has to be like this: 'yyyy-MM-dd HH:mm:ss'
		if nothing provided, it will check value from '.canary.lastcheck' file,
		if .canary.lastcheck file does not exist, it will default to events from last 7 days`)
	flag.StringVar(&whichIncidents, "which", "", "which incidents to fetch? can be one of ('all', or 'unacknowledged')")
	flag.StringVar(&incidentFilter, "filter", "", "filter to apply to incident ('none', or 'dropevents')")

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
	flag.StringVar(&imConsoleTokenFile, "tokenfile", "", "the token file 'canarytools.config' which contains api token and the domain")
	flag.IntVar(&imConsoleAPIFetchInterval, "interval", 0, "alert fetch interval 'in seconds'")

	// OUTPUT MODULES
	// TCP/UDP output module
	flag.IntVar(&omTCPUDPPort, "port", 0, "[OUT|TCP] TCP/UDP port")
	flag.StringVar(&omTCPUDPHost, "host", "", "[OUT|TCP] host")

	// File forward module
	flag.IntVar(&omFileMaxSize, "maxsize", 0, "[OUT|FILE] file max size in megabytes")
	flag.IntVar(&omFileMaxBackups, "maxbackups", 0, "[OUT|FILE] file max number of files to keep")
	flag.IntVar(&omFileMaxAge, "maxage", 0, "[OUT|FILE] file max age in days 'older than this will be deleted'")
	flag.BoolVar(&omFileCompress, "compress", false, "[OUT|FILE] file compress log files?")
	flag.StringVar(&omFileName, "filename", "", "[OUT|FILE] file name")

	// elasticsearch forward module
	flag.StringVar(&omElasticHost, "eshost", "", "[OUT|ELASTIC] elasticsearch host")
	flag.StringVar(&omElasticUser, "esuser", "", "[OUT|ELASTIC] elasticsearch user 'basic auth'")
	flag.StringVar(&omElasticPass, "espass", "", "[OUT|ELASTIC] elasticsearch password 'basic auth'")
	flag.StringVar(&omElasticCloudAPIKey, "escloudapikey", "", "[OUT|ELASTIC] elasticsearch Base64-encoded token for authorization; if set, overrides username and password")
	flag.StringVar(&omElasticCloudID, "escloudid", "", "[OUT|ELASTIC] endpoint for the Elastic Cloud Service 'https://elastic.co/cloud'")
	flag.StringVar(&omElasticIndex, "esindex", "canarychirps", "[OUT|ELASTIC] elasticsearch index")

	// kafka forward module
	flag.StringVar(&omKafkaBrokers, "kafkabrokers", "", `[OUT|KAFKA] kafka brokers "broker:port"
		for multiple brokers, separate using semicolon "broker1:9092;broker2:9092"`)
	flag.StringVar(&omKafkaTopic, "kafkatopic", "", "[OUT|KAFKA] elasticsearch user 'basic auth'")
}

func setDefaultVars(l *log.Logger) {
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
		l.Warn("unsupported log level, or none specified; will set to 'info', ")
		l.SetLevel(log.InfoLevel)
	}
	// setting default values for those that doesn't exist
	// had to do it here instead of flag package to support envrionment vars
	switch thenWhat {
	case "nothing":
	case "ack":
	default:
		l.Warn("'then' is not valid, or not specified; will set to 'nothing'")
		thenWhat = "nothing"
	}

	switch whichIncidents {
	case "all":
	case "unacknowledged":
	default:
		l.Warn("'which' is not valid, or not specified; will set to 'unacknowledged'")
		whichIncidents = "unacknowledged"
	}

	switch incidentFilter {
	case "none":
	case "dropevents":
	default:
		l.Warn("'filter' is not valid, or not specified; will set to 'none'")
		incidentFilter = "none"
	}

	if imConsoleAPIFetchInterval == 0 {
		l.Warn("'interval' is not valid, or not specified; will set to '60 seconds'")
		imConsoleAPIFetchInterval = 60
	}

	// File forward module
	if omFileMaxSize == 0 {
		l.Warn("'maxsize' is not valid, or not specified; will set to '8 Megabytes'")
		omFileMaxSize = 8
	}
	if omFileMaxBackups == 0 {
		l.Warn("'maxbackups' is not valid, or not specified; will set to '14 files'")
		omFileMaxBackups = 14
	}
	if omFileMaxAge == 0 {
		l.Warn("'maxage' is not valid, or not specified; will set to '120 days'")
		omFileMaxAge = 120
	}
	if omFileName == "" {
		l.Warn("'filename' is not valid, or not specified; will set to 'canaryChirps.json'")
		omFileName = "canaryChirps.json"
	}
}

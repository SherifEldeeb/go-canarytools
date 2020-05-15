package main

import (
	"flag"
	"os"
	"strconv"
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
}

func populateVarsFromFlags() {
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

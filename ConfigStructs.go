package canarytools

import "crypto/tls"

// ChirpForwarderConfig contains configs for the forwarder
type ChirpForwarderConfig struct {
	// General flags
	GeneralConfig

	// SSL/TLS Client configs
	// used by TCP & Elastic output
	SSLConfig

	// INPUT MODULES
	// Console API input module
	ConsoleAPIConfig

	// OUTPUT MODULES
	// TCP/UDP output module
	TCPOutConfig

	// File forward module
	FileOutConfig

	// elasticsearch forward module
	ElasticOutConfig

	// kafka forward module
	KafkaOutConfig

	// TLS config
	TLSConfig *tls.Config
}

// GeneralConfig contains general configs
type GeneralConfig struct {
	// General flags
	FeederModule    string // CANARY_FEEDER
	ForwarderModule string // CANARY_OUTPUT
	Loglevel        string // CANARY_LOGLEVEL
	ThenWhat        string // CANARY_THEN
	SinceWhenString string // CANARY_SINCE
	WhichIncidents  string // CANARY_WHICH
	IncidentFilter  string // CANARY_FILTER
}

// SSLConfig contains SSL related configs
type SSLConfig struct {
	// SSL/TLS Client configs
	// used by TCP & Elastic output
	SSLUseSSL       bool   // CANARY_SSL
	SSLSkipInsecure bool   // CANARY_INSECURE
	SSLCA           string // CANARY_SSLCLIENTCA
	SSLKey          string // CANARY_SSLCLIENTKEY
	SSLCert         string // CANARY_SSLCLIENTCERT
}

// ConsoleAPIConfig contains configs for Console API input module
type ConsoleAPIConfig struct {
	// Console API input module
	ImConsoleAPIKey           string // CANARY_APIKEY
	ImConsoleAPIDomain        string // CANARY_DOMAIN
	ImConsoleTokenFile        string // CANARY_TOKENFILE
	ImConsoleAPIFetchInterval int    // CANARY_INTERVAL
}

// TCPOutConfig contains configs for the TCP output module
type TCPOutConfig struct {
	// TCP output module
	OmTCPUDPPort int    // CANARY_PORT
	OmTCPUDPHost string // CANARY_HOST
}

// FileOutConfig contains configs for the File forward module
type FileOutConfig struct {
	// File forward module
	OmFileMaxSize    int    // CANARY_MAXSIZE
	OmFileMaxBackups int    // CANARY_MAXBACKUPS
	OmFileMaxAge     int    // CANARY_MAXAGE
	OmFileCompress   bool   // CANARY_COMPRESS
	OmFileName       string // CANARY_FILENAME
}

// ElasticOutConfig contains configs for the elasticsearch forward module
type ElasticOutConfig struct {
	// elasticsearch forward module
	OmElasticHost        string // CANARY_ESHOST
	OmElasticUser        string // CANARY_ESUSER
	OmElasticPass        string // CANARY_ESPASS
	OmElasticCloudAPIKey string // CANARY_ESCLOUDAPIKEY
	OmElasticCloudID     string // CANARY_ESCLOUDID
	OmElasticIndex       string // CANARY_ESINDEX
}

// KafkaOutConfig contains configs for the kafka forward module
type KafkaOutConfig struct {
	// kafka forward module
	OmKafkaBrokers string // CANARY_KAFKABROKERS
	OmKafkaTopic   string // CANARY_KAFKATOPIC
}

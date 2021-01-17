package canarytools

import "crypto/tls"

// ConsoleAPIConfig contains configs for Console API input module
type ConsoleAPIConfig struct {
	// Console domain hash
	ConsoleAPIDomain string // CANARY_DOMAIN

	// if using console API (can't be used with factory)
	ConsoleAPIKey    string // CANARY_APIKEY
	ConsoleTokenFile string // CANARY_TOKENFILE

	// if using factory auth (can't be used with console api)
	ConsoleFactoryAuth string // CANARY_FACTORYAUTH
	FactoryAuthFile    string // CANARY_FACTORYAUTHFILE

	// OpMode can be "api" or "factory"
	// this should be automatically set by NewClient
	OpMode string

	// TODO: Move to consoleAPIFeeder
	ConsoleAPIFetchInterval int // CANARY_INTERVAL
}

// TokenDropperConfig contains configs for the TokenDropper
type TokenDropperConfig struct {
	ConsoleAPIConfig
	GeneralTokenDropperConfig
}

// GeneralTokenDropperConfig contains general configs for TokenDropper
type GeneralTokenDropperConfig struct {
	// General flags
	FilesCount                 int      // number of files per directory
	RandYearsBack              int      // Randomize dates between Now() and 'years' back
	LocalTokenProxy            bool     // start as a local token proxy?
	DropWhere                  string   // where to drop tokens?
	KindsStr                   string   // comma-separated string with what kind of tokens to drop
	Kinds                      []string // what kind of tokens to drop
	LogLevel                   string   // loglevel
	FlockName                  string   // Name of the flock
	FlockID                    string   // Flock ID
	CreateFlockIfNotExists     bool     // should we create the flock if it didn't exist?
	CreateDirectoryIfNotExists bool     // should we create the directory (DropWhere) if it didn't exist?
	CustomMemo                 string   // custom memo to be added to the default one
	FileName                   string   // the filename of the token, if this is set, count will be one, and there will be some checks to make sure the extension matchs the kind
	RandomizeFilenames         bool     // add random text to filenames to make them unique
	OverwriteFileIfExists      bool     // if a file with same name exists, should we overwrite it?
}

// CanaryDeleterConfig contains configs for the CanaryDeleter
type CanaryDeleterConfig struct {
	ConsoleAPIConfig
	GeneralCanaryDeleterConfig
}

// GeneralCanaryDeleterConfig contains general configs for CanaryDeleter
type GeneralCanaryDeleterConfig struct {
	DeleteWhat            string // what to delete? either alerts or tokens
	IncidentsState        string // State of Incidents ... valid options are "all", "acknowledged" and "unacknowledged"
	FlockName             string // Name of the flock
	FlockID               string // Flock ID
	NodeID                string // Node ID
	FilterType            string // filter using flock or node?
	LogLevel              string
	IncludeUnacknowledged bool
	DumpToJson            bool
	DumpOnly              bool // setting this to true will NOT delete incidents
}

// ChirpForwarderConfig contains configs for the forwarder
type ChirpForwarderConfig struct {
	// General flags
	GeneralChirpForwarderConfig

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

	// SQS forward module
	SQSOutConfig

	// TLS config
	TLSConfig *tls.Config
}

// GeneralChirpForwarderConfig contains general configs for ChirpForwarder
type GeneralChirpForwarderConfig struct {
	// General flags
	FeederModule    string // CANARY_FEEDER
	ForwarderModule string // CANARY_OUTPUT
	Loglevel        string // CANARY_LOGLEVEL
	ThenWhat        string // CANARY_THEN
	SinceWhenString string // CANARY_SINCE
	WhichIncidents  string // CANARY_WHICH
	IncidentFilter  string // CANARY_FILTER
	FlockName       string // CANARY_FLOCKNAME
	FlockID         string // CANARY_FLOCKID
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

// SQSOutConfig contains configs for the kafka forward module
type SQSOutConfig struct {
	// SQS forward module
	OmSQSQueueName string // CANARY_SQSQUEUENAME
}

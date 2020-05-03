package canarytools

// Mostly copy/paste from ref: https://canarytools.readthedocs.io/en/latest/incident_attributes.html

// HTTPEvent is an HTTP Canarytoken triggered Event
type HTTPEvent struct {
	Description string                 `json:"description,omitempty"` // Canarytoken triggered
	Type        string                 `json:"type,omitempty"`        // http
	Canarytoken string                 `json:"canarytoken,omitempty"` // Unique string that acts as the Canarytoken
	Headers     map[string]interface{} `json:"headers,omitempty"`     // Headers is a dict, Only present for HTTP Canarytokens.
	URL         string                 `json:"url,omitempty"`         // URL of the HTTP Canarytoken.
	Logtype     string                 `json:"logtype,omitempty"`     // 17000
}

// WebImageEvent is a Remote Web Image triggered Event
type WebImageEvent struct {
	Description  string                 `json:"description,omitempty"`    // "Remote Web Image"
	Type         string                 `json:"type,omitempty"`           // "web-image"
	Canarytoken  string                 `json:"canarytoken,omitempty"`    // Unique string that acts as the Canarytoken
	Headers      map[string]interface{} `json:"headers,omitempty"`        // Headers is a dict, Only present for HTTP Canarytokens.
	URL          string                 `json:"url,omitempty"`            // URL of the HTTP Canarytoken.
	Logtype      string                 `json:"logtype,omitempty"`        // "17001"
	WebImage     string                 `json:"web_image,omitempty"`      // Byte string of the web image
	WebImageType string                 `json:"web_image_type,omitempty"` // Type of the web image
	WebImageName string                 `json:"web_image_name,omitempty"` // Name of the web image
}

// DocMSWordEvent is a MS Word triggered Event
type DocMSWordEvent struct {
	Description string                 `json:"description,omitempty"` // "Canarytoken triggered"
	Type        string                 `json:"type,omitempty"`        // "doc-msword"
	Canarytoken string                 `json:"canarytoken,omitempty"` // Unique string that acts as the Canarytoken
	Headers     map[string]interface{} `json:"headers,omitempty"`     // Headers is a dict, Only present for HTTP Canarytokens.
	URL         string                 `json:"url,omitempty"`         // URL of the HTTP Canarytoken.
	Logtype     string                 `json:"logtype,omitempty"`     // "17002"
	Doc         string                 `json:"doc,omitempty"`         // Byte String of the tokened document
	DocName     string                 `json:"doc_name,omitempty"`    // Name of the document tokened (or created)
	DocType     string                 `json:"doc_type,omitempty"`    // Type of document chosen (doc or docx)
}

// ClonedWebEvent is a Cloned Website triggered Event
type ClonedWebEvent struct {
	Description string                 `json:"description,omitempty"` // "Cloned Website"
	Type        string                 `json:"type,omitempty"`        // "cloned-web"
	Canarytoken string                 `json:"canarytoken,omitempty"` // Unique string that acts as the Canarytoken
	Headers     map[string]interface{} `json:"headers,omitempty"`     // Headers is a dict, Only present for HTTP Canarytokens.
	URL         string                 `json:"url,omitempty"`         // URL of the HTTP Canarytoken.
	Logtype     string                 `json:"logtype,omitempty"`     // "17003"
	ClonedWeb   string                 `json:"cloned_web,omitempty"`  // Domain that we are tokening
}

// AWSS3Event is an Amazon AWS S3  triggered Event
type AWSS3Event struct {
	Description    string                 `json:"description,omitempty"`       // "Amazon S3"
	Type           string                 `json:"type,omitempty"`              // "aws-s3"
	Canarytoken    string                 `json:"canarytoken,omitempty"`       // Unique string that acts as the Canarytoken
	Headers        map[string]interface{} `json:"headers,omitempty"`           // Headers is a dict, Only present for HTTP Canarytokens.
	URL            string                 `json:"url,omitempty"`               // URL of the HTTP Canarytoken.
	Logtype        string                 `json:"logtype,omitempty"`           // "17005"
	S3SourceBucket string                 `json:"s_3_source_bucket,omitempty"` // bucket that we are tokening
	S3LogBucket    string                 `json:"s_3_log_bucket,omitempty"`    // bucket where logging to stored and monitored
	Online         string                 `json:"online,omitempty"`            // Whether the token is online or not
}

// GoogleDocsEvent is a Google Docs triggered Event
type GoogleDocsEvent struct {
	Description  string                 `json:"description,omitempty"`   // "Google Document"
	Type         string                 `json:"type,omitempty"`          // "google-docs"
	Canarytoken  string                 `json:"canarytoken,omitempty"`   // Unique string that acts as the Canarytoken
	Headers      map[string]interface{} `json:"headers,omitempty"`       // Headers is a dict, Only present for HTTP Canarytokens.
	URL          string                 `json:"url,omitempty"`           // URL of the HTTP Canarytoken.
	Logtype      string                 `json:"logtype,omitempty"`       // "17006"
	DocsLink     string                 `json:"docs_link,omitempty"`     // URL to the google doc
	EmailLink    string                 `json:"email_link,omitempty"`    // URL used for email
	DocumentName string                 `json:"document_name,omitempty"` // Name of the document
}

// GoogleSheetsEvent is a Google Docs triggered Event
type GoogleSheetsEvent struct {
	Description  string                 `json:"description,omitempty"`   // "Google Document"
	Type         string                 `json:"type,omitempty"`          // "google-sheets
	Canarytoken  string                 `json:"canarytoken,omitempty"`   // Unique string that acts as the Canarytoken
	Headers      map[string]interface{} `json:"headers,omitempty"`       // Headers is a dict, Only present for HTTP Canarytokens.
	URL          string                 `json:"url,omitempty"`           // URL of the HTTP Canarytoken.
	Logtype      string                 `json:"logtype,omitempty"`       // "17006"
	DocsLink     string                 `json:"docs_link,omitempty"`     // URL to the google doc
	EmailLink    string                 `json:"email_link,omitempty"`    // URL used for email
	DocumentName string                 `json:"document_name,omitempty"` // Name of the document
}

// QRCodeEvent is a QR Code triggered Event
type QRCodeEvent struct {
	Description string                 `json:"description,omitempty"` //  QR Code
	Type        string                 `json:"type,omitempty"`        //  qr-code
	Canarytoken string                 `json:"canarytoken,omitempty"` // Unique string that acts as the Canarytoken
	Headers     map[string]interface{} `json:"headers,omitempty"`     // Headers is a dict, Only present for HTTP Canarytokens.
	URL         string                 `json:"url,omitempty"`         // URL of the HTTP Canarytoken.
	Logtype     string                 `json:"logtype,omitempty"`     //  17009
	QRCode      string                 `json:"qr_code,omitempty"`     // Byte string of the tokened QR code

}

// SVNEvent is a SVN triggered Event
type SVNEvent struct {
	Description string                 `json:"description,omitempty"` //  SVN Repo
	Type        string                 `json:"type,omitempty"`        //  svn
	Canarytoken string                 `json:"canarytoken,omitempty"` // Unique string that acts as the Canarytoken
	Headers     map[string]interface{} `json:"headers,omitempty"`     // Headers is a dict, Only present for HTTP Canarytokens.
	URL         string                 `json:"url,omitempty"`         // URL of the HTTP Canarytoken.
	Logtype     string                 `json:"logtype,omitempty"`     //  17010

}

// SQLEvent is a SQL triggered Event
type SQLEvent struct {
	Description  string                 `json:"description,omitempty"`   //  SQL Server
	Type         string                 `json:"type,omitempty"`          //  sql
	Canarytoken  string                 `json:"canarytoken,omitempty"`   // Unique string that acts as the Canarytoken
	Headers      map[string]interface{} `json:"headers,omitempty"`       // Headers is a dict, Only present for HTTP Canarytokens.
	URL          string                 `json:"url,omitempty"`           // URL of the HTTP Canarytoken.
	Logtype      string                 `json:"logtype,omitempty"`       //  17011
	TriggerType  string                 `json:"trigger_type,omitempty"`  // SQL trigger type (SELECT, UPDATE, INSERT, DELETE)
	TableName    string                 `json:"table_name,omitempty"`    // SQL table name (trigger_type: UPDATE, INSERT,DELETE)
	TriggerName  string                 `json:"trigger_name,omitempty"`  // SQL trigger name (trigger_type: UPDATE, INSERT,DELETE)
	ViewName     string                 `json:"view_name,omitempty"`     // SQL View name (trigger_type: SELECT)
	FunctionName string                 `json:"function_name,omitempty"` // SQL function name (trigger_type: SELECT)

}

// AWSIDEvent is a AWS ID triggered Event
type AWSIDEvent struct {
	Description     string                 `json:"description,omitempty"`       //  Amazon API Key
	Type            string                 `json:"type,omitempty"`              //  aws-id
	Canarytoken     string                 `json:"canarytoken,omitempty"`       // Unique string that acts as the Canarytoken
	Headers         map[string]interface{} `json:"headers,omitempty"`           // Headers is a dict, Only present for HTTP Canarytokens.
	URL             string                 `json:"url,omitempty"`               // URL of the HTTP Canarytoken.
	Logtype         string                 `json:"logtype,omitempty"`           //  17012
	SecretAccessKey string                 `json:"secret_access_key,omitempty"` // AWS generated secret access key
	AccessKeyID     string                 `json:"access_key_id,omitempty"`     // AWS generated access key ID

}

// FastRedirectEvent is a Fast Redirect triggered Event
type FastRedirectEvent struct {
	Description        string                 `json:"description,omitempty"`          //  Fast HTTP Redirect
	Type               string                 `json:"type,omitempty"`                 //  fast-redirect
	Canarytoken        string                 `json:"canarytoken,omitempty"`          // Unique string that acts as the Canarytoken
	Headers            map[string]interface{} `json:"headers,omitempty"`              // Headers is a dict, Only present for HTTP Canarytokens.
	URL                string                 `json:"url,omitempty"`                  // URL of the HTTP Canarytoken.
	Logtype            string                 `json:"logtype,omitempty"`              //  17016
	BrowserRedirectURL string                 `json:"browser_redirect_url,omitempty"` // Original URL attempted before redirect

}

// SlowRedirectEvent is a Slow Redirect triggered Event
type SlowRedirectEvent struct {
	Description        string                 `json:"description,omitempty"`          //  Slow HTTP Redirect
	Type               string                 `json:"type,omitempty"`                 //  slow-redirect
	Canarytoken        string                 `json:"canarytoken,omitempty"`          // Unique string that acts as the Canarytoken
	Headers            map[string]interface{} `json:"headers,omitempty"`              // Headers is a dict, Only present for HTTP Canarytokens.
	URL                string                 `json:"url,omitempty"`                  // URL of the HTTP Canarytoken.
	Logtype            string                 `json:"logtype,omitempty"`              //  17017
	BrowserRedirectURL string                 `json:"browser_redirect_url,omitempty"` // Original URL attempted before redirect

}

// DNSEvent is a DNS triggered Event
type DNSEvent struct {
	Description string                 `json:"description,omitempty"` //  DNS
	Type        string                 `json:"type,omitempty"`        //  dns
	Canarytoken string                 `json:"canarytoken,omitempty"` // Unique string that acts as the Canarytoken
	Hostname    map[string]interface{} `json:"hostname,omitempty"`    // Hostname of the DNS Canarytoken.
	Logtype     string                 `json:"logtype,omitempty"`     //  16000

}

// DesktopINIEvent is a Desktop ini triggered Event
type DesktopINIEvent struct {
	Description string                 `json:"description,omitempty"` //  Windows Directory Browsing
	Type        string                 `json:"type,omitempty"`        //  windows-dir
	Canarytoken string                 `json:"canarytoken,omitempty"` // Unique string that acts as the Canarytoken
	Hostname    map[string]interface{} `json:"hostname,omitempty"`    // Hostname of the DNS Canarytoken.
	Logtype     string                 `json:"logtype,omitempty"`     //  16006

}

// AdobeReaderPDFEvent is a Adobe Reader PDF triggered Event
type AdobeReaderPDFEvent struct {
	Description string                 `json:"description,omitempty"` //  Acrobat Reader PDF Document
	Type        string                 `json:"type,omitempty"`        //  pdf-acrobat-reader
	Canarytoken string                 `json:"canarytoken,omitempty"` // Unique string that acts as the Canarytoken
	Hostname    map[string]interface{} `json:"hostname,omitempty"`    // Hostname of the DNS Canarytoken.
	Logtype     string                 `json:"logtype,omitempty"`     //  16008

}

// MSWordDocMacroedEvent is a MS Word Doc Macroed triggered Event
type MSWordDocMacroedEvent struct {
	Description string                 `json:"description,omitempty"` //  MS Word .docm Document
	Type        string                 `json:"type,omitempty"`        //  msword-macro
	Canarytoken string                 `json:"canarytoken,omitempty"` // Unique string that acts as the Canarytoken
	Hostname    map[string]interface{} `json:"hostname,omitempty"`    // Hostname of the DNS Canarytoken.
	Logtype     string                 `json:"logtype,omitempty"`     //  16009
	Doc         string                 `json:"doc,omitempty"`         // Byte String of the tokened document
	DocName     string                 `json:"doc_name,omitempty"`    // Name of the document
	DocType     string                 `json:"doc_type,omitempty"`    // Type of document chosen (doc or docx)

}

// MSExcelMacroedEvent is a MS Excel Macroed triggered Event
type MSExcelMacroedEvent struct {
	Description string                 `json:"description,omitempty"` //  MS Excel .xlsm Document
	Type        string                 `json:"type,omitempty"`        //  msexcel-macro
	Canarytoken string                 `json:"canarytoken,omitempty"` // Unique string that acts as the Canarytoken
	Hostname    map[string]interface{} `json:"hostname,omitempty"`    // Hostname of the DNS Canarytoken.
	Logtype     string                 `json:"logtype,omitempty"`     //  16010
	Doc         string                 `json:"doc,omitempty"`         // Byte String of the tokened document
	DocName     string                 `json:"doc_name,omitempty"`    // Name of the document
	DocType     string                 `json:"doc_type,omitempty"`    // Type of document chosen (doc or docx)
}

// Port Scan Events
// There are many types of port scans incidents.
//     A host port scan occurs when a single Canary is port scanned by a single source.
//     A consolidated network port scan occurs when multiple Canaries are scanned by a single source.
//     An NMAP NULL scan was run against the Canary.
//     An NMAP OS scan was run against the Canary.
//     An NMAP XMAS scan was run against the Canary.

// HostPortScanEvent is a Host Port Scan triggered Event
type HostPortScanEvent struct {
	Description string        `json:"description,omitempty"` // Host Port Scan
	Ports       []interface{} `json:"ports,omitempty"`       //A list of ports scanned
	Logtype     string        `json:"logtype,omitempty"`     // 5003
}

// ConsolidatedNetworkPortScanEvent is a Consolidated Network Port Scan triggered Event
type ConsolidatedNetworkPortScanEvent struct {
	Description  string                 `json:"description,omitempty"`   // Host Port Scan
	PortsScanned map[string]interface{} `json:"ports_scanned,omitempty"` //A dictionary of ports scanned and the IP address of the Canaries on which the scan occurred.
	Logtype      string                 `json:"logtype,omitempty"`       // 5007
}

// NMAPNULLScanEvent is a NMAP NULL Scan triggered Event
type NMAPNULLScanEvent struct {
	Description string `json:"description,omitempty"` // NMAP NULL Scan Detected
	Logtype     string `json:"logtype,omitempty"`     // 5005
}

// NMAPOSScanEvent is a NMAP OS Scan triggered Event
type NMAPOSScanEvent struct {
	Description string `json:"description,omitempty"` // NMAP OS Scan Detected
	Logtype     string `json:"logtype,omitempty"`     // 5004
}

// NMAPXMASScanEvent is a NMAP XMAS Scan triggered Event
type NMAPXMASScanEvent struct {
	Description string `json:"description,omitempty"` // NMAP XMAS Scan Detected
	Logtype     string `json:"logtype,omitempty"`     // 5006
}

// NMAPFINScanEvent is a NMAP FIN Scan triggered Event
type NMAPFINScanEvent struct {
	Description string `json:"description,omitempty"` // NMAP FIN Scan Detected
	Logtype     string `json:"logtype,omitempty"`     // 5008
}

// Canary Connection Events

// CanaryDisconnectedEvent is a Canary Disconnected triggered Even
type CanaryDisconnectedEvent struct {
	Description string `json:"description,omitempty"` // Canary Disconnected
	Logtype     string `json:"logtype,omitempty"`     // 1004
}

// FTPIncidentEvent is a FTP Incident triggered Even
type FTPIncidentEvent struct {
	Description string        `json:"description,omitempty"` // FTP Login Attempt
	Logtype     string        `json:"logtype,omitempty"`     // 2000
	Username    []interface{} `json:"username,omitempty"`    // Attacker supplied username.
	Password    []interface{} `json:"password,omitempty"`    // Attacked supplied password.
}

// GitRepositoryCloneAttemptEvent is a Git Repository Clone Attempt triggered Even
// Triggered when an attacker connects to the Canary git service and attempts any repo clone.
type GitRepositoryCloneAttemptEvent struct {
	Description string        `json:"description,omitempty"` // Git Repository Clone Attempt
	Logtype     string        `json:"logtype,omitempty"`     // 19001
	Host        []interface{} `json:"host,omitempty"`        // Git client's view of the Canary's hostname.
	Repo        string        `json:"repo,omitempty"`        // Name of the repository the client attempted to clone.
}

// HTTP Incidents
// Two types of HTTP Incidents:
//     Page loads, triggered by GET requests. They are disabled by default as they're noisy, and needs to be specifically enabled.
//     Login attempts, triggered by GET requests. They are always enabled.

// HTTPPageLoadEvent is a HTTP Page Load triggered Event
type HTTPPageLoadEvent struct {
	Description string        `json:"description,omitempty"` // HTTP Page Load
	Logtype     string        `json:"logtype,omitempty"`     // 3000
	Path        []interface{} `json:"path,omitempty"`        // Web path requested by the source.
	Useragent   string        `json:"useragent,omitempty"`   // Useragent of the source's browser.
	Channel     string        `json:"channel,omitempty"`     // Optional. Set to 'TLS' if an encrypted site is configured, otherwise absent.
}

// HTTPLoginAttemptEvent is a HTTP Login Attempt triggered Event
type HTTPLoginAttemptEvent struct {
	Description string        `json:"description,omitempty"` // HTTP Login Attempt
	Logtype     string        `json:"logtype,omitempty"`     // 3001
	Username    string        `json:"username,omitempty"`    // Attack supplied username.
	Password    string        `json:"password,omitempty"`    // Attacked supplied password.
	Path        []interface{} `json:"path,omitempty"`        // Web path requested by the source.
	Useragent   string        `json:"useragent,omitempty"`   // Useragent of the source's browser.
	Channel     string        `json:"channel,omitempty"`     // Optional. Set to 'TLS' if an encrypted site is configured, otherwise absent.
}

// HTTPProxyRequestEvent is a HTTP Proxy Request triggered Event
// Triggered by any request through the HTTP proxy module.
type HTTPProxyRequestEvent struct {
	Description string `json:"description,omitempty"` // HTTP Proxy Request
	Logtype     string `json:"logtype,omitempty"`     // 7001
	Username    string `json:"username,omitempty"`    // Attack supplied username.
	Password    string `json:"password,omitempty"`    // Attacked supplied password.
	URL         string `json:"url,omitempty"`         // URL requested by the source.
	Useragent   string `json:"useragent,omitempty"`   // Useragent of the source's browser.
}

// MicrosoftSQLServerLoginAttemptEvent is a Microsoft SQL Server Login Attempt triggered Event
// Triggered by any attempt to authenticate to the Microsoft SQL Server module.
// SQL Server supports multiple authentication modes, and the fields that come through depend on the mode.
type MicrosoftSQLServerLoginAttemptEvent struct {
	Description string `json:"description,omitempty"` // MSSQL Login Attempt
	Logtype     string `json:"logtype,omitempty"`     // 9001 - SQL Server authentication.  9002 - Windows authentication.
	Username    string `json:"username,omitempty"`    // Attack supplied username.
	Password    string `json:"password,omitempty"`    // Optional. Attacker supplied database password
	Hostname    string `json:"hostname,omitempty"`    // Optional. Attacker supplied hostname.
	Domainname  string `json:"domainname,omitempty"`  // Optional. Attacker supplied Active Directory name.
}

// ModBusRequestEvent is a ModBus Request triggered Event
// Triggered by any valid ModBus request.
type ModBusRequestEvent struct {
	Description string `json:"description,omitempty"` // ModBus Request
	Logtype     string `json:"logtype,omitempty"`     // 18001 (Modbus Query Function); 18002 (Modbus Read Function); 18003 (Modbus Write Function)
	UnitID      string `json:"unit_id,omitempty"`     // ModBus unit target.
	FuncCode    string `json:"func_code,omitempty"`   // ModBus function code.
	FuncName    string `json:"func_name,omitempty"`   // Optional. ModBus function name, if available.
	SfuncCode   string `json:"sfunc_code,omitempty"`  // Optional. ModBus subfunction code, if available.
	SfuncName   string `json:"sfunc_name,omitempty"`  // Optional. ModBus subfunction name, if available.
}

// NTPMonlistRequestEvent is a NTP Monlist Request triggered Event
// Triggered by the NTP Monlist command.
type NTPMonlistRequestEvent struct {
	Description string `json:"description,omitempty"` // NTP Monlist Request
	Logtype     string `json:"logtype,omitempty"`     // 11001
	NTPCmd      string `json:"ntp_cmd,omitempty"`     // Name of the NTP command sent. Currently is 'monlist'.
	ClientHash  string `json:"client_hash,omitempty"` // Attacker supplied database password hash.
	Salt        string `json:"salt,omitempty"`        // Attacker supplied database password hash salt.
	Password    string `json:"password,omitempty"`    // Recovered password if possible, otherwise '<Password not in common list>'
}

// RedisCommandEvent is a Redis Command triggered Event
// Triggered by an attacker connecting to the Redis service and issuing valid Redis commands.
type RedisCommandEvent struct {
	Description string `json:"description,omitempty"` // Redis Command
	Logtype     string `json:"logtype,omitempty"`     // 21001
	Cmd         string `json:"cmd,omitempty"`         // Redis command issued by the attacker.
	Args        string `json:"args,omitempty"`        // Arguments to the command.
}

// SIPRequestEvent is a SIP Request triggered Event
// Triggered by an attacker connecting to the SIP service and issuing valid SIP request.
type SIPRequestEvent struct {
	Description string                 `json:"description,omitempty"` // SIP Request
	Logtype     string                 `json:"logtype,omitempty"`     // 15001
	Headers     map[string]interface{} `json:"headers,omitempty"`     // Dict of the SIP headers included in the request.
	Args        string                 `json:"args,omitempty"`        // Arguments to the command.
}

// SharedFileOpenedEvent is a Shared File Opened triggered Event
// Triggered by the opening of a file on the Canary's Windows File Share.
type SharedFileOpenedEvent struct {
	Description string `json:"description,omitempty"` // Shared File Opened
	Logtype     string `json:"logtype,omitempty"`     // 5000
	User        string `json:"user,omitempty"`        // Username supplied by the attacker.
	Filename    string `json:"filename,omitempty"`    // Name of file on the Canary that was accessed.
	Auditaction string `json:"auditaction,omitempty"` // Type of file action. Currently only 'pread'.
	Domain      string `json:"domain,omitempty"`      // Name of domain or workgroup.
	Localname   string `json:"localname,omitempty"`   // Windows Name of Canary machine.
	Mode        string `json:"mode,omitempty"`        // 'workgroup' or 'domain'
	Offset      string `json:"offset,omitempty"`      // Starting position of the read.
	Remotename  string `json:"remotename,omitempty"`  // Windows Name of client machine.
	Sharename   string `json:"sharename,omitempty"`   // Name of the share on which the file resides.
	Size        string `json:"size,omitempty"`        // Amount of bytes read.
	SMBarch     string `json:"smbarch,omitempty"`     // Guess of the remote machine's Windows version.
	SMBver      string `json:"smbver,omitempty"`      // Version of the SMB protocol that was used.
	Status      string `json:"status,omitempty"`      // Result of the file read. Currently only 'ok'.
}

// SNMPRequestEvent is a SNMP Request triggered Event
// Triggered by an incoming SNMP query against the Canary.
type SNMPRequestEvent struct {
	Description     string `json:"description,omitempty"`      // SNMP Request
	Logtype         string `json:"logtype,omitempty"`          // 13001
	CommunityString string `json:"community_string,omitempty"` // SNMP community string supplied by attacker.
	Requests        string `json:"requests,omitempty"`         // SNMP OID requested by the attacker.
}

// SSHLoginAttemptEvent is a SSH Login Attempt triggered Event
// Triggered by an attempt to login to the Canary using SSH. Both password-based and key-based authentication is possible.
// It is also possible to configure  Watched Credentials, which says to only alert if the attacker-supplied credentials match a configured list.
type SSHLoginAttemptEvent struct {
	Description        string `json:"description,omitempty"`         // SSH Login Attempt
	Logtype            string `json:"logtype,omitempty"`             // 4002        username string //SSH username
	Password           string `json:"password,omitempty"`            // SSH password
	Localversion       string `json:"localversion,omitempty"`        // SSH server string supplied by canary.
	Remoteversion      string `json:"remoteversion,omitempty"`       // SSH client string supplied by the attacker.
	Key                string `json:"key,omitempty"`                 // SSH key used by attacker.
	WatchedCredentials string `json:"watched_credentials,omitempty"` // Indicates whether this an SSH key watched for.
}

// CustomTCPServiceRequestEvent is a Custom TCP Service Request triggered Event
// The Custom TCP Service module let's the Canary administrator create simple services that either immediately print a banner on connection, or wait for the client to send data before responding.
type CustomTCPServiceRequestEvent struct {
	Description string `json:"description,omitempty"` // Custom TCP Service Request
	BannerID    string `json:"banner_id,omitempty"`   // Multiple banners are supported, the id identifies which banner service was triggered.
	Data        string `json:"data,omitempty"`        // Optional. Attacker's supplied data.
	Function    string `json:"function,omitempty"`    // Indicates which trigger fired, either 'DATA_RECEIVED' for when a banner was sent after the attacker sent data, or 'CONNECTION_MADE' for when a banner was sent immediately on connection.
	Logtype     string `json:"logtype,omitempty"`     // 20001 (Banner set immediately on connection);  20002 (Banner sent after client sent a line)
}

// TFTPRequestEvent is a TFTP Request triggered Event
// Triggered by a TFTP request against the Canary.
type TFTPRequestEvent struct {
	Description string `json:"description,omitempty"` // TFTP Request
	Logtype     string `json:"logtype,omitempty"`     // 10001
	Filename    string `json:"filename,omitempty"`    // Name of file the attacker tried to act on.
	Opcode      string `json:"opcode,omitempty"`      // File action, either 'READ' or 'WRITE'
}

// TelnetLoginAttemptEvent is a Telnet Login Attempt triggered Event
// Triggered by a Telnet authentication attempt.
type TelnetLoginAttemptEvent struct {
	Description string `json:"description,omitempty"` // Telnet Login Attempt
	Logtype     string `json:"logtype,omitempty"`     // 6001
	Username    string `json:"username,omitempty"`    // Attacker supplied username.
	Password    string `json:"password,omitempty"`    // Attacker supplied password.
}

// VNCLoginAttemptEvent is a VNC Login Attempt triggered Event
// Triggered by an attempt to login to Canary's password protected VNC service.
// VNC passwords are not transmitted in the clear. Instead a hashed version is sent. The Canary will test the hashed password against a handful of common passwords to guess the password, but the hash parameters are also reported so the administrator can crack the hash on more powerful rigs.
type VNCLoginAttemptEvent struct {
	Description     string `json:"description,omitempty"`      // VNC Login Attempt
	Logtype         string `json:"logtype,omitempty"`          // 12001
	Password        string `json:"password,omitempty"`         // Cracked password if very weak.
	ServerChallenge string `json:"server_challenge,omitempty"` // VNC password hashing parameter.
	ClientResponse  string `json:"client_response,omitempty"`  // VNC password hashing parameter.
}

// Settings change events

// ConsoleSettingsChangedEvent is a Console Settings Changed triggered Event
// Triggered by a Canary console setting being changed.
type ConsoleSettingsChangedEvent struct {
	Description string `json:"description,omitempty"` // Console Settings Changed
	Logtype     string `json:"logtype,omitempty"`     // 23001
	Settings    string `json:"settings,omitempty"`    // the settings that were changed.
}

// DeviceSettingsChangedEvent is a Device Settings Changed triggered Event
// Triggered by a Canary's settings being changed.
type DeviceSettingsChangedEvent struct {
	Description string `json:"description,omitempty"` // Device Settings Changed
	Logtype     string `json:"logtype,omitempty"`     // 23002
	Settings    string `json:"settings,omitempty"`    // the settings that were changed.
}

// FlockSettingsChangedEvent is a Flock Settings Changed triggered Event
// Triggered by a flock's settings being changed.
type FlockSettingsChangedEvent struct {
	Description string `json:"description,omitempty"` // Flock Settings Changed
	Logtype     string `json:"logtype,omitempty"`     // 23003
	Settings    string `json:"settings,omitempty"`    // the settings that were changed.
}

// RollbackNetworkSettingsEvents is a Rollback Network Settings triggered Event
// Triggered by a Canary rolling back its settings after an unsuccessful attempt to change its network settings.
type RollbackNetworkSettingsEvents struct {
	Description string `json:"description,omitempty"` // Network Settings Roll-back
	Logtype     string `json:"logtype,omitempty"`     // 22001
}

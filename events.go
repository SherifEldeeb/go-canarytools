package canarytools

// ref: https://canarytools.readthedocs.io/en/latest/incident_attributes.html

// HTTPEvent is an HTTP Canarytoken triggered Event
type HTTPEvent struct {
	Description string                 // Canarytoken triggered
	Type        string                 // http
	Canarytoken string                 // Unique string that acts as the Canarytoken
	Headers     map[string]interface{} // Headers is a dict, Only present for HTTP Canarytokens.
	URL         string                 // URL of the HTTP Canarytoken.
	Logtype     string                 // 17000
}

// WebImageEvent is a Remote Web Image triggered Event
type WebImageEvent struct {
	description  string                 // "Remote Web Image"
	Type         string                 // "web-image"
	canarytoken  string                 // Unique string that acts as the Canarytoken
	headers      map[string]interface{} // Headers is a dict, Only present for HTTP Canarytokens.
	url          string                 // URL of the HTTP Canarytoken.
	logtype      string                 // "17001"
	WebImage     string                 // Byte string of the web image
	WebImageType string                 // Type of the web image
	WebImageName string                 // Name of the web image
}

// DocMSWordEvent is a MS Word triggered Event
type DocMSWordEvent struct {
	description string                 // "Canarytoken triggered"
	Type        string                 // "doc-msword"
	canarytoken string                 // Unique string that acts as the Canarytoken
	headers     map[string]interface{} // Headers is a dict, Only present for HTTP Canarytokens.
	url         string                 // URL of the HTTP Canarytoken.
	logtype     string                 // "17002"
	doc         string                 // Byte String of the tokened document
	DocName     string                 // Name of the document tokened (or created)
	DocType     string                 // Type of document chosen (doc or docx)
}

// ClonedWebEvent is a Cloned Website triggered Event
type ClonedWebEvent struct {
	description string                 // "Cloned Website"
	Type        string                 // "cloned-web"
	canarytoken string                 // Unique string that acts as the Canarytoken
	headers     map[string]interface{} // Headers is a dict, Only present for HTTP Canarytokens.
	url         string                 // URL of the HTTP Canarytoken.
	logtype     string                 // "17003"
	ClonedWeb   string                 // Domain that we are tokening
}

// AWSS3Event is an Amazon AWS S3  triggered Event
type AWSS3Event struct {
	description    string                 // "Amazon S3"
	Type           string                 // "aws-s3"
	canarytoken    string                 // Unique string that acts as the Canarytoken
	headers        map[string]interface{} // Headers is a dict, Only present for HTTP Canarytokens.
	url            string                 // URL of the HTTP Canarytoken.
	logtype        string                 // "17005"
	S3SourceBucket string                 // bucket that we are tokening
	S3LogBucket    string                 // bucket where logging to stored and monitored
	online         string                 // Whether the token is online or not
}

// GoogleDocsEvent is a Google Docs triggered Event
type GoogleDocsEvent struct {
	description  string                 // "Google Document"
	Type         string                 // "google-docs"
	canarytoken  string                 // Unique string that acts as the Canarytoken
	headers      map[string]interface{} // Headers is a dict, Only present for HTTP Canarytokens.
	url          string                 // URL of the HTTP Canarytoken.
	logtype      string                 // "17006"
	DocsLink     string                 // url to the google doc
	EmailLink    string                 // url used for email
	DocumentName string                 // Name of the document
}

// GoogleSheetsEvent is a Google Docs triggered Event
type GoogleSheetsEvent struct {
	description  string                 // "Google Document"
	Type         string                 // "google-sheets
	canarytoken  string                 // Unique string that acts as the Canarytoken
	headers      map[string]interface{} // Headers is a dict, Only present for HTTP Canarytokens.
	url          string                 // URL of the HTTP Canarytoken.
	logtype      string                 // "17006"
	DocsLink     string                 // url to the google doc
	EmailLink    string                 // url used for email
	DocumentName string                 // Name of the document
}

// QRCodeEvent is a QR Code triggered Event
type QRCodeEvent struct {
	description string                 // “QR Code”
	Type        string                 // “qr-code”
	canarytoken string                 // Unique string that acts as the Canarytoken
	headers     map[string]interface{} // Headers is a dict, Only present for HTTP Canarytokens.
	url         string                 // URL of the HTTP Canarytoken.
	logType     string                 // “17009”
	QRCode      string                 // Byte string of the tokened QR code

}

// SVNEvent is a SVN triggered Event
type SVNEvent struct {
	description string                 // “SVN Repo”
	Type        string                 // “svn”
	canarytoken string                 // Unique string that acts as the Canarytoken
	headers     map[string]interface{} // Headers is a dict, Only present for HTTP Canarytokens.
	url         string                 // URL of the HTTP Canarytoken.
	logType     string                 // “17010”

}

// SQLEvent is a SQL triggered Event
type SQLEvent struct {
	description   string                 // “SQL Server”
	Type          string                 // “sql”
	canarytoken   string                 // Unique string that acts as the Canarytoken
	headers       map[string]interface{} // Headers is a dict, Only present for HTTP Canarytokens.
	url           string                 // URL of the HTTP Canarytoken.
	logType       string                 // “17011”
	trigger_Type  string                 // SQL trigger type (SELECT, UPDATE, INSERT, DELETE)
	table_name    string                 // SQL table name (trigger_type: UPDATE, INSERT,DELETE)
	trigger_name  string                 // SQL trigger name (trigger_type: UPDATE, INSERT,DELETE)
	view_name     string                 // SQL View name (trigger_type: SELECT)
	function_name string                 // SQL function name (trigger_type: SELECT)

}

// AWSIDEvent is a AWS ID triggered Event
type AWSIDEvent struct {
	description       string                 // “Amazon API Key”
	Type              string                 // “aws-id”
	canarytoken       string                 // Unique string that acts as the Canarytoken
	headers           map[string]interface{} // Headers is a dict, Only present for HTTP Canarytokens.
	url               string                 // URL of the HTTP Canarytoken.
	logType           string                 // “17012”
	secret_access_key string                 // AWS generated secret access key
	access_key_id     string                 // AWS generated access key ID

}

// FastRedirectEvent is a Fast Redirect triggered Event
type FastRedirectEvent struct {
	description          string                 // “Fast HTTP Redirect”
	Type                 string                 // “fast-redirect”
	canarytoken          string                 // Unique string that acts as the Canarytoken
	headers              map[string]interface{} // Headers is a dict, Only present for HTTP Canarytokens.
	url                  string                 // URL of the HTTP Canarytoken.
	logType              string                 // “17016”
	browser_redirect_url string                 // Original url attempted before redirect

}

// SlowRedirectEvent is a Slow Redirect triggered Event
type SlowRedirectEvent struct {
	description          string                 // “Slow HTTP Redirect”
	Type                 string                 // “slow-redirect”
	canarytoken          string                 // Unique string that acts as the Canarytoken
	headers              map[string]interface{} // Headers is a dict, Only present for HTTP Canarytokens.
	url                  string                 // URL of the HTTP Canarytoken.
	logType              string                 // “17017”
	browser_redirect_url string                 // Original url attempted before redirect

}

// DNSEvent is a DNS triggered Event
type DNSEvent struct {
	description string                 // “DNS”
	Type        string                 // “dns”
	canarytoken string                 // Unique string that acts as the Canarytoken
	hostname    map[string]interface{} // Hostname of the DNS Canarytoken.
	logType     string                 // “16000”

}

// DesktopINIEvent is a Desktop ini triggered Event
type DesktopINIEvent struct {
	description string                 // “Windows Directory Browsing”
	Type        string                 // “windows-dir”
	canarytoken string                 // Unique string that acts as the Canarytoken
	hostname    map[string]interface{} // Hostname of the DNS Canarytoken.
	logType     string                 // “16006”

}

// AdobeReaderPDFEvent is a Adobe Reader PDF triggered Event
type AdobeReaderPDFEvent struct {
	description string                 // “Acrobat Reader PDF Document”
	Type        string                 // “pdf-acrobat-reader”
	canarytoken string                 // Unique string that acts as the Canarytoken
	hostname    map[string]interface{} // Hostname of the DNS Canarytoken.
	logType     string                 // “16008”

}

// MSWordDocMacroedEvent is a MS Word Doc Macroed triggered Event
type MSWordDocMacroedEvent struct {
	description string                 // “MS Word .docm Document”
	Type        string                 // “msword-macro”
	canarytoken string                 // Unique string that acts as the Canarytoken
	hostname    map[string]interface{} // Hostname of the DNS Canarytoken.
	logType     string                 // “16009”
	doc         string                 // Byte String of the tokened document
	doc_name    string                 // Name of the document
	doc_Type    string                 // Type of document chosen (doc or docx)

}

// MSExcelMacroedEvent is a MS Excel Macroed triggered Event
type MSExcelMacroedEvent struct {
	description string                 // “MS Excel .xlsm Document”
	Type        string                 // “msexcel-macro”
	canarytoken string                 // Unique string that acts as the Canarytoken
	hostname    map[string]interface{} // Hostname of the DNS Canarytoken.
	logType     string                 // “16010”
	doc         string                 // Byte String of the tokened document
	doc_name    string                 // Name of the document
	doc_Type    string                 // Type of document chosen (doc or docx)
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
	description string        //“Host Port Scan”
	ports       []interface{} //A list of ports scanned
	logtype     string        //“5003”
}

// ConsolidatedNetworkPortScanEvent is a Consolidated Network Port Scan triggered Event
type ConsolidatedNetworkPortScanEvent struct {
	description   string                 //“Host Port Scan”
	ports_scanned map[string]interface{} //A dictionary of ports scanned and the IP address of the Canaries on which the scan occurred.
	logtype       string                 //“5007”
}

// NMAPNULLScanEvent is a NMAP NULL Scan triggered Event
type NMAPNULLScanEvent struct {
	description string //“NMAP NULL Scan Detected”
	logtype     string //“5005”
}

// NMAPOSScanEvent is a NMAP OS Scan triggered Event
type NMAPOSScanEvent struct {
	description string //“NMAP OS Scan Detected”
	logtype     string //“5004”
}

// NMAPXMASScanEvent is a NMAP XMAS Scan triggered Event
type NMAPXMASScanEvent struct {
	description string //“NMAP XMAS Scan Detected”
	logtype     string //“5006”
}

// NMAPFINScanEvent is a NMAP FIN Scan triggered Event
type NMAPFINScanEvent struct {
	description string //“NMAP FIN Scan Detected”
	logtype     string //“5008”
}

// Canary Connection Events
// CanaryDisconnectedEvent is a Canary Disconnected triggered Even
type CanaryDisconnectedEvent struct {
	description – “Canary Disconnected”
	logtype (str) – “1004”
}

// FTPIncidentEvent is a FTP Incident triggered Even
type FTPIncidentEvent struct {
	description – “FTP Login Attempt”
	logtype (str) – “2000”        username (list) – Attacker supplied username.
	password (list – Attacked supplied password.
}

// GitRepositoryCloneAttemptEvent is a Git Repository Clone Attempt triggered Even
// Triggered when an attacker connects to the Canary git service and attempts any repo clone.
type GitRepositoryCloneAttemptEvent struct {
	description – “Git Repository Clone Attempt”
	logtype (str) – “19001”        host (list) – Git client’s view of the Canary’s hostname.
	repo (str) – Name of the repository the client attempted to clone.
}

// HTTP Incidents
// Two types of HTTP Incidents:
//     Page loads, triggered by GET requests. They are disabled by default as they’re noisy, and needs to be specifically enabled.
//     Login attempts, triggered by GET requests. They are always enabled.

// HTTPPageLoadEvent is a HTTP Page Load triggered Event
type HTTPPageLoadEvent struct {
	description – “HTTP Page Load”
	logtype (str) – “3000”        path (list) – Web path requested by the source.
	useragent (str) – Useragent of the source’s browser.
	channel (str) – Optional. Set to ‘TLS’ if an encrypted site is configured, otherwise absent.
}

// HTTPLoginAttemptEvent is a HTTP Login Attempt triggered Event
type HTTPLoginAttemptEvent struct {
	description – “HTTP Login Attempt”
	logtype (str) – “3001”        username (str) – Attack supplied username.
	password (str) – Attacked supplied password.
	path (list) – Web path requested by the source.
	useragent (str) – Useragent of the source’s browser.
	channel (str) – Optional. Set to ‘TLS’ if an encrypted site is configured, otherwise absent.
}

// HTTPProxyRequestEvent is a HTTP Proxy Request triggered Event
// Triggered by any request through the HTTP proxy module.
type HTTPProxyRequestEvent struct {
	description – “HTTP Proxy Request”
	logtype (str) – “7001”        username (str) – Attack supplied username.
	password (str) – Attacked supplied password.
	url (str) – URL requested by the source.
	useragent (str) – Useragent of the source’s browser.
}

// MicrosoftSQLServerLoginAttemptEvent is a Microsoft SQL Server Login Attempt triggered Event
// Triggered by any attempt to authenticate to the Microsoft SQL Server module.
// SQL Server supports multiple authentication modes, and the fields that come through depend on the mode.
type MicrosoftSQLServerLoginAttemptEvent struct {
	description – “MSSQL Login Attempt”
	logtype (str) – “9001” - SQL Server authentication. “9002” - Windows authentication.
	username (str) – Attack supplied username.
	password (str) – Optional. Attacker supplied database password
	hostname (str) – Optional. Attacker supplied hostname.
	domainname (str) – Optional. Attacker supplied Active Directory name.
}

// ModBusRequestEvent is a ModBus Request triggered Event
// Triggered by any valid ModBus request.
type ModBusRequestEvent struct {
	description – “ModBus Request”
	logtype (str) – “18001” (Modbus Query Function)
	logtype (str) – “18002” (Modbus Read Function)
	logtype (str) – “18003” (Modbus Write Function)
	unit_id (str) – ModBus unit target.
	func_code (str) – ModBus function code.
	func_name (str) – Optional. ModBus function name, if available.
	sfunc_code (str) – Optional. ModBus subfunction code, if available.
	sfunc_nmae (str) – Optional. ModBus subfunction name, if available.
}

// NTPMonlistRequestEvent is a NTP Monlist Request triggered Event
// Triggered by the NTP Monlist command.
type NTPMonlistRequestEvent struct {
	description – “NTP Monlist Request”
	logtype (str) – “11001
	ntp_cmd (str) – Name of the NTP command sent. Currently is ‘monlist’.
	client_hash (str) – Attacker supplied database password hash.
	salt (str) – Attacker supplied database password hash salt.
	password (str) – Recovered password if possible, otherwise ‘<Password not in common list>’
}

// RedisCommandEvent is a Redis Command triggered Event
// Triggered by an attacker connecting to the Redis service and issuing valid Redis commands.
type RedisCommandEvent struct {
	description – “Redis Command”
	logtype (str) – “21001        
	cmd (str) – Redis command issued by the attacker.
	args (str) – Arguments to the command.
}

// RedisCommandEvent is a Redis Command triggered Event
// Triggered by an attacker connecting to the SIP service and issuing valid SIP request.

type RedisCommandEvent struct {
	description – “SIP Request”
	logtype (str) – “15001        headers (dict) – Dict of the SIP headers included in the request.
	args (str) – Arguments to the command.
}


package canarytools

// FlockSummaryResponse is the response received from flock/summary endpoiont
type FlockSummaryResponse struct {
	DifferentTokenNum int           `json:"different_token_num,omitempty"`
	DisabledTokens    int           `json:"disabled_tokens,omitempty"`
	IncidentCount     int           `json:"incident_count,omitempty"`
	Result            string        `json:"result,omitempty"`
	TopTokens         []interface{} `json:"top_tokens,omitempty"`
	TotalTokens       int           `json:"total_tokens,omitempty"`
	TriggeredTokens   int           `json:"triggered_tokens,omitempty"`
	Message           string        `json:"message,omitempty"`
}

// FetchAllTokensResponse is the response received from canarytokens/fetch endpoint
type FetchAllTokensResponse struct {
	Result string  `json:"result"`
	Tokens []Token `json:"tokens"`
}

// Token is a Canarytoken definition
type Token struct {
	BrowserRedirectURL    string `json:"browser_redirect_url,omitempty"`
	BrowserScannerEnabled bool   `json:"browser_scanner_enabled,omitempty"`
	Canarytoken           string `json:"canarytoken"`
	Created               string `json:"created"`
	CreatedPrintable      string `json:"created_printable"`
	Enabled               bool   `json:"enabled"`
	FlockID               string `json:"flock_id"`
	Hostname              string `json:"hostname"`
	Key                   string `json:"key"`
	Kind                  string `json:"kind"`
	Memo                  string `json:"memo"`
	NodeID                string `json:"node_id"`
	TriggeredCount        int    `json:"triggered_count"`
	UpdatedID             int    `json:"updated_id"`
	URL                   string `json:"url"`
	AccessKeyID           string `json:"access_key_id,omitempty"`
	FactoryAuth           string `json:"factory_auth,omitempty"`
	Renders               struct {
		ClonedWeb string `json:"cloned-web"`
		AwsID     string `json:"aws-id"`
	} `json:"renders,omitempty"`
	SecretAccessKey string `json:"secret_access_key,omitempty"`
	Username        string `json:"username,omitempty"`
	ClonedWeb       string `json:"cloned_web,omitempty"`
	Online          string `json:"online,omitempty"`
	QRCode          string `json:"qr_code,omitempty"`
}

// PingResponse holds the ping API response.
// if the result is "success", Message will be empty,
// otherwise it will contain failure reason.
type PingResponse struct {
	Message string `json:"message,omitempty"`
	Result  string `json:"result,omitempty"`
}

// GetDevicesResponse is a response for all the device/* API
type GetDevicesResponse struct {
	Devices          []Device `json:"devices,omitempty"`
	Feed             string   `json:"feed,omitempty"`
	Message          string   `json:"message,omitempty"`
	Result           string   `json:"result,omitempty"`
	Updated          string   `json:"updated,omitempty"`
	UpdatedStd       string   `json:"updated_std,omitempty"`
	UpdatedTimestamp int64    `json:"updated_timestamp,omitempty"`
}

// BasicResponse contains fields that are returned with all responses
type BasicResponse struct {
	Message string `json:"message,omitempty"`
	Result  string `json:"result,omitempty"`
}

// TokenCreateResponse is the response received when creating a token
type TokenCreateResponse struct { // TODO: add all possible fields
	Canarytoken struct {
		AccessKeyID           string `json:"access_key_id,omitempty"`
		SecretAccessKey       string `json:"secret_access_key,omitempty"`
		QRCode                string `json:"qr_code,omitempty"`
		ClonedWeb             string `json:"cloned_web,omitempty"`
		BrowserRedirectURL    string `json:"browser_redirect_url,omitempty"`
		BrowserScannerEnabled bool   `json:"browser_scanner_enabled,omitempty"`
		Canarytoken           string `json:"canarytoken,omitempty"`
		Created               string `json:"created,omitempty"`
		CreatedPrintable      string `json:"created_printable,omitempty"`
		Enabled               bool   `json:"enabled,omitempty"`
		FlockID               string `json:"flock_id,omitempty"`
		Hostname              string `json:"hostname,omitempty"`
		Key                   string `json:"key,omitempty"`
		Kind                  string `json:"kind,omitempty"`
		Memo                  string `json:"memo,omitempty"`
		TriggeredCount        int    `json:"triggered_count,omitempty"`
		UpdatedID             int    `json:"updated_id,omitempty"`
		URL                   string `json:"url,omitempty"`
		Username              string `json:"username,omitempty"`
	} `json:"canarytoken,omitempty"`
	Result string `json:"result,omitempty"`
}

// FlocksSummaryResponse is the response received by flocks/summary API endpoint
type FlocksSummaryResponse struct {
	FlocksSummary               map[string]FlockSummary `json:"flocks_summary,omitempty"`
	Result                      string                  `json:"result,omitempty"`
	UnackedDeviceIncidentCounts map[string]int          `json:"unacked_device_incident_counts,omitempty"`
	Message                     string                  `json:"message,omitempty"`
}

// FlockSummary represents a summary for a single flock
type FlockSummary struct {
	Devices        []string `json:"devices,omitempty"`
	DisabledTokens int      `json:"disabled_tokens,omitempty"`
	EnabledTokens  int      `json:"enabled_tokens,omitempty"`
	FlockID        string   `json:"flock_id,omitempty"`
	GlobalSettings struct {
		ConsoleChangeAlertsEnabled bool     `json:"console_change_alerts_enabled,omitempty"`
		DeviceChangeAlertsEnabled  bool     `json:"device_change_alerts_enabled,omitempty"`
		EmailNotifications         bool     `json:"email_notifications,omitempty"`
		SmsNotifications           bool     `json:"sms_notifications,omitempty"`
		SummaryEmails              bool     `json:"summary_emails,omitempty"`
		WebhookNotifications       bool     `json:"webhook_notifications,omitempty"`
		WhitelistEnabled           bool     `json:"whitelist_enabled,omitempty"`
		WhitelistIps               []string `json:"whitelist_ips,omitempty"`
	} `json:"global_settings,omitempty"`
	IncidentCount  int    `json:"incident_count,omitempty"`
	Name           string `json:"name,omitempty"`
	Note           string `json:"note,omitempty"`
	OfflineDevices int    `json:"offline_devices,omitempty"`
	OnlineDevices  int    `json:"online_devices,omitempty"`
	Settings       struct {
		ChangeControl struct {
			DeviceSettingsNotifications string `json:"device_settings_notifications,omitempty"`
			FlockSettingsNotifications  string `json:"flock_settings_notifications,omitempty"`
		} `json:"change_control,omitempty"`
		NotificationInfo struct {
			Emails struct {
				Addresses string `json:"addresses,omitempty"`
				Enabled   string `json:"enabled,omitempty"`
			} `json:"emails,omitempty"`
			Sms struct {
				Enabled string `json:"enabled,omitempty"`
				Numbers string `json:"numbers,omitempty"`
			} `json:"sms,omitempty"`
			SummaryEmail struct {
				Addresses string `json:"addresses,omitempty"`
				Enabled   string `json:"enabled,omitempty"`
			} `json:"summary_email,omitempty"`
		} `json:"notification_info,omitempty"`
		Webhooks struct {
			GenericWebhooks []string `json:"generic_webhooks,omitempty"`
			HipchatWebhooks []string `json:"hipchat_webhooks,omitempty"`
			MsTeams         []string `json:"ms_teams,omitempty"`
			SlackWebhooks   []string `json:"slack_webhooks,omitempty"`
			WebhooksEnabled string   `json:"webhooks_enabled,omitempty"`
		} `json:"webhooks,omitempty"`
		Whitelisting struct {
			HostnameIgnorelisting struct {
				Enabled   string   `json:"enabled,omitempty"`
				Hostnames []string `json:"hostnames,omitempty"`
			} `json:"hostname_ignorelisting,omitempty"`
			InheritGlobalWhitelistIps bool `json:"inherit_global_whitelist_ips,omitempty"`
			IPWhitelisting            struct {
				Enabled string `json:"enabled,omitempty"`
				Ips     string `json:"ips,omitempty"`
			} `json:"ip_whitelisting,omitempty"`
			SrcPortIgnorelisting struct {
				Enabled  string `json:"enabled,omitempty"`
				SrcPorts string `json:"src_ports,omitempty"`
			} `json:"src_port_ignorelisting,omitempty"`
		} `json:"whitelisting,omitempty"`
	} `json:"settings,omitempty"`
	TotalTokens int `json:"total_tokens,omitempty"`
}

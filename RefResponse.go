package canarytools

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

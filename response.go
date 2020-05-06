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

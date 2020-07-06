package canarytools

// GetIncidentsResponse is the response of the incidents/unacknowledged API
type GetIncidentsResponse struct {
	Feed             string     `json:"feed,omitempty"`
	Incidents        []Incident `json:"incidents,omitempty"`
	MaxUpdatedID     int        `json:"max_updated_id,omitempty"` // cursor
	Message          string     `json:"message,omitempty"`
	Result           string     `json:"result,omitempty"`
	Updated          string     `json:"updated,omitempty"`
	UpdatedStd       string     `json:"updated_std,omitempty"`
	UpdatedTimestamp int        `json:"updated_timestamp,omitempty"`
}

// Incident is an incident, returned from the incidnets API
type Incident struct {
	//	Description []Description `json:"description,omitempty"` // TODO: this varies greatly!
	Description map[string]interface{} `json:"description,omitempty"`
	HashID      string                 `json:"hash_id,omitempty"`
	ID          string                 `json:"id,omitempty"`
	Summary     string                 `json:"summary,omitempty"`
	Updated     string                 `json:"updated,omitempty"`
	UpdatedID   int                    `json:"updated_id,omitempty"`
	UpdatedStd  string                 `json:"updated_std,omitempty"`
	UpdatedTime string                 `json:"updated_time,omitempty"`
	// we add this to inform the acker if we want to ACK the incident
	ThenWhat string `json:"then_what,omitempty"`
}

// Description contains details about incidents
type Description struct {
	Acknowledged   string                   `json:"acknowledged,omitempty"`
	Created        string                   `json:"created,omitempty"`
	CreatedStd     string                   `json:"created_std,omitempty"`
	Description    string                   `json:"description,omitempty"`
	DstHost        string                   `json:"dst_host,omitempty"`
	DstPort        string                   `json:"dst_port,omitempty"`
	Events         []map[string]interface{} `json:"events,omitempty"`
	EventsCount    string                   `json:"events_count,omitempty"`
	IPAddress      string                   `json:"ip_address,omitempty"`
	Ippers         string                   `json:"ippers,omitempty"`
	LocalTime      string                   `json:"local_time,omitempty"`
	Logtype        string                   `json:"logtype,omitempty"`
	MACAddress     string                   `json:"mac_address,omitempty"`
	Memo           string                   `json:"memo,omitempty"`
	Name           string                   `json:"name,omitempty"`
	NodeID         string                   `json:"node_id,omitempty"`
	Notified       string                   `json:"notified,omitempty"`
	SrcHost        string                   `json:"src_host,omitempty"`
	SrcHostReverse string                   `json:"src_host_reverse,omitempty"`
	SrcPort        string                   `json:"src_port,omitempty"`
}

// Event is part of Incidents (typically in an array),
// they contain details of the events which triggered the incident.
type Event struct {
	AuditAction  string `json:"audit_action,omitempty"`
	Domain       string `json:"domain,omitempty"`
	FileName     string `json:"file_name,omitempty"`
	LocalName    string `json:"local_name,omitempty"`
	Mode         string `json:"mode,omitempty"`
	Offset       string `json:"offset,omitempty"`
	RemoteName   string `json:"remote_name,omitempty"`
	ShareName    string `json:"share_name,omitempty"`
	Size         string `json:"size,omitempty"`
	SMBArch      string `json:"smb_arch,omitempty"`
	SMBVer       string `json:"smb_ver,omitempty"`
	Status       string `json:"status,omitempty"`
	User         string `json:"user,omitempty"`
	Timestamp    int    `json:"timestamp,omitempty"`
	TimestampStd string `json:"timestamp_std,omitempty"`
}

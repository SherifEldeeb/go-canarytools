package canarytools

// Device represents a canary device
type Device struct {
	Description                    string `json:"description,omitempty"`
	FirstSeen                      string `json:"first_seen,omitempty"`
	FirstSeenStd                   string `json:"first_seen_std,omitempty"`
	GCPProject                     string `json:"gcp_project,omitempty"`
	GCPZone                        string `json:"gcp_zone,omitempty"`
	ID                             string `json:"id,omitempty"`
	IgnoreNotifications            bool   `json:"ignore_notifications,omitempty"`
	IgnoreNotificationsDisconnects bool   `json:"ignore_notifications_disconnects,omitempty"`
	IgnoreNotificationsReconnects  bool   `json:"ignore_notifications_reconnects,omitempty"`
	InstanceID                     string `json:"instance_id,omitempty"`
	IPAddress                      string `json:"ip_address,omitempty"`
	IPPers                         string `json:"ippers,omitempty"`
	LastSeen                       string `json:"last_seen,omitempty"`
	LastSeenStd                    string `json:"last_seen_std,omitempty"`
	Live                           bool   `json:"live,omitempty"`
	LocalTime                      string `json:"local_time,omitempty"`
	Location                       string `json:"location,omitempty"`
	MACAddress                     string `json:"mac_address,omitempty"`
	MigrationStatus                string `json:"migration_status,omitempty"`
	Name                           string `json:"name,omitempty"`
	Netmask                        string `json:"netmask,omitempty"`
	Note                           string `json:"note,omitempty"`
	PublicIP                       string `json:"public_ip,omitempty"`
	RegionID                       string `json:"region_id,omitempty"`
	Sensor                         string `json:"sensor,omitempty"`
	Subnet                         string `json:"subnet,omitempty"`
	Updated                        string `json:"updated,omitempty"`
	UpdatedStd                     string `json:"updated_std,omitempty"`
	UpdatedTimestamp               int64  `json:"updated_timestamp,omitempty"`
	Uptime                         int64  `json:"uptime,omitempty"`
	UptimeAge                      string `json:"uptime_age,omitempty"`
	Version                        string `json:"version,omitempty"`
}

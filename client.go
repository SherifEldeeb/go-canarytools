package canarytools

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Client is a canarytools client, which is used to issue requests to the API
type Client struct {
	domain     string
	apikey     string
	url        string
	httpclient *http.Client
	l          *log.Logger
}

// NewClient creates a new client from domain & API Key
func NewClient(domain, apikey, loglevel string) (c Client, err error) {
	// logging config
	c.l = log.New()
	switch loglevel {
	case "info":
		c.l.SetLevel(log.InfoLevel)
	case "warning":
		c.l.SetLevel(log.WarnLevel)
	case "debug":
		c.l.SetLevel(log.DebugLevel)
	default:
		return c, errors.New("unsupported log level (can be 'info', 'warning' or 'debug')")
	}

	c.httpclient = &http.Client{Timeout: 5 * time.Second} // TODO: provide ability to configure
	c.domain = domain
	c.apikey = apikey
	c.url = fmt.Sprintf("https://%s.canary.tools/api/v1/%%s?auth_token=%s", domain, apikey)

	c.l.Debug("pinging console...")
	err = c.Ping()
	return
}

// api constructs the full URL for API querying,
// wich includes the domain and the API auth token
func (c Client) api(endpoint string) (url string) {
	return fmt.Sprintf(c.url, endpoint)
}

// decodeResponse decodes reponses into target interfaces
func (c Client) decodeResponse(api string, target interface{}) (err error) {
	targetAPI := c.api(api)

	c.l.WithFields(log.Fields{
		"url": targetAPI,
	}).Debug("hitting API")
	resp, err := c.httpclient.Get(targetAPI)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

// Ping tests connection to the console, and validity of connection params
func (c Client) Ping() (err error) {
	var pr PingResponse
	err = c.decodeResponse("ping", &pr)
	if err != nil {
		return
	}

	if pr.Result != "success" {
		return errors.New(pr.Message) // there will be a message, if it failed
	}

	return
}

// getDevices returns  devices
func (c Client) getDevices(which string) (devices []Device, err error) {
	var getdevicesresponse GetDevicesResponse
	err = c.decodeResponse("devices/"+which, &getdevicesresponse)
	if err != nil {
		return
	}

	if getdevicesresponse.Result != "success" {
		return nil, errors.New(getdevicesresponse.Message) // there will be a message, if it failed
	}
	return getdevicesresponse.Devices, nil
}

// GetAllDevices returns all devices
func (c Client) GetAllDevices() (devices []Device, err error) {
	return c.getDevices("all")
}

// GetLiveDevices returns live devices
func (c Client) GetLiveDevices() (devices []Device, err error) {
	return c.getDevices("live")
}

// GetDeadDevices returns live devices
func (c Client) GetDeadDevices() (devices []Device, err error) {
	return c.getDevices("dead")
}

// GetUnacknowledgedIncidents returns all Unacknowledged Incidents since lastID,
// setting lastID to 0 returns all unack'd incidents,
//
func (c Client) GetUnacknowledgedIncidents(lastID int) (incidents []Incident, maxUpdatedID int, err error) {
	var unackIncidents GetIncidentsResponse
	err = c.decodeResponse("incidents/unacknowledged", &unackIncidents)
	if err != nil {
		return
	}

	if unackIncidents.Result != "success" {
		return nil, 0, errors.New(unackIncidents.Message) // there will be a message, if it failed
	}

	return unackIncidents.Incidents, unackIncidents.MaxUpdatedID, nil
}

package canarytools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	log "github.com/sirupsen/logrus"
)

// Client is a canarytools client, which is used to issue requests to the API
type Client struct {
	domain            string
	apikey            string
	baseURL           *url.URL
	httpclient        *http.Client
	l                 *log.Logger
	lastCheck         time.Time
	lastCheckRegister *os.File
	errorCount        int
	fetchInterval     int
	thenWhat          string
	whichIncidents    string
	lastCheckFile     string
}

// NewClient creates a new client from domain & API Key
func NewClient(domain, apikey, thenWhat, sinceWhen, whichIncidents string, fetchInterval int, l *log.Logger) (c *Client, err error) {
	c = &Client{}
	c.l = l
	// create local work folder
	err = os.MkdirAll(".canary", 0755)
	if err != nil {
		return
	}
	c.lastCheckFile = ".canary/lastcheck"
	// time parsing
	var t time.Time
	switch {
	// if valid time has been provided, it superseeds everything else
	// in this case, we either get a valid date, or we fail miserably.
	case sinceWhen != "":
		t, err = time.Parse("2006-01-02 15:04:05", sinceWhen)
		if err != nil {
			c.l.WithFields(log.Fields{
				"err":          err,
				"providedTime": sinceWhen,
			}).Warn("error parsing time from provided value, setting default time (-7days)!")
			t = time.Now().AddDate(0, 0, -7).UTC()
		}
	// if nothing provided, we look for '.canary.lastcheck' file
	case sinceWhen == "":
		if _, err = os.Stat(c.lastCheckFile); err == nil { // file exists, and we have no issues reading it
			var b = []byte{}
			b, err = ioutil.ReadFile(c.lastCheckFile)
			if err != nil {
				return
			}
			// now we shoould have the content in that file
			s := string(b)
			t, err = time.Parse("2006-01-02 15:04:05", s)
			if err != nil {
				c.l.WithFields(log.Fields{
					"err":           err,
					c.lastCheckFile: s,
				}).Warn("error parsing time from lastCheckFile, setting default time (-7days)!")
				t = time.Now().AddDate(0, 0, -7).UTC()
			}
		} else { // file doesn't exist, we default to (today - 7 days).
			t = time.Now().AddDate(0, 0, -7).UTC()
		}
	}
	l.WithField("sinceWhen", t).Info("Events 'sinceWhen' parsed or successfully set")
	c.lastCheck = t
	c.thenWhat = thenWhat
	c.whichIncidents = whichIncidents
	c.fetchInterval = fetchInterval
	c.httpclient = &http.Client{Timeout: 5 * time.Second} // TODO: provide ability to configure
	c.domain = domain
	c.apikey = apikey
	c.baseURL, err = url.Parse(fmt.Sprintf("https://%s.canary.tools/api/v1/", domain))
	if err != nil {
		return
	}
	// c.url = fmt.Sprintf("https://%s.canary.tools/api/v1/%%s?auth_token=%s", domain, apikey)

	c.l.Debug("pinging console...")
	err = c.Ping()
	if err != nil {
		return
	}
	// write lastcheck register
	c.lastCheckRegister, err = os.Create(c.lastCheckFile)
	if err != nil {
		return
	}
	return
}

func (c *Client) WriteLastCheckRegister(t time.Time) (err error) {
	err = c.lastCheckRegister.Truncate(0)
	if err != nil {
		return
	}
	_, err = c.lastCheckRegister.Seek(0, 0)
	if err != nil {
		return
	}
	_, err = c.lastCheckRegister.WriteString(t.Format("2006-01-02 15:04:05"))
	return
}

// api constructs the full URL for API querying, it always adds the API auth
// token, and adds  optional parameters as needed.
func (c Client) api(endpoint string, params *url.Values) (fullURL *url.URL, err error) {
	if endpoint == "" {
		return nil, errors.New("API endpoint has not been provided")
	}

	// if no additional params has been provided, we have to construct one
	if params == nil {
		params = &url.Values{}
	}
	// always add auth token to list of values
	params.Add("auth_token", c.apikey)

	// adding the API endpoint to path
	fullURL, err = url.Parse(c.baseURL.String())
	if err != nil {
		return
	}
	fullURL.Path = path.Join(fullURL.Path, endpoint)

	// building the full query
	fullURL.RawQuery = params.Encode()
	return
}

// decodeResponse decodes reponses into target interfaces
func (c Client) decodeResponse(endpoint, verb string, params *url.Values, target interface{}) (err error) {
	fullURL, err := c.api(endpoint, params)
	if err != nil {
		return
	}

	c.l.WithFields(log.Fields{
		"url":      fullURL.String(), // TODO: remove sensitive data
		"HTTPverb": verb,
	}).Debug("hitting API")

	var resp = &http.Response{}
	switch verb {
	case "GET":
		resp, err = c.httpclient.Get(fullURL.String())
	case "POST":
		resp, err = c.httpclient.Post(fullURL.String(), "application/x-www-form-urlencoded", nil)
	case "DELETE":
		req, _ := http.NewRequest("DELETE", fullURL.String(), nil)
		resp, err = c.httpclient.Do(req)
	}
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

// Ping tests connection to the console, and validity of connection params
func (c Client) Ping() (err error) {
	var pr PingResponse
	err = c.decodeResponse("ping", "GET", nil, &pr)
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
	err = c.decodeResponse("devices/"+which, "GET", nil, &getdevicesresponse)
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

// getIncidents returns all Incidents since time specified, which is either
// "all", or "unacknowledged".
// setting "since" to zero vaule (time.Time{}) returns all incidents.
func (c Client) getIncidents(which string, since time.Time) (incidents []Incident, err error) {
	var inc GetIncidentsResponse
	var ts string
	var tt time.Time

	// this API has an optional parameter (newer_than) which is Timestamp used
	// to filter returned incidents. All incidents created after this timestamp
	// will be returned.
	// Format: ‘yyyy-mm-dd-hh:mm:ss’
	if since.Equal(time.Time{}) {
		// dummy date, but definetly before any incident
		tt, _ = time.Parse(time.RFC3339, "2000-01-02T15:04:05Z")
	} else {
		tt = since
	}
	ts = tt.Format("2006-01-02-15:04:05")
	u := &url.Values{}
	u.Add("newer_than", ts)
	u.Add("shrink", "true")
	err = c.decodeResponse("incidents/"+which, "GET", u, &inc)
	if err != nil {
		return
	}

	if inc.Result != "success" {
		return nil, errors.New(inc.Message) // there will be a message, if it failed
	}

	return inc.Incidents, nil
}

// GetUnacknowledgedIncidents returns all Unacknowledged Incidents since time
// secified, setting "since" to zero vaule (time.Time{}) returns all incidents
func (c Client) GetUnacknowledgedIncidents(since time.Time) (incidents []Incident, err error) {
	return c.getIncidents("unacknowledged", since)
}

// GetAllIncidents returns all Incidents since time
// secified, setting "since" to zero vaule (time.Time{}) returns all incidents,
func (c Client) GetAllIncidents(since time.Time) (incidents []Incident, err error) {
	return c.getIncidents("all", since)
}

// Feed fetches incidents and feeds them to chan
func (c *Client) Feed(incidnetsChan chan<- Incident) {
	for {
		// get all unacked incidents
		c.l.WithFields(log.Fields{
			"lastCheck":      c.lastCheck,
			"whichIncidents": c.whichIncidents,
		}).Debug("getting incidents")
		var incidents = []Incident{}
		var err error
		switch c.whichIncidents {
		case "all":
			incidents, err = c.GetAllIncidents(c.lastCheck)
		case "unacknowledged":
			incidents, err = c.GetUnacknowledgedIncidents(c.lastCheck)
		default:
			c.l.WithFields(log.Fields{
				"lastCheck":      c.lastCheck,
				"whichIncidents": c.whichIncidents,
			}).Fatal("unknown whichIncident")
		}
		if err != nil {
			c.l.Error(err) // TODO: fail gracefully
			time.Sleep(time.Duration(c.fetchInterval) * time.Second)
			continue
		}
		c.lastCheck = time.Now().UTC()
		log.Debugf("found total of %d unacked incidents", len(incidents))
		for _, v := range incidents {
			log.WithFields(log.Fields{
				"UpdatedID": v.UpdatedID,
			}).Debug(v.Summary)
			incidnetsChan <- v
			// if c.thenWhat == "ack" {
			// 	a, ok := v.Description["acknowledged"]
			// 	if ok {
			// 		if a == "False" {
			// 			err = c.AckIncident(v.ID)
			// 		}
			// 	}
			// }
		}
		// update register
		err = c.WriteLastCheckRegister(c.lastCheck)
		if err != nil {
			c.l.WithFields(log.Fields{
				"lastCheck": c.lastCheck,
				"err":       err,
			}).Fatal("error writing lastcheck register file")
		}

		// sleep
		time.Sleep(time.Duration(c.fetchInterval) * time.Second)
	}
}

func (c *Client) AckIncidents(ackedIncident <-chan []byte) {
	for v := range ackedIncident {
		var incident Incident
		err := json.Unmarshal(v, &incident)
		if err != nil {
			c.l.WithFields(log.Fields{
				"source": "ConsoleClient",
				"stage":  "ack",
				"err":    err,
			}).Error("Client error Ack Incident")
			continue
		}
		// do it
		a, ok := incident.Description["acknowledged"]
		if ok {
			// ack, only if it hasn't been ack'd already, and we're supposed to ack
			if a == "False" && c.thenWhat == "ack" {
				err = c.AckIncident(incident.ID)
			}
		}
	}
}

// AckIncident acknowledges incident
func (c *Client) AckIncident(incident string) (err error) {
	c.l.WithFields(log.Fields{
		"source":   "ConsoleClient",
		"stage":    "ack",
		"incident": incident,
	}).Debug("Client Ack Incident")

	uv := &url.Values{}
	uv.Add("incident", incident)

	br := &BasicResponse{}
	err = c.decodeResponse("incident/acknowledge", "POST", uv, br)
	if err != nil {
		c.l.WithFields(log.Fields{
			"source": "ConsoleClient",
			"stage":  "ack",
			"err":    err,
		}).Error("Client error Ack Incident")
		return
	}
	c.l.WithFields(log.Fields{
		"source":   "ConsoleClient",
		"stage":    "ack",
		"incident": incident,
	}).Info("Client successfully Ack'd Incident")

	return
}

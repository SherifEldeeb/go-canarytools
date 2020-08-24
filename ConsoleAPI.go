package canarytools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Client is a canarytools client, which is used to issue requests to the API
type Client struct {
	domain     string
	apikey     string
	baseURL    *url.URL
	httpclient *http.Client
	l          *log.Logger
}

// GetFlockNameFromID retrieves the flock name given its ID
func (c Client) GetFlockNameFromID(flockid string) (flockname string, err error) {
	fsr, err := c.GetFlocksSummary()
	if err != nil {
		return
	}

	for fid, flockSummary := range fsr.FlocksSummary {
		if flockid == fid {
			flockname = flockSummary.Name
			return
		}
	}
	err = fmt.Errorf("flock_id does not exist: %s", flockid)
	return
}

// GetFlockIDFromName retrieves the flock ID given its name
func (c Client) GetFlockIDFromName(flockname string) (flockid string, err error) {
	fsr, err := c.GetFlocksSummary()
	if err != nil {
		return
	}

	for fid, flockSummary := range fsr.FlocksSummary {
		if flockSummary.Name == flockname {
			flockid = fid
			return
		}
	}
	err = fmt.Errorf("flock name does not exist: %s", flockname)
	return
}

// FlockCreate creates a flock
func (c Client) FlockCreate(flockname string) (flockid string, err error) {
	fcr := FlockCreateResponse{}
	u := &url.Values{}
	u.Set("name", flockname)
	err = c.decodeResponse("flock/create", "POST", u, &fcr)
	if err != nil {
		return
	}
	if fcr.Result != "success" {
		err = fmt.Errorf("error creating flock")
	}
	flockid = fcr.FlockID
	return
}

// FlockNameExists checks if a flock exists given its name
func (c Client) FlockNameExists(flockname string) (exists bool, FlockID string, err error) {
	fsr, err := c.GetFlocksSummary()
	if err != nil {
		return
	}
	for fid, flockSummary := range fsr.FlocksSummary {
		if flockSummary.Name == flockname {
			exists = true
			FlockID = fid
			return
		}
	}
	return
}

// FlockIDExists checks if a flock exists given its ID
func (c Client) FlockIDExists(flockid string) (exists bool, err error) {
	fsr, err := c.GetFlocksSummary()
	if err != nil {
		return
	}
	for fid := range fsr.FlocksSummary {
		if fid == flockid {
			exists = true
			return
		}
	}
	return
}

// GetFlocksSummary returns summary for all flocks
func (c Client) GetFlocksSummary() (flocksSummaryResponse FlocksSummaryResponse, err error) {
	flocksSummaryResponse = FlocksSummaryResponse{}
	err = c.decodeResponse("flocks/summary", "GET", nil, &flocksSummaryResponse)
	if err != nil {
		return
	}
	if flocksSummaryResponse.Result != "success" {
		err = fmt.Errorf(flocksSummaryResponse.Message)
	}
	return
}

// GetFlockSummary returns summary for a single flock
func (c Client) GetFlockSummary(flockid string) (flocksummaryresponse FlockSummaryResponse, err error) {
	flocksummaryresponse = FlockSummaryResponse{}
	u := &url.Values{}
	u.Set("flock_id", flockid)
	err = c.decodeResponse("flock/fetch", "GET", u, &flocksummaryresponse)
	if err != nil {
		return
	}
	if flocksummaryresponse.Result != "success" {
		err = fmt.Errorf(flocksummaryresponse.Message)
	}
	return
}

// NewClient creates a new client from domain & API Key
func NewClient(domain, apikey string, l *log.Logger) (c *Client, err error) {
	c = &Client{}
	c.l = l
	c.httpclient = &http.Client{Timeout: 10 * time.Second}
	c.domain = domain
	c.apikey = apikey
	c.baseURL, err = url.Parse(fmt.Sprintf("https://%s.canary.tools/api/v1/", domain))
	if err != nil {
		return
	}

	c.l.Debug("pinging console...")
	err = c.Ping()
	return
}

// FetchCanarytokenAll fetches all canarytokens
func (c Client) FetchCanarytokenAll() (tokens []Token, err error) {
	fetchalltokenresponse := FetchAllTokensResponse{}
	err = c.decodeResponse("canarytokens/fetch", "GET", nil, &fetchalltokenresponse)
	tokens = fetchalltokenresponse.Tokens
	return
}

// DeleteCanarytoken fetches all canarytokens
func (c Client) DeleteCanarytoken(canarytoken string) (err error) {
	br := BasicResponse{}
	u := &url.Values{}
	u.Set("canarytoken", canarytoken)
	err = c.decodeResponse("canarytoken/delete", "POST", u, &br)
	if br.Result != "success" {
		err = fmt.Errorf("Error deleting token %s: %s", canarytoken, br.Message)
	}
	return
}

// DropFileToken drops a file token
func (c Client) DropFileToken(kind, memo, filename, FlockID string, CreateFlockIfNotExists bool) (err error) {
	c.l.WithFields(log.Fields{
		"kind":                   kind,
		"memo":                   memo,
		"flock_id":               FlockID,
		"CreateFlockIfNotExists": CreateFlockIfNotExists,
		"filename":               filename,
	}).Debugf("Generating Token")

	var tcr = TokenCreateResponse{}
	switch kind {
	case "aws-id":
		tcr, err = c.CreateTokenFromAPI(kind, memo, FlockID, nil)
		if err != nil {
			return
		}
		if tcr.Result != "success" {
			err = fmt.Errorf("failed to CreateTokenFromAPI")
			return
		}
		var aswTemplate = `[default]
aws_access_key=%s
aws_secret_access_key=%s
region=us-east-2
output=json
`
		// simple checks
		if !fileExists(filename) {
			// Create the file
			out, err := os.Create(filename)
			if err != nil {
				return err
			}
			defer out.Close()

			// Write the body to file
			_, err = out.WriteString(fmt.Sprintf(aswTemplate, tcr.Canarytoken.AccessKeyID, tcr.Canarytoken.SecretAccessKey))
		} else {
			return fmt.Errorf("file exists: %s", filename)
		}
	case "doc-msword", "pdf-acrobat-reader", "msword-macro", "msexcel-macro":
		tcr, err = c.CreateTokenFromAPI(kind, memo, FlockID, nil)
		if err != nil {
			return
		}
		if tcr.Result != "success" {
			err = fmt.Errorf("failed to CreateTokenFromAPI")
			return
		}
		_, err = c.DownloadTokenFromAPI(tcr.Canarytoken.Canarytoken, filename)
	default:
		err = fmt.Errorf("unsupported Canarytoken: %s", kind)
		return
	}
	return
}

// CreateTokenFromAPI uses the canarytoken/create API endpoint to create a token
func (c Client) CreateTokenFromAPI(kind, memo, FlockID string, additionalParams *url.Values) (tokencreateresponse TokenCreateResponse, err error) {
	tokencreateresponse = TokenCreateResponse{}
	u := &url.Values{}
	if additionalParams != nil {
		u = additionalParams
	}
	u.Set("kind", kind)
	u.Set("memo", memo)
	if FlockID != "" {
		u.Set("flock_id", FlockID)
	}

	switch kind {
	// case "doc-msword", "pdf-acrobat-reader", "msword-macro", "msexcel-macro":
	case "http", "dns", "cloned-web", "doc-msword", "web-image", "windows-dir", "aws-s3", "pdf-acrobat-reader", "msword-macro", "msexcel-macro", "aws-id", "apeeper", "qr-code", "svn", "sql", "fast-redirect", "slow-redirect":
	// TODO: must check additional params per kind
	default:
		return tokencreateresponse, errors.New("unsupported token type: " + kind)
	}

	err = c.decodeResponse("canarytoken/create", "POST", u, &tokencreateresponse)
	return
}

// DownloadTokenFromAPI downloads a file-based token given its ID
func (c Client) DownloadTokenFromAPI(canarytoken, filename string) (n int64, err error) {
	params := &url.Values{}
	params.Set("canarytoken", canarytoken)

	fullURL, err := c.api("canarytoken/download", params)
	if err != nil {
		return
	}
	resp, err := c.httpclient.Get(fullURL.String())
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("DownloadTokenFromAPI returned: %d", resp.StatusCode)
	}

	if !fileExists(filename) {
		// Create the file
		out, err := os.Create(filename)
		if err != nil {
			return 0, err
		}
		defer out.Close()

		// Write the body to file
		n, err = io.Copy(out, resp.Body)
	} else {
		return 0, fmt.Errorf("file exists: %s", filename)
	}
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
	var resp = &http.Response{}
	var fullURL = &url.URL{}
	fullURL, err = c.api(endpoint, nil)
	if err != nil {
		return
	}

	switch verb {
	case http.MethodGet:
		fullURL, err = c.api(endpoint, params)
		if err != nil {
			return
		}
		c.l.WithFields(log.Fields{
			"url":      fullURL.String(), // TODO: remove sensitive data
			"HTTPverb": verb,
		}).Debug("hitting API")
		resp, err = c.httpclient.Get(fullURL.String())
	case http.MethodPost:
		resp, err = c.httpclient.PostForm(fullURL.String(), *params)
	case http.MethodDelete:
		var req = &http.Request{}
		if params != nil {
			req, _ = http.NewRequest(http.MethodDelete, fullURL.String(), strings.NewReader(params.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		} else {
			req, _ = http.NewRequest(http.MethodDelete, fullURL.String(), nil)
		}
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

// ThenWhatIncidents performs "ThenWhat" action on incidents
func (c *Client) ThenWhatIncidents(thenWhat string, incidents <-chan []byte) {
	switch thenWhat {
	case "ack":
		c.AckIncidents(incidents)
	case "delete":
		c.DeleteIncidents(incidents)
	case "nothing":
		for range incidents {
		}
	default:
		c.l.WithField("ThenWhat", thenWhat).Fatal("unsupported ThenWhat")
	}
}

// AckIncidents consumes incidents from an incidents chan,
// and ACKs them if they haven't been ACK'd already
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
			// ack, only if it hasn't been ack'd already
			if a == "False" && incident.ThenWhat == "ack" {
				err = c.AckIncident(incident.ID)
				if err != nil {
					c.l.WithField("err", err).Error("error acking incident")
				}
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

// DeleteIncidents consumes incidents from an incidents chan,
// and deletes them
func (c *Client) DeleteIncidents(incidents <-chan []byte) {
	for v := range incidents {
		var incident Incident
		err := json.Unmarshal(v, &incident)
		if err != nil {
			c.l.WithFields(log.Fields{
				"source": "ConsoleClient",
				"stage":  "delete",
				"err":    err,
			}).Error("Client error delete Incident")
			continue
		}
		err = c.DeleteIncident(incident.ID)
		if err != nil {
			c.l.WithField("err", err).Error("error deleting incident")
		}
	}
}

// DeleteIncident deletes incident
func (c *Client) DeleteIncident(incident string) (err error) {
	c.l.WithFields(log.Fields{
		"source":   "ConsoleClient",
		"stage":    "delete",
		"incident": incident,
	}).Debug("Client Delete Incident")

	uv := &url.Values{}
	uv.Add("incident", incident)

	br := &BasicResponse{}
	err = c.decodeResponse("incident/delete", http.MethodDelete, uv, br)
	if err != nil {
		c.l.WithFields(log.Fields{
			"source": "ConsoleClient",
			"stage":  "delete",
			"err":    err,
		}).Error("Client error delete Incident")
		return
	}
	if br.Result != "success" {
		return fmt.Errorf("error deleting incident:%s: %s", incident, br.Message)
	}
	c.l.WithFields(log.Fields{
		"source":   "ConsoleClient",
		"stage":    "delete",
		"incident": incident,
	}).Info("Client successfully delete Incident")

	return
}

package canarytools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"

	log "github.com/sirupsen/logrus"
)

// Client is a canarytools client, which is used to issue requests to the API
type Client struct {
	domain      string
	apikey      string
	factoryAuth string
	opmode      string // can be "api" or "factory"
	baseURL     *url.URL
	httpclient  *http.Client
	l           *log.Logger
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
	errJsonDecode := c.decodeResponse("flocks/summary", "GET", nil, &flocksSummaryResponse)
	if errJsonDecode != nil {
		c.l.Debugf("error decoding JSON response (shouldn't be a big problem, unless we fail the next check): %s", errJsonDecode)
	}
	if flocksSummaryResponse.Result != "success" {
		err = fmt.Errorf("error getting flocks summary: %s", flocksSummaryResponse.Message)
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

// NewClient creates a new client from domain, auth token and operation mode
// operation mode determines the auth token type:
// "api": auth token is the main console api
// "factory": auth token  is the factory_auth
func NewClient(cfg ConsoleAPIConfig, l *log.Logger) (c *Client, err error) {
	c = &Client{}
	c.l = l
	c.httpclient = &http.Client{Timeout: 180 * time.Second}
	c.domain = cfg.ConsoleAPIDomain
	c.opmode = cfg.OpMode
	switch c.opmode {
	case "api", "":
		c.apikey = cfg.ConsoleAPIKey
		c.opmode = "api"
	case "factory":
		c.factoryAuth = cfg.ConsoleFactoryAuth
	default:
		return nil, fmt.Errorf("unsupported opmode: %s, valid values are 'api' & 'factory'")
	}
	c.baseURL, err = url.Parse(fmt.Sprintf("https://%s.canary.tools/api/v1/", cfg.ConsoleAPIDomain))
	if err != nil {
		return
	}

	if c.opmode == "api" { // factory does not support ping
		c.l.Debug("pinging console...")
		err = c.Ping()
	}
	return
}

// FetchCanarytokenAll fetches all canarytokens
func (c Client) FetchCanarytokenAll() (tokens []Token, err error) {
	fetchalltokenresponse := FetchAllTokensResponse{}
	err = c.decodeResponse("canarytokens/fetch", "GET", nil, &fetchalltokenresponse)
	tokens = fetchalltokenresponse.Tokens
	return
}

// DeleteCanarytoken deletes a canarytoken identified by its ID
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

// DeleteAllIncidents deletes multiple incidents identified by a filter
// https://docs.canary.tools/incidents/actions.html#delete-multiple-incidents
func (c Client) DeleteAllIncidents(includeUnacknowledged bool) (err error) {
	br := BasicResponse{}
	u := &url.Values{}
	u.Set("include_unacknowledged", strconv.FormatBool(includeUnacknowledged))
	err = c.decodeResponse("incidents/delete", "DELETE", u, &br)
	if br.Result != "success" {
		err = fmt.Errorf("Error deleting Incidents: %s", br.Message)
	}
	return
}

// DeleteMultipleIncidents deletes multiple incidents identified by a filter
// https://docs.canary.tools/incidents/actions.html#delete-multiple-incidents
func (c Client) DeleteMultipleIncidents(paramType, paramValue string, includeUnacknowledged bool) (err error) {
	switch paramType {
	case "flock_id", "node_id", "src_host", "older_than", "filter_str", "filter_logtypes":
	default:
		return errors.New("unsupported parameter for DeleteMultipleIncidents: " + paramType)
	}
	br := BasicResponse{}
	u := &url.Values{}
	u.Set(paramType, paramValue)
	u.Set("include_unacknowledged", strconv.FormatBool(includeUnacknowledged))
	err = c.decodeResponse("incidents/delete", "DELETE", u, &br)
	if br.Result != "success" {
		err = fmt.Errorf("Error deleting Incidents %s:%s - %s", paramType, paramValue, br.Message)
	}
	return
}

// DropFileToken drops a file token
func (c Client) DropFileToken(kind, memo, dropWhere, filename, FlockID string, CreateFlockIfNotExists, CreateDirectoryIfNotExists, OverwriteFileIfExists bool) (err error) {
	c.l.WithFields(log.Fields{
		"kind":                   kind,
		"memo":                   memo,
		"flock_id":               FlockID,
		"CreateFlockIfNotExists": CreateFlockIfNotExists,
		"filename":               filename,
		"dropWhere":              dropWhere,
	}).Debugf("Generating Token")

	// check if 'where' directory exists
	// if it doesn't exist, and CreateDirectoryIfNotExists is true, create it
	// if it doesn't exist, and CreateDirectoryIfNotExists is false, error out
	absPath, err := filepath.Abs(dropWhere)
	if err != nil {
		return
	}
	if _, errstat := os.Stat(absPath); os.IsNotExist(errstat) { // it does NOT exist
		if CreateDirectoryIfNotExists {
			os.MkdirAll(absPath, 0755)
		} else {
			err = fmt.Errorf("'where' does not exist, and you told me not to create it ... gonna have to bail out")
			return
		}
	}

	fullTokenPath := filepath.Join(dropWhere, filename)

	var tcr = TokenCreateResponse{}
	switch kind {
	case "windows-dir":
		// tcr, err = c.CreateTokenFromAPI(kind, memo, FlockID, nil)
		// if err != nil {
		// 	return
		// }
		// if tcr.Result != "success" {
		// 	err = fmt.Errorf("failed to CreateTokenFromAPI")
		// 	return
		// }
		// _, err = c.DownloadWindowsDirTokenFromAPI(tcr.Canarytoken.Canarytoken, dropWhere, filename, OverwriteFileIfExists)
		// if err != nil {
		// 	return
		// }

		tcr, err = c.CreateTokenFromAPI(kind, memo, FlockID, nil)
		if err != nil {
			return
		}
		if tcr.Result != "success" {
			err = fmt.Errorf("failed to CreateTokenFromAPI")
			return
		}
		var iniTemplate = "\r\n[.ShellClassInfo]\r\nIconResource=\\\\%%USERNAME%%.%%USERDOMAIN%%.INI.%s\\resource.dll\r\n"
		// simple checks
		exists, errFileExists := fileExists(fullTokenPath)
		if errFileExists != nil {
			return errFileExists
		}

		if !exists || OverwriteFileIfExists {
			if exists {
				c.l.WithField("file", fullTokenPath).Warn("file exists and will be overwritten! ('-overwrite-files' is set to true)")
			}
			// Create the directory
			err = os.MkdirAll(fullTokenPath, 0755)
			if err != nil {
				return err
			}
			// Create the file
			fullTokenPathINI := filepath.Join(fullTokenPath, "desktop.ini")
			out, err := os.Create(fullTokenPathINI)
			if err != nil {
				return err
			}

			// windows loves UTF-16LE with BOM
			var bytes [2]byte
			const BOM = '\ufffe' //LE. for BE '\ufeff'
			bytes[0] = BOM >> 8
			bytes[1] = BOM & 255
			_, err = out.Write(bytes[0:])
			if err != nil {
				return err
			}
			runes := utf16.Encode([]rune(fmt.Sprintf(iniTemplate, tcr.Canarytoken.Hostname)))
			for _, r := range runes {
				bytes[1] = byte(r >> 8)
				bytes[0] = byte(r & 255)
				_, err = out.Write(bytes[0:])
				if err != nil {
					return err
				}
			}
			// Write the body to file
			// _, err = out.WriteString(fmt.Sprintf(iniTemplate, tcr.Canarytoken.Hostname))
			out.Close()
			if err != nil {
				return err
			}
			// set the file to be hidden
			err = SetFileAttributeHiddenAndSystem(fullTokenPathINI)
			if err != nil {
				return err
			}
			// Setting the dir to system
			err = SetFileAttributeSystem(fullTokenPath)
			if err != nil {
				return err
			}
		}
		if exists && !OverwriteFileIfExists { // id DOES exist, and you told me not to overwrite
			return fmt.Errorf("file exists: %s, and '-overwrite-file' is false", fullTokenPath)
		}

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
		exists, errFileExists := fileExists(fullTokenPath)
		if errFileExists != nil {
			return errFileExists
		}

		if !exists || OverwriteFileIfExists {
			if exists {
				c.l.WithField("file", fullTokenPath).Warn("file exists and will be overwritten! ('-overwrite-files' is set to true)")
			}
			// Create the file
			out, err := os.Create(fullTokenPath)
			if err != nil {
				return err
			}
			defer out.Close()

			// Write the body to file
			_, err = out.WriteString(fmt.Sprintf(aswTemplate, tcr.Canarytoken.AccessKeyID, tcr.Canarytoken.SecretAccessKey))
		}
		if exists && !OverwriteFileIfExists { // id DOES exist, and you told me not to overwrite
			return fmt.Errorf("file exists: %s, and '-overwrite-file' is false", fullTokenPath)
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
		_, err = c.DownloadTokenFromAPI(tcr.Canarytoken.Canarytoken, fullTokenPath, OverwriteFileIfExists)
	default:
		err = fmt.Errorf("unsupported Canarytoken: %s", kind)
	}
	return
}

// CreateFactory creates an auth string for the Canarytoken Factory endpoint.
func (c Client) CreateFactory(memo string) (createfactoryresponse CreateFactoryResponse, err error) {
	createfactoryresponse = CreateFactoryResponse{}
	u := &url.Values{}
	u.Set("memo", memo)

	err = c.decodeResponse("canarytoken/create_factory", "POST", u, &createfactoryresponse)
	if err != nil {
		return
	}

	if createfactoryresponse.Result != "success" {
		err = fmt.Errorf("error creating token: %s", err)
	}
	return
}

// CreateTokenFromFactory uses the canarytoken/factory API endpoint to create a token
func (c Client) CreateTokenFromFactory(kind, memo, FlockID string, additionalParams *url.Values) (tokencreateresponse TokenCreateResponse, err error) {
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
	case "http", "dns", "cloned-web", "doc-msword", "web-image", "windows-dir", "pdf-acrobat-reader", "msword-macro", "msexcel-macro", "aws-id", "qr-code", "fast-redirect", "slow-redirect", "slack-api":
	// TODO: must check additional params per kind
	default:
		return tokencreateresponse, errors.New("unsupported token type: " + kind)
	}

	err = c.decodeResponse("canarytoken/factory", "POST", u, &tokencreateresponse)
	if err != nil {
		return
	}
	if tokencreateresponse.Result != "success" {
		err = fmt.Errorf("error creating token: %s", err)
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

	var apiEndpoint string

	switch c.opmode {
	case "api":
		apiEndpoint = "canarytoken/create"
	case "factory":
		apiEndpoint = "canarytoken/factory"
	default:
		err = errors.New("unsupported opmode: " + c.opmode)
		return
	}
	err = c.decodeResponse(apiEndpoint, "POST", u, &tokencreateresponse)
	if err != nil {
		return
	}
	if tokencreateresponse.Result != "success" {
		err = fmt.Errorf("error creating token: %s", err)
	}
	return
}

func (c Client) DownloadWindowsDirTokenFromAPI(canarytoken, dropWhere, filename string, OverwriteFileIfExists bool) (n int64, err error) {
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

	// create temp folder
	tmpdir, err := ioutil.TempDir("", "tokendropper")
	if err != nil {
		return 0, err
	}
	defer os.RemoveAll(tmpdir)
	tmpFilename := filepath.Join(tmpdir, filename+".zip")
	// Create the file
	out, err := os.Create(tmpFilename)
	if err != nil {
		return 0, err
	}
	defer out.Close()
	n, err = io.Copy(out, resp.Body)
	if err != nil {
		return
	}
	// now we have the zip file in a temp folder
	// we need to unzip it.
	c.l.WithFields(log.Fields{
		"tmpFilename": tmpFilename,
		"tmpdir":      tmpdir,
	}).Debug("unzipping windows-dir Canarytoken")
	filenames, err := Unzip(tmpFilename, tmpdir)
	if err != nil {
		return
	}
	for _, filename := range filenames {
		c.l.WithField("file", filename).Debug("file extracted from zip")
	}
	oldpath := filepath.Join(tmpdir, "My Documents")
	newfoldername := filename
	// newpath := filepath.Join(oldpath, newfoldername)

	// full path of token
	tokenFullDirPath := filepath.Join(dropWhere, filename)
	c.l.WithFields(log.Fields{
		"oldpath":          oldpath,
		"newfoldername":    newfoldername,
		"tokenFullDirPath": tokenFullDirPath,
		// "newpath":       newpath,
	}).Debug("renaming default windows-dir default directory")

	// check if new folder does not already exist
	exists, err := fileExists(tokenFullDirPath)
	if err != nil {
		return
	}
	if !exists || OverwriteFileIfExists {
		err = os.MkdirAll(tokenFullDirPath, 0755)
		if err != nil {
			return
		}
		// Read contents

		err = os.Rename(oldpath, tokenFullDirPath)
		if err != nil {
			return
		}
	}
	// defer shoud take or of cleaning up the tmpdir
	return
}

// DownloadTokenFromAPI downloads a file-based token given its ID
func (c Client) DownloadTokenFromAPI(canarytoken, filename string, OverwriteFileIfExists bool) (n int64, err error) {
	params := &url.Values{}
	params.Set("canarytoken", canarytoken)

	var apiEndpoint string

	switch c.opmode {
	case "api":
		apiEndpoint = "canarytoken/download"
	case "factory":
		apiEndpoint = "canarytoken/factory/download"
	default:
		err = errors.New("unsupported opmode: " + c.opmode)
		return
	}

	fullURL, err := c.api(apiEndpoint, params)
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

	exists, err := fileExists(filename)
	if err != nil {
		return
	}
	if !exists || OverwriteFileIfExists {
		if exists {
			c.l.WithField("file", filename).Warn("file exists and will be overwritten! ('-overwrite-files' is set to true)")
		}
		// Create the file
		out, err := os.Create(filename)
		if err != nil {
			return 0, err
		}
		defer out.Close()

		// Write the body to file
		n, err = io.Copy(out, resp.Body)
	}
	if exists && !OverwriteFileIfExists { // it DOES exist, and you told me not to overwrite
		return 0, fmt.Errorf("file exists: %s, and '-overwrite-file' is false", filename)
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

	switch c.opmode {
	case "api":
		// always add auth token to list of values
		params.Add("auth_token", c.apikey)
	case "factory":
		// always add auth token to list of values
		params.Add("factory_auth", c.factoryAuth)
	default:
		c.l.WithField("opmode", c.opmode).Fatal("unsupported client opmode")
	}

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

func (c Client) SearchAllIncidents(state string) (respIncidents []interface{}, err error) {
	return c.searchIncidents("", "", false, state)
}

func (c Client) SearchFilteredIncidents(filter string, id string, state string) (respIncidents []interface{}, err error) {
	return c.searchIncidents(filter, id, true, state)
}

// SearchIncidents returns all Incidents specified by filter
// filter can be "flock_id" or "node_id"
func (c Client) searchIncidents(filter string, id string, withFilter bool, state string) (respIncidents []interface{}, err error) {
	respIncidents = make([]interface{}, 0)
	resp := IncidentSearchResponse{}

	var u = &url.Values{}
	u.Add("limit", "1000")

	switch state {
	case "all", "acknowledged", "unacknowledged":
		u.Add("filter_incident_state", state)
	default:
		err = fmt.Errorf("unsupported Incident State: %s", state)
		return
	}

	if withFilter {
		u.Add(filter, id)
	}

	for {
		err = c.decodeResponse("incidents/search", "GET", u, &resp)
		if err != nil {
			return
		}

		if resp.Result != "success" {
			return nil, errors.New(resp.Message) // there will be a message, if it failed
		}

		c.l.WithFields(log.Fields{
			"page_number": resp.PageNumber,
			"total_pages": resp.TotalPages,
		}).Infof("incidents fetched")
		for _, i := range resp.Incidents {
			respIncidents = append(respIncidents, i)
		}

		// we done?
		if resp.TotalPages == resp.PageNumber {
			break
		}
		u.Set("cursor", resp.Cursor.Next)
		u.Del("limit")
	}

	return
}

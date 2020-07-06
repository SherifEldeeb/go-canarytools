package canarytools

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// ConsoleAPIFeeder is a feeder that fetches Incidents from the console API client
type ConsoleAPIFeeder struct {
	*Client
	lastCheck         time.Time
	lastCheckRegister *os.File
	errorCount        int
	fetchInterval     int
	thenWhat          string
	whichIncidents    string
	lastCheckFile     string
}

// NewConsoleAPIFeeder creates a new ConsolevAPI feeder
func NewConsoleAPIFeeder(domain, apikey, thenWhat, sinceWhen, whichIncidents string, fetchInterval int, l *log.Logger) (c *ConsoleAPIFeeder, err error) {
	c = &ConsoleAPIFeeder{}
	c.Client, err = NewClient(domain, apikey, l)
	if err != nil {
		return
	}
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

func (c *ConsoleAPIFeeder) WriteLastCheckRegister(t time.Time) (err error) {
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

// Feed fetches incidents and feeds them to chan
func (c *ConsoleAPIFeeder) Feed(incidnetsChan chan<- Incident) {
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

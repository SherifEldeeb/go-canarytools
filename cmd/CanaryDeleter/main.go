package main

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	canarytools "github.com/thinkst/go-canarytools"
)

var (
	cfg canarytools.CanaryDeleterConfig
	err error

	// constants for build-time hardcoding of params
	// those could be set at build time to create a self-running executable
	// go build -ldflags "-X main.DOMAIN=$domain_hash  -X main.APIKEY=$api_auth -w -s -linkmode=internal"

	// DOMAIN is the domain hash
	DOMAIN string
	// APIKEY is the main console API auth token
	APIKEY string
	// BUILDTIME is the time the tools was built
	BUILDTIME string
	// SHA1VER is the built sha1
	SHA1VER string
)

func init() {
	populateVarsFromFlags(&cfg)
}

func main() {
	flag.Parse()
	l := log.New()
	l.WithFields(log.Fields{
		"BUILDTIME": BUILDTIME,
		"SHA1VER":   SHA1VER,
	}).Info("starting CanaryDeleter")

	// Finish config logic
	err = finishConfig(&cfg, l)
	if err != nil {
		l.WithField("err", err).Fatal("configuration error")
	}

	// Get a client
	l.Info("building an API client")
	c, err := canarytools.NewClient(cfg.ConsoleAPIConfig, l)
	if err != nil {
		l.Fatal(err)
	}

	if cfg.FlockName != "" {
		l.WithField("FlockName", cfg.FlockName).Info("getting flock_id from FlockName")
		// does the flock exist?
		flockID, err := c.GetFlockIDFromName(cfg.FlockName)
		if err != nil {
			l.Fatal(err)
		}
		cfg.FlockID = flockID
		l.WithField("FlockName", cfg.FlockName).WithField("flock_id", cfg.FlockID).Info("got flock_id")

	}

	var id string
	var filter string
	switch cfg.DeleteWhat {
	case "incidents":
		if cfg.DumpToJson {
			l.Infof("dumping incidents to json file before deleting them")
			l.Info("fetching incidents ... this might take a while")
			switch cfg.FilterType {
			case "flock_id":
				id = cfg.FlockID
				filter = "flock_id"
				l.WithField("flock_id", cfg.FlockID).Info("filtering incidents using Flock ID")
			case "node_id":
				id = cfg.NodeID
				filter = "node_id"
				l.WithField("node_id", cfg.NodeID).Info("filtering incidents using Node ID")
			default:
				l.Fatal("unsupported filter type:" + cfg.FilterType)
			}
			incidents, err := c.SearchIncidents(filter, id)
			if err != nil {
				l.Fatal(err)
			}
			l.WithField("incidents_count", len(incidents)).Info("fetching incidents done!")

			if len(incidents) > 0 {
				filename := "canary-" + time.Now().UTC().Format("2006-01-02_15-04-05") + ".json"
				l.WithField("filename", filename).Infof("opening file for writing")
				f, err := os.Create(filename)
				if err != nil {
					l.Fatal(err)
				}
				defer f.Close()
				for _, i := range incidents {
					j, err := json.Marshal(i)
					if err != nil {
						l.WithField("err", err).Error("error marshaling incident")
					}
					f.Write(j)
					f.Write([]byte("\n"))
				}
			} else {
				l.Info("no incidents found! gonna bail out.")
				os.Exit(0)
			}
		}
		l.WithField("id", id).WithField("filter", filter).Info("deleting all incidents")
		err = c.DeleteMultipleIncidents(filter, id, cfg.IncludeUnacknowledged)
		if err != nil {
			l.Fatal(err)
		}
	case "tokens":
		t, err := c.FetchCanarytokenAll()
		if err != nil {
			l.Fatal(err)
		}
		l.WithField("FlockName", cfg.FlockName).WithField("flock_id", cfg.FlockID).Info("deleting all tokens")
		for _, token := range t {
			if token.FlockID == cfg.FlockID {
				l.Info("deleteing:", token.Canarytoken)
				err = c.DeleteCanarytoken(token.Canarytoken)
				if err != nil {
					l.Error("error:", err)
				}
			}
		}
	default:
		l.Fatal("you have to tell me what to delete using '-what' ... supported values are 'incidents' & 'tokens'")
	}
	l.Info("done!")

}

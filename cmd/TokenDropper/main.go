//go:generate goversioninfo
package main

import (
	"flag"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	canarytools "github.com/thinkst/go-canarytools"
)

var (
	tokendropper canarytools.TokenDropper
	cfg          canarytools.TokenDropperConfig
	err          error
)

func init() {
	rand.Seed(time.Now().UnixNano())
	populateVarsFromFlags(&cfg)
}

func main() {
	flag.Parse()

	// overwrite from flag
	for _, k := range strings.Split(cfg.KindsStr, ",") {
		cfg.Kinds = append(cfg.Kinds, strings.TrimSpace(k))
	}

	// Set LogLevel
	l := log.New()
	switch cfg.LogLevel {
	case "info":
		l.SetLevel(log.InfoLevel)
	case "warning":
		l.SetLevel(log.WarnLevel)
	case "debug":
		l.SetLevel(log.DebugLevel)
	case "trace":
		l.SetLevel(log.TraceLevel)
	default:
		l.Warn("unsupported log level, or none specified; will set to 'info'")
		l.SetLevel(log.InfoLevel)
	}

	// start
	// dpending on the execution environment, sometimes "./" does not get evaluated as "same dir as the exe file"
	// so, till I figure out a better way, we do the following.
	if cfg.DropWhere == "./" {
		p, err := os.Executable()
		if err != nil {
			l.Fatal("couldn't get current directory")
		}
		cfg.DropWhere = filepath.Dir(p) // full path to executable
	}
	err = os.Chdir(cfg.DropWhere)
	if err != nil {
		l.WithField("err", err).Fatal("couldn't change directory")
	}

	l.WithField("where", cfg.DropWhere).Info("Dropping Canarytokens")

	if cfg.FilesCount > 30 {
		l.Warn("File count is > 30 ... will set to 30")
	}

	// try to populte domain hash and API key
	// either from file or params...
	// first, we didn't get api key and domain through flags? let's try to load them from file
	if cfg.ConsoleAPIKey == "" && cfg.ConsoleAPIDomain == "" {
		// if we don't have them, we try to load it from same drectory
		if cfg.ConsoleTokenFile == "" { // if not
			cfg.ConsoleTokenFile = filepath.Join(cfg.DropWhere, "canarytools.config")
		}
		// do we have canarytools.config in same path? get data from it...
		if _, err := os.Stat(cfg.ConsoleTokenFile); os.IsNotExist(err) {
			l.Fatal("canarytools.config does not exist, and we couldn't get domain hash and API key!")
		}
		cfg.ConsoleAPIKey, cfg.ConsoleAPIDomain, err = canarytools.LoadTokenFile(cfg.ConsoleTokenFile)
		if err != nil || cfg.ConsoleAPIDomain == "" || cfg.ConsoleAPIKey == "" {
			l.WithFields(log.Fields{
				"err":    err,
				"api":    cfg.ConsoleAPIKey,
				"domain": cfg.ConsoleAPIDomain,
			}).Fatal("error parsing token file")
		}
	}
	// by now, we should have both key and domain

	// create new canary API client
	c, err := canarytools.NewClient(cfg.ConsoleAPIDomain, cfg.ConsoleAPIKey, l)
	if err != nil {
		l.WithField("err", err).Fatal("error creating canary client")
	}

	// flock related stuff
	// ultimate goal to populate both FlockName & FlockID
	// if provided name exists, retrieve the FlockID,
	// if it doesn't, and CreateFlockIfNotExist id true
	// create it, and set the FlockID
	if cfg.FlockName == "" { // if no flock name provided, use the default one
		cfg.FlockID = "flock:default"
		cfg.FlockName = "Default Flock"
	} else { // we have been given a flockname...
		// does it exist?
		exists, fid, err := c.FlockNameExists(cfg.FlockName)
		if err != nil {
			l.WithField("err", err).Fatal("error checking if flock exists")
		}
		if exists {
			cfg.FlockID = fid
		} else {
			l.WithField("flockname", cfg.FlockName).Info("flock does not exist")
			if cfg.CreateFlockIfNotExists {
				l.WithField("flockname", cfg.FlockName).Info("creating flock")
				cfg.FlockID, err = c.FlockCreate(cfg.FlockName)
				if err != nil {
					l.WithField("err", err).Fatal("error creating flock")
				}
			} else {
				l.WithField("flockname", cfg.FlockName).Fatal("flock doesn't exist, and you told me not to create it")
			}
		}
	}
	// we now should have both FlockName & FlockID
	// let's get this over with...
	log.WithFields(log.Fields{
		"kind":  cfg.KindsStr,
		"count": cfg.FilesCount,
		"flock": cfg.FlockName,
	}).Info("dropping tokens..")
	for i := 0; i < cfg.FilesCount; i++ {
		kind := pick(cfg.Kinds)
		filename, err := GetRandomTokenName(kind)
		if err != nil {
			l.Error(err)
			continue
		}
		l.WithFields(log.Fields{
			"kind":     kind,
			"filename": filename,
		}).Info("Generating Token")
		memo, err := CreateMemo(filename, cfg.CustomMemo)
		if err != nil {
			l.Error(err)
			continue
		}

		l.WithFields(log.Fields{
			"kind":     kind,
			"filename": filename,
			"memo":     memo,
		}).Debug("Generating Token")
		// drop
		filename = filepath.Join(cfg.DropWhere, filename)
		err = c.DropFileToken(kind, memo, filename, cfg.FlockID, cfg.CreateFlockIfNotExists)
		if err != nil {
			l.Error(err)
			continue
		}
		rtime := GetRandomDate(cfg.RandYearsBack)
		err = os.Chtimes(filename, rtime, rtime)
		if err != nil {
			l.WithFields(log.Fields{
				"filename": filename,
				"err":      err,
			}).Error("Error changing file timestamps")
		}
	}
}

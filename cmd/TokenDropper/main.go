//go:generate goversioninfo
package main

import (
	"flag"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/SherifEldeeb/canarytools"
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
		l.Warn("unsupported log level, or none specified; will set to 'info', ")
		l.SetLevel(log.InfoLevel)
	}

	// start
	p, err := os.Executable()
	if err != nil {
		l.Fatal("couldn't get current directory")
	}
	l.Info("running from:", p)

	d := filepath.Dir(p)  // directory
	e := filepath.Base(p) // exe name
	// get count from exe name?
	countFromName, err := strconv.Atoi(e)
	if err == nil {
		if countFromName > 15 {
			l.Warn("count from name > 15 ... will set to 15")
			countFromName = 15
		}
		cfg.FilesCount = countFromName
	}

	err = os.Chdir(d)
	if err != nil {
		l.Fatal("couldn't change directory")
	}

	// try to populte domain hash and API key
	// either from file or params...
	// first, we don't got api key and domain through flags? let's try to load them from ile
	if cfg.ImConsoleAPIKey == "" && cfg.ImConsoleAPIDomain == "" {
		// if we don't have them, we try to load it from same drectory
		if cfg.ImConsoleTokenFile == "" { // if not
			cfg.ImConsoleTokenFile = filepath.Join(d, "canarytools.config")
		}
		// do we have canarytools.config in same path? get data from it...
		if _, err := os.Stat(cfg.ImConsoleTokenFile); os.IsNotExist(err) {
			l.Fatal("canarytools.config does not exist, and we couldn't get domain hash and API key!")
		}
		cfg.ImConsoleAPIKey, cfg.ImConsoleAPIDomain, err = canarytools.LoadTokenFile(cfg.ImConsoleTokenFile)
		if err != nil || cfg.ImConsoleAPIDomain == "" || cfg.ImConsoleAPIKey == "" {
			l.WithFields(log.Fields{
				"err":    err,
				"api":    cfg.ImConsoleAPIKey,
				"domain": cfg.ImConsoleAPIDomain,
			}).Fatal("error parsing token file")
		}
	}
	// by now, we should have both key and domain

	// filename is number of files?
	c, err := canarytools.NewClient(cfg.ImConsoleAPIDomain, cfg.ImConsoleAPIKey, l)
	if err != nil {
		l.Fatal(err)
	}

	fileCount := cfg.FilesCount
	log.Info(fileCount)
	for i := 0; i < fileCount; i++ {
		kind := pick(cfg.Kinds)
		n, err := GetRandomTokenName(kind)
		if err != nil {
			l.Error(err)
			continue
		}
		l.WithFields(log.Fields{
			"kind":     kind,
			"filename": n,
		}).Info("Generating Token")
		memo, err := CreateMemo(n)
		if err != nil {
			l.Error(err)
			continue
		}

		l.WithFields(log.Fields{
			"kind":     kind,
			"filename": n,
			"memo":     memo,
		}).Debug("Generating Token")
		// drop
		n = filepath.Join(cfg.DropWhere, n)
		err = c.DropFileToken(kind, memo, "", n)
		if err != nil {
			l.Error(err)
			continue
		}
		rtime := GetRandomDate(2)
		os.Chtimes(n, rtime, rtime)
	}
}

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-logfmt/logfmt"

	log "github.com/sirupsen/logrus"

	canarytools "github.com/thinkst/go-canarytools"
)

var (
	cfg canarytools.ConsoleAPIConfig
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
	l := log.New()
	l.SetLevel(log.InfoLevel)

	l.WithFields(log.Fields{
		"BUILDTIME": BUILDTIME,
		"SHA1VER":   SHA1VER,
	}).Info("Starting CSV TokenLister")

	flag.Parse()

	// Set all hardcoded info, if provided
	if DOMAIN != "" && cfg.ConsoleAPIDomain == "" { // command line values always Supersede hardcoded ones
		l.Info("Found pre-configured domain hash value")
		cfg.ConsoleAPIDomain = DOMAIN
	}
	if APIKEY != "" && cfg.ConsoleAPIKey == "" { // command line values always Supersede hardcoded ones
		l.Info("Found pre-configured API auth value")
		cfg.ConsoleAPIKey = APIKEY
	}

	// final checks
	// try to populte domain hash and API key
	// either from file or params...
	// first, we didn't get api key and domain through flags? let's try to load them from file
	if cfg.ConsoleAPIKey == "" && cfg.ConsoleAPIDomain == "" {
		// do we have canarytools.config in same path? get data from it...
		if _, err := os.Stat(cfg.ConsoleTokenFile); os.IsNotExist(err) {
			l.Fatal("canarytools.config does not exist, and we couldn't get domain hash and API key")
		}
		cfg.ConsoleAPIKey, cfg.ConsoleAPIDomain, err = canarytools.LoadTokenFile(cfg.ConsoleTokenFile)
		if err != nil || cfg.ConsoleAPIDomain == "" || cfg.ConsoleAPIKey == "" {
			l.Fatalf("error parsing token file: %s", err)
		}
	}

	l.Info("Creating a console API client & pinging your console...")

	c, err := canarytools.NewClient(cfg.ConsoleAPIDomain, cfg.ConsoleAPIKey, cfg.OpMode, l)
	if err != nil {
		l.Fatal(err)
	}

	l.Info("Console response looks good, let's get started... ")

	// get flock names
	var flockMapping = make(map[string]string)
	l.Info("Getting flocks info...")
	flocksummary, err := c.GetFlocksSummary()
	if err != nil {
		l.Fatal(err)
	}
	if flocksummary.Result != "success" {
		l.Fatalf("error fetching flocks summary: %s", flocksummary.Message)
	}

	l.Infof("Found total of %d flock(s)", len(flocksummary.FlocksSummary))
	for flock, summary := range flocksummary.FlocksSummary {
		flockMapping[flock] = summary.Name
	}

	// we're good to go
	l.Info("Fetching tokens (this might take a while) ...")
	t, err := c.FetchCanarytokenAll()
	if err != nil {
		l.Fatal(err)
	}

	l.Infof("Done! fetched a total of: %d tokens", len(t))

	// create file at same location as the program
	filename := "canary-tokens_" + time.Now().Format("2006-01-02_15_04_05") + ".csv"
	l.Infof("Opening file for writing: %s...", filename)
	f, err := os.Create(filename)
	if err != nil {
		l.Fatal(err)
	}
	defer f.Close()

	// create the csv writer
	csvWriter := csv.NewWriter(f)
	defer csvWriter.Flush()

	l.Info("Writing CSV header...")
	// write csv header
	err = csvWriter.Write([]string{
		"Created",
		"Canarytoken",
		"Kind",
		"OriginalFilename",
		"Where",
		"FlockName",
		"Generator",
		"Memo",
		"TDUser",
		"TDHost",
		"TDOS",
		"FullRawMemo",
	})
	if err != nil {
		l.Fatal(err)
	}

	// iterate over tokens
	l.Info("Writing CSV entries\n")
	for _, token := range t {
		l.Debug("adding:", token.Canarytoken)
		fmt.Print(".")
		var Generator, Memo, TDUser, TDHost, OriginalFilename, Where, TDOS string
		// unmarshal memo
		logfmtDecoder := logfmt.NewDecoder(strings.NewReader(token.Memo))
		for logfmtDecoder.ScanRecord() {
			for logfmtDecoder.ScanKeyval() {
				k := string(logfmtDecoder.Key())
				v := string(logfmtDecoder.Value())
				switch k {
				case "Generator":
					Generator = v
				case "Memo":
					Memo = v
				case "TD-User":
					TDUser = v
				case "TD-Host":
					TDHost = v
				case "OriginalFilename":
					OriginalFilename = v
				case "Where":
					Where = v
				case "TD-OS":
					TDOS = v
				}
			}
		}

		err = csvWriter.Write([]string{
			token.CreatedPrintable,
			token.Canarytoken,
			token.Kind,
			OriginalFilename,
			Where,
			flockMapping[token.FlockID],
			Generator,
			Memo,
			TDUser,
			TDHost,
			TDOS,
			token.Memo,
		})
		if err != nil {
			l.Error("error writing csv entry:", err)
		}
	}
	fmt.Println()
	l.Info("Done! enjoy the rest of your day...")
}

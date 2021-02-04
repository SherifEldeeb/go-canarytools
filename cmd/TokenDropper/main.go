//go:generate goversioninfo
package main

import (
	"errors"
	"flag"
	"fmt"
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

	// constants for build-time hardcoding of params
	// those could be set at build time to create a self-running executable
	// go build -ldflags "-X main.DOMAIN=$domain_hash  -X main.APIKEY=$api_auth -w -s -linkmode=internal"

	// DOMAIN is the domain hash
	DOMAIN string
	// APIKEY is the main console API auth token
	APIKEY string
	// FACTORYAUTH is the factory auth token
	FACTORYAUTH string
	// BUILDTIME is the time the tools was built
	BUILDTIME string
	// SHA1VER is the built sha1
	SHA1VER string
)

func init() {
	rand.Seed(time.Now().UnixNano())
	populateVarsFromFlags(&cfg)
}

func main() {
	l := log.New()
	l.WithFields(log.Fields{
		"BUILDTIME": BUILDTIME,
		"SHA1VER":   SHA1VER,
	}).Info("starting TokenDropper")

	flag.Parse()

	// Set LogLevel
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

	// Set all hardcoded info, if provided
	if DOMAIN != "" && cfg.ConsoleAPIDomain == "" { // command line values always Supersede hardcoded ones
		l.Debug("found pre-configured domain hash value")
		cfg.ConsoleAPIDomain = DOMAIN
	}
	if APIKEY != "" && cfg.ConsoleAPIKey == "" { // command line values always Supersede hardcoded ones
		l.Debug("found pre-configured API auth value")
		cfg.ConsoleAPIKey = APIKEY
	}
	if FACTORYAUTH != "" && cfg.ConsoleFactoryAuth == "" { // command line values always Supersede hardcoded ones
		l.Debug("found pre-configured factory auth value")
		cfg.ConsoleFactoryAuth = FACTORYAUTH
	}

	// overwrite kinds from flag after splitting
	for _, k := range strings.Split(cfg.KindsStr, ",") {
		cfg.Kinds = append(cfg.Kinds, strings.TrimSpace(k))
	}

	// Finish config logic
	err = finishConfig(&cfg, l)
	if err != nil {
		l.WithField("err", err).Fatal("configuration error")
	}

	// by now, we should have both key and domain

	// create new canary API client
	c, err := canarytools.NewClient(cfg.ConsoleAPIConfig, l)
	if err != nil {
		l.WithField("err", err).Fatal("error creating canary client")
	}

	// flock related stuff
	// ultimate goal to populate both FlockName & FlockID
	// if provided name exists, retrieve the FlockID,
	// if it doesn't, and CreateFlockIfNotExist id true
	// create it, and set the FlockID
	if cfg.FlockName == "" && cfg.FlockID == "" { // if no flock name or ID provided, use the default one
		cfg.FlockID = "flock:default"
		cfg.FlockName = "Default Flock"
	} else { // we have been given a flockname, or ID ...
		switch cfg.FlockID {
		case "":
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
		default:
			// flockname already provided
		}
	}
	// we now should have both FlockName & FlockID
	// let's get this over with...
	log.WithFields(log.Fields{
		"kind":  cfg.KindsStr,
		"count": cfg.FilesCount,
		"flock": cfg.FlockName,
	}).Info("dropping tokens..")
	// kind := pick(cfg.Kinds)
	for i := 0; i < cfg.FilesCount; i++ {
		for _, kind := range cfg.Kinds {
			var filename string
			if cfg.FileName != "" {
				filename = cfg.FileName
			} else {
				filename, err = GetRandomTokenName(kind, cfg.RandomizeFilenames)
				if err != nil {
					l.Error(err)
					continue
				}
			}
			memo, err := CreateMemo(filename, cfg.DropWhere, cfg.CustomMemo)
			if err != nil {
				l.Error(err)
				continue
			}

			l.WithFields(log.Fields{
				"kind":     kind,
				"filename": filename,
				"memo":     memo,
				"where":    cfg.DropWhere,
			}).Debug("Generating Token")

			// drop
			err = c.DropFileToken(kind, memo, cfg.DropWhere, filename, cfg.FlockID, cfg.CreateFlockIfNotExists, cfg.CreateDirectoryIfNotExists, cfg.OverwriteFileIfExists)
			if err != nil {
				l.Error(err)
				continue
			}

			fullFilePath := filepath.Join(cfg.DropWhere, filename)

			rtime := GetRandomDate(cfg.RandYearsBack)
			l.WithFields(log.Fields{
				"fullFilePath": fullFilePath,
				"fake_time":    rtime.UTC().String(),
			}).Debug("changing timestamps")
			err = os.Chtimes(fullFilePath, rtime, rtime)
			if err != nil {
				l.WithFields(log.Fields{
					"filename": filename,
					"err":      err,
				}).Error("Error changing file timestamps")
			}
		}
	}
}

func finishConfig(cfg *canarytools.TokenDropperConfig, l *log.Logger) (err error) {
	// TODO: a big one ... lots of factory logic
	// start
	// dpending on the execution environment, sometimes "./" does not get evaluated as "same dir as the exe file"
	// so, till I figure out a better way, we do the following.
	if (cfg.FactoryAuthFile != "" || cfg.ConsoleFactoryAuth != "") && (cfg.ConsoleAPIKey != "" || cfg.ConsoleTokenFile != "") {
		return errors.New("can't specify both factory && console API")
	}
	if cfg.FactoryAuthFile != "" || cfg.ConsoleFactoryAuth != "" {
		cfg.OpMode = "factory"
	}
	if cfg.DropWhere == "./" {
		p, err := os.Executable()
		if err != nil {
			return fmt.Errorf("couldn't get current directory")
		}
		cfg.DropWhere = filepath.Dir(p) // full path to executable
	}
	// you shouldn't specify flockname && flock_id
	if cfg.FlockName != "" && cfg.FlockID != "" {
		return errors.New("can't specify both -flock && -flockid")
	}
	if (cfg.FactoryAuthFile != "" || cfg.ConsoleFactoryAuth != "") && cfg.FlockName != "" {
		return errors.New("can't use factory auth with '-flock' name ... you MUST use '-flockid'")
	}
	// TODO: remove from ConsoleAPI.go
	// check if 'where' directory exists
	// if it doesn't exist, and CreateDirectoryIfNotExists is true, create it
	// if it doesn't exist, and CreateDirectoryIfNotExists is false, error out
	cfg.DropWhere, err = filepath.Abs(cfg.DropWhere)
	if err != nil {
		return
	}
	if _, errstat := os.Stat(cfg.DropWhere); os.IsNotExist(errstat) { // it does NOT exist
		if cfg.CreateDirectoryIfNotExists {
			os.MkdirAll(cfg.DropWhere, 0755)
		} else {
			err = fmt.Errorf("'where' does not exist, and you told me not to create it ... gonna have to bail out")
			return
		}
	}

	err = os.Chdir(cfg.DropWhere)
	if err != nil {
		return fmt.Errorf("couldn't change directory: %s", err)
	}

	l.WithField("where", cfg.DropWhere).Info("Dropping Canarytokens")

	if cfg.FilesCount > 20 {
		l.Warn("File count is > 20 ... will set to 20")
		cfg.FilesCount = 20
	}

	// try to populte domain hash and API key
	// either from file or params...
	switch cfg.OpMode {
	case "api":
		// first, we didn't get api key and domain through flags? let's try to load them from file
		if cfg.ConsoleAPIKey == "" && cfg.ConsoleAPIDomain == "" {
			// if we don't have them, we try to load it from same drectory
			if cfg.ConsoleTokenFile == "" { // if not
				cfg.ConsoleTokenFile = filepath.Join(cfg.DropWhere, "canarytools.config")
			}
			// do we have canarytools.config in same path? get data from it...
			if _, err := os.Stat(cfg.ConsoleTokenFile); os.IsNotExist(err) {
				return fmt.Errorf("canarytools.config does not exist, and we couldn't get domain hash and API key")
			}
			cfg.ConsoleAPIKey, cfg.ConsoleAPIDomain, err = canarytools.LoadTokenFile(cfg.ConsoleTokenFile)
			if err != nil || cfg.ConsoleAPIDomain == "" || cfg.ConsoleAPIKey == "" {
				return fmt.Errorf("error parsing token file: %s", err)
			}
		}
	case "factory":
		// first, we didn't get api key and domain through flags? let's try to load them from file
		if cfg.ConsoleFactoryAuth == "" && cfg.ConsoleAPIDomain == "" {
			// if we don't have them, we try to load it from same drectory
			if cfg.FactoryAuthFile == "" { // if not
				cfg.FactoryAuthFile = filepath.Join(cfg.DropWhere, "canarytools.config")
			}
			// do we have canarytools.config in same path? get data from it...
			if _, err := os.Stat(cfg.FactoryAuthFile); os.IsNotExist(err) {
				return fmt.Errorf("canarytools.config does not exist, and we couldn't get domain hash and API key")
			}
			cfg.ConsoleFactoryAuth, cfg.ConsoleAPIDomain, err = canarytools.LoadTokenFile(cfg.FactoryAuthFile)
			if err != nil || cfg.ConsoleAPIDomain == "" || cfg.ConsoleFactoryAuth == "" {
				return fmt.Errorf("error parsing factory auth token file: %s", err)
			}
		}
	default:
		return fmt.Errorf("unknown opmode: %s", cfg.OpMode)
	}

	// checks if they specified the filename
	if cfg.FileName != "" {
		if cfg.RandomizeFilenames {
			l.Warn("filename has been specified, but 'randomize-filenames' is set to 'true' ... will be setting it to 'false'")
			cfg.RandomizeFilenames = false
		}
		if len(cfg.Kinds) > 1 {
			return fmt.Errorf(`you provided multipe 'kind' of tokens, yet specified a filename ... you can't do that
you should either set one 'kind', or remove 'filename' and the tool will generate a random one for you`)
		}
		if cfg.FilesCount != 1 {
			l.Warn("filename has been specified, but 'count' is not '1' ... setting it to 1")
			cfg.FilesCount = 1
		}
		// extension match kind?
		e := strings.TrimPrefix(filepath.Ext(cfg.FileName), ".")

		switch cfg.Kinds[0] { // this should be only one
		case "doc-msword":
			if e != "docx" {
				return fmt.Errorf("kind is 'doc-msword', yet the filename provided does not have a .docx extension")
			}
		case "msexcel-macro":
			if e != "xlsm" {
				return fmt.Errorf("kind is 'msexcel-macro', yet the filename provided does not have a .xlsm extension")
			}
		case "msword-macro":
			if e != "docm" {
				return fmt.Errorf("kind is 'msword-macro', yet the filename provided does not have a .docm extension")
			}
		case "pdf-acrobat-reader":
			if e != "pdf" {
				return fmt.Errorf("kind is 'pdf-acrobat-reader', yet the filename provided does not have a .pdf extension")
			}
		case "aws-id":
			if e != "txt" {
				if e != "" { // no extension?
					return fmt.Errorf("kind is 'aws-id', yet the filename provided does not have either a .txt extension, or no extension")
				}
			}
		}
	}
	return
}

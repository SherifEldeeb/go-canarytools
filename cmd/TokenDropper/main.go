//go:generate goversioninfo
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/go-logfmt/logfmt"

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
	c, err := canarytools.NewClient(cfg.ImConsoleAPIDomain, cfg.ImConsoleAPIKey, l)
	if err != nil {
		l.Fatal(err)
	}

	fileCount := rand.Intn(cfg.MaxFiles-cfg.MinFiles) + cfg.MinFiles + 1
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
		memo, err := CreateMemo()
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
		err = c.DropFileToken(kind, memo, "", n)
		if err != nil {
			l.Error(err)
			continue
		}
		rtime := GetRandomDate(2)
		os.Chtimes(n, rtime, rtime)
	}
}

func pick(s []string) string {
	return s[rand.Intn(len(s))]
}

// CreateMemo creates a meaningful memo to be included during Canarytoken creation
// value is logfmt encoded for easier processing
func CreateMemo() (memo string, err error) {
	keyVals := []interface{}{
		"Generator", "TokenDropper",
	}

	// Add time
	keyVals = append(keyVals, "Timestamp", time.Now().UTC().Format(time.RFC3339))

	// Add username who run the dropper
	u, err := user.Current()
	if err != nil {
		return
	}
	keyVals = append(keyVals, "Username", u.Username)

	// Get Hostname
	hn, err := os.Hostname()
	if err != nil {
		return
	}
	keyVals = append(keyVals, "Hostname", hn)

	lf, err := logfmt.MarshalKeyvals(keyVals...)
	if err != nil {
		return
	}
	memo = string(lf)
	return
}

func GetRandomTokenName(kind string) (name string, err error) {
	var n string // name
	var e string // ext
	switch kind {
	case "aws-id":
		n = pick(awsFileNames)
		e = "txt"
	case "doc-msword":
		n = pick(fileNames)
		e = "docx"
	case "pdf-acrobat-reader":
		n = pick(fileNames)
		e = "pdf"
	case "msword-macro":
		n = pick(fileNames)
		e = "docm"
	case "msexcel-macro":
		n = pick(fileNames)
		e = "xlsm"
	default:
		err = fmt.Errorf("unsupported Canarytoken: %s", kind)
		return
	}
	name = RandomizeName(n, e)
	return
}

package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"runtime"
	"time"

	"github.com/go-logfmt/logfmt"
)

func GetRandomTokenName(kind string, randomizeFilename bool) (name string, err error) {
	var n string // name
	var e string // ext
	switch kind {
	case "aws-id":
		n = pick(awsFileNames)
		e = "txt"
	case "slack-api":
		n = pick(slackFileNames)
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
	case "windows-dir":
		n = pick(windirFileNames)
		// windows-dir do not have extentions
		e = ""
	default:
		err = fmt.Errorf("unsupported Canarytoken: %s", kind)
		return
	}

	if randomizeFilename {
		name = RandomizeName(n, e)
	} else {
		switch e {
		case "":
			name = n
		default:
			name = n + "." + e
		}
	}

	return
}

// CreateMemo creates a meaningful memo to be included during Canarytoken creation
// value is logfmt encoded for easier processing
func CreateMemo(filename, dropWhere, customMemo string, noDefaultMemo bool) (memo string, err error) {
	if !noDefaultMemo {
		keyVals := []interface{}{}

		// custom reminders?
		if customMemo != "" {
			keyVals = append(keyVals, "Memo", customMemo)
		}

		keyVals = append(keyVals, "Generator", "TokenDropper")
		// Add username who run the dropper
		u, errT := user.Current()
		if errT != nil {
			err = errT
			return
		}
		keyVals = append(keyVals, "TD-User", u.Username)

		// Get Hostname
		hn, errT := os.Hostname()
		if errT != nil {
			err = errT
			return
		}
		keyVals = append(keyVals, "TD-Host", hn)

		// Add original filename
		keyVals = append(keyVals, "OriginalFilename", filename)

		// Add 'where' this token has been dropped
		keyVals = append(keyVals, "Where", dropWhere)

		// Add 'OS' where this token has been dropped
		keyVals = append(keyVals, "TD-OS", runtime.GOOS)

		lf, errT := logfmt.MarshalKeyvals(keyVals...)
		if errT != nil {
			err = errT
			return
		}
		memo = string(lf)
		return
	}
	memo = customMemo
	return
}

// TossCoin randomly returns true or false
func TossCoin() bool {
	if rand.Intn(2) == 1 {
		return true
	}
	return false
}

// RandomizeName takes a name and does a set of premutations to it to make it
// look a bit random
func RandomizeName(name string, ext string) (randName string) {
	randName = name

	// add "Copy of " in front of the name
	if TossCoin() {
		randName = "Copy of " + randName
	}

	// add Random date from last year
	if TossCoin() {
		randName = randName + " " + GetRandomDateString(3)
	}

	// add " ({RandDigit})" to end?
	if TossCoin() {
		randNumber := rand.Intn(5)
		randName = randName + fmt.Sprintf(" (%d)", randNumber+1)
	}

	// if ext is "" don't add a .
	if ext == "" {
		return randName
	}
	return randName + "." + ext
}

// GetRandomDate returns a random date between Now and specicifed number of years
func GetRandomDate(years int) (t time.Time) {
	lastYearUnix := time.Now().Add(-1 * time.Duration(years) * time.Second * 86400 * 365).Unix() // 'years' years ago
	nowUnix := time.Now().Unix()
	delta := nowUnix - lastYearUnix

	randomDateUnix := rand.Int63n(delta) + lastYearUnix

	return time.Unix(randomDateUnix, 0)
}

// GetRandomDateString returns a random date as a string between Now
// and specicifed number of years.
// Time formatting is random as well
func GetRandomDateString(years int) (t string) {
	var timeFormats = []string{
		"January 2 2006",
		"January_2_2006",
		"Jan2006",
		"2006Jan2",
		"2006Jan",
		"Jan",
		"January",
		"2006",
		"2006_01_02",
		"2006-01-02",
		"20060102",
		"02-Jan-2006",
		"010206",
		"Jan-02-06",
		"Jan-02-2006",
	}
	f := timeFormats[rand.Intn(len(timeFormats))]
	return GetRandomDate(years).Format(f)
}

func pick(s []string) string {
	return s[rand.Intn(len(s))]
}

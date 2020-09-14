package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"time"

	"github.com/go-logfmt/logfmt"
)

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

// CreateMemo creates a meaningful memo to be included during Canarytoken creation
// value is logfmt encoded for easier processing
func CreateMemo(filename, dropWhere, customMemo string) (memo string, err error) {
	keyVals := []interface{}{
		"Generator", "TokenDropper",
	}

	// custom reminders?
	if customMemo != "" {
		keyVals = append(keyVals, "Memo", customMemo)
	}
	// Add time
	// keyVals = append(keyVals, "Timestamp", time.Now().UTC().Format(time.RFC3339))

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

	// Add original filename
	keyVals = append(keyVals, "OriginalFilename", filename)

	// Add 'where' this token has been dropped
	keyVals = append(keyVals, "Where", dropWhere)

	lf, err := logfmt.MarshalKeyvals(keyVals...)
	if err != nil {
		return
	}
	memo = string(lf)
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

// 	case "http", "dns", "cloned-web", "doc-msword", "web-image", "windows-dir", "aws-s3", "pdf-acrobat-reader", "msword-macro", "msexcel-macro", "aws-id", "apeeper", "qr-code", "svn", "sql", "fast-redirect", "slow-redirect":
// t, err := c.CreateTokenFromAPI("dns", "koko dns", "", nil)
// if err != nil {
// 	l.Fatal(err)
// }

// n, err := c.DownloadTokenFromAPI(t.Canarytoken.Canarytoken, "hamada.docx")
// if err != nil {
// 	l.Fatal(err)
// }
// l.Infof("written %d bytes", n)

// // get flocks
// flockssummary, err := c.GetFlocksSummary()
// if err != nil {
// 	l.Fatal(err)
// }
// for fid, summary := range flockssummary.FlocksSummary {
// 	l.Infof("%s:%s", fid, summary.Name)
// 	test_flock_id, err := c.GetFlockIDFromName(summary.Name)
// 	if err != nil {
// 		l.Fatal(err)
// 	}
// 	test_flock_name, err := c.GetFlockNameFromID(fid)
// 	if err != nil {
// 		l.Fatal(err)
// 	}
// 	l.Infof("[func] %s:%s", test_flock_name, test_flock_id)
// }
// flock:e5d3b65df5438f1b285692ff3c705571
// flock_id, err :=  c.GetFlockIDFromName("Default Flock")

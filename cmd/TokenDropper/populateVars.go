package main

import (
	"flag"

	canarytools "github.com/thinkst/go-canarytools"
)

func populateVarsFromFlags(cfg *canarytools.TokenDropperConfig) {
	// General flags
	flag.StringVar(&cfg.ConsoleAPIDomain, "domain", "", "Canary console domain (hash)")
	flag.IntVar(&cfg.FilesCount, "count", 1, "Number of Canarytoken files to be generated")
	flag.IntVar(&cfg.RandYearsBack, "yearsback", 3, `Randomize dates of modified files between Now() and 'years' back
we do this so generated tokens better blend in`)
	flag.StringVar(&cfg.DropWhere, "where", "./", "where to drop Canarytokens?")
	flag.BoolVar(&cfg.CreateDirectoryIfNotExists, "createdir", true, "Create the directory where tokens should be dropped if it didn't exist?")
	flag.StringVar(&cfg.KindsStr, "kind", "aws-id,doc-msword", "comma separated list of Canarytokens to be generated")
	// "apeeper":"EC2 Meta-data Service",
	// "autoreg-google-docs":"Google Document",
	// "autoreg-google-sheets":"Google Sheet",
	// "aws-id":"Amazon API Key",
	// "aws-s3":"Amazon S3",
	// "cloned-web":"Cloned Website",
	// "dns":"DNS",
	// "doc-msword":"MS Word .docx Document",
	// "fast-redirect":"Fast HTTP Redirect",
	// "google-docs":"Google Document",
	// "google-sheets":"Google Sheet",
	// "googledocs_factorydoc":"Document Factory",
	// "googlesheets_factorydoc":"Document Factory",
	// "http":"Web",
	// "msexcel-macro":"MS Excel .xlsm Document",
	// "msword-macro":"MS Word .docm Document",
	// "office365mail":"Office 365 email token",
	// "pdf-acrobat-reader":"Acrobat Reader PDF Document",
	// "qr-code":"QR Code",
	// "signed-exe":"Signed Exe",
	// "slack-api":"Slack API Key",
	// "slow-redirect":"Slow HTTP Redirect",
	// "web-image":"Remote Web Image",
	// "windows-dir":"Windows Directory Browsing"
	flag.StringVar(&cfg.CustomMemo, "memo", "", `tokens' memo always include 'host', 'user', and 'filename',
use this flag to add custom text to the Canarytoken memo`)
	flag.StringVar(&cfg.LogLevel, "loglevel", "info", "set loglevel, can be one of ('info', 'warning', 'debug' or 'trace')")

	// Tokens can be created using the console API, or Factory
	// This flag specifies how we're gonna roll
	flag.StringVar(&cfg.OpMode, "opmode", "api", "operate using console API or Factory? valid values are 'api' & 'factory'")

	// Creating tokens using Console API? can't be used with factory
	flag.StringVar(&cfg.ConsoleAPIKey, "apikey", "", "API key (can't be specified with '-factoryauth')")
	flag.StringVar(&cfg.ConsoleTokenFile, "tokenfile", "", "the token file 'canarytools.config' which contains auth token and the domain")

	// Creating tokens using FActory? can't be used with API
	// flag.StringVar(&cfg.ConsoleFactoryAuth, "factoryauth", "", "factory authentication key (can't be specified with '-apikey')")
	// flag.StringVar(&cfg.FactoryAuthFile, "factoryauthfile", "", "the factory auth file 'canaryfactoryauth.config' which contains factory auth and the domain")

	// Flock Specific flags
	flag.StringVar(&cfg.FlockName, "flock", "", "created tokens will be part of this flock 'if empty, will be assigned to the default flock'")
	flag.BoolVar(&cfg.CreateFlockIfNotExists, "createflock", true, "Create the flock if it doesn't exist? has to be used with '-flock', and is not suported with factory")

}

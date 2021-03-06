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
	flag.BoolVar(&cfg.OverwriteFileIfExists, "overwrite-files", false, "Overwrite files if they already exist?")
	flag.StringVar(&cfg.KindsStr, "kind", "aws-id,doc-msword", "comma separated list of Canarytokens to be generated")

	// filename
	flag.StringVar(&cfg.FileName, "filename", "", `filename that will be given to the token; if empty a random name will be set.
setting this will make 'count' to be 'one', and will do some checks to make sure extension matches the 'kind' specified`)
	flag.BoolVar(&cfg.RandomizeFilenames, "randomize-filenames", true, "add random text to filenames to make them unique")
	flag.StringVar(&cfg.CustomMemo, "memo", "", `tokens' memo includes 'host', 'user', and 'filename' by default
use this flag to add custom text to the Canarytoken memo
this flag is mandatory if '-no-default-memo' is set to true`)
	flag.BoolVar(&cfg.NoDefaultMemo, "no-default-memo", false, "do not include the default memo (if this set to true, then you MUST specify '-memo')")
	flag.StringVar(&cfg.LogLevel, "loglevel", "info", "set loglevel, can be one of ('info', 'warning', 'debug' or 'trace')")

	// Tokens can be created using the console API, or Factory
	// This flag specifies how we're gonna roll
	flag.StringVar(&cfg.OpMode, "opmode", "api", "operate using console API or Factory? valid values are 'api' & 'factory'")

	// Creating tokens using Console API? can't be used with factory
	flag.StringVar(&cfg.ConsoleAPIKey, "apikey", "", "API key (can't be specified with '-factoryauth')")
	flag.StringVar(&cfg.ConsoleTokenFile, "tokenfile", "", "the token file 'canarytools.config' which contains auth token and the domain")

	// Creating tokens using FActory? can't be used with API
	flag.StringVar(&cfg.ConsoleFactoryAuth, "factoryauth", "", "factory authentication key (can't be specified with '-apikey')")
	flag.StringVar(&cfg.FactoryAuthFile, "factoryauthfile", "", "the factory auth file 'canaryfactoryauth.config' which contains factory auth and the domain")

	// Flock Specific flags
	flag.StringVar(&cfg.FlockName, "flock", "", "created tokens will be part of this flock (can't be specified with -flockid) 'if empty, will be assigned to the default flock'")
	flag.BoolVar(&cfg.CreateFlockIfNotExists, "createflock", true, "Create the flock if it doesn't exist? has to be used with '-flock', and is not suported with factory")

	flag.StringVar(&cfg.FlockID, "flockid", "", "created tokens will be part of this flock_id (can't be specified with -flock) 'if empty, will be assigned to the default flock'")

}

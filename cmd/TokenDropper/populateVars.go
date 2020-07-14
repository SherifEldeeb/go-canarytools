package main

import (
	"flag"

	"github.com/SherifEldeeb/canarytools"
)

func populateVarsFromFlags(cfg *canarytools.TokenDropperConfig) {
	// General flags
	flag.IntVar(&cfg.MinFiles, "minfiles", 5, "Minimum number of Canarytoken files to be generated")
	flag.IntVar(&cfg.MaxFiles, "maxfiles", 15, "Maximum number of Canarytoken files to be generated")
	flag.IntVar(&cfg.RandYearsBack, "yearsback", 3, "Randomize dates between Now() and 'years' back")
	flag.StringVar(&cfg.DropWhere, "where", "./", "where to drop Canarytokens?")
	flag.StringVar(&cfg.KindsStr, "kind", "aws-id,doc-msword,pdf-acrobat-reader,msword-macro,msexcel-macro", "comma separated list of Canarytokens to be generated")
	flag.StringVar(&cfg.ImConsoleAPIKey, "apikey", "", "API Key")
	flag.StringVar(&cfg.ImConsoleAPIDomain, "domain", "", "canarytools domain")
	flag.StringVar(&cfg.LogLevel, "loglevel", "info", "set loglevel, can be one of ('info', 'warning', 'debug' or 'trace')")

	if cfg.MaxFiles < cfg.MinFiles {
		panic("Min files is Larger than Max files :/")
	}
}

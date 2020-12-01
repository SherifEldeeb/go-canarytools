package canarytools

import log "github.com/sirupsen/logrus"

// TokenDropper is the main struct of the droper
// it contains the configurations, and various components
type TokenDropper struct {
	// configs
	cfg TokenDropperConfig

	// console API client
	C *Client

	// logger
	l *log.Logger
}

// NewTokenDropper returns a TokenDropper from TokenDropperConfig.
// it will check for valid configs, and establish a connection to the console.
func NewTokenDropper(cfg TokenDropperConfig, l *log.Logger) (tdropper TokenDropper, err error) {
	tdropper.cfg = cfg
	c, err := NewClient(cfg.ConsoleAPIConfig, l)
	if err != nil {
		return
	}

	tdropper.C = c
	tdropper.l = l

	return
}

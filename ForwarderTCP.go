package canarytools

import (
	"errors"
	"net"
	"strconv"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// TCPForwarder
type TCPForwarder struct {
	host string
	port int
	l    *log.Logger

	// TODO: TLS!
}

func NewTCPForwarder(host string, port int, loglevel string) (tcpforwarder *TCPForwarder, err error) {
	tcpforwarder = &TCPForwarder{}
	// logging config
	tcpforwarder.l = log.New()
	switch loglevel {
	case "info":
		tcpforwarder.l.SetLevel(log.InfoLevel)
	case "warning":
		tcpforwarder.l.SetLevel(log.WarnLevel)
	case "debug":
		tcpforwarder.l.SetLevel(log.DebugLevel)
	default:
		return nil, errors.New("unsupported log level (can be 'info', 'warning' or 'debug')")
	}
	if host == "" {
		return nil, errors.New("host can't be empty")
	}
	if port > 65535 || port < 1 {
		return nil, errors.New("invalid port")
	}
	tcpforwarder.host = host
	tcpforwarder.port = port
	return
}

func (t TCPForwarder) Forward(outChan <-chan []byte) {
	connect := t.host + ":" + strconv.Itoa(t.port)
	c, err := net.Dial("tcp", connect)
	if err != nil {
		return
	}
	defer c.Close()

	for i := range outChan {
		_, err := c.Write(i)
		if err != nil {
			logrus.Errorf("encoding: %s", string(i))
		}
	}
}

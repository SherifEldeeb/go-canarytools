package canarytools

import (
	"errors"
	"net"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// TCPForwarder
type TCPForwarder struct {
	host string
	port int
	l    *log.Logger
	// TODO: TLS!
}

func NewTCPForwarder(host string, port int, l *log.Logger) (tcpforwarder *TCPForwarder, err error) {
	tcpforwarder = &TCPForwarder{}
	tcpforwarder.l = l

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
			t.l.Errorf("encoding: %s", string(i))
		}
	}
}

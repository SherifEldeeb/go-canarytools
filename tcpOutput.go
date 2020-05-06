package canarytools

import (
	"encoding/json"
	"errors"
	"net"
	"strconv"

	"github.com/sirupsen/logrus"
)

// TCPOutput
type TCPOutput struct {
	host string
	port int
	// TODO: TLS!
}

func NewTCPOutput(host string, port int) (tcpoutput *TCPOutput, err error) {
	tcpoutput = &TCPOutput{}
	if host == "" {
		return nil, errors.New("host can't be empty")
	}
	if port > 65535 || port < 1 {
		return nil, errors.New("invalid port")
	}
	tcpoutput.host = host
	tcpoutput.port = port
	return
}

func (t TCPOutput) Out(outChan <-chan Incident) {
	connect := t.host + ":" + strconv.Itoa(t.port)
	c, err := net.Dial("tcp", connect)
	if err != nil {
		return
	}
	defer c.Close()
	enc := json.NewEncoder(c)
	enc.SetEscapeHTML(false)

	for i := range outChan {
		err := enc.Encode(i)
		if err != nil {
			logrus.Errorf("encoding: %#v", i)
		}
	}
}

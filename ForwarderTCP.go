package canarytools

import (
	"crypto/tls"
	"errors"
	"net"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// TCPForwarder
type TCPForwarder struct {
	host      string
	port      int
	l         *log.Logger
	tlsConfig *tls.Config
	sslUseSSL bool
	// TODO: TLS!
}

func NewTCPForwarder(host string, port int, tlsConfig *tls.Config, sslUseSSL bool, l *log.Logger) (tcpforwarder *TCPForwarder, err error) {
	tcpforwarder = &TCPForwarder{}
	tcpforwarder.l = l
	tcpforwarder.sslUseSSL = sslUseSSL
	tcpforwarder.tlsConfig = tlsConfig

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

func (t TCPForwarder) Forward(outChan <-chan []byte, incidentAckerChan chan<- []byte) {
	var conn net.Conn
	var err error

	connect := t.host + ":" + strconv.Itoa(t.port)
	// using SSL/TLS?
	if t.sslUseSSL {
		conn, err = tls.Dial("tcp", connect, t.tlsConfig)
		if err != nil {
			t.l.WithFields(log.Fields{
				"err":  err,
				"host": connect,
			}).Fatalf("error dialing TLS")
		}
	} else {
		conn, err = net.Dial("tcp", connect)
		if err != nil {
			t.l.WithFields(log.Fields{
				"source": "TCPForwarder",
				"stage":  "forward",
				"err":    err,
			}).Fatal("Forward error dialing")
		}
	}
	defer conn.Close()

	for v := range outChan {
		_, err := conn.Write(v)
		if err != nil {
			t.l.WithFields(log.Fields{
				"source": "TCPForwarder",
				"stage":  "forward",
				"err":    err,
			}).Fatal("Forward error writing to socket")
			return
		}
		incidentAckerChan <- v

	}
}

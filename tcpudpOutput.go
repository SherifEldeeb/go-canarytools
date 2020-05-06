package canarytools

// TCPOutput
type TCPOutput struct {
	host string
	port int
	// TODO: TLS!
}

func NewTCPOutput(host string, port int) (tcpoutput TCPOutput, err error) {

}

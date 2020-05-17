package main

import (
	"fmt"

	"github.com/SherifEldeeb/canarytools"
)

const (
	sslCA   = "../ChirpForwarder/certs/ca/ca.crt"
	sslCert = "../ChirpForwarder/certs/chirp/chirp.crt"
	sslKey  = "../ChirpForwarder/certs/chirp/chirp.key"
)

// var (
// 	insecure = flag.Bool("insecure-ssl", false, "Accept/Ignore all server SSL certificates")
// )

func main() {
	a, d, err := canarytools.LoadTokenFile("../ChirpForwarder/canarytools.config")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a, d)
}

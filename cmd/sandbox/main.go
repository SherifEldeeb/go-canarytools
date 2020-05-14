package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
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
	// flag.Parse()
	var insecure = false

	// Load client cert
	clientCert, err := tls.LoadX509KeyPair(sslCert, sslKey)
	if err != nil {
		log.Fatal(err)
	}
	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	// Read in the cert file
	certs, err := ioutil.ReadFile(sslCA)
	if err != nil {
		log.Fatalf("Failed to append %q to RootCAs: %v", sslCA, err)
	}
	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs appended, using system certs only")
	}

	// Trust the augmented cert pool in our client
	tlsConfig := &tls.Config{
		InsecureSkipVerify: insecure,
		RootCAs:            rootCAs,
		Certificates:       []tls.Certificate{clientCert},
	}

	httpTransport := &http.Transport{
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second,
		TLSClientConfig:       tlsConfig,
	}

	elasticConfig := elasticsearch.Config{
		Addresses: []string{"https://192.168.0.165:9200"}, // A list of Elasticsearch nodes to use.
		Username:  "elastic",                              // Username for HTTP Basic Authentication.
		Password:  "elastic",                              // Password for HTTP Basic Authentication.
		Transport: httpTransport,
	}
	c, err := elasticsearch.NewClient(elasticConfig)
	if err != nil {
		log.Fatal(err)
	}
	i, err := c.Info()
	if err != nil {
		log.Fatal(err)
	}
	if i.IsError() {
		log.Fatal(i)
	}
	fmt.Println(i.String())
	i.Body.Close()
}

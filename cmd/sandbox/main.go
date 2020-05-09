package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
)

func main() {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://127.0.0.1:9200"}, // A list of Elasticsearch nodes to use.
		Username:  "elastic",                         // Username for HTTP Basic Authentication.
		Password:  "elastic",                         // Password for HTTP Basic Authentication.
		// elastic cloud specific
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	c, err := elasticsearch.NewClient(cfg)
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
	i.Body.Close()
	fmt.Println(i)
}

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/koderhut/memorynotes/config"
)

var c config.Context

func init() {
	addr := flag.String("addr", "0.0.0.0:44666", "IP and port to bind to")
	prefix := flag.String("path-prefix", "/api", "Path prefix for endpoints")
	domain := flag.String("domain", "localhost", "Domain used for link generation")
	flag.Parse()
	c.IP, c.Port, _ = c.ParseAddr(*addr)
	c.PathPrefix = *prefix
	c.Domain = c.ParseTLD(*domain)
}

func main() {
	router := ConfigRouter(&c)

	log.Printf(">>> safe-notes web service is ready to receive requests on: [%s]\n", c.Addr())

	log.Fatal(http.ListenAndServe(c.Addr(), router))
}

package main

import (
	"context"
	"flag"
	"github.com/koderhut/memorynotes/config"
	"github.com/koderhut/memorynotes/webapp"
	"log"
	"os"
	"os/signal"
	"time"
)

var c config.Parameters

func init() {
	addr := flag.String("addr", "0.0.0.0:44666", "IP and port to bind to")
	prefix := flag.String("path-prefix", "/api", "Path prefix for endpoints")
	domain := flag.String("domain", "localhost", "Domain used for link generation")
	cors := flag.String("cors", "http://localhost:44666", "Allowed hosts for access")

	flag.Parse()

	var err error
	c, err = config.NewConfigParams(*addr, *prefix, *domain, *cors)

	if nil != err {
		panic("Unable to process config parameters. Shutting down!")
	}
}

func main() {
	wait := time.Second * 15
	addr, _ := c.Addr()

	// bootstrap the server and register the routes
	srv := webapp.BootstrapServer(c)

	log.Printf(">>> memory-notes web service is ready to receive requests on: [%s]\n", addr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println(">>> memory-notes web service has shutdown")

	os.Exit(0)
}

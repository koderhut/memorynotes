package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

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
	var wait time.Duration
	srv := ConfigRouter(&c)

	log.Printf(">>> safe-notes web service is ready to receive requests on: [%s]\n", c.Addr())

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

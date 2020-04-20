package config

import (
	"fmt"
	"net"
	"strings"
)

type Context struct {
	IP         string
	Port       string
	PathPrefix string
	Domain     string
}

// ParseAddr parses and makes sure the IP:port provided are correct
func (c *Context) ParseAddr(a string) (string, string, error) {
	host, port, err := net.SplitHostPort(a)

	if err != nil {
		panic("Incorrect IP:Port values provided")
	}

	if host == "" {
		host = "0.0.0.0"
	}

	return host, port, nil
}

// Addr returns the address in IP:port
func (c *Context) Addr() string {
	return fmt.Sprintf("%s:%s", c.IP, c.Port)
}

// ParseTLD
func (c *Context) ParseTLD(d string) string {
	var port string

	if false == strings.EqualFold(c.Port, "80") && false == strings.EqualFold(c.Port, "443") {
		port = fmt.Sprintf(":%s", c.Port)
	}

	return fmt.Sprintf("%s%s", d, port)
}

package main

import (
	"net"
	"net/url"
	"strings"
	"time"
)

func IsConnect(address string, timeout int) bool {
	conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func Address(u *url.URL) string {
	host := u.Host
	if strings.Index(host,":") == -1 {
		host += ":80"
	}
	return host
}

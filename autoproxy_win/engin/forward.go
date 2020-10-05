package engin

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type Forward interface {
	Close() error
	http(r *http.Request) (*http.Response, error)
	https(address string, r *http.Request) (net.Conn, error)
}

type defaultForward struct {
	sync.RWMutex

	tmout int
	address map[string]int
	trans *http.Transport
}

func (d *defaultForward)Close() error {
	d.trans.CloseIdleConnections()
	return nil
}

func (d *defaultForward)http(r *http.Request) (*http.Response, error) {
	return d.trans.RoundTrip(r)
}

func (d *defaultForward)https(address string, r *http.Request) (net.Conn, error) {
	return net.DialTimeout("tcp", address, time.Second * time.Duration(d.tmout) )
}

func NewDefault(timeout int) (Forward, error) {
	return &defaultForward{
		trans: newTransport(timeout, nil),
		tmout: timeout,
		address: make(map[string]int, 0),
	},nil
}
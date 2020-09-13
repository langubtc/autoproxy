package main

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http/httpproxy"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type proxyServer struct {
	sync.RWMutex

	enable  bool
	auth    AuthConfig
	address string
	timeout int

	tlscfg *tls.Config
	client map[string]*http.Transport
	proxyconfig *httpproxy.Config
	proxyfunc func(reqURL *url.URL) (*url.URL, error)
}

type HttpProxyServer struct {
	server *http.Server
	listen  string
	timeout int
	tlscfg *tls.Config

	mode string // auto,proxy,local
	auth []AuthConfig

	local *http.Transport
	proxy []*proxyServer
}

func (proxy *proxyServer)ProxyFunc(r *http.Request) (*url.URL, error)  {
	return proxy.proxyfunc(r.URL)
}

func NewTransport(timeout int,tlscfg *tls.Config) *http.Transport {
	tmout := time.Duration(timeout) * time.Second
	return &http.Transport{
		TLSClientConfig: tlscfg,
		DialContext: (&net.Dialer{
			Timeout:   tmout,
			KeepAlive: tmout,
		}).DialContext,
		MaxIdleConns:          1000,
		IdleConnTimeout:       3*tmout,
		TLSHandshakeTimeout:   tmout/3,
		ExpectContinueTimeout: time.Second }
}

func (proxy *proxyServer)RoundTrip(req *http.Request) (*http.Response, error) {
	proxy.RLock()
	transport := proxy.client[req.Host]
	proxy.RUnlock()

	if transport != nil {
		return transport.RoundTrip(req)
	}

	proxy.Lock()
	transport = NewTransport(proxy.timeout, proxy.tlscfg)
	transport.Proxy = proxy.ProxyFunc
	proxy.client[req.Host] = transport
	proxy.Unlock()

	return transport.RoundTrip(req)
}

func newProxyServer(cfg RemoteConfig) *proxyServer {
	var err error
	proxy := new(proxyServer)

	proxy.tlscfg, err = TlsConfigClient(cfg.Tls, cfg.Address)
	if err != nil {
		Fatal(err.Error())
	}

	proxy.client = make(map[string]*http.Transport,1024)
	proxy.address = cfg.Address
	proxy.auth = cfg.Auth
	proxy.timeout = cfg.Timeout

	scheme := "http"
	if proxy.tlscfg != nil {
		scheme = "https"
	}

	var secondProxy string
	if proxy.auth.UserName != "" && proxy.auth.UserName != "" {
		secondProxy = fmt.Sprintf("%s://%s:%s@%s",
			scheme, proxy.auth.UserName, proxy.auth.Password, proxy.address)
	}else {
		secondProxy = fmt.Sprintf("%s://%s",
			scheme, proxy.address)
	}
	proxy.proxyconfig = &httpproxy.Config{HTTPProxy:secondProxy, HTTPSProxy:secondProxy}
	proxy.proxyfunc = proxy.proxyconfig.ProxyFunc()

	return proxy
}

func NewHttpProxyServer(cfg *Config) *HttpProxyServer{
	var err error

	proxy := new(HttpProxyServer)
	proxy.auth = cfg.Local.Auths
	proxy.listen = cfg.Local.Listen
	proxy.timeout = cfg.Local.Timeout
	proxy.mode = cfg.Local.Mode

	proxy.tlscfg,err = TlsConfigServer(cfg.Local.Tls)
	if err != nil {
		Fatal(err.Error())
	}

	if cfg.Local.Mode == "local" || cfg.Local.Mode == "auto" {
		proxy.local = NewTransport(cfg.Local.Timeout,nil)
	}

	if cfg.Local.Mode == "proxy" || cfg.Local.Mode == "auto" {
		proxy.proxy = make([]*proxyServer,0)
		for _,v := range cfg.Remote {
			proxyone := newProxyServer(v)
			proxy.proxy = append(proxy.proxy, proxyone)
		}
	}

	return proxy
}

func (proxy *HttpProxyServer)Start() error {

	lis, err := net.Listen("tcp",proxy.listen)
	if err != nil {
		Fatal(err.Error())
		return err
	}

	if proxy.tlscfg != nil {
		lis = tls.NewListener(lis, proxy.tlscfg)
	}

	tmout := time.Duration(proxy.timeout) * time.Second

	proxy.server = &http.Server{
		Handler:proxy,
		ReadTimeout:tmout,
		WriteTimeout:tmout,
		TLSConfig:proxy.tlscfg,
	}

	proxy.server.SetKeepAlivesEnabled(true)

	if proxy.tlscfg != nil {
		Infof("Proxy [https://%s] success!\n", proxy.listen)
	}else {
		Infof("Proxy [http://%s] success!\n", proxy.listen)
	}

	Infof("Run mode %s", proxy.mode)

	err = proxy.server.Serve(lis)
	if err != nil {
		Fatal(err.Error())
	}
	return err
}
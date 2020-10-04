package engin

import (
	"crypto/tls"
	"fmt"
	"github.com/astaxie/beego/logs"
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
	proxy *proxyServer

	request int64
}

func (proxy *proxyServer)ProxyFunc(r *http.Request) (*url.URL, error)  {
	return proxy.proxyfunc(r.URL)
}

func NewTransport(timeout int, tlscfg *tls.Config) *http.Transport {
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

func newProxyServer(remote *RemoteConfig) *proxyServer {
	var err error
	proxy := new(proxyServer)

	if remote.TlsEnable {
		proxy.tlscfg, err = TlsConfigClient(remote.Address)
		if err != nil {
			logs.Error(err.Error())
			return nil
		}
	}

	proxy.client = make(map[string]*http.Transport, 1024)
	proxy.address = remote.Address
	proxy.auth = remote.Auth
	proxy.timeout = remote.Timeout

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

type LocalConfig struct {
	Listen     string
	Timeout    int
	Mode       string  // local、auto、proxy
	Auths      []AuthConfig
	TlsEnable  bool
	TlsVersion string
}

type RemoteConfig struct {
	Address   string
	Timeout   int
	Auth      AuthConfig
	TlsEnable bool
}

func NewHttpProxyServer(local *LocalConfig, remote *RemoteConfig) *HttpProxyServer{
	var err error

	proxy := new(HttpProxyServer)
	proxy.auth = local.Auths
	proxy.listen = local.Listen
	proxy.timeout = local.Timeout
	proxy.mode = local.Mode

	if local.TlsEnable {
		proxy.tlscfg, err = TlsConfigServer(nil)
		if err != nil {
			logs.Error(err.Error())
			return nil
		}
	}

	if local.Mode == "local" || local.Mode == "auto" {
		proxy.local = NewTransport(local.Timeout,nil)
	}

	if local.Mode == "proxy" || local.Mode == "auto" {
		proxy.proxy = newProxyServer(remote)
	}

	return proxy
}

func (proxy *HttpProxyServer)Start() error {
	lis, err := net.Listen("tcp", proxy.listen)
	if err != nil {
		logs.Error(err.Error())
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
		logs.Info("Proxy [https://%s] success!\n", proxy.listen)
	}else {
		logs.Info("Proxy [http://%s] success!\n", proxy.listen)
	}

	logs.Info("Run mode %s", proxy.mode)

	err = proxy.server.Serve(lis)
	if err != nil {
		logs.Error(err.Error())
	}
	return err
}

func (proxy *HttpProxyServer)Shutdown()  {
	proxy.server.Close()
	proxy.local.CloseIdleConnections()
}

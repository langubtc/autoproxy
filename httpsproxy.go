package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/lixiangyun/autoproxy/util"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

var httpconnects int32

func getHeaderValue(values []string) string {
	body := bytes.NewBuffer(make([]byte,0))
	for i, v := range values {
		if i  == len(values) - 1 {
			body.WriteString(v)
		}else {
			body.WriteString(v + " ")
		}
	}
	return body.String()
}

func getSecondRequest(r *http.Request) []byte {
	body := bytes.NewBuffer(make([]byte,0))
	body.WriteString(fmt.Sprintf("CONNECT %s %s\r\n",r.URL.Host,r.Proto))
	body.WriteString(fmt.Sprintf("Host: %s\r\n", r.Host))
	for key,value := range r.Header {
		body.WriteString(fmt.Sprintf("%s: %s\r\n",key,getHeaderValue(value)))
	}
	body.WriteString("\r\n")
	return body.Bytes()
}

func HttpConnectsAdd() {
	atomic.AddInt32(&httpconnects,1)
	LogPrefix(fmt.Sprintf("%d",httpconnects))
}

func HttpConnectsDel() {
	atomic.AddInt32(&httpconnects,-1)
	LogPrefix(fmt.Sprintf("%d",httpconnects))
}

func HttpsProxyOneHandler(r *http.Request, proxy *proxyServer) (net.Conn,error) {
	host := util.Address(r.URL)

	Infof("connect to %s use remote proxy %s", host, proxy.address)

	server, err := net.DialTimeout("tcp", proxy.address,
		time.Second*time.Duration(proxy.timeout) )
	if err != nil {
		return nil, fmt.Errorf("connect to remote proxy %s failed, err=%s",
			proxy.address, err.Error())
	}
	defer func() {
		if err != nil {
			server.Close()
		}
	}()

	r.Header.Del("Proxy-Authenticate")
	proxyAddAuth(r, &proxy.auth)

	body := getSecondRequest(r)

	if proxy.tlscfg != nil {
		server = tls.Client(server, proxy.tlscfg)
	}

	cnt, err := server.Write(body)
	if err != nil {
		return nil, fmt.Errorf("write to remote proxy failed! %s",err.Error())
	}

	var readbuf [1024]byte
	cnt, err = server.Read(readbuf[:])
	if err != nil {
		return nil, fmt.Errorf("read from remote proxy failed! %s",err.Error())
	}

	if -1 == strings.Index(string(readbuf[:cnt]),"200") {
		Warnf("read from remote proxy body :%s ",string(readbuf[:cnt]))
	}

	return server, nil
}

func (proxy *HttpProxyServer)HttpsProxyHandler(r *http.Request) (net.Conn,error)  {
	var server net.Conn
	var err error

	for _,v := range proxy.proxy {
		server, err = HttpsProxyOneHandler(r, v)
		if err != nil {
			continue
		}else {
			break
		}
	}
	return server, err
}

func (proxy *HttpProxyServer)HttpsHandler(w http.ResponseWriter, r *http.Request) {
	hij, ok := w.(http.Hijacker)
	if !ok {
		Fatal("httpserver does not support hijacking")
	}
	client, _, e := hij.Hijack()
	if e != nil {
		Fatalf("Cannot hijack connection " + e.Error())
	}

	var server net.Conn
	var err error

	tmout := time.Second * time.Duration(proxy.timeout)

	host := util.Address(r.URL)
	if proxy.mode == "local" {
		server, err = net.DialTimeout("tcp", host, tmout )
	}else if proxy.mode == "proxy" {
		server, err = proxy.HttpsProxyHandler(r)
	}else {
		host := util.Address(r.URL)
		if IsSecondProxy(host) {
			server, err = proxy.HttpsProxyHandler(r)
		}else {
			server, err = net.DialTimeout("tcp", host, tmout)
		}
	}

	if err != nil {
		Errorf("address %s: %s", host, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		client.Close()
		return
	}

	client.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	go func() {
		HttpConnectsAdd()

		wt := NewWaitGroupTimeout(2, tmout, func() {
			Infof("connect from %s to %s timeout",
				server.RemoteAddr(), client.RemoteAddr())

			client.Close()
			server.Close()

		})

		Infof("connect from %s to %s start",
			server.RemoteAddr(), client.RemoteAddr())

		go ConnectCopy(client, server, wt)
		go ConnectCopy(server, client, wt)

		wt.Wait()
		Infof("connect from %s to %s close",
			server.RemoteAddr(), client.RemoteAddr())

		client.Close()
		server.Close()

		HttpConnectsDel()
	}()
}

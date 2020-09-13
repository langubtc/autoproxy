package main

import (
	"encoding/base64"
	"fmt"
	"github.com/lixiangyun/autoproxy/util"
	"io"
	"net/http"
	"strings"
	"time"
)

func proxyAddAuth(r *http.Request, auth *AuthConfig)  {
	if auth.UserName != "" && auth.Password != "" {
		authBody := auth.UserName + ":" + auth.Password
		basic := "Basic " + base64.StdEncoding.EncodeToString([]byte(authBody))
		r.Header.Add("Proxy-Authorization",basic)
		Infof("add auth %s %s",authBody,basic)
	}
}

func proxyAuthInfo(r *http.Request, auth []AuthConfig) bool {
	if auth == nil || len(auth) == 0{
		return true
	}
	value := r.Header.Get("Proxy-Authorization")
	if value == "" {
		Warnf("[%s]no auth form header",r.RemoteAddr)
		return false
	}

	body, err := base64.StdEncoding.DecodeString(value[6:])
	if err != nil {
		Warnf("[%s:%s]auth is illegal",r.RemoteAddr,value)
		return false
	}
	ctx := strings.Split(string(body),":")
	if len(ctx) != 2 {
		Warnf("[%s:%s]auth is illegal",r.RemoteAddr,body)
		return false
	}
	for _,v := range auth {
		if v.UserName == ctx[0] && v.Password == ctx[1] {
			return true
		}
	}
	Warnf("[%s:%s:%s]auth is not exist",
		r.RemoteAddr,ctx[0],ctx[1])
	return false
}

func NoProxyHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second) // 防DOS攻击延时
	Warnf("Request is illegal. RemoteAddr: %s",r.RemoteAddr)
	http.Error(w,
		"This is a proxy server. Does not respond to non-proxy requests.",
		http.StatusInternalServerError)
}

func AuthFailHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(30*time.Second) // 防DOS攻击延时
	Warnf("Request authentication failed. RemoteAddr: %s",r.RemoteAddr)
	http.Error(w,
		"Request authentication failed.",
		http.StatusUnauthorized)
}

func copyHeaders(dst, src http.Header) {
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}

func removeProxyHeaders(r *http.Request)  {
	r.RequestURI = ""
	r.Header.Del("Proxy-Connection")
	r.Header.Del("Proxy-Authenticate")
	r.Header.Del("Proxy-Authorization")
}

func (proxy *HttpProxyServer)proxyRoundTrip(r *http.Request) (*http.Response, error) {
	var rsp *http.Response
	var err error

	for _,v := range proxy.proxy {
		r.Header.Del("Proxy-Authenticate")
		proxyAddAuth(r, &v.auth)

		rsp, err = v.RoundTrip(r)
		if err != nil {
			continue
		}else {
			break
		}
	}

	return rsp, err
}

func DebugReqeust(r *http.Request) {
	var headers string
	for key, value := range r.Header {
		headers += fmt.Sprintf("[%s:%s]",key,value)
	}
	Infof("%s %s %s %s",r.RemoteAddr, r.Method, r.URL.String(), headers)
}

func (proxy *HttpProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	DebugReqeust(r)

	if !proxyAuthInfo(r, proxy.auth) {
		AuthFailHandler(w,r)
		return
	}

	Infof("recv request from %s", r.RemoteAddr)

	if r.Method == "CONNECT" {
		proxy.HttpsHandler(w, r)
		return
	}

	if !r.URL.IsAbs() {
		NoProxyHandler(w, r)
		return
	}

	removeProxyHeaders(r)

	var rsp *http.Response
	var err error

	Infof("transport %s %s start", r.Host, r.URL.String())

	if proxy.mode == "local" {
		rsp, err = proxy.local.RoundTrip(r)
	}else if proxy.mode == "proxy" {
		rsp, err = proxy.proxyRoundTrip(r)
	}else {
		host := util.Address(r.URL)
		if IsSecondProxy(host) {
			rsp, err = proxy.local.RoundTrip(r)
		}else {
			rsp, err = proxy.proxyRoundTrip(r)
		}
	}

	if err != nil {
		Warnf("transport %s %s failed! %s",
			r.Host, r.URL.String(), err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	if rsp == nil {
		errStr := fmt.Sprintf("transport %s read response failed!", r.URL.Host)
		Warnf(errStr)
		http.Error(w, errStr, 500)
		return
	}

	origBody := rsp.Body
	defer origBody.Close()

	copyHeaders(w.Header(), rsp.Header)
	w.WriteHeader(rsp.StatusCode)
	io.Copy(w, rsp.Body)

	Infof("transport %s %s success", r.Host, r.URL.String())
}

package engin

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

func PublicFailDelay() {
	time.Sleep(time.Second * 5) // 防DOS攻击延时
}

type HttpAccess struct {
	Timeout int
	Address string
}

type Access interface {
}

func NoProxyHandler(w http.ResponseWriter, r *http.Request) {
	PublicFailDelay()
	logs.Warn("request is illegal. RemoteAddr: ", r.RemoteAddr)
	http.Error(w,
		"This is a proxy server. Does not respond to non-proxy requests.",
		http.StatusInternalServerError)
}

func AuthFailHandler(w http.ResponseWriter, r *http.Request)  {
	PublicFailDelay()
	logs.Warn("Request authentication failed. RemoteAddr: ", r.RemoteAddr)
	http.Error(w, "Request authentication failed.", http.StatusUnauthorized)
}

func AuthInfoParse(r *http.Request) *AuthInfo {
	value := r.Header.Get("Proxy-Authorization")
	if value == "" {
		return nil
	}
	body, err := base64.StdEncoding.DecodeString(value[6:])
	if err != nil {
		return nil
	}
	ctx := strings.Split(string(body),":")
	if len(ctx) != 2 {
		return nil
	}
	return &AuthInfo{User: ctx[0], Token: ctx[1]}
}

func AuthHttp(r *http.Request) bool {
	return true
}

func (acc *HttpAccess)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if AuthHttp(r) == false {
		AuthFailHandler(w, r)
		return
	}

	if r.Method == "CONNECT" {
		HttpsRoundTripper(w, r)
		return
	}

	if !r.URL.IsAbs() {
		NoProxyHandler(w, r)
		return
	}

	removeProxyHeaders(r)

	rsp, err := HttpRoundTripper(r)
	if err != nil {
		errStr := fmt.Sprintf("transport %s %s failed! %s", r.Host, r.URL.String(), err.Error())
		logs.Warn(errStr)
		http.Error(w, errStr, http.StatusInternalServerError)
		return
	}

	if rsp == nil {
		errStr := fmt.Sprintf("transport %s read response failed!", r.URL.Host)
		logs.Warn(errStr)
		http.Error(w, errStr, http.StatusInternalServerError)
		return
	}

	origBody := rsp.Body
	defer origBody.Close()

	copyHeaders(w.Header(), rsp.Header)
	w.WriteHeader(rsp.StatusCode)

	_, err = io.Copy(w, rsp.Body)
	if err != nil {
		logs.Warn("io copy fail", err.Error())
	}
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

func NewHttpsAccess(addr string, timeout int, config *tls.Config) Access {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logs.Error("listen address fail", addr)
		panic(err.Error())
	}

	if config != nil {
		lis = tls.NewListener(lis, config)
	}

	acc := new(HttpAccess)
	acc.Address = addr
	acc.Timeout = timeout

	tmout := time.Duration(timeout) * time.Second
	httpserver := &http.Server{
		Handler:acc,
		ReadTimeout:tmout,
		WriteTimeout:tmout,
		TLSConfig: config,
	}

	go func() {
		err = httpserver.Serve(lis)
		if err != nil {
			logs.Error("http server fail", err.Error())
			panic(err.Error())
		}
	}()

	if config == nil {
		logs.Info("access http start success.")
	}else {
		logs.Info("access https start success.")
	}

	return acc
}

func NewHttpAccess(addr string, timeout int) Access {
	return NewHttpsAccess(addr, timeout, nil)
}

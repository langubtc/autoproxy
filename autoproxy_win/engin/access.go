package engin

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func PublicFailDelay() {
	time.Sleep(time.Second) // 防DOS攻击延时
}

type HttpAccess struct {
	Timeout int
	Address string
	httpserver *http.Server
	sync.WaitGroup

	requset  uint64
	lastreq  uint64
	flowsize uint64
	lastflow uint64

	authHandler func(auth *AuthInfo) bool
	forwardHandler func(address string, r *http.Request) Forward
}

type Access interface {
	Stat() (uint64,uint64)
	Shutdown() error
	AuthHandlerSet(func(*AuthInfo) bool)
	ForwardHandlerSet(func(address string, r *http.Request) Forward)
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
	http.Error(w,
		"Request authentication failed.",
		http.StatusUnauthorized)
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

func (acc *HttpAccess)AuthHandlerSet(handler func(auth *AuthInfo) bool)  {
	acc.authHandler = handler
}

func (acc *HttpAccess)ForwardHandlerSet(handler func(address string, r *http.Request) Forward )  {
	acc.forwardHandler = handler
}

func (acc *HttpAccess)AuthHttp(r *http.Request) bool {
	if acc.authHandler == nil {
		return true
	}
	return acc.authHandler(AuthInfoParse(r))
}

func (acc *HttpAccess)Stat() (uint64,uint64) {
	tempreq := acc.requset
	tempflow := acc.flowsize

	req := tempreq - acc.lastreq
	flow := tempflow - acc.lastflow

	acc.lastreq = tempreq
	acc.lastflow = tempflow

	return req, flow
}

func (acc *HttpAccess)Shutdown() error {
	context, cencel := context.WithTimeout(context.Background(), 15 * time.Second)
	err := acc.httpserver.Shutdown(context)
	cencel()
	if err != nil {
		logs.Error("http access ready to shut down fail, %s", err.Error())
	}
	acc.Wait()
	return err
}

func DebugReqeust(r *http.Request) {
	var headers string
	for key, value := range r.Header {
		headers += fmt.Sprintf("[%s:%s]",key,value)
	}
	logs.Info("%s %s %s %s",r.RemoteAddr, r.Method, r.URL.String(), headers)
}

func (acc *HttpAccess)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	DebugReqeust(r)

	atomic.AddUint64(&acc.requset, 1)

	if acc.AuthHttp(r) == false {
		AuthFailHandler(w, r)
		return
	}

	if r.Method == "CONNECT" {
		acc.HttpsRoundTripper(w, r)
		return
	}

	if !r.URL.IsAbs() {
		NoProxyHandler(w, r)
		return
	}

	removeProxyHeaders(r)

	rsp, err := acc.HttpRoundTripper(r)
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

	cnt, err := io.Copy(w, rsp.Body)
	if err != nil {
		logs.Warn("io copy fail", err.Error())
	}

	atomic.AddUint64(&acc.flowsize, uint64(cnt))
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

func NewHttpsAccess(addr string, timeout int, tlsEnable bool) (Access, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logs.Error("listen address fail", addr)
		return nil, err
	}

	var config *tls.Config
	if tlsEnable {
		config, err = TlsConfigServer(nil)
		if err != nil {
			logs.Error("make tls config server fail, %s", err.Error())
			return nil, err
		}
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

	acc.httpserver = httpserver

	acc.Add(1)

	go func() {
		defer acc.Done()
		err = httpserver.Serve(lis)
		if err != nil {
			logs.Error("http server ", err.Error())
		}
	}()

	if config == nil {
		logs.Info("access http start success.")
	}else {
		logs.Info("access https start success.")
	}

	return acc,nil
}



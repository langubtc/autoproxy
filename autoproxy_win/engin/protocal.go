package engin

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net"
	"net/http"
)

var forwardDefault Forward

func HttpForward(r *http.Request) (*http.Response, error) {
	return forwardDefault.http(r)
}

func HttpsForward(address string, r *http.Request) (net.Conn, error) {
	return forwardDefault.https(address, r)
}

var HTTPS_CLIENT_CONNECT_FLAG  = []byte("HTTP/1.1 200 Connection Established\r\n\r\n")

func HttpsRoundTripper(w http.ResponseWriter, r *http.Request) {
	hij, ok := w.(http.Hijacker)
	if !ok {
		logs.Error("httpserver does not support hijacking")
	}

	client, _, err := hij.Hijack()
	if err != nil {
		logs.Error("Cannot hijack connection", err.Error())
		panic("golang sdk is too old.")
	}

	address := Address(r.URL)

	err = WriteFull(client, HTTPS_CLIENT_CONNECT_FLAG)
	if err != nil {
		errstr := fmt.Sprintf("client connect %s fail", client.RemoteAddr())
		logs.Error(errstr, err.Error())
		http.Error(w, errstr, http.StatusInternalServerError)

		client.Close()
		return
	}

	server, err := HttpsForward(address, r)
	if err != nil {
		errstr := fmt.Sprintf("can't forward hostname %s", address)
		logs.Error(errstr, err.Error())
		http.Error(w, errstr, http.StatusInternalServerError)

		client.Close()
		return
	}

	go func() {
		ConnectCopyWithTimeout(client, server, 60)
	}()
}

func HttpRoundTripper(r *http.Request) (*http.Response, error) {
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(r.Body)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return HttpForward(r)
}


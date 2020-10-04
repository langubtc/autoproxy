package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	mathrand "math/rand"
)

func VersionGet() string {
	return "v1.1.0"
}

func SaveToFile(name string, body []byte) error {
	return ioutil.WriteFile(name, body, 0664)
}

func GetToken(length int) string {
	token := make([]byte, length)
	bytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!#$%^&*"
	for i:=0; i<length; i++  {
		token[i] = bytes[mathrand.Int()%len(bytes)]
	}
	return string(token)
}

func GetUser(length int) string {
	token := make([]byte, length)
	bytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	for i:=0; i<length; i++  {
		token[i] = bytes[mathrand.Int()%len(bytes)]
	}
	return string(token)
}

func CapSignal(proc func())  {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<- signalChan
		proc()
		os.Exit(-1)
	}()
}

func InterfaceAddsGet(iface *net.Interface) ([]net.IP, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, nil
	}
	ips := make([]net.IP, 0)
	for _, v:= range addrs {
		ipone, _, err:= net.ParseCIDR(v.String())
		if err != nil {
			continue
		}
		if len(ipone) > 0 {
			ips = append(ips, ipone)
		}
	}
	return ips, nil
}

func InterfaceLocalIP(inface *net.Interface) ([]net.IP, error) {
	addrs, err := InterfaceAddsGet(inface)
	if err != nil {
		return nil, err
	}
	var output []net.IP
	for _, v := range addrs {
		if IsIPv4(v) == true {
			output = append(output, v)
		}
	}
	if len(output) == 0 {
		return nil, fmt.Errorf("interface not ipv4 address.")
	}
	return output, nil
}

func IsIPv4(ip net.IP) bool {
	return strings.Index(ip.String(), ".") != -1
}

func ByteViewLite(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%db", size)
	} else if size < (1024 * 1024) {
		return fmt.Sprintf("%.1fkb", float64(size)/float64(1024))
	} else {
		return fmt.Sprintf("%.1fmb", float64(size)/float64(1024*1024))
	}
}

func ByteView(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	} else if size < (1024 * 1024) {
		return fmt.Sprintf("%.1fKB", float64(size)/float64(1024))
	} else if size < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.1fMB", float64(size)/float64(1024*1024))
	} else if size < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.1fGB", float64(size)/float64(1024*1024*1024))
	} else {
		return fmt.Sprintf("%.1fTB", float64(size)/float64(1024*1024*1024*1024))
	}
}

func init()  {
	mathrand.Seed(time.Now().Unix())
}
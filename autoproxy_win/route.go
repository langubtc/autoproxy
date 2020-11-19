package main

import (
	"github.com/astaxie/beego/logs"
	"strings"
	"sync"
)

type RouteCtrl struct {
	sync.RWMutex
	cache map[string]string
	domain []string
}

var routeCtrl *RouteCtrl

func RouteInit() error {
	routeCtrl = new(RouteCtrl)
	routeCtrl.cache = make(map[string]string, 2048)
	routeCtrl.domain = StringClone(DomainList())
	return nil
}

func AddressToDomain(address string) string {
	domain := address
	idx := strings.Index(address, ":")
	if idx != -1 {
		domain = address[:idx]
	}
	return domain
}

func stringCompare(domain string, match string) bool {
	begin := strings.Index(match, "*")
	end := strings.Index(match[begin+1:], "*")
	if end != -1 {
		end += begin+1
	}
	if begin != -1 && end == -1 {
		// suffix match
		return strings.HasSuffix(domain, match[begin+1:])
	}
	if begin == -1 && end != -1 {
		// prefix match
		return strings.HasPrefix(domain, match[:end])
	}
	if begin == -1 && end == -1 {
		// full match
		if domain == match {
			return true
		} else {
			return false
		}
	}
	idx := strings.Index(domain, match[begin+1: end])
	if idx == -1 {
		return false
	}
	return true
}

func RouteUpdate()  {
	routeCtrl.Lock()
	defer routeCtrl.Unlock()

	newList := StringClone(DomainList())
	oldList := routeCtrl.domain

	delList, addList := StringDiff(oldList, newList)

	logs.Info("route update, domain %s delete", delList)
	logs.Info("route update, domain %s add", addList)

	for _, v := range delList {
		for address, value := range routeCtrl.cache {
			if value == v {
				delete(routeCtrl.cache, address)
				logs.Info("domain %s delete, address %s no match", v, address)
			}
		}
	}

	for _, v := range addList {
		for address, value := range routeCtrl.cache {
			if value == "" {
				domain := AddressToDomain(address)
				if stringCompare(domain, v) {
					routeCtrl.cache[address] = v
					logs.Info("domain %s add, address %s match", address, v)
				}
			}
		}
	}

	routeCtrl.domain = newList
}

// address: www.baidu.com:80 or www.baidu.com:443
func routeMatch(address string) string {
	domain := AddressToDomain(address)
	for _, v := range routeCtrl.domain {
		if stringCompare(domain, v) {
			routeCtrl.cache[address] = v
			logs.Info("route address %s match to domain %s", address, v)
			return v
		}
	}
	logs.Info("route address %s no match", address)
	routeCtrl.cache[address] = ""
	return ""
}

func RouteCheck(address string) bool {
	routeCtrl.RLock()
	result, flag := routeCtrl.cache[address]
	routeCtrl.RUnlock()

	if flag == false {
		routeCtrl.Lock()
		result = routeMatch(address)
		routeCtrl.Unlock()
	}

	if result == "" {
		return false
	}
	return true
}
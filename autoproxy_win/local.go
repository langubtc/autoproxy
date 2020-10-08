package main

import (
	"github.com/astaxie/beego/logs"
	"net"
)


type Options struct {
	Name   string
	Detail string
}

func ModeOptions() []*Options {
	return []*Options{
		{"local",LangValue("localforward")},
		{"auto",LangValue("autoforward")},
		{"proxy", LangValue("globalforward")},
	}
}

func ModeOptionGet() string {
	return ModeOptions()[ModeOptionsIdx()].Name
}

func ModeOptionsIdx() int {
	return int(DataIntValueGet("LocalMode"))
}

func ModeOptionsSet(idx int)  {
	err := DataIntValueSet("LocalMode", uint32(idx))
	if err != nil {
		logs.Error(err.Error())
	}
}

func PortOptionGet() int {
	value := DataIntValueGet("LocalPort")
	if value == 0 {
		value = 8080
	}
	return int(value)
}

func PortOptionSet(value int)  {
	err := DataIntValueSet("LocalPort", uint32(value))
	if err != nil {
		logs.Error(err.Error())
	}
}

var ifaceList []string

func IfaceOptions() []string {
	if ifaceList != nil {
		return ifaceList
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		ErrorBoxAction(mainWindow, err.Error())
		return nil
	}
	output := []string{"0.0.0.0"}
	for _, v := range ifaces {
		if v.Flags & net.FlagUp == 0 {
			continue
		}
		address, err := InterfaceLocalIP(&v)
		if err != nil {
			continue
		}
		if len(address) == 0 {
			continue
		}
		output = append(output, address[0].String())
	}
	ifaceList = output
	return output
}

func LocalIfaceOptionsIdx() int {
	ifaces := IfaceOptions()
	ifaceName := DataStringValueGet("LocalIface")
	for idx, v := range ifaces {
		if v == ifaceName {
			return idx
		}
	}
	return 0
}

func LocalIfaceOptionsSet(ifaceName string)  {
	err := DataStringValueSet("LocalIface", ifaceName)
	if err != nil {
		logs.Error(err.Error())
	}
}


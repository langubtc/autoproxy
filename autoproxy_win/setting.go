package main

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"strings"
)

const REGISTER_KEY = "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Internet Settings"

type ProxySetting struct {
	Override []string
	Enable   bool
	Server   string
}

func ProxyEnable() error {
	k, err := registry.OpenKey(registry.CURRENT_USER, REGISTER_KEY, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	value, _, err:= k.GetIntegerValue("ProxyEnable")
	if err != nil {
		return err
	}
	if value == 1 {
		return nil
	}
	return k.SetDWordValue("ProxyEnable", 1)
}

func ProxyDisable() error {
	k, err := registry.OpenKey(registry.CURRENT_USER, REGISTER_KEY, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	value, _, err:= k.GetIntegerValue("ProxyEnable")
	if err != nil {
		return err
	}
	if value == 0 {
		return nil
	}
	return k.SetDWordValue("ProxyEnable", 0)
}

func ProxySettingGet() *ProxySetting {
	k, err := registry.OpenKey(registry.CURRENT_USER, REGISTER_KEY, registry.ALL_ACCESS)
	if err != nil {
		return nil
	}
	var setting ProxySetting

	value, _, err:= k.GetIntegerValue("ProxyEnable")
	if err != nil {
		return nil
	}
	if value == 1 {
		setting.Enable = true
	}

	body, _, err := k.GetStringValue("ProxyServer")
	if err != nil {
		return nil
	}
	setting.Server = body

	body, _, err = k.GetStringValue("ProxyOverride")
	if err != nil {
		return nil
	}
	setting.Override = strings.Split(body, ";")
	return &setting
}

func ProxySettingSet(setting *ProxySetting) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, REGISTER_KEY, registry.ALL_ACCESS)
	if err != nil {
		return nil
	}

	err = k.SetStringValue("ProxyServer", setting.Server)
	if err != nil {
		return err
	}

	err = k.SetStringValue("ProxyOverride", StringList(setting.Override))
	if err != nil {
		return err
	}

	if setting.Enable {
		err = k.SetDWordValue("ProxyEnable", 1)
	} else {
		err = k.SetDWordValue("ProxyEnable", 0)
	}

	return err
}

func ProxyOverride(override []string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, REGISTER_KEY, registry.ALL_ACCESS)
	if err != nil {
		return nil
	}
	return k.SetStringValue("ProxyOverride", StringList(override))
}

func ProxyServer(server string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, REGISTER_KEY, registry.ALL_ACCESS)
	if err != nil {
		return nil
	}
	return k.SetStringValue("ProxyServer", server)
}

func StringList(list []string) string {
	var body string
	for idx,v := range list {
		if idx == len(list) - 1 {
			body += fmt.Sprintf("%s",v)
		}else {
			body += fmt.Sprintf("%s;",v)
		}
	}
	return body
}

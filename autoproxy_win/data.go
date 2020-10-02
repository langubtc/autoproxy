package main

import (
	"golang.org/x/sys/windows/registry"
)

const DATA_KEY = "SOFTWARE\\Autoproxy"

func keyGet() (registry.Key, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, DATA_KEY, registry.ALL_ACCESS)
	if err != nil {
		if err != registry.ErrNotExist {
			return 0, err
		}
		key, _, err = registry.CreateKey(registry.CURRENT_USER, DATA_KEY, registry.ALL_ACCESS)
		if err != nil {
			return 0, err
		}
	}
	return key, nil
}

func DataStringValueGet(name string) string {
	key, err := keyGet()
	if err != nil {
		return ""
	}
	value, _, err := key.GetStringValue(name)
	if err != nil {
		return ""
	}
	return value
}

func DataIntValueGet(name string ) uint32 {
	key, err := keyGet()
	if err != nil {
		return 0
	}
	value, _, err := key.GetIntegerValue(name)
	if err != nil {
		return 0
	}
	return uint32(value)
}

func DataStringValueSet(name string, value string) error {
	key, err := keyGet()
	if err != nil {
		return err
	}
	return key.SetStringValue(name, value)
}

func DataIntValueSet(name string, value uint32) error {
	key, err := keyGet()
	if err != nil {
		return err
	}
	return key.SetDWordValue(name, value)
}
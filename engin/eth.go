package main

import "net"

func EthGet() (map[string]string, error) {

	ips :=  make(map[string]string)

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range interfaces {
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			return nil, err
		}
		addresses, err := byName.Addrs()
		for _, v := range addresses {
			ips[byName.Name] = v.String()
		}
	}
	return ips, nil
}
package main

import (
	"flag"
	"log"
)

var (
	Help   bool
	Debug  bool

	ConfigFile string
)


func init()  {
	flag.StringVar(&ConfigFile,"config","config.yaml","configure file")
	flag.BoolVar(&Debug, "debug",false,"enable debug")
	flag.BoolVar(&Help,"help",false,"usage help")
}

func main()  {
	flag.Parse()
	if Help {
		flag.Usage()
		return
	}

	config, err := LoadConfig(ConfigFile)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = LogInit(config.Log, Debug)
	if err != nil {
		log.Fatalf(err.Error())
	}

	proxy := NewHttpProxyServer(config)
	proxy.Start()
}

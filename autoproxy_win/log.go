package main


import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
)

type logconfig struct {
	Filename string  `json:"filename"`
	Level    int     `json:"level"`
	MaxLines int     `json:"maxlines"`
	MaxSize  int     `json:"maxsize"`
	Daily    bool    `json:"daily"`
	MaxDays  int     `json:"maxdays"`
	Color    bool    `json:"color"`
}

var logCfg = logconfig{Filename: os.Args[0], Level: 7, Daily: true, MaxDays: 30, Color: true}

func LogInit() error {
	logCfg.Filename = fmt.Sprintf("%s%c%s", logDirGet(), os.PathSeparator, "autoproxy.log")
	value, err := json.Marshal(&logCfg)
	if err != nil {
		return err
	}
	err = logs.SetLogger(logs.AdapterFile, string(value))
	if err != nil {
		return err
	}
	return nil
}



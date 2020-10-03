package main

import (
	"github.com/astaxie/beego/logs"
)

func main()  {
	err := FileInit()
	if err != nil {
		logs.Error(err.Error())
	}
	err = LogInit()
	if err != nil {
		logs.Error(err.Error())
	}
	err = BoxInit()
	if err != nil {
		logs.Error(err.Error())
	}
	err = LanguageInit()
	if err != nil {
		logs.Error(err.Error())
	}
	err = AuthInit()
	if err != nil {
		logs.Error(err.Error())
	}
	err = IconInit()
	if err != nil {
		logs.Error(err.Error())
	}
	mainWindows()
}

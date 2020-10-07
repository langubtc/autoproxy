package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var consoleIface *walk.ComboBox
var consoleRemoteProxy *walk.ComboBox
var consoleMode *walk.ComboBox
var consolePort *walk.NumberEdit

func ConsoleEnable(enable bool)  {
	consoleIface.SetEnabled(enable)
	consolePort.SetEnabled(enable)
}

func ConsoleRemoteUpdate()  {
	consoleRemoteProxy.SetModel(RemoteOptions())
	consoleRemoteProxy.SetCurrentIndex(RemoteIndexGet())
}

func ConsoleWidget() []Widget {
	return []Widget{
		Label{
			Text: LangValue("localaddress") + ":",
		},
		ComboBox{
			AssignTo: &consoleIface,
			CurrentIndex:  LocalIfaceOptionsIdx(),
			Model:         IfaceOptions(),
			OnCurrentIndexChanged: func() {
				LocalIfaceOptionsSet(consoleIface.Text())
			},
		},
		Label{
			Text: LangValue("port") + ":",
		},
		NumberEdit{
			AssignTo: &consolePort,
			Value:    float64(PortOptionGet()),
			ToolTipText: "1~65535",
			MaxValue: 65535,
			MinValue: 1,
			OnValueChanged: func() {
				PortOptionSet(int(consolePort.Value()))
			},
		},
		Label{
			Text: LangValue("mode") + ":",
		},
		ComboBox{
			AssignTo: &consoleMode,
			BindingMember: "Name",
			DisplayMember: "Detail",
			CurrentIndex:  ModeOptionsIdx(),
			Model:         ModeOptions(),
			OnCurrentIndexChanged: func() {
				ModeOptionsSet(consoleMode.CurrentIndex())
				go func() {
					ModeUpdate()
				}()
			},
		},
		Label{
			Text: LangValue("remoteproxy") + ":",
		},
		ComboBox{
			AssignTo:      &consoleRemoteProxy,
			CurrentIndex:  RemoteIndexGet(),
			OnCurrentIndexChanged: func() {
				consoleRemoteProxy.SetEnabled(false)
				RemoteIndexSet(consoleRemoteProxy.Text())
				go func() {
					err := RemoteForwardUpdate()
					if err != nil {
						ErrorBoxAction(mainWindow, err.Error())
					}
					consoleRemoteProxy.SetEnabled(true)
				}()
			},
			Model:         RemoteOptions(),
		},
	}
}

func InternalSettingEnable() error {
	ifaceAddr := IfaceOptions()[LocalIfaceOptionsIdx()]
	if ifaceAddr == "0.0.0.0" {
		ifaceAddr = "127.0.0.1"
	}
	address := fmt.Sprintf("%s:%d", ifaceAddr, PortOptionGet())
	err := ProxyServer(address)
	if err != nil {
		logs.Error("setting proxy server fail, %s", err.Error())
		return err
	}
	err = ProxyEnable()
	if err != nil {
		logs.Error("setting proxy enable fail, %s", err.Error())
		return err
	}
	return nil
}

func ButtonWight() []Widget {
	var start *walk.PushButton
	var stop *walk.PushButton

	return []Widget{
		PushButton{
			AssignTo:  &start,
			Text:      LangValue("start"),
			OnClicked: func() {
				start.SetEnabled(false)
				ConsoleEnable(false)
				go func() {
					err := ServerStart()
					if err != nil {
						ErrorBoxAction(mainWindow, err.Error())
						start.SetEnabled(true)
						ConsoleEnable(true)
					} else {
						err = InternalSettingEnable()
						if err != nil {
							ErrorBoxAction(mainWindow, err.Error())
						}
						StatRunningStatus(true)
						stop.SetEnabled(true)
					}
				}()
			},
		},
		PushButton{
			AssignTo:  &stop,
			Enabled:   false,
			Text:      LangValue("stop"),
			OnClicked: func() {
				stop.SetEnabled(false)
				go func() {
					err := ServerShutdown()
					if err != nil {
						ErrorBoxAction(mainWindow, err.Error())
						stop.SetEnabled(true)
					} else {
						err = ProxyDisable()
						if err != nil {
							logs.Error("setting proxy disable fail, %s", err.Error())
							ErrorBoxAction(mainWindow, err.Error())
						}
						StatRunningStatus(false)
						start.SetEnabled(true)
						ConsoleEnable(true)
					}
					ConsoleRemoteUpdate()
				}()
			},
		},
	}
}
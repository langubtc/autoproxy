package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var consoleLocalAddress *walk.Label
var consoleAuth *walk.RadioButton
var consoleRemoteProxy *walk.ComboBox
var consoleMode *walk.ComboBox

func ConsoleUpdate()  {
	consoleLocalAddress.SetText(LocalAddressGet())
	consoleAuth.SetChecked(AuthSwitchGet())
	consoleMode.SetCurrentIndex(ModeOptionsIdx())
}

func ConsoleRemoteUpdate()  {
	consoleRemoteProxy.SetModel(RemoteList())
}

func ConsoleWidget() []Widget {
	return []Widget{
		Label{
			Text: LangValue("localaddress") + ":",
		},
		Label {
			AssignTo: &consoleLocalAddress,
			Text: LocalAddressGet(),
		},
		Label{
			Text: LangValue("whetherauth") + ":",
		},
		RadioButton{
			AssignTo: &consoleAuth,
			OnBoundsChanged: func() {
				consoleAuth.SetChecked(AuthSwitchGet())
			},
			OnClicked: func() {
				consoleAuth.SetChecked(!AuthSwitchGet())
				AuthSwitchSet(!AuthSwitchGet())
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
				if ProtcalOptionsGet() == "local" {
					consoleRemoteProxy.SetEnabled(false)
				} else {
					consoleRemoteProxy.SetEnabled(true)
				}
			},
		},
		Label{
			Text: LangValue("remoteproxy") + ":",
		},
		ComboBox{
			AssignTo:      &consoleRemoteProxy,
			CurrentIndex:  0,
			OnBoundsChanged: func() {
				if ProtcalOptionsGet() == "local" {
					consoleRemoteProxy.SetEnabled(false)
				} else {
					consoleRemoteProxy.SetEnabled(true)
				}
			},
			Model:         RemoteList(),
		},
	}
}

func InternalSettingEnable() error {
	address := fmt.Sprintf("%s:%d",
		IfaceOptions()[LocalIfaceOptionsIdx()],
		PortOptionGet())

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

				go func() {
					err := ServerStart()
					if err != nil {
						ErrorBoxAction(mainWindow, err.Error())
						start.SetEnabled(true)
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
					}
				}()
			},
		},
	}
}
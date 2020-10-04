package main

import (
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
			},
		},
		Label{
			Text: LangValue("remoteproxy") + ":",
		},
		ComboBox{
			AssignTo:      &consoleRemoteProxy,
			CurrentIndex:  0,
			Model:         RemoteList(),
		},
	}
}

func ButtonWight() []Widget {
	return []Widget{
		PushButton{
			Text:     LangValue("start"),
			OnClicked: func() {

			},
		},
		PushButton{
			//Enabled: false,
			Text:     LangValue("stop"),
			OnClicked: func() {

			},
		},
	}
}
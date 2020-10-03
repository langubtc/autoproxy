package main

import (
	. "github.com/lxn/walk/declarative"
)

func ConsoleWidget() []Widget {
	return []Widget{
		Label{
			Text: LangValue("localaddress") + ":",
		},
		Label {
			Text: "http://192.168.3.11:8080",
		},
		Label{
			Text: LangValue("whetherauth") + ":",
		},
		RadioButton{
			OnBoundsChanged: func() {
			},
			OnClicked: func() {
			},
		},
		Label{
			Text: LangValue("mode") + ":",
		},
		ComboBox{
			BindingMember: "Name",
			DisplayMember: "Detail",
			CurrentIndex:  0,
			Model:         ModeOptions(),
		},
		Label{
			Text: LangValue("remoteproxy") + ":",
		},
		ComboBox{
			BindingMember: "Name",
			DisplayMember: "Detail",
			CurrentIndex:  0,
			Model:         RemoteOptions(),
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
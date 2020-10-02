package main

import (
	. "github.com/lxn/walk/declarative"
)

func ConsoleWidget() []Widget {
	return []Widget{
		Label{
			Text: "本地地址:",
		},
		Label {
			Text: "http://192.168.3.11:8080",
		},
		Label{
			Text: "是否认证:",
		},
		RadioButton{
			OnBoundsChanged: func() {
			},
			OnClicked: func() {
			},
		},
		Label{
			Text: "代理模式:",
		},
		ComboBox{
			BindingMember: "Name",
			DisplayMember: "Detail",
			CurrentIndex:  0,
			Model:         ModeOptions(),
		},
		Label{
			Text: "二级代理:",
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
			Text:     "启动服务",
			OnClicked: func() {

			},
		},
		PushButton{
			//Enabled: false,
			Text:     "停止服务",
			OnClicked: func() {

			},
		},
	}
}
package main

import (
	. "github.com/lxn/walk/declarative"
)

type Options struct {
	Name   string
	Detail string
}

func ModeOptions() []*Options {
	return []*Options{
		{"auto","自动转发"},
		{"local","本地模式"},
		{"proxy","全局转发"},
	}
}

func RemoteOptions() []*Options {
	return []*Options{
		{"easymesh.cc:8080","easymesh.cc"},
	}
}

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
package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"net"
)

var ifaceList []string

func IfaceOptions() []string {
	if ifaceList != nil {
		return ifaceList
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		ErrorBoxAction(mainWindow, err.Error())
		return nil
	}
	output := []string{"0.0.0.0"}
	for _, v := range ifaces {
		if v.Flags & net.FlagUp == 0 {
			continue
		}
		address, err := InterfaceLocalIP(&v)
		if err != nil {
			continue
		}
		if len(address) == 0 {
			continue
		}
		output = append(output, address[0].String())
	}
	ifaceList = output
	return output
}

func localWidget() []Widget {
	return []Widget{
		Label{
			Text: "本地地址:",
		},
		ComboBox{
			CurrentIndex:  0,
			Model:         IfaceOptions(),
		},
		Label{
			Text: "端口:",
		},
		LineEdit{
			Text: "8080",
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
			Text: "是否认证:",
		},
		RadioButton{
			OnBoundsChanged: func() {
			},
			OnClicked: func() {
			},
		},
		Label{
			Text: "自动启动:",
		},
		RadioButton{
			OnBoundsChanged: func() {
			},
			OnClicked: func() {
			},
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

func LocalServer()  {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	_, err := Dialog{
		AssignTo: &dlg,
		Title: "本地代理",
		Icon: walk.IconShield(),
		DefaultButton: &acceptPB,
		CancelButton: &cancelPB,
		Size: Size{400, 200},
		MinSize: Size{400, 200},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2, MarginsZero: true},
				Children: localWidget(),
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &acceptPB,
						Text:     "确认",
						OnClicked: func() {

						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "取消",
						OnClicked: func() {
							dlg.Cancel()
						},
					},
				},
			},
		},
	}.Run(mainWindow)

	if err != nil {
		log.Println(err.Error())
	}
}
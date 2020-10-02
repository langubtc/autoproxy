package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"net"
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

func ModeOptionsIdx() int {
	return int(DataIntValueGet("LocalMode"))
}

func ModeOptionsSet(idx int)  {
	err := DataIntValueSet("LocalMode", uint32(idx))
	if err != nil {
		log.Println(err.Error())
	}
}

func ProtocalOptions() []string {
	return []string{
		"HTTP","HTTPS","SOCK5",
	}
}

func ProtocalOptionsIdx() int {
	return int(DataIntValueGet("LocalProtocal"))
}

func ProtcalOptionsSet(idx int)  {
	err := DataIntValueSet("LocalProtocal", uint32(idx))
	if err != nil {
		log.Println(err.Error())
	}
}

func TlsOptions() []string {
	return []string{
		"TLS1.1","TLS1.2","TLS1.3",
	}
}

func TlsOptionsIdx() int {
	return int(DataIntValueGet("LocalTls"))
}

func TlsOptionsSet(idx int)  {
	err := DataIntValueSet("LocalTls", uint32(idx))
	if err != nil {
		log.Println(err.Error())
	}
}

func PortOptionGet() float64 {
	value := DataIntValueGet("LocalPort")
	if value == 0 {
		value = 8080
	}
	return float64(value)
}

func PortOptionSet(value int)  {
	err := DataIntValueSet("LocalPort", uint32(value))
	if err != nil {
		log.Println(err.Error())
	}
}

func RemoteOptions() []*Options {
	return []*Options{
		{"easymesh.cc:8080","easymesh.cc"},
	}
}


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
	var protocal, mode, tls *walk.ComboBox
	var port *walk.NumberEdit

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
		NumberEdit{
			AssignTo: &port,
			Value:    PortOptionGet(),
			MaxValue: 65535,
			MinValue: 1,
			OnValueChanged: func() {
				PortOptionSet(int(port.Value()))
			},
		},
		Label{
			Text: "代理模式:",
		},
		ComboBox{
			AssignTo: &mode,
			BindingMember: "Name",
			DisplayMember: "Detail",
			CurrentIndex:  ModeOptionsIdx(),
			Model:         ModeOptions(),
			OnCurrentIndexChanged: func() {
				ModeOptionsSet(mode.CurrentIndex())
			},
		},
		Label{
			Text: "接入协议:",
		},
		ComboBox{
			AssignTo: &protocal,
			CurrentIndex:  ProtocalOptionsIdx(),
			Model:         ProtocalOptions(),
			OnCurrentIndexChanged: func() {
				ProtcalOptionsSet(protocal.CurrentIndex())
			},
		},
		Label{
			Text: "安全协议:",
		},
		ComboBox{
			AssignTo: &tls,
			CurrentIndex:  TlsOptionsIdx(),
			Model:         TlsOptions(),
			OnCurrentIndexChanged: func() {
				TlsOptionsSet(tls.CurrentIndex())
			},
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
		Size: Size{250, 300},
		MinSize: Size{250, 300},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
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
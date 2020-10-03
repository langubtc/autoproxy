package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"net"
)


type Options struct {
	Name   string
	Detail string
}

func ModeOptions() []*Options {
	return []*Options{
		{"auto",LangValue("autoforward")},
		{"local",LangValue("localforward")},
		{"proxy", LangValue("globalforward")},
	}
}

func ModeOptionsIdx() int {
	return int(DataIntValueGet("LocalMode"))
}

func ModeOptionsSet(idx int)  {
	err := DataIntValueSet("LocalMode", uint32(idx))
	if err != nil {
		logs.Error(err.Error())
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
		logs.Error(err.Error())
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
		logs.Error(err.Error())
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
		logs.Error(err.Error())
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

func LocalIfaceOptionsIdx() int {
	ifaceName := DataStringValueGet("LocalIface")
	for idx, v := range ifaceList {
		if v == ifaceName {
			return idx
		}
	}
	return 0
}

func LocalIfaceOptionsSet(ifaceName string)  {
	err := DataStringValueSet("LocalIface", ifaceName)
	if err != nil {
		logs.Error(err.Error())
	}
}

func AuthSwitchGet() bool {
	if DataIntValueGet("authswtich") > 0 {
		return true
	}
	return false
}

func AuthSwitchSet(flag bool)  {
	if flag {
		DataIntValueSet("authswtich", 1)
	} else {
		DataIntValueSet("authswtich", 0)
	}
}

func localWidget() []Widget {
	var iface, protocal, mode, tls *walk.ComboBox
	var port *walk.NumberEdit
	var auth *walk.RadioButton

	return []Widget{
		Label{
			Text: LangValue("localaddress") + ":",
		},
		ComboBox{
			AssignTo: &iface,
			CurrentIndex:  LocalIfaceOptionsIdx(),
			RightToLeftReading: false,
			Model:         IfaceOptions(),
			OnCurrentIndexChanged: func() {
				LocalIfaceOptionsSet(iface.Text())
			},
		},
		Label{
			Text: LangValue("port") + ":",
		},
		NumberEdit{
			AssignTo: &port,
			Value:    PortOptionGet(),
			ToolTipText: "1~65535",
			MaxValue: 65535,
			MinValue: 1,
			OnValueChanged: func() {
				PortOptionSet(int(port.Value()))
			},
		},
		Label{
			Text: LangValue("mode") + ":",
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
			Text: LangValue("protocal") + ":",
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
			Text: LangValue("tlsversion") + ":",
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
			Text: LangValue("whetherauth") + ":",
		},
		RadioButton{
			AssignTo: &auth,
			OnBoundsChanged: func() {
				auth.SetChecked(AuthSwitchGet())
			},
			OnClicked: func() {
				auth.SetChecked(!AuthSwitchGet())
				AuthSwitchSet(!AuthSwitchGet())
			},
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

func LocalServer()  {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	_, err := Dialog{
		AssignTo: &dlg,
		Title: LangValue("localproxy"),
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
						Text:     LangValue("accpet"),
						OnClicked: func() {

						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      LangValue("cancel"),
						OnClicked: func() {
							dlg.Cancel()
						},
					},
				},
			},
		},
	}.Run(mainWindow)

	if err != nil {
		logs.Error(err.Error())
	}
}
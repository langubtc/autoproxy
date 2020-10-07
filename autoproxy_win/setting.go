package main

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"strings"
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const REGISTER_KEY = "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Internet Settings"

type ProxySetting struct {
	Override []string
	Enable   bool
	Server   string
}

func ProxyEnable() error {
	k, err := registry.OpenKey(registry.CURRENT_USER, REGISTER_KEY, registry.ALL_ACCESS)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	defer k.Close()

	value, _, err:= k.GetIntegerValue("ProxyEnable")
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	if value == 1 {
		return nil
	}
	return k.SetDWordValue("ProxyEnable", 1)
}

func ProxyDisable() error {
	k, err := registry.OpenKey(registry.CURRENT_USER, REGISTER_KEY, registry.ALL_ACCESS)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	defer k.Close()

	value, _, err:= k.GetIntegerValue("ProxyEnable")
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	if value == 0 {
		return nil
	}
	return k.SetDWordValue("ProxyEnable", 0)
}

func ProxySettingGet() *ProxySetting {
	k, err := registry.OpenKey(registry.CURRENT_USER, REGISTER_KEY, registry.ALL_ACCESS)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	defer k.Close()

	var setting ProxySetting

	value, _, err:= k.GetIntegerValue("ProxyEnable")
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	if value == 1 {
		setting.Enable = true
	}

	body, _, err := k.GetStringValue("ProxyServer")
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	setting.Server = body

	body, _, err = k.GetStringValue("ProxyOverride")
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	setting.Override = strings.Split(body, ";")
	return &setting
}

func ProxySettingSet(setting *ProxySetting) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, REGISTER_KEY, registry.ALL_ACCESS)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	defer k.Close()

	err = k.SetStringValue("ProxyServer", setting.Server)
	if err != nil {
		logs.Error(err.Error())
		return err
	}

	err = k.SetStringValue("ProxyOverride", StringList(setting.Override))
	if err != nil {
		logs.Error(err.Error())
		return err
	}

	if setting.Enable {
		err = k.SetDWordValue("ProxyEnable", 1)
	} else {
		err = k.SetDWordValue("ProxyEnable", 0)
	}

	return err
}

func ProxyOverride(override []string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, REGISTER_KEY, registry.ALL_ACCESS)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	defer k.Close()
	return k.SetStringValue("ProxyOverride", StringList(override))
}

func ProxyServer(server string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, REGISTER_KEY, registry.ALL_ACCESS)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	defer k.Close()
	return k.SetStringValue("ProxyServer", server)
}

var proxyServer *walk.LineEdit
var override *walk.TextEdit
var usingproxy *walk.RadioButton
var followlocalservice *walk.RadioButton

var proxysetting *ProxySetting

func InternetSettingWidget() []Widget {
	proxysetting = ProxySettingGet()

	return []Widget{
		Label{
			Text: LangValue("proxyserver") + ":",
		},
		LineEdit{
			AssignTo: &proxyServer,
			Text: proxysetting.Server,
		},
		PushButton{
			Text:     LangValue("syncaddress"),
			OnClicked: func() {
				server := ProxyServerGet()
				proxysetting.Server = server
				proxyServer.SetText(server)
			},
		},
		Label{
			Text: LangValue("override") + ":",
		},
		TextEdit{
			AssignTo: &override,
			Text: StringList(proxysetting.Override),
		},
		HSpacer{

		},
		Label{
			Text: LangValue("usingproxy") + ":",
		},
		RadioButton{
			AssignTo: &usingproxy,
			OnBoundsChanged: func() {
				usingproxy.SetChecked(proxysetting.Enable)
			},
			OnClicked: func() {
				usingproxy.SetChecked(!proxysetting.Enable)
				proxysetting.Enable = !proxysetting.Enable
			},
		},
		HSpacer{

		},
	}
}

func ProxyServerGet() string {
	ifaceAddr := IfaceOptions()[LocalIfaceOptionsIdx()]
	if ifaceAddr == "0.0.0.0" {
		ifaceAddr = "127.0.0.1"
	}
	return fmt.Sprintf("%s:%d", ifaceAddr, PortOptionGet())
}

func InternetSetting()  {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	_, err := Dialog{
		AssignTo: &dlg,
		Title: LangValue("internetsettings"),
		Icon: walk.IconShield(),
		DefaultButton: &acceptPB,
		CancelButton: &cancelPB,
		Size: Size{400, 300},
		MinSize: Size{400, 300},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 3},
				Children: InternetSettingWidget(),
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &acceptPB,
						Text:     LangValue("setting"),
						OnClicked: func() {
							err := ProxySettingSet(proxysetting)
							if err != nil {
								ErrorBoxAction(dlg, err.Error())
							} else {
								InfoBoxAction(dlg, LangValue("settingsuccess"))
								dlg.Accept()
							}
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
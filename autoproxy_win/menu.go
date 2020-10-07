package main

import (
	. "github.com/lxn/walk/declarative"
)

func MenuBarInit() []MenuItem {
	return []MenuItem{
		Menu{
			Text: LangValue("setting"),
			Items: []MenuItem{
				Action{
					Text: LangValue("basesetting"),
					OnTriggered: func() {
						BaseSetting()
					},
				},
				Action{
					Text: LangValue("internetsettings"),
					OnTriggered: func() {
						InternetSetting()
					},
				},
				Action{
					Text: LangValue("runlog"),
					OnTriggered: func() {
						OpenBrowserWeb(logDirGet())
					},
				},
				Separator{},
				Action{
					Text: LangValue("exit"),
					OnTriggered: func() {
						CloseWindows()
					},
				},
			},
		},
		Action{
			Text:     LangValue("forwarddomain"),
			OnTriggered: func() {
				RemodeEdit()
			},
		},
		Action{
			Text: LangValue("remoteproxy"),
			OnTriggered: func() {
				RemoteServer()
			},
		},
		/*
		Menu{
			Text: LangValue("authcred"),
			Items: []MenuItem{
				Action{
					Text:     LangValue("viewcred"),
					OnTriggered: func() {
						AuthView()
					},
				},
				Action{
					Text:     LangValue("addcred"),
					OnTriggered: func() {
						AuthAdd()
					},
				},
			},
		},*/
		Action{
			Text: LangValue("miniwin"),
			OnTriggered: func() {
				Notify()
			},
		},
		Action{
			Text: LangValue("about"),
			OnTriggered: func() {
				AboutAction()
			},
		},
	}
}
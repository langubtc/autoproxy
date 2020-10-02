package main

import (
	. "github.com/lxn/walk/declarative"
)

func MenuBarInit() []MenuItem {
	return []MenuItem{
		Menu{
			Text: "配置",
			Items: []MenuItem{
				Action{
					Text:     "本地服务",
					OnTriggered: func() {
						LocalServer()
					},
				},
				Action{
					Text:     "二级代理",
				},
				Separator{},
				Action{
					Text:        "退出",
					OnTriggered: func() {
						CloseWindows()
					},
				},
			},
		},
		Menu{
			Text: "认证凭证",
			Items: []MenuItem{
				Action{
					Text:     "查看凭证",
					OnTriggered: func() {
						AuthView()
					},
				},
				Action{
					Text:     "添加凭证",
					OnTriggered: func() {
						AuthAdd()
					},
				},
			},
		},
		Action{
			Text:     "最小化窗口",
			OnTriggered: func() {
				Notify()
			},
		},
		Action{
			Text:        "关于",
			OnTriggered: func() {
				AboutAction()
			},
		},
	}
}
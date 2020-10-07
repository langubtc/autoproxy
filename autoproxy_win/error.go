package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func boxAction(from walk.Form, title string, icon *walk.Icon, message string)  {
	var dlg *walk.Dialog
	var cancelPB *walk.PushButton

	_, err := Dialog{
		AssignTo: &dlg,
		Title: title,
		Icon: icon,
		CancelButton: &cancelPB,
		Size: Size{150, 150},
		MinSize: Size{150, 150},
		MaxSize: Size{200, 300},
		Layout:  VBox{},
		Children: []Widget{
			Label{
				Text: message,
			},
			PushButton{
				AssignTo:  &cancelPB,
				Text:      LangValue("accpet"),
				OnClicked: func() {
					dlg.Accept()
				},
			},
		},
	}.Run(from)

	if err != nil {
		logs.Error(err.Error())
	}
}

func ErrorBoxAction(form walk.Form, message string) {
	boxAction(form, LangValue("error"), walk.IconError(), message)
}

func InfoBoxAction(form walk.Form, message string) {
	boxAction(form, LangValue("info"), walk.IconInformation(), message)
}

func ConfirmBoxAction(form walk.Form, message string) {
	boxAction(form, LangValue("confirm"), walk.IconWarning(), message)
}

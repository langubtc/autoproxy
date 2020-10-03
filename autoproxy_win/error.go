package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"time"
)

func boxAction(from walk.Form, title string, icon *walk.Icon, message string, timeout time.Duration)  {
	var dlg *walk.Dialog
	var cancelPB *walk.PushButton

	if timeout > 0 {
		go func() {
			for  {
				time.Sleep(timeout)
				if dlg != nil && dlg.Visible() {
					dlg.Cancel()
					break
				}
			}
		}()
	}

	_, err := Dialog{
		AssignTo: &dlg,
		Title: title,
		Icon: icon,
		CancelButton: &cancelPB,
		Size: Size{150, 150},
		MinSize: Size{150, 150},
		Layout:  VBox{},
		Children: []Widget{
			Label{
				Text: message,
			},
			PushButton{
				AssignTo:  &cancelPB,
				Text:      LangValue("accpet"),
				OnClicked: func() {
					dlg.Cancel()
				},
			},
		},
	}.Run(from)

	if err != nil {
		logs.Error(err.Error())
	}
}

func ErrorBoxAction(form walk.Form, message string) {
	boxAction(form, LangValue("error"), walk.IconError(), message, 0)
}

func InfoBoxAction(form walk.Form, message string) {
	boxAction(form, LangValue("info"), walk.IconInformation(), message, 2*time.Second)
}

func ConfirmBoxAction(form walk.Form, message string) {
	boxAction(form, LangValue("confirm"), walk.IconWarning(), message, 0)
}
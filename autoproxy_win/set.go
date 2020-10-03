package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)


func SettingWidget() []Widget {
	var lang *walk.ComboBox

	return []Widget{
		Label{
			Text: LangValue("langname") + ":",
		},
		ComboBox{
			AssignTo: &lang,
			CurrentIndex:  LangOptionIdx(),
			Model:         LangOptionGet(),
			OnCurrentIndexChanged: func() {
				LangOptionSet(lang.CurrentIndex())
			},
		},
	}
}

func BaseSetting()  {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	_, err := Dialog{
		AssignTo: &dlg,
		Title: LangValue("basesetting"),
		Icon: walk.IconShield(),
		DefaultButton: &acceptPB,
		CancelButton: &cancelPB,
		Size: Size{250, 200},
		MinSize: Size{250, 200},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: SettingWidget(),
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &acceptPB,
						Text:     LangValue("accpet"),
						OnClicked: func() {
							go func() {
								InfoBoxAction(dlg, LangValue("rebootsetting"))
								dlg.Accept()
							}()
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

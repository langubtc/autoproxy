package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"os/exec"
)


func OpenBrowserWeb(url string)  {
	cmd := exec.Command("rundll32","url.dll,FileProtocolHandler", url)
	err := cmd.Run()
	if err != nil {
		log.Printf("run cmd fail, %s", err.Error())
	}
}

var aboutsCtx string

func AboutAction() {
	var ok    *walk.PushButton
	var about *walk.Dialog
	var err error

	if aboutsCtx == "" {
		aboutsCtx, err = BoxFile().String("about.txt")
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	_, err = Dialog{
		AssignTo:      &about,
		Title:         "关于",
		Icon:          walk.IconInformation(),
		DefaultButton: &ok,
		Layout:  VBox{},
		Children: []Widget{
			TextLabel{
				Text: aboutsCtx,
				MinSize: Size{Width: 200, Height: 250},
			},
			ImageView{

			},
			Composite{
				Layout: VBox{},
				Children: []Widget{
					PushButton{
						Text:      "官方网站",
						OnClicked: func() {
							OpenBrowserWeb("https://easymesh.info")
						},
					},
					PushButton{
						Text:      "确认",
						OnClicked: func() { about.Cancel() },
					},
				},
			},
		},
	}.Run(mainWindow)

	if err != nil {
		log.Println(err.Error())
	}
}

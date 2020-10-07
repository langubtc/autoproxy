package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"time"
)

var mainWindow *walk.MainWindow

var mainWindowWidth = 300
var mainWindowHeight = 350

func waitWindows()  {
	for  {
		if mainWindow != nil && mainWindow.Visible() {
			break
		}
		time.Sleep(100*time.Millisecond)
	}
	NotifyInit()
}

func statusUpdate()  {
	StatUpdate(StatGet())
}

func init()  {
	go func() {
		waitWindows()
		for  {
			statusUpdate()
			time.Sleep(2 * time.Second)
		}
	}()
}

var isAuth *walk.RadioButton
var protocal  *walk.RadioButton

func mainWindows() {
	CapSignal(CloseWindows)
	cnt, err := MainWindow{
		Title:   "AutoProxy " + VersionGet(),
		Icon: ICON_Main,
		AssignTo: &mainWindow,
		MinSize: Size{mainWindowWidth, mainWindowHeight},
		Size: Size{mainWindowWidth, mainWindowHeight},
		Layout:  VBox{},
		MenuItems: MenuBarInit(),
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 3},
				Children: StatWidget(),
			},
			Composite{
				Layout: Grid{Columns: 2, MarginsZero: true},
				Children: ConsoleWidget(),
			},
			Composite{
				Layout: Grid{Columns: 2, MarginsZero: true},
				Children: ButtonWight(),
			},
		},
	}.Run()

	if err != nil {
		logs.Error(err.Error())
	} else {
		logs.Info("main windows exit %d", cnt)
	}

	CloseWindows()
}

func CloseWindows()  {
	err := ProxyDisable()
	if err != nil {
		logs.Error(err.Error())
	}
	if mainWindow != nil {
		mainWindow.Close()
		mainWindow = nil
	}
	NotifyExit()
}
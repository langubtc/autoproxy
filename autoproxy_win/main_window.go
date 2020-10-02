package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"time"
)

var mainWindow *walk.MainWindow

var mainWindowWidth = 300
var mainWindowHeight = 450

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
	NotifyUpdate(ICON_Network_Disable, "")
}

func init()  {
	go func() {
		waitWindows()
		for  {
			statusUpdate()
			time.Sleep(time.Second)
		}
	}()
}

var isAuth *walk.RadioButton
var protocal  *walk.RadioButton

func mainWindows() {
	CapSignal(CloseWindows)
	MainWindow{
		Title:   "AutoProxy " + VersionGet(),
		Icon: ICON_Main,
		AssignTo: &mainWindow,
		MinSize: Size{mainWindowWidth, mainWindowHeight},
		Size: Size{mainWindowWidth, mainWindowHeight},
		Layout:  VBox{},
		MenuItems: MenuBarInit(),
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2, MarginsZero: true},
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


}

func CloseWindows()  {
	if mainWindow != nil {
		mainWindow.Close()
	}
	NotifyExit()
}
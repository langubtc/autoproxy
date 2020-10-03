package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	"time"
)

var notify *walk.NotifyIcon

func NotifyUpdate(icon *walk.Icon, flow string)  {
	if notify == nil {
		return
	}
	notify.SetIcon(icon)
	notify.SetToolTip(flow)
}

func NotifyExit()  {
	if notify == nil {
		return
	}
	notify.Dispose()
	notify = nil
}

var lastCheck time.Time

func NotifyInit()  {
	var err error

	notify, err = walk.NewNotifyIcon(mainWindow)
	if err != nil {
		logs.Error("new notify icon fail, %s", err.Error())
		return
	}

	exitBut := walk.NewAction()
	err = exitBut.SetText(LangValue("exit"))
	if err != nil {
		logs.Error("notify new action fail, %s", err.Error())
		return
	}

	exitBut.Triggered().Attach(func() {
		walk.App().Exit(0)
	})

	if err := notify.ContextMenu().Actions().Add(exitBut); err != nil {
		logs.Error("notify add action fail, %s", err.Error())
		return
	}

	notify.MouseUp().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}
		now := time.Now()
		if now.Sub(lastCheck) < 2 * time.Second {
			mainWindow.SetVisible(true)
		}
		lastCheck = now
	})

	notify.SetVisible(true)
}

func Notify()  {
	if notify == nil {
		NotifyInit()
	}
	mainWindow.SetVisible(false)
}
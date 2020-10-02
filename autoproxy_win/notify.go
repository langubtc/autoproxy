package main

import (
	"github.com/lxn/walk"
	"log"
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
}

var lastCheck time.Time

func NotifyInit()  {
	var err error

	notify, err = walk.NewNotifyIcon(mainWindow)
	if err != nil {
		log.Printf("new notify icon fail, %s", err.Error())
		return
	}

	exitBut := walk.NewAction()
	err = exitBut.SetText("&Exit")
	if err != nil {
		log.Printf("notify new action fail, %s", err.Error())
		return
	}

	exitBut.Triggered().Attach(func() {
		walk.App().Exit(0)
	})

	if err := notify.ContextMenu().Actions().Add(exitBut); err != nil {
		log.Printf("notify add action fail, %s", err.Error())
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
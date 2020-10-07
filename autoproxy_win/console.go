package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sync"
	"time"
)

var consoleIface *walk.ComboBox
var consoleRemoteProxy *walk.ComboBox
var consoleMode *walk.ComboBox
var consolePort *walk.NumberEdit

func ConsoleEnable(enable bool)  {
	consoleIface.SetEnabled(enable)
	consolePort.SetEnabled(enable)
}

func ConsoleRemoteUpdate()  {
	consoleRemoteProxy.SetModel(RemoteOptions())
	consoleRemoteProxy.SetCurrentIndex(RemoteIndexGet())
}

func ConsoleWidget() []Widget {
	return []Widget{
		Label{
			Text: LangValue("localaddress") + ":",
		},
		ComboBox{
			AssignTo: &consoleIface,
			CurrentIndex:  LocalIfaceOptionsIdx(),
			Model:         IfaceOptions(),
			OnCurrentIndexChanged: func() {
				LocalIfaceOptionsSet(consoleIface.Text())
			},
		},
		Label{
			Text: LangValue("port") + ":",
		},
		NumberEdit{
			AssignTo: &consolePort,
			Value:    float64(PortOptionGet()),
			ToolTipText: "1~65535",
			MaxValue: 65535,
			MinValue: 1,
			OnValueChanged: func() {
				PortOptionSet(int(consolePort.Value()))
			},
		},
		Label{
			Text: LangValue("mode") + ":",
		},
		ComboBox{
			AssignTo: &consoleMode,
			BindingMember: "Name",
			DisplayMember: "Detail",
			CurrentIndex:  ModeOptionsIdx(),
			Model:         ModeOptions(),
			OnCurrentIndexChanged: func() {
				ModeOptionsSet(consoleMode.CurrentIndex())
				go func() {
					ModeUpdate()
				}()
			},
		},
		Label{
			Text: LangValue("remoteproxy") + ":",
		},
		ComboBox{
			AssignTo:      &consoleRemoteProxy,
			CurrentIndex:  RemoteIndexGet(),
			OnBoundsChanged: func() {
				if len(RemoteList()) == 0 {
					consoleMode.SetCurrentIndex(0)
					ModeOptionsSet(0)
					consoleMode.SetEnabled(false)
				} else {
					consoleMode.SetEnabled(true)
				}
			},
			OnCurrentIndexChanged: func() {
				if len(RemoteList()) == 0 {
					consoleMode.SetCurrentIndex(0)
					ModeOptionsSet(0)
					consoleMode.SetEnabled(false)
				} else {
					consoleMode.SetEnabled(true)
				}

				consoleRemoteProxy.SetEnabled(false)
				RemoteIndexSet(consoleRemoteProxy.Text())
				go func() {
					err := RemoteForwardUpdate()
					if err != nil {
						ErrorBoxAction(mainWindow, err.Error())
					}
					consoleRemoteProxy.SetEnabled(true)
				}()
			},
			Model: RemoteOptions(),
		},
	}
}

func ButtonWight() []Widget {
	var start *walk.PushButton
	var stop *walk.PushButton

	mutex := new(sync.Mutex)

	if AutoRunningGet() {
		go func() {
			for  {
				if start != nil && start.Visible() {
					break
				}
				time.Sleep(100 * time.Millisecond)
			}

			mutex.Lock()
			start.SetEnabled(false)
			ConsoleEnable(false)
			go func() {
				err := ServerStart()
				if err != nil {
					ErrorBoxAction(mainWindow, err.Error())
					start.SetEnabled(true)
					ConsoleEnable(true)
				} else {
					StatRunningStatus(true)
					stop.SetEnabled(true)
				}
				mutex.Unlock()
			}()
		}()
	}

	return []Widget{
		PushButton{
			AssignTo:  &start,
			Text:      LangValue("start"),
			OnClicked: func() {
				mutex.Lock()

				start.SetEnabled(false)
				ConsoleEnable(false)
				go func() {
					err := ServerStart()
					if err != nil {
						ErrorBoxAction(mainWindow, err.Error())
						start.SetEnabled(true)
						ConsoleEnable(true)
					} else {
						StatRunningStatus(true)
						stop.SetEnabled(true)
					}
					mutex.Unlock()
				}()
			},
		},
		PushButton{
			AssignTo:  &stop,
			Enabled:   false,
			Text:      LangValue("stop"),
			OnClicked: func() {
				mutex.Lock()

				stop.SetEnabled(false)
				go func() {
					err := ServerShutdown()
					if err != nil {
						ErrorBoxAction(mainWindow, err.Error())
						stop.SetEnabled(true)
					} else {
						StatRunningStatus(false)
						start.SetEnabled(true)
						ConsoleEnable(true)
					}
					ConsoleRemoteUpdate()

					mutex.Unlock()
				}()
			},
		},
	}
}
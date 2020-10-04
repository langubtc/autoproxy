package main

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type RemoteItem struct {
	Name     string
	Address  string
	Protocal string
	Auth     bool
	User     string
	Password string
}

func TestUrlGet() string {
	url := DataStringValueGet("remotetest")
	if url == "" {
		url = "https://google.com"
	}
	return url
}

func TestUrlSet(url string)  {
	DataStringValueSet("remotetest", url)
}

var remoteCache []RemoteItem

func remoteGet() []RemoteItem {
	if remoteCache == nil {
		list := make([]RemoteItem, 0)
		value := DataStringValueGet("remotelist")
		if value != "" {
			err := json.Unmarshal([]byte(value), &list)
			if err != nil {
				logs.Error("json marshal fail",err.Error())
			}
		}
		remoteCache = list
	}
	return remoteCache
}

func remoteSync()  {
	value, err := json.Marshal(remoteCache)
	if err != nil {
		logs.Error("json marshal fail",err.Error())
	} else {
		DataStringValueSet("remotelist", string(value))
	}
}

func RemoteList() []string {
	var output []string
	list := remoteGet()
	for _, v := range list {
		output = append(output, v.Name)
	}
	if len(output) == 0 {
		output = append(output, "")
	}
	return output
}

func RemoteFind(name string) RemoteItem {
	list := remoteGet()
	for _, v := range list {
		if v.Name == name {
			return v
		}
	}
	return RemoteItem{
		Name: name, Protocal: "HTTPS",
	}
}

func RemoteGet() RemoteItem {
	list := remoteGet()
	if len(list) > 0 {
		return list[0]
	}
	return RemoteItem{
		Protocal: "HTTPS",
	}
}

func RemoteDelete(name string)  {

}

func RemoteUpdate(item RemoteItem) {
	defer remoteSync()
	for i, v := range remoteCache {
		if v.Name == item.Name {
			remoteCache[i] = item
			return
		}
	}
	remoteCache = append(remoteCache, item)
}

var curRemoteItem RemoteItem

func remoteWidget() []Widget {
	var remote, protocal *walk.ComboBox
	var auth *walk.RadioButton
	var user, passwd, address, testurl *walk.LineEdit

	curRemoteItem = RemoteGet()

	updateHandler := func() {
		protocal.SetText(curRemoteItem.Protocal)
		address.SetText(curRemoteItem.Address)
		auth.SetChecked(curRemoteItem.Auth)
		user.SetEnabled(curRemoteItem.Auth)
		passwd.SetEnabled(curRemoteItem.Auth)
		user.SetText(curRemoteItem.User)
		passwd.SetText(curRemoteItem.Password)
	}

	return []Widget{
		Label{
			Text: LangValue("remoteproxy") + ":",
		},
		ComboBox{
			AssignTo: &remote,
			Editable: true,
			CurrentIndex:  0,
			Model:         RemoteList(),
			OnCurrentIndexChanged: func() {
				curRemoteItem = RemoteFind(remote.Text())
				updateHandler()
			},
			OnEditingFinished: func() {
				curRemoteItem = RemoteFind(remote.Text())
				updateHandler()
			},
		},

		Label{
			Text: LangValue("remoteaddress") + ":",
		},

		LineEdit{
			AssignTo: &address,
			Text: curRemoteItem.Address,
			OnEditingFinished: func() {
				curRemoteItem.Address = address.Text()
			},
		},

		Label{
			Text: LangValue("protocal") + ":",
		},
		ComboBox{
			AssignTo: &protocal,
			Model: ProtocalOptions(),
			Value: curRemoteItem.Protocal,
			OnCurrentIndexChanged: func() {
				curRemoteItem.Protocal = protocal.Text()
			},
		},

		Label{
			Text: LangValue("whetherauth") + ":",
		},
		RadioButton{
			AssignTo: &auth,
			OnBoundsChanged: func() {
				auth.SetChecked(curRemoteItem.Auth)
			},
			OnClicked: func() {
				auth.SetChecked(!curRemoteItem.Auth)
				curRemoteItem.Auth = !curRemoteItem.Auth

				user.SetEnabled(curRemoteItem.Auth)
				passwd.SetEnabled(curRemoteItem.Auth)
			},
		},

		Label{
			Text: LangValue("user") + ":",
		},

		LineEdit{
			AssignTo: &user,
			Text: curRemoteItem.User,
			Enabled: curRemoteItem.Auth,
			OnEditingFinished: func() {
				curRemoteItem.User = user.Text()
			},
		},

		Label{
			Text: LangValue("password") + ":",
		},

		LineEdit{
			AssignTo: &passwd,
			Text: curRemoteItem.Password,
			Enabled: curRemoteItem.Auth,
			OnEditingFinished: func() {
				curRemoteItem.Password = passwd.Text()
			},
		},

		PushButton{
			Text: LangValue("test"),
			OnClicked: func() {

			},
		},

		LineEdit{
			AssignTo: &testurl,
			Text: TestUrlGet(),
			OnEditingFinished: func() {
				TestUrlSet(testurl.Text())
			},
		},
	}
}

func RemoteServer()  {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	_, err := Dialog{
		AssignTo: &dlg,
		Title: LangValue("remoteproxy"),
		Icon: walk.IconShield(),
		DefaultButton: &acceptPB,
		CancelButton: &cancelPB,
		Size: Size{250, 300},
		MinSize: Size{250, 300},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: remoteWidget(),
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &acceptPB,
						Text:     LangValue("save"),
						OnClicked: func() {
							if curRemoteItem.Auth {
								if curRemoteItem.User == "" || curRemoteItem.Password == "" {
									ErrorBoxAction(dlg, LangValue("inputuserandpasswd"))
									return
								}
							}
							if curRemoteItem.Name == "" || curRemoteItem.Address == "" {
								ErrorBoxAction(dlg, LangValue("inputnameandaddress"))
								return
							}
							RemoteUpdate(curRemoteItem)
							ConsoleRemoteUpdate()
							dlg.Accept()
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
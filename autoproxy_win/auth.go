package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io/ioutil"
	"sync"
)

type AuthInfo struct {
	User     string
	Password string
}

type AuthCtrl struct {
	sync.RWMutex
	Items []AuthInfo
	cache map[string]*AuthInfo
}

var authCtrl *AuthCtrl

func AuthInit() error {
	authCtrl = new(AuthCtrl)
	authCtrl.Items = make([]AuthInfo, 0)
	authCtrl.cache = make(map[string]*AuthInfo, 1024)

	body, err := ioutil.ReadFile(DEFAULT_HOME + "\\auth.json")
	if err != nil {
		logs.Error("no auth json fail")
		return nil
	}
	err = json.Unmarshal(body, &authCtrl.Items)
	if err != nil {
		return err
	}
	for _, v := range authCtrl.Items {
		authCtrl.cache[v.User] = &AuthInfo{
			User: v.User, Password: v.Password,
		}
	}
	return nil
}

func authSync() error {
	body, err := json.Marshal(authCtrl.Items)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	return SaveToFile(DEFAULT_HOME + "\\auth.json", body)
}

func authAdd(user string, passwd string) error {
	authCtrl.Lock()
	defer authCtrl.Unlock()

	_, flag := authCtrl.cache[user]
	if flag == true {
		return fmt.Errorf(LangValue("userexist"))
	}

	authCtrl.cache[user] = &AuthInfo{
		User: user, Password: passwd,
	}

	authCtrl.Items = append(authCtrl.Items, AuthInfo{
		User: user, Password: passwd,
	})

	return authSync()
}

func authDelete(user string) error {
	authCtrl.Lock()
	defer authCtrl.Unlock()

	_, flag := authCtrl.cache[user]
	if flag == false {
		return fmt.Errorf(LangValue("usernotexist"))
	}
	delete(authCtrl.cache, user)

	for i, v:= range authCtrl.Items {
		if v.User == user {
			authCtrl.Items = append(authCtrl.Items[:i], authCtrl.Items[i+1:]...)
			break
		}
	}

	return authSync()
}

func AuthCheck(user string, passwd string) bool {
	authCtrl.RLock()
	defer authCtrl.RUnlock()

	userInfo, _ := authCtrl.cache[user]
	if userInfo == nil {
		return false
	}
	if userInfo.Password != passwd {
		return false
	}
	return true
}

func AuthGet() []AuthInfo {
	return authCtrl.Items
}

func AuthAdd()  {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var user, passwd *walk.LineEdit

	_, err := Dialog{
		AssignTo: &dlg,
		Title:  LangValue("addcred"),
		Icon: walk.IconShield(),
		DefaultButton: &acceptPB,
		CancelButton: &cancelPB,
		Size: Size{300, 150},
		MinSize: Size{300, 150},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 3},
				Children: []Widget{
					Label{
						Text: LangValue("user") + ":",
					},
					LineEdit{
						AssignTo: &user,
						Text: "",
					},
					PushButton{
						Text:      LangValue("randomgen"),
						OnClicked: func() {
							user.SetText("U"+GetUser(5))
						},
					},
					Label{
						Text: LangValue("password") + ":",
					},
					LineEdit{
						AssignTo: &passwd,
						Text: "",
					},
					PushButton{
						Text:      LangValue("randomgen"),
						OnClicked: func() {
							passwd.SetText(GetToken(16))
						},
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &acceptPB,
						Text:     LangValue("add"),
						OnClicked: func() {
							if user.Text() == "" || passwd.Text() == "" {
								ErrorBoxAction(dlg, LangValue("inputuserandpasswd"))
								return
							}
							err := authAdd(user.Text(), passwd.Text())
							if err != nil {
								ErrorBoxAction(dlg, err.Error())
								return
							}
							InfoBoxAction(dlg, LangValue("addsuccess"))
							dlg.Cancel()
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


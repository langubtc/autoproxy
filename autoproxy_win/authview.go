package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sort"
)


type UserTable struct {
	Index   int
	User    string
	Passwd  string

	checked bool
}

type UserModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*UserTable
}

func (n *UserModel)RowCount() int {
	return len(n.items)
}


func (n *UserModel)Value(row, col int) interface{} {
	item := n.items[row]
	switch col {
	case 0:
		return item.Index
	case 1:
		return item.User
	case 2:
		return item.Passwd
	}
	panic("unexpected col")
}

func (n *UserModel) Checked(row int) bool {
	return n.items[row].checked
}

func (n *UserModel) SetChecked(row int, checked bool) error {
	n.items[row].checked = checked
	return nil
}

func (m *UserModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order
	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]
		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}
			return !ls
		}
		switch m.sortColumn {
		case 0:
			return c(a.Index < b.Index)
		case 1:
			return c(a.User < b.User)
		case 2:
			return c(a.Passwd < b.Passwd)
		case 3:
		}
		panic("unreachable")
	})
	return m.SorterBase.Sort(col, order)
}

var userList *UserModel

func UserListUpdate()  {
	authCtrl.RLock()
	defer authCtrl.RUnlock()

	item := make([]*UserTable, 0)
	for idx, v := range authCtrl.Items {
		item = append(item, &UserTable{
			Index: idx,
			User: v.User,
			Passwd: v.Password,
		})
	}

	userList.items = item
	userList.PublishRowsReset()
	userList.Sort(userList.sortColumn, userList.sortOrder)
}

func init()  {
	userList = new(UserModel)
	userList.items = make([]*UserTable, 0)
}

func AuthDelete(from walk.Form) error {
	var tempList []string
	for _,v := range userList.items {
		if v.checked {
			tempList = append(tempList, v.User)
		}
	}
	if len(tempList) == 0 {
		ErrorBoxAction(from, LangValue("nochoice"))
		return fmt.Errorf("no user select")
	}
	ConfirmBoxAction(from, fmt.Sprintf(LangValue("confirmdelete")+ "%v", tempList))

	for _, v := range tempList {
		authDelete(v)
	}

	InfoBoxAction(from, LangValue("deletesuccess"))
	return nil
}

func AuthView()  {
	UserListUpdate()

	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	_, err := Dialog{
		AssignTo: &dlg,
		Title: LangValue("viewcred"),
		Icon: walk.IconShield(),
		DefaultButton: &acceptPB,
		CancelButton: &cancelPB,
		Size: Size{400, 200},
		MinSize: Size{400, 200},
		Layout:  VBox{},
		Children: []Widget{
			TableView{
				AlternatingRowBG: true,
				ColumnsOrderable: true,
				CheckBoxes: true,
				Columns: []TableViewColumn{
					{Title: "#", Width: 40},
					{Title: LangValue("user"), Width: 80},
					{Title: LangValue("password"), Width: 150},
				},
				StyleCell: func(style *walk.CellStyle) {
					if style.Row()%2 == 0 {
						style.BackgroundColor = walk.RGB(248, 248, 255)
					} else {
						style.BackgroundColor = walk.RGB(220, 220, 220)
					}
				},
				Model:userList,
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &acceptPB,
						Text:     LangValue("delete"),
						OnClicked: func() {
							err := AuthDelete(dlg)
							if err != nil {
								logs.Error(err.Error())
							} else {
								UserListUpdate()
							}
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


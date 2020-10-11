package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"sync"
)

type DomainItem struct {
	Index   int
	Domain  string

	checked bool
}

type DomainModel struct {
	sync.RWMutex

	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder

	items      []*DomainItem
}

func (n *DomainModel)RowCount() int {
	return len(n.items)
}

func (n *DomainModel)Value(row, col int) interface{} {
	item := n.items[row]
	switch col {
	case 0:
		return item.Index
	case 1:
		return item.Domain
	}
	panic("unexpected col")
}

func (n *DomainModel) Checked(row int) bool {
	return n.items[row].checked
}

func (n *DomainModel) SetChecked(row int, checked bool) error {
	n.items[row].checked = checked
	return nil
}

func (m *DomainModel) Sort(col int, order walk.SortOrder) error {
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
			return c(a.Domain < b.Domain)
		}
		panic("unreachable")
	})
	return m.SorterBase.Sort(col, order)
}

var domainTable *DomainModel
var domainList  []string

func DomainSave(list []string) error {
	sort.Strings(list)
	file := fmt.Sprintf("%s/domain.json", DEFAULT_HOME)
	value, err := json.Marshal(list)
	if err != nil {
		logs.Error("json marshal domain list fail, %s", err.Error())
		return err
	}
	return SaveToFile(file, value)
}

func DomainAdd(domain string) error {
	for _, v := range domainList {
		if v == domain {
			return fmt.Errorf("domain %s exist", domain)
		}
	}
	domainList = append(domainList, domain)
	return DomainSave(domainList)
}

func DomainList() []string {
	return domainList
}

func DomainTableUpdate(find string)  {
	item := make([]*DomainItem, 0)
	for idx, v := range domainList {
		if strings.Index(v, find) == -1 {
			continue
		}
		item = append(item, &DomainItem{
			Index: idx, Domain: v,
		})
	}
	domainTable.items = item
	domainTable.PublishRowsReset()
	domainTable.Sort(domainTable.sortColumn, domainTable.sortOrder)
}

func DomainInit() error {
	domainTable = new(DomainModel)
	domainTable.items = make([]*DomainItem, 0)

	domainFile := fmt.Sprintf("%s/domain.json", DEFAULT_HOME)
	_, err := os.Stat(domainFile)
	if err != nil {
		value, err := BoxFile().Bytes("default_domain.json")
		if err != nil {
			logs.Error("open default domain json file fail, %s", err.Error())
			return err
		}
		err = SaveToFile(domainFile, value)
		if err != nil {
			logs.Error("default domain json save to app data dir fail, %s", err.Error())
			return err
		}
	}

	value, err := ioutil.ReadFile(domainFile)
	if err != nil {
		logs.Error("read domain json file from app data dir fail, %s", err.Error())
		return err
	}

	var output []string
	err = json.Unmarshal(value, &output)
	if err != nil {
		logs.Error("json unmarshal domain json fail, %s", err.Error())
		return err
	}

	domainList = output
	return nil
}

func DomainDelete(owner *walk.Dialog) error {
	var deleteList []string
	for _, v := range domainTable.items {
		if v.checked {
			deleteList = append(deleteList, v.Domain)
		}
	}
	if len(deleteList) == 0 {
		return fmt.Errorf(LangValue("nochoiceobject"))
	}

	var remanderList []string
	for _, v := range domainList {
		var exist bool
		for _, v2 := range deleteList {
			if v == v2 {
				exist = true
				break
			}
		}
		if !exist {
			remanderList = append(remanderList, v)
		}
	}

	domainList = remanderList
	DomainSave(remanderList)

	InfoBoxAction(owner, fmt.Sprintf("%v %s", deleteList, LangValue("deletesuccess")))

	return nil
}

func RemodeEdit()  {
	DomainTableUpdate("")

	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var findPB, addPB *walk.PushButton
	var addLine, findLine *walk.LineEdit

	_, err := Dialog{
		AssignTo: &dlg,
		Title: LangValue("forwarddomain"),
		Icon: walk.IconShield(),
		DefaultButton: &acceptPB,
		CancelButton: &cancelPB,
		Size: Size{300, 450},
		MinSize: Size{300, 300},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 3, MarginsZero: true},
				Children: []Widget{
					Label{
						Text: LangValue("domain") + ":",
					},
					LineEdit {
						AssignTo: &addLine,
						Text: "",
					},
					PushButton{
						AssignTo: &addPB,
						Text:     LangValue("add"),
						OnClicked: func() {
							addDomain := addLine.Text()

							if addDomain == "" {
								ErrorBoxAction(dlg, LangValue("inputdomain"))
								return
							}
							err := DomainAdd(addDomain)
							if err != nil {
								ErrorBoxAction(dlg, err.Error())
								return
							}

							go func() {
								InfoBoxAction(dlg, addDomain + " " + LangValue("addsuccess") )
							}()

							addLine.SetText("")
							findLine.SetText("")
							DomainTableUpdate("")
							RouteUpdate()
						},
					},
					Label{
						Text: LangValue("findkey") + ":",
					},
					LineEdit {
						AssignTo: &findLine,
						Text: "",
					},
					PushButton{
						AssignTo: &findPB,
						Text:     LangValue("find"),
						OnClicked: func() {
							DomainTableUpdate(findLine.Text())
						},
					},
				},
			},
			TableView{
				AlternatingRowBG: true,
				ColumnsOrderable: true,
				CheckBoxes: true,
				Columns: []TableViewColumn{
					{Title: "#", Width: 60},
					{Title: LangValue("domain"), Width: 160},
				},
				StyleCell: func(style *walk.CellStyle) {
					if style.Row()%2 == 0 {
						style.BackgroundColor = walk.RGB(248, 248, 255)
					} else {
						style.BackgroundColor = walk.RGB(220, 220, 220)
					}
				},
				Model:domainTable,
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &acceptPB,
						Text:     LangValue("delete"),
						OnClicked: func() {
							err := DomainDelete(dlg)
							if err != nil {
								logs.Error(err.Error())
								ErrorBoxAction(dlg, err.Error())
							} else {
								DomainTableUpdate(findLine.Text())
								RouteUpdate()
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


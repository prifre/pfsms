package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/prifre/pfsms/db"
)

type AppTable struct {
	// Theme holds the current theme
	Theme string
}

type thetable struct {
	tableShowCustomers 	*widget.Table
	tableShowGroups 	*widget.Table
	appTable  	*AppTable
	window      fyne.Window
	app         fyne.App
}

func NewTable(a fyne.App, w fyne.Window,  at *AppTable) *thetable {
	return &thetable{app: a, window: w,  appTable: at}
}
func (s *thetable) buildTableCustomers() *container.Scroll {
//	var data = [][]string{{"A1", "B1"},{"A2", "B2"},{"A3", "B3"},{"A4", "B4"},{"A5", "B5"}}
	d:=new(db.DBtype)
	d.Setupdb()
	dataCustomers,err:=d.ShowCustomers(0,10000)
	if err!=nil {
		fmt.Printf("ShowCustomer failed %s",err.Error())
	}
	listCustomers := widget.NewTable(
	func() (int, int) {
		return len(dataCustomers), len(dataCustomers[0])
	},
	func() fyne.CanvasObject {
		return widget.NewLabel("wide content")
	},
	func(i widget.TableCellID, o fyne.CanvasObject) {
		o.(*widget.Label).SetText(dataCustomers[i.Row][i.Col])
	})
	listCustomers.OnSelected=func(i widget.TableCellID) {
		fmt.Println(i)
	}
	s.tableShowCustomers = listCustomers
	return container.NewScroll(listCustomers)
}
func (s *thetable) buildTableGroups() *container.Scroll {
	//	var data = [][]string{{"A1", "B1"},{"A2", "B2"},{"A3", "B3"},{"A4", "B4"},{"A5", "B5"}}
	d:=new(db.DBtype)
	d.Setupdb()
		dataGroups,err:=d.ShowGroupnames()
	if err!=nil {
		fmt.Printf("ShowGroups failed %s",err.Error())
	}
	listGroups := widget.NewTable(
	func() (int, int) {
		return len(dataGroups), len(dataGroups[0])
	},
	func() fyne.CanvasObject {
		return widget.NewLabel("wide content")
	},
	func(i widget.TableCellID, o fyne.CanvasObject) {
		o.(*widget.Label).SetText(dataGroups[i.Row][i.Col])
	})
	listGroups.OnSelected=func(i widget.TableCellID) {
		fmt.Println(i)
	}
	s.tableShowGroups = listGroups
	return container.NewScroll(listGroups)
}
func (s *thetable) buildTable() *container.Scroll {
	bigContainer:=container.NewScroll(container.NewVBox(s.buildTableCustomers(),s.buildTableGroups(),
	//		&widget.Card{Title: "Data Handling", Content: list},
widget.NewLabel("asdf")))
return bigContainer
}

func (s *thetable) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Table", Icon: theme.GridIcon(), Content: s.buildTable()}
}

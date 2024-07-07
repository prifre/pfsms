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
	tableShow 	*widget.Table
	appTable  	*AppTable
	window      fyne.Window
	app         fyne.App
}

func NewTable(a fyne.App, w fyne.Window,  at *AppTable) *thetable {
	return &thetable{app: a, window: w,  appTable: at}
}
func (s *thetable) buildTable() *container.Scroll {
//	var data = [][]string{{"A1", "B1"},{"A2", "B2"},{"A3", "B3"},{"A4", "B4"},{"A5", "B5"}}
	d:=new(db.DBtype)
	d.Setupdb()
	data,err:=d.ShowCustomers(0,10000)
	if err!=nil {
		fmt.Printf("ShowCustomer failed %s",err.Error())
	}
	list := widget.NewTable(
	func() (int, int) {
		return len(data), len(data[0])
	},
	func() fyne.CanvasObject {
		return widget.NewLabel("wide content")
	},
	func(i widget.TableCellID, o fyne.CanvasObject) {
		o.(*widget.Label).SetText(data[i.Row][i.Col])
	})
	list.OnSelected=func(i widget.TableCellID) {
		fmt.Println(i)
	}
	s.tableShow = list
	return container.NewScroll(list,
//		&widget.Card{Title: "Data Handling", Content: list},
	)
}
func (s *thetable) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Table", Icon: theme.GridIcon(), Content: s.buildTable()}
}

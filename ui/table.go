package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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
	var data = [][]string{}
	for i:=0;i<100;i++ {
		data1:=[]string{fmt.Sprint(i),fmt.Sprintf("A_%d",i)}
		data=append(data,data1)
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
	s.tableShow = list
	return container.NewScroll(list,
//		&widget.Card{Title: "Data Handling", Content: list},
	)
}
func (s *thetable) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Table", Icon: theme.SettingsIcon(), Content: s.buildTable()}
}
	
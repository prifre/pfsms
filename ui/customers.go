package ui

// show database with Customers & Grops
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
func (s *thetable) listCustomers() *widget.Table {
	d:=new(db.DBtype)
	dataCustomers,err:=d.ShowCustomers(0,10000)
	if err!=nil {
		fmt.Printf("ShowCustomer failed %s",err.Error())
	}
	if len(dataCustomers)<=0 {
		dataCustomers=[][]string{{"No customers"}}
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
	listCustomers.SetColumnWidth(0,20)
	listCustomers.SetRowHeight(0,20)
	listCustomers.BaseWidget.Resize(fyne.NewSize(1000,1000))
	return listCustomers
}
func (s *thetable) buildTableCustomers() *container.Scroll {
//	var data = [][]string{{"A1", "B1"},{"A2", "B2"},{"A3", "B3"},{"A4", "B4"},{"A5", "B5"}}
	s.tableShowCustomers = s.listCustomers()
	return container.NewScroll(s.tableShowCustomers)
}
func (s *thetable) buildTableGroups() *container.Scroll {
	//	var data = [][]string{{"A1", "B1"},{"A2", "B2"},{"A3", "B3"},{"A4", "B4"},{"A5", "B5"}}
	d:=new(db.DBtype)
	d.Opendb()
	dataGroups,err:=d.ShowGroupnames()
	if err!=nil {
		fmt.Printf("ShowGroups failed %s",err.Error())
	}
	if len(dataGroups)<=0 {
		dataGroups=[][]string{{"No groups"}}
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
		fmt.Println(i,s.window.Canvas().Size().Width)
	}
	s.tableShowGroups = listGroups
	return container.NewScroll(listGroups)
}

func (s *thetable) buildTable() *container.Scroll {
	gr:=container.NewGridWithColumns(2,
			s.buildTableCustomers().Content,
			s.buildTableGroups())
	var windowSize= fyne.NewSize(700,600)
	s.buildTableCustomers().SetMinSize(fyne.NewSize(windowSize.Width*.9,windowSize.Height*.9))
	s.buildTableGroups().SetMinSize(fyne.NewSize(windowSize.Width*.2,windowSize.Height*.9))
	s.window.Resize(windowSize)
	bigContainer:=container.NewScroll(gr)
return bigContainer
}

func (s *thetable) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Customers", Icon: theme.StorageIcon(), Content: s.buildTable()}
}

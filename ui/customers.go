package ui

// show database with Customers & Grops
import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/prifre/pfsms/pfdatabase"
)

type AppTable struct {
	// Theme holds the current theme
	Theme string
}

type thetable struct {
	tableCustomers 		*widget.Table
	tableGroups 		*widget.Table
	dataCustomers		[][]string
	dataGroups			[]string
	dataAllCustomers	[][]string
	dataAllGroups		[][]string
	appTable  			*AppTable
	window      		fyne.Window
	app         		fyne.App
}

func NewTable(a fyne.App, w fyne.Window,  at *AppTable) *thetable {
	return &thetable{app: a, window: w,  appTable: at}
}
func (s *thetable) buildTableCustomers() *container.Scroll {
	if s.dataAllCustomers==nil {
		s.dataCustomers=new(pfdatabase.DBtype).ShowCustomers()
		s.dataAllCustomers=s.dataCustomers
	}
	s.tableCustomers = widget.NewTable(nil,nil,nil)
	s.tableCustomers.Length = func() (int, int) {	
		if len(s.dataCustomers)<=0 {
			s.dataCustomers = [][]string{{"No customers",""}}
		}
		return len(s.dataCustomers),len(s.dataCustomers[0])	
	}
	s.tableCustomers.CreateCell = func() fyne.CanvasObject {
		return widget.NewLabelWithStyle("123456789012345678901234567890",fyne.TextAlignCenter,fyne.TextStyle{Monospace: true})
	}
	s.tableCustomers.UpdateCell = func(i widget.TableCellID, o fyne.CanvasObject) {
		o.(*widget.Label).SetText(strings.TrimSpace(s.dataCustomers[i.Row][i.Col]))
		o.(*widget.Label).Refresh()
	}
	// s.tableCustomers.OnSelected=func(i widget.TableCellID) {
	// 	fmt.Println("Currently selected Customer: ",i)
	// }
	// s.tableCustomers.MouseIn=func(i widget.TableCellID) {
	// 	sq :=fmt.Sprintf("SELECT CONCAT(phone,' ',firstname,' ',lastname) AS data FROM tblCustomers WHERE ORDER BY phone OFFSET %d LIMIT 1 ASC",i.Row)
	// 	o.(*widget.Label).SetText(GetCustomers(sq))
	// }
	// s.window.CenterOnScreen()
	// s.window.Canvas().Content().Resize()
	// 	// Update the content size based on the new window size
	// 	fmt.Println()
	// }
	s.tableCustomers.SetColumnWidth(0,s.window.Canvas().Content().Size().Width*0.16)
	s.tableCustomers.SetColumnWidth(1,s.window.Canvas().Content().Size().Width*0.17)
	s.tableCustomers.SetColumnWidth(2,s.window.Canvas().Content().Size().Width*0.12)
	return container.NewScroll(s.tableCustomers)
}
func (s *thetable) buildTableGroups() *container.Scroll {
	//	var data = [][]string{{"A1", "B1"},{"A2", "B2"},{"A3", "B3"},{"A4", "B4"},{"A5", "B5"}}
	d:=new(pfdatabase.DBtype)
	s.dataAllGroups =d.ShowAllGroups()
	s.dataAllGroups = append(s.dataAllGroups,[]string{"All customers...",""})
	s.dataGroups =new(pfdatabase.DBtype).ShowGroups()
	s.dataGroups = append(s.dataGroups,"All customers...")
	s.tableGroups = widget.NewTable(nil,nil,nil)
	s.tableGroups.Length = func() (int, int) {	
		return len(s.dataGroups), 1
	}
	s.tableGroups.CreateCell = func() fyne.CanvasObject {	
		return widget.NewLabel("wide content")	
	}
	s.tableGroups.UpdateCell = func(i widget.TableCellID, o fyne.CanvasObject) {
		o.(*widget.Label).SetText(s.dataGroups[i.Row])
	}
	s.tableGroups.OnSelected=func(i widget.TableCellID) {
		var tc [][]string = nil
		for j:=0;j<len(s.dataAllGroups);j++ {
			if s.dataGroups[i.Row] == s.dataAllGroups[j][0] {
				for k:=0;k<len(s.dataAllCustomers);k++ {
					gp:=s.dataAllGroups[j][1]
					cp:=s.dataAllCustomers[k][0]
					if gp == cp {
						tc=append(tc,s.dataAllCustomers[k])
					}
				}
			}
		}
		if tc==nil {
			s.dataCustomers=s.dataAllCustomers
		} else {
			s.dataCustomers=tc
		}	
		s.buildTableCustomers()
		s.tableCustomers.Refresh()
	}
	s.tableGroups.SetColumnWidth(0,s.window.Canvas().Content().Size().Width*0.2)
	return container.NewScroll(s.tableGroups)
}
func (s *thetable) tabItem() *container.TabItem {
	// s.tableGroups.SetColumnWidth(0,512)
	c:= &container.TabItem{Text: "Customers", Icon: theme.StorageIcon(), Content: 
	container.NewScroll(container.NewGridWithColumns(2,
			s.buildTableCustomers(),
			s.buildTableGroups())),
	}
	c.Content.Refresh()
	return c
}

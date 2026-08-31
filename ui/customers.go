package ui

import (
	"slices"
	"strings"

	"pfsms/pfdatabase"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type thetable struct {
	tableCustomers   *widget.Table
	tableGroups      *widget.Table
	dataCustomers    [][]string
	dataGroups       []string
	dataAllCustomers [][]string
	dataAllGroups    [][]string
	window           fyne.Window
	db               *pfdatabase.DBtype
}

func NewCustomers(w fyne.Window) *thetable {
	return &thetable{
		window: w,
		db:     new(pfdatabase.DBtype),
	}
}

func (s *thetable) buildCustomers() fyne.CanvasObject {
	// Grid med 2 kolumner. Tabellerna hanterar sin egen scrollning.
	return container.NewGridWithColumns(2,
		s.buildTableCustomers(),
		s.buildTableGroups(),
	)
}

func (s *thetable) buildTableCustomers() fyne.CanvasObject {
	s.dataCustomers = s.db.ShowCustomers()
	s.dataAllCustomers = s.dataCustomers

	s.tableCustomers = widget.NewTable(
		func() (int, int) {
			if len(s.dataCustomers) == 0 {
				return 0, 0
			}
			return len(s.dataCustomers), len(s.dataCustomers[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			lbl := o.(*widget.Label)
			// Säkra mot index out of bounds om datat ändrats under tiden
			if i.Row < len(s.dataCustomers) && i.Col < len(s.dataCustomers[i.Row]) {
				lbl.SetText(strings.TrimSpace(s.dataCustomers[i.Row][i.Col]))
			} else {
				lbl.SetText("")
			}
		},
	)

	// Sätt rimliga fasta bredder för kolumnerna
	s.tableCustomers.SetColumnWidth(0, 180)
	s.tableCustomers.SetColumnWidth(1, 180)
	s.tableCustomers.SetColumnWidth(2, 120)

	return s.tableCustomers
}

func (s *thetable) buildTableGroups() fyne.CanvasObject {
	s.reloadGroupData()

	s.tableGroups = widget.NewTable(
		func() (int, int) {
			return len(s.dataGroups), 1
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			lbl := o.(*widget.Label)
			if i.Row < len(s.dataGroups) {
				lbl.SetText(s.dataGroups[i.Row])
			} else {
				lbl.SetText("")
			}
		},
	)

	s.tableGroups.OnSelected = func(i widget.TableCellID) {
		if i.Row >= len(s.dataGroups) {
			return
		}

		var tc [][]string
		var ingroup []string

		selectedGroup := s.dataGroups[i.Row]

		if selectedGroup == "All customers..." {
			s.dataCustomers = s.dataAllCustomers
		} else {
			for j := 0; j < len(s.dataAllGroups); j++ {
				if selectedGroup == s.dataAllGroups[j][0] {
					gp := s.dataAllGroups[j][1]
					for k := 0; k < len(s.dataAllCustomers); k++ {
						cp := s.dataAllCustomers[k][0]
						if gp == cp && slices.Index(ingroup, cp) == -1 {
							ingroup = append(ingroup, cp)
							tc = append(tc, s.dataAllCustomers[k])
						}
					}
				}
			}
			s.dataCustomers = tc
		}

		// Rensa markering och uppdatera tabellen
		s.tableCustomers.UnselectAll()
		s.tableCustomers.Refresh()
	}

	s.tableGroups.SetColumnWidth(0, 250)

	return s.tableGroups
}

// Hjälpmetod för att ladda om grupper på ett ställe så att inte "All customers..." tappas bort
func (s *thetable) reloadGroupData() {
	s.dataAllGroups = s.db.ShowAllGroups()
	s.dataAllGroups = append(s.dataAllGroups, []string{"All customers...", ""})

	s.dataGroups = s.db.ShowGroups()
	s.dataGroups = append(s.dataGroups, "All customers...")
}

func (s *thetable) TabItem() *container.TabItem {
	return &container.TabItem{
		Text:    "Customers",
		Icon:    theme.StorageIcon(),
		Content: s.buildCustomers(),
	}
}

func (s *thetable) RefreshData() {
	s.dataAllCustomers = s.db.ShowCustomers()
	s.dataCustomers = s.dataAllCustomers

	s.reloadGroupData()

	if s.tableCustomers != nil {
		s.tableCustomers.Refresh()
	}
	if s.tableGroups != nil {
		s.tableGroups.Refresh()
	}
}

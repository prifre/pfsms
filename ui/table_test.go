package ui

import (
	"testing"
)

func TestTable(t *testing.T) {
		a := app.NewWithID("pfsms")
		w := a.NewWindow("pfsms-gui")
		w.SetContent(ui.Create(a, w))
		w.Resize(fyne.NewSize(700, 600))
		w.SetMaster()
		w.ShowAndRun()
	}
}


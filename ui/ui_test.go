package ui

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func TestGetmemoryinfo(t *testing.T) {
	fmt.Println(Getmemoryinfo())
}

func TestTable(t *testing.T) {
		a := app.NewWithID("pfsms")
		w := a.NewWindow("pfsms-gui")
		appTable := &AppTable{}
		n:=NewTable(a,w,appTable).tabItem()
		c:= &container.AppTabs{Items: []*container.TabItem{		n,		}}
		w.SetContent(c)
		w.Resize(fyne.NewSize(700, 600))
		w.SetMaster()
		w.ShowAndRun()
	}
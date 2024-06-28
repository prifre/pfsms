package main

//Main file

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/prifre/pfsms/ui"
)


func main() {
	a := app.NewWithID("pfsms")
	w := a.NewWindow("pfsms-gui")
	w.SetContent(ui.Create(a, w))
	w.Resize(fyne.NewSize(700, 600))
	w.SetMaster()
	w.ShowAndRun()
}

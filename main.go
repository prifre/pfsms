package main

//Main file

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

const version = "v0.0.2"

func main() {
	a := app.NewWithID("pfsms")
	w := a.NewWindow("pfsms-gui")
	w.SetContent(pfCreate(a, w))
	w.Resize(fyne.NewSize(700, 600))
	w.SetMaster()
	w.ShowAndRun()
}

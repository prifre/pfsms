package main

//Main file

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/prifre/pfsms/ui"
)

func main() {
	var wx,wy float32
	a := app.NewWithID("pfsms")
	w := a.NewWindow("pfsms")
	wx=1024
	wy=768
	w.Canvas().Content().Resize(fyne.NewSize(wx,wy))
	w.Resize(fyne.NewSize(wx,wy))
	w.SetContent(ui.Create(w))
	w.ShowAndRun()
}

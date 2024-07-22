package main

//Main file

//	"github.com/prifre/pfsms/ui"

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/prifre/pfsms/ui"
)


func main() {
	a := app.NewWithID("pfsms")
	w := a.NewWindow("pfsms-gui")
	w.Resize(fyne.NewSize(700, 600))
	w.SetContent(ui.Create(a, w))
	fmt.Printf("%f,%f",w.Canvas().Scale(),w.Canvas().Size().Width)
	w.SetContent(ui.Create(a, w))
	w.SetMaster()
	w.ShowAndRun()
	fmt.Println("SAVING SETTINGS!")
}

package main

//Main file

//	"github.com/prifre/pfsms/ui"

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"github.com/prifre/pfsms/ui"
)

func main() {
	a := app.NewWithID("pfsms")
	w := a.NewWindow("pfsms-gui")
	w.SetContent(ui.Create(a, w))
	w.SetMaster()
	w.ShowAndRun()
	fmt.Println("SAVING SETTINGS!")
}

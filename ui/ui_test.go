package ui

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func TestGetmemoryinfo(t *testing.T) {
	fmt.Println(Getmemoryinfo())
}

func TestTable(t *testing.T) {
		a := app.NewWithID("pfsms")
		w := a.NewWindow("pfsms-gui")
		n:=NewTable(w).tabItem()
		c:= &container.AppTabs{Items: []*container.TabItem{		n,		}}
		w.SetContent(c)
		w.Resize(fyne.NewSize(700, 600))
		w.SetMaster()
		w.ShowAndRun()
	}

func TestX(t *testing.T) {
	myApp := app.New()
	myWindow := myApp.NewWindow("Fyne Example")

 // Create a basic text label
 label := widget.NewLabel("Hello Fyne!")

 // Create a button with a callback function
 button := widget.NewButton("Quit", func() {
  myApp.Quit()
 })

 // Create a simple form with an entry field
 entry := widget.NewEntry()
 form := &widget.Form{
  OnSubmit: func() {
   label.SetText("Entered : " + entry.Text)
  },
 }
 form.Append("Entry : ", entry)

 // Combine widgets into a container
 content := container.NewVBox(
  label,
  form,
  button,
 )

 // Center the content on the screen
 myWindow.SetContent(container.New(layout.NewCenterLayout(), content))

 // Show the window
 myWindow.ShowAndRun()
}
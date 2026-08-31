package ui

import (
	"fmt"
	"testing"

	"pfsms/general"
	"pfsms/pfdatabase"
)

/*
OTHER STUFF:
fmt.Println(fyne.CurrentApp().Metadata().Version)
fmt.Println(fyne.CurrentApp().Metadata().Name)
fmt.Println(fyne.CurrentApp().Metadata().ID)
fmt.Println("Version:",fyne.CurrentApp().Metadata().Version)
fmt.Println("UniqueID:", fyne.CurrentApp().UniqueID())
fmt.Println("BuildType:",fyne.CurrentApp().Settings().BuildType())
fmt.Println("RootID:",fyne.CurrentApp().Storage().RootURI())
fmt.Println("HOME-DIR",GetHomeDir())
*/
func TestGetphonesforgroup(t *testing.T) {
	general.Setupfiles()
	s := new(theform).Getphonesforgroup("PETER FREUND")
	fmt.Println(s)
}
func TestSaveHistory(t *testing.T) {
	general.Setupfiles()
	var sh []string = []string{"20240818121212", "PETER FREUND", "0046736290839", "test"}
	new(pfdatabase.DBtype).SaveHistory(sh)
}

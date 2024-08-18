package ui

import (
	"fmt"
	"os/exec"
	"runtime"
	"testing"

	"fyne.io/fyne/v2"
	"github.com/prifre/pfsms/pfdatabase"
)

func TestImportCustomers(t *testing.T) {
	Setupfiles()
	new(pfdatabase.DBtype).ImportCustomers(fyne.CurrentApp().Preferences().String("customersfile"))
	new(pfdatabase.DBtype).ExportCustomers(fyne.CurrentApp().Preferences().String("customersfile"))
}
func TestImportGroups(t *testing.T) {
	Setupfiles()
	new(pfdatabase.DBtype).ImportGroups(fyne.CurrentApp().Preferences().String("groupsfile"))
	new(pfdatabase.DBtype).ExportGroups(fyne.CurrentApp().Preferences().String("groupsfile"))
}
func TestExportHistory(t *testing.T) {
	Setupfiles()
	new(pfdatabase.DBtype).ExportHistory(fyne.CurrentApp().Preferences().String("historyfile"))
}
func TestOpenLog(t *testing.T) {
	Setupfiles()
	dd := fyne.CurrentApp().Preferences().String("pfsmslog")
	// dd,_:=os.UserHomeDir()
	// dd = fmt.Sprintf("%s%c%s%c%s",dd,os.PathSeparator,"pfsms",os.PathSeparator,fn)
	if runtime.GOOS == "darwin" {
		cmd := `open "` + dd + `"`
		exec.Command("/bin/bash", "-c", cmd).Start()
	} else {
		exec.Command("explorer", dd).Start()
	}
}
func TestMydebug(t *testing.T) {
	dd := fyne.CurrentApp().Preferences().String("debug")
	fmt.Println(dd)
	dd = fyne.CurrentApp().Preferences().String("pfsmslog")
	fmt.Println(dd)
}
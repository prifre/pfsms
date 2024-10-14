package ui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/prifre/pfsms/pfdatabase"
)

var mobilemodels = []string{"Samsung S24", "Samsung S9"}

type pfsettings struct {
	mobileNumber  *widget.Entry
	mobileCountry *widget.Select
	mobileModel   *widget.Select
	mobilePort    *widget.Entry
	mobileAddhash *widget.Check
	btnTestSMS    *widget.Button

	btnOpenDatadir     *widget.Button
	customersfile      *widget.Label
	btnImportCustomers *widget.Button
	btnExportCustomers *widget.Button
	btnDeleteCustomers *widget.Button
	groupsfile         *widget.Label
	btnImportGroups    *widget.Button
	btnExportGroups    *widget.Button
	btnDeleteGroups    *widget.Button
	logfile            *widget.Label
	historyfile        *widget.Label
	btnExportHistory   *widget.Button
	btnOpenLog         *widget.Button
	logtext            *widget.Entry

	window fyne.Window
}

func NewSettings(w fyne.Window) *pfsettings {
	return &pfsettings{ window: w}
}
func (s *pfsettings) buildMobilePart() *fyne.Container {
	var mobileContainer *fyne.Container
	var err error
	var portslist []string
	s.mobileNumber = &widget.Entry{Text: fyne.CurrentApp().Preferences().StringWithFallback("mobilenumber", ""), OnChanged: func(v string) {
		fyne.CurrentApp().Preferences().SetString("mobilenumber", v)
	}}
	portslist, err = GetPortsList()
	if err != nil {
		log.Print("settings.buildUI #1 GetPortsList Error")
	}
	s.mobilePort = &widget.Entry{Text: fyne.CurrentApp().Preferences().StringWithFallback("mobileport", ""), OnChanged: func(v string) {
		fyne.CurrentApp().Preferences().SetString("mobileport", v)
	}}
	s.mobileModel = &widget.Select{Options: mobilemodels, OnChanged: func(sel string) {
		fyne.CurrentApp().Preferences().SetString("mobilemodel", sel)
	}, Selected: fyne.CurrentApp().Preferences().StringWithFallback("mobilemodel", ""),
	}
	allcountries := GetAllCountries()
	s.mobileCountry = &widget.Select{Options: allcountries, OnChanged: func(sel string) {
		fyne.CurrentApp().Preferences().SetString("moilemountry", sel)
	}, Selected: fyne.CurrentApp().Preferences().StringWithFallback("mobilemountry", "Sweden (+46)"),
	}
	s.mobileAddhash = &widget.Check{Text: "Add '#=' and messagenumber to end of sent messages",
		OnChanged: func(sel bool) { fyne.CurrentApp().Preferences().SetBool("addhash", sel) },
		Checked:   fyne.CurrentApp().Preferences().Bool("addhash")}
	s.btnTestSMS = &widget.Button{Text: "Click to send a test sms message to yourself.", OnTapped: func() {
		t := time.Now().Format("2006-01-02 15:04:05")
		testmessage := fmt.Sprintf("This is a short testmessage, sent %s", t)
		pn := fyne.CurrentApp().Preferences().StringWithFallback("mobilenumber", "")
		pn = Fixphonenumber(pn, s.mobileCountry.Selected)
		var sms theform = *new(theform)
		sms.Addhash = fyne.CurrentApp().Preferences().Bool("addHash")
		sms.Comport = fyne.CurrentApp().Preferences().StringWithFallback("mobileport", "COM2")
		sms.SendMessages([]string{pn}, testmessage)
		s.logtext.Text = ReadLastLineWithSeek(fyne.CurrentApp().Preferences().String("pfsmslog"), 6)
		s.logtext.Refresh()
	}}
	mobileContainer = container.NewGridWithColumns(2,
		NewBoldLabel("Your Phone Number"), s.mobileNumber,
		NewBoldLabel("Your Country"), s.mobileCountry,
		NewBoldLabel("Your Phone Model"), s.mobileModel,
		NewBoldLabel("Your Computer Port ("+strings.Join(portslist,", ")+")"), s.mobilePort,
		NewBoldLabel("Add some numbering into messages"), s.mobileAddhash,
		NewBoldLabel("Test mobile settings"), s.btnTestSMS,
	)
	return mobileContainer
}
func (s *pfsettings) buildFilePart() *fyne.Container {
	var fileContainer *fyne.Container
	var err error

	// Open Data directory!
	s.btnOpenDatadir = &widget.Button{Text: "Click to open data directory.", OnTapped: func() {
		path := GetHomeDir()
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "linux":
			// You can change "xdg-open" to your specific file manager if needed
			cmd = exec.Command("xdg-open", path)
		case "darwin":
			cmd = exec.Command("open", path) // For macOS
		case "windows":
			cmd = exec.Command("explorer", path) // For Windows
		default:
			log.Println("#1 buildFilePart unsupported operating system")
		}
		// Run the command
		cmd.Start()
	}}
	// Customers ImportExport
	s.btnImportCustomers = &widget.Button{Text: "Import Customers", OnTapped: func() {
		new(pfdatabase.DBtype).ImportCustomers(fyne.CurrentApp().Preferences().String("customersfile"))
		log.Println("Imported Customers")
	}}
	s.btnExportCustomers = &widget.Button{Text: "Export Customers", OnTapped: func() {
		new(pfdatabase.DBtype).ExportCustomers(fyne.CurrentApp().Preferences().String("customersfile"))
		log.Println("Exported Customers")
	}}
	s.btnDeleteCustomers = &widget.Button{Text: "Delete Customers", OnTapped: func() {
		new(pfdatabase.DBtype).DeleteCustomers()
		log.Println("Deleted Customers")
	}}
	s.customersfile = &widget.Label{Text: fyne.CurrentApp().Preferences().String("customersfile")}

	// Groups ImportExport
	s.btnImportGroups = &widget.Button{Text: "Import Groups", OnTapped: func() {
		var b0 []byte
		b0, err = os.ReadFile(fyne.CurrentApp().Preferences().String("groupsfile")) // SQL to make tables!
		if err != nil {
			log.Println("#1 ImportGroups", err.Error())
		}
		b := string(b0)
		b = strings.Replace(b, "\n", "", -1)
		new(pfdatabase.DBtype).ImportGroups(b)
		log.Println("Imported Groups")
	}}
	s.btnExportGroups = &widget.Button{Text: "Export Groups", OnTapped: func() {
		fn := fyne.CurrentApp().Preferences().StringWithFallback("groupsfile", fyne.CurrentApp().Preferences().String("groupsfile"))
		new(pfdatabase.DBtype).ExportGroups(fn)
		log.Println("Exported Groups")
	}}
	s.btnDeleteGroups = &widget.Button{Text: "Delete Groups", OnTapped: func() {
		new(pfdatabase.DBtype).DeleteGroups()
		log.Println("Deleted Groups")
	}}
	s.groupsfile = &widget.Label{Text: fyne.CurrentApp().Preferences().String("groupsfile")}

	// History ImportExport
	s.btnExportHistory = &widget.Button{Text: "Export History", OnTapped: func() {
		new(pfdatabase.DBtype).ExportHistory(fyne.CurrentApp().Preferences().String("historyfile"))
		log.Println("Exported History")
	}}
	s.btnOpenLog = &widget.Button{Text: "Open Log", OnTapped: func() {
		path := fyne.CurrentApp().Preferences().String("pfsmslog")
		// dd,_:=os.UserHomeDir()
		// dd = fmt.Sprintf("%s%c%s%c%s",dd,os.PathSeparator,"pfsms",os.PathSeparator,fn)
		var cmd *exec.Cmd

		// Determine the operating system
		switch runtime.GOOS {
		case "linux":
			// You can change "xdg-open" to your specific file manager if needed
			cmd = exec.Command("xdg-open", path)
		case "darwin":
			cmd = exec.Command("open", path) // For macOS
		case "windows":
			cmd = exec.Command("explorer", path) // For Windows
		default:
			log.Println("#2 buildFilePart unsupported operating system")
		}
		// Run the command
		cmd.Start()
	}}
	s.logfile= &widget.Label{Text: fyne.CurrentApp().Preferences().String("pfsmslog")}
	s.historyfile= &widget.Label{Text: fyne.CurrentApp().Preferences().String("historyfile")}

	fileContainer = container.NewGridWithColumns(2,
		NewBoldLabel("Location of default datafiles and textfiles:"), s.btnOpenDatadir,
		container.NewGridWithColumns(3, s.btnExportCustomers, s.btnImportCustomers, s.btnDeleteCustomers), s.customersfile,
		container.NewGridWithColumns(3, s.btnExportGroups, s.btnImportGroups, s.btnDeleteGroups), s.groupsfile,
		container.NewGridWithColumns(1, s.btnExportHistory), s.historyfile,
		container.NewGridWithColumns(1,  s.btnOpenLog), s.logfile,
	)
	return fileContainer
}
func (s *pfsettings) buildUI() *container.Scroll {
	// s.themeSelect = &widget.Select{Options: themes, OnChanged: func(tc string) {
	// 	s.app.Preferences().SetString("Theme", checkTheme(tc, s.app))
	// }, Selected: s.appSettings.Theme}
	s.logtext = &widget.Entry{}
	s.logtext = widget.NewMultiLineEntry()
	s.logtext.SetMinRowsVisible(6)
	s.logtext.Text = ReadLastLineWithSeek(fyne.CurrentApp().Preferences().String("pfsmslog"), 6)
	s.logtext.Refresh()

	return container.NewScroll(container.NewVBox(
		&widget.Card{Title: "Mobile Settings", Content: s.buildMobilePart()},
		&widget.Card{Title: "File Settings", Content: s.buildFilePart()},
		s.logtext,
	))
}
func (s *pfsettings) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Settings", Icon: theme.SettingsIcon(), Content: s.buildUI()}
}

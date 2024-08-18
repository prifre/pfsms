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
	"github.com/prifre/pfsms/ariasms"
	"github.com/prifre/pfsms/pfdatabase"
)

var (
	onOffOptions = []string{"On", "Off"}
	mobilemodels = []string{"Samsung S24", "Samsung S9"}
)

type pfsettings struct {
	mobileNumber  *widget.Entry
	mobileCountry *widget.Select
	mobileModel   *widget.Select
	mobilePort    *widget.Select
	mobileAddhash *widget.Check
	btnTestSMS    *widget.Button

	btnOpenDatadir     *widget.Button
	customersfile      *widget.Label
	btnImportCustomers *widget.Button
	btnExportCustomers *widget.Button
	groupsfile         *widget.Label
	btnImportGroups    *widget.Button
	btnExportGroups    *widget.Button
	fileLog            *widget.Label
	btnExportHistory   *widget.Button
	btnOpenLog         *widget.Button
	logtext            *widget.Entry

	window fyne.Window
	app    fyne.App
}

func NewSettings(a fyne.App, w fyne.Window) *pfsettings {
	return &pfsettings{app: a, window: w}
}
func (s *pfsettings) buildMobilePart() *fyne.Container {
	var mobileContainer *fyne.Container
	s.mobileNumber = &widget.Entry{Text: s.app.Preferences().StringWithFallback("mobilenumber", ""), OnChanged: func(v string) {
		s.app.Preferences().SetString("mobilenumber", s.mobileNumber.Text)
	}}
	var sms ariasms.SMStype = *new(ariasms.SMStype)
	p, err := sms.GetPortsList()
	if err != nil {
		log.Print("settings.buildUI #1 GetPortsList Error")
	}
	s.mobilePort = &widget.Select{Options: p, OnChanged: func(sel string) {
		s.app.Preferences().SetString("mobilePort", sel)
	}, Selected: s.app.Preferences().StringWithFallback("mobilePort", ""),
	}
	s.mobileModel = &widget.Select{Options: mobilemodels, OnChanged: func(sel string) {
		s.app.Preferences().SetString("mobileModel", sel)
	}, Selected: s.app.Preferences().StringWithFallback("mobileModel", ""),
	}
	allcountries := GetAllCountries()
	s.mobileCountry = &widget.Select{Options: allcountries, OnChanged: func(sel string) {
		s.app.Preferences().SetString("moileCountry", sel)
	}, Selected: s.app.Preferences().StringWithFallback("mobileCountry", "Sweden (+46)"),
	}
	s.mobileAddhash = &widget.Check{Text: "Add '#=' and messagenumber to end of sent messages",
		OnChanged: func(sel bool) { s.app.Preferences().SetBool("addHash", sel) },
		Checked:   s.app.Preferences().Bool("addHash")}
	s.btnTestSMS = &widget.Button{Text: "Click to send a test sms message to yourself.", OnTapped: func() {
		t := time.Now().Format("2006-01-02 15:04:05")
		testmessage := fmt.Sprintf("This is a short testmessage, sent %s", t)
		pn := s.app.Preferences().StringWithFallback("mobilenumber", "")
		pn = Fixphonenumber(pn, s.mobileCountry.Selected)
		var sms ariasms.SMStype = *new(ariasms.SMStype)
		sms.Addhash = s.app.Preferences().Bool("addHash")
		sms.Comport = s.app.Preferences().StringWithFallback("mobilePort", "COM2")
		sms.SendMessage([]string{pn}, testmessage)
		s.logtext.Text = ReadLastLineWithSeek(fyne.CurrentApp().Preferences().String("pfsmslog"), 6)
		s.logtext.Refresh()
	}}
	mobileContainer = container.NewGridWithColumns(2,
		NewBoldLabel("Your Phone Number"), s.mobileNumber,
		NewBoldLabel("Your Country"), s.mobileCountry,
		NewBoldLabel("Your Phone Model"), s.mobileModel,
		NewBoldLabel("Your Computer Port"), s.mobilePort,
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
		dd := GetHomeDir()
		if runtime.GOOS == "darwin" {
			cmd := `open "` + dd + `"`
			exec.Command("/bin/bash", "-c", cmd).Start()
		} else {
			exec.Command("explorer", dd).Start()
		}
	}}
	// Customers ImportExport
	s.btnImportCustomers = &widget.Button{Text: "Import Customers", OnTapped: func() {
		new(pfdatabase.DBtype).ImportCustomers(fyne.CurrentApp().Preferences().String("customersfile"))
		s.window.SetContent(Create(s.app, s.window))
		// t.tableShowCustomers=t.listCustomers()
		// t.tableShowCustomers.Refresh()
	}}
	s.btnExportCustomers = &widget.Button{Text: "Export Customers", OnTapped: func() {
		new(pfdatabase.DBtype).ExportCustomers(fyne.CurrentApp().Preferences().String("customersfile"))
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
	}}
	s.btnExportGroups = &widget.Button{Text: "Export Groups", OnTapped: func() {
		fn := s.app.Preferences().StringWithFallback("groupsfile", fyne.CurrentApp().Preferences().String("groupsfile"))
		new(pfdatabase.DBtype).ExportGroups(fn)
	}}
	s.groupsfile = &widget.Label{Text: fyne.CurrentApp().Preferences().String("groupsfile")}

	// History ImportExport
	s.btnExportHistory = &widget.Button{Text: "Export History", OnTapped: func() {
		new(pfdatabase.DBtype).ExportHistory(fyne.CurrentApp().Preferences().String("historyfile"))
	}}
	s.btnOpenLog = &widget.Button{Text: "Open Log", OnTapped: func() {
		dd := fyne.CurrentApp().Preferences().String("pfsmslog")
		// dd,_:=os.UserHomeDir()
		// dd = fmt.Sprintf("%s%c%s%c%s",dd,os.PathSeparator,"pfsms",os.PathSeparator,fn)
		if runtime.GOOS == "darwin" {
			cmd := `open "` + dd + `"`
			exec.Command("/bin/bash", "-c", cmd).Start()
		} else {
			exec.Command("explorer", dd).Start()
		}
	}}
	s.fileLog = &widget.Label{Text: fyne.CurrentApp().Preferences().String("pfsmslog")}

	fileContainer = container.NewGridWithColumns(2,
		NewBoldLabel("Location of default datafiles and textfiles:"), s.btnOpenDatadir,
		container.NewGridWithColumns(2, s.btnExportCustomers, s.btnImportCustomers), s.customersfile,
		container.NewGridWithColumns(2, s.btnExportGroups, s.btnImportGroups), s.groupsfile,
		container.NewGridWithColumns(2, s.btnExportHistory, s.btnOpenLog), s.fileLog,
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

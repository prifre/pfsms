package ui

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/prifre/pfsms/ariasms"
)
var (
	onOffOptions 	= []string{"On", "Off"}
	mobilemodels    = []string{"Samsung S24","Samsung S9"}
)
type settings struct {
	mobileNumber  		*widget.Entry
	mobileCountry		*widget.Select
	mobileModel  		*widget.Select
	mobilePort  		*widget.Select
	mobileAddhash		*widget.Check
	btnTestSMS		*widget.Button

	btnOpenDatadir		*widget.Button
	fileCustomers		*widget.Entry
	fileGroups			*widget.Entry

	window  		    fyne.Window
	app        			fyne.App
}
func NewSettings(a fyne.App, w fyne.Window) *settings {
	return &settings{app: a, window: w}
}
func (s *settings) buildUI() *container.Scroll {
	// s.themeSelect = &widget.Select{Options: themes, OnChanged: func(tc string) {
	// 	s.app.Preferences().SetString("Theme", checkTheme(tc, s.app))
	// }, Selected: s.appSettings.Theme}

	s.mobileNumber = &widget.Entry{Text:s.app.Preferences().StringWithFallback("mobilenumber",""),OnChanged: func(v string) {
		s.app.Preferences().SetString("mobilenumber",s.mobileNumber.Text)
	}}
	var sms ariasms.SMStype =*new(ariasms.SMStype)
	p,err:=sms.GetPortsList()
	if err!=nil {
		log.Print("settings.buildUI #1 GetPortsList Error")
	}
	s.mobilePort =&widget.Select{Options: p, OnChanged: func(sel string) {
		s.app.Preferences().SetString("mobilePort", sel)
		}, Selected: s.app.Preferences().StringWithFallback("mobilePort", ""),
	}
	s.mobileModel =&widget.Select{Options: mobilemodels, OnChanged: func(sel string) {
		s.app.Preferences().SetString("mobileModel", sel)
		}, Selected: s.app.Preferences().StringWithFallback("mobileModel", ""),
	}
	allcountries := GetAllCountries()
	s.mobileCountry =&widget.Select{Options: allcountries, OnChanged: func(sel string) {
		s.app.Preferences().SetString("moileCountry", sel)
		}, Selected: s.app.Preferences().StringWithFallback("mobileCountry", "Sweden (+46)"),
	}
	s.mobileAddhash = &widget.Check{Text:"Add '#=' and messagenumber to end of sent messages",
			OnChanged: func(sel bool) {	s.app.Preferences().SetBool("addHash",sel) },
			Checked:s.app.Preferences().Bool("addHash")}
	s.btnTestSMS = & widget.Button{Text:"Click to send a test sms message to yourself.",OnTapped: func ()  {
		t:=time.Now().Format("2006-01-02 15:04:05")
		testmessage:=fmt.Sprintf("This is a short testmessage, sent %s", t)
		pn:= s.app.Preferences().StringWithFallback("mobilenumber","")
		var sms ariasms.SMStype =*new(ariasms.SMStype)
		sms.Addhash=s.app.Preferences().Bool("addHash")
		sms.Comport = s.app.Preferences().StringWithFallback("mobilePort", "COM2")
		sms.SendMessage([]string{pn},testmessage)
	}}
	mobileContainer := container.NewGridWithColumns(2,
		NewBoldLabel("Your Phone Number"),s.mobileNumber,
		NewBoldLabel("Your Country"),s.mobileCountry,
		NewBoldLabel("Your Phone Model"),s.mobileModel,
		NewBoldLabel("Your Computer Port"),s.mobilePort,
		NewBoldLabel("Add some numbering into messages"),s.mobileAddhash,
		NewBoldLabel("Test mobile settings"),s.btnTestSMS,
	)
	s.btnOpenDatadir = &widget.Button{Text:"Click to open data directory.",OnTapped: func ()  {
		dd,_:=os.UserHomeDir()
		dd = fmt.Sprintf("%s%c",dd,os.PathSeparator)
		cmd := exec.Command(`explorer`, `/select,`, dd)
		cmd.Run()
	}}
	s.fileCustomers = &widget.Entry{Text:s.app.Preferences().StringWithFallback("fileCustomers",""),OnChanged: func(v string) {
		s.app.Preferences().SetString("fileCustomers",s.fileCustomers.Text)
	}}
	s.fileGroups = &widget.Entry{Text:s.app.Preferences().StringWithFallback("fileGroups",""),OnChanged: func(v string) {
		s.app.Preferences().SetString("fileGroups",s.fileCustomers.Text)
	}}
	fileContainer := container.NewGridWithColumns(2,
		NewBoldLabel("Customer Import/Export file"), s.fileCustomers,
		NewBoldLabel("Groups Import/Export file"), s.fileGroups,
	)

	return container.NewScroll(container.NewVBox(
		&widget.Card{Title: "Mobile Settings", Content: mobileContainer},
		&widget.Card{Title: "File Settings", Content: fileContainer},
	))
}
func (s *settings) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Settings", Icon: theme.SettingsIcon(), Content: s.buildUI()}
}

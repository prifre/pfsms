package ui

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/prifre/pfsms/ariasms"
	"github.com/prifre/pfsms/db"
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
	btnTest				*widget.Button

	useEmail	     	*widget.RadioGroup
	emailServer  		*widget.Entry
	emailSLabel      	*widget.Label
	emailPort  			*widget.Entry
	emailPortLabel     	*widget.Label
	emailUser  			*widget.Entry
	emailULabel      	*widget.Label
	emailPassword  		*widget.Entry
	emailPLabel      	*widget.Label
	emailFLabel      	*widget.Label
	emailFrequency      *widget.Slider

	window  		    fyne.Window
	app        			fyne.App
}
func NewSettings(a fyne.App, w fyne.Window) *settings {
	return &settings{app: a, window: w}
}
func (s *settings) onUseEmailChanged(selected string) {
//	s.client.OverwriteExisting = selected == "On"
	s.app.Preferences().SetString("UseEmail",selected)
	s.emailServer.Hidden=(s.app.Preferences().StringWithFallback("UseEmail", "Off")=="Off")
	s.emailSLabel.Hidden=(s.app.Preferences().StringWithFallback("UseEmail", "Off")=="Off")
	s.emailPort.Hidden=(s.app.Preferences().StringWithFallback("UseEmail", "Off")=="Off")
	s.emailPortLabel.Hidden=(s.app.Preferences().StringWithFallback("UseEmail", "Off")=="Off")
	s.emailUser.Hidden=(s.app.Preferences().StringWithFallback("UseEmail", "Off")=="Off")
	s.emailULabel.Hidden=(s.app.Preferences().StringWithFallback("UseEmail", "Off")=="Off")
	s.emailPassword.Hidden=(s.app.Preferences().StringWithFallback("UseEmail", "Off")=="Off")
	s.emailPLabel.Hidden=(s.app.Preferences().StringWithFallback("UseEmail", "Off")=="Off")
	s.emailFrequency.Hidden=(s.app.Preferences().StringWithFallback("UseEmail", "Off")=="Off")
	s.emailFLabel.Hidden=(s.app.Preferences().StringWithFallback("UseEmail", "Off")=="Off")
}
func (s *settings) getPassword() string {
	prefPassword:= s.app.Preferences().StringWithFallback("ePassword","")
	d:=new(db.DBtype)
	d.Opendb()
	realPassword,err:=d.DecryptPassword(prefPassword)
	if err!=nil {
		log.Println("getPassword DecryptPassword error")
	}
	return realPassword
}
func (s *settings) setPassword(realPassword string) error {
	d:=new(db.DBtype)
	d.Opendb()
	prefPassword,err:=d.EncryptPassword(realPassword)
	if err!=nil {
		log.Println("setPassWord EncryptPassword error")
	}
	s.app.Preferences().SetString("ePassword",prefPassword)
	return err
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
	s.btnTest = & widget.Button{Text:"Click to send a test sms message to yourself.",OnTapped: func ()  {
		t:=time.Now().Format("2006-01-02 15:04:05")
		testmessage:=fmt.Sprintf("This is a short testmessage, sent %s", t)
		pn:= s.app.Preferences().StringWithFallback("mobilenumber","")
		var sms ariasms.SMStype =*new(ariasms.SMStype)
		sms.Addhash=s.app.Preferences().Bool("addHash")
		sms.Comport = s.app.Preferences().StringWithFallback("mobilePort", "COM2")
		sms.SendMessage([]string{pn},testmessage)
	}}

	s.emailSLabel = &widget.Label{Text: "Email Server", TextStyle: fyne.TextStyle{Bold: true}}
	s.emailServer = &widget.Entry{Text:s.app.Preferences().StringWithFallback("eServer",""),OnChanged: func(v string) {
		s.app.Preferences().SetString("eServer",s.emailServer.Text)
	}}
	s.emailPortLabel = &widget.Label{Text: "Email Server Port", TextStyle: fyne.TextStyle{Bold: true}}
	s.emailPort = &widget.Entry{Text:s.app.Preferences().StringWithFallback("ePort","993"),OnChanged: func(v string) {
		v0:=""
		for i:=0;i<len(v);i++ {

			if (v[i] >= '0' && v[i] <= '9') {
				v0 = v0+string(v[i])
			}	
		}
		s.emailPort.SetText(v0)
		s.app.Preferences().SetString("ePort",s.emailPort.Text)
	}}
	s.emailULabel = &widget.Label{Text: "Email User", TextStyle: fyne.TextStyle{Bold: true}}
	s.emailUser = &widget.Entry{Text:s.app.Preferences().StringWithFallback("eUser",""),OnChanged: func(v string) {
		s.app.Preferences().SetString("eUser",s.emailUser.Text)
	}}
	s.emailPLabel = &widget.Label{Text: "Email Password", TextStyle: fyne.TextStyle{Bold: true}}
	s.emailPassword = &widget.Entry{Text:s.getPassword(),OnChanged: func(v string) {
		s.setPassword(s.emailPassword.Text)
	}}
	s.emailFLabel = &widget.Label{Text:
		fmt.Sprintf("Email frequency (%d min)",int(s.app.Preferences().FloatWithFallback("eFrequency",10))),
		 TextStyle: fyne.TextStyle{Bold: true}}
	s.emailFrequency=&widget.Slider{Value: s.app.Preferences().FloatWithFallback("eFrequency",10),
		Min: 1.0, Max:60.0, Step: 1, 
		OnChanged: func(i float64) {
			s.emailFLabel.SetText(fmt.Sprintf("Email frequency (%d min)",int(i)))
		}, 
		OnChangeEnded: func(i float64) {
			s.emailFrequency.Value = i
			s.app.Preferences().SetFloat("eFrequency",i)
	}}
	s.useEmail = &widget.RadioGroup{Options: onOffOptions, Horizontal: true, Required: true, OnChanged: s.onUseEmailChanged}
	s.useEmail.SetSelected(s.app.Preferences().StringWithFallback("UseEmail", "Off"))
	s.onUseEmailChanged(s.app.Preferences().StringWithFallback("UseEmail", "Off"))

	mobileContainer := container.NewGridWithColumns(2,
		NewBoldLabel("Your Phone Number"),s.mobileNumber,
		NewBoldLabel("Your Country"),s.mobileCountry,
		NewBoldLabel("Your Phone Model"),s.mobileModel,
		NewBoldLabel("Your Computer Port"),s.mobilePort,
		NewBoldLabel("Add some numbering into messages"),s.mobileAddhash,
		NewBoldLabel("Test mobile settings"),s.btnTest,
	)

	dataContainer := container.NewGridWithColumns(2,
		NewBoldLabel("Use Email"), s.useEmail,
		s.emailSLabel, s.emailServer,
		s.emailPortLabel, s.emailPort,
		s.emailULabel, s.emailUser,
		s.emailPLabel, s.emailPassword,
		s.emailFLabel,s.emailFrequency,
	)
	return container.NewScroll(container.NewVBox(
		&widget.Card{Title: "Mobile Settings", Content: mobileContainer},
		&widget.Card{Title: "Email Settings", Content: dataContainer},
	))
}
func (s *settings) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Settings", Icon: theme.SettingsIcon(), Content: s.buildUI()}
}

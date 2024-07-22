package ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/prifre/pfsms/db"
)

var (
	themes       = []string{"Adaptive (requires restart)", "Light", "Dark"}
	onOffOptions = []string{"On", "Off"}
)

// AppSettings contains settings specific to the application
type AppSettings struct {
	// Theme holds the current theme
	Theme string
}

type settings struct {
	themeSelect 		*widget.Select
	mobileNumber  		*widget.Entry
	mobileLabel      	*widget.Label

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

	appSettings *AppSettings
	window      fyne.Window
	app         fyne.App
}

func NewSettings(a fyne.App, w fyne.Window,  as *AppSettings) *settings {
	return &settings{app: a, window: w,  appSettings: as}
}

func (s *settings) onThemeChanged(selected string) {
	s.app.Preferences().SetString("Theme", checkTheme(selected, s.app))
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
		log.Println("getPassWord DecryptPassword error")
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
	s.themeSelect = &widget.Select{Options: themes, OnChanged: s.onThemeChanged, Selected: s.appSettings.Theme}

	s.mobileLabel = &widget.Label{Text: "Your Mobile#", TextStyle: fyne.TextStyle{Bold: true}}
	s.mobileNumber = &widget.Entry{Text:s.app.Preferences().StringWithFallback("mobilenumber",""),OnChanged: func(v string) {
		s.app.Preferences().SetString("mobilenumber",s.mobileNumber.Text)
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

	interfaceContainer := container.NewGridWithColumns(2,
		newBoldLabel("Application Theme"), s.themeSelect,
	)
	mobileContainer := container.NewGridWithColumns(2,
	 s.mobileLabel,s.mobileNumber,
	)

	dataContainer := container.NewGridWithColumns(2,
		newBoldLabel("Use Email"), s.useEmail,
		s.emailSLabel, s.emailServer,
		s.emailPortLabel, s.emailPort,
		s.emailULabel, s.emailUser,
		s.emailPLabel, s.emailPassword,
		s.emailFLabel,s.emailFrequency,
	)
	return container.NewScroll(container.NewVBox(
		&widget.Card{Title: "User Interface", Content: interfaceContainer},
		&widget.Card{Title: "Mobile Settings", Content: mobileContainer},
		&widget.Card{Title: "Email Settings", Content: dataContainer},
	))
}

func (s *settings) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Settings", Icon: theme.SettingsIcon(), Content: s.buildUI()}
}

func checkTheme(themec string, a fyne.App) string {
	switch themec {
	case "Light":
		//lint:ignore SA1019 Not quite ready for removal on Linux.
		a.Settings().SetTheme(theme.LightTheme())
	case "Dark":
		//lint:ignore SA1019 Not quite ready for removal on Linux.
		a.Settings().SetTheme(theme.DarkTheme())
	}
	return themec
}

func newBoldLabel(text string) *widget.Label {
	return &widget.Label{Text: text, TextStyle: fyne.TextStyle{Bold: true}}
}

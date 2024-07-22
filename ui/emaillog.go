package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	pfemail "github.com/prifre/pfsms/email"
)

type AppEmail struct {
	// Theme holds the current theme
	Theme string
}

type thelog struct {
	btnStart	*widget.Button
	btnCheck	*widget.Button
	logtext				*widget.Label
	appEmail 	*AppEmail
	window      	fyne.Window
	app         	fyne.App
}

func NewEmaillog(a fyne.App, w fyne.Window,  el *AppEmail) *thelog {
	return &thelog{app: a, window: w,  appEmail: el}
}
func (s *thelog) buildLog() *container.Scroll {
	var e pfemail.Etype
	var err error
	s.logtext = &widget.Label{}
	s.btnCheck = &widget.Button{Text:"Check Email",OnTapped: func() {
		st :=new(settings)
		e.SetupEmail(s.app.Preferences().StringWithFallback("eServer",""),
		s.app.Preferences().StringWithFallback("eUser",""),
		st.getPassword(),
		s.app.Preferences().StringWithFallback("ePort",""))
		s.Addtolog(fmt.Sprintf("\r\n%s%s",time.Now(),"Checking email..."+s.app.Preferences().StringWithFallback("eUser","")))
		err=e.Checkemaillogin()
		if err!=nil {
			s.Addtolog("\r\nResult of email login check: %s"+err.Error())
		} else {
			s.Addtolog(fmt.Sprintf("\r\n%s%s",time.Now(),"Email check ok."))
		}
	}}
	s.btnStart = &widget.Button{Text:"Start Email",OnTapped: func() {
		st :=new(settings)
		e.SetupEmail(s.app.Preferences().StringWithFallback("eServer",""),
		s.app.Preferences().StringWithFallback("eUser",""),
		st.getPassword(),
		s.app.Preferences().StringWithFallback("ePort",""))
		fmt.Println("START CHECKING EMAIL!!")
		s.Addtolog(fmt.Sprintf("\r\n%s%s",time.Now(),"Handling mail..."))
		et:=e.Getallmailmovetosmsfolder()
		if et!=nil {
			s.Addtolog(et.Text)
		} else {
			s.Addtolog(fmt.Sprintf("\r\n%s%s",time.Now(),"No mail..."))
		}
	}}
	return container.NewScroll(		
		container.NewVBox(
			container.NewHBox(s.btnCheck,s.btnStart),
			s.logtext,
	))
}

func (s *thelog) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Mail Log", Icon: theme.ComputerIcon(), Content: s.buildLog()}
}
func (s *thelog) Addtolog (m string) {
	s.logtext.SetText(fmt.Sprintf("%s%s",s.logtext.Text ,m ))
}


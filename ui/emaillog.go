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
	s.logtext = &widget.Label{}
	s.btnStart = &widget.Button{Text:"Start Email",OnTapped: func() {
		e.SetupEmail(s.app.Preferences().StringWithFallback("eServer",""),
		s.app.Preferences().StringWithFallback("eUser",""),
		s.app.Preferences().StringWithFallback("ePassword",""),
		s.app.Preferences().StringWithFallback("ePort",""))
		fmt.Println("START CHECKING EMAIL!!")
		s.Addtolog(fmt.Sprintf("\r\n%s%s",time.Now(),"Checking mail..."))
		et:=e.Getonemail()
		if et!=nil {
			s.Addtolog(et.Text)
		} else {
			s.Addtolog(fmt.Sprintf("\r\n%s%s",time.Now(),"No mail..."))
		}
	}}
	return container.NewScroll(		
		container.NewVBox(
			s.btnStart,
			s.logtext,
	))
}

func (s *thelog) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Mail Log", Icon: theme.ComputerIcon(), Content: s.buildLog()}
}
func (s *thelog) Addtolog (m string) {
	s.logtext.SetText(fmt.Sprintf("%s%s",s.logtext.Text ,m ))
}


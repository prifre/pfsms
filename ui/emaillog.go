package ui

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/prifre/pfsms/db"
	pfemail "github.com/prifre/pfsms/email"
)

type theemaillog struct {
	btnStart			*widget.Button
	btnCheck			*widget.Button
	btnImport 			*widget.Button
	logtext				*widget.Label
	window      		fyne.Window
	app         		fyne.App
}

func NewEmaillog(a fyne.App, w fyne.Window) *theemaillog {
	return &theemaillog{app: a, window: w}
}
func (s *theemaillog) buildLog() *container.Scroll {
	var e pfemail.Etype
	var err error
	s.logtext = &widget.Label{}
	s.btnCheck = &widget.Button{Text:"Check Email",OnTapped: func() {
		d:=new(db.DBtype)
		d.Opendb()
		p,_:=d.DecryptPassword(s.app.Preferences().StringWithFallback("ePassword",""))
		e.SetupEmail(s.app.Preferences().StringWithFallback("eServer",""),
		s.app.Preferences().StringWithFallback("eUser",""),p,
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
		d:=new(db.DBtype)
		d.Opendb()
		p,_:=d.DecryptPassword(s.app.Preferences().StringWithFallback("ePassword",""))
		e.SetupEmail(s.app.Preferences().StringWithFallback("eServer",""),
		s.app.Preferences().StringWithFallback("eUser",""),p,
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
	s.btnImport = widget.NewButton("Import customers",func() {
		var importfilename string
		dialog.ShowFileOpen(func (f fyne.URIReadCloser,err error) {
			if err==nil && f!=nil {
				importfilename=f.URI().String()
				importfilename = strings.Replace(importfilename,"file://","",-1)
				if importfilename>"" {
					db:=new(db.DBtype)
					db.ImportCustomers(importfilename)
				}
			}
		},s.window)		
	})

	return container.NewScroll(		
		container.NewVBox(
			container.NewHBox(s.btnCheck,s.btnStart,s.btnImport),
			s.logtext,
	))
}

func (s *theemaillog) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Mail", Icon: theme.ComputerIcon(), Content: s.buildLog()}
}
func (s *theemaillog) Addtolog (m string) {
	s.logtext.SetText(fmt.Sprintf("%s%s",s.logtext.Text ,m ))
	err :=Appendtotextfile("emaillog.txt",m)
	if err!=nil {
		fmt.Println("#1 Addtolog Appendtotextfile failed")
	}
}


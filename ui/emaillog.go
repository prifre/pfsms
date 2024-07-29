package ui

import (
	"fmt"
	"log"
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
	btnExport 			*widget.Button

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

	logtext				*widget.Label
	window      		fyne.Window
	app         		fyne.App
}

func NewEmaillog(a fyne.App, w fyne.Window) *theemaillog {
	return &theemaillog{app: a, window: w}
}
func (s *theemaillog) buildLog() *container.Scroll {
	var err error

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

	emailContainer := container.NewGridWithColumns(2,
		NewBoldLabel("Use Email"), s.useEmail,
		s.emailSLabel, s.emailServer,
		s.emailPortLabel, s.emailPort,
		s.emailULabel, s.emailUser,
		s.emailPLabel, s.emailPassword,
		s.emailFLabel,s.emailFrequency,
	)


	s.logtext = &widget.Label{}
	s.btnCheck = &widget.Button{Text:"Check Email",OnTapped: func() {
		d:=new(db.DBtype)
		d.Opendb()
		p,_:=d.DecryptPassword(s.app.Preferences().StringWithFallback("ePassword",""))
		e:=new(pfemail.Etype)
		e.SetupEmail(s.app.Preferences().StringWithFallback("eServer",""),
		s.app.Preferences().StringWithFallback("eUser",""),p,
		s.app.Preferences().StringWithFallback("ePort",""))
		m:=fmt.Sprintf("\r\n%s %s %s",time.Now().Format("2006-01-02 15:04:05"),"Checking email...",s.app.Preferences().StringWithFallback("eUser",""))
		Appendtotextfile("emaillog.txt",m)
		err=e.Checkemaillogin()
		if err!=nil {
			m:="\r\nResult of email login check: %s"+err.Error()
			Appendtotextfile("emaillog.txt",m)
		} else {
			m:=fmt.Sprintf("\r\n%s %s",time.Now().Format("2006-01-02 15:04:05"),"Email check ok.")
			Appendtotextfile("emaillog.txt",m)
		}
		s.logtext.Text += m
		s.logtext.Refresh()
	}}
	s.btnStart = &widget.Button{Text:"Start Email",OnTapped: func() {
		d:=new(db.DBtype)
		d.Opendb()
		p,_:=d.DecryptPassword(s.app.Preferences().StringWithFallback("ePassword",""))
		e:=new(pfemail.Etype)
		e.SetupEmail(s.app.Preferences().StringWithFallback("eServer",""),
		s.app.Preferences().StringWithFallback("eUser",""),p,
		s.app.Preferences().StringWithFallback("ePort",""))
		fmt.Println("START CHECKING EMAIL!!")
		m:=fmt.Sprintf("\r\n%s %s",time.Now().Format("2006-01-02 15:04:05"),"Handling mail...")
		Appendtotextfile("emaillog.txt",m)
		et:=e.Getallmailmovetosmsfolder()
		if et!=nil {
			Appendtotextfile("emaillog.txt",fmt.Sprintf("\r\n%s %s",time.Now().Format("2006-01-02 15:04:05"),"Got Email!!!"))
		} else {
			Appendtotextfile("emaillog.txt",fmt.Sprintf("\r\n%s %s",time.Now().Format("2006-01-02 15:04:05"),"No mail..."))
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
	s.btnExport = widget.NewButton("Export customers",func() {
		var exportfilename string
		dialog.ShowFileSave(func (f fyne.URIWriteCloser,err error) {
			if err==nil && f!=nil {
				exportfilename=f.URI().String()
				exportfilename = strings.Replace(exportfilename,"file://","",-1)
				if exportfilename>"" {
					db:=new(db.DBtype)
					db.ExportCustomers(exportfilename)
				}
			}
		},s.window)
	})

	return container.NewScroll(		
		container.NewVBox(
			&widget.Card{Title: "Email Settings", Content: emailContainer},
			container.NewHBox(s.btnCheck,s.btnStart,s.btnImport,s.btnExport),
			s.logtext,
	))
}

func (s *theemaillog) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Email", Icon: theme.ComputerIcon(), Content: s.buildLog()}
}
func (s *theemaillog) Addtolog (m string) {
	s.logtext.SetText(fmt.Sprintf("%s%s",s.logtext.Text ,m ))
	err :=Appendtotextfile("emaillog.txt",m)
	if err!=nil {
		fmt.Println("#1 Addtolog Appendtotextfile failed")
	}
}
func (s *theemaillog) onUseEmailChanged(selected string) {
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
	func (s *theemaillog) getPassword() string {
		prefPassword:= s.app.Preferences().StringWithFallback("ePassword","")
		d:=new(db.DBtype)
		d.Opendb()
		realPassword,err:=d.DecryptPassword(prefPassword)
		if err!=nil {
			log.Println("getPassword DecryptPassword error")
		}
		return realPassword
	}
	func (s *theemaillog) setPassword(realPassword string) error {
		d:=new(db.DBtype)
		d.Opendb()
		prefPassword,err:=d.EncryptPassword(realPassword)
		if err!=nil {
			log.Println("setPassWord EncryptPassword error")
		}
		s.app.Preferences().SetString("ePassword",prefPassword)
		return err
	}
	

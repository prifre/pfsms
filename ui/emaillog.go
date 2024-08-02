package ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/prifre/pfsms/db"
	pfemail "github.com/prifre/pfsms/email"
)

type theemaillog struct {
	btnStart			*widget.Button
	btnCheck			*widget.Button

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
	p:=s.getPassword()
	s.emailPassword = &widget.Entry{Text:p,OnChanged: func(v string) {
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
	s.logtext =&widget.Label{}
	var txt string
	txt,err = ReadLastLineWithSeek("emaillog.txt",18)
	if err!=nil {
		log.Println("#1 buildLog!", err.Error())
	}
	s.logtext.Text=txt
	s.btnCheck = &widget.Button{Text:"Check Email",OnTapped: func() {
		e:=new(pfemail.Etype)
		err = e.Checkemaillogin()
		if err!=nil {
			Appendtotextfile("emaillog.txt","Email check failed.\r\n")
		} else {
			Appendtotextfile("emaillog.txt","Email check ok.\r\n")
		}
		var m string
		m,err = ReadLastLineWithSeek("emaillog.txt",18)
		if err!=nil {
			log.Println("#1 buildLog!", err.Error())
		}
		s.logtext.Text = m
		s.logtext.Refresh()
	}}
	s.btnStart = &widget.Button{Text:"Start Email",OnTapped: func() {
		e:=new(pfemail.Etype)
		Appendtotextfile("emaillog.txt","Handling mail...\r\n")
		err = e.Login()
		if err!=nil {
			log.Println("Login failed!")
		}
		m0:=e.Getallsmsmail()
		if m0!=nil {
			Appendtotextfile("emaillog.txt","Got Email!!!\r\n")
			m :=  "SUBJECT:"+ m0[0].Envelope.Subject
			// m += "SENDER: "+ m0[0].Envelope.Sender.Address
			// fmt.Println("\r\nFlags: ", m0[0].Flags)		} else {
			Appendtotextfile("emaillog.txt",m+"\r\n")
		} else {
			Appendtotextfile("emaillog.txt","No mail...\r\n")
		}
		err = e.Moveallsmsmail()
		if err!=nil {
			Appendtotextfile("emaillog.txt","Move SMS Mail failed.\r\n")
		} else {
			Appendtotextfile("emaillog.txt","Moved SMS mail to sms folder.\r\n")
		}
		var m string
		m,err = ReadLastLineWithSeek("emaillog.txt",18)
		if err!=nil {
			log.Println("#2 buildLog!", err.Error())
		}
		s.logtext.Text = m
		s.logtext.Refresh()
	}}

	return container.NewScroll(		
		container.NewVBox(
			&widget.Card{Title: "Email Settings", Content: emailContainer},
			container.NewHBox(s.btnCheck,s.btnStart),
			s.logtext,
	))
}

func (s *theemaillog) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Email", Icon: theme.ComputerIcon(), Content: s.buildLog()}
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
		var hash,realPassword string
		var err error
		hash,err =pfdatabase.MakeHash()
		if err!=nil {
			log.Println("buildLog onUseEmailChanged MakeHash error ",err.Error())
		}
		realPassword,err=pfdatabase.DecryptPassword(prefPassword,hash)
		if err!=nil {
			log.Println("getPassword onEmailChanged DecryptPassword error")
		}
		return realPassword
	}
	func (s *theemaillog) setPassword(realPassword string) error {
		var hash,prefPassword string
		var err error
		hash,err =pfdatabase.MakeHash()
		if err!=nil {
			log.Println("buildLog setPassord MakeHash error ",err.Error())
		}
		prefPassword,err=pfdatabase.EncryptPassword(realPassword,hash)
		if err!=nil {
			log.Println("setPassWord EncryptPassword error")
		}
		s.app.Preferences().SetString("ePassword",prefPassword)
		return err
	}
	

package ui

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/prifre/pfsms/pfdatabase"
	"github.com/prifre/pfsms/pfemail"
)

type theemail struct {
	btnStart *widget.Button
	btnCheck *widget.Button

	useEmail       *widget.RadioGroup
	emailServer    *widget.Entry
	emailSLabel    *widget.Label
	emailPort      *widget.Entry
	emailPortLabel *widget.Label
	emailUser      *widget.Entry
	emailULabel    *widget.Label
	emailPassword  *widget.Entry
	emailPLabel    *widget.Label
	emailFLabel    *widget.Label
	emailFrequency *widget.Slider

	logtext *widget.Label
	window  fyne.Window
}

func NewEmail(w fyne.Window) *theemail {
	return &theemail{window: w}
}
func (s *theemail) buildUI() *container.Scroll {
	var err error

	s.emailSLabel = &widget.Label{Text: "Email Server", TextStyle: fyne.TextStyle{Bold: true}}
	s.emailServer = &widget.Entry{Text: fyne.CurrentApp().Preferences().StringWithFallback("eserver", ""), OnChanged: func(v string) {
		fyne.CurrentApp().Preferences().SetString("eserver", s.emailServer.Text)
	}}
	s.emailPortLabel = &widget.Label{Text: "Email Server Port", TextStyle: fyne.TextStyle{Bold: true}}
	s.emailPort = &widget.Entry{Text: fyne.CurrentApp().Preferences().StringWithFallback("eport", "993"), OnChanged: func(v string) {
		v0 := ""
		for i := 0; i < len(v); i++ {

			if v[i] >= '0' && v[i] <= '9' {
				v0 = v0 + string(v[i])
			}
		}
		s.emailPort.SetText(v0)
		fyne.CurrentApp().Preferences().SetString("eport", s.emailPort.Text)
	}}
	s.emailULabel = &widget.Label{Text: "Email User", TextStyle: fyne.TextStyle{Bold: true}}
	s.emailUser = &widget.Entry{Text: fyne.CurrentApp().Preferences().StringWithFallback("euser", ""), OnChanged: func(v string) {
		fyne.CurrentApp().Preferences().SetString("euser", s.emailUser.Text)
	}}
	s.emailPLabel = &widget.Label{Text: "Email Password", TextStyle: fyne.TextStyle{Bold: true}}
	p := s.getPassword()
	s.emailPassword = &widget.Entry{Text: p, OnChanged: func(v string) {
		s.setPassword(s.emailPassword.Text)
	}}
	s.emailPassword.Password = true
	s.emailFLabel = &widget.Label{Text: fmt.Sprintf("Email frequency (%d min)", int(fyne.CurrentApp().Preferences().FloatWithFallback("efrequency", 10))),
		TextStyle: fyne.TextStyle{Bold: true}}
	s.emailFrequency = &widget.Slider{Value: fyne.CurrentApp().Preferences().FloatWithFallback("efrequency", 10),
		Min: 1.0, Max: 60.0, Step: 1,
		OnChanged: func(i float64) {
			s.emailFLabel.SetText(fmt.Sprintf("Email frequency (%d min)", int(i)))
		},
		OnChangeEnded: func(i float64) {
			s.emailFrequency.Value = i
			fyne.CurrentApp().Preferences().SetFloat("efrequency", i)
		}}
	s.useEmail = &widget.RadioGroup{Options: []string{"On", "Off"}, Horizontal: true, Required: true, OnChanged: s.onUseEmailChanged}
	s.useEmail.SetSelected(fyne.CurrentApp().Preferences().StringWithFallback("useemail", "Off"))
	s.onUseEmailChanged(fyne.CurrentApp().Preferences().StringWithFallback("useemail", "Off"))

	emailContainer := container.NewGridWithColumns(2,
		&widget.Label{Text: "Use Email", TextStyle: fyne.TextStyle{Bold: true}}, s.useEmail,
		s.emailSLabel, s.emailServer,
		s.emailPortLabel, s.emailPort,
		s.emailULabel, s.emailUser,
		s.emailPLabel, s.emailPassword,
		s.emailFLabel, s.emailFrequency,
	)
	m := ReadLastLineWithSeek(fyne.CurrentApp().Preferences().String("pfsmslog"), 10)
	s.logtext = &widget.Label{Text: m}
	s.btnCheck = &widget.Button{Text: "Check Email", OnTapped: func() {
		e := new(pfemail.Etype)
		err = e.Checkemaillogin()
		if err != nil {
			log.Println("Email check failed.")
		} else {
			log.Println("Email check ok.")
		}
		m := ReadLastLineWithSeek(fyne.CurrentApp().Preferences().String("pfsmslog"), 10)
		s.logtext.SetText(m)
		s.logtext.Refresh()
	}}
	s.btnStart = &widget.Button{Text: "Start Email", OnTapped: func() {
		e := new(pfemail.Etype)
		log.Println("Handling mail...")
		err = e.Login()
		if err != nil {
			log.Println("Login failed!")
		}
		m0 := e.Getallsmsmail()
		if m0 != nil {
			log.Println("Got Email!!!")
			// m := "SUBJECT:" + m0[0].Envelope.Subject
			// m += "SENDER: "+ m0[0].Envelope.Sender.Address
			// fmt.Println("\r\nFlags: ", m0[0].Flags)		} else {
			log.Println(m + "\r\n")
		} else {
			log.Println("No mail...")
		}
		err = e.Moveallsmsmail()
		if err != nil {
			log.Println("Move SMS Mail failed.")
		} else {
			log.Println("Moved SMS mail to sms folder.")
		}
		m := ReadLastLineWithSeek(fyne.CurrentApp().Preferences().String("pfsmslog"), 10)
		s.logtext.SetText(m)
		s.logtext.Refresh()
	}}

	return container.NewScroll(
		container.NewVBox(
			&widget.Card{Title: "Email Settings", Content: emailContainer},
			container.NewHBox(s.btnCheck, s.btnStart),
			s.logtext,
		))
}

func (s *theemail) onUseEmailChanged(selected string) {
	//	s.client.OverwriteExisting = selected == "On"
	fyne.CurrentApp().Preferences().SetString("useemail", selected)
	s.emailServer.Hidden = (fyne.CurrentApp().Preferences().StringWithFallback("useemail", "Off") == "Off")
	s.emailSLabel.Hidden = (fyne.CurrentApp().Preferences().StringWithFallback("useemail", "Off") == "Off")
	s.emailPort.Hidden = (fyne.CurrentApp().Preferences().StringWithFallback("useemail", "Off") == "Off")
	s.emailPortLabel.Hidden = (fyne.CurrentApp().Preferences().StringWithFallback("useemail", "Off") == "Off")
	s.emailUser.Hidden = (fyne.CurrentApp().Preferences().StringWithFallback("useemail", "Off") == "Off")
	s.emailULabel.Hidden = (fyne.CurrentApp().Preferences().StringWithFallback("useemail", "Off") == "Off")
	s.emailPassword.Hidden = (fyne.CurrentApp().Preferences().StringWithFallback("useemail", "Off") == "Off")
	s.emailPLabel.Hidden = (fyne.CurrentApp().Preferences().StringWithFallback("useemail", "Off") == "Off")
	s.emailFrequency.Hidden = (fyne.CurrentApp().Preferences().StringWithFallback("useemail", "Off") == "Off")
	s.emailFLabel.Hidden = (fyne.CurrentApp().Preferences().StringWithFallback("useemail", "Off") == "Off")
}
func (s *theemail) getPassword() string {
	prefPassword := fyne.CurrentApp().Preferences().StringWithFallback("epassword", "")
	var hash, realPassword string
	var err error
	hash, err = pfdatabase.MakeHash()
	if err != nil {
		log.Println("buildLog onUseEmailChanged MakeHash error ", err.Error())
	}
	realPassword, err = pfdatabase.DecryptPassword(prefPassword, hash)
	if err != nil {
		log.Println("getPassword onEmailChanged DecryptPassword error")
	}
	return realPassword
}
func (s *theemail) setPassword(realPassword string) error {
	var hash, prefPassword string
	var err error
	hash, err = pfdatabase.MakeHash()
	if err != nil {
		log.Println("buildLog setPassord MakeHash error ", err.Error())
	}
	prefPassword, err = pfdatabase.EncryptPassword(realPassword, hash)
	if err != nil {
		log.Println("setPassWord EncryptPassword error")
	}
	fyne.CurrentApp().Preferences().SetString("epassword", prefPassword)
	return err
}
func (s *theemail) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Email", Icon: theme.ComputerIcon(), Content: s.buildUI()}
}

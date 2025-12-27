package ui

// Special program to send a lots of sms using a mobile phone!
// uses logfilename for logging
// uses phonenumbersfilename to specify file with phonenumbers
// 2024-01-21 working!!!!
// 2024-03-10 switched to newer serial driver, implemented support for S24U and model selection
// got it working with Samsung S24Ultra! speed 14s/sms using timeout = Millisecond*700

import (
	"fmt"
	"log"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/prifre/pfsms/pfdatabase"
	"github.com/prifre/pfsms/pfmobile"
	"go.bug.st/serial"
)
const loglines = 10
type theform struct {
	form           *widget.Form
	phone          *widget.Entry
	groupname      *widget.Entry
	message        *widget.Entry
	btnSaveGroup   *widget.Button
	btnDeleteGroup *widget.Button
	btnSubmit      *widget.Button
	groupSelect    *widget.Select
	dataAllGroups  [][]string
	logtext        *widget.Label
	window         fyne.Window
	mydebug   bool
	Comport   string
	timeout   time.Duration
	starttime time.Time
	Addhash   bool
}
func NewMessages(w fyne.Window) *theform {
	return &theform{window: w}
}
func (s *theform) buildUI() *container.Scroll {
	s.groupname = &widget.Entry{}
	s.groupname.Text = fyne.CurrentApp().Preferences().StringWithFallback("groupname","")
	s.groupname.OnChanged = func(v string) {
		fyne.CurrentApp().Preferences().SetString("groupname",s.groupname.Text)
	}
	s.phone = &widget.Entry{}
	s.phone = widget.NewMultiLineEntry()
	s.phone.Wrapping = fyne.TextWrap(fyne.TextWrapWord)
	s.phone.SetMinRowsVisible(6)
	s.phone.Text = fyne.CurrentApp().Preferences().String("messagephone")
	s.phone.OnChanged = func(v string) {
		 fyne.CurrentApp().Preferences().SetString("messagephone",s.phone.Text)
	}
	s.dataAllGroups = new(pfdatabase.DBtype).ShowAllGroups()
	s.groupSelect = &widget.Select{Options: new(pfdatabase.DBtype).ShowGroups(),
		Selected: "",
		OnChanged: func(v string) {
			s.phone.Text = s.Getphonesforgroup(v)
			fyne.CurrentApp().Preferences().SetString("messagephone",s.phone.Text)
			s.phone.Refresh()
			s.groupname.Text = v
			fyne.CurrentApp().Preferences().SetString("groupname",s.groupname.Text)
			s.groupname.Refresh()
		},
	}
	s.btnSaveGroup = &widget.Button{Text: "Save Group", OnTapped: func() {
		// whe4n pasting, ensure so onlyu "0123456789+, " survives...
		if len(s.phone.Text) < 5 {
			return
		}
		z := s.phone.Text
		z = strings.Replace(z, "\r", ", ", -1)
		z2 := ""
		for i := 0; i < len(z); i++ {
			if strings.Contains("0123456789+, ", string(z[i])) {
				z2 += string(z[i])
			}
		}
		z = z2
		for strings.Contains(" ,", string(z[len(z)-1])) {
			z = z[:len(z)-2]
		}
		for strings.Contains(" ,", string(z[0])) {
			z = z[1:]
		}
		s.phone.Text = z
		if strings.Contains(s.phone.Text, ",") && len(s.groupname.Text) > 1 {
			new(pfdatabase.DBtype).SaveGroup(s.groupname.Text, s.phone.Text)
		}
		s.dataAllGroups = new(pfdatabase.DBtype).ShowAllGroups()
		s.groupSelect.Options = new(pfdatabase.DBtype).ShowGroups()
		s.groupSelect.SetSelected(s.groupname.Text)
		s.groupSelect.Refresh()
		s.phone.Text = s.Getphonesforgroup(s.groupSelect.Selected)
		s.phone.Refresh()
	}}
	s.btnDeleteGroup = &widget.Button{Text: "Delete Group", OnTapped: func() {
		new(pfdatabase.DBtype).DeleteGroup(s.groupSelect.Selected)
		s.dataAllGroups = new(pfdatabase.DBtype).ShowAllGroups()
		s.groupSelect.Options = new(pfdatabase.DBtype).ShowGroups()
		s.groupSelect.Selected = ""
		s.groupSelect.Refresh()
		s.groupname.SetText("")
		s.groupname.Refresh()
	}}
	GroupsInfo := "To use multiple mobile numbers, separate them with commas or Enter.\r\n"
	GroupsInfo += "Click Save Group to reuse in future."
	MessageInfo := "To insert firstname and/or lastname, use <<Fname>> and <<Lname>> in message."
	s.message = &widget.Entry{}
	s.message = widget.NewMultiLineEntry()
	s.message.Wrapping = fyne.TextWrap(fyne.TextWrapWord)
	s.message.SetMinRowsVisible(8)
	s.message.Text = fyne.CurrentApp().Preferences().String("message")
	s.message.OnChanged = func (v string) {
		fyne.CurrentApp().Preferences().SetString("message",s.message.Text)
	}
	s.btnSubmit = &widget.Button{Text: "Click to send message", OnTapped: func() {
		s.HandleSendsms(s.phone.Text, s.groupname.Text, s.message.Text)
	}}
	s.form = &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "", Widget: NewBoldLabel(GroupsInfo)},
			{Text: "Groups", Widget: container.NewGridWithColumns(3, s.groupSelect, s.btnSaveGroup, s.btnDeleteGroup)},
			{Text: "Groupname", Widget: s.groupname},
			{Text: "Phone", Widget: s.phone},
			{Text: "", Widget: NewBoldLabel(MessageInfo)},
			{Text: "Message", Widget: s.message},
			{Text: "", Widget: s.btnSubmit},
		},
	}
	s.logtext = &widget.Label{Text: ShowShortLines(ReadLastLineWithSeek(fyne.CurrentApp().Preferences().String("pfsmslog"),loglines))}
	s.form.Refresh()
	return container.NewScroll(
		container.NewVBox(
			s.form,
			s.logtext,
		))
}
func (s *theform) Getphonesforgroup(v string) string {
	var np string
	dag := new(pfdatabase.DBtype).ShowAllGroups()
	for i := 0; i < len(dag); i++ {
		if dag[i][0] == v {
			if np > "" {
				np = np + ", " + dag[i][1]
			} else {
				np = dag[i][1]
			}
		}
	}
	return np
}
func  (s *theform) HandleSendsms(phone,groupname, msg string) {
	// split phone into \r\n and ","
	ph := phone
	ph = strings.Replace(ph, "\r", ",", -1)
	ph = strings.Replace(ph, "\n", ",", -1)
	ph = strings.Replace(ph, ",,", ",", -1)
	ph = strings.Replace(ph, ",,", ",", -1)
	ph2 := ""
	for i := 0; i < len(ph); i++ {
		if strings.Contains("0123456789+,", string(ph[i])) {
			ph2 = ph2 + string(ph[i])
		}
	}
	countrycode := fyne.CurrentApp().Preferences().StringWithFallback("mobilecountry", "Sweden(+46)")
	p2:=""
	for _,p :=range(strings.Split(ph2,",")) {
		if len(p)<5 {
			continue
		}
		if p2>"" {
			p2+=","
		}
		p2 += Fixphonenumber(p, countrycode)
		p2 += "\t" + new(pfdatabase.DBtype).GetFname(p) + "\t" + new(pfdatabase.DBtype).GetLname(p)
	}
	s.SendMessages(strings.Split(p2,","), msg)
}
func (s *theform) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Messages", Icon: theme.MailSendIcon(), Content: s.buildUI()}
}

func (s *theform) SendMessages(phonenumbers []string, message string) error {
	// Replace with the correct serial port of the modem
	if s.Comport=="" {
		s.Comport = fyne.CurrentApp().Preferences().StringWithFallback("mobileport", "COM2")
		s.Addhash = fyne.CurrentApp().Preferences().Bool("addhash")
	}
	var sendtext, phoneNumber string
	var failures, success int
	var result error
	s.mydebug = true
	// s.Setuplog()
	s.starttime = time.Now()
	s.timeout = time.Millisecond * 700
	message = strings.TrimSpace(message)
	log.Printf("Got %d phonenumbers to send ok.\r\n", len(phonenumbers))
	modemresetfail := 0
	for !pfmobile.Modemreset(s.Comport) && modemresetfail < 10 {
		log.Println("--------------------MODEMRESET FAIL: ", modemresetfail)
		modemresetfail++
		s.logtext = &widget.Label{Text: ShowShortLines(ReadLastLineWithSeek(fyne.CurrentApp().Preferences().String("pfsmslog"),loglines))}
	}
	for i, record := range phonenumbers {
		rec := strings.Split(record, "\t")
		phoneNumber = rec[0]
		sendtext = message
		if strings.Contains(sendtext, "<<Fname>>") || strings.Contains(sendtext, "<<Lname>>") {
			sendtext = strings.Replace(sendtext, "<<Fname>>", rec[1], -1)
			sendtext = strings.Replace(sendtext, "<<Lname>>", rec[2], -1)
		}
		if s.Addhash {
			sendtext = fmt.Sprintf(sendtext+"\r\n#=%d", i+1)
		}
		sentok := false
		for !sentok {
			sentok = pfmobile.SendSMS(s.Comport,phoneNumber, sendtext)
			if !sentok {
				log.Println("--------------------SENDSMS FAILED")
				modemresetfail := 0
				for !pfmobile.Modemreset(s.Comport) && modemresetfail < 10 {
					log.Println("--------------------MODEMRESET FAIL: ", modemresetfail)
					modemresetfail++
					s.logtext = &widget.Label{Text: ShowShortLines(ReadLastLineWithSeek(fyne.CurrentApp().Preferences().String("pfsmslog"),loglines))}
				}
				if modemresetfail > 8 {
					return nil
				}
				log.Println("--------------------MODEMRESET OK")
				failures++
			}
			s.logtext = &widget.Label{Text: ShowShortLines(ReadLastLineWithSeek(fyne.CurrentApp().Preferences().String("pfsmslog"),loglines))}
		}
		success++
		m:=fmt.Sprintf("Message %d/%d to phone %s sent!", i+1, len(phonenumbers), phoneNumber)
		if failures>0 {
			m += fmt.Sprintf(" (failures: %d)", failures)
		}
		log.Println(m)
		tstamp := time.Now().Format("20060102150405")
		//SaveHistory([]string {tstamp,groupname,phone,message})
		new(pfdatabase.DBtype).SaveHistory([]string{tstamp, s.groupname.Text, phoneNumber, sendtext})
		s.logtext.Text = ReadLastLineWithSeek(fyne.CurrentApp().Preferences().String("pfsmslog"), 8)
		s.logtext.Refresh()
		if !s.mydebug {
			log.Printf("%s Message %d/%d to phone %s sent! (failures: %d)\r\n", time.Now().Format("2006-01-02 15:04:05"), i+1, len(phonenumbers), phoneNumber, failures)
		}
	}
	if !s.mydebug {
		log.Printf("RESULT OF SMS SENDING: Failures: %d Success: %d\r\n", failures, success)
		s1 := s.starttime.Format("2006-01-02 15:04:05")
		s2 := time.Now().Format("2006-01-02 15:04:05")
		log.Printf("Started: %s  Finished: %s  Duration: %s\r\n", s1, s2, time.Since(s.starttime))
		log.Printf("Speed: %ds/sms\r\n", int(time.Since(s.starttime).Seconds())/len(phonenumbers))
	}
	return result
}

func GetPortsList() ([]string, error) {
	return serial.GetPortsList()
}
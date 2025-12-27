package ui

// Special program to send a lots of sms using a mobile phone!
// uses logfilename for logging
// uses phonenumbersfilename to specify file with phonenumbers
// 2024-01-21 working!!!!
// 2024-03-10 switched to newer serial driver, implemented support for S24U and model selection
// got it working with Samsung S24Ultra! speed 14s/sms using timeout = Millisecond*700

import (
	"errors"
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
	var err error
	s.mydebug = true
	// s.Setuplog()
	s.starttime = time.Now()
	s.timeout = time.Millisecond * 700
	message = strings.TrimSpace(message)
	log.Printf("Got %d phonenumbers to send ok.\r\n", len(phonenumbers))
	modemresetfail := 0
	for (pfmobile.Modemreset("")!=nil) && modemresetfail < 10 {
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
		err = pfmobile.SendSMS(nil,phoneNumber, sendtext)
		if err!=nil {
			log.Println("--------------------SENDSMS FAILED")
			return errors.New("SendSMS failed: " + err.Error())
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
	return nil
}

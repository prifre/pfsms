package ui

import (
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/prifre/pfsms/ariasms"
)

type theform struct {
	form			*widget.Form
	phone	 		*widget.Entry
	reference 		*widget.Entry
	message	 		*widget.Entry
	btnSubmit 		*widget.Button
	logtext			*widget.Label
	window      	fyne.Window
	app         	fyne.App
}

func NewMessages(a fyne.App, w fyne.Window) *theform {
	return &theform{app: a, window: w}
}
func (s *theform) buildForm() *container.Scroll {
	var err error
	s.reference = &widget.Entry{}
	s.reference.Text="test1"
	s.phone = &widget.Entry{}
	s.phone=widget.NewMultiLineEntry()
	s.phone.Wrapping=fyne.TextWrap(fyne.TextWrapWord)
	s.phone.SetMinRowsVisible(6)
	s.phone.Text="0046736290839"
	s.message = &widget.Entry{}
	s.message=widget.NewMultiLineEntry()
	s.message.Wrapping=fyne.TextWrap(fyne.TextWrapWord)
	s.message.SetMinRowsVisible(8)
	s.message.Text="ett litet test!"
	s.btnSubmit = & widget.Button{Text:"Click to send message",OnTapped: func ()  {
		s.HandleSendsms(s.phone.Text,s.reference.Text,s.message.Text)
	}}
	s.form = &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Reference", Widget: s.reference},
			{Text: "Phone", Widget: s.phone},
			{Text: "Messagetext", Widget: s.message},
			{Text:"",Widget:s.btnSubmit},
		},}
		// sms:=new(ariasms.SMStype)
		// sms.SendMessage([]string{s.phone.Text},s.message.Text)		
	info:="To use multiple mobile numbers, separate them with commas or Enter.\r\n"
	info +="To insert firstname and/or lastname, use <<Fname>> and <<Lname>> in message"
	s.logtext=&widget.Label{Text: " "}
	var txt string
	txt,err = ReadLastLineWithSeek("smslog.txt",10)
	if err!=nil {
		log.Println("#1 buildLog!", err.Error())
	}
	s.logtext.Text=txt
	return container.NewScroll(		
		container.NewVBox(
			widget.NewLabel(info),
			s.form,
			s.logtext,
	))
}
func (s *theform) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Messages", Icon: theme.MailSendIcon(), Content: s.buildForm()}
}
func (s *theform) HandleSendsms(p,t,m string) {
	// split phone into \r\n and ","
	ph := s.phone.Text
	ph = strings.Replace(ph,"\r",",",-1)
	ph = strings.Replace(ph,"\n",",",-1)
	ph = strings.Replace(ph,",,",",",-1)
	ph = strings.Replace(ph,",,",",",-1)
	ph2:=""
	for i:=0;i<len(ph);i++ {
		if strings.Contains("0123456789+,",string(ph[i])) {
			ph2= ph2+string(ph[i])
		}
	}
	p1:=strings.Split(ph2,",")
	countrycode := s.app.Preferences().StringWithFallback("mobileCountry", "Sweden(+46)")
	for i:=0;i<len(p1);i++ {
		p1[i]=Fixphonenumber(p1[i],countrycode)
	}
	var sms ariasms.SMStype =*new(ariasms.SMStype)
	sms.Comport = s.app.Preferences().StringWithFallback("mobilePort", "COM2")
	sms.Addhash=s.app.Preferences().Bool("addHash")
	err := sms.SendMessage(p1,m)
	if err!=nil {
		Appendtotextfile("smslog.txt","#1 HandleSendsms Failed to send messages "+err.Error())
	} else {
		for i:=0;i<len(p1);i++ {
			Appendtotextfile("smslog.txt"," Sent "+s.message.Text+" to "+p1[i]+"\r\n")
		}
	}
	var txt string
	txt,err = ReadLastLineWithSeek("smslog.txt",10)
	if err!=nil {
		log.Println("#1 buildLog!", err.Error())
	}
	s.logtext.Text=txt
	s.logtext.Refresh()
}
package ui

import (
	"fmt"
	"strings"
	"time"

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
	window      	fyne.Window
	app         	fyne.App
}

func NewMessages(a fyne.App, w fyne.Window) *theform {
	return &theform{app: a, window: w}
}
func (s *theform) buildForm() *container.Scroll {
	s.reference = &widget.Entry{}
	s.reference.Text="test1"
	s.phone = &widget.Entry{}
	s.phone=widget.NewMultiLineEntry()
	s.phone.Wrapping=fyne.TextWrap(fyne.TextWrapWord)
	s.phone.SetMinRowsVisible(8)
	s.phone.Text="0046736290839"
	s.message = &widget.Entry{}
	s.message=widget.NewMultiLineEntry()
	s.message.Wrapping=fyne.TextWrap(fyne.TextWrapWord)
	s.message.SetMinRowsVisible(12)
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
	return container.NewScroll(		
		container.NewVBox(
			widget.NewLabel(info),
			s.form,
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
	p1:=strings.Split(ph,",")
	countrycode := s.app.Preferences().StringWithFallback("mobileCountry", "Sweden(+46)")
	for i:=0;i<len(p1);i++ {
		p1[i]=Fixphonenumber(p1[i],countrycode)
	}
	var sms ariasms.SMStype =*new(ariasms.SMStype)
	sms.Comport = s.app.Preferences().StringWithFallback("mobilePort", "COM2")
	sms.Addhash=s.app.Preferences().Bool("addHash")
	sms.SendMessage(p1,m)
	for i:=0;i<len(p1);i++ {
		Appendtotextfile("smslog.txt",fmt.Sprintf("\r\n%s %s",time.Now().Format("2006-01-02 15:04:05")," Sent "+s.message.Text+" to "+p1[i]))
	}
}
